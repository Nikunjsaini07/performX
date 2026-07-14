import React from 'react';
import { Star } from 'lucide-react';
import { formatRating, toNumber } from '@/lib/format';

interface RatingBadgeProps {
  value: unknown;
  votes?: unknown;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  showStar?: boolean;
  className?: string;
}

const SIZES: Record<NonNullable<RatingBadgeProps['size']>, string> = {
  sm: 'text-xs px-1.5 py-0.5 gap-1',
  md: 'text-sm px-1.5 py-0.5 gap-1',
  lg: 'text-xl px-2.5 py-1 gap-1.5',
  xl: 'text-3xl px-3 py-1.5 gap-2',
};

const STAR_SIZE: Record<NonNullable<RatingBadgeProps['size']>, number> = {
  sm: 11,
  md: 12,
  lg: 16,
  xl: 20,
};

/**
 * IMDb-style amber rating. Amber is used ONLY for score numbers.
 */
export default function RatingBadge({
  value,
  votes,
  size = 'md',
  showStar = true,
  className = '',
}: RatingBadgeProps) {
  const n = toNumber(value);
  const voteCount = toNumber(votes);

  return (
    <span className={`inline-flex items-baseline ${className}`}>
      <span className={`inline-flex items-center rounded-lg font-bold stat-number text-rating ${SIZES[size]}`}>
        {showStar && <Star size={STAR_SIZE[size]} className="fill-rating text-rating shrink-0 self-center" />}
        {n == null ? '—' : formatRating(value)}
        {size !== 'sm' && size !== 'md' && (
          <span className="text-muted-foreground font-medium text-[0.4em] self-end mb-1">/10</span>
        )}
      </span>
      {voteCount != null && voteCount > 0 && (size === 'md' || size === 'lg' || size === 'xl') && (
        <span className="ml-2 text-xs font-medium text-muted-foreground self-center">
          {voteCount.toLocaleString()} {voteCount === 1 ? 'vote' : 'votes'}
        </span>
      )}
    </span>
  );
}
