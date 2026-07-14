import React from 'react';
import { initials } from '@/lib/format';

interface PlayerTileProps {
  photoUrl?: string;
  name?: string;
  jerseyNumber?: number;
  className?: string;
  rounded?: string;
}

/**
 * Premium headshot-on-gradient tile. We only have transparent player headshot
 * PNGs, so the player sits bottom-aligned on a team/accent-tinted radial
 * gradient with subtle grain — like a Sorare / Apple TV person tile.
 */
export default function PlayerTile({
  photoUrl,
  name,
  jerseyNumber,
  className = '',
  rounded = 'rounded-[14px]',
}: PlayerTileProps) {
  return (
    <div className={`player-tile ${rounded} ${className}`}>
      {jerseyNumber != null && (
        <span className="absolute top-3 left-3 z-10 text-[11px] font-bold stat-number text-foreground/70 bg-black/30 backdrop-blur-sm rounded-md px-1.5 py-0.5 border border-white/10">
          #{jerseyNumber}
        </span>
      )}

      {photoUrl ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={photoUrl}
          alt={name || 'Player'}
          className="absolute inset-0 h-full w-full object-contain object-bottom drop-shadow-[0_8px_24px_rgba(0,0,0,0.5)]"
          loading="lazy"
        />
      ) : (
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-4xl font-bold font-display text-foreground/25">{initials(name)}</span>
        </div>
      )}
    </div>
  );
}
