package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/google/uuid"
)

func TestDoubleEntryBalancingInvariant(t *testing.T) {
	householdID := uuid.New()
	checkingID := uuid.New()
	categoryID := uuid.New()
	amount := domain.NewBRL(15000) // R$ 150,00
	now := time.Now()

	// 1. Valid Expense Transaction (Debits == Credits)
	tx, err := domain.NewExpense(householdID, checkingID, categoryID, amount, "Supermercado", now, domain.SourceManual)
	if err != nil {
		t.Fatalf("expected valid expense creation, got error: %v", err)
	}

	if tx.TotalDebits() != tx.TotalCredits() {
		t.Errorf("expected debits (%d) to equal credits (%d)", tx.TotalDebits(), tx.TotalCredits())
	}
	if tx.TotalDebits() != 15000 {
		t.Errorf("expected total debits 15000, got %d", tx.TotalDebits())
	}

	// 2. Artificially Unbalanced Transaction Must Fail
	unbalancedTx := &domain.LedgerTransaction{
		ID:          uuid.New(),
		HouseholdID: householdID,
		OccurredAt:  now,
		Description: "Unbalanced Hack",
		Status:      domain.StatusPosted,
		Entries: []domain.LedgerEntry{
			{
				ID:          uuid.New(),
				EntryType:   domain.EntryTypeDebit,
				AmountMinor: 10000,
			},
			{
				ID:          uuid.New(),
				EntryType:   domain.EntryTypeCredit,
				AmountMinor: 5000, // Differs by 5000 cents!
			},
		},
	}

	err = unbalancedTx.Validate()
	if err == nil {
		t.Fatal("expected error validating unbalanced transaction, got nil")
	}
	if !errors.Is(err, domain.ErrUnbalancedLedger) {
		t.Errorf("expected ErrUnbalancedLedger, got %v", err)
	}
}

func TestTransferIsNotExpense(t *testing.T) {
	// FINANCIAL INVARIANT: Moving R$ 1.000 from Checking to Savings must NOT produce R$ 1.000 of household spending.
	householdID := uuid.New()
	checkingID := uuid.New()
	savingsID := uuid.New()
	amount := domain.NewBRL(100000) // R$ 1.000,00

	tx, err := domain.NewTransfer(householdID, checkingID, savingsID, amount, "Reserva de Emergência", time.Now())
	if err != nil {
		t.Fatalf("unexpected error creating transfer: %v", err)
	}

	// Double-entry validation
	if err := tx.Validate(); err != nil {
		t.Fatalf("transfer transaction should be balanced: %v", err)
	}

	// The net expense MUST be 0
	if tx.NetExpense() != 0 {
		t.Errorf("CRITICAL VIOLATION: Transfer produced NetExpense = %d; MUST BE 0", tx.NetExpense())
	}

	// Net income MUST also be 0
	if tx.NetIncome() != 0 {
		t.Errorf("CRITICAL VIOLATION: Transfer produced NetIncome = %d; MUST BE 0", tx.NetIncome())
	}
}

func TestCreditCardPaymentIsNotExpenseTwice(t *testing.T) {
	// FINANCIAL INVARIANT: A card purchase of R$ 200 and later paying the bill of R$ 200 must not total R$ 400 in spending.
	householdID := uuid.New()
	checkingID := uuid.New()
	cardID := uuid.New()
	categoryID := uuid.New()

	// 1. Purchase on Credit Card (R$ 200,00)
	purchaseTx, err := domain.NewExpense(householdID, cardID, categoryID, domain.NewBRL(20000), "Jantar Restaurante", time.Now(), domain.SourceManual)
	if err != nil {
		t.Fatalf("failed to create card expense: %v", err)
	}

	if purchaseTx.NetExpense() != 20000 {
		t.Errorf("expected purchase NetExpense = 20000, got %d", purchaseTx.NetExpense())
	}

	// 2. Paying the Credit Card bill from Checking (R$ 200,00)
	paymentTx, err := domain.NewCreditCardPayment(householdID, checkingID, cardID, domain.NewBRL(20000), "Pagamento Fatura Nubank", time.Now())
	if err != nil {
		t.Fatalf("failed to create card bill payment: %v", err)
	}

	if err := paymentTx.Validate(); err != nil {
		t.Fatalf("card payment must be balanced: %v", err)
	}

	// The payment transaction must NOT count as an expense
	if paymentTx.NetExpense() != 0 {
		t.Errorf("CRITICAL VIOLATION: Card bill payment produced NetExpense = %d; MUST BE 0", paymentTx.NetExpense())
	}

	// Total household spending across purchase + bill payment must be exactly R$ 200 (20000 cents), NOT R$ 400
	totalSpending := purchaseTx.NetExpense() + paymentTx.NetExpense()
	if totalSpending != 20000 {
		t.Errorf("CRITICAL VIOLATION: Total household spending was %d; expected exactly 20000", totalSpending)
	}
}

func TestTransactionReversal(t *testing.T) {
	householdID := uuid.New()
	checkingID := uuid.New()
	categoryID := uuid.New()
	amount := domain.NewBRL(5000) // R$ 50,00

	tx, err := domain.NewExpense(householdID, checkingID, categoryID, amount, "Farmácia", time.Now(), domain.SourceManual)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	revTx, err := tx.Reverse("Compra cancelada / estorno")
	if err != nil {
		t.Fatalf("unexpected error on reverse: %v", err)
	}

	if err := revTx.Validate(); err != nil {
		t.Fatalf("reversed transaction must be balanced: %v", err)
	}

	if tx.Status != domain.StatusReversed {
		t.Errorf("original transaction status should be reversed, got %s", tx.Status)
	}

	// The sum of entries of original + reversal must cancel out to zero
	var combinedDebitSum int64
	for _, e := range tx.Entries {
		if e.EntryType == domain.EntryTypeDebit {
			combinedDebitSum += e.AmountMinor
		}
	}
	for _, e := range revTx.Entries {
		if e.EntryType == domain.EntryTypeCredit {
			combinedDebitSum -= e.AmountMinor
		}
	}

	if combinedDebitSum != 0 {
		t.Errorf("expected combined balance after reversal to be 0, got %d", combinedDebitSum)
	}
}
