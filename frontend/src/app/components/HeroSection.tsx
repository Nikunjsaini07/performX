import React from 'react';
import Link from 'next/link';
import { ArrowRight, Database, Users, Zap, Star } from 'lucide-react';
import AppImage from '@/components/ui/AppImage';

export default function HeroSection() {
  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center overflow-hidden">
      {/* Stadium Background — massive night stadium packed with audience */}
      <div className="absolute inset-0 z-0">
        <AppImage
          src="https://images.unsplash.com/photo-1702411854093-fc63f696202e"
          alt="Massive night stadium packed with audience under floodlights during FIFA 2026"
          fill
          priority
          className="object-cover object-center scale-105"
          sizes="100vw" />
        
        {/* Deep cinematic overlay — preserves crowd energy while keeping text readable */}
        <div className="absolute inset-0 bg-gradient-to-b from-[#050508]/85 via-[#050508]/60 to-[#050508]/98" />
        {/* Atmospheric side vignettes */}
        <div className="absolute inset-0 bg-gradient-to-r from-[#050508]/80 via-transparent to-[#050508]/80" />
        {/* Extra darkness layer for depth */}
        <div className="absolute inset-0 bg-[#050508]/30" />
        {/* Floodlight glow effect */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[80%] h-[50%] bg-[#c9a84c]/5 blur-[120px] rounded-full pointer-events-none" />
      </div>

      {/* Noise grain for premium feel */}
      <div className="absolute inset-0 z-[1] opacity-[0.025] noise-overlay pointer-events-none" />

      {/* Content */}
      <div className="relative z-10 max-w-screen-2xl mx-auto px-6 lg:px-10 w-full flex flex-col items-center text-center pt-24 pb-20">
        {/* Archive label */}
        <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-border bg-card/60 backdrop-blur-sm mb-8">
          <span className="w-1.5 h-1.5 rounded-full bg-accent animate-pulse" />
          <span className="text-xs font-semibold tracking-widest uppercase text-muted-foreground">
            FIFA World Cup 2026 — Performance Archive
          </span>
        </div>

        {/* Hero headline */}
        <h1 className="font-display text-5xl md:text-7xl lg:text-8xl font-bold leading-tight tracking-tight max-w-5xl mb-6">
          <span className="text-foreground">FIFA 2026,</span>
          <br />
          <span className="text-gradient-gold italic">remembered</span>
          <br />
          <span className="text-foreground">through performances.</span>
        </h1>

        <p className="text-base md:text-lg text-muted-foreground max-w-2xl mb-10 leading-relaxed">
          Track every match, every player performance, and every moment fans will remember.
          96 matches. 421 players. 734 performances archived.
        </p>

        {/* CTA Buttons */}
        <div className="flex flex-col sm:flex-row items-center gap-4 mb-16">
          <Link href="/performances" className="btn-primary px-7 py-3 text-base gap-2">
            Browse Performances
            <ArrowRight size={16} />
          </Link>
          <Link href="/matches" className="btn-ghost px-7 py-3 text-base">
            Explore Matches
          </Link>
        </div>

        {/* Archive Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-8 w-full max-w-3xl">
          {[
          { icon: <Database size={18} />, value: '96', label: 'Matches Archived' },
          { icon: <Users size={18} />, value: '421', label: 'Players Tracked' },
          { icon: <Zap size={18} />, value: '734', label: 'Performances Logged' },
          { icon: <Star size={18} />, value: '4,404', label: 'Performance Stats' }]?.
          map((stat) =>
          <div
            key={`hero-stat-${stat?.label}`}
            className="flex flex-col items-center gap-1 px-4 py-4 rounded-xl border border-border bg-card/40 backdrop-blur-sm">
            
              <span className="text-primary mb-1">{stat?.icon}</span>
              <span className="stat-number text-2xl font-bold text-foreground">{stat?.value}</span>
              <span className="text-xs text-muted-foreground font-medium text-center">{stat?.label}</span>
            </div>
          )}
        </div>
      </div>

      {/* Bottom fade into page */}
      <div className="absolute bottom-0 left-0 right-0 h-40 bg-gradient-to-t from-background to-transparent z-10" />
    </section>);

}