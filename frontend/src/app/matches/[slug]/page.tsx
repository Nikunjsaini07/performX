import React from 'react';
import Link from 'next/link';
import Navbar from '@/components/navbar';
import ArchiveFooter from '@/app/components/ArchiveFooter';
import RatingPanel from '@/app/components/RatingPanel';
import ReviewsSection from '@/app/components/ReviewsSection';
import { getMatch, getMatchPerformances, getMatchReviews, Match, Performance, Review } from '@/lib/api';
import { CalendarDays, MapPin, Star, Trophy } from 'lucide-react';

interface MatchPageProps {
  params: Promise<{ slug: string }>;
}

function penaltyValue(value?: Match['home_penalty_score']) {
  if (typeof value === 'number') return value;
  if (value && typeof value === 'object' && value.Valid) return value.Int32;
  return null;
}

function teamLogo(url?: string, name?: string) {
  if (!url) return <span className="text-lg font-black">{name?.slice(0, 3).toUpperCase()}</span>;
  return <img src={url} alt="" className="h-full w-full object-cover" />;
}

function lineupRows(performances: Performance[], teamName: string) {
  return performances
    .filter((performance) => performance.team_name === teamName)
    .sort((a, b) => Number(b.is_starter) - Number(a.is_starter) || (a.jersey_number || 99) - (b.jersey_number || 99));
}

