package intelligence

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// Cadence represents subscription frequency.
type Cadence string

const (
	CadenceWeekly  Cadence = "WEEKLY"
	CadenceMonthly Cadence = "MONTHLY"
	CadenceAnnual  Cadence = "ANNUAL"
)

// Subscription represents a detected recurring charge.
type Subscription struct {
	ID              string       `json:"id"`
	MerchantName    string       `json:"merchantName"`
	CategoryName    string       `json:"categoryName"`
	Cadence         Cadence      `json:"cadence"`
	Amount          domain.Money `json:"amount"`
	LastChargedAt   time.Time    `json:"lastChargedAt"`
	NextExpectedAt  time.Time    `json:"nextExpectedAt"`
	OccurrenceCount int          `json:"occurrenceCount"`
	AnnualizedCost  domain.Money `json:"annualizedCost"`
	PriceDriftNote  string       `json:"priceDriftNote,omitempty"`
	Confidence      float64      `json:"confidence"`
}

// TransactionRecord is a minimal input for recurring analysis.
type TransactionRecord struct {
	Description string
	Category    string
	AmountMinor int64
	Currency    string
	Date        time.Time
}

// SubscriptionDetector analyzes transaction history to find recurring subscriptions.
type SubscriptionDetector struct{}

// NewSubscriptionDetector creates a new detector instance.
func NewSubscriptionDetector() *SubscriptionDetector {
	return &SubscriptionDetector{}
}

// Detect analyzes historical transactions and returns detected subscriptions.
func (d *SubscriptionDetector) Detect(txs []TransactionRecord) []Subscription {
	// Group transactions by normalized merchant name
	groups := make(map[string][]TransactionRecord)
	for _, tx := range txs {
		if tx.AmountMinor <= 0 {
			continue // only analyze expenses
		}
		norm := normalizeMerchant(tx.Description)
		groups[norm] = append(groups[norm], tx)
	}

	var results []Subscription

	for norm, items := range groups {
		if len(items) < 2 {
			continue // Need at least 2 occurrences to establish recurrence
		}

		// Sort ascending by date
		sort.Slice(items, func(i, j int) bool {
			return items[i].Date.Before(items[j].Date)
		})

		// Calculate intervals in days
		var intervals []float64
		var totalAmount int64
		for i := 1; i < len(items); i++ {
			days := items[i].Date.Sub(items[i-1].Date).Hours() / 24.0
			intervals = append(intervals, days)
		}
		for _, it := range items {
			totalAmount += it.AmountMinor
		}

		avgInterval := average(intervals)
		stdInterval := stdDev(intervals, avgInterval)

		// Determine cadence
		var cadence Cadence
		var expectedDays int
		var confidence float64 = 0.85

		if avgInterval >= 25 && avgInterval <= 35 && stdInterval <= 6.0 {
			cadence = CadenceMonthly
			expectedDays = 30
			confidence = 0.95
		} else if avgInterval >= 6 && avgInterval <= 8 && stdInterval <= 2.5 {
			cadence = CadenceWeekly
			expectedDays = 7
			confidence = 0.90
		} else if avgInterval >= 350 && avgInterval <= 380 {
			cadence = CadenceAnnual
			expectedDays = 365
			confidence = 0.88
		} else {
			continue // Not a regular pattern
		}

		lastTx := items[len(items)-1]
		firstTx := items[0]
		lastAmount := domain.Money{AmountMinor: lastTx.AmountMinor, Currency: lastTx.Currency}
		nextExpected := lastTx.Date.AddDate(0, 0, expectedDays)

		// Check for price drift
		driftNote := ""
		if lastTx.AmountMinor > firstTx.AmountMinor {
			diff := domain.Money{AmountMinor: lastTx.AmountMinor - firstTx.AmountMinor, Currency: lastTx.Currency}
			driftNote = "Mensalidade subiu " + diff.FormatBRL() + " em relação ao início"
		}

		// Annualized cost calculation
		var annualizedMinor int64
		switch cadence {
		case CadenceWeekly:
			annualizedMinor = lastTx.AmountMinor * 52
		case CadenceMonthly:
			annualizedMinor = lastTx.AmountMinor * 12
		case CadenceAnnual:
			annualizedMinor = lastTx.AmountMinor
		}

		results = append(results, Subscription{
			ID:              "sub-" + norm,
			MerchantName:    capitalizeTitle(norm),
			CategoryName:    lastTx.Category,
			Cadence:         cadence,
			Amount:          lastAmount,
			LastChargedAt:   lastTx.Date,
			NextExpectedAt:  nextExpected,
			OccurrenceCount: len(items),
			AnnualizedCost:  domain.Money{AmountMinor: annualizedMinor, Currency: lastTx.Currency},
			PriceDriftNote:  driftNote,
			Confidence:      confidence,
		})
	}

	// Sort results by annualized cost descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].AnnualizedCost.AmountMinor > results[j].AnnualizedCost.AmountMinor
	})

	return results
}

func normalizeMerchant(desc string) string {
	s := strings.ToLower(strings.TrimSpace(desc))
	for _, noise := range []string{"*br", "*brasil", "pag*", "dl*", "mensalidade", "assinatura"} {
		s = strings.ReplaceAll(s, noise, "")
	}
	s = strings.TrimSpace(s)
	if strings.Contains(s, "netflix") {
		return "netflix"
	}
	if strings.Contains(s, "spotify") {
		return "spotify"
	}
	if strings.Contains(s, "smart fit") || strings.Contains(s, "smartfit") {
		return "smart fit"
	}
	if strings.Contains(s, "amazon prime") || strings.Contains(s, "prime video") {
		return "amazon prime"
	}
	if strings.Contains(s, "claro") {
		return "claro internet"
	}
	if strings.Contains(s, "vivo") {
		return "vivo fibra"
	}
	if strings.Contains(s, "apple") || strings.Contains(s, "icloud") {
		return "apple services / icloud"
	}
	return s
}

func average(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stdDev(vals []float64, avg float64) float64 {
	if len(vals) <= 1 {
		return 0
	}
	var sumSquares float64
	for _, v := range vals {
		diff := v - avg
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(vals)))
}

func capitalizeTitle(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
