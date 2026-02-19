import type { Metadata, Viewport } from 'next';
import { Inter } from 'next/font/google';
import { BottomNav } from '@/components/ui/bottom-nav';
import { AuthProvider } from '@/lib/auth-context';
import '@/styles/globals.css';

const inter = Inter({ subsets: ['latin'], variable: '--font-inter' });

export const metadata: Metadata = {
  title: 'Exchange Travel Planner',
  description: 'Mobile-first exchange and weekend trip planner',
  manifest: '/manifest.webmanifest'
};

export const viewport: Viewport = {
  themeColor: '#FAFAF8',
  width: 'device-width',
  initialScale: 1
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={inter.variable}>
      <body className="font-sans antialiased">
        <AuthProvider>
          <main className="mx-auto min-h-screen max-w-lg px-4 pb-24 pt-6">{children}</main>
          <BottomNav />
        </AuthProvider>
      </body>
    </html>
  );
}
