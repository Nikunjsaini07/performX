import React from 'react';
import Link from 'next/link';
import { ArrowRight, MapPin, Trophy, Clock } from 'lucide-react';
import { getMatches, Match } from '@/lib/api';

const stageBadgeClasses: Record<string, string> = {
  gold: 'bg-primary/20 text-primary border border-primary/40',
  silver: 'bg-muted text-muted-foreground border border-border',
  default: 'bg-muted/50 text-muted-foreground border border-border/50',
};

// Helper to determine badge color based on round
function getStageColor(round: string) {
  const r = (round || '').toLowerCase();
  if (r.includes('final') && !r.includes('semi') && !r.includes('quarter')) return 'gold';
  if (r.includes('semi') || r.includes('quarter')) return 'silver';
  return 'default';
}

function renderFlag(logoUrl?: string) {
  if (!logoUrl) return '🏳️';
  if (logoUrl.startsWith('http')) {
    return <img src={logoUrl} alt="flag" className="w-12 h-8 object-cover rounded-sm shadow-sm" />;
  }
  // It's an emoji
  return <span className="text-5xl leading-none">{logoUrl}</span>;
}

function penaltyValue(value?: Match['home_penalty_score']) {
  if (typeof value === 'number') return value;
  if (value && typeof value === 'object' && value.Valid) return value.Int32;
  return null;
}

