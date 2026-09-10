'use client';

import { useState } from 'react';
import { 
  UploadCloud, 
  FileSpreadsheet, 
  CheckCircle, 
  AlertTriangle, 
  ArrowRight, 
  Sparkles, 
  HelpCircle,
  Database,
  Undo2
} from 'lucide-react';

interface MappingItem {
  sourceColumn: string;
  targetField: string;
  confidence: number;
  sampleValue: string;
  confirmed: boolean;
}

export default function ImportsPage() {
  const [stage, setStage] = useState<'upload' | 'mapping' | 'complete'>('upload');
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  const [mappings, setMappings] = useState<MappingItem[]>([
    {
      sourceColumn: 'Data',
      targetField: 'transaction.occurred_at',
      confidence: 0.96,
      sampleValue: '12/08/2026',
      confirmed: true,
    },
    {
      sourceColumn: 'Descrição',
      targetField: 'transaction.description',
      confidence: 0.98,
      sampleValue: 'Supermercado Pão de Açúcar',
      confirmed: true,
    },
    {
      sourceColumn: 'Valor Pg.',
      targetField: 'transaction.amount_minor',
      confidence: 0.95,
      sampleValue: 'R$ 283,20',
      confirmed: true,
    },
    {
      sourceColumn: 'Categoria',
      targetField: 'transaction.category_id',
      confidence: 0.91,
      sampleValue: 'Alimentação',
      confirmed: true,
    },
    {
      sourceColumn: 'Cartão M.',
      targetField: 'transaction.account_id',
      confidence: 0.78,
      sampleValue: 'Nubank Platinum',
      confirmed: false, // Requires user review!
    },
  ]);

  const handleSimulateUpload = () => {
    setSelectedFile('Planilha_Financeira_Familia_2026.xlsx');
    setStage('mapping');
  };

  const toggleConfirm = (index: number) => {
    const updated = [...mappings];
    updated[index].confirmed = !updated[index].confirmed;
    setMappings(updated);
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-slate-900 tracking-tight">Centro de Importação & Migração Excel</h2>
        <p className="text-sm text-slate-500">
          Migre sua planilha familiar histórica ou extratos bancários (XLSX, CSV, OFX) para o Razão de Partidas Dobradas.
        </p>
      </div>

      {stage === 'upload' && (
        <div className="bg-white rounded-xl border border-slate-200 p-8 shadow-xs">
          <div className="max-w-xl mx-auto text-center space-y-4">
            <div className="mx-auto h-16 w-16 rounded-full bg-emerald-50 text-emerald-600 flex items-center justify-center">
              <UploadCloud className="h-8 w-8" />
            </div>

            <div>
              <h3 className="text-lg font-semibold text-slate-900">Arraste sua planilha ou extrato bancário</h3>
              <p className="text-xs text-slate-500 mt-1">
                Formatos aceitos: Excel (.xlsx, .xls), Extrato bancário (.ofx) ou Valores separados por vírgula (.csv)
              </p>
            </div>

            <div className="pt-2">
              <button
                onClick={handleSimulateUpload}
                className="px-5 py-2.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold shadow-xs transition-all cursor-pointer"
              >
                Selecionar Arquivo da Família
              </button>
            </div>

            {/* Quick Demo Preload */}
            <div className="pt-6 border-t border-slate-100">
              <p className="text-xs text-slate-400 font-medium mb-2">Ou experimente com a planilha de teste sintética:</p>
              <button
                onClick={handleSimulateUpload}
                className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-slate-200 bg-slate-50 hover:bg-slate-100 text-xs font-medium text-slate-700 transition-colors cursor-pointer"
              >
                <FileSpreadsheet className="h-4 w-4 text-emerald-600" />
                Carregar "Planilha_Financeira_Familia_2026.xlsx" (182 linhas)
              </button>
            </div>
          </div>
        </div>
      )}

      {stage === 'mapping' && (
        <div className="space-y-6">
          {/* File Analysis Summary */}
          <div className="bg-white rounded-xl border border-slate-200 p-6 shadow-xs flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="p-3 bg-emerald-50 text-emerald-700 rounded-lg">
                <FileSpreadsheet className="h-6 w-6" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900">{selectedFile}</h3>
                <p className="text-xs text-slate-500">
                  Aba detectada: <span className="font-semibold text-slate-700">Gastos 2026</span> • 182 linhas • Moeda: BRL (R$) • Data: DD/MM/AAAA
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <span className="text-xs bg-emerald-100 text-emerald-800 px-2.5 py-1 rounded-full font-semibold flex items-center gap-1">
                <Sparkles className="h-3.5 w-3.5" />
                Mapeamento IA Concluído
              </span>
            </div>
          </div>

          {/* AI Column Mapping Review */}
          <div className="bg-white rounded-xl border border-slate-200 shadow-xs overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100">
              <h4 className="font-semibold text-slate-900 text-sm">Revisão de Mapeamento de Colunas</h4>
              <p className="text-xs text-slate-500">
                A IA analisou as colunas da planilha e propôs o mapeamento para o razão contábil. Confirme os campos com menor confiança.
              </p>
            </div>

            <div className="divide-y divide-slate-100">
              {mappings.map((m, idx) => (
                <div key={m.sourceColumn} className="px-6 py-4 flex items-center justify-between hover:bg-slate-50/60">
                  <div className="flex items-center gap-4">
                    <div className="w-36 font-mono text-sm font-semibold text-slate-900">
                      "{m.sourceColumn}"
                    </div>
                    <ArrowRight className="h-4 w-4 text-slate-400" />
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-xs font-semibold bg-slate-100 text-slate-800 px-2 py-0.5 rounded">
                          {m.targetField}
                        </span>
                        <span className={`text-[11px] font-semibold px-2 py-0.5 rounded-full ${
                          m.confidence >= 0.90
                            ? 'bg-emerald-50 text-emerald-700'
                            : 'bg-amber-50 text-amber-700 border border-amber-200'
                        }`}>
                          {(m.confidence * 100).toFixed(0)}% Confiança
                        </span>
                      </div>
                      <p className="text-xs text-slate-500 mt-1">
                        Exemplo encontrado na planilha: <span className="font-medium text-slate-700">"{m.sampleValue}"</span>
                      </p>
                    </div>
                  </div>

                  <div>
                    {m.confidence < 0.85 && !m.confirmed ? (
                      <button
                        onClick={() => toggleConfirm(idx)}
                        className="px-3 py-1.5 bg-amber-500 hover:bg-amber-600 text-white rounded text-xs font-semibold shadow-xs flex items-center gap-1 cursor-pointer"
                      >
                        <AlertTriangle className="h-3.5 w-3.5" />
                        Confirmar Coluna
                      </button>
                    ) : (
                      <span className="text-xs text-emerald-700 font-semibold flex items-center gap-1">
                        <CheckCircle className="h-4 w-4" />
                        Validado
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>

            {/* Reconciliation & Duplicates Alert */}
            <div className="bg-slate-50 p-4 border-t border-slate-200 flex items-center justify-between text-xs text-slate-600">
              <div className="flex items-center gap-2">
                <CheckCircle className="h-4 w-4 text-emerald-600" />
                <span>180 linhas prontas para importação imediata.</span>
                <span className="text-slate-400">|</span>
                <span className="text-amber-700 font-medium">2 possíveis duplicatas identificadas e isoladas para revisão.</span>
              </div>

              <div className="flex items-center gap-3">
                <button
                  onClick={() => setStage('upload')}
                  className="px-3 py-1.5 text-slate-600 hover:text-slate-900 font-medium cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  onClick={() => setStage('complete')}
                  className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-semibold shadow-xs flex items-center gap-1.5 cursor-pointer"
                >
                  <Database className="h-4 w-4" />
                  Efetivar Lançamento no Razão
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {stage === 'complete' && (
        <div className="bg-white rounded-xl border border-slate-200 p-8 shadow-xs text-center space-y-4">
          <div className="mx-auto h-14 w-14 rounded-full bg-emerald-100 text-emerald-700 flex items-center justify-center">
            <CheckCircle className="h-7 w-7" />
          </div>

          <div>
            <h3 className="text-xl font-bold text-slate-900">Importação Concluída com Sucesso!</h3>
            <p className="text-sm text-slate-500 mt-1">
              180 transações foram lançadas no razão contábil com partidas dobradas equilibradas.
            </p>
          </div>

          <div className="max-w-md mx-auto bg-slate-50 rounded-lg p-4 text-left text-xs space-y-2 border border-slate-100">
            <div className="flex justify-between">
              <span className="text-slate-500">Total de Linhas Processadas:</span>
              <span className="font-semibold text-slate-900">182</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-500">Lançamentos Criados:</span>
              <span className="font-semibold text-emerald-700">180 (100% equilibrados)</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-500">Duplicatas Prevenidas:</span>
              <span className="font-semibold text-amber-700">2</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-500">Linhagem de Origem:</span>
              <span className="font-mono text-slate-700">import_xlsx#batch_20260910</span>
            </div>
          </div>

          <div className="pt-4 flex items-center justify-center gap-3">
            <button
              onClick={() => setStage('upload')}
              className="px-4 py-2 border border-slate-200 rounded-lg text-xs font-semibold text-slate-700 hover:bg-slate-50 flex items-center gap-1.5 cursor-pointer"
            >
              <Undo2 className="h-3.5 w-3.5" />
              Reverter Importação (Rollback)
            </button>
            <a
              href="/transactions"
              className="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-semibold shadow-xs"
            >
              Ver Transações Importadas →
            </a>
          </div>
        </div>
      )}
    </div>
  );
}
