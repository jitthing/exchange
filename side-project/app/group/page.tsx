'use client';

import { useState } from 'react';
import { SectionTitle } from '@/components/section-title';
import { getTrip, shareTrip } from '@/lib/api';
import { Trip } from '@/lib/types';

export default function GroupPage() {
  const [trip, setTrip] = useState<Trip | null>(null);
  const [member, setMember] = useState('friend-1');
  const [error, setError] = useState('');

  async function loadTrip() {
    try {
      setTrip(await getTrip('trip-1'));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load trip');
    }
  }

  async function addMember() {
    if (!member.trim()) return;
    try {
      const updated = await shareTrip('trip-1', [member.trim()]);
      setTrip(updated);
      setMember('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to share trip');
    }
  }

  return (
    <div className="space-y-4">
      <SectionTitle title="Group Planning" subtitle="Share a trip and keep everyone on the same page." />

      <section className="card space-y-3">
        <button className="button" onClick={loadTrip}>
          Load shared trip
        </button>

        {trip ? (
          <>
            <h2 className="text-lg font-semibold">{trip.destination}</h2>
            <p className="text-sm text-slate-600">Estimated total EUR {trip.estimatedCost.toFixed(0)}</p>
            <ul className="space-y-2">
              {trip.members.map((item) => (
                <li key={item} className="rounded-xl border border-slate-200 px-3 py-2 text-sm">
                  {item}
                </li>
              ))}
            </ul>
            <div className="flex gap-2">
              <input className="input" value={member} onChange={(event) => setMember(event.target.value)} placeholder="member id" />
              <button className="button-secondary" onClick={addMember}>
                Add
              </button>
            </div>
          </>
        ) : null}
      </section>

      {error ? <p className="rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p> : null}
    </div>
  );
}
