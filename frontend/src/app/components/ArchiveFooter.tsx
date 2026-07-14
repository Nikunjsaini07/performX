import React from 'react';
import Link from 'next/link';
import AppLogo from '@/components/ui/AppLogo';

export default function ArchiveFooter() {
  return (
    <footer className="border-t border-border bg-card/30 py-12">
      <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
        <div className="flex flex-col md:flex-row items-start justify-between gap-8 mb-8">
          <div>
            <div className="flex items-center gap-2.5 mb-3">
              <AppLogo size={28} />
              <span className="font-bold text-base text-foreground">
                Perform<span className="text-gradient-gold">X</span>
              </span>
            </div>
            <p className="text-sm text-muted-foreground max-w-xs">
              The definitive FIFA 2026 performance archive. Every match, every player, every moment.
            </p>
          </div>
          <nav className="flex flex-col gap-2">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1">Archive</p>
            {[
              { label: 'Home', href: '/' },
              { label: 'Matches', href: '/matches' },
              { label: 'Performances', href: '/performances' },
            ]?.map((link) => (
              <Link
                key={`footer-${link?.href}`}
                href={link?.href}
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                {link?.label}
              </Link>
            ))}
          </nav>
        </div>
        <div className="pt-6 border-t border-border flex flex-col sm:flex-row items-center justify-between gap-3">
          <p className="text-xs text-muted-foreground">
            © 2026 PerformX. An independent fan archive. Not affiliated with FIFA.
          </p>
          <p className="text-xs text-muted-foreground">
            96 matches · 421 players · 734 performances
          </p>
        </div>
      </div>
    </footer>
  );
}