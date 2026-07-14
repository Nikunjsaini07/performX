'use client';

import React, { useState } from 'react';
import { Star } from 'lucide-react';

export default function RatingPanel({ label = 'Your Rating' }: { label?: string }) {
  const [rating, setRating] = useState<number | null>(null);

  return (
    <div className="archive-card p-5">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <p className="section-label mb-1">{label}</p>
          <h3 className="font-display text-2xl font-bold text-foreground">Log your score</h3>
          <p className="text-sm text-muted-foreground mt-1">Sign in to save a 1-10 rating and write a review.</p>
        </div>
        <div className="flex flex-wrap items-center gap-1.5">
          {Array.from({ length: 10 }, (_, i) => i + 1).map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => setRating(value)}
              className={`w-9 h-9 rounded-md border text-sm font-bold transition-all focus:outline-none focus:ring-2 focus:ring-primary/60 ${
                rating === value
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-muted text-muted-foreground hover:border-primary/50 hover:text-foreground'
              }`}
              aria-label={`Rate ${value} out of 10`}
            >
              {value}
            </button>
          ))}
        </div>
      </div>
      {rating && (
        <div className="mt-4 flex items-center gap-2 text-sm text-primary">
          <Star size={15} className="fill-primary" />
          Draft rating selected: {rating}/10
        </div>
      )}
    </div>
  );
}
