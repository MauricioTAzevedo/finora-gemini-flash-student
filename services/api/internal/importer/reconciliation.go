package importer

import (
	"fmt"
	"math"
	"strings"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
)

// ReconcileTransaction evaluates an ingested record against existing ledger transactions to detect duplicates or candidate matches.
func ReconcileTransaction(record CanonicalIngestionRecord, existing []domain.LedgerTransaction) ReconciliationResult {
	bestScore := 0.0
	var bestMatch *domain.LedgerTransaction
	bestReason := "No matching transaction found in ledger"

	recordAmountAbs := int64(math.Abs(float64(record.Amount.AmountMinor)))

	for _, tx := range existing {
		// Calculate amount score (Max 50 points)
		txAmountAbs := int64(math.Abs(float64(tx.TotalDebits())))
		if txAmountAbs == 0 {
			txAmountAbs = int64(math.Abs(float64(tx.TotalCredits())))
		}

		if txAmountAbs != recordAmountAbs {
			continue // In financial reconciliation, different amounts rarely match
		}

		amountScore := 0.50

		// Date proximity score (Max 30 points)
		dayDiff := int(math.Abs(record.OccurredAt.Sub(tx.OccurredAt).Hours() / 24))
		dateScore := 0.0
		switch dayDiff {
		case 0:
			dateScore = 0.30
		case 1:
			dateScore = 0.25
		case 2:
			dateScore = 0.18
		case 3:
			dateScore = 0.10
		default:
			dateScore = 0.0
		}

		// Text / Merchant similarity score (Max 20 points)
		textScore := computeTextSimilarity(record.Description, tx.Description) * 0.20

		totalScore := amountScore + dateScore + textScore

		if totalScore > bestScore {
			bestScore = totalScore
			bestMatch = &tx
			bestReason = fmt.Sprintf(
				"Amount match (R$ %.2f) + Date distance (%d days) + Description similarity (%.0f%%)",
				float64(recordAmountAbs)/100.0,
				dayDiff,
				textScore*500,
			)
		}
	}

	if bestMatch != nil {
		if bestScore >= 0.85 {
			return ReconciliationResult{
				Record:               record,
				Status:               MatchStatusDuplicate,
				Score:                bestScore,
				MatchedTransactionID: bestMatch.ID.String(),
				Reason:               fmt.Sprintf("Definite duplicate of existing transaction #%s: %s", bestMatch.ID.String()[:8], bestReason),
			}
		} else if bestScore >= 0.65 {
			return ReconciliationResult{
				Record:               record,
				Status:               MatchStatusCandidate,
				Score:                bestScore,
				MatchedTransactionID: bestMatch.ID.String(),
				Reason:               fmt.Sprintf("Probable candidate match with transaction #%s: %s", bestMatch.ID.String()[:8], bestReason),
			}
		}
	}

	return ReconciliationResult{
		Record: record,
		Status: MatchStatusNew,
		Score:  bestScore,
		Reason: "Inédita: nenhuma transação equivalente encontrada no razão",
	}
}

func computeTextSimilarity(s1, s2 string) float64 {
	w1 := tokenize(s1)
	w2 := tokenize(s2)

	if len(w1) == 0 || len(w2) == 0 {
		return 0.0
	}

	matches := 0.0
	for _, a := range w1 {
		for _, b := range w2 {
			if a == b || strings.Contains(a, b) || strings.Contains(b, a) || isMerchantAlias(a, b) {
				matches += 1.0
				break
			}
		}
	}

	return matches / math.Max(float64(len(w1)), float64(len(w2)))
}

func isMerchantAlias(a, b string) bool {
	aliases := map[string]string{
		"amzn": "amazon",
		"mktp": "marketplace",
		"posto": "combustivel",
		"cpfl": "energia",
	}

	if val, ok := aliases[a]; ok && (val == b || strings.Contains(b, val)) {
		return true
	}
	if val, ok := aliases[b]; ok && (val == a || strings.Contains(a, val)) {
		return true
	}
	if len(a) >= 3 && len(b) >= 3 && (strings.HasPrefix(a, b[:3]) || strings.HasPrefix(b, a[:3])) {
		return true
	}
	return false
}

func tokenize(s string) []string {
	cleaned := strings.ToLower(s)
	// Remove common punctuation
	replacer := strings.NewReplacer("*", " ", "-", " ", ".", " ", "/", " ", ",", " ")
	cleaned = replacer.Replace(cleaned)
	words := strings.Fields(cleaned)

	var filtered []string
	for _, w := range words {
		// Ignore short stopwords
		if len(w) > 2 && w != "ltda" && w != "sao" && w != "paulo" && w != "brasil" {
			filtered = append(filtered, w)
		}
	}
	return filtered
}
