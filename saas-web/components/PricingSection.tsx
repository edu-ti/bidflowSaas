'use client';

import { Card } from '@/components/Card';
import { Badge } from '@/components/Badge';
import { Button } from '@/components/Button';
import { postCheckout } from '@/lib/api';
import { useState } from 'react';

const plans = [
  { name: 'Starter', price: '49' },
  { name: 'Pro', price: '149' },
  { name: 'Enterprise', price: '399' },
];

export function PricingSection() {
  const [loading, setLoading] = useState<string | null>(null);

  async function handleCheckout(planName: string) {
    setLoading(planName);
    try {
      const data = await postCheckout(planName.toLowerCase());
      if (data.checkout_url) {
        window.location.href = data.checkout_url;
      } else {
        console.log('Backend response missing checkout_url:', data);
      }
    } catch (err) {
      console.log('Fallback - Error calling checkout API:', err);
      // Fallback behavior if backend is not ready
      alert(`Backend indisponível para assinar o plano ${planName}. Verifique o console.`);
    } finally {
      setLoading(null);
    }
  }

  return (
    <section className="max-w-4xl mx-auto w-full px-4 text-center">
      <h2 className="text-3xl font-bold mb-8 text-emerald-400">Planos Disponíveis</h2>
      <div className="grid md:grid-cols-3 gap-6">
        {plans.map(({ name, price }) => (
          <Card key={name} className="flex flex-col items-center p-8 bg-gray-900 border-gray-800 hover:border-emerald-500/30 transition-colors">
            <Badge variant="low" className="mb-4 text-sm bg-gray-800 text-gray-300">
              {name}
            </Badge>
            <p className="text-4xl font-bold mb-1">R$ {price}</p>
            <p className="text-gray-500 text-sm mb-6">/mês, cobrado anualmente</p>
            <Button
              onClick={() => handleCheckout(name)}
              disabled={loading === name}
              className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white w-full py-3 rounded-xl transition-all font-medium"
            >
              {loading === name ? 'Processando...' : `Assinar ${name}`}
            </Button>
          </Card>
        ))}
      </div>
    </section>
  );
}
