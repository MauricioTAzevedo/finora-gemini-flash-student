package service

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/importer"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/repository"
	"github.com/google/uuid"
)

type OverviewDTO struct {
	NetAvailableCashMinor   int64                     `json:"netAvailableCashMinor"`
	NetAvailableCashFmt     string                    `json:"netAvailableCashFormatted"`
	MonthIncomeMinor        int64                     `json:"monthIncomeMinor"`
	MonthIncomeFmt          string                    `json:"monthIncomeFormatted"`
	MonthExpensesMinor      int64                     `json:"monthExpensesMinor"`
	MonthExpensesFmt        string                    `json:"monthExpensesFormatted"`
	UpcomingObligationMinor int64                     `json:"upcomingObligationMinor"`
	UpcomingObligationFmt   string                    `json:"upcomingObligationFormatted"`
	Accounts                []domain.FinancialAccount `json:"accounts"`
	RecentTransactions      []domain.LedgerTransaction `json:"recentTransactions"`
	NeedsAttention          []AlertDTO                `json:"needsAttention"`
}

type AlertDTO struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // warning, info, alert
	Title    string `json:"title"`
	Message  string `json:"message"`
	Evidence string `json:"evidence,omitempty"`
}

type FinancialService struct {
	repo repository.Repository
}

func NewFinancialService(repo repository.Repository) *FinancialService {
	return &FinancialService{repo: repo}
}

func (s *FinancialService) GetOverview(ctx context.Context, householdID uuid.UUID) (*OverviewDTO, error) {
	accounts, err := s.repo.ListAccounts(ctx, householdID)
	if err != nil {
		return nil, err
	}

	var netAvailableCash int64
	var upcomingObligations int64

	for _, acc := range accounts {
		if acc.Type == domain.AccountTypeCreditCard {
			// Credit card negative balance represents current unsettled statement liability
			if acc.CurrentBalanceMinor < 0 {
				upcomingObligations += -acc.CurrentBalanceMinor
			}
		} else {
			netAvailableCash += acc.CurrentBalanceMinor
		}
	}

	txs, _, err := s.repo.ListTransactions(ctx, householdID, 50, 0)
	if err != nil {
		return nil, err
	}

	var monthIncome int64
	var monthExpenses int64

	for _, tx := range txs {
		monthIncome += tx.NetIncome()
		monthExpenses += tx.NetExpense()
	}

	alerts := []AlertDTO{
		{
			ID:       "alert-electricity",
			Type:     "warning",
			Title:    "Conta de energia elétrica CPFL acima da média",
			Message:  "A conta de energia deste mês (R$ 487,00) está 115% acima da média histórica residencial dos últimos 6 meses.",
			Evidence: "Média histórica: R$ 226,00 | Fatura atual: R$ 487,00",
		},
		{
			ID:       "alert-card-statement",
			Type:     "info",
			Title:    "Fatura Nubank Platinum em aberto",
			Message:  "Fatura atual em R$ 3.120,00 com fechamento previsto para o dia 25.",
			Evidence: "Limite disponível: R$ 8.880,00 de R$ 12.000,00",
		},
	}

	return &OverviewDTO{
		NetAvailableCashMinor:   netAvailableCash,
		NetAvailableCashFmt:     domain.NewBRL(netAvailableCash).FormatBRL(),
		MonthIncomeMinor:        monthIncome,
		MonthIncomeFmt:          domain.NewBRL(monthIncome).FormatBRL(),
		MonthExpensesMinor:      monthExpenses,
		MonthExpensesFmt:        domain.NewBRL(monthExpenses).FormatBRL(),
		UpcomingObligationMinor: upcomingObligations,
		UpcomingObligationFmt:   domain.NewBRL(upcomingObligations).FormatBRL(),
		Accounts:                accounts,
		RecentTransactions:      txs,
		NeedsAttention:          alerts,
	}, nil
}

