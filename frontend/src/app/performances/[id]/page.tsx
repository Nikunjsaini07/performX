import React from 'react';
import Link from 'next/link';
import Navbar from '@/components/navbar';
import ArchiveFooter from '@/app/components/ArchiveFooter';
import RatingPanel from '@/app/components/RatingPanel';
import ReviewsSection from '@/app/components/ReviewsSection';
import { getPerformance, getPerformanceReviews, getPerformanceStats, Performance, Review, StatRow } from '@/lib/api';
import { BadgeCheck, Clock, Flag, Star } from 'lucide-react';

interface PerformancePageProps {
  params: Promise<{ id: string }>;
}

function displayStatValue(value?: number | string, suffix = '') {
  if (value === undefined || value === null || value === '') return '-';
  return `${value}${suffix}`;
}

export default async function PerformancePage({ params }: PerformancePageProps) {
  const { id } = await params;

  let performance: Performance | null = null;
  let stats: StatRow[] = [];
  let reviews: Review[] = [];
  try {
    const fetchedPerformance = await getPerformance(id);
    const [fetchedStats, fetchedReviews] = await Promise.all([
      getPerformanceStats(id).catch(() => []),
      getPerformanceReviews(id).catch(() => []),
    ]);
    performance = fetchedPerformance;
    stats = fetchedStats;
    reviews = fetchedReviews;
  } catch (error) {
    console.error(`Failed to fetch performance ${id}:`, error);
  }

  if (!performance) {
    return (
      <main className="min-h-screen bg-background">
        <Navbar />
        <div className="flex min-h-[60vh] items-center justify-center px-6">
          <div className="text-center">
            <h1 className="mb-4 text-2xl font-bold">Performance Not Found</h1>
            <Link href="/performances" className="text-primary hover:underline">Return to Performances</Link>
          </div>
        </div>
        <ArchiveFooter />
      </main>
    );
  }

  const communityRating = Number(performance.average_rating || 0);
  const photo = performance.cover_image_url || performance.player_photo_url;
  const keyStats = [
    ['Goals', displayStatValue(performance.goals)],
    ['Assists', displayStatValue(performance.assists)],
    ['Minutes', displayStatValue(performance.minutes_played || 90)],
    ['Pass Accuracy', displayStatValue(performance.passes_accuracy, '%')],
    ['Dribbles', displayStatValue(performance.dribbles)],
  ];

  return (
    <main className="min-h-screen bg-background">
      <Navbar />

      <section className="relative min-h-[650px] overflow-hidden pt-28">
        <div className="absolute inset-0">
          <img src="/stadium.jpg" alt="" className="h-full w-full object-cover" />
          <div className="absolute inset-0 bg-gradient-to-b from-[#060708]/70 via-[#060708]/62 to-background" />
          <div className="absolute inset-x-0 bottom-0 h-60 bg-gradient-to-t from-background to-transparent" />
        </div>

        <div className="relative z-10 mx-auto grid min-h-[540px] max-w-screen-2xl items-end gap-10 px-6 pb-14 md:grid-cols-[280px_1fr] lg:px-10">
          <div className="w-56 overflow-hidden rounded-lg border border-white/15 bg-card shadow-2xl md:w-72">
            <div className="aspect-[4/5] bg-muted">
              {photo ? <img src={photo} alt={performance.player_name} className="h-full w-full object-cover object-top" /> : <div className="flex h-full items-center justify-center text-5xl text-muted-foreground">PX</div>}
            </div>
          </div>

          <div>
            <div className="mb-5 flex flex-wrap items-center gap-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              <span className="text-primary">{performance.team_name || 'World Cup 2026'}</span>
              {performance.jersey_number && <span className="inline-flex items-center gap-1.5"><BadgeCheck size={13} />#{performance.jersey_number}</span>}
              <span className="inline-flex items-center gap-1.5"><Clock size={13} />{performance.minutes_played || 90} minutes</span>
              {performance.flag_emoji && <span className="inline-flex items-center gap-1.5"><Flag size={13} />{performance.flag_emoji}</span>}
            </div>

            <h1 className="font-display text-5xl font-bold leading-tight text-foreground md:text-7xl">{performance.player_name}</h1>
            <p className="mt-3 max-w-3xl text-lg text-muted-foreground">{performance.title}</p>

            <div className="mt-8 flex flex-wrap items-end gap-5">
              <div>
                <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Community Rating</p>
                <div className="mt-1 flex items-center gap-3">
                  <Star size={28} className="fill-primary text-primary" />
                  <span className="stat-number text-7xl font-black text-primary">{communityRating.toFixed(1)}</span>
                  <span className="pb-3 text-sm text-muted-foreground">/10 from {performance.total_votes || 0} votes</span>
                </div>
              </div>
              {performance.match_title && (
                <Link href="/matches" className="mb-3 rounded-lg border border-border bg-card px-4 py-3 text-sm text-muted-foreground transition-colors hover:border-primary/40 hover:text-foreground">
                  {performance.match_title} · {performance.match_round || 'Group Stage'}
                </Link>
              )}
            </div>
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-screen-2xl space-y-12 px-6 py-12 lg:px-10">
        <section>
          <p className="section-label mb-2">Statistical Breakdown</p>
          <div className="grid grid-cols-2 gap-3 md:grid-cols-5">
            {keyStats.map(([label, value]) => (
              <div key={label} className="archive-card p-5">
                <p className="text-xs uppercase tracking-widest text-muted-foreground">{label}</p>
                <p className="stat-number mt-3 text-3xl font-black text-foreground">{value}</p>
              </div>
            ))}
          </div>

          {stats.length > 0 && (
            <div className="mt-5 grid grid-cols-2 gap-3 md:grid-cols-4 lg:grid-cols-6">
              {stats.map((stat, index) => (
                <div key={stat.id || `${stat.stat_name}-${index}`} className="rounded-lg border border-border bg-muted/40 p-4">
                  <p className="text-[11px] uppercase tracking-widest text-muted-foreground">{stat.stat_name || stat.stat_short_name || stat.stat_type_id}</p>
                  <p className="stat-number mt-2 text-xl font-bold text-foreground">{stat.value}{stat.stat_unit || ''}</p>
                </div>
              ))}
            </div>
          )}
        </section>

        <RatingPanel label="Your Performance Rating" />
        <ReviewsSection title="Performance Reviews" reviews={reviews} />
      </section>

      <ArchiveFooter />
    </main>
  );
}
