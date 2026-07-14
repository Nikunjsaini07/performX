import React from 'react';
import { Heart, MessageCircle, Send, Star } from 'lucide-react';
import { Review } from '@/lib/api';

function reviewBody(review: Review) {
  return review.review_text || review.content || review.comment || 'No written note was added with this rating.';
}

function reviewerName(review: Review) {
  return review.author_name || review.user_name || review.username || 'Archive member';
}

export default function ReviewsSection({ title = 'Community Reviews', reviews }: { title?: string; reviews: Review[] }) {
  return (
    <section className="space-y-5">
      <div className="flex items-end justify-between gap-4">
        <div>
          <p className="section-label mb-1">Reviews</p>
          <h2 className="font-display text-3xl font-bold text-foreground">{title}</h2>
        </div>
      </div>

      <div className="archive-card p-4">
        <textarea
          className="w-full min-h-24 resize-none rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
          placeholder="Write a review after signing in..."
        />
        <div className="mt-3 flex justify-end">
          <button className="btn-primary" type="button">
            <Send size={15} />
            Write a Review
          </button>
        </div>
      </div>

      {reviews.length === 0 ? (
        <div className="archive-card p-10 text-center text-muted-foreground">No reviews yet. The first thoughtful review will set the tone here.</div>
      ) : (
        <div className="space-y-4">
          {reviews.map((review, index) => {
            const name = reviewerName(review);
            const initials = name.slice(0, 2).toUpperCase();
            return (
              <article key={review.id || `${name}-${index}`} className="archive-card p-5">
                <div className="flex items-start gap-4">
                  <div className="w-11 h-11 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-xs font-bold text-primary shrink-0">
                    {review.author_avatar ? <img src={review.author_avatar} alt="" className="w-full h-full rounded-full object-cover" /> : initials}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2 mb-2">
                      <h3 className="font-semibold text-foreground">{name}</h3>
                      <span className="inline-flex items-center gap-1 rounded bg-primary/10 px-2 py-0.5 text-xs font-bold text-primary">
                        <Star size={12} className="fill-primary" />
                        {Number(review.rating || 0).toFixed(1)}
                      </span>
                      {review.created_at && <span className="text-xs text-muted-foreground">{new Date(review.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</span>}
                    </div>
                    {review.title && <p className="mb-2 text-sm font-semibold text-foreground">{review.title}</p>}
                    <p className="text-sm leading-relaxed text-muted-foreground">{reviewBody(review)}</p>
                    <div className="mt-4 flex items-center gap-4 text-xs text-muted-foreground">
                      <button className="inline-flex items-center gap-1.5 hover:text-primary transition-colors" type="button">
                        <Heart size={14} />
                        {review.like_count || 0}
                      </button>
                      <button className="inline-flex items-center gap-1.5 hover:text-primary transition-colors" type="button">
                        <MessageCircle size={14} />
                        {review.comment_count || 0} replies
                      </button>
                    </div>
                  </div>
                </div>
              </article>
            );
          })}
        </div>
      )}
    </section>
  );
}
