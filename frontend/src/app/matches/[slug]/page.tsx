import React from 'react';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import { CalendarDays, MapPin } from 'lucide-react';
import { getMatch, getMatchPerformances, getMatchReviews } from '@/lib/data';
import { formatDateTime, penaltyValue, toNumber } from '@/lib/format';
import type { Performance } from '@/lib/api';
import DetailBackdrop from '@/components/layout/DetailBackdrop';
import BackButton from '@/components/layout/BackButton';
import TeamCrest from '@/components/TeamCrest';
import RatingBadge from '@/components/features/RatingBadge';
import RatingWidget from '@/components/features/RatingWidget';
import ReviewsSection from '@/components/features/ReviewsSection';
import PlayerTile from '@/components/features/PlayerTile';

export const revalidate = 60;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await Promise.resolve(params);
  const match = await getMatch(slug);
  return {
    title: match ? match.title : 'Match',
    description: match?.tagline || match?.description,
  };
}

function PerformanceRow({ p }: { p: Performance }) {
  const rating = toNumber(p.average_rating) ?? toNumber(p.average_rating_2) ?? toNumber(p.provider_rating);
  const goals = toNumber(p.goals);
  const assists = toNumber(p.assists);
  const minutes = toNumber(p.minutes_played);
  const jersey = toNumber(p.jersey_number);

  return (
    <Link
      href={`/performances/${p.slug || p.id}`}
      className="group flex items-center gap-3 rounded-[14px] border border-border bg-surface px-3 py-2.5 transition-all hover:-translate-y-0.5 hover:border-primary/35"
    >
      <div className="h-12 w-12 shrink-0 overflow-hidden rounded-lg">
        <PlayerTile photoUrl={p.player_photo_url} name={p.player_name} rounded="rounded-lg" className="h-full w-full" />
      </div>
      <div className="min-w-0 flex-1">
        <p className="line-clamp-1 text-sm font-semibold text-foreground transition-colors group-hover:text-primary">
          {jersey != null && <span className="mr-1.5 text-muted-foreground stat-number">#{jersey}</span>}
          {p.player_name}
        </p>
        <div className="mt-0.5 flex flex-wrap items-center gap-x-2.5 text-[11px] text-muted-foreground">
          {minutes != null && minutes > 0 && <span>{minutes}&apos;</span>}
          {goals != null && goals > 0 && <span className="text-foreground">{goals} G</span>}
          {assists != null && assists > 0 && <span className="text-foreground">{assists} A</span>}
          {p.position && <span>{p.position}</span>}
        </div>
      </div>
      {rating != null ? (
        <RatingBadge value={rating} size="md" />
      ) : (
        <span className="text-[11px] text-muted-foreground/60">—</span>
      )}
    </Link>
  );
}

export default async function MatchDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await Promise.resolve(params);

  const match = await getMatch(slug);
  if (!match) notFound();

  const [performances, reviews] = await Promise.all([
    getMatchPerformances(slug),
    getMatchReviews(slug),
  ]);

  const homePens = penaltyValue(match.home_penalty_score);
  const awayPens = penaltyValue(match.away_penalty_score);
  const hasPens = homePens != null && awayPens != null;

  const homePerf = performances.filter((p) => p.team_name === match.home_team_name);
  const awayPerf = performances.filter((p) => p.team_name === match.away_team_name);
  const otherPerf = performances.filter(
    (p) => p.team_name !== match.home_team_name && p.team_name !== match.away_team_name,
  );

  return (
    <main>
      {/* ─── Cinematic header ──────────────────────────────────── */}
      <section className="relative">
        <DetailBackdrop />
        <div className="container-max container-px relative pt-32">
          <BackButton className="mb-6" />
          <div className="flex flex-wrap items-center gap-2.5">
            <span className="text-xs font-semibold uppercase tracking-[0.2em] text-foreground/90 [text-shadow:0_2px_8px_rgba(0,0,0,0.6)]">
              {match.round || 'Match'}
            </span>
            {match.group && (
              <span className="text-xs font-semibold uppercase tracking-[0.2em] text-foreground/90 [text-shadow:0_2px_8px_rgba(0,0,0,0.6)]">
                {match.group}
              </span>
            )}
          </div>

          {/* Scoreline */}
          <div className="mt-8 grid grid-cols-3 items-center gap-4 sm:gap-8">
            <div className="flex flex-col items-center gap-3 text-center sm:flex-row sm:justify-end sm:text-right">
              <span className="order-2 font-display text-lg font-bold text-foreground sm:order-1 sm:text-2xl">
                {match.home_team_name}
              </span>
              <TeamCrest
                name={match.home_team_name}
                shortName={match.home_team_short_name}
                logoUrl={match.home_team_logo_url}
                size={64}
                className="order-1 sm:order-2"
              />
            </div>

            <div className="flex flex-col items-center">
              <div className="flex items-baseline gap-3 font-display text-5xl font-bold stat-number text-foreground sm:text-6xl">
                <span className="inline-flex items-baseline gap-1">
                  <span>{toNumber(match.home_score) ?? 0}</span>
                  {hasPens && (
                    <span className="text-sm font-medium text-muted-foreground md:text-base">({homePens})</span>
                  )}
                </span>
                <span className="text-muted-foreground/40">–</span>
                <span className="inline-flex items-baseline gap-1">
                  <span>{toNumber(match.away_score) ?? 0}</span>
                  {hasPens && (
                    <span className="text-sm font-medium text-muted-foreground md:text-base">({awayPens})</span>
                  )}
                </span>
              </div>
            </div>

            <div className="flex flex-col items-center gap-3 text-center sm:flex-row sm:justify-start sm:text-left">
              <TeamCrest
                name={match.away_team_name}
                shortName={match.away_team_short_name}
                logoUrl={match.away_team_logo_url}
                size={64}
              />
              <span className="font-display text-lg font-bold text-foreground sm:text-2xl">
                {match.away_team_name}
              </span>
            </div>
          </div>

          {/* Meta + rating */}
          <div className="mt-8 flex flex-wrap items-center justify-center gap-x-6 gap-y-3 pt-6">
            <span className="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
              <CalendarDays size={15} /> {formatDateTime(match.utc_datetime)}
            </span>
            {match.venue && (
              <span className="inline-flex items-center gap-1.5 text-sm text-muted-foreground">
                <MapPin size={15} /> {match.venue}
              </span>
            )}
            <span className="inline-flex items-center gap-2">
              <RatingBadge value={match.average_rating} votes={match.total_votes} size="md" />
            </span>
          </div>
        </div>
      </section>

      {/* ─── Body ──────────────────────────────────────────────── */}
      <div className="container-max container-px relative grid grid-cols-1 gap-10 py-12 lg:grid-cols-[1fr_360px]">
        <div className="space-y-12">
          {(match.tagline || match.description) && (
            <section>
              {match.tagline && (
                <p className="font-display text-xl font-semibold leading-snug text-foreground">{match.tagline}</p>
              )}
              {match.description && (
                <p className="mt-3 whitespace-pre-line leading-relaxed text-muted-foreground">{match.description}</p>
              )}
            </section>
          )}

          {/* Player performances */}
          {performances.length > 0 && (
            <section>
              <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Player Performances</h2>
              <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                {homePerf.length > 0 && (
                  <div>
                    <div className="mb-3 flex items-center gap-2">
                      <TeamCrest name={match.home_team_name} logoUrl={match.home_team_logo_url} size={20} />
                      <h3 className="text-sm font-semibold text-foreground">{match.home_team_name}</h3>
                    </div>
                    <div className="space-y-2">
                      {homePerf.map((p) => (
                        <PerformanceRow key={p.id} p={p} />
                      ))}
                    </div>
                  </div>
                )}
                {awayPerf.length > 0 && (
                  <div>
                    <div className="mb-3 flex items-center gap-2">
                      <TeamCrest name={match.away_team_name} logoUrl={match.away_team_logo_url} size={20} />
                      <h3 className="text-sm font-semibold text-foreground">{match.away_team_name}</h3>
                    </div>
                    <div className="space-y-2">
                      {awayPerf.map((p) => (
                        <PerformanceRow key={p.id} p={p} />
                      ))}
                    </div>
                  </div>
                )}
                {otherPerf.length > 0 && (
                  <div className="md:col-span-2">
                    <div className="space-y-2">
                      {otherPerf.map((p) => (
                        <PerformanceRow key={p.id} p={p} />
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </section>
          )}

          {/* Rating Widget - Mobile Only */}
          <section className="lg:hidden">
            <RatingWidget
              kind="match"
              entityKey={slug}
              initialAverage={match.average_rating}
              initialVotes={match.total_votes}
            />
          </section>

          {/* Reviews */}
          <section>
            <ReviewsSection kind="match" entityKey={slug} initialReviews={reviews} />
          </section>
        </div>

        {/* Sidebar */}
        <aside className="hidden lg:sticky lg:top-24 lg:block lg:self-start">
          <RatingWidget
            kind="match"
            entityKey={slug}
            initialAverage={match.average_rating}
            initialVotes={match.total_votes}
          />
        </aside>
      </div>
    </main>
  );
}
