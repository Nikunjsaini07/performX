'use client';

import { useState, useEffect, useRef } from 'react';
import { useAuth } from '@/lib/auth-context';
import { useRouter } from 'next/navigation';
import { updateProfile, updateUsername, updateAvatar, uploadToCloudinary } from '@/lib/api';
import { Camera, Loader2, User, AtSign, FileText, Check, X } from 'lucide-react';

export default function SettingsPage() {
  const { user, token, isAuthenticated } = useAuth();
  const router = useRouter();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [displayName, setDisplayName] = useState('');
  const [username, setUsername] = useState('');
  const [bio, setBio] = useState('');
  const [avatarUrl, setAvatarUrl] = useState('');

  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [checkingUsername, setCheckingUsername] = useState(false);
  const [usernameAvailable, setUsernameAvailable] = useState<boolean | null>(null);
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/');
      return;
    }
    if (user) {
      setDisplayName(user.display_name || '');
      setUsername(user.username || '');
      setBio(user.bio || '');
      setAvatarUrl(user.avatar_url || '');
    }
  }, [user, isAuthenticated, router]);

  // Check username availability as user types
  useEffect(() => {
    const checkUsername = async () => {
      if (!username || username === user?.username || username.length < 3) {
        setUsernameAvailable(null);
        return;
      }

      setCheckingUsername(true);
      try {
        // Simple check by trying to fetch the profile - if it doesn't exist, it's available
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/users/${username}`);
        setUsernameAvailable(!res.ok); // Available if user doesn't exist (404)
      } catch {
        setUsernameAvailable(null);
      } finally {
        setCheckingUsername(false);
      }
    };

    const timeoutId = setTimeout(checkUsername, 500); // Debounce
    return () => clearTimeout(timeoutId);
  }, [username, user?.username]);

  const handleAvatarClick = () => {
    fileInputRef.current?.click();
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !token) return;

    setUploading(true);
    setError('');
    try {
      const url = await uploadToCloudinary(file);
      await updateAvatar(token, url);
      setAvatarUrl(url);
      setSuccess('Avatar updated successfully!');
      setTimeout(() => setSuccess(''), 3000);
      window.location.reload(); // Refresh to update navbar
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to upload avatar');
    } finally {
      setUploading(false);
    }
  };

  const handleSaveProfile = async () => {
    if (!token) return;

    setLoading(true);
    setError('');
    try {
      await updateProfile(token, {
        display_name: displayName,
        bio: bio,
      });
      setSuccess('Profile updated successfully!');
      setTimeout(() => setSuccess(''), 3000);
      window.location.reload(); // Refresh to update UI
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update profile');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveUsername = async () => {
    if (!token || !username) return;

    // Basic validation
    if (username.length < 3) {
      setError('Username must be at least 3 characters');
      return;
    }
    if (username.length > 20) {
      setError('Username must be less than 20 characters');
      return;
    }
    if (!/^[a-z0-9_]+$/.test(username)) {
      setError('Username can only contain lowercase letters, numbers, and underscores');
      return;
    }

    setLoading(true);
    setError('');
    try {
      await updateUsername(token, username);
      setSuccess('Username updated successfully!');
      setTimeout(() => setSuccess(''), 3000);
      window.location.reload(); // Refresh to update UI
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Failed to update username';
      if (errorMsg.includes('already exists') || errorMsg.includes('Conflict')) {
        setError('This username is already taken. Please choose another one.');
      } else {
        setError(errorMsg);
      }
    } finally {
      setLoading(false);
    }
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <main className="container-max container-px py-8">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="font-display text-3xl font-bold text-foreground">Settings</h1>
          <p className="text-sm text-muted-foreground mt-2">Manage your account settings and profile</p>
        </div>

        {/* Success/Error Messages */}
        {success && (
          <div className="mb-6 flex items-center gap-2 rounded-lg border border-green-500/20 bg-green-500/10 px-4 py-3 text-sm text-green-400">
            <Check size={16} />
            {success}
          </div>
        )}
        {error && (
          <div className="mb-6 flex items-center gap-2 rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-400">
            <X size={16} />
            {error}
          </div>
        )}

        {/* Avatar Section */}
        <section className="card-shell p-6 mb-6">
          <h2 className="font-display text-xl font-bold text-foreground mb-4">Profile Picture</h2>
          <div className="flex items-center gap-6">
            <div className="relative">
              <div className="h-24 w-24 rounded-full overflow-hidden bg-surface-2 flex items-center justify-center">
                {avatarUrl ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={avatarUrl} alt="Avatar" className="h-full w-full object-cover" />
                ) : (
                  <span className="text-2xl font-bold text-muted-foreground">
                    {displayName ? displayName[0].toUpperCase() : 'U'}
                  </span>
                )}
              </div>
              <button
                onClick={handleAvatarClick}
                disabled={uploading}
                className="absolute bottom-0 right-0 p-2 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
                aria-label="Upload avatar"
              >
                {uploading ? <Loader2 size={16} className="animate-spin" /> : <Camera size={16} />}
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                onChange={handleAvatarChange}
                className="hidden"
              />
            </div>
            <div>
              <p className="text-sm text-foreground font-medium mb-1">Change your avatar</p>
              <p className="text-xs text-muted-foreground">JPG, PNG or GIF. Max size 5MB.</p>
            </div>
          </div>
        </section>

        {/* Display Name & Bio */}
        <section className="card-shell p-6 mb-6">
          <h2 className="font-display text-xl font-bold text-foreground mb-4">Profile Information</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
                <User size={12} className="inline mr-1" />
                Display Name
              </label>
              <input
                type="text"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                className="w-full rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-colors"
                placeholder="Your display name"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
                <FileText size={12} className="inline mr-1" />
                Bio
              </label>
              <textarea
                value={bio}
                onChange={(e) => setBio(e.target.value)}
                rows={4}
                className="w-full rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-colors resize-none"
                placeholder="Tell us about yourself..."
              />
            </div>
            <button
              onClick={handleSaveProfile}
              disabled={loading}
              className="btn-primary disabled:opacity-50"
            >
              {loading ? <Loader2 size={16} className="animate-spin" /> : 'Save Changes'}
            </button>
          </div>
        </section>

        {/* Username */}
        <section className="card-shell p-6">
          <h2 className="font-display text-xl font-bold text-foreground mb-4">Username</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-2">
                <AtSign size={12} className="inline mr-1" />
                Username
              </label>
              <div className="relative">
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value.toLowerCase().replace(/[^a-z0-9_]/g, ''))}
                  className={`w-full rounded-lg border ${
                    usernameAvailable === false ? 'border-red-500 focus:ring-red-500' : 
                    usernameAvailable === true ? 'border-green-500 focus:ring-green-500' : 
                    'border-border focus:ring-primary/50'
                  } bg-muted px-4 py-3 pr-10 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:border-primary transition-colors`}
                  placeholder="username"
                />
                {checkingUsername && (
                  <Loader2 size={16} className="absolute right-3 top-1/2 -translate-y-1/2 animate-spin text-muted-foreground" />
                )}
                {!checkingUsername && usernameAvailable === true && username !== user?.username && (
                  <Check size={16} className="absolute right-3 top-1/2 -translate-y-1/2 text-green-500" />
                )}
                {!checkingUsername && usernameAvailable === false && (
                  <X size={16} className="absolute right-3 top-1/2 -translate-y-1/2 text-red-500" />
                )}
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                Your profile will be available at: performx.com/u/{username || 'username'}
              </p>
              {usernameAvailable === false && username !== user?.username && (
                <p className="text-xs text-red-400 mt-1">This username is already taken</p>
              )}
              {usernameAvailable === true && username !== user?.username && (
                <p className="text-xs text-green-400 mt-1">This username is available!</p>
              )}
            </div>
            <button
              onClick={handleSaveUsername}
              disabled={loading || !username || username === user?.username || usernameAvailable === false}
              className="btn-primary disabled:opacity-50"
            >
              {loading ? <Loader2 size={16} className="animate-spin" /> : 'Update Username'}
            </button>
          </div>
        </section>
      </div>
    </main>
  );
}
