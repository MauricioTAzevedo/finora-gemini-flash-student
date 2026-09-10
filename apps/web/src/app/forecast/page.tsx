'use client';

import { useState } from 'react';
import { formatBRL } from '@/lib/format';
import { 
  TrendingUp, 
  HelpCircle, 
  CheckCircle, 
  AlertTriangle, 
  Calculator, 
  ArrowRight,
  ShieldAlert,
  Zap
} from 'lucide-react';

export default function ForecastPage() {
  const [productName, setProductName] = useState('Geladeira Frost Free Inox');
  const [totalPrice, setTotalPrice] = useState('4800');
  const [installments, setInstallments] = useState(12);

  // Baseline data from Família Silva Demo
  const baseline = {
    startingBalanceMinor: 2422000, // R$ 24.220,00
    lowestBalanceMinor: 1870000,   // R$ 18.700,00 on day 15
    endBalanceMinor: 2720000,      // R$ 27.200,00
    reserveTargetMinor: 500000,    // R$ 5.000,00
  };

  const parsedTotal = parseFloat(totalPrice) || 0;
  const monthlyInstallment = installments > 0 ? parsedTotal / installments : 0;
  const monthlyInstallmentMinor = Math.round(monthlyInstallment * 100);

  // Simulation calculation
  const simulatedLowestBalanceMinor = baseline.lowestBalanceMinor - monthlyInstallmentMinor;
  const canAfford = simulatedLowestBalanceMinor >= baseline.reserveTargetMinor;
  const safetyMarginMinor = simulatedLowestBalanceMinor - baseline.reserveTargetMinor;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Projeções & Gêmeo Financeiro Digital</h2>
        <p className="text-sm text-slate-500">
          Projeções determinísticas de fluxo de caixa baseadas no razão contábil e simulador de decisões de compra parcelada.
        </p>
      </div>

      {/* Baseline Forecast Card */}
      <div className="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold text-slate-900 text-base">Projeção Base dos Próximos 30 Dias</h3>
            <span className="text-[11px] font-semibold bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full">
              Cálculo Determinístico
            </span>
          </div>
          <span className="text-xs text-slate-400">Horizonte: 30 dias</span>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 pt-2">
          <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
            <span className="text-xs text-slate-500 font-medium">Saldo Inicial Atual</span>
            <p className="text-xl font-bold text-slate-900 mt-1">{formatBRL(baseline.startingBalanceMinor)}</p>
          </div>

          <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
            <span className="text-xs text-slate-500 font-medium">Menor Saldo Projetado (Vale)</span>
            <p className="text-xl font-bold text-indigo-700 mt-1">{formatBRL(baseline.lowestBalanceMinor)}</p>
            <span className="text-[11px] text-slate-400">Previsão no dia 15 deste mês</span>
          </div>

          <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
            <span className="text-xs text-slate-500 font-medium">Saldo Final Projetado</span>
            <p className="text-xl font-bold text-emerald-700 mt-1">{formatBRL(baseline.endBalanceMinor)}</p>
            <span className="text-[11px] text-emerald-600 font-medium">+R$ 2.980,00 no ciclo</span>
          </div>

          <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
            <span className="text-xs text-slate-500 font-medium">Reserva Mínima Familiar</span>
            <p className="text-xl font-bold text-slate-700 mt-1">{formatBRL(baseline.reserveTargetMinor)}</p>
            <span className="text-[11px] text-slate-400">Margem segura atual: +R$ 13.700,00</span>
          </div>
        </div>
      </div>

      {/* What-If Scenario: "Can We Afford This?" */}
      <div className="bg-white rounded-xl border border-slate-200 p-6 shadow-xs space-y-6">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-indigo-50 text-indigo-600 rounded-lg">
            <Calculator className="h-6 w-6" />
          </div>
          <div>
            <h3 className="text-lg font-bold text-slate-900">Simulador: "Podemos Comprar Isso?"</h3>
            <p className="text-xs text-slate-500">
              Avalie o impacto de uma compra parcelada contra as obrigações futuras sem alterar o razão contábil.
            </p>
          </div>
        </div>

        {/* Inputs */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-100">
          <div>
            <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              O que você deseja comprar?
            </label>
            <input
              type="text"
              value={productName}
              onChange={(e) => setProductName(e.target.value)}
              className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Valor Total (R$)
            </label>
            <input
              type="number"
              value={totalPrice}
              onChange={(e) => setTotalPrice(e.target.value)}
              className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 font-mono focus:ring-2 focus:ring-emerald-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
              Número de Parcelas
            </label>
            <select
              value={installments}
              onChange={(e) => setInstallments(parseInt(e.target.value))}
              className="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
            >
              <option value={1}>À vista (1x)</option>
              <option value={3}>3x sem juros</option>
              <option value={6}>6x sem juros</option>
              <option value={10}>10x sem juros</option>
              <option value={12}>12x sem juros</option>
              <option value={18}>18x</option>
              <option value={24}>24x</option>
            </select>
          </div>
        </div>

        {/* Simulation Results */}
        <div className={`p-5 rounded-xl border ${
          canAfford ? 'bg-emerald-50/60 border-emerald-200' : 'bg-rose-50/60 border-rose-200'
        }`}>
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-3">
              {canAfford ? (
                <CheckCircle className="h-6 w-6 text-emerald-600 shrink-0 mt-0.5" />
              ) : (
                <ShieldAlert className="h-6 w-6 text-rose-600 shrink-0 mt-0.5" />
              )}
              <div>
                <h4 className="text-base font-bold text-slate-900">
                  {canAfford ? 'Simulação Favorável: Cabe no Orçamento Familiar' : 'Alerta: Risco de Comprometimento de Reserva'}
                </h4>
                <p className="text-sm text-slate-700 mt-1">
                  A compra de <span className="font-semibold text-slate-900">"{productName}"</span> adiciona{' '}
                  <span className="font-semibold text-slate-900">{formatBRL(monthlyInstallmentMinor)}/mês</span> durante{' '}
                  <span className="font-semibold text-slate-900">{installments} parcelas</span>.
                </p>
                <p className="text-xs text-slate-600 mt-2">
                  O menor saldo de caixa projetado nos próximos 6 meses passará de{' '}
                  <span className="font-semibold text-slate-800">{formatBRL(baseline.lowestBalanceMinor)}</span> para{' '}
                  <span className="font-semibold text-slate-800">{formatBRL(simulatedLowestBalanceMinor)}</span>.{' '}
                  {canAfford ? (
                    <span className="text-emerald-800 font-medium">
                      Permanece R$ {(safetyMarginMinor / 100).toLocaleString('pt-BR', { minimumFractionDigits: 2 })} acima da sua reserva mínima de segurança.
                    </span>
                  ) : (
                    <span className="text-rose-800 font-medium">
                      Viola a reserva mínima de segurança de R$ {(baseline.reserveTargetMinor / 100).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}.
                    </span>
                  )}
                </p>
              </div>
            </div>

            <span className={`text-xs font-bold px-3 py-1 rounded-full ${
              canAfford ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'
            }`}>
              {canAfford ? 'Suporte à Decisão: Seguro' : 'Atenção Necessária'}
            </span>
          </div>

          <div className="mt-4 pt-3 border-t border-slate-200/60 flex items-center justify-between text-xs text-slate-500">
            <span>* Esta análise é determinística e visa apoiar a tomada de decisão da família, não constituindo conselho financeiro autoritativo.</span>
            <span className="font-medium text-slate-700">Digital Twin Engine v1.0</span>
          </div>
        </div>
      </div>
    </div>
  );
}
