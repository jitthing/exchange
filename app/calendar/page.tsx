'use client';

import { useEffect, useState } from 'react';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select } from '@/components/ui/select';
import { evaluateConflicts, getTravelWindows } from '@/lib/api';
import { ConflictAlert, TravelWindow } from '@/lib/types';

function scoreBadgeVariant(score: number) {
  if (score >= 80) return 'safe' as const;
  if (score >= 50) return 'warning' as const;
  return 'danger' as const;
}

function severityBadgeVariant(severity: string) {
  if (severity === 'high-risk') return 'danger' as const;
  if (severity === 'warning') return 'warning' as const;
  return 'info' as const;
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' });
}

export default function CalendarPage() {
  const [windows, setWindows] = useState<TravelWindow[]>([]);
  const [selected, setSelected] = useState<string>('');
  const [alerts, setAlerts] = useState<ConflictAlert[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    getTravelWindows()
      .then((payload) => {
        setWindows(payload.windows);
        if (payload.windows[0]) setSelected(payload.windows[0].id);
      })
      .catch((err: Error) => setError(err.message));
  }, []);

  useEffect(() => {
    if (!selected) return;
    evaluateConflicts(selected)
      .then((payload) => setAlerts(payload.alerts))
      .catch((err: Error) => setError(err.message));
  }, [selected]);

  return (
    <div className="space-y-6">
      <SectionHeader title="Calendar" subtitle="Travel windows and academic conflicts" />

      {error ? (
        <div className="badge-danger rounded-md p-3 text-small">{error}</div>
      ) : null}

      <Card>
        <Select
          label="Travel window"
          value={selected}
          onChange={(e) => setSelected(e.target.value)}
        >
          {windows.map((w) => (
            <option key={w.id} value={w.id}>
              {formatDate(w.startDate)} — {formatDate(w.endDate)}
            </option>
          ))}
        </Select>
      </Card>

      {/* Window cards */}
      <div className="space-y-3">
        {windows.map((w) => (
          <Card
            key={w.id}
            className={`cursor-pointer transition-all ${selected === w.id ? 'ring-2 ring-primary/30' : ''}`}
            shadow={selected === w.id ? 'raised' : 'subtle'}
            onClick={() => setSelected(w.id)}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-h3 text-heading">
                  {formatDate(w.startDate)} — {formatDate(w.endDate)}
                </p>
                <p className="mt-0.5 text-caption text-muted">{w.conflicts.length} potential conflicts</p>
              </div>
              <Badge variant={scoreBadgeVariant(w.score)}>{w.score}</Badge>
            </div>
          </Card>
        ))}
      </div>

      {/* Conflict alerts */}
      <div>
        <h2 className="mb-3 text-h3 text-heading">Conflict Alerts</h2>
        {alerts.length === 0 ? (
          <Card shadow="subtle">
            <p className="text-center text-small text-muted">✅ No conflicts in this window</p>
          </Card>
        ) : (
          <div className="space-y-2">
            {alerts.map((alert) => (
              <Card key={alert.relatedEventId} shadow="subtle">
                <div className="flex items-start justify-between gap-3">
                  <p className="text-small text-body">{alert.reason}</p>
                  <Badge variant={severityBadgeVariant(alert.severity)}>{alert.severity}</Badge>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
