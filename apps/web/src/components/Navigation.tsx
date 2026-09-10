'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { 
  LayoutDashboard, 
  Receipt, 
  CreditCard, 
  UploadCloud, 
  Sparkles, 
  Users,
  LineChart 
} from 'lucide-react';

const navItems = [
  { href: '/', label: 'Visão Geral', icon: LayoutDashboard },
  { href: '/transactions', label: 'Transações', icon: Receipt },
  { href: '/accounts', label: 'Contas & Cartões', icon: CreditCard },
  { href: '/imports', label: 'Centro de Importação', icon: UploadCloud },
  { href: '/forecast', label: 'Projeções & Cenários', icon: LineChart },
];

export function Navigation() {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-white border-r border-slate-200 min-h-screen p-4 flex flex-col justify-between">
      <div>
        {/* Brand */}
        <div className="flex items-center gap-2 mb-6 px-2">
          <div className="h-8 w-8 rounded-lg bg-emerald-600 text-white flex items-center justify-center font-bold text-lg">
            F
          </div>
          <div>
            <h1 className="font-semibold text-slate-900 leading-tight">Finora</h1>
            <p className="text-xs text-slate-500 font-medium">Household Intelligence</p>
          </div>
        </div>

        {/* Household Switcher Badge */}
        <div className="bg-slate-50 border border-slate-200 rounded-lg p-3 mb-6 flex items-center gap-2.5">
          <div className="h-7 w-7 rounded-full bg-slate-200 flex items-center justify-center text-slate-600">
            <Users className="h-4 w-4" />
          </div>
          <div className="overflow-hidden">
            <p className="text-xs font-semibold text-slate-900 truncate">Família Silva</p>
            <p className="text-[11px] text-emerald-700 font-medium">BRL (R$) • Principal</p>
          </div>
        </div>

        {/* Nav Links */}
        <nav className="space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-emerald-50 text-emerald-800'
                    : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                }`}
              >
                <Icon className={`h-4 w-4 ${isActive ? 'text-emerald-700' : 'text-slate-400'}`} />
                {item.label}
              </Link>
            );
          })}
        </nav>
      </div>

      {/* Footer Info */}
      <div className="pt-4 border-t border-slate-100">
        <div className="flex items-center gap-2 px-2 py-1 text-xs text-slate-500">
          <Sparkles className="h-3.5 w-3.5 text-emerald-600" />
          <span>Double-Entry Ledger v1.0</span>
        </div>
      </div>
    </aside>
  );
}
