'use client';

import { useFetch } from '@/hooks/useFetch';
import { Card } from '@/components/Card';
import { BrainCircuit, TrendingUp, Target, BarChart3, AlertCircle } from 'lucide-react';
import Link from 'next/link';
import { Button } from '@/components/Button';

export default function AIPage() {
  const { data: insightsData, loading, error } = useFetch<any>('/api/v1/ai/insights');

  const isLoading = loading || !insightsData;

  const getPercentage = (val: number) => {
    return val ? `${(val * 100).toFixed(1)}%` : '--';
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold text-gray-100 flex items-center gap-3">
            <BrainCircuit className="text-emerald-400" size={32} />
            Inteligência Artificial
          </h1>
          <p className="text-gray-400 mt-1">Sua performance preditiva e histórico de análises</p>
        </div>
        <Link href="/dashboard/ai/history">
          <Button className="bg-gray-800 text-gray-300 border-gray-700 hover:bg-gray-700">Ver Histórico de Análises</Button>
        </Link>
      </div>

      {error && (
        <div className="p-4 bg-red-900/20 border border-red-500/50 rounded-xl text-red-200">
          Erro ao carregar os dados de IA: {error}
        </div>
      )}

      {/* Cards Row */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="p-6 border border-emerald-500/10 hover:border-emerald-500/30 transition-colors">
          <div className="flex items-center gap-4 mb-4">
            <div className="p-3 bg-emerald-500/10 rounded-lg">
              <TrendingUp className="text-emerald-400" size={24} />
            </div>
            <h3 className="text-lg font-semibold text-gray-200">Taxa de Sucesso</h3>
          </div>
          {isLoading ? (
            <div className="h-10 bg-gray-700 animate-pulse rounded max-w-24"></div>
          ) : (
            <p className="text-4xl font-bold text-emerald-400">{getPercentage(insightsData?.success_rate)}</p>
          )}
          <p className="text-sm text-gray-400 mt-2">Das licitações avaliadas favoravelmente</p>
        </Card>

        <Card className="p-6 border border-blue-500/10 hover:border-blue-500/30 transition-colors">
          <div className="flex items-center gap-4 mb-4">
            <div className="p-3 bg-blue-500/10 rounded-lg">
              <Target className="text-blue-400" size={24} />
            </div>
            <h3 className="text-lg font-semibold text-gray-200">Score Médio</h3>
          </div>
          {isLoading ? (
            <div className="h-10 bg-gray-700 animate-pulse rounded max-w-24"></div>
          ) : (
             <p className="text-4xl font-bold text-blue-400">{insightsData?.average_score ? (insightsData.average_score * 100).toFixed(0) : '--'}/100</p>
          )}
          <p className="text-sm text-gray-400 mt-2">Aderência média aos editais no radar</p>
        </Card>

        <Card className="p-6 border border-purple-500/10 hover:border-purple-500/30 transition-colors">
          <div className="flex items-center gap-4 mb-4">
            <div className="p-3 bg-purple-500/10 rounded-lg">
              <BarChart3 className="text-purple-400" size={24} />
            </div>
            <h3 className="text-lg font-semibold text-gray-200">Confiança IA</h3>
          </div>
           {isLoading ? (
            <div className="h-10 bg-gray-700 animate-pulse rounded max-w-24"></div>
          ) : (
             <p className="text-4xl font-bold text-purple-400">{getPercentage(insightsData?.average_confidence)}</p>
          )}
          <p className="text-sm text-gray-400 mt-2">Nível médio de confiabilidade do modelo</p>
        </Card>
      </div>

      <Card className="p-8 border border-gray-800 bg-gray-800/50 flex flex-col items-center justify-center min-h-64 text-center">
        <AlertCircle className="text-gray-500 mb-4" size={48} />
        <h3 className="text-xl font-medium text-gray-300 mb-2">Gráficos em Desenvolvimento</h3>
        <p className="text-gray-400 max-w-md mx-auto">
          A visualização gráfica das análises, clusterização de similares e linha do tempo de scores em `recharts` será disponibilizada na próxima atualização da plataforma.
        </p>
      </Card>
    </div>
  );
}
