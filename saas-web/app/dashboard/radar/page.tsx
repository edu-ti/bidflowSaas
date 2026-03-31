'use client';

import { RadarCard, type RadarEdital } from '@/components/RadarCard';
import { useFetch } from '@/hooks/useFetch';

export default function RadarPage() {
  const { data, loading, error } = useFetch<{ editais: RadarEdital[] }>('/api/v1/radar');

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-100">Radar Completo</h1>
        <p className="text-gray-400 mt-1">Todos os editais capturados pela IA</p>
      </div>

      {loading ? (
        <div className="grid md:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map((i) => (
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
      ) : error ? (
        <div className="p-4 bg-red-900/20 border border-red-500/50 rounded-xl text-red-200">
          Erro ao carregar radar: {error}
        </div>
      ) : data?.editais?.length === 0 ? (
        <div className="p-8 text-center text-gray-400 bg-gray-800 rounded-2xl border border-gray-700 border-dashed">
          Nenhum edital no radar.
        </div>
      ) : (
        <div className="grid md:grid-cols-3 gap-6">
          {data?.editais?.map((edital) => (
            <RadarCard key={edital.id} edital={edital} />
          ))}
        </div>
      )}
    </div>
  );
}
