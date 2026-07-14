import React from 'react';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import { Clock, Shirt, ArrowUpRight } from 'lucide-react';
import { getPerformance, getPerformanceStats, getPerformanceReviews, getMatch } from '@/lib/data';
import { toNumber } from '@/lib/format';
import type { StatRow } from '@/lib/api';
import DetailBackdrop from '@/components/DetailBackdrop';
import BackButton from '@/components/BackButton';
import TeamCrest from '@/components/TeamCrest';
import RatingBadge from '@/components/RatingBadge';
import RatingWidget from '@/components/RatingWidget';
import ReviewsSection from '@/components/ReviewsSection';
import PlayerTile from '@/components/PlayerTile';

export const revalidate = 60;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}): Promise<Metadata> {
  const { id } = await Promise.resolve(params);
  const p = await getPerformance(id);
  return {
    title: p ? `${p.player_name} — Performance` : 'Performance',
    description: p?.description,
  };
}

function prettifyStat(row: StatRow): string {
  const raw = row.stat_name || row.stat_short_name || 'Stat';
  return raw
    .replace(/[_-]+/g, ' ')
    .replace(/\b\w/g, (c) => c.toUpperCase())
    .trim();
}

function statValue(row: StatRow): string {
  const n = toNumber(row.value);
  const val = n != null ? (Number.isInteger(n) ? String(n) : n.toFixed(1)) : String(row.value ?? '—');
  return row.stat_unit ? `${val}${row.stat_unit === '%' ? '' : ' '}${row.stat_unit}` : val;
}

export default async function PerformanceDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id: idOrSlug } = await Promise.resolve(params);

  // The route accepts either the raw id or the human-readable slug — the
  // backend resolves whichever was given to the same performance record.
  const performance = await getPerformance(idOrSlug);
  if (!performance) notFound();

  // Sub-resource calls (stats/reviews/rating) always need the internal id.
  const id = performance.id;

  const [stats, reviews] = await Promise.all([
    getPerformanceStats(id),
    getPerformanceReviews(id),
  ]);

  const p = performance;
  const rating = toNumber(p.average_rating) ?? toNumber(p.average_rating_2) ?? toNumber(p.provider_rating);
  const minutes = toNumber(p.minutes_played);
  const jersey = toNumber(p.jersey_number);

  return (
    <main>
      {/* ─── Cinematic header ──────────────────────────────────── */}
      <section className="relative">
        <DetailBackdrop />
        <div className="container-max container-px relative pt-32">
          <BackButton className="mb-6" />
        </div>
        <div className="container-max container-px relative flex flex-col items-start gap-6 sm:flex-row sm:items-end sm:gap-8">
          {/* Headshot tile */}
          <div className="w-40 shrink-0 sm:w-52">
            <div className="card-shell overflow-hidden">
              <PlayerTile
                photoUrl={p.player_photo_url}
                name={p.player_name}
                jerseyNumber={jersey ?? undefined}
                rounded="rounded-none"
                className="aspect-[4/5] w-full"
              />
            </div>
          </div>

          {/* Meta */}
          <div className="min-w-0 flex-1 pb-1">
            <div className="mb-3 flex flex-wrap items-center gap-2">
              {p.team_name && (
                <span className="inline-flex items-center gap-1.5 rounded-full border border-border bg-surface/60 px-2.5 py-1 text-xs font-medium text-foreground backdrop-blur-sm">
                  <TeamCrest name={p.team_name} logoUrl={p.team_logo_url} flagEmoji={p.flag_emoji || p.flag} size={16} />
                  {p.team_name}
                </span>
              )}
              {p.match_round && <span className="stage-badge">{p.match_round}</span>}
              {/* Show opponent */}
              {p.match_title && (
                <span className="inline-flex items-center gap-1 rounded-full border border-border bg-surface/60 px-2.5 py-1 text-xs font-medium text-muted-foreground backdrop-blur-sm">
                  <span className="text-[10px]">vs</span>
                  <span className="font-semibold text-foreground">
                    {p.match_title.split(' vs ').find(team => team !== p.team_name) || p.match_title}
                  </span>
                </span>
              )}
            </div>

            <h1 className="font-display text-[clamp(2rem,5vw,3.25rem)] font-bold leading-tight text-foreground">
              {p.player_name}
            </h1>

            {p.tagline && (
              <p className="mt-2 font-display text-lg font-semibold leading-snug text-foreground/90">
                {p.tagline}
              </p>
            )}

            {(p.match_slug) && (
              <Link
                href={`/matches/${p.match_slug}`}
                className="group mt-2 inline-flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-primary"
              >
                View full match
                <ArrowUpRight size={14} className="transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
              </Link>
            )}

            {/* Rating + quick meta */}
            <div className="mt-6 flex flex-wrap items-center gap-x-6 gap-y-3 border-t border-border/60 pt-6">
              <RatingBadge value={rating} votes={p.total_votes} size="xl" />
              {minutes != null && minutes > 0 && (
                <span className="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
                  <Clock size={15} /> {minutes} minutes
                </span>
              )}
              {jersey != null && (
                <span className="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
                  <Shirt size={15} /> #{jersey}
                </span>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* ─── Body ──────────────────────────────────────────────── */}
      <div className="container-max container-px relative z-10 grid grid-cols-1 gap-10 py-12 lg:grid-cols-[1fr_360px]">
        <div className="space-y-12">
          {p.description && (
            <section>
              <p className="whitespace-pre-line leading-relaxed text-muted-foreground">{p.description}</p>
            </section>
          )}

          {/* Stats grid */}
          {stats.length > 0 && (
            <section>
              <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Match Stats</h2>
              <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
                {stats.map((s, i) => (
                  <div key={s.id || `${s.stat_name}-${i}`} className="card-shell p-4">
                    <p className="font-display text-2xl font-bold stat-number text-foreground">{statValue(s)}</p>
                    <p className="mt-1 line-clamp-2 text-xs font-medium uppercase tracking-wide text-muted-foreground">
                      {prettifyStat(s)}
                    </p>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* Rating Widget - Mobile Only */}
          <section className="lg:hidden">
            <RatingWidget
              kind="performance"
              entityKey={id}
              initialAverage={rating}
              initialVotes={p.total_votes}
            />
          </section>

          {/* Reviews */}
          <section>
            <ReviewsSection kind="performance" entityKey={id} initialReviews={reviews} />
          </section>
        </div>

        {/* Sidebar */}
        <aside className="hidden lg:sticky lg:top-24 lg:block lg:self-start">
          <RatingWidget
            kind="performance"
            entityKey={id}
            initialAverage={rating}
            initialVotes={p.total_votes}
          />
        </aside>
      </div>
    </main>
  );
}
