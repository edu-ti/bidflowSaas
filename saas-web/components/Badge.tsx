import React from 'react';

interface BadgeProps {
  variant: 'high' | 'medium' | 'low';
  children?: React.ReactNode;
  className?: string;
}

const variantClasses = {
  high: 'bg-red-600 text-white',
  medium: 'bg-amber-500 text-white',
  low: 'bg-green-600 text-white',
};

export const Badge: React.FC<BadgeProps> = ({ variant, children, className }) => {
  return (
    <span
      className={`px-2 py-1 text-xs font-medium rounded ${variantClasses[variant]} ${className}`}
    >
      {children}
    </span>
  );
};
