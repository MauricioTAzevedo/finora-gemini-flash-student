# ADR-002: PostgreSQL as Source of Truth

## Status
Accepted

## Context
Financial software cannot afford eventual consistency or silent corruption for posted financial entries. We need ACID guarantees, foreign key constraints, check constraints, and battle-tested migration tooling.

## Decision
PostgreSQL 16 is chosen as the canonical transactional database. Derived stores (Redis for caching, pgvector/embeddings for search, NATS for async messaging) remain secondary and rebuildable from the primary database state.

## Consequences
- ACID transactions guarantee that ledger balancing rules are atomically verified.
- Outbox pattern enables reliable event publication without distributed two-phase commit.
