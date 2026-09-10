# ADR-001: Modular Monolith as Initial Architecture

## Status
Accepted

## Context
A household financial management platform requires strong transactional guarantees across accounts, ledger transactions, entries, categories, and audit events. Decomposing into 10+ microservices upfront introduces network latency, distributed transaction complexity (2PC/Saga), operational overhead, and deployment friction without immediate scalability benefits.

## Decision
We adopt a **modular monolith** written in Go for the core API and ledger engine, with a separate Python service dedicated specifically to AI document parsing and inference tasks where Python's ecosystem is uniquely suited.

## Consequences
- Single transactional boundary in PostgreSQL for all financial mutations.
- Fast, straightforward local development via Docker or native runtimes.
- Clear internal package boundaries (`internal/domain/ledger`, `internal/service`, `internal/http`) that can be extracted cleanly in the future if required.
