'use client';

import { useState, useEffect } from 'react';
import { 
  Sparkles, 
  Search, 
  AlertTriangle, 
  CheckCircle2, 
  TrendingUp, 
  Calendar, 
  CreditCard, 
  Repeat, 
  ShieldAlert,
  HelpCircle,
  ArrowRight,
  Filter
} from 'lucide-react';

const API_BASE = 'http://localhost:8080/api/v1';
const HEADERS = {
  'Content-Type': 'application/json',
  'X-Household-Id': 'b0000000-0000-0000-0000-000000000001',
};

interface Subscription {
  id: string;
  merchantName: string;
  categoryName: string;
  cadence: string;
  amount: { formatted: string; amountMinor: number };
  annualizedCost: { formatted: string; amountMinor: number };
  lastChargedAt: string;
  nextExpectedAt: string;
  occurrenceCount: number;
  priceDriftNote?: string;
  confidence: number;
}

interface Anomaly {
  id: string;
  type: string;
  severity: string;
  title: string;
  description: string;
  amount: { formatted: string; amountMinor: number };
  historicalMean: { formatted: string; amountMinor: number };
  zScore: number;
  percentageDiff: number;
  categoryName: string;
  merchantName: string;
  date: string;
  actionAdvice: string;
}

interface QueryResult {
  filter: {
    category?: string;
    merchant?: string;
    explanation: string;
    naturalQuery: string;
  };
  totalCount: number;
  totalAmount: { formatted: string; amountMinor: number };
  summaryText: string;
  transactions: Array<{
    Description: string;
    Category: string;
    AmountMinor: number;
    Currency: string;
    Date: string;
  }>;
}

