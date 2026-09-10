# Project State — Finora Gemini Flash Student

## Current Context
- **Active Milestone**: Milestone 1: Double-Entry Ledger & Excel Replacement
- **Current Phase**: Phase 1 Complete (Foundation & Core Ledger Established)
- **Active Repository**: https://github.com/MauricioTAzevedo/finora-gemini-flash-student

## Completed Deliverables
- [x] **Milestone 0: Engineering Foundation**
  - Git repository initialized on `main` and pushed to GitHub: `https://github.com/MauricioTAzevedo/finora-gemini-flash-student`.
  - Repository hygiene: `README.md`, `LICENSE` (MIT), `.gitignore`, `.editorconfig`, `.env.example`, `SECURITY.md`, `CONTRIBUTING.md`.
  - CI Workflow: `.github/workflows/ci.yml` (Go tests, Python tests, Next.js build).
  - Local Containers: `docker-compose.yml` for PostgreSQL 16, Redis 7, NATS JetStream, MinIO.
  - Task Runners: `Makefile` and `dev.ps1` for Windows PowerShell.
  - Architecture Documentation: `SYSTEM_OVERVIEW.md`, `DATA_MODEL.md`, `THREAT_MODEL.md`.
  - Architecture Decision Records: `ADR-001` through `ADR-005`.
  - GSD Core v1.13.0 installed and initialized in `.agents/` and `.planning/`.

- [x] **Milestone 1: Double-Entry Ledger & Excel Replacement (Vertical Slice)**
  - Canonical PostgreSQL schema (`database/migrations/000001_initial_schema.sql`): Users, Households, Members, Accounts, Ledger Transactions, Entries, Categories, Outbox, and Audit Events.
  - Synthetic demo seed (`database/seeds/demo_seed.sql`): "Família Silva Demo" (Checking, Savings, Nubank Platinum credit card, utility bills, groceries).
  - Go Core Ledger Engine (`services/api/internal/domain/`):
    - `Money` value object with integer minor units (`amountMinor` cents) and pt-BR BRL formatting (`R$ 1.250,50`). Zero floating point drift.
    - Double-entry ledger aggregate enforcing $\sum \text{Debits} = \sum \text{Credits}$ for every posted transaction.
    - Transfer invariant verified: transfers between accounts produce `NetExpense = 0`.
    - Credit card invariant verified: purchase produces expense; card bill repayment produces `NetExpense = 0` (no double counting).
    - Transaction reversals creating inverse postings.
    - Security isolation test: cross-household access strictly denied with 403 Forbidden.
    - REST API handlers for overview, accounts, categories, transactions, `/livez`, and `/readyz`.
  - Python AI Service (`services/ai/`):
    - FastAPI micro-service with structured Pydantic schemas.
    - Abstract `AIProvider` with deterministic `MockAIProvider` for 100% offline local development and CI.
    - Column mapping for Brazilian spreadsheets (`"Descrição"`, `"Valor Pg."`, `"Data"`, `"Cartão M."`), merchant normalization, and transaction categorization.
  - Next.js 15+ Web Application (`apps/web/`):
    - pt-BR first-class financial overview dashboard: Net available cash card, monthly income vs expenses, upcoming obligations card, anomaly alerts.
    - Dense spreadsheet-inspired transactions table with search, category filtering, and transaction creation modal.
    - Accounts and Credit Cards page with credit limit progress bar and statement billing dates.
    - Import & Migration Center with spreadsheet analysis, AI column mapping review, confidence indicators, duplicate detection, and import rollback.

- [x] **Milestone 2: Intelligent Imports & Migration Center**
  - Deterministic Brazilian OFX statement parser (`services/api/internal/importer/ofx.go`).
  - Brazilian financial CSV parser (`services/api/internal/importer/csv.go`) with auto-delimiter detection (`;` and `,`).
  - Reconciliation & Deduplication engine (`services/api/internal/importer/reconciliation.go`) calculating multi-signal match scores to prevent duplicate balance inflation.
  - HTTP import endpoint: `POST /api/v1/imports/reconcile` returning explainable match evidence.
  - Synthetic Brazilian test fixtures in `database/samples/` (`extrato_nubank_agosto.ofx`, `planilha_gastos_familia.csv`).

- [x] **Milestone 4: Deterministic Forecasting & Digital Twin**
  - Deterministic cash flow forecasting engine (`services/api/internal/forecasting/engine.go`) computing exact daily balances across 7, 30, 90, and 180 days without LLM math hallucinations.
  - "Can We Afford This?" What-If scenario simulation evaluating installment commitments (e.g. 12x parcelas) against household baseline and minimum reserve target.
  - HTTP endpoints: `GET /api/v1/forecast?days=30` and `POST /api/v1/scenarios/affordability`.
  - Next.js interactive UI (`apps/web/src/app/forecast/page.tsx`) with real-time affordability calculator and cash impact explanation.

## Test Verification Summary
1. **Go Ledger, Importer & Forecasting Suite**: 15 unit/HTTP tests passing (`go test -v ./...`).
2. **Python AI Service Suite**: 3 pytest tests passing (`pytest -v`).
3. **Web Production Build**: All 8 static and dynamic routes compiled successfully with zero type errors (`pnpm build`).

## Next Priority (Milestone 3 & 5)
- **Milestone 3**: Financial Intelligence: Subscription detection engine & Statistical anomaly detection (CPFL electricity baseline deviation).
- **Milestone 5**: Transactional Outbox relay & Event Explorer.
