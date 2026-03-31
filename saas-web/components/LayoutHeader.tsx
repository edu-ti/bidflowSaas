'use client';

import { Bell, User } from 'lucide-react';
import Link from 'next/link';

export const LayoutHeader: React.FC = () => {
  return (
    <header className="flex items-center justify-between px-6 py-4 bg-gray-850 border-b border-gray-700">
      <div>
        <h2 className="text-lg font-semibold text-gray-100">BidFlow</h2>
        <p className="text-xs text-gray-400">AI Copilot para Licitações</p>
      </div>
      <div className="flex items-center gap-4">
        <Link href="/dashboard/notifications" className="relative text-gray-400 hover:text-emerald-400 transition-colors" aria-label="Notificações">
          <Bell size={20} />
          <span className="absolute -top-1 -right-1 w-2 h-2 bg-emerald-500 rounded-full" />
        </Link>
        <div className="relative group">
          <button className="flex items-center gap-2 text-gray-400 hover:text-gray-100 transition-colors" aria-label="Perfil">
            <div className="w-8 h-8 rounded-full bg-emerald-700 flex items-center justify-center">
              <User size={16} className="text-white" />
            </div>
          </button>
          <div className="absolute right-0 mt-2 w-48 bg-gray-800 rounded-xl border border-gray-700 shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
            <div className="p-2">
              <Link href="/dashboard/billing" className="block px-4 py-2 text-sm text-gray-300 hover:bg-gray-700 rounded-lg">Meu Plano</Link>
              <button 
                onClick={() => {
                  localStorage.removeItem('token');
                  window.location.href = '/login';
                }}
                className="block w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-gray-700 rounded-lg"
              >
                Sair da conta
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};
