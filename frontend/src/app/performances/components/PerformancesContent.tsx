'use client';

import React, { useState, useMemo } from 'react';
import {
  Filter,
  ChevronLeft,
  ChevronRight,
  ChevronDown,
  ChevronUp,
  Star,
  TrendingUp,
  Search,
  SlidersHorizontal,
} from 'lucide-react';

import { Performance } from '@/lib/api';

type SortKey = 'providerScore' | 'goals' | 'assists' | 'passAcc' | 'dribbles';
type SortDir = 'desc' | 'asc';

const stages = ['All Stages', 'Final', 'Semi-final', 'Quarter-final', 'Round of 16', 'Group Stage'];
const positions = ['All Positions', 'CF', 'LW', 'RW', 'AM', 'CM', 'DM', 'CB', 'LB', 'RB', 'GK'];
const ITEMS_PER_PAGE = 10;

function getScoreClass(score: number) {
  if (score >= 9.0) return 'score-high';
  if (score >= 7.0) return '';
  return 'score-low';
}

const stageBadgeMap: Record<string, string> = {
  Final: 'bg-primary/20 text-primary border border-primary/40',
  'Semi-final': 'bg-accent/10 text-accent border border-accent/30',
  'Quarter-final': 'bg-blue-500/10 text-blue-400 border border-blue-400/30',
  'Round of 16': 'bg-muted text-muted-foreground border border-border',
  'Group Stage': 'bg-muted text-muted-foreground border border-border',
};

interface PerformancesContentProps {
  initialPerformances?: Performance[];
  title?: string;
  description?: string;
  hideHeader?: boolean;
}

