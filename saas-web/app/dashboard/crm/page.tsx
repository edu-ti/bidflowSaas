'use client';

import { Card } from '@/components/Card';
import { Briefcase, MoreHorizontal, Plus } from 'lucide-react';

const stages = [
  { id: '1', name: 'Prospectando', count: 3, color: 'bg-blue-500/20 text-blue-400 border-blue-500/30' },
  { id: '2', name: 'Contato', count: 5, color: 'bg-amber-500/20 text-amber-400 border-amber-500/30' },
  { id: '3', name: 'Proposta', count: 2, color: 'bg-purple-500/20 text-purple-400 border-purple-500/30' },
  { id: '4', name: 'Negociação', count: 1, color: 'bg-orange-500/20 text-orange-400 border-orange-500/30' },
  { id: '5', name: 'Fechado', count: 4, color: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' },
];

export default function CRMPage() {
  return (
    <div className="space-y-8h-full flex flex-col h-full">
      <div className="flex justify-between items-end mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-100 flex items-center gap-3">
            <Briefcase className="text-emerald-400" size={32} />
            CRM Licitações
          </h1>
          <p className="text-gray-400 mt-1">Gerencie suas oportunidades de negócio governamentais</p>
        </div>
        <button className="bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors">
          <Plus size={18} />
          Nova Oportunidade
        </button>
      </div>

      {/* Kanban Board Mock */}
      <div className="flex-1 overflow-x-auto pb-4">
        <div className="flex gap-6 min-w-max h-full">
          {stages.map((stage) => (
            <div key={stage.id} className="w-80 flex flex-col bg-gray-900/50 rounded-xl rounded-t-2xl">
              <div className={`p-4 border-b-2 ${stage.color} rounded-t-xl bg-gray-800/80 sticky top-0 flex justify-between items-center`}>
                <h3 className="font-semibold text-gray-200">{stage.name}</h3>
                <span className="text-xs font-bold leading-none bg-gray-700 px-2 py-1 rounded-full text-gray-300">
                  {stage.count}
                </span>
              </div>
              
              <div className="flex-1 p-4 space-y-4 overflow-y-auto">
                <Card className="p-4 bg-gray-800 border-gray-700 cursor-pointer hover:border-emerald-500/50 hover:bg-gray-750 transition-all">
                   <div className="flex justify-between items-start mb-2">
                     <span className="text-xs text-gray-400 font-mono">#MOCK-001</span>
                     <button className="text-gray-500 hover:text-gray-300"><MoreHorizontal size={16} /></button>
                   </div>
                   <h4 className="font-semibold text-gray-200 text-sm mb-2 line-clamp-2">Exemplo de Licitação de TI (Mock Data)</h4>
                   <div className="flex items-center gap-2">
                     <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
                     <span className="text-xs text-gray-400">R$ 150.000,00</span>
                   </div>
                </Card>
                 <Card className="p-4 bg-gray-800 border-gray-700 border-dashed flex justify-center items-center opacity-50">
                    <span className="text-gray-500 text-sm">Mais itens...</span>
                 </Card>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
