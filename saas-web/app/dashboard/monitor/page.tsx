'use client';

import { useRealtime, type MonitorEvent } from '@/hooks/useRealtime';
import { Card } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { Clock, AlertCircle, Info, AlertTriangle } from 'lucide-react';

export default function MonitorPage() {
  const { events, loading, error } = useRealtime(5000);

  const getPriorityIcon = (severity: string) => {
    switch (severity?.toUpperCase()) {
      case 'HIGH':
        return <AlertTriangle className="text-red-500" size={20} />;
      case 'MEDIUM':
        return <AlertCircle className="text-amber-500" size={20} />;
      default:
        return <Info className="text-blue-500" size={20} />;
    }
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-100 flex items-center gap-3">
          Monitor de Pregão 
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
          </span>
        </h1>
        <p className="text-gray-400 mt-1">Acompanhamento de lances e mensagens em tempo real</p>
      </div>

      {loading && events.length === 0 ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
             <div key={i} className="bg-gray-800 rounded-xl p-4 h-24 animate-pulse">
                <div className="h-5 bg-gray-700 rounded w-1/4 mb-3" />
                <div className="h-4 bg-gray-700 rounded w-1/2" />
             </div>
          ))}
        </div>
      ) : error ? (
        <div className="p-4 bg-red-900/20 border border-red-500/50 rounded-xl text-red-200">
          Erro de conexão: {error}
        </div>
      ) : events.length === 0 ? (
        <div className="p-8 text-center text-gray-400 bg-gray-800 rounded-2xl border border-gray-700 border-dashed flex flex-col items-center justify-center">
          <Clock className="w-12 h-12 mb-4 text-gray-500" />
          <p>Aguardando novos eventos de pregão...</p>
        </div>
      ) : (
        <div className="space-y-4">
          {events.map((event) => (
            <Card key={event.id} className="p-5 flex items-start gap-4 hover:bg-gray-800/80 transition-colors">
              <div className="mt-1">
                {getPriorityIcon(event.severity)}
              </div>
              <div className="flex-1">
                <div className="flex justify-between items-start mb-1">
                  <div className="flex items-center gap-2">
                    <Badge variant={(event.severity?.toLowerCase() || 'low') as any}>
                      {event.type}
                    </Badge>
                    <span className="text-xs text-gray-500 font-medium">
                      {new Date(event.timestamp).toLocaleTimeString('pt-BR')}
                    </span>
                  </div>
                </div>
                <p className="text-gray-200 mt-2 text-sm leading-relaxed">{event.message}</p>
                {event.metadata?.suggested_action && (
                  <div className="mt-3 inline-flex items-center gap-2 px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 text-blue-400 text-xs rounded-lg font-medium">
                    Sugestão da IA: {event.metadata.suggested_action}
                  </div>
                )}
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
