import React from 'react';
import Link from 'next/link';
import { NAV_LINKS, SITE } from '@/lib/site';

export default function Footer() {
  return (
    <footer className="relative mt-8 border-t border-border">
      <div className="glow-lime pointer-events-none absolute inset-x-0 top-0 h-40" />
      <div className="container-max container-px relative py-14">
        <div className="flex flex-col justify-between gap-10 md:flex-row">
          <div className="max-w-sm">
            <Link href="/" className="font-display text-2xl font-bold tracking-tight">
              {SITE.name.replace(/X$/, '')}
              <span className="text-gradient-lime">X</span>
            </Link>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{SITE.description}</p>
          </div>

          <div className="flex gap-16">
            <div>
              <h3 className="mb-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
                Explore
              </h3>
              <ul className="space-y-2.5">
                {NAV_LINKS.map((l) => (
                  <li key={l.href}>
                    <Link href={l.href} className="text-sm text-foreground/80 transition-colors hover:text-primary">
                      {l.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="mb-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">
                Tournament
              </h3>
              <ul className="space-y-2.5 text-sm text-foreground/80">
                <li>FIFA World Cup 2026</li>
                <li>USA · Canada · Mexico</li>
                <li>48 Teams</li>
              </ul>
            </div>
          </div>
        </div>

        <div className="mt-12 flex flex-col items-center justify-between gap-3 border-t border-border pt-6 text-xs text-muted-foreground sm:flex-row">
          <p>© {new Date().getFullYear()} {SITE.name}. Community ratings for the beautiful game.</p>
          <p>{SITE.tagline}</p>
        </div>
      </div>
    </footer>
  );
}
