'use client';

import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', id, ...props }, ref) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-');
    return (
      <div className="space-y-1">
        {label ? (
          <label htmlFor={inputId} className="block text-small font-medium text-heading">
            {label}
          </label>
        ) : null}
        <input
          ref={ref}
          id={inputId}
          className={`input ${error ? 'border-danger' : ''} ${className}`}
          {...props}
        />
        {error ? <p data-testid="field-error" className="text-caption text-danger-600">{error}</p> : null}
      </div>
    );
  }
);

Input.displayName = 'Input';
