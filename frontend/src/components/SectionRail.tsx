import React from 'react';
import Link from 'next/link';
import { ArrowRight } from 'lucide-react';

interface SectionRailProps {
  title: string;
  label?: string;
  viewAllHref?: string;
  children: React.ReactNode;
  /** When false, hides the whole rail (used for empty data). */
  show?: boolean;
}

/**
 * Apple-TV-style horizontal scroll rail: a section header with an optional
 * "view all" link and a snap-scrolling row that wraps gracefully on mobile.
 */
export default function SectionRail({
  title,
  label,
  viewAllHref,
  children,
  show = true,
}: SectionRailProps) {
  if (!show) return null;

  return (
    <section className="container-max container-px overflow-visible py-8">
      <div className="mb-5 flex items-end justify-between gap-4">
        <div>
          {label && <span className="section-label mb-2 block">{label}</span>}
          <h2 className="font-display text-[1.75rem] font-bold text-foreground">{title}</h2>
        </div>
        {viewAllHref && (
          <Link
            href={viewAllHref}
            className="group inline-flex shrink-0 items-center gap-1.5 text-sm font-semibold text-muted-foreground transition-colors hover:text-primary"
          >
            View all
            <ArrowRight size={15} className="transition-transform group-hover:translate-x-0.5" />
          </Link>
        )}
      </div>

      <div className="hide-scrollbar snap-x-rail -mx-6 flex gap-4 overflow-x-auto overflow-y-visible px-6 pb-2 pt-1 lg:-mx-10 lg:px-10">
        {children}
      </div>
    </section>
  );
}
