# ADR-003: Integer Minor Units for Money Representation

## Status
Accepted

## Context
IEEE 754 floating-point arithmetic introduces rounding errors (e.g. `0.1 + 0.2 = 0.30000000000000004`). In financial accounting, fractional cent errors accumulate and cause ledger reconciliation discrepancies.

## Decision
All monetary values are represented as integer minor units (`BIGINT` representing cents) paired with an explicit 3-letter currency code (e.g., `amountMinor: 12990, currency: "BRL"` represents R$ 129,90).

## Consequences
- Zero floating-point arithmetic in domain calculations and database schemas.
- Explicit formatting functions format minor units into localized display strings (e.g., `R$ 1.299,00` for pt-BR).
