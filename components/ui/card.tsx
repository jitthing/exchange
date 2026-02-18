import { HTMLAttributes, ReactNode } from 'react';

type ShadowLevel = 'subtle' | 'medium' | 'raised';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  title?: string;
  subtitle?: string;
  shadow?: ShadowLevel;
  children: ReactNode;
}

const shadowClasses: Record<ShadowLevel, string> = {
  subtle: 'shadow-subtle',
  medium: 'shadow-medium',
  raised: 'shadow-raised',
};

export function Card({ title, subtitle, shadow = 'medium', children, className = '', ...props }: CardProps) {
  return (
    <div className={`rounded-lg bg-white p-5 ${shadowClasses[shadow]} ${className}`} {...props}>
      {title ? (
        <div className="mb-3">
          <h3 className="text-h3 text-heading">{title}</h3>
          {subtitle ? <p className="mt-0.5 text-small text-muted">{subtitle}</p> : null}
        </div>
      ) : null}
      {children}
    </div>
  );
}
