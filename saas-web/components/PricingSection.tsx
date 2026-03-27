'use client';

import { Card } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { Button } from '@/components/Button';

const plans = [
  { name: 'Starter', price: '49' },
  { name: 'Pro', price: '149' },
  { name: 'Enterprise', price: '399' },
];

export function PricingSection() {
  async function handleCheckout() {
    const res = await fetch('/api/v1/billing/checkout', { method: 'POST' });
    const data = await res.json();
    if (data.url) window.location.href = data.url;
  }

  return (
    <section className="max-w-4xl w-full">
      <h2 className="text-3xl font-bold text-center mb-8 text-emerald-400">Planos</h2>
      <div className="grid md:grid-cols-3 gap-6">
        {plans.map(({ name, price }) => (
          <Card key={name} className="flex flex-col items-center p-8">
            <Badge variant="low" className="mb-4 text-sm">
              {name}
            </Badge>
            <p className="text-2xl font-semibold mb-4">R$ {price} / mês</p>
            <Button
              onClick={handleCheckout}
              className="bg-gradient-to-r from-gradientStart to-gradientEnd text-white px-6 py-2 rounded-2xl"
            >
              Assinar {name}
            </Button>
          </Card>
        ))}
      </div>
    </section>
  );
}
