import React from 'react';
import Link from 'next/link';
import { ArrowRight, Star, Zap } from 'lucide-react';
import { getTopRatedPerformances, Performance } from '@/lib/api';

function getScoreClass(score: number) {
  if (score >= 9.0) return 'score-high';
  if (score >= 7.0) return '';
  return 'score-low';
}

const rankColors = ['text-primary', 'text-muted-foreground', 'text-[#cd7f32]'];

export default async function TopPerformances() {
  let performances: Performance[] = [];
  try {
    performances = await getTopRatedPerformances(6);
  } catch (error) {
    console.error("Failed to fetch top performances:", error);
  }

  return (
    <section className="py-24 bg-background relative">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[60%] h-px bg-gradient-to-r from-transparent via-border to-transparent" />

      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        {/* Header */}
        <div className="flex items-end justify-between mb-12">
          <div>
            <p className="text-[#c9a84c] text-sm font-bold tracking-widest uppercase mb-2">Recent Update</p>
            <h2 className="font-display text-3xl md:text-5xl font-bold text-white tracking-tight" style={{ fontFamily: "Georgia, serif" }}>
              Recent Performances
            </h2>
            <p className="text-blue-400 text-sm md:text-base mt-3 font-medium">
              The latest standout performances from the tournament
            </p>
          </div>
          <Link
            href="/performances"
            className="hidden sm:flex items-center gap-1.5 text-sm text-muted-foreground hover:text-primary transition-colors group"
          >
            All performances
            <ArrowRight size={14} className="group-hover:translate-x-0.5 transition-transform" />
          </Link>
        </div>

        {/* Portrait Performance Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
          {performances.length === 0 ? (
            <div className="col-span-full py-20 text-center text-muted-foreground">
              No performances found. Ensure the API is running.
            </div>
          ) : (
            performances.map((perf, index) => {
              // We initialize average_rating to the provider's rating in the database, so we just use that directly.
              const score = perf.average_rating ? Number(perf.average_rating) : 0;
              
              return (
                <Link
                  href={`/performances/${perf.id}`}
                  key={perf.id}
                  className="group block"
                >
                  <div
                    className="relative overflow-hidden rounded-xl border border-[#333] bg-[#111] transition-all duration-300 hover:border-[#c9a84c]/50 hover:-translate-y-1 hover:shadow-xl hover:shadow-[#c9a84c]/10 flex flex-col font-sans"
                    style={{ minHeight: '380px' }}
                  >
                    {/* Top Half: Player Image / Avatar */}
                    <div className="w-full h-48 bg-gradient-to-b from-[#1a1a1a] to-[#111] relative border-b border-[#222]">
                      {/* Score badge — top right */}
                      <div className="absolute top-3 right-3 z-10">
                        <div className="flex items-center justify-center bg-[#c9a84c] text-black font-black text-sm px-2.5 py-1 rounded-md shadow-md">
                          {score > 0 ? score.toFixed(1) : '-'}
                        </div>
                      </div>

                      {perf.cover_image_url ? (
                        <img 
                          src={perf.cover_image_url} 
                          alt={perf.player_name}
                          className="w-full h-full object-cover object-top opacity-90 group-hover:opacity-100 transition-opacity"
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-6xl opacity-20">
                          👤
                        </div>
                      )}
                    </div>

                    {/* Bottom Half: Player Info area */}
                    <div className="flex flex-col items-center pt-5 pb-4 px-4 flex-1 bg-[#151515]">
                      {/* Player name */}
                      <h3 className="font-display text-white text-xl uppercase font-bold tracking-wide text-center leading-tight mb-1">
                        {perf.player_name}
                      </h3>
                      
                      {/* Match title */}
                      <p className="text-[#888] text-sm text-center mb-0.5">
                        {perf.match_title}
                      </p>
                      
                      {/* Round */}
                      <p className="text-[#666] text-xs text-center mb-4">
                        {perf.match_round || 'Group Stage'}
                      </p>

                      {/* Quote/Tagline */}
                      <p className="text-[#aaa] text-sm text-center italic leading-relaxed line-clamp-3 px-2 mt-auto" style={{ fontFamily: "Georgia, serif" }}>
                        "{perf.title}"
                      </p>
                    </div>

                    {/* Footer Section */}
                    <div className="p-4 border-t border-[#333] bg-[#111] flex items-center justify-between mt-auto">
                      <div className="flex items-center gap-2">
                        <span className="text-[#c9a84c] text-lg">★</span>
                        <div className="flex items-baseline gap-1">
                          <span className="text-white font-bold text-lg">{score > 0 ? score.toFixed(1) : '0'}</span>
                          <span className="text-[#666] text-xs font-semibold">/10</span>
                        </div>
                        <span className="text-[#666] text-xs ml-1">({perf.total_votes || 0})</span>
                      </div>
                      
                      <div className="px-4 py-1.5 border border-[#c9a84c] text-[#c9a84c] rounded-lg text-sm font-semibold hover:bg-[#c9a84c] hover:text-black transition-colors flex items-center gap-1.5">
                        <span className="text-sm leading-none">☆</span> Rate
                      </div>
                    </div>
                  </div>
                </Link>
              );
            })
          )}
        </div>

        <div className="flex sm:hidden justify-center mt-6">
          <Link href="/performances" className="btn-ghost text-sm flex items-center gap-1.5">
            All performances
            <ArrowRight size={14} />
          </Link>
        </div>
      </div>
    </section>
  );
}
