export default function DashboardOverview() {
  // Hardcoded mock values simulating GET /api/v1/billing/overview
  const plan = {
    name: "Pro Suite",
    status: "active",
    renewalDate: "April 26, 2026",
  };

  const limits = [
    { label: "CRM Leads", used: 845, max: 1000 },
    { label: "Licitacao Editais", used: 12, max: 50 },
    { label: "Propostas Generated", used: 25, max: 200 },
  ];

  return (
    <div className="max-w-5xl mx-auto space-y-8 animate-in fade-in duration-500">
      <header>
        <h1 className="text-3xl font-bold tracking-tight">Welcome back!</h1>
        <p className="text-neutral-400 mt-1">Here is a summary of your workspace usage this cycle.</p>
      </header>

      {/* Quick Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-neutral-900 border border-neutral-800 p-6 rounded-2xl shadow-sm">
          <p className="text-sm font-medium text-neutral-400">Current Plan</p>
          <div className="mt-2 flex items-center gap-3">
            <span className="text-2xl font-bold text-white">{plan.name}</span>
            <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              {plan.status.toUpperCase()}
            </span>
          </div>
        </div>
        
        <div className="bg-neutral-900 border border-neutral-800 p-6 rounded-2xl shadow-sm">
          <p className="text-sm font-medium text-neutral-400">Renewal Date</p>
          <p className="mt-2 text-2xl font-bold text-white">{plan.renewalDate}</p>
        </div>

        <div className="bg-neutral-900 border border-neutral-800 p-6 rounded-2xl shadow-sm bg-gradient-to-br from-indigo-500/10 to-transparent">
          <p className="text-sm font-medium text-indigo-300">Total Usage</p>
          <p className="mt-2 text-2xl font-bold text-indigo-100">
            {Math.round((limits[0].used / limits[0].max) * 100)}% Maxed
          </p>
        </div>
      </div>

      {/* Usage Bars */}
      <div className="bg-neutral-900 border border-neutral-800 rounded-2xl p-6 shadow-sm mt-8">
        <h3 className="text-lg font-bold text-white mb-6">Resource Allocation</h3>
        <div className="space-y-6">
          {limits.map((item, idx) => {
            const percentage = Math.min((item.used / item.max) * 100, 100);
            const isNearing = percentage > 80;
            return (
              <div key={idx}>
                <div className="flex justify-between text-sm mb-2">
                  <span className="font-medium text-neutral-200">{item.label}</span>
                  <span className="text-neutral-400 font-mono">
                    {item.used} / {item.max} ({Math.round(percentage)}%)
                  </span>
                </div>
                <div className="w-full bg-neutral-800 rounded-full h-2.5 overflow-hidden">
                  <div
                    className={"h-2.5 rounded-full transition-all duration-1000 ease-out " + (isNearing ? "bg-amber-500" : "bg-blue-500")}
                    style={{ width: `${percentage}%` }}
                  ></div>
                </div>
                {isNearing && (
                  <p className="text-xs text-amber-500/80 font-medium mt-1.5 flex items-center gap-1">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
                    Nearing limit. Consider upgrading.
                  </p>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
