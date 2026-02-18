import { ReactNode } from 'react';

type BadgeVariant = 'safe' | 'warning' | 'danger' | 'info';

interface BadgeProps {
  variant: BadgeVariant;
  children: ReactNode;
}

const variantClasses: Record<BadgeVariant, string> = {
  safe: 'badge-safe',
  warning: 'badge-warning',
  danger: 'badge-danger',
  info: 'badge-info',
};

export function Badge({ variant, children }: BadgeProps) {
  return <span className={variantClasses[variant]}>{children}</span>;
}
