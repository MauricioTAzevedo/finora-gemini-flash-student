package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/google/uuid"
)

// BenchmarkLedgerPosting10000 tests high-throughput double-entry validation over 10,000 transactions.
func BenchmarkLedgerPosting10000(b *testing.B) {
	householdID := uuid.New()
	checkingID := uuid.New()
	categoryID := uuid.New()
	now := time.Now().UTC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var totalDebits int64
		var totalCredits int64

		for j := 0; j < 10000; j++ {
			amountMinor := int64(1000 + (j % 5000))
			amount := domain.NewBRL(amountMinor)

			tx, err := domain.NewExpense(
				householdID,
				checkingID,
				categoryID,
				amount,
				fmt.Sprintf("Supermercado Compra #%d", j),
				now,
				domain.SourceManual,
			)
			if err != nil {
				b.Fatalf("failed to create expense transaction: %v", err)
			}

			if err := tx.Validate(); err != nil {
				b.Fatalf("transaction %d failed double-entry validation: %v", j, err)
			}

			totalDebits += tx.TotalDebits()
			totalCredits += tx.TotalCredits()
		}

		if totalDebits != totalCredits {
			b.Fatalf("invariant violated: totalDebits %d != totalCredits %d", totalDebits, totalCredits)
		}
	}
}
