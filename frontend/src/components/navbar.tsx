'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import AppLogo from '@/components/ui/AppLogo';
import { Menu, X, User, LogIn } from 'lucide-react';

const navLinks = [
  { label: 'Home', href: '/' },
  { label: 'Matches', href: '/matches' },
  { label: 'Performances', href: '/performances' },
  { label: 'Teams', href: '/teams' },
];

export default function Navbar() {
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 40);
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <>
      <header
        className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
          scrolled
            ? 'bg-background/95 backdrop-blur-md border-b border-border' :'bg-transparent'
        }`}
      >
        <div className="max-w-screen-2xl mx-auto px-6 lg:px-10">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link href="/" className="flex items-center gap-2.5 group">
              <AppLogo size={32} />
              <span className="font-bold text-lg tracking-tight text-foreground">
                Perform<span className="text-gradient-gold">X</span>
              </span>
            </Link>

            {/* Desktop Nav */}
            <nav className="hidden md:flex items-center gap-8">
              {navLinks?.map((link) => (
                <Link
                  key={`nav-${link?.href}`}
                  href={link?.href}
                  className={`nav-link text-sm font-medium transition-colors duration-150 ${
                    pathname === link?.href
                      ? 'text-foreground'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {link?.label}
                </Link>
              ))}
            </nav>

            {/* Desktop Actions */}
            <div className="hidden md:flex items-center gap-3">
              <button className="btn-ghost text-sm flex items-center gap-2">
                <LogIn size={15} />
                Sign In
              </button>
              <button className="btn-primary text-sm flex items-center gap-2">
                <User size={15} />
                Join Archive
              </button>
            </div>

            {/* Mobile Hamburger */}
            <button
              className="md:hidden p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
              onClick={() => setMobileOpen(true)}
              aria-label="Open menu"
            >
              <Menu size={22} />
            </button>
          </div>
        </div>
      </header>
      {/* Mobile Drawer */}
      {mobileOpen && (
        <div className="fixed inset-0 z-[100] md:hidden">
          <div
            className="absolute inset-0 bg-background/80 backdrop-blur-sm"
            onClick={() => setMobileOpen(false)}
          />
          <div className="absolute top-0 right-0 bottom-0 w-72 bg-card border-l border-border flex flex-col">
            <div className="flex items-center justify-between px-5 h-16 border-b border-border">
              <span className="font-bold text-base">
                Perform<span className="text-gradient-gold">X</span>
              </span>
              <button
                className="p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
                onClick={() => setMobileOpen(false)}
                aria-label="Close menu"
              >
                <X size={20} />
              </button>
            </div>
            <nav className="flex flex-col gap-1 p-4 flex-1">
              {navLinks?.map((link) => (
                <Link
                  key={`mobile-nav-${link?.href}`}
                  href={link?.href}
                  onClick={() => setMobileOpen(false)}
                  className={`px-4 py-3 rounded-lg text-sm font-medium transition-colors ${
                    pathname === link?.href
                      ? 'bg-muted text-foreground'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                  }`}
                >
                  {link?.label}
                </Link>
              ))}
            </nav>
            <div className="p-4 border-t border-border flex flex-col gap-2">
              <button className="btn-ghost w-full justify-center">
                <LogIn size={15} />
                Sign In
              </button>
              <button className="btn-primary w-full justify-center">
                <User size={15} />
                Join Archive
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
