package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnbalancedLedger      = errors.New("ledger transaction is unbalanced: sum of debits must equal sum of credits")
	ErrEmptyTransaction      = errors.New("ledger transaction must contain at least two entries")
	ErrNegativeEntryAmount   = errors.New("ledger entry amount must be strictly positive")
	ErrInvalidTransaction    = errors.New("invalid ledger transaction state")
	ErrSameAccountTransfer   = errors.New("cannot transfer between the same account")
	ErrAlreadyReversed       = errors.New("transaction has already been reversed")
)

type EntryType string

const (
	EntryTypeDebit  EntryType = "debit"
	EntryTypeCredit EntryType = "credit"
)

type TransactionStatus string

const (
	StatusDraft    TransactionStatus = "draft"
	StatusPosted   TransactionStatus = "posted"
	StatusReversed TransactionStatus = "reversed"
)

type TransactionSource string

const (
	SourceManual    TransactionSource = "manual"
	SourceImportCSV TransactionSource = "import_csv"
	SourceImportXLSX TransactionSource = "import_xlsx"
	SourceImportOFX TransactionSource = "import_ofx"
	SourceImportPDF TransactionSource = "import_pdf"
)

// LedgerEntry represents a single posting in a double-entry transaction.
type LedgerEntry struct {
	ID            uuid.UUID  `json:"id"`
	TransactionID uuid.UUID  `json:"transactionId"`
	AccountID     *uuid.UUID `json:"accountId,omitempty"`     // Target asset/liability account (nil if external category equity/income/expense)
	EntryType     EntryType  `json:"entryType"`              // debit or credit
	AmountMinor   int64      `json:"amountMinor"`            // Always positive
	CategoryID    *uuid.UUID `json:"categoryId,omitempty"`    // Category for income/expense reporting
	Description   string     `json:"description,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

// LedgerTransaction is the aggregate root of balanced journal entries.
type LedgerTransaction struct {
	ID          uuid.UUID              `json:"id"`
	HouseholdID uuid.UUID              `json:"householdId"`
	OccurredAt  time.Time              `json:"occurredAt"`
	Description string                 `json:"description"`
	Status      TransactionStatus      `json:"status"`
	Source      TransactionSource      `json:"source"`
	ReferenceID string                 `json:"referenceId,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Entries     []LedgerEntry          `json:"entries"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// Validate checks the fundamental double-entry invariant: Sum(debits) == Sum(credits).
func (tx *LedgerTransaction) Validate() error {
	if len(tx.Entries) < 2 {
		return ErrEmptyTransaction
	}

	var sumDebits int64
	var sumCredits int64

	for _, entry := range tx.Entries {
		if entry.AmountMinor <= 0 {
			return ErrNegativeEntryAmount
		}
		switch entry.EntryType {
		case EntryTypeDebit:
			sumDebits += entry.AmountMinor
		case EntryTypeCredit:
			sumCredits += entry.AmountMinor
		default:
			return fmt.Errorf("unknown entry type: %s", entry.EntryType)
		}
	}

	if sumDebits != sumCredits {
		return fmt.Errorf("%w: debits=%d, credits=%d (difference=%d)", ErrUnbalancedLedger, sumDebits, sumCredits, sumDebits-sumCredits)
	}

	return nil
}

// TotalDebits calculates total debits in the transaction.
func (tx *LedgerTransaction) TotalDebits() int64 {
	var total int64
	for _, e := range tx.Entries {
		if e.EntryType == EntryTypeDebit {
			total += e.AmountMinor
		}
	}
	return total
}

// TotalCredits calculates total credits in the transaction.
func (tx *LedgerTransaction) TotalCredits() int64 {
	var total int64
	for _, e := range tx.Entries {
		if e.EntryType == EntryTypeCredit {
			total += e.AmountMinor
		}
	}
	return total
}

// NetExpense returns the total expense recognized by this transaction.
// An expense is recognized when a category debit occurs without an asset debit.
// Inter-account transfers and card bill repayments produce NetExpense = 0.
func (tx *LedgerTransaction) NetExpense() int64 {
	var expense int64
	for _, e := range tx.Entries {
		// Category debit indicates an expense recognition
		if e.EntryType == EntryTypeDebit && e.CategoryID != nil {
			expense += e.AmountMinor
		}
	}
	return expense
}

// NetIncome returns the total income recognized by this transaction.
func (tx *LedgerTransaction) NetIncome() int64 {
	var income int64
	for _, e := range tx.Entries {
		// Category credit indicates an income recognition
		if e.EntryType == EntryTypeCredit && e.CategoryID != nil {
			income += e.AmountMinor
		}
	}
	return income
}

// NewExpense creates a standard balanced expense transaction.
// e.g. Paying R$ 100,00 for Groceries from Checking:
// - Debit: Expense (Category: Groceries) +10000
// - Credit: Asset (Checking Account) +10000
func NewExpense(
	householdID uuid.UUID,
	accountID uuid.UUID,
	categoryID uuid.UUID,
	amount Money,
	description string,
	occurredAt time.Time,
	source TransactionSource,
) (*LedgerTransaction, error) {
	if amount.AmountMinor <= 0 {
		return nil, ErrNegativeEntryAmount
	}

	txID := uuid.New()
	now := time.Now().UTC()

	tx := &LedgerTransaction{
		ID:          txID,
		HouseholdID: householdID,
		OccurredAt:  occurredAt,
		Description: description,
		Status:      StatusPosted,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
		Entries: []LedgerEntry{
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     nil, // Expense nominal account
				EntryType:     EntryTypeDebit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    &categoryID,
				Description:   description,
				CreatedAt:     now,
			},
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &accountID, // Asset/Liability source account
				EntryType:     EntryTypeCredit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil,
				Description:   description,
				CreatedAt:     now,
			},
		},
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}
	return tx, nil
}

