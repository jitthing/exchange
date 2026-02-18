'use client';

import { SelectHTMLAttributes, forwardRef } from 'react';

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, error, children, className = '', id, ...props }, ref) => {
    const selectId = id || label?.toLowerCase().replace(/\s+/g, '-');
    return (
      <div className="space-y-1">
        {label ? (
          <label htmlFor={selectId} className="block text-small font-medium text-heading">
            {label}
          </label>
        ) : null}
        <select
          ref={ref}
          id={selectId}
          className={`input appearance-none bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2212%22%20height%3D%2212%22%20viewBox%3D%220%200%2012%2012%22%3E%3Cpath%20fill%3D%22%239B9B9B%22%20d%3D%22M6%208L1%203h10z%22%2F%3E%3C%2Fsvg%3E')] bg-[length:12px] bg-[right_12px_center] bg-no-repeat pr-8 ${error ? 'border-danger' : ''} ${className}`}
          {...props}
        >
          {children}
        </select>
        {error ? <p className="text-caption text-danger-600">{error}</p> : null}
      </div>
    );
  }
);

Select.displayName = 'Select';
