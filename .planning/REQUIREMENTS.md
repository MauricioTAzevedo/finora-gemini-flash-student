# Requirements Specification — Finora Gemini Flash Student

## Milestone 0: Foundation & Engineering Baseline
- [x] **REQ-M0-01**: Initialize Git repository and GitHub remote repository under `MauricioTAzevedo/finora-gemini-flash-student`.
- [x] **REQ-M0-02**: Establish repository hygiene (README, LICENSE, .gitignore, .editorconfig, .env.example, SECURITY, CONTRIBUTING, CI workflow).
- [x] **REQ-M0-03**: GSD Core project memory initialization (`.planning/`).
- [x] **REQ-M0-04**: Architectural ADRs and system documentation.
- [x] **REQ-M0-05**: Docker Compose definitions for PostgreSQL 16, Redis, NATS, MinIO.

## Milestone 1: Double-Entry Ledger & Excel Replacement
- [x] **REQ-M1-01**: Money representation using integer minor units (BRL cents) with zero binary floating-point drift.
- [x] **REQ-M1-02**: Core double-entry ledger aggregate enforcing $\sum \text{Debits} = \sum \text{Credits}$ for every posted transaction.
- [x] **REQ-M1-03**: Transfer invariants: inter-account transfers do not increment household expenses.
- [x] **REQ-M1-04**: Credit card accounting: purchase creates liability; payment clears liability without double-counting expenses.
- [x] **REQ-M1-05**: Multi-tenant authorization enforcing strict isolation by `household_id`.
- [x] **REQ-M1-06**: PostgreSQL schema migrations for Users, Households, Members, Accounts, Transactions, Entries, Categories, AuditEvents.
- [x] **REQ-M1-07**: Go 1.26 REST API service exposing health checks (`/livez`, `/readyz`), overview summary, and transaction endpoints.
- [x] **REQ-M1-08**: Next.js 15+ Web application providing overview dashboard, accounts manager, and dense transaction table.
- [x] **REQ-M1-09**: Synthetic demo seed ("Família Silva Demo") providing realistic Brazilian financial data.

## Milestone 2: Intelligent Imports & Migration Center
- [ ] **REQ-M2-01**: Ingestion pipeline for CSV, XLSX, and OFX bank statements.
- [ ] **REQ-M2-02**: Python AI service for ambiguous column mapping and confidence scoring.
- [ ] **REQ-M2-03**: Merchant normalization rules and fuzzy string matching.
- [ ] **REQ-M2-04**: Duplicate detection and reconciliation scoring.

## Milestone 3: Financial Intelligence & Briefs
- [ ] **REQ-M3-01**: Subscription and recurring charge detection.
- [ ] **REQ-M3-02**: Statistical anomaly detection for out-of-band utility bills.
- [ ] **REQ-M3-03**: Natural language querying with structured tool-calling.

## Milestone 4: Deterministic Forecasting & Digital Twin
- [ ] **REQ-M4-01**: 7-day, 30-day, 90-day cash flow forecast engine.
- [ ] **REQ-M4-02**: "Can we afford this?" installment simulator.
- [ ] **REQ-M4-03**: Goal planning integration with baseline forecasts.

## Milestone 5: Automations & Distributed Events
- [ ] **REQ-M5-01**: Transactional outbox relay to NATS JetStream.
- [ ] **REQ-M5-02**: Automation rule engine and DSL.
- [ ] **REQ-M5-03**: Developer Event Explorer.

## Milestone 6: Production Hardening
- [ ] **REQ-M6-01**: OpenTelemetry distributed tracing and metrics.
- [ ] **REQ-M6-02**: Automated backup and disaster recovery validation.
- [ ] **REQ-M6-03**: End-to-end performance benchmarks (100k+ transactions).
