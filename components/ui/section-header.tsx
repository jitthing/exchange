interface SectionHeaderProps {
  title: string;
  subtitle?: string;
}

export function SectionHeader({ title, subtitle }: SectionHeaderProps) {
  return (
    <header className="mb-6">
      <h1 className="text-h1 text-heading">{title}</h1>
      {subtitle ? <p className="mt-1 text-body text-muted">{subtitle}</p> : null}
    </header>
  );
}
