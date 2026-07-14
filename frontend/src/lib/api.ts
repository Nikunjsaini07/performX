export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Match {
  id: string;
  title: string;
  slug: string;
  description?: string;
  round?: string;
  utc_datetime: string;
  venue?: string;
  cover_image_url?: string;
  home_score: number;
  away_score: number;
  home_penalty_score?: number | { Int32: number; Valid: boolean };
  away_penalty_score?: number | { Int32: number; Valid: boolean };
  home_team_name: string;
  home_team_short_name: string;
  home_team_logo_url?: string;
  away_team_name: string;
  away_team_short_name: string;
  away_team_logo_url?: string;
  tagline?: string;
  average_rating?: number;
  total_votes?: number;
  group?: string;
}

export interface Player {
  id: string;
  name: string;
  photo_url: string;
  team_name?: string;
  position?: string;
  jersey_number?: number;
}

export interface Performance {
  id: string;
  title: string;
  description?: string;
  match_id: string;
  player_id: string;
  player_team_id: string;
  cover_image_url?: string;
  player_photo_url?: string;
  player_slug?: string;
  player_name: string;
  match_title: string;
  match_slug?: string;
  team_name?: string;
  team_logo_url?: string;
  flag_emoji?: string;
  flag?: string;
  match_round?: string;
  jersey_number?: number;
  is_starter?: boolean;
  captain?: boolean;
  position?: string;
  goals?: number;
  assists?: number;
  dribbles?: number;
  passes_accuracy?: number;
  average_rating?: number;
  average_rating_2?: number;
  provider_rating?: number;
  total_votes?: number;
  minutes_played?: number;
}

export interface Team {
  id: string;
  name: string;
  short_name: string;
  slug: string;
  type?: string;
  logo_url?: string;
  flag_emoji?: string;
}

export interface Review {
  id: string;
  user_id: string;
  rating: number;
  title?: string;
  content?: string;
  comment?: string;
  review_text?: string;
  created_at: string;
  user_name?: string;
  username?: string;
  author_name?: string;
  author_avatar?: string;
  like_count?: number;
  comment_count?: number;
  // Potentially joined fields depending on review type
  match_title?: string;
  player_name?: string;
}

export interface StatRow {
  id?: string;
  stat_name?: string;
  stat_short_name?: string;
  stat_type_id?: string;
  value: number | string;
  stat_unit?: string;
}

export interface TrendingResponse<T> {
  data: {
    entity_id: string;
    entity_type: string;
    score: number;
    rank: number;
    entity: T;
  }[];
  meta: {
    window: string;
    limit: number;
    type: string;
  };
}

export async function getTrendingMatches(limit = 4): Promise<TrendingResponse<Match>> {
  const res = await fetch(`${API_URL}/trending/matches?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch trending matches');
  }
  return res.json();
}

export async function getTrendingPerformances(limit = 4): Promise<TrendingResponse<Performance>> {
  const res = await fetch(`${API_URL}/trending/performances?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch trending performances');
  }
  return res.json();
}

export async function getMatches(limit = 8): Promise<Match[]> {
  const res = await fetch(`${API_URL}/matches?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch matches');
  }
  return res.json();
}

export async function getMatch(slug: string): Promise<Match> {
  const res = await fetch(`${API_URL}/matches/${slug}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch match ${slug}`);
  }
  return res.json();
}

export async function getPerformances(limit = 6): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/performances?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch performances');
  }
  return res.json();
}

export async function getTopRatedPerformances(limit = 6): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/performances/top-rated?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch top rated performances');
  }
  return res.json();
}

export async function getPerformance(id: string): Promise<Performance> {
  const res = await fetch(`${API_URL}/performances/${id}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch performance ${id}`);
  }
  return res.json();
}

export async function getPerformanceStats(id: string): Promise<StatRow[]> {
  const res = await fetch(`${API_URL}/performances/${id}/stats`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch stats for performance ${id}`);
  }
  return res.json();
}

export async function getPerformanceReviews(id: string): Promise<Review[]> {
  const res = await fetch(`${API_URL}/performances/${id}/reviews`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch reviews for performance ${id}`);
  }
  return res.json();
}

export async function getMatchReviews(slug: string): Promise<Review[]> {
  const res = await fetch(`${API_URL}/matches/${slug}/reviews?limit=12`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch reviews for match ${slug}`);
  }
  return res.json();
}

export async function getMatchPerformances(slug: string): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/matches/${slug}/performances`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch performances for match ${slug}`);
  }
  return res.json();
}

export async function getTeams(limit = 4): Promise<Team[]> {
  const res = await fetch(`${API_URL}/teams?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch teams');
  }
  return res.json();
}

export async function getTrendingPlayers(limit = 3): Promise<TrendingResponse<Player>> {
  const res = await fetch(`${API_URL}/trending/players?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch trending players');
  }
  return res.json();
}

export async function getTrendingReviews(limit = 4): Promise<TrendingResponse<Review>> {
  const res = await fetch(`${API_URL}/trending/reviews?limit=${limit}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    throw new Error('Failed to fetch trending reviews');
  }
  return res.json();
}
