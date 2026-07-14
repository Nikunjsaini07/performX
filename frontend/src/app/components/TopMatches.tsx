import React from 'react';
import Link from 'next/link';
import { ArrowRight, Star } from 'lucide-react';
import { getMatches, Match } from '@/lib/api';

function penaltyValue(value?: Match['home_penalty_score']) {
  if (typeof value === 'number') return value;
  if (value && typeof value === 'object' && value.Valid) return value.Int32;
  return null;
}

export default async function TopMatches() {
  let matches: Match[] = [];
  try {
    matches = (await getMatches(96))
      .filter((match) => Number(match.average_rating || 0) > 0)
      .sort((a, b) => Number(b.average_rating || 0) - Number(a.average_rating || 0))
      .slice(0, 4);
  } catch (error) {
    console.error('Failed to fetch top matches:', error);
  }

  return (
    <section className="py-20 bg-background">
      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        <div className="flex items-end justify-between mb-10">
          <div>
            <p className="section-label mb-2">Community Table</p>
            <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground">Top Matches</h2>
            <p className="text-sm text-muted-foreground mt-2">The fixtures fans rated highest after the final whistle.</p>
          </div>
          <Link href="/matches" className="hidden sm:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-primary transition-colors">
            Full archive
            <ArrowRight size={14} />
          </Link>
        </div>

        {matches.length === 0 ? (
          <div className="archive-card p-10 text-center text-muted-foreground">No rated matches yet. Ratings will appear here once the archive has votes.</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4">
            {matches.map((match, index) => {
              const homePen = penaltyValue(match.home_penalty_score);
              const awayPen = penaltyValue(match.away_penalty_score);
              return (
                <Link key={match.id} href={`/matches/${match.slug}`} className="archive-card group p-5 hover:border-primary/40">
                  <div className="flex items-center justify-between mb-5">
                    <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">#{index + 1}</span>
                    <span className="inline-flex items-center gap-1 rounded-md bg-primary/10 px-2 py-1 text-sm font-bold text-primary">
                      <Star size={13} className="fill-primary" />
                      {Number(match.average_rating || 0).toFixed(1)}
                    </span>
                  </div>
                  <div className="space-y-3">
                    {[{ name: match.home_team_name, logo: match.home_team_logo_url, score: match.home_score, pen: homePen }, { name: match.away_team_name, logo: match.away_team_logo_url, score: match.away_score, pen: awayPen }].map((team) => (
                      <div key={team.name} className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-md bg-muted border border-border flex items-center justify-center overflow-hidden">
                          {team.logo ? <img src={team.logo} alt="" className="w-full h-full object-cover" /> : <span className="text-xs">{team.name.slice(0, 3)}</span>}
                        </div>
                        <span className="flex-1 text-sm font-semibold text-foreground truncate">{team.name}</span>
                        <span className="stat-number text-2xl font-black text-foreground">{team.score ?? 0}</span>
                        {team.pen !== null && <span className="stat-number text-xs text-primary">({team.pen})</span>}
                      </div>
                    ))}
                  </div>
                  <div className="mt-5 pt-4 border-t border-border">
                    <p className="text-xs text-muted-foreground line-clamp-2">{match.tagline || match.title}</p>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </div>
    </section>
  );
}
