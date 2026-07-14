import React from 'react';
import Navbar from '@/components/navbar';
import PerformancesContent from '@/app/performances/components/PerformancesContent';
import ArchiveFooter from '@/app/components/ArchiveFooter';
import { getPerformances, Performance } from '@/lib/api';

export default async function PerformancesPage() {
  let performances: Performance[] = [];
  try {
    performances = await getPerformances(100);
  } catch (error) {
    console.error("Failed to fetch performances for page:", error);
  }

  return (
    <main className="min-h-screen bg-background">
      <Navbar />
      <PerformancesContent initialPerformances={performances} />
      <ArchiveFooter />
    </main>
  );
}
