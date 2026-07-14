import React from 'react';
import { initials } from '@/lib/format';

interface TeamCrestProps {
  name?: string;
  shortName?: string;
  logoUrl?: string;
  flagEmoji?: string;
  size?: number;
  className?: string;
}

/**
 * Team crest with a graceful fallback chain:
 * logo_url (Cloudinary flag PNG) → flag_emoji (unicode) → initials monogram.
 */
export default function TeamCrest({
  name,
  shortName,
  logoUrl,
  flagEmoji,
  size = 40,
  className = '',
}: TeamCrestProps) {
  const label = name || shortName || 'Team';
  const dimension = { width: size, height: size };

  if (logoUrl) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={logoUrl}
        alt={`${label} crest`}
        style={dimension}
        className={`object-contain ${className}`}
        loading="lazy"
      />
    );
  }

  if (flagEmoji) {
    return (
      <span
        style={{ fontSize: size * 0.82, lineHeight: 1, width: size, height: size }}
        className={`inline-flex items-center justify-center ${className}`}
        aria-label={`${label} flag`}
        role="img"
      >
        {flagEmoji}
      </span>
    );
  }

  return (
    <span
      style={{ ...dimension, fontSize: size * 0.36 }}
      className={`inline-flex items-center justify-center rounded-full bg-surface-2 font-bold text-muted-foreground border border-border ${className}`}
      aria-label={label}
    >
      {initials(label)}
    </span>
  );
}
