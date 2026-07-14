'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { ChevronLeft } from 'lucide-react';

interface BackButtonProps {
  className?: string;
  label?: string;
}

/**
 * Subtle "Back" control for detail pages. Uses browser history so it returns
 * the user to wherever they came from (list, home rail, profile, etc.).
 */
export default function BackButton({ className = '', label = 'Back' }: BackButtonProps) {
  const router = useRouter();

  return (
    <button
      type="button"
      onClick={() => router.back()}
      className={`inline-flex items-center gap-1 rounded-full border border-border bg-surface/60 px-3 py-1.5 text-sm font-medium text-muted-foreground backdrop-blur-sm transition-colors hover:border-primary/40 hover:text-foreground ${className}`}
    >
      <ChevronLeft size={16} />
      {label}
    </button>
  );
}
