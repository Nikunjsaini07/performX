import React from 'react';
import Link from 'next/link';
import { ArrowRight } from 'lucide-react';
import { getTrendingPlayers } from '@/lib/api';
import StatLeadersChart from './StatLeadersChart';

export default async function StatLeaders() {
  let chartData: any[] = [];
  try {
    const res = await getTrendingPlayers(8);
    if (res && res.data) {
      chartData = res.data.map(item => ({
        player: item.entity.name,
        team: item.entity.team_name || 'Unknown',
        value: item.score,
        flag: item.entity.photo_url ? '👤' : '⚽',
      }));
    }
  } catch (error) {
    console.error("Failed to fetch trending players:", error);
  }

  return (
    <section className="py-24 bg-background relative">
      {/* Atmospheric background */}
      <div className="absolute inset-0 pointer-events-none overflow-hidden">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] bg-primary/3 blur-[150px] rounded-full" />
      </div>
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[60%] h-px bg-gradient-to-r from-transparent via-border to-transparent" />

      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10 relative z-10">
        <div className="flex items-end justify-between mb-12">
          <div>
            <p className="section-label mb-2">Tournament Stats</p>
            <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground">
              Trending Players
            </h2>
            <p className="text-sm text-muted-foreground mt-2">
              The most discussed players right now based on our algorithm
            </p>
          </div>
          <Link
            href="/performances"
            className="hidden sm:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-primary transition-colors group"
          >
            Full stats
            <ArrowRight size={14} className="group-hover:translate-x-0.5 transition-transform" />
          </Link>
        </div>

        <StatLeadersChart data={chartData} />
      </div>
    </section>
  );
}