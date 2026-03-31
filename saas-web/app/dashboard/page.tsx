'use client';

import { RadarCard, type RadarEdital } from '@/components/RadarCard';
import { Card } from '@/components/Card';
import { useFetch } from '@/hooks/useFetch';

export default function DashboardPage() {
  const { data: radarData, loading: radarLoading, error: radarError } = useFetch<{ editais: RadarEdital[] }>('/api/v1/radar');
  const { data: insightsData, loading: insightsLoading } = useFetch<any>('/api/v1/ai/insights');

  // Show 3 items max on dashboard for Radar
  const editaisToShow = radarData?.editais?.slice(0, 3) || [];

  const stats = [
    { label: 'Editais em Radar', value: radarData?.editais?.length || 0, color: 'text-emerald-400' },
    { label: 'Score Médio', value: insightsData?.average_score ? `${(insightsData.average_score * 100).toFixed(0)}/100` : '--', color: 'text-amber-400' },
    { label: 'Win Probability Médio', value: insightsData?.average_win_probability ? `${(insightsData.average_win_probability * 100).toFixed(1)}%` : '--', color: 'text-blue-400' },
    { label: 'Valor Mapeado', value: radarData?.editais ? `R$ ${(radarData.editais.reduce((acc, curr) => acc + curr.value, 0) / 1000).toLocaleString('pt-BR', { maximumFractionDigits: 0 })}k` : '--', color: 'text-purple-400' },
  ];

  return (
    <div className="space-y-8">
      {/* Page header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-100">Visão Geral</h1>
        <p className="text-gray-400 mt-1">Seus editais mais relevantes de hoje</p>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {stats.map((s) => (
          <Card key={s.label} className="p-6 text-center">
            {insightsLoading || radarLoading ? (
              <div className="h-9 bg-gray-700 animate-pulse rounded max-w-24 mx-auto mb-1"></div>
            ) : (
              <p className={`text-3xl font-bold ${s.color}`}>{s.value}</p>
            )}
            <p className="text-sm text-gray-400 mt-1">{s.label}</p>
          </Card>
        ))}
      </div>

      {/* Radar section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-200 mb-4">Radar de Oportunidades (Recentes)</h2>
        {(radarLoading) ? (
          <div className="grid md:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-gray-800 rounded-2xl p-6 h-[340px] animate-pulse">
                <div className="h-6 bg-gray-700 rounded w-3/4 mb-4" />
                <div className="h-4 bg-gray-700 rounded w-1/2 mb-2" />
                <div className="h-4 bg-gray-700 rounded w-1/3 mb-4" />
                <div className="space-y-2 mt-8">
                   <div className="h-4 bg-gray-700 rounded w-full" />
                   <div className="h-4 bg-gray-700 rounded w-5/6" />
                </div>
              </div>
            ))}
          </div>
        ) : radarError ? (
          <div className="p-4 bg-red-900/20 border border-red-500/50 rounded-xl text-red-200">
            Erro ao carregar radar: {radarError}
          </div>
        ) : editaisToShow.length === 0 ? (
          <div className="p-8 text-center text-gray-400 bg-gray-800 rounded-2xl border border-gray-700 border-dashed">
            Nenhum edital encontrado no radar hoje.
          </div>
        ) : (
          <div className="grid md:grid-cols-3 gap-6">
            {editaisToShow.map((edital) => (
              <RadarCard key={edital.id} edital={edital} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
