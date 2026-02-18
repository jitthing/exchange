import { SectionHeader } from '@/components/ui/section-header';
import { Card } from '@/components/ui/card';
import { API_BASE_URL } from '@/lib/config';

const settings = [
  { label: 'API Base URL', value: API_BASE_URL },
  { label: 'Region Focus', value: 'Europe (Schengen)' },
  { label: 'User Profile', value: 'demo-user' },
  { label: 'Booking Behavior', value: 'Deep-link to provider websites' },
];

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <SectionHeader title="Settings" subtitle="App configuration" />

      <Card>
        <div className="divide-y divide-neutral-200">
          {settings.map((item) => (
            <div key={item.label} className="flex items-center justify-between py-3 first:pt-0 last:pb-0">
              <span className="text-small text-muted">{item.label}</span>
              <span className="text-small font-medium text-heading">{item.value}</span>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