export default async function FeaturedMatches() {
  // We'll fetch from /matches for now since trending requires active cron/voting
  let matches: Match[] = [];
  try {
    matches = await getMatches(4);
  } catch (error) {
    console.error("Failed to fetch featured matches:", error);
  }

  // Fallback styling array
  const bgAccents = [
    'from-[#c9a84c]/20',
    'from-[#4ade80]/10',
    'from-[#60a5fa]/10',
    'from-[#c9a84c]/15',
  ];

  return (
    <section className="py-24 bg-background relative">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[60%] h-px bg-gradient-to-r from-transparent via-border to-transparent" />

      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        <div className="flex items-end justify-between mb-12">
          <div>
            <p className="text-[#c9a84c] text-sm font-bold tracking-widest uppercase mb-2">Recent</p>
            <h2 className="font-display text-3xl md:text-5xl font-bold text-white tracking-tight" style={{ fontFamily: "Georgia, serif" }}>
              Recent Matches
            </h2>
            <p className="text-blue-400 text-sm md:text-base mt-3 font-medium">
              The latest fixtures from FIFA 2026
            </p>
          </div>
          <Link
            href="/matches"
            className="hidden sm:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-primary transition-colors group"
          >
            All matches
            <ArrowRight size={14} className="group-hover:translate-x-0.5 transition-transform" />
          </Link>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5">
          {matches.length === 0 ? (
             <div className="col-span-full py-20 text-center text-muted-foreground">
               No matches found. Ensure the API is running.
             </div>
          ) : (
            matches.map((match, idx) => {
              const dateObj = new Date(match.utc_datetime);
              const dateStr = dateObj.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
              const bgAccent = bgAccents[idx % bgAccents.length];
              const stageColor = getStageColor(match.round || '');
              const homePen = penaltyValue(match.home_penalty_score);
              const awayPen = penaltyValue(match.away_penalty_score);

              return (
                <Link
                  href={`/matches/${match.slug}`}
                  key={match.id}
                  className="group block"
                >
                  <div className="relative overflow-hidden rounded-2xl border border-[#c9a84c]/20 bg-gradient-to-b from-[#1a1a1a] to-[#111] transition-all duration-300 hover:border-white hover:-translate-y-1 hover:shadow-[0_0_30px_rgba(255,255,255,0.15)] flex flex-col h-full font-sans cursor-pointer">
                    
                    {/* Top Header */}
                    <div className="flex items-center justify-between px-4 py-3 border-b border-[#c9a84c]/10 bg-black/40">
                      <span className="text-[#c9a84c] text-[11px] font-black tracking-[0.2em] uppercase">FIFA World Cup 2026</span>
                      <span className="text-gray-400 text-xs font-medium bg-[#222] px-2 py-0.5 rounded-full">{match.round || 'Group stage'}</span>
                    </div>

                    {/* Flags & Scores Grid */}
                    <div className="grid grid-cols-2 grid-rows-2 h-[180px]">
                      {/* Top Left: Home Flag */}
                      <div className="w-full h-full relative border-r border-b border-[#222] bg-[#000]">
                        {match.home_team_logo_url ? (
                          <div className="absolute inset-0 bg-cover bg-center opacity-90 group-hover:opacity-100 transition-opacity" style={{ backgroundImage: `url(${match.home_team_logo_url})` }} />
                        ) : (
                          <div className="w-full h-full bg-[#222] flex items-center justify-center text-sm text-gray-500">No Flag</div>
                        )}
                        <div className="absolute inset-0 bg-gradient-to-r from-black/60 via-transparent to-transparent opacity-50" />
                      </div>
                      
                      {/* Top Right: Home Score */}
                      <div className="w-full h-full bg-[#111] flex flex-col items-center justify-center border-b border-[#222] relative">
                        <span className="text-[#c9a84c] text-6xl font-black font-display tracking-tighter drop-shadow-lg">{match.home_score ?? 0}</span>
                        {homePen !== null && (
                          <span className="absolute bottom-2 right-3 text-[#c9a84c]/80 text-sm font-bold bg-[#c9a84c]/10 px-1.5 py-0.5 rounded">({homePen})</span>
                        )}
                      </div>
                      
                      {/* Bottom Left: Away Score */}
                      <div className="w-full h-full bg-[#111] flex flex-col items-center justify-center border-r border-[#222] relative">
                        <span className="text-[#c9a84c] text-6xl font-black font-display tracking-tighter drop-shadow-lg">{match.away_score ?? 0}</span>
                        {awayPen !== null && (
                          <span className="absolute bottom-2 left-3 text-[#c9a84c]/80 text-sm font-bold bg-[#c9a84c]/10 px-1.5 py-0.5 rounded">({awayPen})</span>
                        )}
                      </div>
                      
                      {/* Bottom Right: Away Flag */}
                      <div className="w-full h-full relative border-[#222] bg-[#000]">
                        {match.away_team_logo_url ? (
                          <div className="absolute inset-0 bg-cover bg-center opacity-90 group-hover:opacity-100 transition-opacity" style={{ backgroundImage: `url(${match.away_team_logo_url})` }} />
                        ) : (
                          <div className="w-full h-full bg-[#222] flex items-center justify-center text-sm text-gray-500">No Flag</div>
                        )}
                        <div className="absolute inset-0 bg-gradient-to-l from-black/60 via-transparent to-transparent opacity-50" />
                      </div>
                    </div>

                    {/* Middle Section */}
                    <div className="p-5 flex flex-col items-center flex-grow bg-gradient-to-b from-[#111] to-[#151515] relative">
                      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-24 h-px bg-gradient-to-r from-transparent via-[#c9a84c]/30 to-transparent" />
                      <p className="text-gray-300 italic text-sm text-center mb-5 line-clamp-2 leading-relaxed" style={{ fontFamily: "Georgia, serif" }}>
                        "{match.tagline || 'A thrilling matchup awaits as these two nations clash.'}"
                      </p>
                      <h3 className="text-white text-lg font-black tracking-widest uppercase text-center flex flex-wrap justify-center items-center gap-2">
                        <span>{match.home_team_name}</span> 
                        <span className="text-[#c9a84c] text-sm opacity-80">VS</span> 
                        <span>{match.away_team_name}</span>
                      </h3>
                      <p className="text-[#c9a84c]/60 text-xs mt-2 font-semibold uppercase tracking-wider">{dateStr}</p>
                    </div>

                    {/* Footer Section */}
                    <div className="p-4 border-t border-[#c9a84c]/10 bg-[#0a0a0a] flex items-center justify-between group-hover:bg-[#111] transition-colors">
                      <div className="flex items-center gap-2">
                        <span className="text-[#c9a84c] text-lg drop-shadow-[0_0_8px_rgba(201,168,76,0.5)]">★</span>
                        <div className="flex items-baseline gap-1">
                          <span className="text-white font-bold text-lg">0</span>
                          <span className="text-gray-500 text-xs font-semibold">/10</span>
                        </div>
                        <div className="w-20 h-1 bg-[#222] rounded-full ml-2 overflow-hidden hidden sm:block border border-[#333]">
                          <div className="h-full bg-gradient-to-r from-[#c9a84c] to-[#f5d57b] rounded-full" style={{ width: '0%' }}></div>
                        </div>
                      </div>
                      
                      <div className="px-4 py-1.5 border border-[#c9a84c]/50 text-[#c9a84c] rounded text-xs font-bold uppercase tracking-wider hover:bg-[#c9a84c] hover:text-black hover:border-[#c9a84c] transition-all flex items-center gap-1.5 shadow-[0_0_15px_rgba(201,168,76,0.1)] hover:shadow-[0_0_20px_rgba(201,168,76,0.4)]">
                        <span className="text-[10px] leading-none">☆</span> Rate
                      </div>
                    </div>
                  </div>
                </Link>
              );
            })
          )}
        </div>

        <div className="flex sm:hidden justify-center mt-6">
          <Link href="/matches" className="btn-ghost text-sm flex items-center gap-1.5">
            View all matches
            <ArrowRight size={14} />
          </Link>
        </div>
      </div>
    </section>
  );
}
