package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	internalHttp "github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/http"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/repository"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/service"
	"github.com/google/uuid"
)

func TestCrossHouseholdTenantIsolationSecurity(t *testing.T) {
	// Household A: Authorized Demo Household
	householdA := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	// Household B: Alien Household (Attacker or another user's household)
	householdB := uuid.New()

	repo := repository.NewMemoryRepository()
	fs := service.NewFinancialService(repo)
	router := internalHttp.NewRouter(fs, householdA)

	// Case 1: Client sends request targeting Household B (Unauthorized Tenant Access)
	req := httptest.NewRequest("GET", "/api/v1/overview", nil)
	req.Header.Set("X-Household-ID", householdB.String())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// SECURITY INVARIANT: Must return 403 Forbidden
	if rr.Code != http.StatusForbidden {
		t.Fatalf("CRITICAL SECURITY FAILURE: Cross-household access returned status %d; expected %d Forbidden", rr.Code, http.StatusForbidden)
	}

	// Case 2: Authorized Household A request must succeed with 200 OK
	reqAuth := httptest.NewRequest("GET", "/api/v1/overview", nil)
	reqAuth.Header.Set("X-Household-ID", householdA.String())

	rrAuth := httptest.NewRecorder()
	router.ServeHTTP(rrAuth, reqAuth)

	if rrAuth.Code != http.StatusOK {
		t.Fatalf("Authorized household request failed with status %d; expected %d OK", rrAuth.Code, http.StatusOK)
	}
}

func TestHealthEndpoints(t *testing.T) {
	householdID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	repo := repository.NewMemoryRepository()
	fs := service.NewFinancialService(repo)
	router := internalHttp.NewRouter(fs, householdID)

	// Test /livez
	reqLive := httptest.NewRequest("GET", "/livez", nil)
	rrLive := httptest.NewRecorder()
	router.ServeHTTP(rrLive, reqLive)

	if rrLive.Code != http.StatusOK {
		t.Errorf("/livez returned %d, want 200", rrLive.Code)
	}

	// Test /readyz
	reqReady := httptest.NewRequest("GET", "/readyz", nil)
	rrReady := httptest.NewRecorder()
	router.ServeHTTP(rrReady, reqReady)

	if rrReady.Code != http.StatusOK {
		t.Errorf("/readyz returned %d, want 200", rrReady.Code)
	}
}
