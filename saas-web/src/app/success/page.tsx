import Link from "next/link";

export default function StripeSuccess() {
  return (
    <div className="min-h-screen bg-neutral-950 flex flex-col items-center justify-center p-6 text-white text-center">
      <div className="w-20 h-20 bg-emerald-500/10 text-emerald-400 rounded-full flex items-center justify-center mb-8 shadow-[0_0_50px_rgba(16,185,129,0.2)] animate-in zoom-in duration-500">
        <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="3" d="M5 13l4 4L19 7"></path></svg>
      </div>
      <h1 className="text-4xl font-bold mb-4 animate-in slide-in-from-bottom-4 duration-500 delay-100">Payment Successful</h1>
      <p className="text-neutral-400 max-w-md mx-auto mb-10 text-lg animate-in slide-in-from-bottom-4 duration-500 delay-200">
        Your subscription has been activated securely via Stripe! Your account metrics and limits are now unlocked globally.
      </p>
      <Link href="/dashboard" className="px-8 py-3 bg-white text-black font-bold rounded-lg hover:bg-neutral-200 transition animate-in fade-in duration-500 delay-300">
        Go to Dashboard
      </Link>
    </div>
  );
}
