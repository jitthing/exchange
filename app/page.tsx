import Link from 'next/link';
import { KpiCard } from '@/components/kpi-card';
import { SectionTitle } from '@/components/section-title';
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

export default async function HomePage() {
  const { forecast, windows } = await getData();
  const nextWindow = windows[0];

  return (
    <div className="space-y-6">
      <SectionTitle
        title="Exchange Travel Planner"
        subtitle="Plan weekend trips without losing track of budget, deadlines, and group details."
      />

      <section className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <KpiCard
          label="Projected spend"
          value={`EUR ${forecast.projectedMonthlySpend.toFixed(0)}`}
          hint="Current month including trips"
        />
        <KpiCard
          label="Remaining"
          value={`EUR ${forecast.remainingBudget.toFixed(0)}`}
          hint={`Budget status: ${forecast.affordability}`}
        />
        <KpiCard
          label="Next safe window"
          value={nextWindow ? `${nextWindow.startDate} to ${nextWindow.endDate}` : 'No windows'}
          hint={nextWindow ? `Score ${nextWindow.score}` : 'Import your calendar'}
        />
      </section>

      <section className="card space-y-3">
        <h2 className="text-lg font-semibold">Quick Actions</h2>
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          <Link href="/discover" className="button text-center">
            Run Weekend Optimizer
          </Link>
          <Link href="/calendar" className="button-secondary text-center">
            Review Academic Conflicts
          </Link>
        </div>
      </section>
    </div>
  );
}
