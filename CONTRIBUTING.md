# Contributing to Finora Gemini Flash Student

Welcome! Finora is an AI-native household financial intelligence platform designed with production-grade engineering principles.

## Development Workflow

1. **GSD Development Operating System**: Finora uses GSD (`.planning/`) to structure milestones, phases, and verifiable requirements.
2. **Monorepo Architecture**:
   - `services/api`: Go 1.26 core API, double-entry ledger, accounting invariants.
   - `apps/web`: Next.js 15+ TypeScript frontend with Tailwind CSS.
   - `services/ai`: Python 3.13 FastAPI service with structured AI models and fallback providers.
   - `database`: PostgreSQL migrations and synthetic demo seeds.

## Pull Request Guidelines

- All tests must pass:
  - Go: `cd services/api && go test -v -race ./...`
  - Python: `cd services/ai && pytest`
  - Web: `cd apps/web && pnpm lint && pnpm build`
- No binary floating-point numbers for currency amounts (`amountMinor` integer cents).
- Maintain double-entry balancing rules ($\sum \text{Debits} = \sum \text{Credits}$).
- Provide tests for any new accounting or authorization rules.
- Follow Conventional Commits: `feat(ledger):`, `fix(auth):`, `test(imports):`, `docs(adr):`.
