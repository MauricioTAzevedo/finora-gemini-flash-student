package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
)

type Repository interface {
	// Households & Members
	GetHousehold(ctx context.Context, id uuid.UUID) (*domain.Household, error)
	GetHouseholdMember(ctx context.Context, householdID, userID uuid.UUID) (*domain.HouseholdMember, error)
	ListHouseholdMembers(ctx context.Context, householdID uuid.UUID) ([]domain.HouseholdMember, error)

	// Financial Accounts
	CreateAccount(ctx context.Context, account *domain.FinancialAccount) error
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.FinancialAccount, error)
	ListAccounts(ctx context.Context, householdID uuid.UUID) ([]domain.FinancialAccount, error)
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, deltaMinor int64) error

	// Categories
	ListCategories(ctx context.Context, householdID uuid.UUID) ([]domain.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)

	// Transactions & Entries
	RecordTransaction(ctx context.Context, tx *domain.LedgerTransaction) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*domain.LedgerTransaction, error)
	ListTransactions(ctx context.Context, householdID uuid.UUID, limit, offset int) ([]domain.LedgerTransaction, int, error)

	// Audit Trail
	RecordAuditEvent(ctx context.Context, event *domain.AuditEvent) error
	ListAuditEvents(ctx context.Context, householdID uuid.UUID, limit int) ([]domain.AuditEvent, error)
}

// MemoryRepository is an in-memory thread-safe implementation of Repository for testing and lightweight local demo.
type MemoryRepository struct {
	mu           sync.RWMutex
	households   map[uuid.UUID]*domain.Household
	members      map[string]*domain.HouseholdMember // key: householdID:userID
	accounts     map[uuid.UUID]*domain.FinancialAccount
	categories   map[uuid.UUID]*domain.Category
	transactions map[uuid.UUID]*domain.LedgerTransaction
	auditEvents  []domain.AuditEvent
}

func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		households:   make(map[uuid.UUID]*domain.Household),
		members:      make(map[string]*domain.HouseholdMember),
		accounts:     make(map[uuid.UUID]*domain.FinancialAccount),
		categories:   make(map[uuid.UUID]*domain.Category),
		transactions: make(map[uuid.UUID]*domain.LedgerTransaction),
		auditEvents:  make([]domain.AuditEvent, 0),
	}
	repo.seedDemoData()
	return repo
}

