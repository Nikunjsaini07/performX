import React from 'react';
import { BACKGROUNDS } from '@/lib/site';

interface DetailBackdropProps {
  /** Optional entity-specific cover image; falls back to the shared backdrop. */
  imageUrl?: string;
  className?: string;
  /**
   * Controls how tall/positioned the backdrop wrapper is. Defaults to a
   * generous height that bleeds well past the header content so the
   * internal vignette has room to fade smoothly into `var(--background)`
   * instead of cutting off abruptly at the end of the header.
   */
  heightClassName?: string;
}

/**
 * Shared cinematic backdrop for match & performance detail pages.
 * Letterboxd-style: image behind the header with a heavy duotone/dark
 * treatment (dark gradient + lime radial glow + grain) fading into the page.
 */
export default function DetailBackdrop({
  imageUrl,
  className = '',
  heightClassName = 'absolute inset-x-0 top-0 h-[90vh] min-h-[780px]',
}: DetailBackdropProps) {
  const src = imageUrl || BACKGROUNDS.detail;
  return (
    <div className={`detail-backdrop ${heightClassName} ${className}`} aria-hidden="true">

      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={src}
        alt=""
        className="h-full w-full object-cover opacity-40"
        style={{ filter: 'grayscale(0.35) contrast(1.05)' }}
      />
    </div>
  );
}