// NewIncome creates a standard balanced income transaction.
// e.g. Receiving R$ 5.000,00 Salary into Checking:
// - Debit: Asset (Checking Account) +500000
// - Credit: Income (Category: Salary) +500000
func NewIncome(
	householdID uuid.UUID,
	accountID uuid.UUID,
	categoryID uuid.UUID,
	amount Money,
	description string,
	occurredAt time.Time,
	source TransactionSource,
) (*LedgerTransaction, error) {
	if amount.AmountMinor <= 0 {
		return nil, ErrNegativeEntryAmount
	}

	txID := uuid.New()
	now := time.Now().UTC()

	tx := &LedgerTransaction{
		ID:          txID,
		HouseholdID: householdID,
		OccurredAt:  occurredAt,
		Description: description,
		Status:      StatusPosted,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
		Entries: []LedgerEntry{
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &accountID, // Asset account increases
				EntryType:     EntryTypeDebit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil,
				Description:   description,
				CreatedAt:     now,
			},
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     nil, // Income nominal account
				EntryType:     EntryTypeCredit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    &categoryID,
				Description:   description,
				CreatedAt:     now,
			},
		},
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}
	return tx, nil
}

// NewTransfer creates a balanced transfer between two household accounts.
// CRITICAL FINANCIAL INVARIANT: Transfers are NOT expenses.
// - Credit: fromAccount (decreases source asset)
// - Debit: toAccount (increases destination asset)
// - Category: nil
// - NetExpense() == 0
func NewTransfer(
	householdID uuid.UUID,
	fromAccountID uuid.UUID,
	toAccountID uuid.UUID,
	amount Money,
	description string,
	occurredAt time.Time,
) (*LedgerTransaction, error) {
	if fromAccountID == toAccountID {
		return nil, ErrSameAccountTransfer
	}
	if amount.AmountMinor <= 0 {
		return nil, ErrNegativeEntryAmount
	}

	txID := uuid.New()
	now := time.Now().UTC()

	tx := &LedgerTransaction{
		ID:          txID,
		HouseholdID: householdID,
		OccurredAt:  occurredAt,
		Description: description,
		Status:      StatusPosted,
		Source:      SourceManual,
		CreatedAt:   now,
		UpdatedAt:   now,
		Entries: []LedgerEntry{
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &toAccountID, // Debit to destination asset
				EntryType:     EntryTypeDebit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil, // Zero expense
				Description:   fmt.Sprintf("Transfer from %s", fromAccountID),
				CreatedAt:     now,
			},
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &fromAccountID, // Credit to source asset
				EntryType:     EntryTypeCredit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil, // Zero expense
				Description:   fmt.Sprintf("Transfer to %s", toAccountID),
				CreatedAt:     now,
			},
		},
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}
	return tx, nil
}

