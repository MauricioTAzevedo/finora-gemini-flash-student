package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// AnomalySeverity indicates urgency.
type AnomalySeverity string

const (
	SeverityInfo     AnomalySeverity = "INFO"
	SeverityWarning  AnomalySeverity = "WARNING"
	SeverityCritical AnomalySeverity = "CRITICAL"
)

// AnomalyType classifies the detected anomaly.
type AnomalyType string

const (
	AnomalySpikeUtility   AnomalyType = "SPIKE_UTILITY"
	AnomalyCategorySpike  AnomalyType = "SPIKE_CATEGORY"
	AnomalyPotentialDupe  AnomalyType = "POTENTIAL_DUPLICATE"
	AnomalyUnusualMerchant AnomalyType = "UNUSUAL_MERCHANT"
)

// Anomaly represents an explainable financial anomaly.
type Anomaly struct {
	ID             string          `json:"id"`
	Type           AnomalyType     `json:"type"`
	Severity       AnomalySeverity `json:"severity"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Amount         domain.Money    `json:"amount"`
	HistoricalMean domain.Money    `json:"historicalMean"`
	ZScore         float64         `json:"zScore"`
	PercentageDiff float64         `json:"percentageDiff"`
	TransactionID  string          `json:"transactionId,omitempty"`
	CategoryName   string          `json:"categoryName"`
	MerchantName   string          `json:"merchantName"`
	Date           time.Time       `json:"date"`
	ActionAdvice   string          `json:"actionAdvice"`
}

// AnomalyDetector analyzes transactions for statistical spikes and duplicates.
type AnomalyDetector struct {
	zScoreThreshold float64
}

// NewAnomalyDetector creates a new detector with standard 1.8 sigma threshold for small samples.
func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{zScoreThreshold: 1.8}
}

// Detect evaluates transactions for statistical anomalies.
func (ad *AnomalyDetector) Detect(txs []TransactionRecord) []Anomaly {
	var anomalies []Anomaly

	// Group transactions by category to compute historical baselines
	categoryMap := make(map[string][]int64)
	for _, tx := range txs {
		if tx.AmountMinor > 0 {
			categoryMap[tx.Category] = append(categoryMap[tx.Category], tx.AmountMinor)
		}
	}

	// 1. Check for duplicate charges on the same day with the same amount
	dupeKeyMap := make(map[string][]TransactionRecord)
	for _, tx := range txs {
		if tx.AmountMinor > 0 {
			key := fmt.Sprintf("%s|%s|%d", tx.Date.Format("2006-01-02"), strings.ToLower(tx.Description), tx.AmountMinor)
			dupeKeyMap[key] = append(dupeKeyMap[key], tx)
		}
	}

	for _, list := range dupeKeyMap {
		if len(list) > 1 {
			tx := list[0]
			m := domain.Money{AmountMinor: tx.AmountMinor, Currency: tx.Currency}
			anomalies = append(anomalies, Anomaly{
				ID:           fmt.Sprintf("anom-dupe-%d-%s", tx.AmountMinor, tx.Date.Format("20060102")),
				Type:         AnomalyPotentialDupe,
				Severity:     SeverityWarning,
				Title:        "Possível Cobrança Duplicada",
				Description:  fmt.Sprintf("Identificamos %d cobranças idênticas de %s em %s no mesmo dia.", len(list), m.FormatBRL(), tx.Description),
				Amount:       m,
				CategoryName: tx.Category,
				MerchantName: tx.Description,
				Date:         tx.Date,
				ActionAdvice: "Verifique o extrato do seu cartão para confirmar se não houve duplicidade de débito pela maquininha.",
			})
		}
	}

	// 2. Check for statistical category spikes (Z-score >= 2.0)
	for _, tx := range txs {
		if tx.AmountMinor <= 0 {
			continue
		}

		amounts := categoryMap[tx.Category]
		if len(amounts) < 3 {
			continue // Insufficient data for statistical baseline
		}

		mean := meanInt64(amounts)
		sd := stdDevInt64(amounts, mean)

		if sd > 0 {
			z := float64(tx.AmountMinor-mean) / float64(sd)
			if z >= ad.zScoreThreshold {
				pct := float64(tx.AmountMinor-mean) / float64(mean) * 100.0

				sev := SeverityWarning
				if z >= 3.0 || pct > 80.0 {
					sev = SeverityCritical
				}

				anomType := AnomalyCategorySpike
				advice := "Revise os gastos recentes desta categoria para manter o orçamento planejado."

				isUtility := strings.Contains(strings.ToLower(tx.Category), "moradia") ||
					strings.Contains(strings.ToLower(tx.Category), "utilidades") ||
					strings.Contains(strings.ToLower(tx.Description), "cpfl") ||
					strings.Contains(strings.ToLower(tx.Description), "sabesp") ||
					strings.Contains(strings.ToLower(tx.Description), "enel")

				if isUtility {
					anomType = AnomalySpikeUtility
					advice = "Verifique bandeira tarifária ou possível fuga/vazamento de energia/água no imóvel."
				}

				anomalies = append(anomalies, Anomaly{
					ID:             fmt.Sprintf("anom-spike-%d-%s", tx.AmountMinor, tx.Date.Format("20060102")),
					Type:           anomType,
					Severity:       sev,
					Title:          fmt.Sprintf("Desvio Significativo em %s (+%.0f%%)", tx.Category, pct),
					Description:    fmt.Sprintf("Gasto de %s em '%s' está %.0f%% acima da sua média de %s (desvio de %.1fσ).", domain.Money{AmountMinor: tx.AmountMinor, Currency: tx.Currency}.FormatBRL(), tx.Description, pct, domain.Money{AmountMinor: mean, Currency: tx.Currency}.FormatBRL(), z),
					Amount:         domain.Money{AmountMinor: tx.AmountMinor, Currency: tx.Currency},
					HistoricalMean: domain.Money{AmountMinor: mean, Currency: tx.Currency},
					ZScore:         math.Round(z*10) / 10,
					PercentageDiff: math.Round(pct),
					CategoryName:   tx.Category,
					MerchantName:   tx.Description,
					Date:           tx.Date,
					ActionAdvice:   advice,
				})
			}
		}
	}

	// Sort by severity (Critical first), then date descending
	sort.Slice(anomalies, func(i, j int) bool {
		if anomalies[i].Severity != anomalies[j].Severity {
			return anomalies[i].Severity == SeverityCritical
		}
		return anomalies[i].Date.After(anomalies[j].Date)
	})

	return anomalies
}

func meanInt64(vals []int64) int64 {
	if len(vals) == 0 {
		return 0
	}
	var sum int64
	for _, v := range vals {
		sum += v
	}
	return sum / int64(len(vals))
}

func stdDevInt64(vals []int64, mean int64) float64 {
	if len(vals) <= 1 {
		return 0
	}
	var sumSquares float64
	for _, v := range vals {
		diff := float64(v - mean)
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(vals)))
}
