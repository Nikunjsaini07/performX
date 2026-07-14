'use client';

import React, { useState, useMemo } from 'react';
import { Filter, ChevronLeft, ChevronRight, Trophy, MapPin, Calendar, Users } from 'lucide-react';
import { Match } from '@/lib/api';

const stages = ['All Stages', 'Final', 'Semi-final', 'Quarter-final', 'Round of 16', 'Group Stage'];
const groups = ['All Groups', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L'];
const ITEMS_PER_PAGE = 12;

const stageBadgeMap: Record<string, string> = {
  Final: 'bg-primary/20 text-primary border border-primary/40',
  'Semi-final': 'bg-accent/10 text-accent border border-accent/30',
  'Quarter-final': 'bg-blue-500/10 text-blue-400 border border-blue-400/30',
  'Round of 16': 'bg-muted text-muted-foreground border border-border',
  'Group Stage': 'bg-muted text-muted-foreground border border-border',
};

function penaltyValue(value?: Match['home_penalty_score']) {
  if (typeof value === 'number') return value;
  if (value && typeof value === 'object' && value.Valid) return value.Int32;
  return null;
}

export default function MatchesContent({ initialMatches }: { initialMatches: Match[] }) {
  const [activeStage, setActiveStage] = useState('All Stages');
  const [activeGroup, setActiveGroup] = useState('All Groups');
  const [currentPage, setCurrentPage] = useState(1);

  const filtered = useMemo(() => {
    return initialMatches.filter((m) => {
      const stageOk = activeStage === 'All Stages' || m.round === activeStage || (!m.round && activeStage === 'Group Stage');
      const groupOk = activeGroup === 'All Groups' || m.group === activeGroup;
      return stageOk && groupOk;
    });
  }, [activeStage, activeGroup]);

  const totalPages = Math.ceil(filtered.length / ITEMS_PER_PAGE);
  const paginated = filtered.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE);

  const handleStageChange = (stage: string) => {
    setActiveStage(stage);
    setCurrentPage(1);
    if (stage !== 'Group Stage') setActiveGroup('All Groups');
  };

  return (
    <div className="pt-24 pb-20 min-h-screen">
      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        {/* Page Header */}
        <div className="mb-10">
          <p className="section-label mb-2">FIFA World Cup 2026</p>
          <h1 className="font-display text-4xl md:text-5xl font-bold text-foreground mb-3">
            Match Archive
          </h1>
          <p className="text-muted-foreground text-base">
            All 96 tournament matches — scorelines, venues, and top performances
          </p>
        </div>

        {/* Archive Summary Bar */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-10">
          {[
            { label: 'Total Matches', value: '96' },
            { label: 'Teams', value: '48' },
            { label: 'Stages', value: '6' },
            { label: 'Showing', value: `${filtered.length}` },
          ].map((s) => (
            <div key={`summary-${s.label}`} className="archive-card px-4 py-3 flex items-center justify-between">
              <span className="text-sm text-muted-foreground">{s.label}</span>
              <span className="stat-number text-xl font-bold text-foreground">{s.value}</span>
            </div>
          ))}
        </div>

        {/* Filters */}
        <div className="flex flex-col gap-4 mb-8">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="flex items-center gap-1.5 text-xs text-muted-foreground mr-2">
              <Filter size={13} />
              Stage
            </span>
            {stages.map((stage) => (
              <button
                key={`stage-${stage}`}
                onClick={() => handleStageChange(stage)}
                className={`filter-chip ${activeStage === stage ? 'active' : ''}`}
              >
                {stage}
              </button>
            ))}
          </div>
          {activeStage === 'Group Stage' && (
            <div className="flex items-center gap-2 flex-wrap">
              <span className="flex items-center gap-1.5 text-xs text-muted-foreground mr-2">
                <Filter size={13} />
                Group
              </span>
              {groups.map((group) => (
                <button
                  key={`group-${group}`}
                  onClick={() => { setActiveGroup(group); setCurrentPage(1); }}
                  className={`filter-chip ${activeGroup === group ? 'active' : ''}`}
                >
                  {group}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Match Grid */}
        {paginated.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <Calendar size={40} className="text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold text-foreground mb-2">No matches found</h3>
            <p className="text-sm text-muted-foreground">
              Try adjusting the stage or group filter to find matches.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-5 mb-10">
            {paginated.map((match) => {
              const dateObj = new Date(match.utc_datetime);
              const dateStr = dateObj.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
              const homePen = penaltyValue(match.home_penalty_score);
              const awayPen = penaltyValue(match.away_penalty_score);
              return (
              <div
                key={match.id}
                className="archive-card match-card-shine card-glow-gold group cursor-pointer"
                onClick={() => window.location.href = `/matches/${match.slug}`}
              >
                <div className="p-5">
                  {/* Stage + Date */}
                  <div className="flex items-center justify-between mb-4">
                    <span className={`stage-badge text-xs ${stageBadgeMap[match.round || 'Group Stage'] || 'bg-muted text-muted-foreground'}`}>
                      {match.round || 'Group Stage'}
                    </span>
                    <span className="text-xs text-muted-foreground stat-number">{dateStr}</span>
                  </div>

                  {/* Teams + Score */}
                  <div className="flex items-center justify-between gap-3 mb-3">
                    <div className="flex flex-col items-center gap-1.5 flex-1">
                      <span className="text-3xl text-center leading-none">{match.home_team_logo_url || '🏳️'}</span>
                      <span className="text-xs font-semibold text-foreground text-center leading-tight">
                        {match.home_team_name}
                      </span>
                    </div>
                    <div className="flex flex-col items-center gap-1">
                      <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-muted">
                        <span className="stat-number text-2xl font-bold text-foreground">
                          {match.home_score ?? 0}
                        </span>
                        <span className="text-muted-foreground">–</span>
                        <span className="stat-number text-2xl font-bold text-foreground">
                          {match.away_score ?? 0}
                        </span>
                      </div>
                      {homePen !== null && awayPen !== null && (
                        <span className="text-[10px] font-semibold text-accent stat-number tracking-wider mt-1">
                          ({homePen}) PEN ({awayPen})
                        </span>
                      )}
                      {match.tagline && (
                        <span className="text-[10px] text-accent/80 font-medium truncate max-w-[80px]">{match.tagline}</span>
                      )}
                    </div>
                    <div className="flex flex-col items-center gap-1.5 flex-1">
                      <span className="text-3xl text-center leading-none">{match.away_team_logo_url || '🏳️'}</span>
                      <span className="text-xs font-semibold text-foreground text-center leading-tight">
                        {match.away_team_name}
                      </span>
                    </div>
                  </div>

                  {/* Venue */}
                  <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-4 border-t border-border pt-3">
                    <MapPin size={11} />
                    <span className="truncate">{match.venue || 'TBD Venue'}</span>
                  </div>
                </div>
              </div>
            )})}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground stat-number">
              Showing {(currentPage - 1) * ITEMS_PER_PAGE + 1}–
              {Math.min(currentPage * ITEMS_PER_PAGE, filtered.length)} of {filtered.length} matches
            </p>
            <div className="flex items-center gap-2">
              <button
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                aria-label="Previous page"
              >
                <ChevronLeft size={16} />
              </button>
              {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                <button
                  key={`page-${page}`}
                  onClick={() => setCurrentPage(page)}
                  className={`w-9 h-9 rounded-lg text-sm font-medium transition-colors ${
                    page === currentPage
                      ? 'bg-primary text-primary-foreground'
                      : 'border border-border text-muted-foreground hover:text-foreground hover:bg-muted'
                  }`}
                >
                  {page}
                </button>
              ))}
              <button
                onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
                className="p-2 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-muted disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                aria-label="Next page"
              >
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
