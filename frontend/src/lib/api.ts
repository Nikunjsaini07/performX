export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// ─── Interfaces ──────────────────────────────────────────────────────────────

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
  slug?: string;
}

export interface Performance {
  id: string;
  title: string;
  description?: string;
  /** Human-readable slug (e.g. "ricardo-rodriguez-vs-argentina-r8e6zh"). Prefer
   * this over `id` for URLs — falls back to `id` if absent. */
  slug?: string;
  /** Short punchy one-liner, distinct from the longer `description`. */
  tagline?: string;
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

export interface User {
  id: string;
  username: string;
  display_name: string;
  email: string;
  bio?: string;
  avatar_url?: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface Review {
  id: string;
  user_id: string;
  rating?: number;
  title?: string;
  content?: string;
  comment?: string;
  review_text?: string;
  body?: string;
  created_at: string;
  user_name?: string;
  username?: string;
  display_name?: string;
  author_name?: string;
  author_username?: string;
  author_display_name?: string;
  author_avatar?: string;
  author_avatar_url?: string;
  avatar_url?: string;
  like_count?: number;
  likes_count?: number;
  comment_count?: number;
  liked_by_me?: boolean;
  match_title?: string;
  match_slug?: string;
  player_name?: string;
  performance_id?: string;
}

export interface Comment {
  id: string;
  review_id: string;
  user_id: string;
  body: string;
  created_at: string;
  username?: string;
  display_name?: string;
  avatar_url?: string;
  like_count?: number;
  liked_by_me?: boolean;
}

export interface StatRow {
  id?: string;
  stat_name?: string;
  stat_short_name?: string;
  stat_type_id?: string;
  value: number | string;
  stat_unit?: string;
}

export interface RatingResponse {
  match_id?: string;
  performance_id?: string;
  user_id: string;
  rating: number;
  average_rating?: number;
  total_votes?: number;
}

export interface UserRating {
  rating: number;
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

// ─── Helpers ─────────────────────────────────────────────────────────────────

function authHeaders(token?: string | null): HeadersInit {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  return headers;
}

// ─── Auth API ────────────────────────────────────────────────────────────────

export async function loginUser(email: string, password: string): Promise<AuthTokens> {
  const res = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Login failed');
  }
  return res.json();
}

export async function registerUser(
  username: string,
  display_name: string,
  email: string,
  password: string,
): Promise<{ message: string; email: string; otp_code_dev?: string }> {
  const res = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, display_name, email, password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Registration failed');
  }
  return res.json();
}

export async function verifyOtp(
  email: string,
  otp_code: string,
  purpose: 'REGISTER' | 'LOGIN',
): Promise<AuthTokens> {
  const res = await fetch(`${API_URL}/auth/verify-otp`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, otp_code, purpose }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'OTP verification failed');
  }
  return res.json();
}

export async function refreshAccessToken(refresh_token: string): Promise<{ access_token: string; refresh_token: string }> {
  const res = await fetch(`${API_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token }),
  });
  if (!res.ok) throw new Error('Session expired');
  return res.json();
}

export async function logoutUser(token: string, refresh_token: string): Promise<void> {
  await fetch(`${API_URL}/auth/logout`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ refresh_token }),
  });
}

export async function forgotPassword(email: string): Promise<void> {
  const res = await fetch(`${API_URL}/auth/forgot-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to send reset code');
  }
}

export async function resetPassword(email: string, otp_code: string, new_password: string): Promise<void> {
  const res = await fetch(`${API_URL}/auth/reset-password`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, otp_code, new_password }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to reset password');
  }
}

