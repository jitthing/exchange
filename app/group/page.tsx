'use client';

import { useState } from 'react';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { getTrip, shareTrip } from '@/lib/api';
import { Trip } from '@/lib/types';

export default function GroupPage() {
  const [trip, setTrip] = useState<Trip | null>(null);
  const [member, setMember] = useState('friend-1');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function loadTrip() {
    setLoading(true);
    try {
      setTrip(await getTrip('trip-1'));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load trip');
    } finally {
      setLoading(false);
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
    <div className="space-y-6">
      <SectionHeader title="Group" subtitle="Plan trips together" />

      {!trip ? (
        <Card shadow="medium" className="text-center">
          <p className="mb-4 text-body text-muted">Load a shared trip to get started</p>
          <Button onClick={loadTrip} loading={loading}>Load Shared Trip</Button>
        </Card>
      ) : (
        <div className="space-y-4">
          <Card shadow="raised">
            <div className="flex items-start justify-between">
              <div>
                <h2 className="text-h2 text-heading">{trip.destination}</h2>
                <p className="mt-1 text-small text-muted">
                  Estimated total <span className="font-semibold text-heading">€{trip.estimatedCost.toFixed(0)}</span>
                </p>
              </div>
              <Badge variant="info">{trip.members.length} members</Badge>
            </div>

            <div className="mt-4">
              <p className="mb-2 text-caption font-medium uppercase tracking-wider text-muted">Members</p>
              <div className="flex flex-wrap gap-2">
                {trip.members.map((m) => (
                  <span
                    key={m}
                    className="inline-flex items-center rounded-full bg-primary-50 px-3 py-1 text-small font-medium text-primary"
                  >
                    👤 {m}
                  </span>
                ))}
              </div>
            </div>
          </Card>

          <Card title="Add Member">
            <div className="flex gap-2">
              <div className="flex-1">
                <Input
                  value={member}
                  onChange={(e) => setMember(e.target.value)}
                  placeholder="Enter member ID"
                />
              </div>
              <Button variant="secondary" onClick={addMember}>Add</Button>
            </div>
          </Card>
        </div>
      )}

      {error ? <div className="badge-danger rounded-md p-3 text-small">{error}</div> : null}
    </div>
  );
}
