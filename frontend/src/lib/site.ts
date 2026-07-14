// Central site configuration: brand, navigation, and the two shared background
// images (hosted on Cloudinary).

export const SITE = {
  name: 'PerformX',
  tagline: 'Rate the World Cup. One match, one performance at a time.',
  description:
    'PerformX is the community archive for FIFA World Cup 2026 — rate matches and player performances, write reviews, and see what the world is watching.',
};

// Cloudinary-hosted backgrounds (uploaded under performx/backgrounds).
export const BACKGROUNDS = {
  // Full-bleed home hero: stadium sunset with FIFA World Cup 2026.
  home: 'https://res.cloudinary.com/dbcflpua9/image/upload/f_auto,q_auto/v1783955892/performx/backgrounds/home-hero-stadium.png',
  // Atmospheric pitch used behind the home rails section (below the hero).
  ambience:
    'https://res.cloudinary.com/dbcflpua9/image/upload/v1783875940/performx/backgrounds/ambience-pitch.jpg',
  // Shared Letterboxd-style backdrop for every match & performance detail page.
  detail:
    'https://res.cloudinary.com/dbcflpua9/image/upload/v1783875939/performx/backgrounds/detail-stadium.jpg',
};

export const NAV_LINKS = [
  { label: 'Home', href: '/' },
  { label: 'Matches', href: '/matches' },
  { label: 'Performances', href: '/performances' },
  { label: 'Teams', href: '/teams' },
];
