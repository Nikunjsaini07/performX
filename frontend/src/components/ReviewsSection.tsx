'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Heart, MessageCircle, Loader2, Send, PenLine } from 'lucide-react';
import { useAuth } from '@/lib/auth-context';
import {
  submitMatchReview,
  submitPerformanceReview,
  likeMatchReview,
  unlikeMatchReview,
  likePerformanceReview,
  unlikePerformanceReview,
  commentOnMatchReview,
  commentOnPerformanceReview,
  getMatchReviewComments,
  getPerformanceReviewComments,
  type Review,
  type Comment,
} from '@/lib/api';
import { reviewText, reviewAuthorName, timeAgo, initials, toNumber } from '@/lib/format';
import AuthModal from '@/app/components/AuthModal';

interface ReviewsSectionProps {
  kind: 'match' | 'performance';
  entityKey: string;
  initialReviews: Review[];
}

interface LocalReview extends Review {
  _liked?: boolean;
  _likes?: number;
  _comments?: Comment[];
  _commentCount?: number;
  _commentsLoaded?: boolean;
  _commentsOpen?: boolean;
}

function Avatar({ name, url, size = 36 }: { name: string; url?: string; size?: number }) {
  return (
    <span
      style={{ width: size, height: size }}
      className="flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-surface-2 text-xs font-bold text-muted-foreground"
    >
      {url ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={url} alt={name} className="h-full w-full object-cover" loading="lazy" />
      ) : (
        initials(name)
      )}
    </span>
  );
}