// NewCreditCardPayment creates a transaction representing a credit card statement repayment.
// CRITICAL FINANCIAL INVARIANT: Credit card payments are NOT expenses twice.
// - Credit: checkingAccount (decreases checking balance)
// - Debit: cardAccount (decreases card liability balance)
// - NetExpense() == 0
func NewCreditCardPayment(
	householdID uuid.UUID,
	checkingAccountID uuid.UUID,
	cardAccountID uuid.UUID,
	amount Money,
	description string,
	occurredAt time.Time,
) (*LedgerTransaction, error) {
	if checkingAccountID == cardAccountID {
		return nil, ErrSameAccountTransfer
	}
	if amount.AmountMinor <= 0 {
		return nil, ErrNegativeEntryAmount
	}

	txID := uuid.New()
	now := time.Now().UTC()

	tx := &LedgerTransaction{
		ID:          txID,
		HouseholdID: householdID,
		OccurredAt:  occurredAt,
		Description: description,
		Status:      StatusPosted,
		Source:      SourceManual,
		CreatedAt:   now,
		UpdatedAt:   now,
		Entries: []LedgerEntry{
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &cardAccountID, // Debit reduces card liability
				EntryType:     EntryTypeDebit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil,
				Description:   "Credit card statement payment",
				CreatedAt:     now,
			},
			{
				ID:            uuid.New(),
				TransactionID: txID,
				AccountID:     &checkingAccountID, // Credit reduces checking cash
				EntryType:     EntryTypeCredit,
				AmountMinor:   amount.AmountMinor,
				CategoryID:    nil,
				Description:   "Debit checking for card payment",
				CreatedAt:     now,
			},
		},
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}
	return tx, nil
}

// Reverse creates an inverse transaction that mathematically undoes this transaction's effects.
func (tx *LedgerTransaction) Reverse(reason string) (*LedgerTransaction, error) {
	if tx.Status == StatusReversed {
		return nil, ErrAlreadyReversed
	}

	revID := uuid.New()
	now := time.Now().UTC()

	var reversedEntries []LedgerEntry
	for _, e := range tx.Entries {
		oppositeType := EntryTypeCredit
		if e.EntryType == EntryTypeCredit {
			oppositeType = EntryTypeDebit
		}

		reversedEntries = append(reversedEntries, LedgerEntry{
			ID:            uuid.New(),
			TransactionID: revID,
			AccountID:     e.AccountID,
			EntryType:     oppositeType,
			AmountMinor:   e.AmountMinor,
			CategoryID:    e.CategoryID,
			Description:   fmt.Sprintf("Reversal: %s", reason),
			CreatedAt:     now,
		})
	}

	revTx := &LedgerTransaction{
		ID:          revID,
		HouseholdID: tx.HouseholdID,
		OccurredAt:  now,
		Description: fmt.Sprintf("Reversal of %s (%s)", tx.Description, reason),
		Status:      StatusReversed,
		Source:      tx.Source,
		ReferenceID: tx.ID.String(),
		Entries:     reversedEntries,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := revTx.Validate(); err != nil {
		return nil, err
	}

	tx.Status = StatusReversed
	tx.UpdatedAt = now

	return revTx, nil
}
