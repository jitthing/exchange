'use client';

import { useEffect, useState } from 'react';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';
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

function affordabilityColor(a: string) {
  if (a === 'green') return 'bg-success-50';
  if (a === 'amber') return 'bg-warning-50';
  return 'bg-danger-50';
}

function affordabilityBadge(a: string) {
  if (a === 'green') return 'safe' as const;
  if (a === 'amber') return 'warning' as const;
  return 'danger' as const;
}

export default function BudgetPage() {
  const [entries, setEntries] = useState<BudgetEntry[]>([]);
  const [forecast, setForecast] = useState<ForecastResult | null>(null);
  const [form, setForm] = useState(formSeed);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  async function load() {
    try {
      const [entryPayload, forecastPayload] = await Promise.all([getBudgetEntries(), getForecast('trip-1')]);
      setEntries(entryPayload.entries);
      setForecast(forecastPayload);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load budget');
    }
  }

  useEffect(() => { load(); }, []);

  function validateBudget(): boolean {
    const errors: Record<string, string> = {};
    if (!form.category) errors.category = 'Category is required.';
    if (Number(form.amount) <= 0) errors.amount = 'Amount must be greater than 0.';
    if (!form.date) errors.date = 'Date is required.';
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!validateBudget()) return;
    try {
      await createBudgetEntry({ ...form, amount: Number(form.amount), category: form.category as 'living' | 'travel' });
      setForm((old) => ({ ...old, amount: 0, note: '' }));
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save entry');
    }
  }

  return (
    <div className="space-y-6">
      <SectionHeader title="Budget" subtitle="Track living and travel expenses" />

      {forecast ? (
        <Card className={`${affordabilityColor(forecast.affordability)} border-0`} shadow="subtle">
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <p className="text-caption text-muted">Projected</p>
              <p className="mt-1 text-h3 font-semibold text-heading">€{forecast.projectedMonthlySpend.toFixed(0)}</p>
            </div>
            <div>
              <p className="text-caption text-muted">Remaining</p>
              <p className="mt-1 text-h3 font-semibold text-heading">€{forecast.remainingBudget.toFixed(0)}</p>
            </div>
            <div>
              <p className="text-caption text-muted">Status</p>
              <div className="mt-2">
                <Badge variant={affordabilityBadge(forecast.affordability)}>
                  {forecast.affordability}
                </Badge>
              </div>
            </div>
          </div>
        </Card>
      ) : null}

      <Card title="Add Entry">
        <form className="space-y-3" onSubmit={onSubmit}>
          <div className="grid grid-cols-2 gap-3">
            <Select
              label="Category"
              value={form.category}
              error={fieldErrors.category}
              onChange={(e) => setForm((old) => ({ ...old, category: e.target.value }))}
            >
              <option value="living">Living</option>
              <option value="travel">Travel</option>
            </Select>
            <Input
              label="Amount (€)"
              type="number"
              min={0}
              value={form.amount}
              error={fieldErrors.amount}
              onChange={(e) => setForm((old) => ({ ...old, amount: Number(e.target.value) }))}
            />
          </div>
          <Input
            label="Date"
            type="date"
            value={form.date}
            error={fieldErrors.date}
            onChange={(e) => setForm((old) => ({ ...old, date: e.target.value }))}
          />
          <Input
            label="Note"
            placeholder="Optional note"
            value={form.note}
            onChange={(e) => setForm((old) => ({ ...old, note: e.target.value }))}
          />
          <Button type="submit" className="w-full">Add Entry</Button>
        </form>
      </Card>

      <div>
        <h2 className="mb-3 text-h3 text-heading">Recent Entries</h2>
        {entries.length === 0 ? (
          <Card shadow="subtle">
            <p className="text-center text-small text-muted">No entries yet</p>
          </Card>
        ) : (
          <div className="space-y-2">
            {entries.map((entry) => (
              <Card key={entry.id} shadow="subtle">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Badge variant={entry.category === 'travel' ? 'info' : 'safe'}>
                      {entry.category}
                    </Badge>
                    {entry.note ? <span className="text-small text-body">{entry.note}</span> : null}
                  </div>
                  <div className="text-right">
                    <p className="text-small font-semibold text-heading">€{entry.amount.toFixed(0)}</p>
                    <p className="text-caption text-muted">{entry.date}</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>

      {error ? <div className="badge-danger rounded-md p-3 text-small">{error}</div> : null}
    </div>
  );
}
