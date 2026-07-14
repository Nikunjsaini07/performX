import React from 'react';
import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import { getTeams, getTeamMatches, getTeamPerformances } from '@/lib/data';
import type { Match, Performance } from '@/lib/api';
import TeamCrest from '@/components/TeamCrest';
import BackButton from '@/components/BackButton';
import MatchCard from '@/components/MatchCard';
import PerformanceCard from '@/components/PerformanceCard';

export const revalidate = 60;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await Promise.resolve(params);
  const teams = await getTeams(64);
  const team = teams.find((t) => t.slug === slug);
  return { title: team ? team.name : 'Team' };
}

export default async function TeamDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await Promise.resolve(params);

  const [teams, teamMatches, teamPerformances] = await Promise.all([
    getTeams(64),
    getTeamMatches(slug, 20, 0),
    getTeamPerformances(slug, 24, 0),
  ]);

  const team = teams.find((t) => t.slug === slug);
  if (!team) notFound();

  // Adapt the lighter team-scoped shapes to the cards' expected props.
  const matches: Match[] = (teamMatches || [])
    .filter((m) => m && m.id)
    .map((m) => ({
      id: m.id,
      slug: m.slug,
      title: m.title,
      utc_datetime: m.utc_datetime,
      venue: m.venue,
      home_score: m.home_score,
      away_score: m.away_score,
      home_team_name: m.home_team_name,
      home_team_short_name: m.home_team_name,
      home_team_logo_url: m.home_team_logo_url,
      away_team_name: m.away_team_name,
      away_team_short_name: m.away_team_name,
      away_team_logo_url: m.away_team_logo_url,
    }));

  const performances: Performance[] = (teamPerformances || [])
    .filter((p) => p && p.performance_id)
    .map((p) => ({
      id: p.performance_id,
      title: p.performance_title,
      match_id: p.match_id,
      player_id: p.player_id,
      player_team_id: team.id,
      cover_image_url: p.performance_cover_image,
      player_photo_url: p.player_photo_url,
      player_slug: p.player_slug,
      player_name: p.player_name,
      match_title: p.match_title,
      match_slug: p.match_slug,
      team_name: team.name,
      team_logo_url: team.logo_url,
      flag_emoji: team.flag_emoji,
      average_rating: p.average_rating,
      minutes_played: p.minutes_played,
    }));

  return (
    <main>
      <section className="mesh-bg grain relative">
        <div className="container-max container-px relative pb-14 pt-24">
          <BackButton className="mb-6" />
          <div className="flex items-center gap-5">
            <div className="card-shell flex h-24 w-24 items-center justify-center p-3">
              <TeamCrest name={team.name} shortName={team.short_name} logoUrl={team.logo_url} flagEmoji={team.flag_emoji} size={72} />
            </div>
            <div>
              <h1 className="font-display text-[clamp(2rem,4vw,3rem)] font-bold text-foreground">{team.name}</h1>
              <p className="mt-1 text-muted-foreground">
                {team.short_name}
                {team.type ? ` · ${team.type}` : ''} · FIFA World Cup 2026
              </p>
            </div>
          </div>
        </div>
      </section>

      <div className="container-max container-px space-y-12 py-12">
        {matches.length > 0 && (
          <section>
            <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Matches</h2>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {matches.map((m) => (
                <MatchCard key={m.id} match={m} />
              ))}
            </div>
          </section>
        )}

        {performances.length > 0 && (
          <section>
            <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Performances</h2>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
              {performances.map((p) => (
                <PerformanceCard key={p.id} performance={p} />
              ))}
            </div>
          </section>
        )}

        {matches.length === 0 && performances.length === 0 && (
          <div className="card-shell flex flex-col items-center gap-2 px-6 py-16 text-center">
            <p className="font-display text-lg font-bold text-foreground">Nothing here yet</p>
            <p className="max-w-sm text-sm text-muted-foreground">
              Matches and performances for {team.name} will appear here as the tournament unfolds.
            </p>
          </div>
        )}
      </div>
    </main>
  );
}
