'use client';

import { useEffect, useState } from 'react';
import { SectionTitle } from '@/components/section-title';
import { createBudgetEntry, getBudgetEntries, getForecast } from '@/lib/api';
import { BudgetEntry, ForecastResult } from '@/lib/types';

const formSeed = {
  userId: 'demo-user',
  category: 'travel',
  amount: 0,
  currency: 'EUR',
  date: new Date().toISOString().slice(0, 10),
  note: ''
};

export default function BudgetPage() {
  const [entries, setEntries] = useState<BudgetEntry[]>([]);
  const [forecast, setForecast] = useState<ForecastResult | null>(null);
  const [form, setForm] = useState(formSeed);
  const [error, setError] = useState('');

  async function load() {
    try {
      const [entryPayload, forecastPayload] = await Promise.all([getBudgetEntries(), getForecast('trip-1')]);
      setEntries(entryPayload.entries);
      setForecast(forecastPayload);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load budget');
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await createBudgetEntry({ ...form, amount: Number(form.amount), category: form.category as 'living' | 'travel' });
      setForm((old) => ({ ...old, amount: 0, note: '' }));
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save entry');
    }
  }

  return (
    <div className="space-y-4">
      <SectionTitle title="Budget Tracker" subtitle="Keep travel and living expenses separate." />

      {forecast ? (
        <section className="card grid grid-cols-3 gap-2 text-center">
          <div>
            <p className="text-xs text-slate-500">Projected</p>
            <p className="text-lg font-semibold">EUR {forecast.projectedMonthlySpend.toFixed(0)}</p>
          </div>
          <div>
            <p className="text-xs text-slate-500">Remaining</p>
            <p className="text-lg font-semibold">EUR {forecast.remainingBudget.toFixed(0)}</p>
          </div>
          <div>
            <p className="text-xs text-slate-500">Status</p>
            <p className="text-lg font-semibold capitalize">{forecast.affordability}</p>
          </div>
        </section>
      ) : null}

      <form className="card space-y-3" onSubmit={onSubmit}>
        <div className="grid grid-cols-2 gap-3">
          <select
            className="input"
            value={form.category}
            onChange={(event) => setForm((old) => ({ ...old, category: event.target.value }))}
          >
            <option value="living">Living</option>
            <option value="travel">Travel</option>
          </select>
          <input
            className="input"
            type="number"
            min={0}
            placeholder="Amount"
            value={form.amount}
            onChange={(event) => setForm((old) => ({ ...old, amount: Number(event.target.value) }))}
          />
        </div>
        <input
          className="input"
          type="date"
          value={form.date}
          onChange={(event) => setForm((old) => ({ ...old, date: event.target.value }))}
        />
        <input
          className="input"
          placeholder="Note"
          value={form.note}
          onChange={(event) => setForm((old) => ({ ...old, note: event.target.value }))}
        />
        <button className="button" type="submit">
          Add entry
        </button>
      </form>

      <section className="card">
        <h2 className="text-lg font-semibold">Recent entries</h2>
        {entries.length === 0 ? <p className="mt-2 text-sm text-slate-600">No entries yet.</p> : null}
        <ul className="mt-3 space-y-2">
          {entries.map((entry) => (
            <li key={entry.id} className="flex items-center justify-between rounded-xl border border-slate-200 p-3 text-sm">
              <span className="capitalize">{entry.category}</span>
              <span>EUR {entry.amount.toFixed(0)}</span>
              <span className="text-slate-500">{entry.date}</span>
            </li>
          ))}
        </ul>
      </section>

      {error ? <p className="rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p> : null}
    </div>
  );
}
