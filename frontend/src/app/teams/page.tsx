import React from 'react';
import Link from 'next/link';
import type { Metadata } from 'next';
import { Flag } from 'lucide-react';
import { getTeams } from '@/lib/data';
import TeamCrest from '@/components/TeamCrest';

export const revalidate = 60;

export const metadata: Metadata = {
  title: 'Teams',
  description: 'All 48 nations competing at the FIFA World Cup 2026.',
};

export default async function TeamsPage() {
  const teams = await getTeams(64);

  return (
    <main className="container-max container-px pb-12 pt-24">
      <header className="mb-8">
        <span className="section-label mb-2 block">48 nations, one trophy</span>
        <h1 className="font-display text-[clamp(2rem,4vw,2.75rem)] font-bold text-foreground">Teams</h1>
        <p className="mt-2 max-w-2xl text-muted-foreground">
          Every nation at the World Cup. Follow your side through the group stage and into the knockouts.
        </p>
      </header>

      {teams.length === 0 ? (
        <div className="card-shell flex flex-col items-center gap-3 px-6 py-20 text-center">
          <Flag size={30} className="text-muted-foreground/50" />
          <p className="font-display text-lg font-bold text-foreground">No teams yet</p>
          <p className="max-w-sm text-sm text-muted-foreground">The team list will appear here once published.</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
          {teams.map((t) => (
            <Link
              key={t.id}
              href={`/teams/${t.slug}`}
              className="card-shell card-lift group flex flex-col items-center gap-3 p-5 text-center"
            >
              <TeamCrest name={t.name} shortName={t.short_name} logoUrl={t.logo_url} flagEmoji={t.flag_emoji} size={56} />
              <div>
                <p className="line-clamp-1 text-sm font-semibold text-foreground transition-colors group-hover:text-primary">
                  {t.name}
                </p>
                {t.short_name && t.short_name !== t.name && (
                  <p className="text-xs text-muted-foreground">{t.short_name}</p>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}
