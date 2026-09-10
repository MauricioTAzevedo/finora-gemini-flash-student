package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalHttp "github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/http"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/repository"
	"github.com/MauricioTAzevedo/finora-gemini-flash-student/services/api/internal/service"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("=================================================================")
	log.Println("   FINORA GEMINI FLASH STUDENT — CORE API SERVER")
	log.Println("   Autonomous Household Financial Intelligence Platform")
	log.Println("=================================================================")

	// Initialize repository (Memory repository with Família Silva pre-seeded)
	repo := repository.NewMemoryRepository()
	financialService := service.NewFinancialService(repo)

	// Demo Household ID: Família Silva
	demoHouseholdID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")

	router := internalHttp.NewRouter(financialService, demoHouseholdID)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[SERVER] Listening on http://localhost:%s\n", port)
		log.Printf("[SERVER] Health check available at http://localhost:%s/livez\n", port)
		log.Printf("[SERVER] Overview API at http://localhost:%s/api/v1/overview\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[SERVER] Fatal server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[SERVER] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[SERVER] Server forced to shutdown: %v\n", err)
	}

	log.Println("[SERVER] Server exited cleanly.")
}
