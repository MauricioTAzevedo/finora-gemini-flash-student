export interface FinancialAccount {
  id: string;
  householdId: string;
  name: string;
  type: 'checking' | 'savings' | 'credit_card' | 'cash' | 'investment';
  currency: string;
  initialBalanceMinor: number;
  currentBalanceMinor: number;
  creditLimitMinor?: number;
  statementClosingDay?: number;
  statementDueDay?: number;
  isActive: boolean;
}

export interface LedgerEntry {
  id: string;
  transactionId: string;
  accountId?: string;
  entryType: 'debit' | 'credit';
  amountMinor: number;
  categoryId?: string;
  description?: string;
}

export interface LedgerTransaction {
  id: string;
  householdId: string;
  occurredAt: string;
  description: string;
  status: string;
  source: string;
  entries: LedgerEntry[];
}

export interface AlertDTO {
  id: string;
  type: 'warning' | 'info' | 'alert';
  title: string;
  message: string;
  evidence?: string;
}

export interface OverviewData {
  netAvailableCashMinor: number;
  netAvailableCashFormatted: string;
  monthIncomeMinor: number;
  monthIncomeFormatted: string;
  monthExpensesMinor: number;
  monthExpensesFormatted: string;
  upcomingObligationMinor: number;
  upcomingObligationFormatted: string;
  accounts: FinancialAccount[];
  recentTransactions: LedgerTransaction[];
  needsAttention: AlertDTO[];
}

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const DEMO_HOUSEHOLD_ID = 'b0000000-0000-0000-0000-000000000001';

// Seeded fallback state for zero-configuration client rendering
const SEED_FALLBACK: OverviewData = {
  netAvailableCashMinor: 2422000,
  netAvailableCashFormatted: 'R$ 24.220,00',
  monthIncomeMinor: 850000,
  monthIncomeFormatted: 'R$ 8.500,00',
  monthExpensesMinor: 100750,
  monthExpensesFormatted: 'R$ 1.007,50',
  upcomingObligationMinor: 312000,
  upcomingObligationFormatted: 'R$ 3.120,00',
  accounts: [
    {
      id: 'c0000000-0000-0000-0000-000000000001',
      householdId: DEMO_HOUSEHOLD_ID,
      name: 'Nubank Conta Principal',
      type: 'checking',
      currency: 'BRL',
      initialBalanceMinor: 842000,
      currentBalanceMinor: 842000,
      isActive: true,
    },
    {
      id: 'c0000000-0000-0000-0000-000000000002',
      householdId: DEMO_HOUSEHOLD_ID,
      name: 'Reserva de Emergência',
      type: 'savings',
      currency: 'BRL',
      initialBalanceMinor: 1580000,
      currentBalanceMinor: 1580000,
      isActive: true,
    },
    {
      id: 'c0000000-0000-0000-0000-000000000003',
      householdId: DEMO_HOUSEHOLD_ID,
      name: 'Nubank Platinum',
      type: 'credit_card',
      currency: 'BRL',
      initialBalanceMinor: 0,
      currentBalanceMinor: -312000,
      creditLimitMinor: 1200000,
      statementClosingDay: 25,
      statementDueDay: 5,
      isActive: true,
    },
  ],
  recentTransactions: [
    {
      id: 'tx-1',
      householdId: DEMO_HOUSEHOLD_ID,
      occurredAt: new Date(Date.now() - 86400000 * 1).toISOString(),
      description: 'Conta de Energia Elétrica (CPFL)',
      status: 'posted',
      source: 'manual',
      entries: [
        { id: 'e-1', transactionId: 'tx-1', entryType: 'debit', amountMinor: 48700, categoryId: 'd0000000-0000-0000-0000-000000000002' },
        { id: 'e-2', transactionId: 'tx-1', entryType: 'credit', amountMinor: 48700, accountId: 'c0000000-0000-0000-0000-000000000001' }
      ]
    },
    {
      id: 'tx-2',
      householdId: DEMO_HOUSEHOLD_ID,
      occurredAt: new Date(Date.now() - 86400000 * 2).toISOString(),
      description: 'Supermercado Horizonte',
      status: 'posted',
      source: 'manual',
      entries: [
        { id: 'e-3', transactionId: 'tx-2', entryType: 'debit', amountMinor: 34050, categoryId: 'd0000000-0000-0000-0000-000000000001' },
        { id: 'e-4', transactionId: 'tx-2', entryType: 'credit', amountMinor: 34050, accountId: 'c0000000-0000-0000-0000-000000000001' }
      ]
    },
    {
      id: 'tx-3',
      householdId: DEMO_HOUSEHOLD_ID,
      occurredAt: new Date(Date.now() - 86400000 * 3).toISOString(),
      description: 'Posto Central - Gasolina',
      status: 'posted',
      source: 'manual',
      entries: [
        { id: 'e-5', transactionId: 'tx-3', entryType: 'debit', amountMinor: 18000, categoryId: 'd0000000-0000-0000-0000-000000000003' },
        { id: 'e-6', transactionId: 'tx-3', entryType: 'credit', amountMinor: 18000, accountId: 'c0000000-0000-0000-0000-000000000003' }
      ]
    },
    {
      id: 'tx-4',
      householdId: DEMO_HOUSEHOLD_ID,
      occurredAt: new Date(Date.now() - 86400000 * 10).toISOString(),
      description: 'Salário Mensal',
      status: 'posted',
      source: 'manual',
      entries: [
        { id: 'e-7', transactionId: 'tx-4', entryType: 'debit', amountMinor: 850000, accountId: 'c0000000-0000-0000-0000-000000000001' },
        { id: 'e-8', transactionId: 'tx-4', entryType: 'credit', amountMinor: 850000, categoryId: 'd0000000-0000-0000-0000-000000000005' }
      ]
    }
  ],
  needsAttention: [
    {
      id: 'alert-electricity',
      type: 'warning',
      title: 'Conta de energia CPFL 115% acima da média',
      message: 'A conta de energia deste mês (R$ 487,00) está significativamente superior à média dos últimos 6 meses.',
      evidence: 'Média histórica: R$ 226,00 | Atual: R$ 487,00'
    },
    {
      id: 'alert-card-statement',
      type: 'info',
      title: 'Fatura Nubank Platinum em aberto',
      message: 'Fatura atual em R$ 3.120,00 com fechamento em 25 deste mês.',
      evidence: 'Limite disponível: R$ 8.880,00 de R$ 12.000,00'
    }
  ]
};

export async function fetchOverview(): Promise<OverviewData> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/overview`, {
      headers: {
        'X-Household-ID': DEMO_HOUSEHOLD_ID,
      },
      cache: 'no-store',
    });
    if (!res.ok) {
      return SEED_FALLBACK;
    }
    return await res.json();
  } catch {
    return SEED_FALLBACK;
  }
}
