import React from 'react';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import { Star, PenLine, Heart, CalendarDays } from 'lucide-react';
import { getUserProfile, getUserReviews, getUserRatings } from '@/lib/data';
import { initials, formatDate, timeAgo, toNumber } from '@/lib/format';
import type { UserReviewRow, UserRatingRow } from '@/lib/api';
import RatingBadge from '@/components/RatingBadge';
import BackButton from '@/components/BackButton';

export const revalidate = 30;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ username: string }>;
}): Promise<Metadata> {
  const { username } = await Promise.resolve(params);
  const profile = await getUserProfile(username);
  return { title: profile ? `${profile.display_name} (@${profile.username})` : 'Profile' };
}

function entityHref(type: 'match' | 'performance', entityId: string, matchSlug?: string): string | null {
  if (type === 'performance') return `/performances/${entityId}`;
  if (type === 'match' && matchSlug) return `/matches/${matchSlug}`;
  return null;
}

export default async function ProfilePage({ params }: { params: Promise<{ username: string }> }) {
  const { username } = await Promise.resolve(params);

  const profile = await getUserProfile(username);
  if (!profile) notFound();

  const [reviews, ratings] = await Promise.all([
    getUserReviews(username, 20),
    getUserRatings(username, 20),
  ]);

  const stats = [
    { label: 'Ratings', value: toNumber(profile.rating_count) ?? ratings.length, icon: Star },
    { label: 'Reviews', value: toNumber(profile.review_count) ?? reviews.length, icon: PenLine },
    { label: 'Likes', value: toNumber(profile.likes_received_count) ?? 0, icon: Heart },
  ];

  return (
    <main>
      {/* Header */}
      <section className="mesh-bg grain relative">
        <div className="container-max container-px relative pb-14 pt-24">
          <BackButton className="mb-6" />
          <div className="flex flex-col items-start gap-5 sm:flex-row sm:items-center">
            <span className="flex h-24 w-24 shrink-0 items-center justify-center overflow-hidden rounded-full border border-border bg-surface-2 font-display text-3xl font-bold text-primary">
              {profile.avatar_url ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={profile.avatar_url} alt={profile.display_name} className="h-full w-full object-cover" />
              ) : (
                initials(profile.display_name || profile.username)
              )}
            </span>
            <div className="min-w-0">
              <h1 className="font-display text-[clamp(1.75rem,4vw,2.5rem)] font-bold text-foreground">
                {profile.display_name}
              </h1>
              <p className="text-muted-foreground">@{profile.username}</p>
              {profile.bio && <p className="mt-3 max-w-xl leading-relaxed text-foreground/80">{profile.bio}</p>}
              {profile.created_at && (
                <p className="mt-3 inline-flex items-center gap-1.5 text-sm text-muted-foreground">
                  <CalendarDays size={14} /> Joined {formatDate(profile.created_at)}
                </p>
              )}
            </div>
          </div>

          {/* Stat boxes */}
          <div className="mt-8 grid max-w-md grid-cols-3 gap-3">
            {stats.map((s) => (
              <div key={s.label} className="card-shell flex flex-col items-center gap-1 py-4">
                <s.icon size={16} className="text-primary" />
                <span className="font-display text-2xl font-bold stat-number text-foreground">{s.value}</span>
                <span className="text-xs text-muted-foreground">{s.label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Body */}
      <div className="container-max container-px grid grid-cols-1 gap-10 py-12 lg:grid-cols-2">
        {/* Recent reviews */}
        <section>
          <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Recent Reviews</h2>
          {reviews.length === 0 ? (
            <div className="card-shell px-5 py-10 text-center text-sm text-muted-foreground">
              No reviews yet.
            </div>
          ) : (
            <div className="space-y-3">
              {reviews.map((r: UserReviewRow) => {
                const href = entityHref(r.review_type, r.entity_id, r.match_slug);
                const body = (
                  <article className="card-shell card-lift p-4">
                    <div className="mb-1.5 flex items-center gap-2">
                      <span className="stage-badge">{r.review_type}</span>
                      <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
                    </div>
                    {r.title && <h3 className="font-display text-sm font-bold text-foreground">{r.title}</h3>}
                    {r.content && <p className="mt-1 line-clamp-3 text-sm text-muted-foreground">{r.content}</p>}
                  </article>
                );
                return href ? (
                  <Link key={r.id} href={href} className="block">
                    {body}
                  </Link>
                ) : (
                  <div key={r.id}>{body}</div>
                );
              })}
            </div>
          )}
        </section>

        {/* Recent ratings */}
        <section>
          <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">Recent Ratings</h2>
          {ratings.length === 0 ? (
            <div className="card-shell px-5 py-10 text-center text-sm text-muted-foreground">
              No ratings yet.
            </div>
          ) : (
            <div className="space-y-3">
              {ratings.map((r: UserRatingRow) => {
                const href = entityHref(r.rating_type, r.entity_id, r.match_slug);
                const body = (
                  <article className="card-shell card-lift flex items-center justify-between gap-3 p-4">
                    <div className="flex items-center gap-2">
                      <span className="stage-badge">{r.rating_type}</span>
                      <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
                    </div>
                    <RatingBadge value={r.rating} size="md" />
                  </article>
                );
                return href ? (
                  <Link key={r.id} href={href} className="block">
                    {body}
                  </Link>
                ) : (
                  <div key={r.id}>{body}</div>
                );
              })}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
