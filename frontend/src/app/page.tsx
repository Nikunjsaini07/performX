import React from 'react';
import Navbar from '@/components/navbar';
import HeroSection from '@/app/components/HeroSection';
import FeaturedMatches from '@/app/components/FeaturedMatches';
import TopPerformances from '@/app/components/TopPerformances';
import TopMatches from '@/app/components/TopMatches';
import StatLeaders from '@/app/components/StatLeaders';
import TeamsGrid from '@/app/components/TeamsGrid';
import RecentReviews from '@/app/components/RecentReviews';
import ArchiveFooter from '@/app/components/ArchiveFooter';

export default function HomePage() {
  return (
    <main className="min-h-screen bg-background">
      <Navbar />
      <HeroSection />
      <FeaturedMatches />
      <TopPerformances />
      <TopMatches />
      <StatLeaders />
      <RecentReviews />
      <TeamsGrid />
      <ArchiveFooter />
    </main>
  );
}
