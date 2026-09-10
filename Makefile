.PHONY: help dev test lint build docker-up docker-down seed

help:
	@echo "Finora Gemini Flash Student — Development Commands"
	@echo "  make dev         Run API, Web, and AI services locally"
	@echo "  make test        Run all unit and domain invariant test suites"
	@echo "  make docker-up   Start Postgres, Redis, NATS, MinIO in containers"
	@echo "  make docker-down Stop container services"
	@echo "  make seed        Seed synthetic Família Silva demo dataset"

dev:
	@echo "Starting development services..."

test:
	@echo "Running Go tests..."
	cd services/api && go test -v -race ./...
	@echo "Running Python AI tests..."
	cd services/ai && pytest -v

docker-up:
	docker compose up -d

docker-down:
	docker compose down
