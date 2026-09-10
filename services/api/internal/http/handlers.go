package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/domain"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/forecasting"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/service"
	"github.com/google/uuid"
)

type Handler struct {
	financialService *service.FinancialService
}

func NewHandler(fs *service.FinancialService) *Handler {
	return &Handler{financialService: fs}
}

// LivezHandler checks process liveness.
func (h *Handler) LivezHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"service":   "finora-api",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadyzHandler checks process readiness.
func (h *Handler) ReadyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ready",
		"service":   "finora-api",
		"database":  "connected",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// OverviewHandler returns the financial snapshot for the authenticated household.
func (h *Handler) OverviewHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	overview, err := h.financialService.GetOverview(r.Context(), householdID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(overview)
}

// ListAccountsHandler returns all financial accounts.
func (h *Handler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	accounts, err := h.financialService.ListAccounts(r.Context(), householdID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts": accounts,
		"total":    len(accounts),
	})
}

// ListCategoriesHandler returns available categories.
func (h *Handler) ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	categories, err := h.financialService.ListCategories(r.Context(), householdID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
	})
}

// ListTransactionsHandler returns paginated ledger transactions.
func (h *Handler) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	txs, total, err := h.financialService.ListTransactions(r.Context(), householdID, limit, offset)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": txs,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

type CreateExpenseRequest struct {
	AccountID   string `json:"accountId"`
	CategoryID  string `json:"categoryId"`
	AmountMinor int64  `json:"amountMinor"`
	Description string `json:"description"`
}

// CreateExpenseHandler adds a new expense transaction.
func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", err.Error(), reqID)
		return
	}

	accID, err := uuid.Parse(req.AccountID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_ACCOUNT_ID", "Invalid account ID format", reqID)
		return
	}

	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_CATEGORY_ID", "Invalid category ID format", reqID)
		return
	}

	if req.AmountMinor <= 0 {
		writeJSONError(w, http.StatusBadRequest, "INVALID_AMOUNT", "Amount must be strictly positive integer minor units", reqID)
		return
	}

	money := domain.NewBRL(req.AmountMinor)
	tx, err := h.financialService.RecordExpense(
		r.Context(),
		householdID,
		accID,
		catID,
		money,
		req.Description,
		time.Now().UTC(),
		domain.SourceManual,
	)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "TRANSACTION_REJECTED", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(tx)
}

// ReconcileImportHandler handles file reconciliation for uploaded CSV or OFX statements.
func (h *Handler) ReconcileImportHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	report, err := h.financialService.ReconcileImport(r.Context(), householdID, format, r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "RECONCILIATION_FAILED", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(report)
}

// ForecastHandler computes deterministic day-by-day cash flow projections.
func (h *Handler) ForecastHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val > 0 && val <= 365 {
			days = val
		}
	}

	resp, err := h.financialService.GetForecast(r.Context(), householdID, days)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// AffordabilityScenarioHandler simulates a hypothetical installment purchase.
func (h *Handler) AffordabilityScenarioHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(RequestIDKey).(string)
	householdID, ok := r.Context().Value(HouseholdIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing household context", reqID)
		return
	}

	var scenario forecasting.Scenario
	if err := json.NewDecoder(r.Body).Decode(&scenario); err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", err.Error(), reqID)
		return
	}

	if scenario.InstallmentCount <= 0 {
		scenario.InstallmentCount = 1
	}
	if scenario.InstallmentMonthlyMinor <= 0 && scenario.TotalAmountMinor > 0 {
		scenario.InstallmentMonthlyMinor = scenario.TotalAmountMinor / int64(scenario.InstallmentCount)
	}

	comp, err := h.financialService.EvaluateScenario(r.Context(), householdID, scenario)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), reqID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comp)
}
