'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Heart, MessageCircle, Quote, Trash2, Edit2, X, Check, Loader2 } from 'lucide-react';
import type { Review } from '@/lib/api';
import { reviewText, reviewAuthorName, timeAgo, initials, toNumber } from '@/lib/format';
import { useAuth } from '@/contexts/AuthContext';
import { deleteMatchReview, deletePerformanceReview, updateMatchReview, updatePerformanceReview } from '@/lib/api';

interface ReviewCardProps {
  review: Review & {
    author_username?: string;
    author_display_name?: string;
    author_avatar_url?: string;
    likes_count?: number;
    match_slug?: string;
    performance_id?: string;
  };
  className?: string;
  onDelete?: () => void;
}

export default function ReviewCard({ review: r, className = '', onDelete }: ReviewCardProps) {
  const { user, token } = useAuth();
  const [editing, setEditing] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [editTitle, setEditTitle] = useState(r.title || '');
  const [editContent, setEditContent] = useState(reviewText(r));
  const [loading, setLoading] = useState(false);

  console.log('ReviewCard data:', {
    author_display_name: r.author_display_name,
    author_username: r.author_username,
    author_avatar_url: r.author_avatar_url,
    likes_count: r.likes_count,
    comment_count: r.comment_count,
    display_name: r.display_name,
    username: r.username,
    avatar_url: r.avatar_url,
    like_count: r.like_count
  });

  const text = reviewText(r);
  const name = r.author_display_name || reviewAuthorName(r);
  const username = r.author_username || r.username;
  const avatar = r.author_avatar_url || r.avatar_url;
  const likes = toNumber(r.likes_count) ?? toNumber(r.like_count) ?? 0;
  const comments = toNumber(r.comment_count) ?? 0;

  const isOwner = user && (r.user_id === user.id);

  const context = r.match_title || r.player_name;
  const contextHref = r.match_slug
    ? `/matches/${r.match_slug}`
    : r.performance_id
      ? `/performances/${r.performance_id}`
      : undefined;

  const handleDelete = async () => {
    if (!token || !confirm('Are you sure you want to delete this review?')) return;
    
    setLoading(true);
    try {
      if (r.match_slug) {
        await deleteMatchReview(token, r.match_slug, r.id);
      } else if (r.performance_id) {
        await deletePerformanceReview(token, r.performance_id, r.id);
      }
      onDelete?.();
      window.location.reload();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete review');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveEdit = async () => {
    if (!token) return;
    
    setLoading(true);
    try {
      if (r.match_slug) {
        await updateMatchReview(token, r.match_slug, r.id, { title: editTitle, content: editContent });
      } else if (r.performance_id) {
        await updatePerformanceReview(token, r.performance_id, r.id, { title: editTitle, content: editContent });
      }
      setEditing(false);
      window.location.reload();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update review');
    } finally {
      setLoading(false);
    }
  };

  const Avatar = (
    <span className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-surface-2 text-xs font-bold text-muted-foreground">
      {avatar ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={avatar} alt={name} className="h-full w-full object-cover" loading="lazy" />
      ) : (
        initials(name)
      )}
    </span>
  );

  if (editing) {
    return (
      <article className={`card-shell flex flex-col p-5 ${className}`}>
        <div className="space-y-3">
          <input
            type="text"
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            placeholder="Review title (optional)"
            className="w-full rounded-lg border border-border bg-muted px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
          />
          <textarea
            value={editContent}
            onChange={(e) => setEditContent(e.target.value)}
            rows={4}
            placeholder="Your review..."
            className="w-full rounded-lg border border-border bg-muted px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 resize-none"
          />
          <div className="flex gap-2">
            <button onClick={handleSaveEdit} disabled={loading} className="btn-primary text-xs !py-1.5 !px-3">
              {loading ? <Loader2 size={12} className="animate-spin" /> : <><Check size={12} /> Save</>}
            </button>
            <button onClick={() => setEditing(false)} className="btn-ghost text-xs !py-1.5 !px-3">
              <X size={12} /> Cancel
            </button>
          </div>
        </div>
      </article>
    );
  }

  return (
    <article className={`card-shell flex flex-col p-5 transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-lg ${className}`}>
      <div className="flex items-start justify-between gap-2 mb-3">
        <Quote size={20} className="text-primary/40" />
        {isOwner && !deleting && (
          <div className="flex gap-1">
            <button onClick={() => setEditing(true)} className="p-1.5 rounded-lg text-muted-foreground hover:text-primary hover:bg-surface-2 transition-colors" title="Edit review">
              <Edit2 size={14} />
            </button>
            <button onClick={() => setDeleting(true)} className="p-1.5 rounded-lg text-muted-foreground hover:text-red-400 hover:bg-surface-2 transition-colors" title="Delete review">
              <Trash2 size={14} />
            </button>
          </div>
        )}
        {deleting && (
          <div className="flex gap-1">
            <button onClick={handleDelete} disabled={loading} className="p-1.5 rounded-lg text-red-400 hover:bg-red-500/10 transition-colors" title="Confirm delete">
              {loading ? <Loader2 size={14} className="animate-spin" /> : <Check size={14} />}
            </button>
            <button onClick={() => setDeleting(false)} className="p-1.5 rounded-lg text-muted-foreground hover:bg-surface-2 transition-colors" title="Cancel">
              <X size={14} />
            </button>
          </div>
        )}
      </div>

      {r.title && (
        <h3 className="mb-1.5 line-clamp-1 font-display text-base font-bold text-foreground">{r.title}</h3>
      )}

      <p className="line-clamp-4 flex-1 text-sm leading-relaxed text-muted-foreground">
        {text || 'No written review.'}
      </p>

      {context && (
        contextHref ? (
          <Link
            href={contextHref}
            className="mt-3 line-clamp-1 text-xs font-medium text-primary/80 transition-colors hover:text-primary"
          >
            {context}
          </Link>
        ) : (
          <span className="mt-3 line-clamp-1 text-xs font-medium text-muted-foreground">{context}</span>
        )
      )}

      {/* Footer */}
      <div className="mt-4 flex items-center justify-between gap-3 border-t border-border/70 pt-3">
        {username ? (
          <Link href={`/u/${username}`} className="group flex min-w-0 items-center gap-2">
            {Avatar}
            <span className="flex min-w-0 flex-col">
              <span className="line-clamp-1 text-xs font-semibold text-foreground transition-colors group-hover:text-primary">
                {name}
              </span>
              <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
            </span>
          </Link>
        ) : (
          <span className="flex min-w-0 items-center gap-2">
            {Avatar}
            <span className="flex min-w-0 flex-col">
              <span className="line-clamp-1 text-xs font-semibold text-foreground">{name}</span>
              <span className="text-[11px] text-muted-foreground">{timeAgo(r.created_at)}</span>
            </span>
          </span>
        )}

        <div className="flex shrink-0 items-center gap-3 text-xs text-muted-foreground">
          <span className="inline-flex items-center gap-1">
            <Heart size={13} /> {likes}
          </span>
          <span className="inline-flex items-center gap-1">
            <MessageCircle size={13} /> {comments}
          </span>
        </div>
      </div>
    </article>
  );
}
