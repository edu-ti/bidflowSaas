"use client";

import { useState } from "react";

export default function SignupFlow() {
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);

  const simulateStep = (next: number) => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      if (next === 4) {
        window.location.href = "/dashboard";
      } else {
        setStep(next);
      }
    }, 1000);
  };

  return (
    <div className="min-h-screen bg-neutral-950 flex flex-col items-center justify-center p-6 text-white relative overflow-hidden">
      {/* Background aesthetic */}
      <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20 pointer-events-none mix-blend-overlay"></div>
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-indigo-500/10 blur-[100px] rounded-full point-events-none -z-10"></div>

      <div className="w-full max-w-md bg-neutral-900 border border-neutral-800 rounded-3xl p-8 shadow-2xl z-10 animate-in zoom-in-95 duration-500">
        {/* Progress Bar */}
        <div className="flex gap-2 mb-8">
          {[1, 2, 3].map((s) => (
            <div key={s} className={`h-1 flex-1 rounded-full transition-colors ${s <= step ? 'bg-indigo-500' : 'bg-neutral-800'}`}></div>
          ))}
        </div>

        {step === 1 && (
          <div className="space-y-6 animate-in slide-in-from-right-4">
            <h1 className="text-2xl font-bold tracking-tight">Create your account</h1>
            <p className="text-neutral-400 text-sm">Join BidFlow to supercharge your B2G pipeline.</p>
            <input type="email" placeholder="Email Address" className="w-full bg-neutral-950 border border-neutral-800 rounded-xl px-4 py-3 text-white placeholder:text-neutral-600 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all" />
            <input type="password" placeholder="Password" className="w-full bg-neutral-950 border border-neutral-800 rounded-xl px-4 py-3 text-white placeholder:text-neutral-600 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all" />
            <button onClick={() => simulateStep(2)} disabled={loading} className="w-full bg-white text-black font-bold py-3 rounded-xl hover:bg-neutral-200 transition">
              {loading ? "Authenticating..." : "Continue"}
            </button>
          </div>
        )}

        {step === 2 && (
          <div className="space-y-6 animate-in slide-in-from-right-4">
            <h1 className="text-2xl font-bold tracking-tight">Name your Workspace</h1>
            <p className="text-neutral-400 text-sm">Your tenant namespace for all shared data securely isolated from others.</p>
            <input type="text" placeholder="Workspace Name (e.g. Acme Corp)" className="w-full bg-neutral-950 border border-neutral-800 rounded-xl px-4 py-3 text-white placeholder:text-neutral-600 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all" />
            <button onClick={() => simulateStep(3)} disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-500 text-white font-bold py-3 rounded-xl transition shadow-[0_0_20px_rgba(79,70,229,0.3)]">
              {loading ? "Creating Tenant..." : "Create Tenant"}
            </button>
          </div>
        )}

        {step === 3 && (
          <div className="space-y-6 animate-in slide-in-from-right-4">
            <h1 className="text-2xl font-bold tracking-tight">Setup Billing</h1>
            <p className="text-neutral-400 text-sm">You selected the <strong>Pro Plan</strong>. You won&apos;t be charged for 14 days.</p>
            
            <div className="p-4 bg-neutral-950 border border-neutral-800 rounded-xl mb-6">
              <div className="flex justify-between items-center text-sm font-medium mb-2">
                <span className="text-neutral-400">Total due today</span>
                <span className="text-white">$0.00</span>
              </div>
              <div className="flex justify-between items-center text-sm font-medium">
                <span className="text-neutral-400">After trial</span>
                <span className="text-white">$49.00 / mo</span>
              </div>
            </div>

            <button onClick={() => simulateStep(4)} disabled={loading} className="w-full bg-blue-600 hover:bg-blue-500 text-white font-bold py-3 rounded-xl transition shadow-[0_0_20px_rgba(37,99,235,0.3)] flex items-center justify-center gap-2">
              {loading ? (
                "Forwarding to Stripe..."
              ) : (
                <>
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path></svg>
                  Mock Stripe Checkout
                </>
              )}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
