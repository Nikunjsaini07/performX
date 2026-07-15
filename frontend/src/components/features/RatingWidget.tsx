'use client';

import React, { useEffect, useState } from 'react';
import { Loader2, Check, Star } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import {
  submitMatchRating,
  submitPerformanceRating,
  getMyMatchRating,
  getMyPerformanceRating,
  type RatingResponse,
} from '@/lib/api';
import { toNumber, formatRating } from '@/lib/format';
import AuthModal from '@/components/features/auth/AuthModal';

interface RatingWidgetProps {
  kind: 'match' | 'performance';
  /** slug for match, id for performance */
  entityKey: string;
  initialAverage?: unknown;
  initialVotes?: unknown;
}

export default function RatingWidget({ kind, entityKey, initialAverage, initialVotes }: RatingWidgetProps) {
  const { token, isAuthenticated } = useAuth();

  const [average, setAverage] = useState<number | null>(toNumber(initialAverage));
  const [votes, setVotes] = useState<number | null>(toNumber(initialVotes));
  const [myRating, setMyRating] = useState<number | null>(null);
  const [hover, setHover] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [justSaved, setJustSaved] = useState(false);
  const [error, setError] = useState('');
  const [authOpen, setAuthOpen] = useState(false);

  useEffect(() => {
    let active = true;
    if (!isAuthenticated || !token) {
      setMyRating(null);
      return;
    }
    (async () => {
      try {
        const res =
          kind === 'match'
            ? await getMyMatchRating(entityKey, token)
            : await getMyPerformanceRating(entityKey, token);
        if (active && res) setMyRating(toNumber(res.rating));
      } catch {
        /* no existing rating */
      }
    })();
    return () => {
      active = false;
    };
  }, [isAuthenticated, token, kind, entityKey]);

  const handleRate = async (value: number) => {
    if (!isAuthenticated || !token) {
      setAuthOpen(true);
      return;
    }
    setError('');
    setSubmitting(true);
    const prev = myRating;
    setMyRating(value); // optimistic
    try {
      const res: RatingResponse =
        kind === 'match'
          ? await submitMatchRating(entityKey, value, token)
          : await submitPerformanceRating(entityKey, value, token);
      const newAvg = toNumber(res.average_rating);
      const newVotes = toNumber(res.total_votes);
      if (newAvg != null) setAverage(newAvg);
      if (newVotes != null) setVotes(newVotes);
      setJustSaved(true);
      setTimeout(() => setJustSaved(false), 1800);
    } catch (err) {
      setMyRating(prev);
      setError(err instanceof Error ? err.message : 'Could not save rating');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="py-6 lg:py-0 lg:px-0">
      <div className="mb-4 flex items-center justify-between gap-3">
        <div>
          <h3 className="font-display text-lg font-bold text-foreground">Your Rating</h3>
          <p className="text-xs text-muted-foreground">
            {isAuthenticated ? 'Tap a number to rate out of 10' : 'Sign in to rate this ' + kind}
          </p>
        </div>
        <div className="text-right">
          <div className="flex items-center gap-1 font-display text-2xl font-bold stat-number text-rating">
            <Star size={18} className="fill-rating text-rating" />
            {average != null ? formatRating(average) : '—'}
          </div>
          <p className="text-[11px] text-muted-foreground">
            {votes != null && votes > 0 ? `${votes.toLocaleString()} ${votes === 1 ? 'vote' : 'votes'}` : 'No votes yet'}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-10 gap-1.5">
        {Array.from({ length: 10 }, (_, i) => i + 1).map((n) => {
          const activeVal = hover ?? myRating ?? 0;
          const isSelected = myRating === n;
          const isLit = n <= activeVal;
          return (
            <button
              key={n}
              disabled={submitting}
              onMouseEnter={() => setHover(n)}
              onMouseLeave={() => setHover(null)}
              onClick={() => handleRate(n)}
              className={`flex h-11 items-center justify-center rounded-lg border text-sm font-bold stat-number transition-all disabled:opacity-60 ${
                isSelected
                  ? 'border-rating bg-rating/20 text-rating'
                  : isLit
                    ? 'border-rating/40 bg-rating/10 text-rating'
                    : 'border-border bg-surface-2 text-muted-foreground hover:border-rating/40 hover:text-foreground'
              }`}
              aria-label={`Rate ${n} out of 10`}
            >
              {n}
            </button>
          );
        })}
      </div>

      <div className="mt-3 flex min-h-[20px] items-center gap-2 text-xs">
        {submitting && (
          <span className="inline-flex items-center gap-1.5 text-muted-foreground">
            <Loader2 size={13} className="animate-spin" /> Saving…
          </span>
        )}
        {justSaved && !submitting && (
          <span className="inline-flex items-center gap-1.5 text-primary">
            <Check size={13} /> Rating saved{myRating != null ? ` — you rated ${myRating}/10` : ''}
          </span>
        )}
        {!submitting && !justSaved && myRating != null && (
          <span className="text-muted-foreground">You rated this {kind} <span className="font-semibold text-rating">{myRating}/10</span></span>
        )}
        {error && <span className="text-red-400">{error}</span>}
      </div>

      <AuthModal isOpen={authOpen} onClose={() => setAuthOpen(false)} defaultTab="login" />
    </div>
  );
}
