'use client';

import { useFetch } from '@/hooks/useFetch';
import { Card } from '@/components/Card';
import { History, Send, MessageSquare } from 'lucide-react';
import { postFeedback } from '@/lib/api';
import { useState } from 'react';
import { toast } from 'sonner';

export default function AIHistoryPage() {
  const { data: historyData, loading, error } = useFetch<any[]>('/api/v1/ai/history');
  const [feedbackState, setFeedbackState] = useState<Record<string, string>>({});
  const [loadingFeedback, setLoadingFeedback] = useState<Record<string, boolean>>({});

  const handleFeedbackSubmit = async (historyId: string) => {
    const fbText = feedbackState[historyId];
    if (!fbText || fbText.trim() === '') return;

    setLoadingFeedback(prev => ({ ...prev, [historyId]: true }));
    try {
      await postFeedback({ history_id: historyId, feedback: fbText });
      toast.success('Feedback enviado para aprendizado da IA');
      setFeedbackState(prev => ({ ...prev, [historyId]: '' }));
    } catch (err) {
      toast.error('Erro ao enviar feedback');
    } finally {
      setLoadingFeedback(prev => ({ ...prev, [historyId]: false }));
    }
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold text-gray-100 flex items-center gap-3">
          <History className="text-emerald-400" size={32} />
          Histórico de IA
        </h1>
        <p className="text-gray-400 mt-1">Análises passadas e calibração de modelos</p>
      </div>

      <Card className="overflow-hidden bg-gray-900 border-gray-800">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-gray-400">
            <thead className="text-xs uppercase bg-gray-800 text-gray-300 border-b border-gray-700">
              <tr>
                <th scope="col" className="px-6 py-4 rounded-tl-xl">Edital (Id)</th>
                <th scope="col" className="px-6 py-4 text-center">AI Score</th>
                <th scope="col" className="px-6 py-4 text-center">Smart Score</th>
                <th scope="col" className="px-6 py-4 text-center">Confidence</th>
                <th scope="col" className="px-6 py-4 text-right rounded-tr-xl">Ação / Qualificação</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-6 py-12 text-center text-gray-500">
                    <div className="flex justify-center mb-2"><div className="w-6 h-6 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin"></div></div>
                    Carregando histórico...
                  </td>
                </tr>
              ) : error ? (
                <tr>
                   <td colSpan={5} className="px-6 py-6 text-center text-red-400">
                    Erro ao buscar histórico: {error}
                  </td>
                </tr>
              ) : historyData?.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-12 text-center text-gray-500">
                    Nenhum histórico encontrado.
                  </td>
                </tr>
              ) : (
                historyData?.map((item, idx) => (
                  <tr key={item.id || idx} className="border-b border-gray-800 hover:bg-gray-800/40">
                    <td className="px-6 py-4 font-medium text-gray-200 whitespace-nowrap">
                      #{item.edital_id}
                    </td>
                    <td className="px-6 py-4 text-center">
                       <span className="px-2 py-1 bg-gray-800 rounded font-bold text-blue-400">{(item.ai_score * 100).toFixed(0)}</span>
                    </td>
                    <td className="px-6 py-4 text-center">
                       <span className="px-2 py-1 bg-gray-800 rounded font-bold text-amber-400">{(item.smart_score * 100).toFixed(0)}</span>
                    </td>
                    <td className="px-6 py-4 text-center">
                       <span className="px-2 py-1 bg-gray-800 rounded font-bold text-purple-400 font-mono">{(item.confidence * 100).toFixed(1)}%</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <input
                          type="text"
                          value={feedbackState[item.id] || ''}
                          onChange={(e) => setFeedbackState(prev => ({ ...prev, [item.id]: e.target.value }))}
                          placeholder="Correção de ia..."
                          className="bg-gray-800 border border-gray-700 text-xs text-gray-200 rounded px-2 py-1.5 focus:outline-none focus:border-emerald-500 max-w-32"
                        />
                        <button
                          onClick={() => handleFeedbackSubmit(item.id)}
                          disabled={!feedbackState[item.id] || loadingFeedback[item.id]}
                          className="bg-emerald-600 disabled:bg-gray-700 text-white rounded p-1.5 hover:bg-emerald-500 transition-colors"
                        >
                          <Send size={14} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