export default async function MatchPage({ params }: MatchPageProps) {
  const { slug } = await params;

  let match: Match | null = null;
  let performances: Performance[] = [];
  let reviews: Review[] = [];
  try {
    const fetchedMatch = await getMatch(slug);
    const [fetchedPerformances, fetchedReviews] = await Promise.all([
      getMatchPerformances(slug).catch(() => []),
      getMatchReviews(slug).catch(() => []),
    ]);
    match = fetchedMatch;
    performances = fetchedPerformances;
    reviews = fetchedReviews;
  } catch (error) {
    console.error(`Failed to fetch match ${slug} details:`, error);
  }

  if (!match) {
    return (
      <main className="min-h-screen bg-background">
        <Navbar />
        <div className="flex min-h-[60vh] items-center justify-center px-6">
          <div className="text-center">
            <h1 className="mb-4 text-2xl font-bold">Match Not Found</h1>
            <Link href="/matches" className="text-primary hover:underline">Return to Matches</Link>
          </div>
        </div>
        <ArchiveFooter />
      </main>
    );
  }

  const homePen = penaltyValue(match.home_penalty_score);
  const awayPen = penaltyValue(match.away_penalty_score);
  const date = new Date(match.utc_datetime);
  const homeRows = lineupRows(performances, match.home_team_name);
  const awayRows = lineupRows(performances, match.away_team_name);

  return (
    <main className="min-h-screen bg-background">
      <Navbar />

      <section className="relative min-h-[620px] overflow-hidden pt-28">
        <div className="absolute inset-0">
          <img src="/stadium.jpg" alt="" className="h-full w-full object-cover" />
          <div className="absolute inset-0 bg-gradient-to-b from-[#060708]/70 via-[#060708]/58 to-background" />
          <div className="absolute inset-x-0 bottom-0 h-56 bg-gradient-to-t from-background to-transparent" />
        </div>

        <div className="relative z-10 mx-auto flex min-h-[520px] max-w-screen-2xl flex-col justify-end px-6 pb-14 lg:px-10">
          <div className="mb-8 flex flex-wrap items-center gap-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            <span className="text-primary">{match.round || 'World Cup Match'}</span>
            {match.venue && <span className="inline-flex items-center gap-1.5"><MapPin size={13} />{match.venue}</span>}
            <span className="inline-flex items-center gap-1.5"><CalendarDays size={13} />{date.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}</span>
          </div>

          <div className="grid gap-8 lg:grid-cols-[1fr_auto_1fr] lg:items-end">
            <div className="flex items-center gap-4 lg:justify-start">
              <div className="h-20 w-20 overflow-hidden rounded-lg border border-white/15 bg-black/35 md:h-28 md:w-28">
                {teamLogo(match.home_team_logo_url, match.home_team_name)}
              </div>
              <h1 className="font-display text-4xl font-bold text-foreground md:text-6xl">{match.home_team_name}</h1>
            </div>

            <div className="text-left lg:text-center">
              <div className="stat-number text-6xl font-black tracking-tight text-foreground md:text-8xl">
                {match.home_score ?? 0}<span className="mx-3 text-primary">-</span>{match.away_score ?? 0}
              </div>
              {homePen !== null && awayPen !== null && (
                <p className="mt-2 text-sm font-bold uppercase tracking-widest text-primary">({homePen}) penalties ({awayPen})</p>
              )}
              <div className="mt-5 inline-flex items-center gap-3 rounded-lg border border-primary/30 bg-primary/10 px-4 py-3">
                <Star size={18} className="fill-primary text-primary" />
                <span className="stat-number text-3xl font-black text-primary">{Number(match.average_rating || 0).toFixed(1)}</span>
                <span className="text-xs text-muted-foreground">{match.total_votes || 0} votes</span>
              </div>
            </div>

            <div className="flex items-center gap-4 lg:justify-end lg:text-right">
              <h1 className="font-display text-4xl font-bold text-foreground md:text-6xl">{match.away_team_name}</h1>
              <div className="h-20 w-20 overflow-hidden rounded-lg border border-white/15 bg-black/35 md:h-28 md:w-28">
                {teamLogo(match.away_team_logo_url, match.away_team_name)}
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-screen-2xl space-y-12 px-6 py-12 lg:px-10">
        <div className="grid gap-8 lg:grid-cols-[1.2fr_0.8fr]">
          <div>
            <p className="section-label mb-2">Match Story</p>
            <h2 className="font-display text-3xl font-bold text-foreground">{match.tagline || match.title}</h2>
            <p className="mt-4 max-w-3xl text-base leading-relaxed text-muted-foreground">
              {match.description || `${match.home_team_name} and ${match.away_team_name} meet in a World Cup fixture built for debate: scoreline, momentum, and the player performances that shaped the night.`}
            </p>
          </div>
          <div className="grid grid-cols-2 gap-3">
            {[
              ['Venue', match.venue || 'TBD'],
              ['Stage', match.round || 'Group Stage'],
              ['Kickoff', date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })],
              ['Performances', performances.length.toString()],
            ].map(([label, value]) => (
              <div key={label} className="archive-card p-4">
                <p className="text-xs uppercase tracking-widest text-muted-foreground">{label}</p>
                <p className="mt-2 text-sm font-bold text-foreground">{value}</p>
              </div>
            ))}
          </div>
        </div>

        <RatingPanel label="Your Match Rating" />

        <section>
          <div className="mb-5 flex items-center gap-2">
            <Trophy size={18} className="text-primary" />
            <h2 className="font-display text-3xl font-bold text-foreground">Lineups & Performances</h2>
          </div>
          <div className="grid gap-5 lg:grid-cols-2">
            {[
              [match.home_team_name, homeRows],
              [match.away_team_name, awayRows],
            ].map(([team, rows]) => (
              <div key={team as string} className="archive-card overflow-hidden">
                <div className="border-b border-border px-5 py-4">
                  <h3 className="font-bold text-foreground">{team as string}</h3>
                </div>
                {(rows as Performance[]).length === 0 ? (
                  <div className="p-6 text-sm text-muted-foreground">No player performances logged yet.</div>
                ) : (
                  <div className="divide-y divide-border">
                    {(rows as Performance[]).map((performance) => (
                      <Link key={performance.id} href={`/performances/${performance.id}`} className="grid grid-cols-[44px_1fr_auto_auto] items-center gap-3 px-5 py-3 hover:bg-primary/5 transition-colors">
                        <span className="stat-number text-xs text-muted-foreground">#{performance.jersey_number || '-'}</span>
                        <span className="min-w-0">
                          <span className="block truncate text-sm font-semibold text-foreground">{performance.player_name}</span>
                          <span className="text-xs text-muted-foreground">{performance.minutes_played || 0} min</span>
                        </span>
                        {performance.is_starter && <span className="rounded bg-muted px-2 py-1 text-[10px] font-bold uppercase text-muted-foreground">XI</span>}
                        <span className="performance-score w-10 h-9">{Number(performance.average_rating || performance.average_rating_2 || 0).toFixed(1)}</span>
                      </Link>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </section>

        <ReviewsSection title="Match Reviews" reviews={reviews} />
      </section>

      <ArchiveFooter />
    </main>
  );
}
