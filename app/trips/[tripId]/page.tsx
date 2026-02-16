import { notFound } from 'next/navigation';
import { SectionTitle } from '@/components/section-title';
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
    <div className="space-y-4">
      <SectionTitle title={`${trip.destination} Trip`} subtitle="Shared itinerary, members, and estimated spend." />
      <section className="card">
        <p className="text-sm text-slate-600">Estimated cost: EUR {trip.estimatedCost.toFixed(0)}</p>
      </section>
      <section className="card">
        <h2 className="text-lg font-semibold">Itinerary</h2>
        <ul className="mt-2 space-y-2 text-sm text-slate-700">
          {trip.itinerary.map((item) => (
            <li key={item} className="rounded-xl border border-slate-200 px-3 py-2">
              {item}
            </li>
          ))}
        </ul>
      </section>
      <section className="card">
        <h2 className="text-lg font-semibold">Members</h2>
        <ul className="mt-2 flex flex-wrap gap-2 text-sm">
          {trip.members.map((member) => (
            <li key={member} className="rounded-full bg-slate-100 px-3 py-1">
              {member}
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
