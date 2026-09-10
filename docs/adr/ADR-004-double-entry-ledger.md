# ADR-004: Double-Entry Ledger Engine

## Status
Accepted

## Context
Simple single-entry expense trackers cannot accurately model transfers between accounts, credit card debt accumulation vs statement repayment, refunds, or asset tracking without introducing balance inconsistencies and double-counted spending.

## Decision
Finora models all financial events via an accounting-grade double-entry ledger internally. Every `LedgerTransaction` consists of balanced `LedgerEntry` records satisfying:
$$\sum \text{Debits} = \sum \text{Credits}$$

The consumer UI abstracts away accounting terminology (debit/credit) into intuitive concepts (income, expense, transfer, payment), while the backend strictly enforces double-entry invariants.

## Consequences
- Transfers between Checking and Savings credit one asset and debit another, resulting in zero net household expense.
- Credit card purchases credit liability and debit expense. Statement payments debit liability and credit checking, preventing double-counting expenses.
- Auditability is mathematically guaranteed.