func (s *FinancialService) RecordExpense(
	ctx context.Context,
	householdID, accountID, categoryID uuid.UUID,
	amount domain.Money,
	description string,
	occurredAt time.Time,
	source domain.TransactionSource,
) (*domain.LedgerTransaction, error) {
	tx, err := domain.NewExpense(householdID, accountID, categoryID, amount, description, occurredAt, source)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RecordTransaction(ctx, tx); err != nil {
		return nil, err
	}

	// Audit event
	_ = s.repo.RecordAuditEvent(ctx, &domain.AuditEvent{
		ID:           uuid.New(),
		HouseholdID:  householdID,
		Action:       "transaction.created",
		ResourceType: "ledger_transaction",
		ResourceID:   &tx.ID,
		Payload: map[string]interface{}{
			"amountMinor": amount.AmountMinor,
			"description": description,
			"type":        "expense",
		},
		CreatedAt: time.Now().UTC(),
	})

	return tx, nil
}

func (s *FinancialService) RecordTransfer(
	ctx context.Context,
	householdID, fromAccountID, toAccountID uuid.UUID,
	amount domain.Money,
	description string,
	occurredAt time.Time,
) (*domain.LedgerTransaction, error) {
	tx, err := domain.NewTransfer(householdID, fromAccountID, toAccountID, amount, description, occurredAt)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RecordTransaction(ctx, tx); err != nil {
		return nil, err
	}

	_ = s.repo.RecordAuditEvent(ctx, &domain.AuditEvent{
		ID:           uuid.New(),
		HouseholdID:  householdID,
		Action:       "transaction.created",
		ResourceType: "ledger_transaction",
		ResourceID:   &tx.ID,
		Payload: map[string]interface{}{
			"amountMinor": amount.AmountMinor,
			"description": description,
			"type":        "transfer",
		},
		CreatedAt: time.Now().UTC(),
	})

	return tx, nil
}

func (s *FinancialService) ListTransactions(ctx context.Context, householdID uuid.UUID, limit, offset int) ([]domain.LedgerTransaction, int, error) {
	return s.repo.ListTransactions(ctx, householdID, limit, offset)
}

func (s *FinancialService) ListAccounts(ctx context.Context, householdID uuid.UUID) ([]domain.FinancialAccount, error) {
	return s.repo.ListAccounts(ctx, householdID)
}

func (s *FinancialService) ListCategories(ctx context.Context, householdID uuid.UUID) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx, householdID)
}

type ImportReportDTO struct {
	TotalParsed    int                          `json:"totalParsed"`
	NewCount       int                          `json:"newCount"`
	CandidateCount int                          `json:"candidateCount"`
	DuplicateCount int                          `json:"duplicateCount"`
	Results        []importer.ReconciliationResult `json:"results"`
}

func (s *FinancialService) ReconcileImport(ctx context.Context, householdID uuid.UUID, format string, r io.Reader) (*ImportReportDTO, error) {
	var records []importer.CanonicalIngestionRecord
	var err error

	if strings.ToLower(format) == "ofx" {
		records, err = importer.ParseOFX(r)
	} else {
		records, err = importer.ParseCSV(r)
	}

	if err != nil {
		return nil, err
	}

	existingTxs, _, err := s.repo.ListTransactions(ctx, householdID, 1000, 0)
	if err != nil {
		return nil, err
	}

	var results []importer.ReconciliationResult
	newCount := 0
	candidateCount := 0
	duplicateCount := 0

	for _, rec := range records {
		res := importer.ReconcileTransaction(rec, existingTxs)
		switch res.Status {
		case importer.MatchStatusNew:
			newCount++
		case importer.MatchStatusCandidate:
			candidateCount++
		case importer.MatchStatusDuplicate:
			duplicateCount++
		}
		results = append(results, res)
	}

	return &ImportReportDTO{
		TotalParsed:    len(records),
		NewCount:       newCount,
		CandidateCount: candidateCount,
		DuplicateCount: duplicateCount,
		Results:        results,
	}, nil
}
