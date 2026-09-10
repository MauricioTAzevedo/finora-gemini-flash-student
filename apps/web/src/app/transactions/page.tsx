'use client';

import { useState } from 'react';
import { formatDateBR, formatBRL } from '@/lib/format';
import { Plus, Search, Filter, ArrowDownLeft, ArrowUpRight, CheckCircle2 } from 'lucide-react';

interface TransactionRow {
  id: string;
  date: string;
  description: string;
  category: string;
  account: string;
  amountMinor: number;
  type: 'expense' | 'income' | 'transfer';
  source: string;
}

const INITIAL_ROWS: TransactionRow[] = [
  {
    id: '1',
    date: new Date(Date.now() - 86400000 * 1).toISOString(),
    description: 'Conta de Energia Elétrica (CPFL)',
    category: 'Moradia & Energia',
    account: 'Nubank Conta Principal',
    amountMinor: 48700,
    type: 'expense',
    source: 'manual',
  },
  {
    id: '2',
    date: new Date(Date.now() - 86400000 * 2).toISOString(),
    description: 'Supermercado Horizonte',
    category: 'Alimentação & Supermercado',
    account: 'Nubank Conta Principal',
    amountMinor: 34050,
    type: 'expense',
    source: 'import_xlsx',
  },
  {
    id: '3',
    date: new Date(Date.now() - 86400000 * 3).toISOString(),
    description: 'Posto Central - Gasolina',
    category: 'Transporte & Combustível',
    account: 'Nubank Platinum (Cartão)',
    amountMinor: 18000,
    type: 'expense',
    source: 'import_ofx',
  },
  {
    id: '4',
    date: new Date(Date.now() - 86400000 * 5).toISOString(),
    description: 'Drogasil Loja 92',
    category: 'Saúde & Farmácia',
    account: 'Nubank Conta Principal',
    amountMinor: 8990,
    type: 'expense',
    source: 'manual',
  },
  {
    id: '5',
    date: new Date(Date.now() - 86400000 * 10).toISOString(),
    description: 'Salário Mensal',
    category: 'Salário & Renda',
    account: 'Nubank Conta Principal',
    amountMinor: 850000,
    type: 'income',
    source: 'manual',
  },
];

export default function TransactionsPage() {
  const [rows, setRows] = useState<TransactionRow[]>(INITIAL_ROWS);
  const [search, setSearch] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Form State for new expense
  const [description, setDescription] = useState('');
  const [amountStr, setAmountStr] = useState('');
  const [category, setCategory] = useState('Alimentação & Supermercado');
  const [account, setAccount] = useState('Nubank Conta Principal');

  const filteredRows = rows.filter(r => 
    r.description.toLowerCase().includes(search.toLowerCase()) ||
    r.category.toLowerCase().includes(search.toLowerCase())
  );

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    const cleanAmount = parseFloat(amountStr.replace(',', '.'));
    if (isNaN(cleanAmount) || cleanAmount <= 0) return;

    const newTx: TransactionRow = {
      id: String(Date.now()),
      date: new Date().toISOString(),
      description: description || 'Lançamento Manual',
      category,
      account,
      amountMinor: Math.round(cleanAmount * 100),
      type: 'expense',
      source: 'manual',
    };

    setRows([newTx, ...rows]);
    setIsModalOpen(false);
    setDescription('');
    setAmountStr('');
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Transações do Razão</h2>
          <p className="text-sm text-slate-500">Extrato detalhado com linhagem e partidas dobradas equilibradas.</p>
        </div>
        <button
          onClick={() => setIsModalOpen(true)}
          className="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg bg-emerald-600 text-white text-sm font-medium hover:bg-emerald-700 transition-colors shadow-xs"
        >
          <Plus className="h-4 w-4" />
          Nova Transação
        </button>
      </div>

      {/* Filter / Search Bar */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="h-4 w-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Buscar por descrição, estabelecimento ou categoria..."
            className="w-full pl-9 pr-4 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500"
          />
        </div>
        <button className="flex items-center gap-2 px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm text-slate-700 hover:bg-slate-50">
          <Filter className="h-4 w-4 text-slate-500" />
          Filtros
        </button>
      </div>

      {/* Dense Spreadsheet-Inspired Table */}
      <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-semibold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="px-5 py-3">Data</th>
                <th className="px-5 py-3">Descrição</th>
                <th className="px-5 py-3">Categoria</th>
                <th className="px-5 py-3">Conta / Cartão</th>
                <th className="px-5 py-3">Origem</th>
                <th className="px-5 py-3 text-right">Valor (BRL)</th>
                <th className="px-5 py-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filteredRows.map((tx) => {
                const isIncome = tx.type === 'income';
                return (
                  <tr key={tx.id} className="hover:bg-slate-50/70 transition-colors">
                    <td className="px-5 py-3 whitespace-nowrap text-slate-600 font-mono text-xs">
                      {formatDateBR(tx.date)}
                    </td>
                    <td className="px-5 py-3 text-slate-900 font-medium">
                      {tx.description}
                    </td>
                    <td className="px-5 py-3 whitespace-nowrap">
                      <span className="text-xs bg-slate-100 text-slate-800 px-2 py-0.5 rounded font-medium">
                        {tx.category}
                      </span>
                    </td>
                    <td className="px-5 py-3 whitespace-nowrap text-xs text-slate-600">
                      {tx.account}
                    </td>
                    <td className="px-5 py-3 whitespace-nowrap">
                      <span className="text-[11px] font-mono bg-slate-100 text-slate-600 px-1.5 py-0.5 rounded">
                        {tx.source}
                      </span>
                    </td>
                    <td className={`px-5 py-3 whitespace-nowrap text-right font-semibold font-mono text-sm ${
                      isIncome ? 'text-emerald-700' : 'text-slate-900'
                    }`}>
                      {isIncome ? '+' : '-'}{formatBRL(tx.amountMinor)}
                    </td>
                    <td className="px-5 py-3 whitespace-nowrap text-center">
                      <span className="inline-flex items-center gap-1 text-xs text-emerald-700 font-medium">
                        <CheckCircle2 className="h-3.5 w-3.5" />
                        Equilibrado
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* New Transaction Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl border border-slate-200 shadow-xl max-w-md w-full p-6">
            <h3 className="text-lg font-bold text-slate-900 mb-4">Lançar Nova Transação no Razão</h3>
            <form onSubmit={handleCreate} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Descrição
                </label>
                <input
                  type="text"
                  required
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Ex: Padaria Central, Combustível, Farmácia"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Valor em Reais (R$)
                </label>
                <input
                  type="text"
                  required
                  value={amountStr}
                  onChange={(e) => setAmountStr(e.target.value)}
                  placeholder="Ex: 45,50"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm font-mono focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Categoria
                </label>
                <select
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                >
                  <option>Alimentação & Supermercado</option>
                  <option>Moradia & Energia</option>
                  <option>Transporte & Combustível</option>
                  <option>Saúde & Farmácia</option>
                  <option>Assinaturas & Lazer</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">
                  Conta de Origem
                </label>
                <select
                  value={account}
                  onChange={(e) => setAccount(e.target.value)}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:ring-2 focus:ring-emerald-500 focus:outline-none"
                >
                  <option>Nubank Conta Principal</option>
                  <option>Nubank Platinum (Cartão)</option>
                  <option>Reserva de Emergência</option>
                </select>
              </div>

              <div className="flex items-center justify-end gap-2 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 text-sm text-slate-600 hover:text-slate-900 font-medium"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white text-sm font-medium rounded-lg shadow-xs"
                >
                  Registrar Lançamento
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
