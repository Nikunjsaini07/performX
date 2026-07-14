import React from 'react';
import Link from 'next/link';
import { ArrowRight, Shield } from 'lucide-react';
import { getTeams, Team } from '@/lib/api';

const stageColorMap: Record<string, string> = {
  CLUB: 'text-primary border-primary/40 bg-primary/10',
  NATIONAL: 'text-accent border-accent/30 bg-accent/10',
};

const stageGlowMap: Record<string, string> = {
  CLUB: 'group-hover:border-primary/40 group-hover:shadow-[0_8px_32px_rgba(201,168,76,0.15)]',
  NATIONAL: 'group-hover:border-accent/30 group-hover:shadow-[0_8px_32px_rgba(74,222,128,0.1)]',
};

export default async function TeamsGrid() {
  let teams: Team[] = [];
  try {
    teams = await getTeams(12);
  } catch (error) {
    console.error("Failed to fetch teams:", error);
  }

  return (
    <section className="py-24 bg-background relative">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[60%] h-px bg-gradient-to-r from-transparent via-border to-transparent" />

      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        <div className="flex items-end justify-between mb-12">
          <div>
            <p className="section-label mb-2">Explore</p>
            <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground">
              Top Teams
            </h2>
            <p className="text-sm text-muted-foreground mt-2">
              Every participating nation with archived performances
            </p>
          </div>
          <div className="hidden sm:flex items-center gap-2 text-xs text-muted-foreground">
            <Shield size={12} className="text-primary" />
            <span>Discover more</span>
          </div>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-6 gap-4">
          {teams.length === 0 ? (
            <div className="col-span-full py-20 text-center text-muted-foreground">
              No teams found. Ensure the API is running.
            </div>
          ) : (
            teams.map((team) => {
              // We'll use team.type if available from the backend, otherwise default to "NATIONAL"
              const type = (team as any).type || "NATIONAL";
              
              return (
                <Link
                  href={`/performances?team=${team.slug}`}
                  key={team.id}
                  className={`group block relative overflow-hidden rounded-2xl border border-border bg-card transition-all duration-300 group-hover:-translate-y-1 ${stageGlowMap[type] || ''}`}
                >
                  <div className="p-5 flex flex-col items-center text-center gap-2.5">
                    {/* Logo or Flag */}
                    <span className="text-5xl leading-none mb-1 group-hover:scale-110 transition-transform duration-200 h-14 flex items-center justify-center">
                      {team.logo_url ? (
                        <img src={team.logo_url} alt={team.name} className="max-h-full max-w-full object-contain" />
                      ) : (
                        team.flag_emoji || '⚽'
                      )}
                    </span>

                    {/* Name */}
                    <p className="text-sm font-bold text-foreground leading-tight">{team.name}</p>

                    {/* Type/Stage */}
                    <span className={`stage-badge text-xs ${stageColorMap[type] || 'text-muted-foreground border-border bg-muted'}`}>
                      {type}
                    </span>

                    {/* Short Name label */}
                    <div className="absolute top-2.5 right-2.5 w-7 h-5 rounded bg-muted border border-border flex items-center justify-center">
                      <span className="text-[9px] font-bold text-muted-foreground">{team.short_name || '...'}</span>
                    </div>
                  </div>
                </Link>
              );
            })
          )}
        </div>

        <div className="flex justify-center mt-10">
          <Link href="/teams" className="btn-ghost flex items-center gap-1.5 text-sm">
            Browse all teams
            <ArrowRight size={14} />
          </Link>
        </div>
      </div>
    </section>
  );
}
