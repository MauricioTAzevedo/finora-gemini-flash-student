import { formatBRL } from '@/lib/format';
import { CreditCard, Landmark, PiggyBank, Calendar, AlertCircle } from 'lucide-react';

export default function AccountsPage() {
  const accounts = [
    {
      id: '1',
      name: 'Nubank Conta Principal',
      type: 'checking',
      balanceMinor: 842000,
      description: 'Conta corrente para gastos do dia a dia e PIX',
      icon: Landmark,
      color: 'text-purple-600 bg-purple-50',
    },
    {
      id: '2',
      name: 'Reserva de Emergência',
      type: 'savings',
      balanceMinor: 1580000,
      description: 'Tesouro Selic / CDB 100% CDI',
      icon: PiggyBank,
      color: 'text-emerald-600 bg-emerald-50',
    },
    {
      id: '3',
      name: 'Nubank Platinum',
      type: 'credit_card',
      currentBillMinor: 312000,
      limitMinor: 1200000,
      closingDay: 25,
      dueDay: 5,
      description: 'Cartão compartilhado do casal',
      icon: CreditCard,
      color: 'text-indigo-600 bg-indigo-50',
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Contas & Cartões</h2>
        <p className="text-sm text-slate-500">Gestão dos ativos de liquidez e passivos de cartão da família.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
        {accounts.map((acc) => {
          const Icon = acc.icon;
          const isCard = acc.type === 'credit_card';

          return (
            <div
              key={acc.id}
              className="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex flex-col justify-between"
            >
              <div>
                <div className="flex items-center justify-between mb-4">
                  <div className={`p-2.5 rounded-lg ${acc.color}`}>
                    <Icon className="h-5 w-5" />
                  </div>
                  <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                    {acc.type === 'credit_card' ? 'Cartão de Crédito' : acc.type === 'savings' ? 'Poupança / Reserva' : 'Conta Corrente'}
                  </span>
                </div>

                <h3 className="font-semibold text-slate-900 text-base">{acc.name}</h3>
                <p className="text-xs text-slate-500 mt-0.5">{acc.description}</p>
              </div>

              <div className="mt-6 pt-4 border-t border-slate-100">
                {isCard ? (
                  <div className="space-y-3">
                    <div>
                      <span className="text-xs text-slate-500 font-medium">Fatura Atual (Passivo)</span>
                      <p className="text-2xl font-bold text-slate-900">{formatBRL(acc.currentBillMinor || 0)}</p>
                    </div>

                    {/* Progress bar */}
                    <div>
                      <div className="flex justify-between text-xs text-slate-500 mb-1">
                        <span>Limite Utilizado</span>
                        <span className="font-semibold text-slate-700">
                          {formatBRL(acc.currentBillMinor || 0)} de {formatBRL(acc.limitMinor || 0)}
                        </span>
                      </div>
                      <div className="h-2 w-full bg-slate-100 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-indigo-600 rounded-full"
                          style={{ width: `${((acc.currentBillMinor || 0) / (acc.limitMinor || 1)) * 100}%` }}
                        />
                      </div>
                    </div>

                    <div className="flex items-center gap-4 text-xs text-slate-600 bg-slate-50 p-2.5 rounded-lg">
                      <div className="flex items-center gap-1.5">
                        <Calendar className="h-3.5 w-3.5 text-slate-400" />
                        <span>Fecha dia {acc.closingDay}</span>
                      </div>
                      <div className="flex items-center gap-1.5">
                        <AlertCircle className="h-3.5 w-3.5 text-amber-500" />
                        <span>Vence dia {acc.dueDay}</span>
                      </div>
                    </div>
                  </div>
                ) : (
                  <div>
                    <span className="text-xs text-slate-500 font-medium">Saldo Disponível</span>
                    <p className="text-2xl font-bold text-emerald-700">{formatBRL(acc.balanceMinor || 0)}</p>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
