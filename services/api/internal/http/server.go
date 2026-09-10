package http

import (
	"net/http"

	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/service"
	"github.com/google/uuid"
)

func NewRouter(fs *service.FinancialService, validHouseholdID uuid.UUID) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(fs)

	// Liveness & Readiness checks (No auth required)
	mux.HandleFunc("GET /livez", h.LivezHandler)
	mux.HandleFunc("GET /readyz", h.ReadyzHandler)

	// Protected API Routes (scoped to household)
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /overview", h.OverviewHandler)
	apiMux.HandleFunc("GET /accounts", h.ListAccountsHandler)
	apiMux.HandleFunc("GET /categories", h.ListCategoriesHandler)
	apiMux.HandleFunc("GET /transactions", h.ListTransactionsHandler)
	apiMux.HandleFunc("POST /transactions/expense", h.CreateExpenseHandler)

	// Wrap apiMux with HouseholdAuthMiddleware
	authWrapper := HouseholdAuthMiddleware(validHouseholdID)(apiMux)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", authWrapper))

	// Global Middlewares
	var handler http.Handler = mux
	handler = CORSMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	return handler
}