export async function getMe(token: string): Promise<User> {
  const res = await fetch(`${API_URL}/me`, { headers: authHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch user');
  return res.json();
}

export async function updateProfile(token: string, data: { display_name?: string; bio?: string }): Promise<User> {
  const res = await fetch(`${API_URL}/me`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update profile');
  return res.json();
}

export async function updateUsername(token: string, username: string): Promise<User> {
  const res = await fetch(`${API_URL}/me/username`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify({ username }),
  });
  if (!res.ok) throw new Error('Failed to update username');
  return res.json();
}

export async function updateAvatar(token: string, avatar_url: string): Promise<User> {
  const res = await fetch(`${API_URL}/me/avatar`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify({ avatar_url }),
  });
  if (!res.ok) throw new Error('Failed to update avatar');
  return res.json();
}

export async function uploadToCloudinary(file: File): Promise<string> {
  const cloudName = process.env.NEXT_PUBLIC_CLOUDINARY_CLOUD_NAME || 'dbcflpua9';
  const uploadPreset = process.env.NEXT_PUBLIC_CLOUDINARY_UPLOAD_PRESET || 'unsigned_preset';

  const formData = new FormData();
  formData.append('file', file);
  formData.append('upload_preset', uploadPreset);

  const res = await fetch(`https://api.cloudinary.com/v1_1/${cloudName}/image/upload`, {
    method: 'POST',
    body: formData,
  });

  if (!res.ok) throw new Error('Failed to upload image');
  const data = await res.json();
  return data.secure_url;
}

// Delete review functions
export async function deleteMatchReview(token: string, matchSlug: string, reviewId: string): Promise<void> {
  const res = await fetch(`${API_URL}/matches/${matchSlug}/reviews/${reviewId}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  if (!res.ok) throw new Error('Failed to delete review');
}

export async function deletePerformanceReview(token: string, performanceId: string, reviewId: string): Promise<void> {
  const res = await fetch(`${API_URL}/performances/${performanceId}/reviews/${reviewId}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  if (!res.ok) throw new Error('Failed to delete review');
}

// Update review functions
export async function updateMatchReview(token: string, matchSlug: string, reviewId: string, data: { title?: string; content?: string }): Promise<void> {
  const res = await fetch(`${API_URL}/matches/${matchSlug}/reviews/${reviewId}`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update review');
}

export async function updatePerformanceReview(token: string, performanceId: string, reviewId: string, data: { title?: string; content?: string }): Promise<void> {
  const res = await fetch(`${API_URL}/performances/${performanceId}/reviews/${reviewId}`, {
    method: 'PATCH',
    headers: authHeaders(token),
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update review');
}

// ─── Public Reads — Overview Stats ───────────────────────────────────────────

export interface OverviewStats {
  match_count: number;
  performance_count: number;
  team_count: number;
  rating_count: number;
  review_count: number;
}

export async function getOverviewStats(): Promise<OverviewStats> {
  const res = await fetch(`${API_URL}/stats/overview`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch overview stats');
  return res.json();
}

// ─── Public Reads — Matches ───────────────────────────────────────────────────

export async function getTrendingMatches(limit = 4): Promise<TrendingResponse<Match>> {
  const res = await fetch(`${API_URL}/trending/matches?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch trending matches');
  return res.json();
}

export async function getMatches(limit = 8, offset = 0): Promise<Match[]> {
  const res = await fetch(`${API_URL}/matches?limit=${limit}&offset=${offset}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch matches');
  return res.json();
}

export async function getMatch(slug: string): Promise<Match> {
  const res = await fetch(`${API_URL}/matches/${slug}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error(`Failed to fetch match ${slug}`);
  return res.json();
}

export async function getMatchReviews(slug: string, limit = 20): Promise<Review[]> {
  const res = await fetch(`${API_URL}/matches/${slug}/reviews?limit=${limit}`, { next: { revalidate: 30 } });
  if (!res.ok) throw new Error(`Failed to fetch reviews for match ${slug}`);
  return res.json();
}

export async function getMatchPerformances(slug: string): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/matches/${slug}/performances`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error(`Failed to fetch performances for match ${slug}`);
  return res.json();
}

export async function getMatchReviewComments(reviewId: string, limit = 20): Promise<Comment[]> {
  const res = await fetch(`${API_URL}/match-reviews/${reviewId}/comments?limit=${limit}`, { cache: 'no-store' });
  if (!res.ok) return [];
  return res.json();
}

// ─── Public Reads — Performances ─────────────────────────────────────────────

export async function getTrendingPerformances(limit = 4): Promise<TrendingResponse<Performance>> {
  const res = await fetch(`${API_URL}/trending/performances?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch trending performances');
  return res.json();
}

export async function getPerformances(limit = 6, offset = 0): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/performances?limit=${limit}&offset=${offset}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch performances');
  return res.json();
}

export async function getTopRatedPerformances(limit = 6): Promise<Performance[]> {
  const res = await fetch(`${API_URL}/performances/top-rated?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch top rated performances');
  return res.json();
}

export async function getPerformance(id: string): Promise<Performance> {
  const res = await fetch(`${API_URL}/performances/${id}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error(`Failed to fetch performance ${id}`);
  return res.json();
}

export async function getPerformanceStats(id: string): Promise<StatRow[]> {
  const res = await fetch(`${API_URL}/performances/${id}/stats`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error(`Failed to fetch stats for performance ${id}`);
  return res.json();
}

export async function getPerformanceReviews(id: string, limit = 20): Promise<Review[]> {
  const res = await fetch(`${API_URL}/performances/${id}/reviews?limit=${limit}`, { next: { revalidate: 30 } });
  if (!res.ok) throw new Error(`Failed to fetch reviews for performance ${id}`);
  return res.json();
}

export async function getPerformanceReviewComments(reviewId: string, limit = 20): Promise<Comment[]> {
  const res = await fetch(`${API_URL}/performance-reviews/${reviewId}/comments?limit=${limit}`, { cache: 'no-store' });
  if (!res.ok) return [];
  return res.json();
}

// ─── Public Reads — Teams & Players ──────────────────────────────────────────

export async function getTeams(limit = 4): Promise<Team[]> {
  const res = await fetch(`${API_URL}/teams?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch teams');
  return res.json();
}

// A team's fixtures (lighter shape than the full Match record).
export interface TeamMatch {
  id: string;
  slug: string;
  title: string;
  utc_datetime: string;
  home_score: number;
  away_score: number;
  home_team_name: string;
  home_team_logo_url?: string;
  home_team_flag_emoji?: string;
  away_team_name: string;
  away_team_logo_url?: string;
  away_team_flag_emoji?: string;
  venue?: string;
}

// A team's player performances (join shape from /teams/{slug}/performances).
export interface TeamPerformance {
  performance_id: string;
  performance_title: string;
  performance_cover_image?: string;
  minutes_played?: number;
  average_rating?: number;
  player_id: string;
  player_name: string;
  player_slug?: string;
  player_photo_url?: string;
  match_id: string;
  match_title: string;
  match_slug?: string;
  match_utc_datetime?: string;
}

export async function getTeamMatches(slug: string, limit = 20, offset = 0): Promise<TeamMatch[]> {
  const res = await fetch(`${API_URL}/teams/${slug}/matches?limit=${limit}&offset=${offset}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) throw new Error(`Failed to fetch matches for team ${slug}`);
  return res.json();
}

export async function getTeamPerformances(slug: string, limit = 24, offset = 0): Promise<TeamPerformance[]> {
  const res = await fetch(`${API_URL}/teams/${slug}/performances?limit=${limit}&offset=${offset}`, {
    next: { revalidate: 60 },
  });
  if (!res.ok) throw new Error(`Failed to fetch performances for team ${slug}`);
  return res.json();
}

export async function getTrendingPlayers(limit = 3): Promise<TrendingResponse<Player>> {
  const res = await fetch(`${API_URL}/trending/players?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch trending players');
  return res.json();
}

export async function getTrendingReviews(limit = 4): Promise<TrendingResponse<Review>> {
  const res = await fetch(`${API_URL}/trending/reviews?limit=${limit}`, { next: { revalidate: 60 } });
  if (!res.ok) throw new Error('Failed to fetch trending reviews');
  return res.json();
}

// ─── Protected — Ratings ─────────────────────────────────────────────────────

export async function submitMatchRating(slug: string, rating: number, token: string): Promise<RatingResponse> {
  const res = await fetch(`${API_URL}/matches/${slug}/rating`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ rating }),
  });
  if (!res.ok) {
    // Try PATCH if already rated
    const patch = await fetch(`${API_URL}/matches/${slug}/rating`, {
      method: 'PATCH',
      headers: authHeaders(token),
      body: JSON.stringify({ rating }),
    });
    if (!patch.ok) {
      const err = await patch.json().catch(() => ({}));
      throw new Error(err.message || 'Failed to submit rating');
    }
    return patch.json();
  }
  return res.json();
}

export async function submitPerformanceRating(id: string, rating: number, token: string): Promise<RatingResponse> {
  const res = await fetch(`${API_URL}/performances/${id}/rating`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ rating }),
  });
  if (!res.ok) {
    const patch = await fetch(`${API_URL}/performances/${id}/rating`, {
      method: 'PATCH',
      headers: authHeaders(token),
      body: JSON.stringify({ rating }),
    });
    if (!patch.ok) {
      const err = await patch.json().catch(() => ({}));
      throw new Error(err.message || 'Failed to submit rating');
    }
    return patch.json();
  }
  return res.json();
}

export async function getMyMatchRating(slug: string, token: string): Promise<UserRating | null> {
  const res = await fetch(`${API_URL}/matches/${slug}/ratings/me`, { headers: authHeaders(token) });
  if (res.status === 404 || !res.ok) return null;
  return res.json();
}

export async function getMyPerformanceRating(id: string, token: string): Promise<UserRating | null> {
  const res = await fetch(`${API_URL}/performances/${id}/ratings/me`, { headers: authHeaders(token) });
  if (res.status === 404 || !res.ok) return null;
  return res.json();
}

// ─── Protected — Reviews ──────────────────────────────────────────────────────

export async function submitMatchReview(
  slug: string,
  content: string,
  token: string,
  title?: string,
): Promise<Review> {
  const res = await fetch(`${API_URL}/matches/${slug}/reviews`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ content, ...(title ? { title } : {}) }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to submit review');
  }
  return res.json();
}

export async function submitPerformanceReview(
  id: string,
  content: string,
  token: string,
  title?: string,
): Promise<Review> {
  const res = await fetch(`${API_URL}/performances/${id}/reviews`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ content, ...(title ? { title } : {}) }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to submit review');
  }
  return res.json();
}

// ─── Protected — Likes ────────────────────────────────────────────────────────

export async function likeMatchReview(reviewId: string, token: string): Promise<void> {
  await fetch(`${API_URL}/match-reviews/${reviewId}/like`, {
    method: 'POST',
    headers: authHeaders(token),
  });
}

export async function unlikeMatchReview(reviewId: string, token: string): Promise<void> {
  await fetch(`${API_URL}/match-reviews/${reviewId}/like`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
}

export async function likePerformanceReview(reviewId: string, token: string): Promise<void> {
  await fetch(`${API_URL}/performance-reviews/${reviewId}/like`, {
    method: 'POST',
    headers: authHeaders(token),
  });
}

export async function unlikePerformanceReview(reviewId: string, token: string): Promise<void> {
  await fetch(`${API_URL}/performance-reviews/${reviewId}/like`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
}

// ─── Protected — Comments ─────────────────────────────────────────────────────

export async function commentOnMatchReview(reviewId: string, body: string, token: string): Promise<Comment> {
  const res = await fetch(`${API_URL}/match-reviews/${reviewId}/comments`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ body }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to post comment');
  }
  return res.json();
}

export async function commentOnPerformanceReview(reviewId: string, body: string, token: string): Promise<Comment> {
  const res = await fetch(`${API_URL}/performance-reviews/${reviewId}/comments`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify({ body }),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || 'Failed to post comment');
  }
  return res.json();
}

// ─── Public Reads — Users / Profiles ─────────────────────────────────────────

export interface UserProfile {
  id: string;
  username: string;
  display_name: string;
  email?: string;
  bio?: string;
  avatar_url?: string;
  created_at?: string;
  review_count?: number;
  rating_count?: number;
  likes_received_count?: number;
}

export interface UserReviewRow {
  review_type: 'match' | 'performance';
  id: string;
  entity_id: string;
  title?: string;
  content?: string;
  created_at: string;
  updated_at?: string;
}

export interface UserRatingRow {
  rating_type: 'match' | 'performance';
  id: string;
  entity_id: string;
  rating: number;
  created_at: string;
  updated_at?: string;
}

export async function getUserProfile(username: string): Promise<UserProfile> {
  const res = await fetch(`${API_URL}/users/${encodeURIComponent(username)}`, { next: { revalidate: 30 } });
  if (!res.ok) throw new Error(`Failed to fetch profile for ${username}`);
  return res.json();
}

export async function getUserReviews(username: string, limit = 20): Promise<UserReviewRow[]> {
  const res = await fetch(`${API_URL}/users/${encodeURIComponent(username)}/reviews?limit=${limit}`, {
    next: { revalidate: 30 },
  });
  if (!res.ok) return [];
  return res.json();
}

export async function getUserRatings(username: string, limit = 20): Promise<UserRatingRow[]> {
  const res = await fetch(`${API_URL}/users/${encodeURIComponent(username)}/ratings?limit=${limit}`, {
    next: { revalidate: 30 },
  });
  if (!res.ok) return [];
  return res.json();
}
