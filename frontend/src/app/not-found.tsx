import React from 'react';
import Link from 'next/link';
import { Home } from 'lucide-react';

export default function NotFound() {
  return (
    <main className="mesh-bg grain relative flex min-h-[70vh] items-center justify-center">
      <div className="container-max container-px relative flex flex-col items-center text-center">
        <span className="font-display text-[clamp(4rem,14vw,9rem)] font-bold leading-none text-gradient-lime">
          404
        </span>
        <h1 className="mt-2 font-display text-2xl font-bold text-foreground">Off the pitch</h1>
        <p className="mt-2 max-w-md text-muted-foreground">
          We couldn&apos;t find that page. It may have been moved, or the match hasn&apos;t kicked off yet.
        </p>
        <Link href="/" className="btn-primary mt-8">
          <Home size={16} /> Back home
        </Link>
      </div>
    </main>
  );
}
