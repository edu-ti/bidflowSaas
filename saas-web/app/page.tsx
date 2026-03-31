import Link from 'next/link';
import { Button } from '@/components/Button';
import { Card } from '@/components/Card';
import { PricingSection } from '@/components/PricingSection';
import { ArrowRight, Search, Shield, Zap, Sparkles, BarChart3, Globe } from 'lucide-react';

export default function LandingPage() {
  return (
    <div className="flex flex-col min-h-screen bg-gray-950 text-gray-100 overflow-x-hidden">
      {/* Dynamic Background Effects */}
      <div className="fixed inset-0 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-emerald-500/10 blur-[120px] rounded-full" />
        <div className="absolute bottom-[10%] right-[-10%] w-[35%] h-[35%] bg-purple-600/10 blur-[120px] rounded-full" />
        <div className="absolute top-[20%] right-[10%] w-[25%] h-[25%] bg-blue-500/10 blur-[100px] rounded-full" />
        <div className="absolute inset-0 bg-grid-white opacity-20" />
      </div>

      {/* Navigation */}
      <nav className="relative z-10 flex items-center justify-between px-8 py-6 max-w-7xl mx-auto w-full">
        <div className="flex items-center gap-2 group">
          <div className="w-10 h-10 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center shadow-lg group-hover:rotate-12 transition-transform duration-300">
            <Sparkles className="text-white" size={20} />
          </div>
          <span className="text-2xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-400">
            BidFlow
          </span>
        </div>
        <div className="hidden md:flex items-center gap-8 text-sm font-medium text-gray-400">
          <a href="#features" className="hover:text-emerald-400 transition-colors">Funcionalidades</a>
          <a href="#pricing" className="hover:text-emerald-400 transition-colors">Preços</a>
          <Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Entrar</Link>
          <Link href="/dashboard">
            <Button className="bg-white !text-black hover:bg-gray-200 px-6 font-semibold">
              Get Started
            </Button>
          </Link>
        </div>
      </nav>

      <main className="relative z-10 flex flex-col items-center">
        {/* Hero Section */}
        <section className="relative pt-20 pb-32 flex flex-col items-center text-center px-6 max-w-5xl mx-auto">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-medium mb-8 animate-fade-in">
            <BadgeIcon />
            <span>Inteligência Artificial de Próxima Geração</span>
          </div>
          
          <h1 className="text-6xl md:text-7xl font-extrabold leading-[1.1] mb-8 tracking-tight text-balance">
            Vença licitações com o poder da <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 to-emerald-600">IA Generativa</span>
          </h1>
          
          <p className="text-xl text-gray-400 max-w-2xl mb-12 text-balance leading-relaxed">
            BidFlow automatiza a descoberta, análise e qualificação de editais públicos. 
            Aumente sua taxa de vitória em até 40% com insights estratégicos em tempo real.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 items-center">
            <Link href="/dashboard">
              <Button className="h-14 px-10 text-lg bg-emerald-500 hover:bg-emerald-400 transition-all shadow-[0_0_20px_rgba(16,185,129,0.3)]">
                Começar agora <ArrowRight className="ml-2" size={20} />
              </Button>
            </Link>
            <Link href="/dashboard/radar">
              <Button className="h-14 px-10 text-lg bg-gray-900 border border-gray-800 hover:bg-gray-850 hover:border-gray-700 transition-all">
                Ver demonstração
              </Button>
            </Link>
          </div>
        </section>

        {/* Features Slider / Grid */}
        <section id="features" className="py-24 px-8 max-w-7xl mx-auto w-full">
          <div className="mb-16 text-center">
            <h2 className="text-base font-semibold text-emerald-500 uppercase tracking-widest mb-4">Potencialize seus resultados</h2>
            <h3 className="text-4xl font-bold">Arquitetura construída para escala</h3>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <Card className="group">
              <div className="w-12 h-12 bg-emerald-500/10 rounded-lg flex items-center justify-center mb-6 group-hover:bg-emerald-500 transition-colors duration-500">
                <Search className="text-emerald-400 group-hover:text-white transition-colors" size={24} />
              </div>
              <h4 className="text-xl font-bold mb-3">Radar Inteligente</h4>
              <p className="text-gray-400 leading-relaxed">
                Varredura automática de milhares de portais de compras públicas. Receba apenas o que é relevante para o seu negócio.
              </p>
            </Card>

            <Card className="group">
              <div className="w-12 h-12 bg-blue-500/10 rounded-lg flex items-center justify-center mb-6 group-hover:bg-blue-500 transition-colors duration-500">
                <BarChart3 className="text-blue-400 group-hover:text-white transition-colors" size={24} />
              </div>
              <h4 className="text-xl font-bold mb-3">Análise Preditiva</h4>
              <p className="text-gray-400 leading-relaxed">
                Nossa IA analisa similaridade, riscos e chances de vitória baseada no seu histórico e nos requisitos técnicos do edital.
              </p>
            </Card>

            <Card className="group">
              <div className="w-12 h-12 bg-purple-500/10 rounded-lg flex items-center justify-center mb-6 group-hover:bg-purple-500 transition-colors duration-500">
                <Zap className="text-purple-400 group-hover:text-white transition-colors" size={24} />
              </div>
              <h4 className="text-xl font-bold mb-3">Insights Instantâneos</h4>
              <p className="text-gray-400 leading-relaxed">
                Resumos executivos de editais complexos gerados em segundos. Economize horas de leitura técnica e burocrática.
              </p>
            </Card>
          </div>
        </section>

        {/* Pricing */}
        <div id="pricing" className="py-24 w-full">
          <PricingSection />
        </div>

        {/* CTA Footer */}
        <section className="py-32 w-full flex flex-col items-center">
          <Card className="max-w-4xl w-full p-12 bg-gradient-to-br from-emerald-900/40 to-blue-900/20 border-emerald-500/20 text-center relative overflow-hidden">
            <div className="relative z-10">
              <h2 className="text-4xl font-bold mb-6">Pronto para dominar o mercado público?</h2>
              <p className="text-lg text-gray-400 mb-10 max-w-xl mx-auto">
                Junte-se a centenas de empresas que já estão usando BidFlow para crescer. 
                Teste grátis por 14 dias.
              </p>
              <Link href="/dashboard">
                <Button className="h-14 px-12 text-lg bg-white !text-black hover:bg-gray-200 font-bold">
                  Criar conta gratuita
                </Button>
              </Link>
            </div>
            {/* Decoration */}
            <div className="absolute top-0 right-0 w-64 h-64 bg-emerald-500/10 blur-[80px] rounded-full" />
            <div className="absolute bottom-0 left-0 w-48 h-48 bg-blue-500/10 blur-[60px] rounded-full" />
          </Card>
        </section>
      </main>

      <footer className="relative z-10 border-t border-gray-900 bg-gray-950/80 backdrop-blur-md pt-16 pb-12 px-8">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-start gap-12">
          <div className="space-y-4 max-w-xs">
            <div className="flex items-center gap-2">
              <Sparkles className="text-emerald-400" size={20} />
              <span className="text-xl font-bold">BidFlow</span>
            </div>
            <p className="text-sm text-gray-500">
              A próxima fronteira da inteligência competitiva em compras públicas. 
              Tecnologia brasileira, visão global.
            </p>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-12">
            <div className="space-y-4">
              <h5 className="font-semibold text-gray-100">Produto</h5>
              <ul className="space-y-2 text-sm text-gray-500">
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Funcionalidades</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Segurança</Link></li>
                <li><Link href="#pricing" className="hover:text-emerald-400 transition-colors">Pricing</Link></li>
              </ul>
            </div>
            <div className="space-y-4">
              <h5 className="font-semibold text-gray-100">Recursos</h5>
              <ul className="space-y-2 text-sm text-gray-500">
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Blog</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Documentação</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Suporte</Link></li>
              </ul>
            </div>
            <div className="space-y-4">
              <h5 className="font-semibold text-gray-100">Empresa</h5>
              <ul className="space-y-2 text-sm text-gray-500">
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Sobre</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Carreiras</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Legal</Link></li>
              </ul>
            </div>
          </div>
        </div>
        <div className="max-w-7xl mx-auto mt-16 pt-8 border-t border-gray-900 flex flex-col md:flex-row justify-between items-center gap-4 text-xs text-gray-600">
          <p>© 2025 BidFlow SaaS. Todos os direitos reservados.</p>
          <div className="flex gap-6">
            <Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Terms</Link>
            <Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Privacy</Link>
            <Link href="/dashboard" className="hover:text-emerald-400 transition-colors">Cookies</Link>
          </div>
        </div>
      </footer>
    </div>
  );
}

function BadgeIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M6 1V2M6 10V11M11 6H10M2 6H1M9.535 2.465L8.828 3.172M3.172 8.828L2.465 9.535M9.535 9.535L8.828 8.828M3.172 3.172L2.465 2.465" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
  );
}
