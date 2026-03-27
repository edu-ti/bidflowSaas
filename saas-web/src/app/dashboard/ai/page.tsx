"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";

interface AIInsights {
  total_analyses: number;
  success_rate: number;
  average_score: number;
  enter_ratio: number;
  total_feedbacks: number;
}

export default function AIInsightsDashboard() {
  const [timeframe, setTimeframe] = useState<number>(30);
  const [data, setData] = useState<AIInsights | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchInsights();
  }, [timeframe]);

  const fetchInsights = async () => {
    setLoading(true);
    try {
      // Setup base URL or auth headers mapped to the actual Next.js utility config
      const token = localStorage.getItem("token") || "";
      const res = await fetch(`/api/v1/ai/insights?timeframe=${timeframe}`, {
        headers: { "Authorization": `Bearer ${token}` }
      });
      if (res.ok) {
        const json = await res.json();
        setData(json);
      }
    } catch (err) {
      console.error("Failed to load insights", err);
    }
    setLoading(false);
  };

  return (
    <div className="p-8 max-w-7xl mx-auto space-y-8 bg-gray-50 min-h-screen">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight">AI Insights Dashboard</h1>
          <p className="text-gray-500 mt-2">Monitor predictive models and success patterns.</p>
        </div>
        <div className="flex items-center space-x-4">
          <select 
            className="border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 text-sm p-2"
            value={timeframe} 
            onChange={(e) => setTimeframe(Number(e.target.value))}
          >
            <option value={30}>Last 30 Days</option>
            <option value={90}>Last 90 Days</option>
            <option value={0}>All Time</option>
          </select>
          <Link href="/dashboard/ai/history" className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 transition">
            View History
          </Link>
        </div>
      </div>

      {loading ? (
        <div className="animate-pulse space-y-4">
          <div className="h-32 bg-gray-200 rounded-xl w-full"></div>
        </div>
      ) : data ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <MetricCard 
            title="Total Analyses" 
            value={data.total_analyses.toString()} 
            sub="Items evaluated by AI"
            color="bg-blue-50"
          />
          <MetricCard 
            title="Success Rate" 
            value={`${(data.success_rate * 100).toFixed(1)}%`} 
            sub={`Based on ${data.total_feedbacks} valid feedbacks`}
            color="bg-green-50"
          />
          <MetricCard 
            title="Enter Ratio" 
            value={`${(data.enter_ratio * 100).toFixed(1)}%`} 
            sub="Opportunities marked ENTER"
            color="bg-purple-50"
          />
          <MetricCard 
            title="Avg AI Quality" 
            value={`${data.average_score.toFixed(0)}/100`} 
            sub="Internal AI heuristic baseline"
            color="bg-orange-50"
          />
        </div>
      ) : (
        <p className="text-gray-500">No data available for the selected timeframe.</p>
      )}
      
      {/* Charts placeholder for next phase integration with Recharts */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 h-96 flex flex-col items-center justify-center text-gray-400">
        <svg className="w-16 h-16 mb-4 text-gray-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <span className="font-medium">Performance Charts will render here.</span>
        <span className="text-sm">Install Recharts to map success rate over time.</span>
      </div>
    </div>
  );
}

function MetricCard({ title, value, sub, color }: { title: string, value: string, sub: string, color: string }) {
  return (
    <div className={`p-6 rounded-2xl ${color} border border-opacity-50 flex flex-col justify-between hover:shadow-md transition-shadow`}>
      <h3 className="text-sm font-semibold text-gray-600 uppercase tracking-wider">{title}</h3>
      <div className="mt-4 flex items-baseline text-4xl font-extrabold text-gray-900">
        {value}
      </div>
      <p className="mt-2 text-sm text-gray-500">{sub}</p>
    </div>
  );
}
