'use client';

import { useState } from 'react';
import Link from 'next/link';
import { SectionTitle } from '@/components/section-title';
import { optimizeTrips } from '@/lib/api';
import { TripOption } from '@/lib/types';

const initialForm = {
  budgetCap: 280,
  maxTravelHours: 5,
  partySize: 2,
  style: 'culture' as const,
  windowId: 'w-1',
  departureCity: 'Berlin'
};

export default function DiscoverPage() {
  const [form, setForm] = useState(initialForm);
  const [results, setResults] = useState<TripOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError('');
    try {
      const payload = await optimizeTrips(form);
      setResults(payload.options);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to optimize');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-4">
      <SectionTitle title="Weekend Optimizer" subtitle="Suggests destinations by budget, timing, and style." />

      <form className="card space-y-3" onSubmit={onSubmit}>
        <div>
          <label className="mb-1 block text-sm font-medium">Departure city</label>
          <input
            className="input"
            value={form.departureCity}
            onChange={(event) => setForm((old) => ({ ...old, departureCity: event.target.value }))}
          />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="mb-1 block text-sm font-medium">Budget cap (EUR)</label>
            <input
              className="input"
              type="number"
              value={form.budgetCap}
              onChange={(event) => setForm((old) => ({ ...old, budgetCap: Number(event.target.value) }))}
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">Max travel hours</label>
            <input
              className="input"
              type="number"
              value={form.maxTravelHours}
              onChange={(event) => setForm((old) => ({ ...old, maxTravelHours: Number(event.target.value) }))}
            />
          </div>
        </div>
        <button className="button" type="submit" disabled={loading}>
          {loading ? 'Optimizing...' : 'Find trips'}
        </button>
        {error ? <p className="text-sm text-red-700">{error}</p> : null}
      </form>

      <section className="space-y-3">
        {results.map((option) => (
          <article key={option.id} className="card space-y-2">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">{option.destination}</h2>
              <span className="text-sm font-medium">EUR {option.totalEstimatedCost.toFixed(0)}</span>
            </div>
            <p className="text-xs text-slate-500">{option.reasonTags.join(' | ')}</p>
            <div className="text-sm text-slate-600">
              Best transport: {option.transportOptions[0]?.provider} ({option.transportOptions[0]?.durationHours}h)
            </div>
            <Link href="/trips/trip-1" className="button-secondary inline-block">
              View sample trip detail
            </Link>
          </article>
        ))}
      </section>
    </div>
  );
}
