import React, { useState } from 'react';
import { Badge } from '@/components/Badge';
import { cn } from '@/lib/utils';
import { ArrowRight, Loader2 } from 'lucide-react';
import { followRadar } from '@/lib/api';
import { toast } from 'sonner';

export interface RadarEdital {
  id: string;
  title: string;
  value: number;
  deadline: string; // ISO string
  priority_label: 'HIGH' | 'MEDIUM' | 'LOW';
  win_probability: number; // 0-1
  similarity: number; // 0-1
  why_relevant: string;
}

export const RadarCard: React.FC<{ edital: RadarEdital }> = ({ edital }) => {
  const [isFollowing, setIsFollowing] = useState(false);

  const handleFollow = async () => {
    setIsFollowing(true);
    try {
      await followRadar(edital.id);
      toast.success('Edital acompanhado', {
        description: `O edital ${edital.title.slice(0, 30)}... foi adicionado ao seu CRM.`,
      });
    } catch (err: any) {
      toast.error('Falha ao seguir', {
        description: err.response?.data?.message || err.message || 'Erro desconhecido',
      });
    } finally {
      setIsFollowing(false);
    }
  };

  const priorityColor =
    edital.priority_label === 'HIGH'
      ? 'bg-red-600'
      : edital.priority_label === 'MEDIUM'
      ? 'bg-amber-500'
      : 'bg-gray-600';

  return (
    <div
      className={cn(
        'bg-gray-800 rounded-2xl shadow-soft p-6 hover:shadow-lg hover:scale-[1.02] transition transform duration-200',
        'flex flex-col justify-between h-full'
      )}
    >
      <div>
        <h3 className="text-xl font-semibold text-emerald-400 mb-2 line-clamp-2">
          {edital.title}
        </h3>
        <p className="text-sm text-gray-300 mb-2">
          Valor: R$ {edital.value.toLocaleString('pt-BR')}
        </p>
        <p className="text-sm text-gray-300 mb-2">
          Prazo: {new Date(edital.deadline).toLocaleDateString('pt-BR')}
        </p>
        <Badge variant={edital.priority_label.toLowerCase() as any} className="mb-2">
          {edital.priority_label}
        </Badge>
        <div className="mt-4">
          <div className="text-sm text-gray-400 mb-1">
            Probabilidade de vitória: {(edital.win_probability * 100).toFixed(1)}%
          </div>
          <div className="w-full bg-gray-700 rounded-full h-2 mb-2">
            <div
              className="bg-emerald-500 h-2 rounded-full"
              style={{ width: `${edital.win_probability * 100}%` }}
            />
          </div>
          <div className="text-sm text-gray-400 mb-1">
            Similaridade: {(edital.similarity * 100).toFixed(1)}%
          </div>
          <div className="w-full bg-gray-700 rounded-full h-2 mb-2">
            <div
              className="bg-blue-500 h-2 rounded-full"
              style={{ width: `${edital.similarity * 100}%` }}
            />
          </div>
        </div>
        <p className="text-sm text-gray-300 mt-3 line-clamp-3">
          {edital.why_relevant}
        </p>
      </div>
      <div className="mt-4 flex justify-between items-center">
        <a
          href={`/dashboard/radar/${edital.id}`}
          className="text-emerald-500 hover:underline flex items-center"
        >
          Ver análise completa <ArrowRight className="ml-1" size={14} />
        </a>
        <button 
          onClick={handleFollow}
          disabled={isFollowing}
          className="bg-gray-700 disabled:opacity-50 hover:bg-gray-600 text-gray-200 px-3 py-1 flex items-center gap-2 rounded-lg text-sm transition-colors"
        >
          {isFollowing ? <Loader2 size={14} className="animate-spin" /> : null}
          {isFollowing ? 'Processando...' : 'Seguir edital'}
        </button>
      </div>
    </div>
  );
};
