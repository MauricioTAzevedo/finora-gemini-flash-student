# Security Policy

## Reporting Security Vulnerabilities

We take the security and integrity of Finora Gemini Flash Student seriously. Financial data confidentiality, system auditability, and cross-household isolation are paramount.

If you discover a security vulnerability:

1. **Do not open a public GitHub issue.**
2. Send a detailed report to `security@finora.local` (or reach out privately via GitHub security advisories).
3. Include:
   - Description of the vulnerability and impact.
   - Exact steps or proof-of-concept (PoC) to reproduce.
   - Affected components (`api`, `web`, `ai`, `database`).
4. We will acknowledge receipt within 48 hours and work with you on a coordinated disclosure timeline.

## Core Security Invariants

- **Multi-Tenant Isolation**: Every financial query must strictly scope to the authenticated user's authorized `household_id`.
- **Untrusted File Ingestion**: Uploaded spreadsheets (XLSX, CSV, OFX, PDF) are untrusted inputs. Formula injection, zip bombs, and prompt injections are treated as adversarial.
- **Zero Raw LLM SQL Execution**: LLMs never generate raw SQL directly against the financial database. All AI interactions use validated intent models, deterministic schemas, and server-side authorization checks.
- **Cryptographic Foundations**: Industry-standard cryptographic primitives (bcrypt/argon2id, standard HMAC-SHA256 tokens) are enforced; custom cryptographic schemes are strictly prohibited.
