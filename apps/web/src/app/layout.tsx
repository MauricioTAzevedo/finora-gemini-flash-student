import type { Metadata } from 'next';
import './globals.css';
import { Navigation } from '@/components/Navigation';

export const metadata: Metadata = {
  title: 'Finora — Inteligência Financeira Familiar',
  description: 'Substituto inteligente de planilhas financeiras para famílias.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR">
      <body className="flex min-h-screen bg-slate-50 text-slate-900 antialiased">
        <Navigation />
        <main className="flex-1 p-8 overflow-y-auto max-w-7xl">
          {children}
        </main>
      </body>
    </html>
  );
}
