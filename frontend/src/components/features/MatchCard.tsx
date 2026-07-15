import React from 'react';
import Link from 'next/link';
import { CalendarDays } from 'lucide-react';
import type { Match } from '@/lib/api';
import { formatDate, penaltyValue, toNumber } from '@/lib/format';
import TeamCrest from '@/components/TeamCrest';
import RatingBadge from '@/components/features/RatingBadge';

interface MatchCardProps {
  match: Match;
  className?: string;
}

export default function MatchCard({ match, className = '' }: MatchCardProps) {
  const homePens = penaltyValue(match.home_penalty_score);
  const awayPens = penaltyValue(match.away_penalty_score);
  const hasPens = homePens != null && awayPens != null;
  const rating = toNumber(match.average_rating);
  const votes = toNumber(match.total_votes);

  return (
    <Link
      href={`/matches/${match.slug}`}
      className={`card-shell card-lift group flex flex-col ${className}`}
    >
      {/* Top strip: round + date */}
      <div className="flex items-center justify-between gap-2 border-b border-border/70 px-4 py-2.5">
        <span className="stage-badge max-w-[60%] truncate">{match.round || 'Match'}</span>
        <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
          <CalendarDays size={13} />
          {formatDate(match.utc_datetime)}
        </span>
      </div>

      {/* Matchup */}
      <div className="flex items-center justify-between gap-3 px-4 py-5">
        <div className="flex min-w-0 flex-1 flex-col items-center gap-2 text-center">
          <TeamCrest
            name={match.home_team_name}
            shortName={match.home_team_short_name}
            logoUrl={match.home_team_logo_url}
            size={44}
          />
          <span className="line-clamp-1 w-full text-xs font-semibold text-foreground">
            {match.home_team_short_name || match.home_team_name}
          </span>
        </div>

        <div className="flex shrink-0 flex-col items-center">
          <div className="flex items-baseline gap-2 font-display text-3xl font-bold stat-number text-foreground">
            <span className="inline-flex items-baseline gap-0.5">
              <span>{toNumber(match.home_score) ?? 0}</span>
              {hasPens && <span className="text-sm font-medium text-muted-foreground">({homePens})</span>}
            </span>
            <span className="text-muted-foreground/50">–</span>
            <span className="inline-flex items-baseline gap-0.5">
              <span>{toNumber(match.away_score) ?? 0}</span>
              {hasPens && <span className="text-sm font-medium text-muted-foreground">({awayPens})</span>}
            </span>
          </div>
        </div>

        <div className="flex min-w-0 flex-1 flex-col items-center gap-2 text-center">
          <TeamCrest
            name={match.away_team_name}
            shortName={match.away_team_short_name}
            logoUrl={match.away_team_logo_url}
            size={44}
          />
          <span className="line-clamp-1 w-full text-xs font-semibold text-foreground">
            {match.away_team_short_name || match.away_team_name}
          </span>
        </div>
      </div>

      {/* Footer: title/tagline + rating */}
      <div className="mt-auto flex items-center justify-between gap-3 border-t border-border/70 px-4 py-3">
        <p className="line-clamp-1 text-sm text-muted-foreground transition-colors group-hover:text-foreground">
          {match.tagline || match.title}
        </p>
        {rating != null ? (
          <span className="flex shrink-0 items-center gap-1.5">
            <RatingBadge value={match.average_rating} size="md" />
            {votes != null && votes > 0 && (
              <span className="text-[11px] font-medium text-muted-foreground">
                · {votes.toLocaleString()} {votes === 1 ? 'vote' : 'votes'}
              </span>
            )}
          </span>
        ) : (
          <span className="text-xs text-muted-foreground/60">Not rated</span>
        )}
      </div>
    </Link>
  );
}
