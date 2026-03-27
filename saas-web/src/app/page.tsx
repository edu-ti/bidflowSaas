import Link from "next/link";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-neutral-950 text-neutral-100 overflow-x-hidden selection:bg-indigo-500/30">
      {/* Navigation */}
      <nav className="fixed top-0 inset-x-0 h-16 bg-neutral-950/80 backdrop-blur-md border-b border-white/5 z-50 flex items-center justify-between px-8">
        <div className="text-xl font-black tracking-tight text-white flex items-center gap-2">
          <div className="w-6 h-6 rounded-full bg-gradient-to-tr from-indigo-500 to-blue-400"></div>
          BidFlow
        </div>
        <div className="flex items-center gap-6 text-sm font-medium hidden md:flex">
          <a href="#features" className="text-neutral-400 hover:text-white transition">Features</a>
          <a href="#pricing" className="text-neutral-400 hover:text-white transition">Pricing</a>
          <div className="w-px h-4 bg-white/10"></div>
        </div>
        <div className="flex items-center gap-4 text-sm font-medium">
          <Link href="/login" className="text-neutral-300 hover:text-white transition">Sign In</Link>
          <Link href="/signup" className="px-4 py-2 bg-white text-black hover:bg-neutral-200 transition-colors rounded-full font-bold">
            Start Free Trial
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="relative pt-32 pb-20 sm:pt-40 sm:pb-24 px-6 lg:px-8 max-w-7xl mx-auto text-center">
        {/* Abstract Glow */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[400px] bg-indigo-500/20 blur-[120px] rounded-full pointer-events-none -z-10"></div>
        
        <h1 className="text-5xl sm:text-7xl font-extrabold tracking-tight text-white mb-8 animate-in slide-in-from-bottom-8 duration-700">
          The ultimate OS for <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-blue-400 to-cyan-400">
            Tender Management
          </span>
        </h1>
        <p className="mx-auto max-w-2xl text-lg sm:text-xl text-neutral-400 mb-10 leading-relaxed animate-in slide-in-from-bottom-10 py-2 duration-1000">
          Supercharge your government contracting pipeline. BidFlow unites hyper-focused CRM tooling, automated Edital parsing, and AI-driven insights into one beautiful platform.
        </p>
        
        <div className="flex items-center justify-center gap-4 animate-in fade-in duration-1000 delay-300">
          <Link href="/signup" className="px-8 py-4 bg-indigo-600 hover:bg-indigo-500 text-white rounded-full font-bold text-lg transition-all hover:scale-105 shadow-[0_0_40px_rgba(79,70,229,0.4)]">
            Start Your Free Trial
          </Link>
          <a href="#demo" className="px-8 py-4 bg-white/5 border border-white/10 hover:bg-white/10 text-white rounded-full font-bold text-lg transition-all">
            Book a Demo
          </a>
        </div>
      </main>

      {/* Social Proof */}
      <section className="border-y border-white/5 bg-white/[0.02] py-12">
        <div className="max-w-7xl mx-auto px-6 text-center">
          <p className="text-sm font-semibold tracking-wider text-neutral-500 uppercase mb-8">Trusted by leading contractors</p>
          <div className="flex flex-wrap justify-center gap-12 opacity-50 grayscale">
            <span className="text-2xl font-black">ACME Corp</span>
            <span className="text-2xl font-black tracking-widest">GLOBEX</span>
            <span className="text-2xl font-black italic">Soylent</span>
            <span className="text-2xl font-black font-serif">Initech</span>
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="py-24 sm:py-32 max-w-7xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold text-white mb-4">Everything you need to win.</h2>
          <p className="text-neutral-400 text-lg max-w-2xl mx-auto">Stop managing tenders in spreadsheets. We bring your entire workflow onto a single pane of glass.</p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {[
            { 
              title: "Hyper-Focused CRM", 
              desc: "Track every lead, customer, and opportunity specifically modeled for B2G sales cycles.",
              icon: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
            },
            { 
              title: "Automated Licitações", 
              desc: "Ingest Editais instantly. Track Propostas, automate status updates, and view daily outcomes.",
              icon: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            },
            { 
              title: "AI Bidding Insights", 
              desc: "Run predictive margin analysis and let AI parse 500-page government docs in seconds.",
              icon: "M13 10V3L4 14h7v7l9-11h-7z"
            }
          ].map((feature, i) => (
            <div key={i} className="bg-neutral-900 border border-neutral-800 p-8 rounded-3xl hover:border-neutral-700 transition">
              <div className="w-12 h-12 bg-indigo-500/10 text-indigo-400 rounded-2xl flex items-center justify-center mb-6">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d={feature.icon}></path></svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">{feature.title}</h3>
              <p className="text-neutral-400 leading-relaxed">{feature.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Pricing */}
      <section id="pricing" className="py-24 bg-neutral-900 border-t border-white/5 relative overflow-hidden">
        <div className="absolute bottom-0 right-0 w-[500px] h-[500px] bg-blue-500/10 blur-[120px] rounded-full pointer-events-none -z-10"></div>
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-white mb-4">Simple, transparent pricing.</h2>
            <p className="text-neutral-400 text-lg">Start for free, upgrade when you need more power.</p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {/* Starter */}
            <div className="bg-neutral-950 border border-neutral-800 rounded-3xl p-8">
              <h3 className="text-xl font-semibold text-white">Starter</h3>
              <div className="mt-4 flex items-baseline text-5xl font-extrabold text-white">
                $0
                <span className="ml-1 text-xl font-medium text-neutral-500">/mo</span>
              </div>
              <p className="mt-4 text-neutral-400">Perfect for exploring the platform.</p>
              <ul className="mt-8 space-y-4 text-sm text-neutral-300">
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> 100 CRM Leads</li>
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> 5 Editais per month</li>
                <li className="flex items-center gap-3"><span className="text-neutral-600">✗</span> AI Insights</li>
              </ul>
              <Link href="/signup?plan=starter" className="mt-8 block w-full py-3 px-6 bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl text-center font-semibold text-white transition">
                Get Started
              </Link>
            </div>

            {/* Pro */}
            <div className="bg-neutral-800 border-2 border-indigo-500 rounded-3xl p-8 relative transform md:-translate-y-4 shadow-2xl shadow-indigo-500/10">
              <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-indigo-500 text-white px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wide">
                Most Popular
              </div>
              <h3 className="text-xl font-semibold text-white">Pro</h3>
              <div className="mt-4 flex items-baseline text-5xl font-extrabold text-white">
                $49
                <span className="ml-1 text-xl font-medium text-neutral-400">/mo</span>
              </div>
              <p className="mt-4 text-neutral-300">For scaling agencies winning bids.</p>
              <ul className="mt-8 space-y-4 text-sm text-neutral-200">
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> 1,000 CRM Leads</li>
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> Unlimited Editais</li>
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> Basic AI Parsing</li>
              </ul>
              <Link href="/signup?plan=pro" className="mt-8 block w-full py-3 px-6 bg-indigo-500 hover:bg-indigo-400 rounded-xl text-center font-semibold text-white shadow-lg shadow-indigo-500/25 transition">
                Start 14-Day Trial
              </Link>
            </div>

            {/* Enterprise */}
            <div className="bg-neutral-950 border border-neutral-800 rounded-3xl p-8">
              <h3 className="text-xl font-semibold text-white">Enterprise</h3>
              <div className="mt-4 flex items-baseline text-5xl font-extrabold text-white">
                $199
                <span className="ml-1 text-xl font-medium text-neutral-500">/mo</span>
              </div>
              <p className="mt-4 text-neutral-400">Unlimited power for massive operations.</p>
              <ul className="mt-8 space-y-4 text-sm text-neutral-300">
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> Unlimited Data</li>
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> Advanced Custom AI Agent</li>
                <li className="flex items-center gap-3"><span className="text-indigo-400">✓</span> Priority API Support</li>
              </ul>
              <Link href="/signup?plan=enterprise" className="mt-8 block w-full py-3 px-6 bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl text-center font-semibold text-white transition">
                Contact Sales
              </Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
