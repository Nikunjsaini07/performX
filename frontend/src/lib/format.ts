// Shared formatting helpers used across the UI.

import type { Match } from '@/lib/api';

// Ratings and numeric fields can arrive as number, numeric-string, or pgx
// wrapper objects ({ Int32, Valid } / { Float64, Valid }). Normalize to number.
export function toNumber(value: unknown): number | null {
  if (value == null) return null;
  if (typeof value === 'number') return Number.isFinite(value) ? value : null;
  if (typeof value === 'string') {
    const n = parseFloat(value);
    return Number.isFinite(n) ? n : null;
  }
  if (typeof value === 'object') {
    const o = value as Record<string, unknown>;
    if ('Valid' in o && o.Valid === false) return null;
    if ('Int32' in o) return toNumber(o.Int32);
    if ('Float64' in o) return toNumber(o.Float64);
    if ('Int64' in o) return toNumber(o.Int64);
  }
  return null;
}

export function formatRating(value: unknown): string {
  const n = toNumber(value);
  if (n == null) return '—';
  return n.toFixed(1);
}

// Rating tier for colour coding (matches .performance-score modifiers).
export function ratingTier(value: unknown): 'high' | 'mid' | 'low' | 'none' {
  const n = toNumber(value);
  if (n == null) return 'none';
  if (n >= 8) return 'high';
  if (n >= 6.5) return 'mid';
  return 'low';
}

export function penaltyValue(value?: Match['home_penalty_score']): number | null {
  return toNumber(value);
}

export function formatDate(iso?: string): string {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}

export function formatDateTime(iso?: string): string {
  if (!iso) return '';
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  });
}

export function timeAgo(iso?: string): string {
  if (!iso) return '';
  const d = new Date(iso).getTime();
  if (isNaN(d)) return '';
  const secs = Math.floor((Date.now() - d) / 1000);
  const units: [number, string][] = [
    [31536000, 'y'],
    [2592000, 'mo'],
    [604800, 'w'],
    [86400, 'd'],
    [3600, 'h'],
    [60, 'm'],
  ];
  for (const [s, label] of units) {
    const v = Math.floor(secs / s);
    if (v >= 1) return `${v}${label} ago`;
  }
  return 'just now';
}

export function initials(name?: string): string {
  if (!name) return '?';
  return name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((p) => p[0]?.toUpperCase() ?? '')
    .join('');
}

// Reviews arrive with several possible text/author field names — normalize.
export function reviewText(r: {
  content?: string;
  comment?: string;
  review_text?: string;
  body?: string;
}): string {
  return r.content || r.comment || r.review_text || r.body || '';
}

export function reviewAuthorName(r: {
  author_display_name?: string;
  author_username?: string;
  display_name?: string;
  author_name?: string;
  user_name?: string;
  username?: string;
}): string {
  return r.author_display_name || r.display_name || r.author_name || r.author_username || r.user_name || r.username || 'Anonymous';
}
