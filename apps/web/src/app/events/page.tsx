'use client';

import { useState, useEffect } from 'react';
import { 
  Activity, 
  CheckCircle, 
  Clock, 
  AlertCircle, 
  Copy, 
  ChevronDown, 
  ChevronRight, 
  RefreshCw,
  GitCommit,
  Layers,
  Database
} from 'lucide-react';

const API_BASE = 'http://localhost:8080/api/v1';
const HEADERS = {
  'Content-Type': 'application/json',
  'X-Household-Id': 'b0000000-0000-0000-0000-000000000001',
};

interface OutboxEvent {
  id: string;
  aggregateType: string;
  aggregateId: string;
  eventType: string;
  correlationId: string;
  causationId: string;
  householdId: string;
  payload: any;
  createdAt: string;
  status: string;
  retryCount: number;
  publishedAt?: string;
}

export default function EventExplorerPage() {
  const [events, setEvents] = useState<OutboxEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [expandedIds, setExpandedIds] = useState<Record<string, boolean>>({});

  const loadEvents = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/events`, { headers: HEADERS });
      if (res.ok) {
        const data = await res.json();
        setEvents(data.events || []);
      }
    } catch (err) {
      console.error('Failed to load outbox events:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadEvents();
  }, []);

  const toggleExpand = (id: string) => {
    setExpandedIds(prev => ({ ...prev, [id]: !prev[id] }));
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'PUBLISHED':
        return (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-emerald-100 text-emerald-800">
            <CheckCircle className="h-3 w-3 text-emerald-600" />
            PUBLISHED
          </span>
        );
      case 'PENDING':
        return (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-100 text-amber-800">
            <Clock className="h-3 w-3 text-amber-600" />
            PENDING
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[11px] font-semibold bg-rose-100 text-rose-800">
            <AlertCircle className="h-3 w-3 text-rose-600" />
            {status}
          </span>
        );
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <span className="px-2 py-0.5 text-xs font-semibold rounded bg-indigo-100 text-indigo-800 flex items-center gap-1">
              <Activity className="h-3 w-3" />
              Milestone 5: Transactional Outbox & Event Architecture
            </span>
          </div>
          <h1 className="text-2xl font-bold text-slate-900">Event Explorer & Auditoria de Causalidade</h1>
          <p className="text-sm text-slate-500">
            Trilha cronológica imutável de eventos de domínio gravados atomicamente via Transactional Outbox.
          </p>
        </div>

        <button
          onClick={loadEvents}
          className="px-3.5 py-2 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 text-xs font-semibold rounded-lg shadow-xs flex items-center gap-2 cursor-pointer transition-colors"
        >
          <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin text-indigo-600' : 'text-slate-500'}`} />
          Atualizar Eventos
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl border border-slate-200 p-4 shadow-xs flex items-center gap-3">
          <div className="h-10 w-10 rounded-lg bg-indigo-50 flex items-center justify-center text-indigo-600">
            <Database className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-medium">Total de Eventos</p>
            <p className="text-xl font-bold text-slate-900">{events.length}</p>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-slate-200 p-4 shadow-xs flex items-center gap-3">
          <div className="h-10 w-10 rounded-lg bg-emerald-50 flex items-center justify-center text-emerald-600">
            <CheckCircle className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-medium">Publicados (Relay)</p>
            <p className="text-xl font-bold text-emerald-700">
              {events.filter(e => e.status === 'PUBLISHED').length}
            </p>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-slate-200 p-4 shadow-xs flex items-center gap-3">
          <div className="h-10 w-10 rounded-lg bg-purple-50 flex items-center justify-center text-purple-600">
            <Layers className="h-5 w-5" />
          </div>
          <div>
            <p className="text-xs text-slate-500 font-medium">Aggregates Rastreáveis</p>
            <p className="text-xl font-bold text-purple-700">
              {new Set(events.map(e => e.aggregateType)).size}
            </p>
          </div>
        </div>
      </div>

      {/* Events Table / Feed */}
      <div className="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
        <div className="px-5 py-4 border-b border-slate-200 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <GitCommit className="h-4 w-4 text-indigo-600" />
            <h2 className="text-sm font-bold text-slate-900">Event Stream (Mais recentes primeiro)</h2>
          </div>
          <span className="text-xs text-slate-400">Idempotência & Correlation ID garantidos</span>
        </div>

        <div className="divide-y divide-slate-200">
          {events.map((evt) => {
            const isExpanded = expandedIds[evt.id];
            return (
              <div key={evt.id} className="p-4 hover:bg-slate-50/50 transition-colors">
                <div 
                  className="flex flex-col md:flex-row md:items-center justify-between gap-3 cursor-pointer"
                  onClick={() => toggleExpand(evt.id)}
                >
                  <div className="flex items-start gap-3">
                    <button className="text-slate-400 hover:text-slate-600 mt-0.5">
                      {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                    </button>
                    <div>
                      <div className="flex items-center gap-2 mb-1 flex-wrap">
                        <span className="text-xs font-mono font-bold text-slate-900">{evt.eventType}</span>
                        <span className="px-2 py-0.5 rounded text-[10px] font-semibold bg-slate-100 text-slate-600 border border-slate-200">
                          {evt.aggregateType}
                        </span>
                        {getStatusBadge(evt.status)}
                      </div>
                      <div className="flex items-center gap-3 text-xs text-slate-500 font-mono">
                        <span>ID: {evt.id.substring(0, 8)}...</span>
                        <span>Corr: {evt.correlationId}</span>
                        <span>Caus: {evt.causationId}</span>
                      </div>
                    </div>
                  </div>

                  <div className="text-right text-xs text-slate-500 pl-7 md:pl-0">
                    <div>{new Date(evt.createdAt).toLocaleString('pt-BR')}</div>
                    {evt.publishedAt && (
                      <span className="text-[10px] text-emerald-600 font-medium">Relay instantâneo</span>
                    )}
                  </div>
                </div>

                {/* Expanded JSON payload view */}
                {isExpanded && (
                  <div className="mt-3 pl-7 pt-3 border-t border-slate-100">
                    <p className="text-xs font-semibold text-slate-700 mb-1">Payload JSON:</p>
                    <pre className="p-3 bg-slate-900 text-emerald-400 rounded-lg text-xs font-mono overflow-x-auto">
                      {JSON.stringify(evt.payload, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