export default function PerformancesContent({ 
  initialPerformances = [],
  title = "Performance Archive",
  description = "Browse, sort, and filter through every logged performance. Discover hidden gems and debate player ratings.",
  hideHeader = false
}: PerformancesContentProps) {
  const [activeStage, setActiveStage] = useState('All Stages');
  const [activePosition, setActivePosition] = useState('All Positions');
  const [sortKey, setSortKey] = useState<SortKey>('providerScore');
  const [sortDir, setSortDir] = useState<SortDir>('desc');
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [expandedRow, setExpandedRow] = useState<string | null>(null);
  const [showFilters, setShowFilters] = useState(false);

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir((d) => (d === 'desc' ? 'asc' : 'desc'));
    } else {
      setSortKey(key);
      setSortDir('desc');
    }
    setCurrentPage(1);
  };

  const filtered = useMemo(() => {
    return (initialPerformances || [])
      .filter((p) => {
        const stage = p.match_round || 'Group Stage';
        const stageOk = activeStage === 'All Stages' || stage === activeStage;
        const posOk = true; // Positions are no longer in DB, so just skip this filter for now
        const searchOk =
          searchQuery === '' ||
          p.player_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
          p.team_name?.toLowerCase().includes(searchQuery.toLowerCase());
        return stageOk && posOk && searchOk;
      })
      .sort((a, b) => {
        let av = 0;
        let bv = 0;
        
        switch(sortKey) {
          case 'providerScore':
            av = a.provider_rating || 0;
            bv = b.provider_rating || 0;
            break;
          case 'goals':
            av = a.goals || 0;
            bv = b.goals || 0;
            break;
          case 'assists':
            av = a.assists || 0;
            bv = b.assists || 0;
            break;
          case 'passAcc':
            av = a.passes_accuracy || 0;
            bv = b.passes_accuracy || 0;
            break;
          case 'dribbles':
            av = a.dribbles || 0;
            bv = b.dribbles || 0;
            break;
        }

        return sortDir === 'desc' ? bv - av : av - bv;
      });
  }, [initialPerformances, activeStage, activePosition, sortKey, sortDir, searchQuery]);

  const totalPages = Math.ceil(filtered.length / ITEMS_PER_PAGE);
  const paginated = filtered.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE);

  function SortIcon({ col }: { col: SortKey }) {
    if (sortKey !== col) return <ChevronDown size={13} className="opacity-30" />;
    return sortDir === 'desc' ? (
      <ChevronDown size={13} className="text-primary" />
    ) : (
      <ChevronUp size={13} className="text-primary" />
    );
  }

  return (
    <div className={`${hideHeader ? 'pt-8' : 'pt-24'} pb-20 min-h-screen`}>
      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        {/* Page Header */}
        {!hideHeader && (
          <div className="mb-10">
            <p className="section-label mb-2">{initialPerformances.length} Performances Logged</p>
            <h1 className="font-display text-4xl md:text-5xl font-bold text-foreground mb-3">
              {title}
            </h1>
            <p className="text-muted-foreground text-base">
              {description}
            </p>
          </div>
        )}

        {/* Search + Filter Toggle */}
        <div className="flex items-center gap-3 mb-4">
          <div className="relative flex-1 max-w-sm">
            <Search size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search player or team..."
              value={searchQuery}
              onChange={(e) => { setSearchQuery(e.target.value); setCurrentPage(1); }}
              className="w-full pl-9 pr-4 py-2.5 rounded-lg bg-muted border border-border text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring transition-colors"
            />
          </div>
          <button
            onClick={() => setShowFilters((f) => !f)}
            className={`btn-ghost flex items-center gap-2 ${showFilters ? 'border-primary/40 text-primary bg-primary/5' : ''}`}
          >
            <SlidersHorizontal size={15} />
            Filters
          </button>
        </div>

        {/* Expandable Filters */}
        {showFilters && (
          <div className="archive-card p-5 mb-6 flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Stage</span>
              <div className="flex flex-wrap gap-2">
                {stages.map((stage) => (
                  <button
                    key={`filter-stage-${stage}`}
                    onClick={() => { setActiveStage(stage); setCurrentPage(1); }}
                    className={`filter-chip ${activeStage === stage ? 'active' : ''}`}
                  >
                    {stage}
                  </button>
                ))}
              </div>
            </div>
            <div className="flex flex-col gap-2">
              <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">Position</span>
              <div className="flex flex-wrap gap-2">
                {positions.map((pos) => (
                  <button
                    key={`filter-pos-${pos}`}
                    onClick={() => { setActivePosition(pos); setCurrentPage(1); }}
                    className={`filter-chip ${activePosition === pos ? 'active' : ''}`}
                  >
                    {pos}
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Sort Bar */}
        <div className="flex items-center gap-2 mb-4 overflow-x-auto scrollbar-dark pb-1">
          <span className="text-xs text-muted-foreground flex-shrink-0">Sort by:</span>
          {(
            [
              { key: 'providerScore', label: 'Score' },
              { key: 'goals', label: 'Goals' },
              { key: 'assists', label: 'Assists' },
              { key: 'passAcc', label: 'Pass Acc.' },
              { key: 'dribbles', label: 'Dribbles' },
            ] as { key: SortKey; label: string }[]
          ).map((s) => (
            <button
              key={`sort-${s.key}`}
              onClick={() => handleSort(s.key)}
              className={`filter-chip flex items-center gap-1 flex-shrink-0 ${sortKey === s.key ? 'active' : ''}`}
            >
              {s.label}
              <SortIcon col={s.key} />
            </button>
          ))}
        </div>

        {/* Results count */}
        <p className="text-xs text-muted-foreground mb-4 stat-number">
          {filtered.length} performances found
          {activeStage !== 'All Stages' ? ` · ${activeStage}` : ''}
          {activePosition !== 'All Positions' ? ` · ${activePosition}` : ''}
        </p>

        {/* Table */}
        {paginated.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-24 text-center archive-card">
            <TrendingUp size={40} className="text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold text-foreground mb-2">No performances found</h3>
            <p className="text-sm text-muted-foreground max-w-sm">
              Try clearing the search or adjusting your filters to browse the full archive.
            </p>
            <button
              onClick={() => { setSearchQuery(''); setActiveStage('All Stages'); setActivePosition('All Positions'); }}
              className="btn-ghost mt-4 text-sm"
            >
              Clear all filters
            </button>
          </div>
        ) : (
          <div className="archive-card overflow-hidden mb-8">
            {/* Table Header */}
            <div className="hidden md:grid grid-cols-[2fr_1.2fr_1fr_0.8fr_0.8fr_0.8fr_0.8fr_0.8fr_1fr_1fr] gap-3 px-5 py-3 border-b border-border">
              <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Player</span>
              <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Match</span>
              <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Stage</span>
              <button onClick={() => handleSort('goals')} className="flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground hover:text-foreground transition-colors">
                G <SortIcon col="goals" />
              </button>
              <button onClick={() => handleSort('assists')} className="flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground hover:text-foreground transition-colors">
                A <SortIcon col="assists" />
              </button>
              <button onClick={() => handleSort('dribbles')} className="flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground hover:text-foreground transition-colors">
                Dr. <SortIcon col="dribbles" />
              </button>
              <button onClick={() => handleSort('passAcc')} className="flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground hover:text-foreground transition-colors">
                Pass% <SortIcon col="passAcc" />
              </button>
              <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Min</span>
              <button onClick={() => handleSort('providerScore')} className="flex items-center gap-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground hover:text-foreground transition-colors">
                Score <SortIcon col="providerScore" />
              </button>
              <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">Community</span>
            </div>

            {/* Table Rows */}
            <div>
              {paginated.map((perf, rowIndex) => (
                <div key={perf.id}>
                  {/* Row */}
                  <div
                    className={`performance-row-hover cursor-pointer transition-colors ${rowIndex % 2 === 0 ? '' : 'bg-white/[0.01]'} ${expandedRow === perf.id ? 'bg-primary/5' : ''}`}
                    onClick={() => setExpandedRow(expandedRow === perf.id ? null : perf.id)}
                  >
                    {/* Mobile card layout */}
                    <div className="md:hidden p-4 flex items-start gap-3 border-b border-border">
                      <div className="w-9 h-9 rounded-full bg-muted border border-border flex items-center justify-center text-lg flex-shrink-0">
                        {perf.flag}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2">
                          <div>
                            <p className="text-sm font-semibold text-foreground">{perf.player_name}</p>
                            <p className="text-xs text-muted-foreground">{perf.team_name} · N/A · {perf.match_title}</p>
                          </div>
                          <div className={`performance-score w-9 h-9 text-xs flex-shrink-0 ${getScoreClass(perf.provider_rating || 0)}`}>
                            {perf.provider_rating}
                          </div>
                        </div>
                        <div className="flex items-center gap-3 mt-2">
                          <span className="stat-pill">{perf.goals}G</span>
                          <span className="stat-pill">{perf.assists}A</span>
                          <span className="stat-pill">{perf.passes_accuracy}%</span>
                          <span className={`stage-badge text-xs ${stageBadgeMap[perf.match_round || 'Group Stage'] || ''}`}>{perf.match_round}</span>
                        </div>
                      </div>
                    </div>

                    {/* Desktop table row */}
                    <div className="hidden md:grid grid-cols-[2fr_1.2fr_1fr_0.8fr_0.8fr_0.8fr_0.8fr_0.8fr_1fr_1fr] gap-3 px-5 py-4 border-b border-border items-center">
                      {/* Player */}
                      <div className="flex items-center gap-3 min-w-0">
                        <div className="w-8 h-8 rounded-full bg-muted border border-border flex items-center justify-center text-base flex-shrink-0">
                          {perf.flag_emoji || '🌍'}
                        </div>
                        <div className="min-w-0">
                          <p className="text-sm font-semibold text-foreground truncate">{perf.player_name}</p>
                          <p className="text-xs text-muted-foreground truncate">{perf.team_name} · N/A</p>
                        </div>
                      </div>
                      {/* Match */}
                      <span className="text-sm text-muted-foreground truncate">{perf.match_title}</span>
                      {/* Stage */}
                      <span className={`stage-badge text-xs w-fit ${stageBadgeMap[perf.match_round || 'Group Stage'] || ''}`}>{perf.match_round}</span>
                      {/* Goals */}
                      <span className="stat-number text-sm font-semibold text-foreground">{perf.goals}</span>
                      {/* Assists */}
                      <span className="stat-number text-sm font-semibold text-foreground">{perf.assists}</span>
                      {/* Dribbles */}
                      <span className="stat-number text-sm text-muted-foreground">{perf.dribbles}</span>
                      {/* Pass Acc */}
                      <span className="stat-number text-sm text-muted-foreground">{perf.passes_accuracy}%</span>
                      {/* Minutes */}
                      <span className="stat-number text-sm text-muted-foreground">{perf.minutes_played || 90}</span>
                      {/* Provider Score */}
                      <div className={`performance-score ${getScoreClass(perf.provider_rating || 0)}`}>
                        {perf.provider_rating}
                      </div>
                      {/* Community */}
                      <div className="flex items-center gap-1.5">
                        {(perf.average_rating || 0) > 0 ? (
                          <>
                            <Star size={11} className="text-primary fill-primary" />
                            <span className="stat-number text-sm font-semibold text-primary">{perf.average_rating}</span>
                            <span className="text-xs text-muted-foreground">({perf.total_votes})</span>
                          </>
                        ) : (
                          <span className="text-xs text-muted-foreground">—</span>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Expanded Row Detail */}
                  {expandedRow === perf.id && (
                    <div className="bg-primary/5 border-b border-border px-5 py-5">
                      <div className="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-6 gap-4">
                        {[
                          { label: 'Goals', value: perf.goals },
                          { label: 'Assists', value: perf.assists },
                          { label: 'Pass Acc.', value: `${perf.passes_accuracy}%` },
                          { label: 'Dribbles', value: perf.dribbles },
                          { label: 'Minutes', value: perf.minutes_played || 90 },
                          { label: 'Provider Score', value: perf.provider_rating },
                          { label: 'Community', value: (perf.average_rating || 0) > 0 ? perf.average_rating : '—' },
                          { label: 'Reviews', value: perf.total_votes || 0 },
                          { label: 'Stage', value: perf.match_round },
                        ].map((stat) => (
                          <div
                            key={`expanded-${perf.id}-${stat.label}`}
                            className="flex flex-col gap-1 p-3 rounded-lg bg-card border border-border"
                          >
                            <span className="text-xs text-muted-foreground font-medium">{stat.label}</span>
                            <span className="stat-number text-base font-bold text-foreground">{stat.value}</span>
                          </div>
                        ))}
                      </div>
                      <div className="mt-4 flex items-center gap-3">
                        <button className="btn-primary text-xs px-4 py-2 flex items-center gap-1.5">
                          <Star size={13} />
                          Rate this Performance
                        </button>
                        <span className="text-xs text-muted-foreground">Sign in to rate and review</span>
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground stat-number">
              Showing {(currentPage - 1) * ITEMS_PER_PAGE + 1}–
              {Math.min(currentPage * ITEMS_PER_PAGE, filtered.length)} of {filtered.length}
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
                  key={`perf-page-${page}`}
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