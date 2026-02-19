import Link from 'next/link';
import { KpiCard } from '@/components/ui/kpi-card';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { API_BASE_URL } from '@/lib/config';

type Forecast = {
  projectedMonthlySpend: number;
  remainingBudget: number;
  affordability: 'green' | 'amber' | 'red';
};

type TravelWindow = {
  id: string;
  startDate: string;
  endDate: string;
  score: number;
};

async function getData(): Promise<{ forecast: Forecast; windows: TravelWindow[] }> {
  const [forecastResponse, windowsResponse] = await Promise.all([
    fetch(`${API_BASE_URL}/api/budget/forecast?userId=demo-user`, { cache: 'no-store' }),
    fetch(`${API_BASE_URL}/api/travel-windows`, { cache: 'no-store' })
  ]);

  const forecast = (await forecastResponse.json()) as Forecast;
  const windowsPayload = (await windowsResponse.json()) as { windows: TravelWindow[] };
  return { forecast, windows: windowsPayload.windows };
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short' });
}

export default async function HomePage() {
  const { forecast, windows } = await getData();
  const nextWindow = windows[0];

  const affordabilityBadge = {
    green: { variant: 'safe' as const, label: 'On track' },
    amber: { variant: 'warning' as const, label: 'Watch it' },
    red: { variant: 'danger' as const, label: 'Over budget' },
  }[forecast.affordability];

  return (
    <div className="space-y-6">
      <div>
        <p className="text-body text-muted">Good to see you 👋</p>
        <h1 data-testid="page-title" className="text-display text-heading">Travel Planner</h1>
      </div>

      <section className="grid grid-cols-2 gap-3">
        <KpiCard
          label="Projected"
          value={`€${forecast.projectedMonthlySpend.toFixed(0)}`}
          hint="This month"
          icon="📊"
          affordability={forecast.affordability}
        />
        <KpiCard
          label="Remaining"
          value={`€${forecast.remainingBudget.toFixed(0)}`}
          hint={affordabilityBadge.label}
          icon="💰"
          affordability={forecast.affordability}
        />
      </section>

      {nextWindow ? (
        <Card shadow="raised" className="relative overflow-hidden">
          <div className="absolute right-4 top-4">
            <Badge variant="info">Score {nextWindow.score}</Badge>
          </div>
          <p className="text-caption font-medium uppercase tracking-wider text-muted">Next travel window</p>
          <p className="mt-2 text-h2 text-heading">
            {formatDate(nextWindow.startDate)} — {formatDate(nextWindow.endDate)}
          </p>
          <p className="mt-1 text-small text-muted">Best upcoming window for a trip</p>
          <Link
            href="/calendar"
            className="mt-3 inline-block text-small font-medium text-primary hover:underline"
          >
            View all windows →
          </Link>
        </Card>
      ) : null}

      <section className="grid grid-cols-2 gap-3">
        <Link href="/discover" className="button-accent text-center">
          🧭 Find Trips
        </Link>
        <Link href="/calendar" className="button-secondary text-center">
          📅 Calendar
        </Link>
      </section>
    </div>
  );
}
