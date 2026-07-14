import React from 'react';
import Navbar from '@/components/navbar';
import MatchesContent from '@/app/matches/components/MatchesContent';
import ArchiveFooter from '@/app/components/ArchiveFooter';

import { getMatches, Match } from '@/lib/api';

export default async function MatchesPage() {
  let matches: Match[] = [];
  try {
    matches = await getMatches(100);
  } catch (error) {
    console.error("Failed to fetch matches:", error);
  }

  return (
    <main className="min-h-screen bg-background">
      <Navbar />
      <MatchesContent initialMatches={matches} />
      <ArchiveFooter />
    </main>
  );
}
