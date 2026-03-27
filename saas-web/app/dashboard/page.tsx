import { RadarCard, type RadarEdital } from '@/components/RadarCard';
import { Card } from '@/components/Card';

// Mock data – replace with real API fetch when backend is ready
const mockEditais: RadarEdital[] = [
  {
    id: '1',
    title: 'Pregão Eletrônico nº 001/2025 – Fornecimento de Software de Gestão Municipal',
    value: 480000,
    deadline: '2025-04-15T18:00:00Z',
    priority_label: 'HIGH',
    win_probability: 0.78,
    similarity: 0.91,
    why_relevant:
      'Altíssima aderência com o portfólio de soluções SaaS da empresa. Requisitos técnicos alinhados com certificações atuais.',
  },
  {
    id: '2',
    title: 'Tomada de Preços nº 003/2025 – Serviços de TI e Suporte Técnico',
    value: 210000,
    deadline: '2025-04-22T17:00:00Z',
    priority_label: 'MEDIUM',
    win_probability: 0.55,
    similarity: 0.73,
    why_relevant:
      'Compatível com serviços de helpdesk e suporte remoto. Concorrência moderada esperada.',
  },
  {
    id: '3',
    title: 'Pregão Presencial nº 012/2025 – Licenças de Software Educacional',
    value: 95000,
    deadline: '2025-05-03T14:00:00Z',
    priority_label: 'LOW',
    win_probability: 0.32,
    similarity: 0.48,
    why_relevant:
      'Possibilidade de participação com parceiro educacional. Risco elevado de margens reduzidas.',
  },
];

const stats = [
  { label: 'Editais ativos', value: '24', color: 'text-emerald-400' },
  { label: 'Análises IA pendentes', value: '7', color: 'text-amber-400' },
  { label: 'Probabilidade média', value: '61%', color: 'text-blue-400' },
  { label: 'Valor total em radar', value: 'R$ 785k', color: 'text-purple-400' },
];

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      {/* Page header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-100">Visão Geral</h1>
        <p className="text-gray-400 mt-1">Seus editais mais relevantes de hoje</p>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {stats.map((s) => (
          <Card key={s.label} className="p-6 text-center">
            <p className={`text-3xl font-bold ${s.color}`}>{s.value}</p>
            <p className="text-sm text-gray-400 mt-1">{s.label}</p>
          </Card>
        ))}
      </div>

      {/* Radar section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-200 mb-4">Radar de Oportunidades</h2>
        <div className="grid md:grid-cols-3 gap-6">
          {mockEditais.map((edital) => (
            <RadarCard key={edital.id} edital={edital} />
          ))}
        </div>
      </div>
    </div>
  );
}
