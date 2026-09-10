# Finora Gemini Flash Student — Project Definition

## 1. Executive Summary

Finora Gemini Flash Student is an autonomous, AI-native household financial intelligence platform designed to replace Excel spreadsheets for Brazilian families and scale into a deterministic, double-entry financial operating system.

## 2. Core Problem & Thesis

- **The Problem**: Families track money across fragmented bank apps, credit cards, PDF bills, and fragile Excel spreadsheets. Existing apps are either simplistic expense trackers or bloated enterprise tools that fail to handle real-world Brazilian financial patterns (installments/parcelas, boleto deadlines, PIX transactions, and statement billing cycles).
- **Product Thesis**: Finora transforms fragmented household financial data into a unified, mathematically balanced financial ledger and decision engine that helps a family understand where their money went, where it is currently committed, and where it is projected to go.
- **Engineering Thesis**: Everything financially relevant becomes a traceable financial event that can be reconciled, explained, projected, and audited without relying on binary floating-point math or unverified LLM hallucinations.

## 3. Technology Stack & Architectural Decisions

- **Web Frontend**: Next.js 15+ (App Router), React 19, TypeScript, Tailwind CSS, Radix UI. First-class pt-BR locale and BRL formatting.
- **Backend Core**: Go 1.26 modular monolith enforcing double-entry ledger invariants, integer minor units (`amountMinor` in cents), multi-tenant household isolation, and transactional outbox.
- **AI & Intelligence Service**: Python 3.13 FastAPI micro-service with structured Pydantic schemas, deterministic `MockAIProvider` for local dev/testing, and Gemini model provider integration.
- **Database**: PostgreSQL 16 as the canonical source of truth with explicit SQL migrations.
- **Infrastructure**: Docker Compose for local PostgreSQL, Redis, NATS JetStream, and MinIO.

## 4. Key Financial Invariants

1. Money is strictly modeled as integer minor units (`amountMinor: 12990, currency: "BRL"`).
2. Every posted ledger transaction balances: $\sum \text{Debits} = \sum \text{Credits}$.
3. Inter-account transfers credit an asset and debit an asset; net household expenditure is zero.
4. Credit card purchase + statement repayment does not double-count household expenses.
5. All queries and mutations are isolated by `household_id` and authorized against the authenticated user's membership.
