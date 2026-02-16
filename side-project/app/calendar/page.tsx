'use client';

import { useEffect, useState } from 'react';
import { SectionTitle } from '@/components/section-title';
import { evaluateConflicts, getTravelWindows } from '@/lib/api';
import { ConflictAlert, TravelWindow } from '@/lib/types';

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
    <div className="space-y-4">
      <SectionTitle title="Calendar Sync" subtitle="Detected travel windows and study overlap checks." />

      {error ? <p className="rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p> : null}

      <section className="card space-y-3">
        <label className="text-sm font-medium" htmlFor="window">
          Travel window
        </label>
        <select
          id="window"
          className="input"
          value={selected}
          onChange={(event) => setSelected(event.target.value)}
        >
          {windows.map((window) => (
            <option key={window.id} value={window.id}>
              {window.startDate} to {window.endDate} (score {window.score})
            </option>
          ))}
        </select>
      </section>

      <section className="card">
        <h2 className="text-lg font-semibold">Conflict Alerts</h2>
        {alerts.length === 0 ? (
          <p className="mt-2 text-sm text-slate-600">No conflicts in selected window.</p>
        ) : (
          <ul className="mt-3 space-y-2">
            {alerts.map((alert) => (
              <li key={alert.relatedEventId} className="rounded-xl border border-slate-200 p-3">
                <p className="text-sm font-semibold capitalize">{alert.severity}</p>
                <p className="text-sm text-slate-600">{alert.reason}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
