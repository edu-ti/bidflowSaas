'use client';

import { useFetch } from '@/hooks/useFetch';
import { Card } from '@/components/Card';
import { Bell, CheckSquare, Loader2 } from 'lucide-react';
import { readNotification } from '@/lib/api';
import { useState } from 'react';
import { toast } from 'sonner';

export interface Notification {
  id: string;
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
}

export default function NotificationsPage() {
  const { data, loading, error } = useFetch<Notification[]>('/api/v1/notifications');
  const [reading, setReading] = useState<Record<string, boolean>>({});
  const [localRead, setLocalRead] = useState<Record<string, boolean>>({});

  const handleRead = async (id: string) => {
    setReading(prev => ({ ...prev, [id]: true }));
    try {
      await readNotification(id);
      setLocalRead(prev => ({ ...prev, [id]: true }));
      toast.success('Marcada como lida');
    } catch (err) {
      toast.error('Erro ao marcar notificação');
    } finally {
      setReading(prev => ({ ...prev, [id]: false }));
    }
  };

  const notifications = data || [];

  return (
    <div className="space-y-8 max-w-4xl">
      <div>
        <h1 className="text-3xl font-bold text-gray-100 flex items-center gap-3">
          <Bell className="text-emerald-400" size={28} />
          Notificações
        </h1>
        <p className="text-gray-400 mt-1">Alertas e atualizações do sistema</p>
      </div>

      {loading ? (
        <div className="space-y-4">
          {[1, 2, 3, 4].map((i) => (
             <div key={i} className="bg-gray-800 rounded-xl p-5 h-24 animate-pulse">
                <div className="h-5 bg-gray-700 rounded w-1/4 mb-3" />
                <div className="h-4 bg-gray-700 rounded w-1/2" />
             </div>
          ))}
        </div>
      ) : error ? (
        <div className="p-4 bg-red-900/20 border border-red-500/50 rounded-xl text-red-200">
          Erro ao buscar notificações: {error}
        </div>
      ) : notifications.length === 0 ? (
        <div className="p-8 text-center text-gray-400 bg-gray-800 rounded-2xl border border-gray-700 border-dashed">
          Você não tem notificações no momento.
        </div>
      ) : (
        <div className="space-y-4">
          {notifications.map((notif) => {
            const isRead = notif.read || localRead[notif.id];
            return (
              <Card key={notif.id} className={`p-5 transition-colors ${isRead ? 'opacity-60 bg-gray-900 border-transparent shadow-none' : 'border-emerald-500/20'}`}>
                <div className="flex justify-between items-start gap-4">
                  <div>
                    <h3 className={`font-medium ${isRead ? 'text-gray-400' : 'text-gray-200'}`}>{notif.title}</h3>
                    <p className={`text-sm mt-1 leading-relaxed ${isRead ? 'text-gray-500' : 'text-gray-400'}`}>
                      {notif.message}
                    </p>
                    <time className="text-xs text-gray-600 block mt-3 font-medium">
                      {new Date(notif.timestamp).toLocaleString('pt-BR')}
                    </time>
                  </div>
                  {!isRead && (
                    <button
                      onClick={() => handleRead(notif.id)}
                      disabled={reading[notif.id]}
                      className="p-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded flex-shrink-0 transition-colors"
                      title="Marcar como lida"
                    >
                      {reading[notif.id] ? <Loader2 size={16} className="animate-spin" /> : <CheckSquare size={16} />}
                    </button>
                  )}
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
