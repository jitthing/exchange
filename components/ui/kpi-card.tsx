import { ReactNode } from 'react';

type Affordability = 'green' | 'amber' | 'red';

interface KpiCardProps {
  label: string;
  value: string;
  hint?: string;
  affordability?: Affordability;
  icon?: ReactNode;
}

const affordabilityStyles: Record<Affordability, string> = {
  green: 'bg-success-50 border border-success/20',
  amber: 'bg-warning-50 border border-warning/20',
  red: 'bg-danger-50 border border-danger/20',
};

const affordabilityValueColor: Record<Affordability, string> = {
  green: 'text-success-600',
  amber: 'text-warning-600',
  red: 'text-danger-600',
};

export function KpiCard({ label, value, hint, affordability, icon }: KpiCardProps) {
  const bgClass = affordability ? affordabilityStyles[affordability] : 'bg-white shadow-medium';
  const valueColor = affordability ? affordabilityValueColor[affordability] : 'text-heading';

  return (
    <article data-testid="kpi-card" className={`rounded-lg p-4 ${bgClass}`}>
      <div className="flex items-center gap-2">
        {icon ? <span className="text-lg">{icon}</span> : null}
        <p className="text-caption font-medium uppercase tracking-wider text-muted">{label}</p>
      </div>
      <p className={`mt-2 text-h2 font-semibold ${valueColor}`}>{value}</p>
      {hint ? <p className="mt-1 text-caption text-muted">{hint}</p> : null}
    </article>
  );
}
