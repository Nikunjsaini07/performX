import React from 'react';
import Link from 'next/link';
import type { Performance } from '@/lib/api';
import { toNumber } from '@/lib/format';
import TeamCrest from '@/components/TeamCrest';
import RatingBadge from '@/components/features/RatingBadge';
import PlayerTile from '@/components/features/PlayerTile';

interface PerformanceCardProps {
  performance: Performance;
  className?: string;
  /** Denser layout for home rails (smaller type + tighter padding). */
  compact?: boolean;
}

export default function PerformanceCard({ performance: p, className = '', compact = false }: PerformanceCardProps) {
  const rating = toNumber(p.average_rating) ?? toNumber(p.average_rating_2) ?? toNumber(p.provider_rating);
  const votes = toNumber(p.total_votes);
  const goals = toNumber(p.goals);
  const assists = toNumber(p.assists);
  const minutes = toNumber(p.minutes_played);
  // Prefer the transparent player headshot, fall back to a cover image, then initials.
  const photo = p.player_photo_url || p.cover_image_url;
  // Prefer the human-readable slug for URLs; fall back to the raw id.
  const href = `/performances/${p.slug || p.id}`;

  return (
    <Link href={href} className={`card-shell card-lift group flex flex-col ${className}`}>
      {/* Headshot tile */}
      <div className="card-media aspect-[4/5]">
        <PlayerTile
          photoUrl={photo}
          name={p.player_name}
          jerseyNumber={toNumber(p.jersey_number) ?? undefined}
          rounded="rounded-none"
          className="h-full w-full"
        />
        {rating != null && (
          <span className="absolute right-2.5 top-2.5 z-10 rounded-lg bg-black/45 px-1.5 py-0.5 backdrop-blur-sm">
            <RatingBadge value={rating} size="sm" />
          </span>
        )}
        {/* bottom fade for legibility */}
        <div className="pointer-events-none absolute inset-x-0 bottom-0 h-16 bg-gradient-to-t from-card to-transparent" />
      </div>

      {/* Body */}
      <div className={`flex flex-1 flex-col ${compact ? 'gap-1 p-3' : 'gap-2 p-4'}`}>
        {/* Match info header - show opponent */}
        {p.match_title && (
          <div className="flex items-center gap-1.5">
            <span className="text-[11px] font-medium text-muted-foreground">vs</span>
            <span className="line-clamp-1 text-xs font-semibold text-foreground">
              {/* Extract opponent from match title (e.g., "Argentina vs Egypt" -> "Egypt" if player is from Argentina) */}
              {p.match_title.split(' vs ').find(team => team !== p.team_name) || p.match_title}
            </span>
          </div>
        )}

        <div className="flex items-center gap-2">
          <TeamCrest
            name={p.team_name}
            logoUrl={p.team_logo_url}
            flagEmoji={p.flag_emoji || p.flag}
            size={compact ? 16 : 18}
          />
          <span className="line-clamp-1 text-xs font-medium text-muted-foreground">
            {p.team_name || '—'}
          </span>
        </div>

        <h3
          className={`line-clamp-1 font-display font-bold text-foreground transition-colors group-hover:text-primary ${
            compact ? 'text-sm' : 'text-base'
          }`}
        >
          {p.player_name}
        </h3>

        {!compact && p.tagline && (
          <p className="line-clamp-1 text-xs text-muted-foreground">{p.tagline}</p>
        )}

        {/* Stat footer */}
        <div className="mt-auto flex flex-wrap items-center gap-x-3 gap-y-1 pt-2 text-[11px] font-medium text-muted-foreground">
          {goals != null && goals > 0 && <span className="text-foreground">{goals} G</span>}
          {assists != null && assists > 0 && <span className="text-foreground">{assists} A</span>}
          {minutes != null && minutes > 0 && <span>{minutes}&apos;</span>}
          {votes != null && votes > 0 && (
            <span>{votes.toLocaleString()} {votes === 1 ? 'vote' : 'votes'}</span>
          )}
        </div>
      </div>
    </Link>
  );
}
