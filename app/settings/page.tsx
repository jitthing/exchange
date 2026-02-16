import { SectionTitle } from '@/components/section-title';
import { API_BASE_URL } from '@/lib/config';

export default function SettingsPage() {
  return (
    <div className="space-y-4">
      <SectionTitle title="Settings" subtitle="Environment and data-source defaults." />
      <section className="card space-y-2 text-sm text-slate-700">
        <p>
          API base URL: <span className="font-semibold">{API_BASE_URL}</span>
        </p>
        <p>Region focus: Europe (Schengen)</p>
        <p>Default user profile: demo-user</p>
        <p>Booking behavior: deep-link to provider websites</p>
      </section>
    </div>
  );
}
