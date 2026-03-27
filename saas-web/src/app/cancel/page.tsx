import Link from "next/link";

export default function StripeCancel() {
  return (
    <div className="min-h-screen bg-neutral-950 flex flex-col items-center justify-center p-6 text-white text-center">
      <div className="w-20 h-20 bg-amber-500/10 text-amber-400 rounded-full flex items-center justify-center mb-8 shadow-[0_0_50px_rgba(245,158,11,0.2)] animate-in zoom-in duration-500">
        <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M6 18L18 6M6 6l12 12"></path></svg>
      </div>
      <h1 className="text-4xl font-bold mb-4 animate-in slide-in-from-bottom-4 duration-500 delay-100">Checkout Canceled</h1>
      <p className="text-neutral-400 max-w-md mx-auto mb-10 text-lg animate-in slide-in-from-bottom-4 duration-500 delay-200">
        You can securely try completing your payment anytime. The active SaaS subscription metrics revert to Trial boundaries until payment succeeds.
      </p>
      <Link href="/dashboard/billing" className="px-8 py-3 bg-white/10 text-white font-bold rounded-lg border border-white/20 hover:bg-white/20 transition animate-in fade-in duration-500 delay-300">
        Return to Billing View
      </Link>
    </div>
  );
}