export default function ReviewsSection({ kind, entityKey, initialReviews }: ReviewsSectionProps) {
  const { token, isAuthenticated, user } = useAuth();

  const [reviews, setReviews] = useState<LocalReview[]>(() =>
    initialReviews.map((r) => ({
      ...r,
      _liked: !!r.liked_by_me,
      _likes: toNumber(r.likes_count) ?? toNumber(r.like_count) ?? 0,
      _commentCount: toNumber(r.comment_count) ?? 0,
      // Normalize author fields for consistent access
      username: r.author_username || r.username,
      display_name: r.author_display_name || r.display_name,
      avatar_url: r.author_avatar_url || r.avatar_url,
    })),
  );

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [posting, setPosting] = useState(false);
  const [error, setError] = useState('');
  const [authOpen, setAuthOpen] = useState(false);
  const [commentDrafts, setCommentDrafts] = useState<Record<string, string>>({});
  const [commentBusy, setCommentBusy] = useState<Record<string, boolean>>({});
  const [visibleCount, setVisibleCount] = useState(5);

  const REVIEWS_PER_PAGE = 5;
  const visibleReviews = reviews.slice(0, visibleCount);
  const hasMore = visibleCount < reviews.length;

  const update = (id: string, patch: Partial<LocalReview>) =>
    setReviews((rs) => rs.map((r) => (r.id === id ? { ...r, ...patch } : r)));

  const handlePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isAuthenticated || !token) {
      setAuthOpen(true);
      return;
    }
    if (!content.trim()) return;
    setPosting(true);
    setError('');
    try {
      const created =
        kind === 'match'
          ? await submitMatchReview(entityKey, content.trim(), token, title.trim() || undefined)
          : await submitPerformanceReview(entityKey, content.trim(), token, title.trim() || undefined);

      const local: LocalReview = {
        ...created,
        content: created.content || content.trim(),
        title: created.title || title.trim() || undefined,
        created_at: created.created_at || new Date().toISOString(),
        display_name: created.display_name || user?.display_name,
        username: created.username || user?.username,
        avatar_url: created.avatar_url || user?.avatar_url,
        _liked: false,
        _likes: 0,
        _commentCount: 0,
      };
      setReviews((rs) => [local, ...rs]);
      setVisibleCount((c) => c + 1); // Show the new review immediately
      setContent('');
      setTitle('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not post review');
    } finally {
      setPosting(false);
    }
  };

  const toggleLike = async (r: LocalReview) => {
    if (!isAuthenticated || !token) {
      setAuthOpen(true);
      return;
    }
    const nextLiked = !r._liked;
    update(r.id, { _liked: nextLiked, _likes: Math.max(0, (r._likes ?? 0) + (nextLiked ? 1 : -1)) });
    try {
      if (kind === 'match') {
        nextLiked ? await likeMatchReview(r.id, token) : await unlikeMatchReview(r.id, token);
      } else {
        nextLiked ? await likePerformanceReview(r.id, token) : await unlikePerformanceReview(r.id, token);
      }
    } catch {
      // revert on failure
      update(r.id, { _liked: r._liked, _likes: r._likes });
    }
  };

  const toggleComments = async (r: LocalReview) => {
    const open = !r._commentsOpen;
    update(r.id, { _commentsOpen: open });
    if (open && !r._commentsLoaded) {
      try {
        const comments =
          kind === 'match'
            ? await getMatchReviewComments(r.id)
            : await getPerformanceReviewComments(r.id);
        update(r.id, { _comments: comments, _commentsLoaded: true, _commentCount: comments.length });
      } catch {
        update(r.id, { _comments: [], _commentsLoaded: true });
      }
    }
  };

  const postComment = async (r: LocalReview) => {
    if (!isAuthenticated || !token) {
      setAuthOpen(true);
      return;
    }
    const body = (commentDrafts[r.id] || '').trim();
    if (!body) return;
    setCommentBusy((b) => ({ ...b, [r.id]: true }));
    try {
      const created =
        kind === 'match'
          ? await commentOnMatchReview(r.id, body, token)
          : await commentOnPerformanceReview(r.id, body, token);
      const local: Comment = {
        ...created,
        body: created.body || body,
        created_at: created.created_at || new Date().toISOString(),
        display_name: created.display_name || user?.display_name,
        username: created.username || user?.username,
        avatar_url: created.avatar_url || user?.avatar_url,
      };
      update(r.id, {
        _comments: [...(r._comments ?? []), local],
        _commentCount: (r._commentCount ?? 0) + 1,
        _commentsLoaded: true,
        _commentsOpen: true,
      });
      setCommentDrafts((d) => ({ ...d, [r.id]: '' }));
    } catch {
      /* ignore */
    } finally {
      setCommentBusy((b) => ({ ...b, [r.id]: false }));
    }
  };

  return (
    <div>
      <h2 className="mb-5 font-display text-[1.75rem] font-bold text-foreground">
        Reviews <span className="text-muted-foreground">({reviews.length})</span>
      </h2>

      {/* Write box */}
      {isAuthenticated ? (
        <form onSubmit={handlePost} className="card-shell mb-8 p-5">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Review title (optional)"
            className="input-field mb-3"
            maxLength={120}
          />
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder={`Share your take on this ${kind}…`}
            rows={3}
            className="input-field resize-y"
          />
          {error && <p className="mt-2 text-sm text-red-400">{error}</p>}
          <div className="mt-3 flex justify-end">
            <button type="submit" disabled={posting || !content.trim()} className="btn-primary disabled:opacity-50">
              {posting ? <Loader2 size={15} className="animate-spin" /> : <Send size={15} />}
              Post review
            </button>
          </div>
        </form>
      ) : (
        <button
          onClick={() => setAuthOpen(true)}
          className="card-shell mb-8 flex w-full items-center gap-3 p-5 text-left transition-colors hover:border-primary/40"
        >
          <PenLine size={18} className="text-primary" />
          <span className="text-sm text-muted-foreground">
            <span className="font-semibold text-foreground">Sign in</span> to write a review, like and comment.
          </span>
        </button>
      )}

      {/* List */}
      {reviews.length === 0 ? (
        <div className="card-shell flex flex-col items-center gap-2 px-6 py-12 text-center">
          <MessageCircle size={26} className="text-muted-foreground/50" />
          <p className="font-medium text-foreground">No reviews yet</p>
          <p className="text-sm text-muted-foreground">Be the first to share what you thought.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {visibleReviews.map((r) => {
            const name = reviewAuthorName(r);
            const text = reviewText(r);
            return (
              <article key={r.id} className="card-shell p-5 transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg">
                <div className="flex items-center gap-3">
                  {r.username ? (
                    <Link href={`/u/${r.username}`} className="group flex items-center gap-3">
                      <Avatar name={name} url={r.avatar_url} />
                      <span className="flex flex-col">
                        <span className="text-sm font-semibold text-foreground transition-colors group-hover:text-primary">
                          {name}
                        </span>
                        <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
                      </span>
                    </Link>
                  ) : (
                    <span className="flex items-center gap-3">
                      <Avatar name={name} url={r.avatar_url} />
                      <span className="flex flex-col">
                        <span className="text-sm font-semibold text-foreground">{name}</span>
                        <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
                      </span>
                    </span>
                  )}
                </div>

                {r.title && <h3 className="mt-3 font-display text-base font-bold text-foreground">{r.title}</h3>}
                {text && <p className="mt-2 whitespace-pre-line text-sm leading-relaxed text-muted-foreground">{text}</p>}

                {/* Actions */}
                <div className="mt-4 flex items-center gap-4 text-sm">
                  <button
                    onClick={() => toggleLike(r)}
                    className={`inline-flex items-center gap-1.5 transition-colors ${
                      r._liked ? 'text-primary' : 'text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    <Heart size={15} className={r._liked ? 'fill-primary' : ''} /> {r._likes ?? 0}
                  </button>
                  <button
                    onClick={() => toggleComments(r)}
                    className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground"
                  >
                    <MessageCircle size={15} /> {r._commentCount ?? 0}
                  </button>
                </div>

                {/* Comments */}
                {r._commentsOpen && (
                  <div className="mt-4 space-y-3 border-t border-border/70 pt-4">
                    {r._commentsLoaded === false || (!r._commentsLoaded && !r._comments) ? (
                      <p className="text-xs text-muted-foreground">
                        <Loader2 size={13} className="mr-1 inline animate-spin" /> Loading comments…
                      </p>
                    ) : (r._comments ?? []).length === 0 ? (
                      <p className="text-xs text-muted-foreground">No comments yet.</p>
                    ) : (
                      (r._comments ?? []).map((c) => {
                        const cName = c.display_name || c.username || 'Anonymous';
                        return (
                          <div key={c.id} className="flex gap-2.5">
                            <Avatar name={cName} url={c.avatar_url} size={28} />
                            <div className="min-w-0">
                              <div className="flex items-baseline gap-2">
                                {c.username ? (
                                  <Link href={`/u/${c.username}`} className="text-xs font-semibold text-foreground hover:text-primary">
                                    {cName}
                                  </Link>
                                ) : (
                                  <span className="text-xs font-semibold text-foreground">{cName}</span>
                                )}
                                <span className="text-[10px] text-muted-foreground">{timeAgo(c.created_at)}</span>
                              </div>
                              <p className="text-sm text-muted-foreground">{c.body}</p>
                            </div>
                          </div>
                        );
                      })
                    )}

                    {/* Add comment */}
                    <div className="flex gap-2 pt-1">
                      <input
                        value={commentDrafts[r.id] || ''}
                        onChange={(e) => setCommentDrafts((d) => ({ ...d, [r.id]: e.target.value }))}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') { e.preventDefault(); postComment(r); }
                        }}
                        placeholder={isAuthenticated ? 'Add a comment…' : 'Sign in to comment'}
                        className="input-field !py-2 text-sm"
                      />
                      <button
                        onClick={() => postComment(r)}
                        disabled={commentBusy[r.id]}
                        className="btn-ghost !px-3 disabled:opacity-50"
                        aria-label="Post comment"
                      >
                        {commentBusy[r.id] ? <Loader2 size={15} className="animate-spin" /> : <Send size={15} />}
                      </button>
                    </div>
                  </div>
                )}
              </article>
            );
          })}
        </div>
      )}

      {/* Load More Button */}
      {hasMore && (
        <div className="mt-6 flex justify-center">
          <button
            onClick={() => setVisibleCount((c) => c + REVIEWS_PER_PAGE)}
            className="btn-ghost px-6 py-2.5 text-sm font-medium"
          >
            Load More ({reviews.length - visibleCount} remaining)
          </button>
        </div>
      )}

      <AuthModal isOpen={authOpen} onClose={() => setAuthOpen(false)} defaultTab="login" />
    </div>
  );
}
