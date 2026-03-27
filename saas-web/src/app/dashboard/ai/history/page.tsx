"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";

interface AIHistoryItem {
  edital_id: string;
  status: string;
  score: number;
  smart_score: number;
  confidence: number;
  decision: string;
  risk_level: string;
  real_result: string | null;
  created_at: string;
}

export default function AIHistoryPage() {
  const [history, setHistory] = useState<AIHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchHistory();
  }, []);

  const fetchHistory = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem("token") || "";
      const res = await fetch(`/api/v1/ai/history`, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (res.ok) {
        const json = await res.json();
        setHistory(json || []);
      }
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const submitFeedback = async (editalId: string, aiDecision: string, result: string) => {
    try {
      const token = localStorage.getItem("token") || "";
      await fetch('/api/v1/ai/feedback', {
        method: 'POST',
        headers: { 
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          edital_id: editalId,
          ai_decision: aiDecision,
          real_result: result
        })
      });
      // Refresh state locally
      setHistory(prev => prev.map(item => item.edital_id === editalId ? {...item, real_result: result} : item));
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="p-8 max-w-[1400px] mx-auto space-y-8 bg-gray-50 min-h-screen">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight">AI History & Feedback</h1>
          <p className="text-gray-500 mt-2">Provide real outcomes to train the model dynamically.</p>
        </div>
        <Link href="/dashboard/ai" className="text-blue-600 hover:text-blue-800 font-medium transition">
          &larr; Back to Insights
        </Link>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Edital</th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Dynamic Score</th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Confidence</th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Decision</th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Outcome Feedback</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-gray-500">Loading history...</td>
              </tr>
            ) : history.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-gray-500">No AI analysis history found.</td>
              </tr>
            ) : (
              history.map((item) => (
                <tr key={item.edital_id} className="hover:bg-gray-50 transition">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-mono">
                    {item.edital_id.substring(0, 8)}...
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center space-x-2">
                       <span className={`text-sm font-bold ${item.smart_score > 70 ? 'text-green-600' : item.smart_score > 40 ? 'text-orange-500' : 'text-red-600'}`}>
                         {item.smart_score}/100
                       </span>
                       <span className="text-xs text-gray-400" title="Base AI Score">({item.score})</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="w-full bg-gray-200 rounded-full h-2.5 max-w-[100px] mt-2">
                      <div className={`h-2.5 rounded-full ${item.confidence > 0.7 ? 'bg-blue-600' : item.confidence > 0.3 ? 'bg-blue-400' : 'bg-gray-400'}`} style={{ width: `${item.confidence * 100}%` }}></div>
                    </div>
                    <span className="text-xs text-gray-500 mt-1 block">{(item.confidence * 100).toFixed(0)}% metrics</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${item.decision === 'ENTER' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                      {item.decision}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {item.real_result ? (
                      <span className={`font-bold ${item.real_result === 'WON' ? 'text-green-600' : item.real_result === 'LOST' ? 'text-red-600' : 'text-gray-500'}`}>
                        {item.real_result}
                      </span>
                    ) : (
                      <div className="flex space-x-2">
                        <button onClick={() => submitFeedback(item.edital_id, item.decision, 'WON')} className="text-green-600 hover:text-green-900 bg-green-50 px-2 py-1 rounded">WON</button>
                        <button onClick={() => submitFeedback(item.edital_id, item.decision, 'LOST')} className="text-red-600 hover:text-red-900 bg-red-50 px-2 py-1 rounded">LOST</button>
                        <button onClick={() => submitFeedback(item.edital_id, item.decision, 'NOT_PARTICIPATED')} className="text-gray-600 hover:text-gray-900 bg-gray-100 px-2 py-1 rounded">NO_PARTICIPATION</button>
                      </div>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
