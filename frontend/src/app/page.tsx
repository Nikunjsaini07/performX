import React from 'react';
import Link from 'next/link';
import { ArrowRight, TrendingUp, Star, Trophy } from 'lucide-react';
import {
  getMatches,
  getTrendingMatches,
  getTrendingPerformances,
  getTrendingPlayersRanked,
  getTrendingReviews,
  getTeams,
  getOverviewStats,
} from '@/lib/data';
import { SITE, BACKGROUNDS } from '@/lib/site';
import MatchCard from '@/components/features/MatchCard';
import PerformanceCard from '@/components/features/PerformanceCard';
import ReviewCard from '@/components/features/ReviewCard';
import SectionRail from '@/components/layout/SectionRail';
import TeamCrest from '@/components/TeamCrest';
import { initials } from '@/lib/format';

export const revalidate = 0;

export default async function HomePage() {
  const [recentMatches, topPerformances, topMatches, trendingPlayers, trendingReviews, teams, overview] =
    await Promise.all([
      getMatches(12, 0),
      getTrendingPerformances(12),
      getTrendingMatches(12),
      getTrendingPlayersRanked(8),
      getTrendingReviews(6),
      getTeams(24), // Reduced from 48 to 24 for faster loading
      getOverviewStats(),
    ]);

  console.log('Trending Reviews:', JSON.stringify(trendingReviews, null, 2));

  const maxScore = trendingPlayers.reduce((m, p) => Math.max(m, p.score || 0), 0) || 1;

  // Real totals from /stats/overview (no hardcoding, no "+").
  const stats = [
    { label: 'Matches', value: overview.match_count },
    { label: 'Performances', value: overview.performance_count },
    { label: 'Teams', value: overview.team_count },
    { label: 'Ratings', value: overview.rating_count },
  ].filter((s) => (s.value ?? 0) > 0);

  return (
    <main>
      {/* ─── Hero ─────────────────────────────────────────────── */}
      <section className="relative overflow-hidden">
        {/* Background Image with Fade */}
        <div className="detail-backdrop absolute inset-0">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            src={BACKGROUNDS.home}
            alt=""
            className="h-full w-full object-cover object-center"
          />
          {/* Additional lighter overlay on lower section for better text readability */}
          <div className="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-black/30" />
        </div>
        
        <div className="container-max container-px relative py-16 sm:py-20 md:py-24 pt-52 sm:pt-64 md:pt-72">
          <div className="max-w-2xl">
            <span className="section-label mb-4 text-white/100">
              <TrendingUp size={14} /> FIFA World Cup 2026 · USA · Canada · Mexico
            </span>
            <h1 className="font-display text-[clamp(2rem,4.5vw,3.25rem)] font-bold leading-[1.05] text-white">
              Rate every match.
              <br />
              Review every <span className="text-gradient-lime">performance</span>.
            </h1>
            <p className="mt-4 max-w-lg text-sm leading-relaxed text-white/75 sm:text-base">
              {SITE.name} is the community archive for the World Cup — score matches and players
              out of 10, write reviews, and see exactly what the world is watching.
            </p>
            <div className="mt-6 flex flex-wrap items-center gap-3">
              <Link href="/matches" className="btn-primary">
                Browse matches <ArrowRight size={16} />
              </Link>
              <Link href="/performances" className="btn-ghost-light">
                Explore performances
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* ─── Rails, sitting on an atmospheric stadium backdrop ──── */}
      <div className="relative overflow-visible">
        <div className="ambience-bg grain" aria-hidden="true">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img src={BACKGROUNDS.ambience} alt="" />
        </div>
        <div className="relative overflow-visible">

      {/* ─── Recent Matches ────────────────────────────────────── */}
      <SectionRail
        label="Fresh off the pitch"
        title="Recent Matches"
        viewAllHref="/matches"
        show={recentMatches.length > 0}
      >
        {recentMatches.map((m) => (
          <MatchCard key={m.id} match={m} className="w-[300px] shrink-0 snap-start" />
        ))}
      </SectionRail>

      {/* ─── Top Performances ──────────────────────────────────── */}
      <SectionRail
        label="Player of the week"
        title="Top Performances"
        viewAllHref="/performances"
        show={topPerformances.length > 0}
      >
        {topPerformances.filter(Boolean).map((p) => (
          <PerformanceCard key={p.id} performance={p} compact className="w-[160px] shrink-0 snap-start" />
        ))}
      </SectionRail>

      {/* ─── Top Matches ───────────────────────────────────────── */}
      <SectionRail
        label="Most rated"
        title="Top Matches"
        viewAllHref="/matches"
        show={topMatches.length > 0}
      >
        {topMatches.map((m) => (
          <MatchCard key={m.id} match={m} className="w-[300px] shrink-0 snap-start" />
        ))}
      </SectionRail>

      {/* ─── Trending Players (ranked) ─────────────────────────── */}
      {trendingPlayers.length > 0 && (
        <section className="container-max container-px py-8">
          <div className="mb-5">
            <span className="section-label mb-2 block">
              <Star size={14} /> This week
            </span>
            <h2 className="font-display text-[1.75rem] font-bold text-foreground">Trending Players</h2>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            {trendingPlayers.map((tp, i) => (
              <Link 
                key={tp.player.id} 
                href={`/players/${tp.player.slug || tp.player.id}`}
                className="card-shell flex items-center gap-4 p-3.5 transition-all duration-300 hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
              >
                <span className="w-6 shrink-0 text-center font-display text-lg font-bold stat-number text-muted-foreground">
                  {i + 1}
                </span>
                <span className="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden rounded-full bg-surface-2 text-sm font-bold text-muted-foreground">
                  {tp.player.photo_url ? (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img src={tp.player.photo_url} alt={tp.player.name} className="h-full w-full object-cover object-top" />
                  ) : (
                    initials(tp.player.name)
                  )}
                </span>
                <div className="min-w-0 flex-1">
                  <p className="line-clamp-1 text-sm font-semibold text-foreground">{tp.player.name}</p>
                  <div className="mt-0.5 flex items-center gap-2 text-xs text-muted-foreground">
                    {tp.player.team_name && <span className="line-clamp-1">{tp.player.team_name}</span>}
                    {tp.player.position && tp.player.team_name && <span>·</span>}
                    {tp.player.position && <span>{tp.player.position}</span>}
                    {tp.player.jersey_number && (
                      <>
                        <span>·</span>
                        <span>#{tp.player.jersey_number}</span>
                      </>
                    )}
                  </div>
                  <div className="mt-2 flex items-center gap-2">
                    <div className="h-1.5 flex-1 overflow-hidden rounded-full bg-surface-2">
                      <div
                        className="h-full rounded-full bg-primary"
                        style={{ width: `${Math.max(8, ((tp.score || 0) / maxScore) * 100)}%` }}
                      />
                    </div>
                    <span className="text-xs font-semibold text-primary" title="Trending score based on ratings, votes, and recency">
                      {tp.score.toFixed(0)} pts
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* ─── Trending Reviews ──────────────────────────────────── */}
      {trendingReviews.length > 0 && (
        <section className="container-max container-px py-8">
          <div className="mb-5">
            <span className="section-label mb-2 block">From the community</span>
            <h2 className="font-display text-[1.75rem] font-bold text-foreground">Trending Reviews</h2>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {trendingReviews.slice(0, 6).map((r) => (
              <ReviewCard key={r.id} review={r} />
            ))}
          </div>
        </section>
      )}

      {/* ─── Teams strip ───────────────────────────────────────── */}
      {teams.length > 0 && (
        <section className="container-max container-px py-8 pb-16">
          <div className="mb-5 flex items-end justify-between gap-4">
            <div>
              <span className="section-label mb-2 block">
                <Trophy size={14} /> {(overview.team_count || teams.length).toLocaleString()} nations
              </span>
              <h2 className="font-display text-[1.75rem] font-bold text-foreground">Teams</h2>
            </div>
            <Link
              href="/teams"
              className="group inline-flex shrink-0 items-center gap-1.5 text-sm font-semibold text-muted-foreground transition-colors hover:text-primary"
            >
              View all <ArrowRight size={15} className="transition-transform group-hover:translate-x-0.5" />
            </Link>
          </div>
          <div className="hide-scrollbar snap-x-rail -mx-6 flex gap-3 overflow-x-auto overflow-y-visible px-6 pb-2 pt-1 lg:-mx-10 lg:px-10">
            {teams.map((t) => (
              <Link
                key={t.id}
                href={`/teams/${t.slug}`}
                className="card-shell card-lift flex w-[130px] shrink-0 snap-start flex-col items-center gap-2.5 p-4 text-center"
              >
                <TeamCrest name={t.name} shortName={t.short_name} logoUrl={t.logo_url} flagEmoji={t.flag_emoji} size={44} />
                <span className="line-clamp-1 w-full text-xs font-semibold text-foreground">
                  {t.short_name || t.name}
                </span>
              </Link>
            ))}
          </div>
        </section>
      )}
        </div>
      </div>
    </main>
  );
}
