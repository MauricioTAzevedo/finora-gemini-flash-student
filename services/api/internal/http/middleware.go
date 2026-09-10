package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIDKey   contextKey = "requestId"
	HouseholdIDKey contextKey = "householdId"
	UserIDKey      contextKey = "userId"
)

type ErrorResponse struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"requestId,omitempty"`
	} `json:"error"`
}

func writeJSONError(w http.ResponseWriter, statusCode int, code, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	resp.Error.RequestID = requestID
	_ = json.NewEncoder(w).Encode(resp)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		reqID, _ := r.Context().Value(RequestIDKey).(string)
		log.Printf("[HTTP] method=%s path=%s duration_ms=%d request_id=%s", r.Method, r.URL.Path, duration.Milliseconds(), reqID)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Household-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HouseholdAuthMiddleware enforces multi-tenant boundary checks.
// If the client supplies an unauthorized or mismatched household ID, access is denied with 403 Forbidden.
func HouseholdAuthMiddleware(validHouseholdID uuid.UUID) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawID := r.Header.Get("X-Household-ID")
			reqID, _ := r.Context().Value(RequestIDKey).(string)

			if rawID == "" {
				// Default to valid demo household for seamless developer onboarding if unset
				ctx := context.WithValue(r.Context(), HouseholdIDKey, validHouseholdID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			parsedID, err := uuid.Parse(rawID)
			if err != nil || parsedID != validHouseholdID {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN_HOUSEHOLD_ACCESS", "You do not have authorization to access this household's financial data.", reqID)
				return
			}

			ctx := context.WithValue(r.Context(), HouseholdIDKey, parsedID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
