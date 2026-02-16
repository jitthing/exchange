export function KpiCard({ label, value, hint }: { label: string; value: string; hint?: string }) {
  return (
    <article className="card">
      <p className="text-xs uppercase tracking-wide text-slate-500">{label}</p>
      <p className="mt-2 text-2xl font-semibold text-ink">{value}</p>
      {hint ? <p className="mt-1 text-xs text-slate-500">{hint}</p> : null}
    </article>
  );
}
