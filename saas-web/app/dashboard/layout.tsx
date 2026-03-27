import '@/app/globals.css';
import { LayoutSidebar } from '@/components/LayoutSidebar';
import { LayoutHeader } from '@/components/LayoutHeader';
import { AuthProvider } from '@/hooks/useAuth';

export const metadata = {
  title: 'BidFlow Dashboard',
};

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <div className="flex min-h-screen bg-gray-900 text-gray-100">
        <LayoutSidebar />
        <div className="flex flex-col flex-1 ml-4">
          <LayoutHeader />
          <main className="flex-1 p-6 overflow-y-auto">{children}</main>
        </div>
      </div>
    </AuthProvider>
  );
}