func (r *MemoryRepository) seedDemoData() {
	householdID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	userID := uuid.MustParse("a0000000-0000-0000-0000-000000000001")

	now := time.Now().UTC()

	r.households[householdID] = &domain.Household{
		ID:        householdID,
		Name:      "Família Silva Demo",
		Currency:  "BRL",
		Locale:    "pt-BR",
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.members[householdID.String()+":"+userID.String()] = &domain.HouseholdMember{
		ID:          uuid.New(),
		HouseholdID: householdID,
		UserID:      userID,
		Role:        domain.RoleOwner,
		CreatedAt:   now,
	}

	// Checking: Nubank Conta (R$ 8.420,00)
	checkingID := uuid.MustParse("c0000000-0000-0000-0000-000000000001")
	r.accounts[checkingID] = &domain.FinancialAccount{
		ID:                  checkingID,
		HouseholdID:         householdID,
		Name:                "Nubank Conta Principal",
		Type:                domain.AccountTypeChecking,
		Currency:            "BRL",
		InitialBalanceMinor: 842000,
		CurrentBalanceMinor: 842000,
		IsActive:            true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Savings: Reserva de Emergência (R$ 15.800,00)
	savingsID := uuid.MustParse("c0000000-0000-0000-0000-000000000002")
	r.accounts[savingsID] = &domain.FinancialAccount{
		ID:                  savingsID,
		HouseholdID:         householdID,
		Name:                "Reserva de Emergência",
		Type:                domain.AccountTypeSavings,
		Currency:            "BRL",
		InitialBalanceMinor: 1580000,
		CurrentBalanceMinor: 1580000,
		IsActive:            true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Credit Card: Nubank Platinum (-R$ 3.120,00 current bill)
	cardID := uuid.MustParse("c0000000-0000-0000-0000-000000000003")
	r.accounts[cardID] = &domain.FinancialAccount{
		ID:                  cardID,
		HouseholdID:         householdID,
		Name:                "Nubank Platinum",
		Type:                domain.AccountTypeCreditCard,
		Currency:            "BRL",
		InitialBalanceMinor: 0,
		CurrentBalanceMinor: -312000,
		CreditLimitMinor:    1200000,
		ClosingDay:          25,
		DueDay:              5,
		IsActive:            true,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Categories
	catGroceries := uuid.MustParse("d0000000-0000-0000-0000-000000000001")
	catHousing := uuid.MustParse("d0000000-0000-0000-0000-000000000002")
	catTransport := uuid.MustParse("d0000000-0000-0000-0000-000000000003")
	catHealth := uuid.MustParse("d0000000-0000-0000-0000-000000000004")
	catSalary := uuid.MustParse("d0000000-0000-0000-0000-000000000005")

	r.categories[catGroceries] = &domain.Category{ID: catGroceries, HouseholdID: householdID, Name: "Alimentação & Supermercado", Icon: "shopping-cart", Color: "#10B981"}
	r.categories[catHousing] = &domain.Category{ID: catHousing, HouseholdID: householdID, Name: "Moradia & Energia", Icon: "home", Color: "#3B82F6"}
	r.categories[catTransport] = &domain.Category{ID: catTransport, HouseholdID: householdID, Name: "Transporte & Combustível", Icon: "car", Color: "#F59E0B"}
	r.categories[catHealth] = &domain.Category{ID: catHealth, HouseholdID: householdID, Name: "Saúde & Farmácia", Icon: "heart", Color: "#EF4444"}
	r.categories[catSalary] = &domain.Category{ID: catSalary, HouseholdID: householdID, Name: "Salário & Renda", Icon: "briefcase", Color: "#059669"}

	// Seed recent transactions
	tx1, _ := domain.NewExpense(householdID, checkingID, catGroceries, domain.NewBRL(34050), "Supermercado Horizonte", now.AddDate(0, 0, -2), domain.SourceManual)
	tx2, _ := domain.NewExpense(householdID, cardID, catTransport, domain.NewBRL(18000), "Posto Central - Gasolina", now.AddDate(0, 0, -3), domain.SourceManual)
	tx3, _ := domain.NewIncome(householdID, checkingID, catSalary, domain.NewBRL(850000), "Salário Mensal", now.AddDate(0, 0, -10), domain.SourceManual)
	// Anomaly bill: Electricity bill higher than normal (R$ 487,00)
	tx4, _ := domain.NewExpense(householdID, checkingID, catHousing, domain.NewBRL(48700), "Conta de Energia Elétrica (CPFL)", now.AddDate(0, 0, -1), domain.SourceManual)

	r.transactions[tx1.ID] = tx1
	r.transactions[tx2.ID] = tx2
	r.transactions[tx3.ID] = tx3
	r.transactions[tx4.ID] = tx4
}

func (r *MemoryRepository) GetHousehold(ctx context.Context, id uuid.UUID) (*domain.Household, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.households[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (r *MemoryRepository) GetHouseholdMember(ctx context.Context, householdID, userID uuid.UUID) (*domain.HouseholdMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := householdID.String() + ":" + userID.String()
	m, ok := r.members[key]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (r *MemoryRepository) ListHouseholdMembers(ctx context.Context, householdID uuid.UUID) ([]domain.HouseholdMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.HouseholdMember
	for _, m := range r.members {
		if m.HouseholdID == householdID {
			list = append(list, *m)
		}
	}
	return list, nil
}

func (r *MemoryRepository) CreateAccount(ctx context.Context, account *domain.FinancialAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[account.ID] = account
	return nil
}

func (r *MemoryRepository) GetAccount(ctx context.Context, id uuid.UUID) (*domain.FinancialAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.accounts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (r *MemoryRepository) ListAccounts(ctx context.Context, householdID uuid.UUID) ([]domain.FinancialAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.FinancialAccount
	for _, a := range r.accounts {
		if a.HouseholdID == householdID {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (r *MemoryRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, deltaMinor int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.accounts[id]
	if !ok {
		return ErrNotFound
	}
	a.CurrentBalanceMinor += deltaMinor
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *MemoryRepository) ListCategories(ctx context.Context, householdID uuid.UUID) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.Category
	for _, c := range r.categories {
		if c.HouseholdID == householdID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (r *MemoryRepository) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (r *MemoryRepository) RecordTransaction(ctx context.Context, tx *domain.LedgerTransaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions[tx.ID] = tx

	// Apply entries to account balances
	for _, entry := range tx.Entries {
		if entry.AccountID != nil {
			acc, ok := r.accounts[*entry.AccountID]
			if ok {
				if acc.Type == domain.AccountTypeCreditCard {
					// For credit card, credit increases liability (more debt = negative balance)
					if entry.EntryType == domain.EntryTypeCredit {
						acc.CurrentBalanceMinor -= entry.AmountMinor
					} else {
						acc.CurrentBalanceMinor += entry.AmountMinor
					}
				} else {
					// For asset accounts (Checking, Savings, Cash): Debit increases asset (+), Credit decreases asset (-)
					if entry.EntryType == domain.EntryTypeDebit {
						acc.CurrentBalanceMinor += entry.AmountMinor
					} else {
						acc.CurrentBalanceMinor -= entry.AmountMinor
					}
				}
				acc.UpdatedAt = time.Now().UTC()
			}
		}
	}

	return nil
}

func (r *MemoryRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*domain.LedgerTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tx, ok := r.transactions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return tx, nil
}

func (r *MemoryRepository) ListTransactions(ctx context.Context, householdID uuid.UUID, limit, offset int) ([]domain.LedgerTransaction, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []domain.LedgerTransaction
	for _, tx := range r.transactions {
		if tx.HouseholdID == householdID {
			matched = append(matched, *tx)
		}
	}

	total := len(matched)
	if offset >= total {
		return []domain.LedgerTransaction{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return matched[offset:end], total, nil
}

func (r *MemoryRepository) RecordAuditEvent(ctx context.Context, event *domain.AuditEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.auditEvents = append(r.auditEvents, *event)
	return nil
}

func (r *MemoryRepository) ListAuditEvents(ctx context.Context, householdID uuid.UUID, limit int) ([]domain.AuditEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []domain.AuditEvent
	for i := len(r.auditEvents) - 1; i >= 0 && len(list) < limit; i-- {
		if r.auditEvents[i].HouseholdID == householdID {
			list = append(list, r.auditEvents[i])
		}
	}
	return list, nil
}
