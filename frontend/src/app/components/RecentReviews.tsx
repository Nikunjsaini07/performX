import React from 'react';
import { MessageSquare, Quote } from 'lucide-react';
import { getTrendingReviews } from '@/lib/api';

function ScoreBar({ score }: { score: number }) {
  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-1 rounded-full bg-muted overflow-hidden">
        <div
          className="h-full rounded-full bg-gradient-to-r from-primary to-accent transition-all duration-500"
          style={{ width: `${score * 10}%` }}
        />
      </div>
      <span className="stat-number text-xs font-bold text-primary w-8 text-right">{score}/10</span>
    </div>
  );
}

function timeSince(dateString: string) {
  const date = new Date(dateString);
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);
  
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
  return `${Math.floor(seconds / 86400)} days ago`;
}

export default async function RecentReviews() {
  let reviews: any[] = [];
  try {
    const res = await getTrendingReviews(6);
    if (res && res.data) {
      reviews = res.data;
    }
  } catch (error) {
    console.error("Failed to fetch recent reviews:", error);
  }

  return (
    <section className="py-24 bg-background relative">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[60%] h-px bg-gradient-to-r from-transparent via-border to-transparent" />

      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        <div className="flex items-end justify-between mb-12">
          <div>
            <p className="section-label mb-2">Community</p>
            <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground">
              Trending Reviews
            </h2>
            <p className="text-sm text-muted-foreground mt-2">
              Most discussed community reviews from the archive
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
          {reviews.length === 0 ? (
            <div className="col-span-full py-20 text-center text-muted-foreground">
              No recent reviews found. Ensure the API is running.
            </div>
          ) : (
            reviews.map((item) => {
              const review = item.entity;
              const reviewerName = review.author_name || 'Anonymous User';
              const avatar = reviewerName.substring(0, 2).toUpperCase();
              
              // Handle differences between match reviews and performance reviews
              const title = review.player_name ? review.player_name : (review.match_title || 'Review');
              const subtitle = review.player_name ? (review.match_title || 'Performance') : 'Match Review';
              const flag = '💬';

              return (
                <div
                  key={review.id || item.entity_id}
                  className="relative overflow-hidden rounded-2xl border border-border bg-card p-5 flex flex-col gap-4 transition-all duration-300 hover:border-primary/25 hover:-translate-y-0.5 hover:shadow-[0_16px_48px_rgba(0,0,0,0.5)]"
                >
                  {/* Quote icon */}
                  <div className="absolute top-4 right-4 opacity-10">
                    <Quote size={32} className="text-primary" />
                  </div>

                  {/* Context */}
                  <div className="flex items-start gap-3">
                    <div className="w-11 h-11 rounded-full bg-muted border border-border flex items-center justify-center text-2xl flex-shrink-0">
                      {flag}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-bold text-foreground truncate">{title}</p>
                      <p className="text-xs text-muted-foreground truncate">{subtitle}</p>
                    </div>
                    <div className="performance-score flex-shrink-0 w-9 h-9 text-xs">
                      {review.rating || 0}
                    </div>
                  </div>

                  {/* Score bar */}
                  <ScoreBar score={review.rating || 0} />

                  {/* Review text */}
                  <p className="text-sm text-muted-foreground leading-relaxed line-clamp-3 flex-1">
                    &ldquo;{review.review_text || 'No comments provided for this rating.'}&rdquo;
                  </p>

                  {/* Reviewer */}
                  <div className="flex items-center justify-between pt-3 border-t border-border">
                    <div className="flex items-center gap-2">
                      <div className="w-7 h-7 rounded-full bg-primary/15 border border-primary/25 flex items-center justify-center">
                        <span className="text-xs font-bold text-primary">{avatar}</span>
                      </div>
                      <span className="text-xs font-semibold text-muted-foreground">{reviewerName}</span>
                    </div>
                    <span className="text-xs text-muted-foreground">{review.created_at ? timeSince(review.created_at) : ''}</span>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* CTA to join */}
        <div className="mt-14 relative overflow-hidden flex flex-col items-center text-center gap-5 py-14 px-6 rounded-2xl border border-border bg-card/50">
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-accent/5 pointer-events-none" />
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[40%] h-px bg-gradient-to-r from-transparent via-primary/40 to-transparent" />

          <div className="relative z-10 w-14 h-14 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center">
            <MessageSquare size={24} className="text-primary" />
          </div>
          <div className="relative z-10">
            <h3 className="font-display font-bold text-foreground text-2xl mb-2">
              Rate performances you watched
            </h3>
            <p className="text-sm text-muted-foreground max-w-md leading-relaxed">
              Create a free account to rate and review any of the archived performances.
              Your ratings contribute to the community score.
            </p>
          </div>
          <button className="relative z-10 btn-primary px-8 py-3 text-sm">
            Join the Archive
          </button>
        </div>
      </div>
    </section>
  );
}