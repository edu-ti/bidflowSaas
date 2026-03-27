"use client";

import { useState } from "react";

export default function BillingPage() {
  const [loading, setLoading] = useState(false);

  const handleUpgrade = async () => {
    setLoading(true);
    // Simulate POST /api/v1/billing/checkout
    setTimeout(() => {
      window.location.href = "/success";
    }, 1500);
  };

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in duration-500">
      <header>
        <h1 className="text-3xl font-bold tracking-tight">Billing & Plans</h1>
        <p className="text-neutral-400 mt-1">Manage your active SaaS subscriptions and invoices.</p>
      </header>

      {/* Active Subscription Box */}
      <div className="bg-neutral-900 border border-neutral-800 rounded-3xl p-8 relative overflow-hidden">
        {/* Glow effect */}
        <div className="absolute top-0 right-0 p-32 bg-blue-500/10 blur-[100px] pointer-events-none rounded-full"></div>
        
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div>
            <h2 className="text-xl font-bold text-white flex items-center gap-3">
              Pro Suite
              <span className="px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                ACTIVE
              </span>
            </h2>
            <p className="text-neutral-400 mt-2 text-sm leading-relaxed max-w-md">
              You are currently on the Pro plan extending CRM bounds to 1,000 leads and enabling unlimited AI insights.
            </p>
            <div className="mt-6 flex flex-col sm:flex-row gap-4">
              <button 
                onClick={handleUpgrade}
                disabled={loading}
                className="px-6 py-2.5 bg-blue-600 hover:bg-blue-500 text-white font-semibold rounded-lg shadow-[0_0_20px_rgba(37,99,235,0.3)] transition-all flex items-center justify-center gap-2 group disabled:opacity-70 disabled:cursor-not-allowed"
              >
                {loading ? "Forwarding..." : "Upgrade to Enterprise"}
                {!loading && <svg className="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path></svg>}
              </button>
              <button className="px-6 py-2.5 bg-neutral-800 border border-neutral-700 hover:bg-red-500/10 hover:text-red-400 hover:border-red-500/30 text-white font-semibold rounded-lg transition-colors">
                Cancel Subscription
              </button>
            </div>
          </div>

          <div className="bg-neutral-950 border border-neutral-800 p-5 rounded-2xl min-w-[240px]">
            <p className="text-sm font-medium text-neutral-400">Next Invoice</p>
            <p className="text-3xl font-black text-white mt-1">$49<span className="text-base font-normal text-neutral-500">/mo</span></p>
            <div className="h-px bg-neutral-800 w-full my-4"></div>
            <p className="text-xs text-neutral-500 font-medium">Auto-renews April 26, 2026</p>
          </div>
        </div>
      </div>

      {/* Warning Box */}
      <div className="bg-amber-500/10 border border-amber-500/20 rounded-2xl p-5 flex items-start gap-4">
        <div className="p-2 bg-amber-500/20 rounded-lg text-amber-400">
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        </div>
        <div>
          <h4 className="text-amber-100 font-semibold">CRM Limit Approaching</h4>
          <p className="text-sm text-amber-200/70 mt-1">
            You have used 84% of your monthly CRM Lead capacity. Once you hit 1,000 leads, the system will block new entries until the next billing cycle.
          </p>
        </div>
      </div>
    </div>
  );
}
