import { notFound } from 'next/navigation';
import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { API_BASE_URL } from '@/lib/config';

type Trip = {
  id: string;
  destination: string;
  itinerary: string[];
  members: string[];
  estimatedCost: number;
};

type PageProps = {
  params: { tripId: string };
};

async function fetchTrip(tripId: string): Promise<Trip | null> {
  const response = await fetch(`${API_BASE_URL}/api/trips/${tripId}`, { cache: 'no-store' });
  if (response.status === 404) return null;
  if (!response.ok) throw new Error('Unable to load trip');
  return (await response.json()) as Trip;
}

export default async function TripDetailPage({ params }: PageProps) {
  const trip = await fetchTrip(params.tripId);
  if (!trip) return notFound();

  return (
    <div className="space-y-6">
      <SectionHeader title={`${trip.destination}`} subtitle="Trip details and itinerary" />

      <Card shadow="raised">
        <div className="flex items-center justify-between">
          <p className="text-caption font-medium uppercase tracking-wider text-muted">Estimated Cost</p>
          <p className="text-h2 font-semibold text-heading">€{trip.estimatedCost.toFixed(0)}</p>
        </div>
      </Card>

      <div>
        <h2 className="mb-3 text-h3 text-heading">Itinerary</h2>
        <div className="space-y-0">
          {trip.itinerary.map((item, i) => (
            <div key={item} className="flex gap-3">
              <div className="flex flex-col items-center">
                <div className="flex h-7 w-7 items-center justify-center rounded-full bg-primary text-caption font-semibold text-white">
                  {i + 1}
                </div>
                {i < trip.itinerary.length - 1 ? (
                  <div className="w-0.5 flex-1 bg-neutral-200" />
                ) : null}
              </div>
              <Card shadow="subtle" className="mb-3 flex-1">
                <p className="text-small text-body">{item}</p>
              </Card>
            </div>
          ))}
        </div>
      </div>

      <div>
        <h2 className="mb-3 text-h3 text-heading">Members</h2>
        <div className="flex flex-wrap gap-2">
          {trip.members.map((member) => (
            <span
              key={member}
              className="inline-flex items-center rounded-full bg-primary-50 px-3 py-1.5 text-small font-medium text-primary"
            >
              👤 {member}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
