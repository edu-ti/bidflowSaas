import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({ children, className }) => {
  return (
    <div
      className={`bg-gray-800 rounded-2xl shadow-soft p-6 ${className}`}
    >
      {children}
    </div>
  );
};
