# ADR-005: AI Provider Abstraction and Deterministic Fallback

## Status
Accepted

## Context
Relying solely on external cloud LLM APIs makes local development, CI test suites, and offline developer onboarding brittle and costly. Furthermore, financial transactions must remain fully functional even if an AI provider experiences an outage.

## Decision
All AI functionality is isolated behind an abstract provider interface (`AIProvider`) with strict Pydantic schemas for structured outputs. A deterministic `MockAIProvider` is implemented for tests and local development.

AI interprets messy data (e.g. spreadsheet column headers, ambiguous merchant strings), but deterministic engines calculate, the ledger records, and policies authorize.

## Consequences
- 100% offline development and automated CI testing with zero API keys required.
- Plug-and-play support for Google Gemini, OpenAI, or local models.
- Graceful degradation if external AI is unavailable.
