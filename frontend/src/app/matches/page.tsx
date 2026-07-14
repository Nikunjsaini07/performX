import React from 'react';
import type { Metadata } from 'next';
import { CalendarSearch } from 'lucide-react';
import { getMatches, getOverviewStats } from '@/lib/data';
import MatchCard from '@/components/MatchCard';
import Pagination from '@/components/Pagination';

export const revalidate = 60;

export const metadata: Metadata = {
  title: 'Matches',
  description: 'Browse and rate every FIFA World Cup 2026 match.',
};

const PAGE_SIZE = 24;

export default async function MatchesPage({
  searchParams,
}: {
  searchParams: Promise<{ page?: string }>;
}) {
  const sp = await Promise.resolve(searchParams);
  const page = Math.max(1, parseInt(sp?.page || '1', 10) || 1);
  const offset = (page - 1) * PAGE_SIZE;

  const [matches, overview] = await Promise.all([
    getMatches(PAGE_SIZE, offset),
    getOverviewStats(),
  ]);
  const totalPages = Math.max(1, Math.ceil((overview.match_count || 0) / PAGE_SIZE));

  return (
    <main className="container-max container-px pb-12 pt-24">
      <header className="mb-8">
        <span className="section-label mb-2 block">The archive</span>
        <h1 className="font-display text-[clamp(2rem,4vw,2.75rem)] font-bold text-foreground">Matches</h1>
        <p className="mt-2 max-w-2xl text-muted-foreground">
          Every fixture from the World Cup — tap in to rate the game and read the community&apos;s verdict.
        </p>
      </header>

      {matches.length === 0 ? (
        <EmptyState />
      ) : (
        <>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {matches.filter(Boolean).map((m) => (
              <MatchCard key={m.id} match={m} />
            ))}
          </div>

          <Pagination currentPage={page} totalPages={totalPages} basePath="/matches" />
        </>
      )}
    </main>
  );
}

function EmptyState() {
  return (
    <div className="card-shell flex flex-col items-center gap-3 px-6 py-20 text-center">
      <CalendarSearch size={30} className="text-muted-foreground/50" />
      <p className="font-display text-lg font-bold text-foreground">No matches to show yet</p>
      <p className="max-w-sm text-sm text-muted-foreground">
        Fixtures will appear here as soon as they&apos;re published. Check back closer to kickoff.
      </p>
    </div>
  );
}
