import React from 'react';
import type { Metadata } from 'next';
import { Users } from 'lucide-react';
import { getPerformances, getOverviewStats } from '@/lib/data';
import PerformanceCard from '@/components/PerformanceCard';
import Pagination from '@/components/Pagination';

export const revalidate = 60;

export const metadata: Metadata = {
  title: 'Performances',
  description: 'Rate individual player performances across the FIFA World Cup 2026.',
};

const PAGE_SIZE = 24;

export default async function PerformancesPage({
  searchParams,
}: {
  searchParams: Promise<{ page?: string }>;
}) {
  const sp = await Promise.resolve(searchParams);
  const page = Math.max(1, parseInt(sp?.page || '1', 10) || 1);
  const offset = (page - 1) * PAGE_SIZE;

  const [rawPerformances, overview] = await Promise.all([
    getPerformances(PAGE_SIZE, offset),
    getOverviewStats(),
  ]);
  // Only render real performances — guard against any null/undefined rows that
  // would otherwise show up as an empty trailing card in the grid.
  const performances = (rawPerformances || []).filter((p) => p && p.id);
  const totalPages = Math.max(1, Math.ceil((overview.performance_count || 0) / PAGE_SIZE));

  return (
    <main className="container-max container-px pb-12 pt-24">
      <header className="mb-8">
        <span className="section-label mb-2 block">Individual brilliance</span>
        <h1 className="font-display text-[clamp(2rem,4vw,2.75rem)] font-bold text-foreground">Performances</h1>
        <p className="mt-2 max-w-2xl text-muted-foreground">
          Every player display, rated out of 10 by the community. Find the moments that defined the tournament.
        </p>
      </header>

      {performances.length === 0 ? (
        <div className="card-shell flex flex-col items-center gap-3 px-6 py-20 text-center">
          <Users size={30} className="text-muted-foreground/50" />
          <p className="font-display text-lg font-bold text-foreground">No performances yet</p>
          <p className="max-w-sm text-sm text-muted-foreground">
            Player performances are added after each match. They&apos;ll show up here.
          </p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
            {performances.map((p) => (
              <PerformanceCard key={p.id} performance={p} />
            ))}
          </div>

          <Pagination currentPage={page} totalPages={totalPages} basePath="/performances" />
        </>
      )}
    </main>
  );
}
