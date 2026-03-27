import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({ children, className }) => {
  return (
    <div
      className={`bg-gray-800/50 backdrop-blur-sm border border-gray-700/50 rounded-2xl shadow-soft p-6 hover:border-emerald-500/30 transition-all duration-300 ${className}`}
    >
      {children}
    </div>
  );
};
