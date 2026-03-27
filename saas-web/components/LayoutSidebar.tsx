import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Search, Monitor, Bell, Settings, CreditCard, LogOut } from 'lucide-react';
import { useAuth } from '@/hooks/useAuth';

const navItems = [
  { name: 'Dashboard', href: '/dashboard', icon: Home },
  { name: 'Radar', href: '/dashboard/radar', icon: Search },
  { name: 'Monitor', href: '/dashboard/monitor', icon: Monitor },
  { name: 'Notificações', href: '/dashboard/notifications', icon: Bell },
  { name: 'IA', href: '/dashboard/ai', icon: Settings },
  { name: 'Billing', href: '/dashboard/billing', icon: CreditCard },
];

export const LayoutSidebar: React.FC = () => {
  const pathname = usePathname();
  const { logout } = useAuth();

  return (
    <aside className="w-64 bg-gray-800 text-gray-100 flex flex-col p-4 rounded-2xl shadow-soft h-screen">
      <nav className="flex-1 space-y-2">
        {navItems.map((item) => {
          const Icon = item.icon;
          const active = pathname?.startsWith(item.href);
          return (
            <Link key={item.name} href={item.href} className={
              `flex items-center p-2 rounded-lg hover:bg-gray-700 transition-colors ${active ? 'bg-gray-700' : ''}`
            }>
              <Icon className="mr-3" size={20} />
              <span>{item.name}</span>
            </Link>
          );
        })}
      </nav>
      <button
        onClick={logout}
        className="flex items-center p-2 rounded-lg hover:bg-gray-700 transition-colors"
      >
        <LogOut className="mr-3" size={20} />
        <span>Sair</span>
      </button>
    </aside>
  );
};
