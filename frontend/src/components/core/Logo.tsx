import React from 'react';

interface LogoProps {
  className?: string;
  size?: number;
}

export default function Logo({ className = '', size = 32 }: LogoProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
    >
      {/* Circle background */}
      <circle cx="50" cy="50" r="48" fill="currentColor" opacity="0.1" />
      
      {/* Football/Soccer ball pattern */}
      <circle cx="50" cy="50" r="40" stroke="currentColor" strokeWidth="3" fill="none" />
      
      {/* Pentagon in center */}
      <path
        d="M50 20 L65 35 L60 55 L40 55 L35 35 Z"
        fill="currentColor"
        opacity="0.2"
      />
      
      {/* P letter (left side) */}
      <path
        d="M30 35 L30 65 M30 35 L42 35 C47 35 50 38 50 43 C50 48 47 51 42 51 L30 51"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      />
      
      {/* X letter (right side) */}
      <path
        d="M55 35 L70 65 M70 35 L55 65"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
