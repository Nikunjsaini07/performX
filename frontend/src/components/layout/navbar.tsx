'use client';

import React, { useEffect, useState, useRef } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Menu, X, User as UserIcon, LogOut, ChevronDown } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import { NAV_LINKS, SITE } from '@/lib/site';
import { initials } from '@/lib/format';
import AuthModal from '@/components/features/auth/AuthModal';
import Logo from '@/components/core/Logo';

export default function Navbar() {
  const pathname = usePathname();
  const { user, isAuthenticated, logout, isLoading } = useAuth();

  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [authOpen, setAuthOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Match and Performance detail routes + Home page
  // keep the navbar permanently transparent — a solid/blurred bar on scroll
  // looks jarring over the cinematic backdrop there.
  const pathSegments = pathname.split('/').filter(Boolean);
  const isCinematicRoute = pathname === '/' || (pathSegments.length === 2 && (pathSegments[0] === 'matches' || pathSegments[0] === 'performances') && pathSegments[1] !== '');

  useEffect(() => {
    if (isCinematicRoute) {
      setScrolled(false);
      return;
    }
    const onScroll = () => setScrolled(window.scrollY > 40);
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, [isCinematicRoute]);

  useEffect(() => {
    setMobileOpen(false);
    setMenuOpen(false);
  }, [pathname]);

  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setMenuOpen(false);
    };
    document.addEventListener('mousedown', onClick);
    return () => document.removeEventListener('mousedown', onClick);
  }, []);

  const isActive = (href: string) =>
    href === '/' ? pathname === '/' : pathname.startsWith(href);

  return (
    <>
      <header
        className={`${isCinematicRoute ? 'absolute' : 'fixed'} inset-x-0 top-0 z-[100] transition-all duration-300 ${
          scrolled
            ? 'border-b border-border bg-background/80 backdrop-blur-xl'
            : 'border-b border-transparent bg-transparent'
        }`}
      >
        <nav className="container-max container-px flex h-16 items-center justify-between gap-4">
          {/* Brand */}
          <Link 
            href="/" 
            className={`flex items-center gap-2 font-display text-xl font-bold tracking-tight transition-colors ${isCinematicRoute && !scrolled ? 'text-white hover:text-white/90' : ''}`}
            style={{ transform: 'translateX(20px)' }}
          >
            <span style={isCinematicRoute && !scrolled ? { color: '#fbfcfa' } : undefined}>
              {SITE.name.replace(/X$/, '')}
            </span>
            <span className="-ml-2 text-gradient-lime">X</span>
          </Link>

          {/* Desktop links */}
          <div className="hidden items-center gap-8 md:flex">
            {NAV_LINKS.map((l) => (
              <Link 
                key={l.href} 
                href={l.href} 
                className={`nav-link ${isActive(l.href) ? 'active' : ''} ${isCinematicRoute && !scrolled ? 'nav-link-light' : ''}`}
              >
                {l.label}
              </Link>
            ))}
          </div>

          {/* Auth */}
          <div className="flex items-center gap-3">
            {!isLoading && isAuthenticated && user ? (
              <div className="relative" ref={menuRef}>
                <button
                  onClick={() => setMenuOpen((v) => !v)}
                  className="flex items-center gap-2 rounded-full border border-border bg-surface-2 py-1 pl-1 pr-2.5 transition-colors hover:border-primary/40"
                >
                  <span className="flex h-7 w-7 items-center justify-center overflow-hidden rounded-full bg-surface text-[11px] font-bold text-primary">
                    {user.avatar_url ? (
                      // eslint-disable-next-line @next/next/no-img-element
                      <img src={user.avatar_url} alt={user.display_name} className="h-full w-full object-cover" />
                    ) : (
                      initials(user.display_name || user.username)
                    )}
                  </span>
                  <span className="hidden max-w-[120px] truncate text-sm font-medium sm:block">
                    {user.display_name || user.username}
                  </span>
                  <ChevronDown size={14} className="text-muted-foreground" />
                </button>

                {menuOpen && (
                  <div className="absolute right-0 mt-2 w-52 overflow-hidden rounded-xl border border-border bg-card shadow-2xl">
                    <div className="border-b border-border px-4 py-3">
                      <p className="truncate text-sm font-semibold text-foreground">{user.display_name}</p>
                      <p className="truncate text-xs text-muted-foreground">@{user.username}</p>
                    </div>
                    <Link
                      href={`/u/${user.username}`}
                      className="flex items-center gap-2.5 px-4 py-2.5 text-sm text-foreground transition-colors hover:bg-surface-2"
                    >
                      <UserIcon size={15} /> Profile
                    </Link>
                    <Link
                      href="/settings"
                      className="flex items-center gap-2.5 px-4 py-2.5 text-sm text-foreground transition-colors hover:bg-surface-2"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>
                      Settings
                    </Link>
                    <button
                      onClick={() => { logout(); setMenuOpen(false); }}
                      className="flex w-full items-center gap-2.5 px-4 py-2.5 text-sm text-red-400 transition-colors hover:bg-surface-2"
                    >
                      <LogOut size={15} /> Sign out
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <button onClick={() => setAuthOpen(true)} className="btn-primary hidden sm:inline-flex">
                Sign in
              </button>
            )}

            {/* Mobile toggle */}
            <button
              onClick={() => setMobileOpen((v) => !v)}
              className="btn-ghost !px-2.5 md:hidden"
              aria-label="Toggle menu"
            >
              {mobileOpen ? <X size={18} /> : <Menu size={18} />}
            </button>
          </div>
        </nav>

        {/* Mobile menu */}
        {mobileOpen && (
          <div className="border-t border-border bg-background/95 backdrop-blur-xl md:hidden">
            <div className="container-max container-px flex flex-col gap-1 py-4">
              {NAV_LINKS.map((l) => (
                <Link
                  key={l.href}
                  href={l.href}
                  className={`rounded-lg px-3 py-2.5 text-sm font-medium transition-colors ${
                    isActive(l.href) ? 'bg-surface-2 text-foreground' : 'text-muted-foreground hover:bg-surface-2'
                  }`}
                >
                  {l.label}
                </Link>
              ))}
              {!isAuthenticated && (
                <button
                  onClick={() => { setAuthOpen(true); setMobileOpen(false); }}
                  className="btn-primary mt-2 w-full"
                >
                  Sign in
                </button>
              )}
            </div>
          </div>
        )}
      </header>

      {/* No spacer: the fixed header floats transparently over hero/backdrop
          content at the top of the page and turns solid on scroll. Pages that
          have no hero add their own top padding to clear the navbar. */}

      <AuthModal isOpen={authOpen} onClose={() => setAuthOpen(false)} defaultTab="login" />
    </>
  );
}
