import React from 'react';
import Link from 'next/link';
import Navbar from '@/components/navbar';
import ArchiveFooter from '@/app/components/ArchiveFooter';
import { getTeams, Team } from '@/lib/api';
import { Search, Shield } from 'lucide-react';

export default async function TeamsPage() {
  let teams: Team[] = [];
  try {
    teams = await getTeams(60);
  } catch (error) {
    console.error('Failed to fetch teams:', error);
  }

  return (
    <main className="min-h-screen bg-background">
      <Navbar />
      <section className="pt-28 pb-20">
        <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
          <div className="mb-10 flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
            <div>
              <p className="section-label mb-2">48 Nations</p>
              <h1 className="font-display text-4xl md:text-5xl font-bold text-foreground">Teams</h1>
              <p className="mt-3 max-w-2xl text-sm text-muted-foreground">
                Browse every World Cup side, then jump into the performances and matches that define their tournament.
              </p>
            </div>
            <div className="relative w-full md:w-80">
              <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
              <input
                className="w-full rounded-lg border border-border bg-muted py-2.5 pl-9 pr-4 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
                placeholder="Search is coming with live filters"
                disabled
              />
            </div>
          </div>

          {teams.length === 0 ? (
            <div className="archive-card p-12 text-center text-muted-foreground">No teams found. Ensure the API is running.</div>
          ) : (
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
              {teams.map((team) => (
                <Link
                  key={team.id}
                  href={`/performances?team=${team.slug}`}
                  className="archive-card group p-5 text-center hover:border-primary/35 focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center overflow-hidden rounded-lg border border-border bg-muted">
                    {team.logo_url ? <img src={team.logo_url} alt="" className="h-full w-full object-cover" /> : <span className="text-4xl">{team.flag_emoji || team.short_name}</span>}
                  </div>
                  <h2 className="text-sm font-bold text-foreground leading-tight">{team.name}</h2>
                  <div className="mt-3 flex items-center justify-center gap-2 text-xs text-muted-foreground">
                    <Shield size={12} className="text-primary" />
                    <span>{team.short_name}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </section>
      <ArchiveFooter />
    </main>
  );
}
