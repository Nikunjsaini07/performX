import React from 'react';
import Link from 'next/link';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  /** Base path, e.g. "/matches" — page is appended as ?page=N. */
  basePath: string;
}

// Build a compact list of page tokens with ellipsis, e.g. [1, '…', 4, 5, 6, '…', 20].
function buildPages(current: number, total: number): (number | 'ellipsis')[] {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const pages: (number | 'ellipsis')[] = [1];
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);

  if (start > 2) pages.push('ellipsis');
  for (let p = start; p <= end; p++) pages.push(p);
  if (end < total - 1) pages.push('ellipsis');

  pages.push(total);
  return pages;
}

/**
 * Numbered pagination: ‹ Previous  1 2 3 … N  Next ›
 * Current page is highlighted in the lime accent; Previous/Next disable at ends.
 */
export default function Pagination({ currentPage, totalPages, basePath }: PaginationProps) {
  if (totalPages <= 1) return null;

  const page = Math.min(Math.max(1, currentPage), totalPages);
  const hasPrev = page > 1;
  const hasNext = page < totalPages;
  const href = (p: number) => (p <= 1 ? basePath : `${basePath}?page=${p}`);
  const pages = buildPages(page, totalPages);

  return (
    <nav className="mt-10 flex flex-wrap items-center justify-center gap-1.5" aria-label="Pagination">
      {hasPrev ? (
        <Link href={href(page - 1)} className="btn-ghost !px-3" aria-label="Previous page">
          <ChevronLeft size={16} /> Previous
        </Link>
      ) : (
        <span className="btn-ghost !px-3 pointer-events-none opacity-40" aria-disabled="true">
          <ChevronLeft size={16} /> Previous
        </span>
      )}

      <div className="flex items-center gap-1">
        {pages.map((p, i) =>
          p === 'ellipsis' ? (
            <span key={`e${i}`} className="px-2 text-sm text-muted-foreground/60 select-none">
              …
            </span>
          ) : p === page ? (
            <span
              key={p}
              aria-current="page"
              className="inline-flex h-9 min-w-9 items-center justify-center rounded-full bg-primary px-3 text-sm font-bold text-primary-foreground"
            >
              {p}
            </span>
          ) : (
            <Link
              key={p}
              href={href(p)}
              className="inline-flex h-9 min-w-9 items-center justify-center rounded-full border border-border px-3 text-sm font-medium text-muted-foreground transition-colors hover:border-primary/40 hover:text-foreground"
            >
              {p}
            </Link>
          ),
        )}
      </div>

      {hasNext ? (
        <Link href={href(page + 1)} className="btn-ghost !px-3" aria-label="Next page">
          Next <ChevronRight size={16} />
        </Link>
      ) : (
        <span className="btn-ghost !px-3 pointer-events-none opacity-40" aria-disabled="true">
          Next <ChevronRight size={16} />
        </span>
      )}
    </nav>
  );
}
