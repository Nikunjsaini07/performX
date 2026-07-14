// Server-side safe data helpers. Every API call is wrapped so that a failed
// fetch (e.g. the API being offline during build) returns an empty/safe
// fallback instead of throwing — keeping pages and the build resilient.

import * as api from '@/lib/api';
import type {
  Match,
  Performance,
  Team,
  TeamMatch,
  TeamPerformance,
  Player,
  Review,
  StatRow,
  TrendingResponse,
  UserProfile,
  UserReviewRow,
  UserRatingRow,
  OverviewStats,
} from '@/lib/api';

async function safe<T>(fn: () => Promise<T>, fallback: T): Promise<T> {
  try {
    return await fn();
  } catch {
    return fallback;
  }
}

/** Map a trending response to a clean, null-filtered list of entities. */
function mapTrending<T>(res: TrendingResponse<T> | null | undefined): T[] {
  if (!res || !Array.isArray(res.data)) return [];
  return res.data
    .map((d) => d?.entity)
    .filter((e): e is T => e != null);
}

// ─── Matches ────────────────────────────────────────────────────────────────

export const getMatches = (limit = 24, offset = 0): Promise<Match[]> =>
  safe(() => api.getMatches(limit, offset), []);

export const getMatch = (slug: string): Promise<Match | null> =>
  safe<Match | null>(() => api.getMatch(slug), null);

export const getMatchPerformances = (slug: string): Promise<Performance[]> =>
  safe(() => api.getMatchPerformances(slug), []);

export const getMatchReviews = (slug: string, limit = 20): Promise<Review[]> =>
  safe(() => api.getMatchReviews(slug, limit), []);

export const getTrendingMatches = async (limit = 12): Promise<Match[]> =>
  mapTrending(await safe(() => api.getTrendingMatches(limit), null));

// ─── Performances ─────────────────────────────────────────────────────────────

export const getPerformances = (limit = 24, offset = 0): Promise<Performance[]> =>
  safe(() => api.getPerformances(limit, offset), []);

export const getPerformance = (id: string): Promise<Performance | null> =>
  safe<Performance | null>(() => api.getPerformance(id), null);

export const getPerformanceStats = (id: string): Promise<StatRow[]> =>
  safe(() => api.getPerformanceStats(id), []);

export const getPerformanceReviews = (id: string, limit = 20): Promise<Review[]> =>
  safe(() => api.getPerformanceReviews(id, limit), []);

export const getTrendingPerformances = async (limit = 12): Promise<Performance[]> =>
  mapTrending(await safe(() => api.getTrendingPerformances(limit), null));

export const getTopRatedPerformances = (limit = 12): Promise<Performance[]> =>
  safe(() => api.getTopRatedPerformances(limit), []);

// ─── Teams & Players ──────────────────────────────────────────────────────────

export const getTeams = (limit = 64): Promise<Team[]> =>
  safe(() => api.getTeams(limit), []);

export const getTeamMatches = (slug: string, limit = 20, offset = 0): Promise<TeamMatch[]> =>
  safe(() => api.getTeamMatches(slug, limit, offset), []);

export const getTeamPerformances = (slug: string, limit = 24, offset = 0): Promise<TeamPerformance[]> =>
  safe(() => api.getTeamPerformances(slug, limit, offset), []);

// ─── Overview stats ───────────────────────────────────────────────────────────

const EMPTY_STATS: OverviewStats = {
  match_count: 0,
  performance_count: 0,
  team_count: 0,
  rating_count: 0,
  review_count: 0,
};

export const getOverviewStats = (): Promise<OverviewStats> =>
  safe(() => api.getOverviewStats(), EMPTY_STATS);

export const getTrendingPlayers = async (limit = 10): Promise<Player[]> =>
  mapTrending(await safe(() => api.getTrendingPlayers(limit), null));

export const getTrendingReviews = async (limit = 6): Promise<Review[]> =>
  mapTrending(await safe(() => api.getTrendingReviews(limit), null));

// ─── Users ──────────────────────────────────────────────────────────────────

export const getUserProfile = (username: string): Promise<UserProfile | null> =>
  safe<UserProfile | null>(() => api.getUserProfile(username), null);

export const getUserReviews = (username: string, limit = 20): Promise<UserReviewRow[]> =>
  safe(() => api.getUserReviews(username, limit), []);

export const getUserRatings = (username: string, limit = 20): Promise<UserRatingRow[]> =>
  safe(() => api.getUserRatings(username, limit), []);

// Also expose the trending player ranking (with score) for ranked lists.
export const getTrendingPlayersRanked = async (
  limit = 10,
): Promise<{ player: Player; score: number; rank: number }[]> => {
  const res = await safe(() => api.getTrendingPlayers(limit), null);
  if (!res || !Array.isArray(res.data)) return [];
  return res.data
    .filter((d) => d?.entity != null)
    .map((d) => ({ player: d.entity, score: d.score, rank: d.rank }));
};
