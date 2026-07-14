'use client';

import { useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/lib/auth-context';
import { Loader2 } from 'lucide-react';

function CallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setAuthData } = useAuth();

  useEffect(() => {
    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const userId = searchParams.get('user_id');
    const username = searchParams.get('username');
    const displayName = searchParams.get('display_name');
    const email = searchParams.get('email');
    const authError = searchParams.get('auth_error');

    if (authError) {
      console.error('Google OAuth error:', authError);
      router.push('/?auth_error=' + authError);
      return;
    }

    if (accessToken && refreshToken && userId && username && displayName && email) {
      // Store tokens in localStorage (using the same keys as auth-context)
      localStorage.setItem('px_access_token', accessToken);
      localStorage.setItem('px_refresh_token', refreshToken);

      // Update auth context
      setAuthData({
        access_token: accessToken,
        refresh_token: refreshToken,
        user: {
          id: userId,
          username: username,
          display_name: displayName,
          email: email,
          bio: '',
          avatar_url: '',
        },
      });

      // Redirect to home page
      router.push('/');
    } else {
      console.error('Missing OAuth parameters');
      router.push('/?auth_error=missing_parameters');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Run only once on mount

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <Loader2 size={48} className="animate-spin text-primary mx-auto mb-4" />
        <h2 className="text-xl font-bold text-foreground">Signing you in...</h2>
        <p className="text-sm text-muted-foreground mt-2">Please wait a moment</p>
      </div>
    </div>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 size={48} className="animate-spin text-primary" />
      </div>
    }>
      <CallbackContent />
    </Suspense>
  );
}
