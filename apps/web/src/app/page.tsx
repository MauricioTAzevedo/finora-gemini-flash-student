import { fetchOverview } from '@/lib/api';
import { formatDateBR } from '@/lib/format';
import { 
  TrendingUp, 
  TrendingDown, 
  Wallet, 
  AlertTriangle, 
  Info, 
  ArrowUpRight,
  ShieldCheck 
} from 'lucide-react';
import Link from 'next/link';

export default async function OverviewPage() {
  const overview = await fetchOverview();

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Visão Geral Familiar</h2>
          <p className="text-sm text-slate-500">Situação consolidada das finanças da Família Silva.</p>
        </div>
        <div className="flex items-center gap-2">
          <Link
            href="/imports"
            className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg bg-emerald-600 text-white text-sm font-medium hover:bg-emerald-700 transition-colors shadow-xs"
          >
            <ArrowUpRight className="h-4 w-4" />
            Importar Planilha / Extrato
          </Link>
        </div>
      </div>

      {/* Metric Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Net Available Cash */}
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
          <div className="flex items-center justify-between text-slate-500 mb-2">
            <span className="text-xs font-semibold uppercase tracking-wider">Disponível Líquido</span>
            <Wallet className="h-4 w-4 text-emerald-600" />
          </div>
          <p className="text-2xl font-bold text-slate-900">{overview.netAvailableCashFormatted}</p>
          <p className="text-xs text-slate-500 mt-1">Contas correntes + Reservas</p>
        </div>

        {/* Month Income */}
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
          <div className="flex items-center justify-between text-slate-500 mb-2">
            <span className="text-xs font-semibold uppercase tracking-wider">Entradas no Mês</span>
            <TrendingUp className="h-4 w-4 text-emerald-600" />
          </div>
          <p className="text-2xl font-bold text-emerald-700">{overview.monthIncomeFormatted}</p>
          <p className="text-xs text-slate-500 mt-1">Salários e rendimentos</p>
        </div>

        {/* Month Expenses */}
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
          <div className="flex items-center justify-between text-slate-500 mb-2">
            <span className="text-xs font-semibold uppercase tracking-wider">Gastos no Mês</span>
            <TrendingDown className="h-4 w-4 text-rose-600" />
          </div>
          <p className="text-2xl font-bold text-slate-900">{overview.monthExpensesFormatted}</p>
          <p className="text-xs text-slate-500 mt-1">Alimentação, contas e serviços</p>
        </div>

        {/* Upcoming Obligations */}
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
          <div className="flex items-center justify-between text-slate-500 mb-2">
            <span className="text-xs font-semibold uppercase tracking-wider">Compromissos a Vencer</span>
            <AlertTriangle className="h-4 w-4 text-amber-600" />
          </div>
          <p className="text-2xl font-bold text-amber-700">{overview.upcomingObligationFormatted}</p>
          <p className="text-xs text-slate-500 mt-1">Faturas de cartão em aberto</p>
        </div>
      </div>

      {/* Needs Attention / Anomalies */}
      {overview.needsAttention && overview.needsAttention.length > 0 && (
        <div className="space-y-3">
          <h3 className="text-sm font-semibold text-slate-700 uppercase tracking-wider">Atenção Necessária</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {overview.needsAttention.map((alert) => (
              <div
                key={alert.id}
                className={`p-4 rounded-xl border flex items-start gap-3 ${
                  alert.type === 'warning'
                    ? 'bg-amber-50/70 border-amber-200 text-amber-900'
                    : 'bg-blue-50/70 border-blue-200 text-blue-900'
                }`}
              >
                {alert.type === 'warning' ? (
                  <AlertTriangle className="h-5 w-5 text-amber-600 shrink-0 mt-0.5" />
                ) : (
                  <Info className="h-5 w-5 text-blue-600 shrink-0 mt-0.5" />
                )}
                <div>
                  <h4 className="text-sm font-semibold">{alert.title}</h4>
                  <p className="text-xs mt-1 text-slate-700">{alert.message}</p>
                  {alert.evidence && (
                    <span className="inline-block mt-2 text-[11px] font-mono bg-white/75 px-2 py-0.5 rounded border border-slate-200 text-slate-700">
                      Evidência: {alert.evidence}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Recent Ledger Transactions */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h3 className="text-base font-semibold text-slate-900">Últimos Lançamentos do Razão</h3>
            <span className="text-xs bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full font-medium">
              Partidas Dobradas
            </span>
          </div>
          <Link href="/transactions" className="text-xs font-medium text-emerald-700 hover:text-emerald-800">
            Ver todas as transações →
          </Link>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="px-6 py-3">Data</th>
                <th className="px-6 py-3">Descrição</th>
                <th className="px-6 py-3">Origem</th>
                <th className="px-6 py-3 text-right">Valor</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {overview.recentTransactions.map((tx) => {
                const isIncome = tx.description.includes('Salário');
                return (
                  <tr key={tx.id} className="hover:bg-slate-50/60 transition-colors">
                    <td className="px-6 py-3.5 whitespace-nowrap text-slate-600 font-medium">
                      {formatDateBR(tx.occurredAt)}
                    </td>
                    <td className="px-6 py-3.5 text-slate-900 font-medium">
                      {tx.description}
                    </td>
                    <td className="px-6 py-3.5 whitespace-nowrap">
                      <span className="text-xs font-medium bg-slate-100 text-slate-700 px-2 py-0.5 rounded">
                        {tx.source}
                      </span>
                    </td>
                    <td className={`px-6 py-3.5 whitespace-nowrap text-right font-semibold ${
                      isIncome ? 'text-emerald-700' : 'text-slate-900'
                    }`}>
                      {isIncome ? '+' : '-'}R$ {((tx.entries[0]?.amountMinor || 0) / 100).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
