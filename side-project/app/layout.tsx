import type { Metadata, Viewport } from 'next';
import Link from 'next/link';
import '@/styles/globals.css';

export const metadata: Metadata = {
  title: 'Exchange Travel Planner',
  description: 'Mobile-first exchange and weekend trip planner',
  manifest: '/manifest.webmanifest'
};

export const viewport: Viewport = {
  themeColor: '#0b1324',
  width: 'device-width',
  initialScale: 1
};

const nav = [
  { href: '/', label: 'Home' },
  { href: '/calendar', label: 'Calendar' },
  { href: '/discover', label: 'Discover' },
  { href: '/budget', label: 'Budget' },
  { href: '/group', label: 'Group' },
  { href: '/settings', label: 'Settings' }
];

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <main className="mx-auto min-h-screen max-w-5xl px-4 pb-24 pt-5">{children}</main>
        <nav className="fixed bottom-0 left-0 right-0 border-t border-slate-200 bg-white/95 backdrop-blur">
          <div className="mx-auto grid max-w-5xl grid-cols-6 gap-1 px-2 py-2">
            {nav.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="rounded-lg px-2 py-2 text-center text-xs font-medium text-slate-700 hover:bg-slate-100"
              >
                {item.label}
              </Link>
            ))}
          </div>
        </nav>
      </body>
    </html>
  );
}
