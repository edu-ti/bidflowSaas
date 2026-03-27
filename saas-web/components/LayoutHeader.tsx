'use client';

import { Bell, User } from 'lucide-react';

export const LayoutHeader: React.FC = () => {
  return (
    <header className="flex items-center justify-between px-6 py-4 bg-gray-850 border-b border-gray-700">
      <div>
        <h2 className="text-lg font-semibold text-gray-100">BidFlow</h2>
        <p className="text-xs text-gray-400">AI Copilot para Licitações</p>
      </div>
      <div className="flex items-center gap-4">
        <button className="relative text-gray-400 hover:text-emerald-400 transition-colors" aria-label="Notificações">
          <Bell size={20} />
          <span className="absolute -top-1 -right-1 w-2 h-2 bg-emerald-500 rounded-full" />
        </button>
        <button className="flex items-center gap-2 text-gray-400 hover:text-gray-100 transition-colors" aria-label="Perfil">
          <div className="w-8 h-8 rounded-full bg-emerald-700 flex items-center justify-center">
            <User size={16} className="text-white" />
          </div>
        </button>
      </div>
    </header>
  );
};
