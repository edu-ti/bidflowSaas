import Link from 'next/link';
import { Button } from '@/components/Button';
import { Card } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { ArrowRight } from 'lucide-react';

export default function LandingPage() {
  return (
    <main className="flex flex-col items-center px-4 py-12 space-y-16">
      {/* Hero Section */}
      <section className="text-center max-w-3xl">
        <h1 className="text-5xl font-bold mb-4 text-emerald-500">
          Seu copiloto de licitações com IA
        </h1>
        <p className="text-xl text-gray-300 mb-8">
          Descubra oportunidades, analise editais e aumente suas chances de vitória
        </p>
        <Link href="/dashboard">
          <Button className="bg-emerald-600 hover:bg-emerald-500 text-white px-8 py-3 rounded-2xl transition-colors">
            Começar agora <ArrowRight className="inline-block ml-2" size={18} />
          </Button>
        </Link>
      </section>

      {/* Features */}
      <section className="grid md:grid-cols-3 gap-8 max-w-5xl">
        <Card className="text-center">
          <h3 className="text-2xl font-semibold mb-2 text-emerald-400">Radar Inteligente</h3>
          <p className="text-gray-300">Descubra editais relevantes usando IA e preferências do tenant.</p>
        </Card>
        <Card className="text-center">
          <h3 className="text-2xl font-semibold mb-2 text-emerald-400">Análise com IA</h3>
          <p className="text-gray-300">Obtenha insights, riscos e recomendações automáticas.</p>
        </Card>
        <Card className="text-center">
          <h3 className="text-2xl font-semibold mb-2 text-emerald-400">Monitoramento em tempo real</h3>
          <p className="text-gray-300">Fique por dentro de mudanças de pregão e mensagens urgentes.</p>
        </Card>
      </section>

      {/* Pricing */}
      <section className="max-w-4xl w-full">
        <h2 className="text-3xl font-bold text-center mb-8 text-emerald-400">Planos</h2>
        <div className="grid md:grid-cols-3 gap-6">
          {['Starter', 'Pro', 'Enterprise'].map((plan) => (
            <Card key={plan} className="flex flex-col items-center p-8">
              <Badge variant="low" className="mb-4 text-sm">
                {plan}
              </Badge>
              <p className="text-2xl font-semibold mb-4">R$ {plan === 'Starter' ? '49' : plan === 'Pro' ? '149' : '399'} / mês</p>
              <Button
                onClick={async () => {
                  // Simple checkout call – replace with real integration later
                  const res = await fetch('/api/v1/billing/checkout', { method: 'POST' });
                  const data = await res.json();
                  if (data.url) window.location.href = data.url;
                }}
                className="bg-gradient-to-r from-gradientStart to-gradientEnd text-white px-6 py-2 rounded-2xl"
              >
                Assinar {plan}
              </Button>
            </Card>
          ))}
        </div>
      </section>
    </main>
  );
}
