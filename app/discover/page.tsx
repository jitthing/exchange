'use client';

import { useState } from 'react';
import Link from 'next/link';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
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

function riskBadgeVariant(risk: string) {
  if (risk === 'high-risk') return 'danger' as const;
  if (risk === 'warning') return 'warning' as const;
  return 'safe' as const;
}

export default function DiscoverPage() {
  const [form, setForm] = useState(initialForm);
  const [results, setResults] = useState<TripOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  function validate(): boolean {
    const errors: Record<string, string> = {};
    if (!form.departureCity.trim()) errors.departureCity = 'Departure city is required.';
    if (form.budgetCap <= 0) errors.budgetCap = 'Budget cap must be greater than 0.';
    if (form.maxTravelHours <= 0) errors.maxTravelHours = 'Max travel hours must be greater than 0.';
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!validate()) return;
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
    <div className="space-y-6">
      <SectionHeader title="Discover" subtitle="Find your next weekend trip" />

      <Card>
        <form className="space-y-4" onSubmit={onSubmit}>
          <Input
            label="Departure city"
            value={form.departureCity}
            error={fieldErrors.departureCity}
            onChange={(e) => setForm((old) => ({ ...old, departureCity: e.target.value }))}
          />
          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Budget cap (€)"
              type="number"
              value={form.budgetCap}
              error={fieldErrors.budgetCap}
              onChange={(e) => setForm((old) => ({ ...old, budgetCap: Number(e.target.value) }))}
            />
            <Input
              label="Max travel hours"
              type="number"
              value={form.maxTravelHours}
              error={fieldErrors.maxTravelHours}
              onChange={(e) => setForm((old) => ({ ...old, maxTravelHours: Number(e.target.value) }))}
            />
          </div>
          <Button type="submit" variant="accent" loading={loading} className="w-full">
            🧭 Find Trips
          </Button>
          {error ? <p className="text-small text-danger">{error}</p> : null}
        </form>
      </Card>

      <div className="space-y-3">
        {results.map((option) => (
          <Card key={option.id} shadow="medium">
            <div className="flex items-start justify-between">
              <div>
                <h2 className="text-h2 text-heading">{option.destination}</h2>
                <div className="mt-2 flex flex-wrap gap-1.5">
                  {option.reasonTags.map((tag) => (
                    <Badge key={tag} variant="info">{tag}</Badge>
                  ))}
                </div>
              </div>
              <div className="text-right">
                <p className="text-h3 font-semibold text-heading">€{option.totalEstimatedCost.toFixed(0)}</p>
                <Badge variant={riskBadgeVariant(option.riskLevel)}>{option.riskLevel}</Badge>
              </div>
            </div>

            {option.transportOptions[0] ? (
              <div className="mt-3 flex items-center gap-2 text-small text-muted">
                <span>🚆</span>
                <span>
                  {option.transportOptions[0].provider} · {option.transportOptions[0].durationHours}h · €{option.transportOptions[0].price}
                </span>
              </div>
            ) : null}

            <Link
              href="/trips/trip-1"
              className="mt-3 inline-block text-small font-medium text-primary hover:underline"
            >
              View details →
            </Link>
          </Card>
        ))}
      </div>
    </div>
  );
}