export default function IntelligencePage() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [anomalies, setAnomalies] = useState<Anomaly[]>([]);
  const [loading, setLoading] = useState(true);

  // Natural language query state
  const [searchQuery, setSearchQuery] = useState('Quanto gastamos com mercado nos últimos 3 meses?');
  const [queryResult, setQueryResult] = useState<QueryResult | null>(null);
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    async function loadIntelligence() {
      try {
        const [subsRes, anomsRes] = await Promise.all([
          fetch(`${API_BASE}/intelligence/subscriptions`, { headers: HEADERS }),
          fetch(`${API_BASE}/intelligence/anomalies`, { headers: HEADERS })
        ]);

        if (subsRes.ok) {
          const data = await subsRes.json();
          setSubscriptions(data.subscriptions || []);
        }
        if (anomsRes.ok) {
          const data = await anomsRes.json();
          setAnomalies(data.anomalies || []);
        }
      } catch (err) {
        console.error('Failed to load intelligence data:', err);
      } finally {
        setLoading(false);
      }
    }

    loadIntelligence();
    handleSearch('Quanto gastamos com mercado nos últimos 3 meses?');
  }, []);

  const handleSearch = async (queryText: string) => {
    if (!queryText.trim()) return;
    setIsSearching(true);
    setSearchQuery(queryText);

    try {
      const res = await fetch(`${API_BASE}/intelligence/query`, {
        method: 'POST',
        headers: HEADERS,
        body: JSON.stringify({ query: queryText })
      });
      if (res.ok) {
        const data = await res.json();
        setQueryResult(data);
      }
    } catch (err) {
      console.error('NLP Query execution failed:', err);
    } finally {
      setIsSearching(false);
    }
  };

  const sampleQueries = [
    'Quanto gastamos com mercado nos últimos 3 meses?',
    'Gastos acima de 200 no transporte',
    'Conta de luz nos últimos 6 meses',
    'Gastos com restaurantes neste mês'
  ];

  const formatBRL = (cents: number) => {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(cents / 100);
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <div className="flex items-center gap-2 mb-1">
          <span className="px-2 py-0.5 text-xs font-semibold rounded bg-purple-100 text-purple-800 flex items-center gap-1">
            <Sparkles className="h-3 w-3" />
            Milestone 3: Inteligência Financeira Autônoma
          </span>
        </div>
        <h1 className="text-2xl font-bold text-slate-900">Inteligência, Recorrências & Anomalias</h1>
        <p className="text-sm text-slate-500">
          Detecção de assinaturas com análise de variação de preços, diagnóstico estatístico de contas (Z-score) e busca em linguagem natural via AST segura.
        </p>
      </div>

      {/* 1. Natural Language Query Interface */}
      <div className="bg-white rounded-xl border border-slate-200 p-6 shadow-xs">
        <div className="flex items-center gap-2 mb-3">
          <Search className="h-5 w-5 text-emerald-600" />
          <h2 className="text-base font-bold text-slate-900">Busca em Linguagem Natural (DSL Segura)</h2>
        </div>
        <p className="text-xs text-slate-500 mb-4">
          Faça perguntas em português sobre os gastos da família. O interpretador converte sua frase em uma árvore de filtros estrita, sem executar comandos SQL brutos.
        </p>

        {/* Input Bar */}
        <form 
          onSubmit={(e) => {
            e.preventDefault();
            handleSearch(searchQuery);
          }}
          className="flex gap-2 mb-3"
        >
          <div className="relative flex-1">
            <Search className="absolute left-3 top-3 h-4 w-4 text-slate-400" />
            <input 
              type="text" 
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Ex: Quanto gastamos com supermercado nos últimos 3 meses?"
              className="w-full pl-9 pr-4 py-2.5 text-sm rounded-lg border border-slate-300 focus:outline-hidden focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900"
            />
          </div>
          <button 
            type="submit"
            disabled={isSearching}
            className="px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white text-sm font-semibold rounded-lg transition-colors flex items-center gap-2 disabled:opacity-50 cursor-pointer"
          >
            {isSearching ? 'Interpretando...' : 'Consultar'}
            <ArrowRight className="h-4 w-4" />
          </button>
        </form>

        {/* Quick Query Chips */}
        <div className="flex flex-wrap gap-2 mb-6">
          <span className="text-xs text-slate-400 flex items-center gap-1 self-center">Sugestões:</span>
          {sampleQueries.map((q, idx) => (
            <button
              key={idx}
              type="button"
              onClick={() => handleSearch(q)}
              className="text-xs px-2.5 py-1 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium transition-colors cursor-pointer"
            >
              {q}
            </button>
          ))}
        </div>

        {/* Query Result Box */}
        {queryResult && (
          <div className="bg-slate-50 rounded-lg p-4 border border-slate-200">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-2 mb-3 pb-3 border-b border-slate-200">
              <div className="flex items-center gap-2">
                <Filter className="h-4 w-4 text-emerald-600" />
                <span className="text-xs font-semibold text-slate-700">{queryResult.filter.explanation}</span>
              </div>
              <div className="text-right">
                <span className="text-xs text-slate-500 mr-2">Total no período:</span>
                <span className="text-base font-bold text-slate-900">{queryResult.totalAmount.formatted}</span>
              </div>
            </div>

            <p className="text-sm font-medium text-slate-800 mb-3">{queryResult.summaryText}</p>

            {queryResult.transactions.length > 0 && (
              <div className="overflow-x-auto">
                <table className="w-full text-xs text-left">
                  <thead className="bg-white text-slate-500 uppercase tracking-wider font-semibold border-b border-slate-200">
                    <tr>
                      <th className="py-2 px-3">Data</th>
                      <th className="py-2 px-3">Estabelecimento / Descrição</th>
                      <th className="py-2 px-3">Categoria</th>
                      <th className="py-2 px-3 text-right">Valor</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-200 bg-white">
                    {queryResult.transactions.map((tx, i) => (
                      <tr key={i} className="hover:bg-slate-50">
                        <td className="py-2 px-3 text-slate-600 whitespace-nowrap">
                          {new Date(tx.Date).toLocaleDateString('pt-BR')}
                        </td>
                        <td className="py-2 px-3 font-medium text-slate-900">{tx.Description}</td>
                        <td className="py-2 px-3 text-slate-600">{tx.Category}</td>
                        <td className="py-2 px-3 text-right font-bold text-slate-900">
                          {formatBRL(tx.AmountMinor)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {/* 2. Statistical Anomalies Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <ShieldAlert className="h-5 w-5 text-amber-600" />
            <h2 className="text-lg font-bold text-slate-900">Detecção Estatística de Anomalias (Z-Score & Baselines)</h2>
          </div>
          <span className="text-xs px-2.5 py-1 rounded-full bg-amber-100 text-amber-800 font-semibold">
            {anomalies.length} anomalia(s) detectada(s)
          </span>
        </div>

        {anomalies.length === 0 ? (
          <div className="bg-white rounded-xl border border-slate-200 p-6 text-center text-slate-500">
            Nenhuma anomalia financeira detectada nas faturas recentes.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {anomalies.map((anom) => (
              <div 
                key={anom.id}
                className="bg-white rounded-xl border border-amber-200 p-5 shadow-xs flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="px-2 py-0.5 text-xs font-semibold rounded bg-amber-100 text-amber-800 flex items-center gap-1">
                      <AlertTriangle className="h-3 w-3 text-amber-600" />
                      Desvio de {anom.zScore}σ ({anom.percentageDiff}% acima da média)
                    </span>
                    <span className="text-xs text-slate-400 font-medium">
                      {new Date(anom.date).toLocaleDateString('pt-BR')}
                    </span>
                  </div>
                  <h3 className="text-base font-bold text-slate-900 mb-1">{anom.title}</h3>
                  <p className="text-sm text-slate-600 mb-3">{anom.description}</p>
                </div>

                <div className="bg-amber-50 border border-amber-200 rounded-lg p-3 text-xs text-amber-900 flex items-start gap-2">
                  <HelpCircle className="h-4 w-4 text-amber-700 shrink-0 mt-0.5" />
                  <div>
                    <span className="font-semibold">Ação Recomendada: </span>
                    {anom.actionAdvice}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 3. Subscriptions & Recurring Charges Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Repeat className="h-5 w-5 text-purple-600" />
            <h2 className="text-lg font-bold text-slate-900">Assinaturas e Mensalidades Recorrentes</h2>
          </div>
          <span className="text-xs px-2.5 py-1 rounded-full bg-purple-100 text-purple-800 font-semibold">
            {subscriptions.length} assinaturas mapeadas
          </span>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {subscriptions.map((sub) => (
            <div 
              key={sub.id} 
              className="bg-white rounded-xl border border-slate-200 p-5 shadow-xs flex flex-col justify-between"
            >
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-[11px] font-semibold uppercase tracking-wider text-purple-700 bg-purple-50 px-2 py-0.5 rounded">
                    {sub.cadence === 'MONTHLY' ? 'Mensal' : sub.cadence}
                  </span>
                  <span className="text-[11px] text-emerald-700 font-medium flex items-center gap-1">
                    <CheckCircle2 className="h-3 w-3" />
                    {Math.round(sub.confidence * 100)}% confiança
                  </span>
                </div>

                <h3 className="text-base font-bold text-slate-900">{sub.merchantName}</h3>
                <p className="text-xs text-slate-500 mb-3">{sub.categoryName}</p>

                <div className="space-y-1 mb-3">
                  <div className="text-2xl font-bold text-slate-900">
                    {sub.amount.formatted}
                    <span className="text-xs font-normal text-slate-500"> /mês</span>
                  </div>
                  <div className="text-xs text-slate-500">
                    Custo anual: <span className="font-semibold text-slate-700">{sub.annualizedCost.formatted}</span>
                  </div>
                </div>

                {sub.priceDriftNote && (
                  <div className="bg-rose-50 border border-rose-200 rounded p-2 text-xs text-rose-800 font-medium mb-3 flex items-center gap-1.5">
                    <TrendingUp className="h-3.5 w-3.5 text-rose-600 shrink-0" />
                    <span>{sub.priceDriftNote}</span>
                  </div>
                )}
              </div>

              <div className="pt-3 border-t border-slate-100 text-[11px] text-slate-500 flex items-center justify-between">
                <span>Próxima renovação:</span>
                <span className="font-semibold text-slate-700">
                  {new Date(sub.nextExpectedAt).toLocaleDateString('pt-BR')}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
