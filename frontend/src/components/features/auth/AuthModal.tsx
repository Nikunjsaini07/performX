'use client';

import React, { useState, useEffect, useRef } from 'react';
import { X, Mail, Lock, User, AtSign, Eye, EyeOff, Loader2, CheckCircle2, KeyRound, ArrowLeft } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';

// ─── Types ────────────────────────────────────────────────────────────────────

type Tab = 'login' | 'register';
type Step = 'form' | 'otp' | 'success' | 'forgot-email' | 'forgot-otp' | 'forgot-newpw' | 'forgot-success';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  defaultTab?: Tab;
}

// ─── OTP Input ────────────────────────────────────────────────────────────────

function OtpInput({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const inputs = useRef<(HTMLInputElement | null)[]>([]);

  const digits = value.padEnd(6, '').split('').slice(0, 6);

  const handleKeyDown = (i: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace') {
      if (digits[i]) {
        const next = digits.slice();
        next[i] = '';
        onChange(next.join(''));
      } else if (i > 0) {
        inputs.current[i - 1]?.focus();
        const next = digits.slice();
        next[i - 1] = '';
        onChange(next.join(''));
      }
    }
  };

  const handleInput = (i: number, e: React.ChangeEvent<HTMLInputElement>) => {
    const ch = e.target.value.replace(/\D/g, '').slice(-1);
    const next = digits.slice();
    next[i] = ch;
    onChange(next.join('').trim());
    if (ch && i < 5) inputs.current[i + 1]?.focus();
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    const pasted = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6);
    onChange(pasted);
    const lastFilled = Math.min(pasted.length, 5);
    inputs.current[lastFilled]?.focus();
    e.preventDefault();
  };

  return (
    <div className="flex gap-2 justify-center" onPaste={handlePaste}>
      {[0, 1, 2, 3, 4, 5].map((i) => (
        <input
          key={i}
          ref={(el) => { inputs.current[i] = el; }}
          type="text"
          inputMode="numeric"
          maxLength={1}
          value={digits[i] || ''}
          onChange={(e) => handleInput(i, e)}
          onKeyDown={(e) => handleKeyDown(i, e)}
          className="w-11 h-14 text-center text-xl font-bold rounded-lg border border-border bg-muted text-foreground focus:outline-none focus:ring-2 focus:ring-primary/60 focus:border-primary transition-colors"
        />
      ))}
    </div>
  );
}

// ─── Modal ────────────────────────────────────────────────────────────────────

export default function AuthModal({ isOpen, onClose, defaultTab = 'login' }: AuthModalProps) {
  const { login, register, confirmOtp, forgotPassword, resetPassword } = useAuth();

  const [tab, setTab] = useState<Tab>(defaultTab);
  const [step, setStep] = useState<Step>('form');
  const [pendingEmail, setPendingEmail] = useState('');
  const [pendingPurpose, setPendingPurpose] = useState<'LOGIN' | 'REGISTER'>('LOGIN');
  const [otpValue, setOtpValue] = useState('');

  // Login form
  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [showLoginPw, setShowLoginPw] = useState(false);

  // Register form
  const [regUsername, setRegUsername] = useState('');
  const [regDisplay, setRegDisplay] = useState('');
  const [regEmail, setRegEmail] = useState('');
  const [regPassword, setRegPassword] = useState('');
  const [showRegPw, setShowRegPw] = useState(false);

  // Forgot Password
  const [forgotEmail, setForgotEmail] = useState('');
  const [forgotOtp, setForgotOtp] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [showNewPw, setShowNewPw] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen) {
      setTab(defaultTab);
      setStep('form');
      setError('');
      setOtpValue('');
      setForgotEmail('');
      setForgotOtp('');
      setNewPassword('');
    }
  }, [isOpen, defaultTab]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [onClose]);

  if (!isOpen) return null;

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(loginEmail, loginPassword);
      setStep('success');
      setTimeout(() => onClose(), 1200);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const result = await register(regUsername, regDisplay, regEmail, regPassword);
      setPendingEmail(result.email);
      setPendingPurpose('REGISTER');
      setStep('otp');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  const handleOtpSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (otpValue.length < 6) { setError('Enter the full 6-digit code'); return; }
    setError('');
    setLoading(true);
    try {
      await confirmOtp(pendingEmail, otpValue, pendingPurpose);
      setStep('success');
      setTimeout(() => onClose(), 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid code');
    } finally {
      setLoading(false);
    }
  };

  const handleForgotEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await forgotPassword(forgotEmail);
      setStep('forgot-otp');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to send reset code');
    } finally {
      setLoading(false);
    }
  };

  const handleForgotOtpSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (forgotOtp.length < 6) { setError('Enter the full 6-digit code'); return; }
    setError('');
    setStep('forgot-newpw');
  };

  const handleResetPasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword.length < 6) { setError('Password must be at least 6 characters'); return; }
    setError('');
    setLoading(true);
    try {
      await resetPassword(forgotEmail, forgotOtp, newPassword);
      setStep('forgot-success');
      setTimeout(() => { setStep('form'); setTab('login'); }, 2500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to reset password');
    } finally {
      setLoading(false);
    }
  };

  const inputClass = 'w-full rounded-lg border border-border bg-muted px-4 py-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-colors';

  return (
    <div className="fixed inset-0 z-[200] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-background/80 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-card shadow-2xl">
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div className="flex items-center gap-2">
            <span className="font-bold text-lg text-foreground">
              Perform<span className="text-gradient-gold">X</span>
            </span>
          </div>
          <button onClick={onClose} className="p-2 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors" aria-label="Close">
            <X size={18} />
          </button>
        </div>

        <div className="p-6">
          {step === 'success' && (
            <div className="flex flex-col items-center gap-4 py-8 text-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-accent/15">
                <CheckCircle2 size={32} className="text-accent" />
              </div>
              <div>
                <h2 className="font-display text-xl font-bold text-foreground">Welcome to PerformX!</h2>
                <p className="mt-1 text-sm text-muted-foreground">You are now signed in.</p>
              </div>
            </div>
          )}

          {step === 'otp' && (
            <form onSubmit={handleOtpSubmit} className="space-y-5">
              <div className="text-center">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/15 mb-3">
                  <Mail size={22} className="text-primary" />
                </div>
                <h2 className="font-display text-xl font-bold text-foreground">Check your email</h2>
                <p className="mt-1 text-sm text-muted-foreground">
                  We sent a 6-digit code to <span className="text-foreground font-medium">{pendingEmail}</span>
                </p>
              </div>
              <OtpInput value={otpValue} onChange={setOtpValue} />
              {error && <p className="text-sm text-red-400 text-center">{error}</p>}
              <button type="submit" disabled={loading || otpValue.length < 6} className="btn-primary w-full justify-center disabled:opacity-50">
                {loading ? <Loader2 size={16} className="animate-spin" /> : 'Verify Code'}
              </button>
              <button type="button" onClick={() => { setStep('form'); setOtpValue(''); setError(''); }} className="w-full text-center text-sm text-muted-foreground hover:text-foreground transition-colors">
                ← Go back
              </button>
            </form>
          )}

          {step === 'forgot-email' && (
            <form onSubmit={handleForgotEmailSubmit} className="space-y-5">
              <div className="text-center">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/15 mb-3">
                  <KeyRound size={22} className="text-primary" />
                </div>
                <h2 className="font-display text-xl font-bold text-foreground">Forgot Password?</h2>
                <p className="mt-1 text-sm text-muted-foreground">Enter your email and we&apos;ll send you a reset code.</p>
              </div>
              <div>
                <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Email</label>
                <div className="relative">
                  <Mail size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                  <input type="email" required placeholder="you@example.com" value={forgotEmail} onChange={(e) => setForgotEmail(e.target.value)} className={`${inputClass} pl-9`} />
                </div>
              </div>
              {error && <p className="text-sm text-red-400">{error}</p>}
              <button type="submit" disabled={loading} className="btn-primary w-full justify-center disabled:opacity-50">
                {loading ? <Loader2 size={16} className="animate-spin" /> : 'Send Reset Code'}
              </button>
              <button type="button" onClick={() => { setStep('form'); setError(''); }} className="w-full flex items-center justify-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors">
                <ArrowLeft size={14} /> Back to Sign In
              </button>
            </form>
          )}

          {step === 'forgot-otp' && (
            <form onSubmit={handleForgotOtpSubmit} className="space-y-5">
              <div className="text-center">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/15 mb-3">
                  <Mail size={22} className="text-primary" />
                </div>
                <h2 className="font-display text-xl font-bold text-foreground">Enter Reset Code</h2>
                <p className="mt-1 text-sm text-muted-foreground">
                  We sent a 6-digit code to <span className="text-foreground font-medium">{forgotEmail}</span>
                </p>
              </div>
              <OtpInput value={forgotOtp} onChange={setForgotOtp} />
              {error && <p className="text-sm text-red-400 text-center">{error}</p>}
              <button type="submit" disabled={loading || forgotOtp.length < 6} className="btn-primary w-full justify-center disabled:opacity-50">
                {loading ? <Loader2 size={16} className="animate-spin" /> : 'Verify Code'}
              </button>
              <button type="button" onClick={() => { setStep('forgot-email'); setForgotOtp(''); setError(''); }} className="w-full flex items-center justify-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors">
                <ArrowLeft size={14} /> Change email
              </button>
            </form>
          )}

          {step === 'forgot-newpw' && (
            <form onSubmit={handleResetPasswordSubmit} className="space-y-5">
              <div className="text-center">
                <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/15 mb-3">
                  <Lock size={22} className="text-primary" />
                </div>
                <h2 className="font-display text-xl font-bold text-foreground">Set New Password</h2>
                <p className="mt-1 text-sm text-muted-foreground">Choose a strong new password for your account.</p>
              </div>
              <div>
                <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">New Password</label>
                <div className="relative">
                  <Lock size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                  <input type={showNewPw ? 'text' : 'password'} required minLength={6} placeholder="Min. 6 characters" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} className={`${inputClass} pl-9 pr-10`} />
                  <button type="button" onClick={() => setShowNewPw((v) => !v)} className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
                    {showNewPw ? <EyeOff size={15} /> : <Eye size={15} />}
                  </button>
                </div>
              </div>
              {error && <p className="text-sm text-red-400">{error}</p>}
              <button type="submit" disabled={loading} className="btn-primary w-full justify-center disabled:opacity-50">
                {loading ? <Loader2 size={16} className="animate-spin" /> : 'Reset Password'}
              </button>
            </form>
          )}

          {step === 'forgot-success' && (
            <div className="flex flex-col items-center gap-4 py-8 text-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-accent/15">
                <CheckCircle2 size={32} className="text-accent" />
              </div>
              <div>
                <h2 className="font-display text-xl font-bold text-foreground">Password Reset!</h2>
                <p className="mt-1 text-sm text-muted-foreground">Your password has been updated. Redirecting to sign in…</p>
              </div>
            </div>
          )}

          {step === 'form' && (
            <>
              <div className="flex rounded-lg border border-border bg-muted p-1 mb-6">
                {(['login', 'register'] as Tab[]).map((t) => (
                  <button
                    key={t}
                    onClick={() => { setTab(t); setError(''); }}
                    className={`flex-1 rounded-md py-2 text-sm font-medium transition-all ${
                      tab === t ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    {t === 'login' ? 'Sign In' : 'Create Account'}
                  </button>
                ))}
              </div>

              {tab === 'login' && (
                <form onSubmit={handleLoginSubmit} className="space-y-4">
                  <button
                    type="button"
                    onClick={() => {
                      const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
                      const redirectUri = `${process.env.NEXT_PUBLIC_API_URL}/auth/google/callback`;
                      const scope = 'email profile';
                      const authUrl = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&response_type=code&scope=${encodeURIComponent(scope)}&access_type=offline&prompt=consent`;
                      window.location.href = authUrl;
                    }}
                    className="w-full flex items-center justify-center gap-3 rounded-lg border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-colors hover:bg-surface-2"
                  >
                    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M17.64 9.20443C17.64 8.56625 17.5827 7.95262 17.4764 7.36353H9V10.8449H13.8436C13.635 11.9699 13.0009 12.9231 12.0477 13.5613V15.8194H14.9564C16.6582 14.2526 17.64 11.9453 17.64 9.20443Z" fill="#4285F4"/>
                      <path d="M8.99976 18C11.4298 18 13.467 17.1941 14.9561 15.8195L12.0475 13.5613C11.2416 14.1013 10.2107 14.4204 8.99976 14.4204C6.65567 14.4204 4.67158 12.8372 3.96385 10.71H0.957031V13.0418C2.43794 15.9831 5.48158 18 8.99976 18Z" fill="#34A853"/>
                      <path d="M3.96409 10.7098C3.78409 10.1698 3.68182 9.59301 3.68182 8.99983C3.68182 8.40665 3.78409 7.82983 3.96409 7.28983V4.95801H0.957273C0.347727 6.17301 0 7.54755 0 8.99983C0 10.4521 0.347727 11.8266 0.957273 13.0416L3.96409 10.7098Z" fill="#FBBC05"/>
                      <path d="M8.99976 3.57955C10.3211 3.57955 11.5075 4.03364 12.4402 4.92545L15.0216 2.34409C13.4629 0.891818 11.4257 0 8.99976 0C5.48158 0 2.43794 2.01682 0.957031 4.95818L3.96385 7.29C4.67158 5.16273 6.65567 3.57955 8.99976 3.57955Z" fill="#EA4335"/>
                    </svg>
                    Continue with Google
                  </button>
                  <div className="relative">
                    <div className="absolute inset-0 flex items-center"><div className="w-full border-t border-border"></div></div>
                    <div className="relative flex justify-center text-xs"><span className="bg-card px-2 text-muted-foreground">Or continue with email</span></div>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Email</label>
                    <div className="relative">
                      <Mail size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                      <input type="email" required placeholder="you@example.com" value={loginEmail} onChange={(e) => setLoginEmail(e.target.value)} className={`${inputClass} pl-9`} />
                    </div>
                  </div>
                  <div>
                    <div className="flex items-center justify-between mb-1.5">
                      <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground">Password</label>
                      <button type="button" onClick={() => { setForgotEmail(loginEmail); setStep('forgot-email'); setError(''); }} className="text-xs text-primary hover:underline">
                        Forgot password?
                      </button>
                    </div>
                    <div className="relative">
                      <Lock size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                      <input type={showLoginPw ? 'text' : 'password'} required placeholder="••••••••" value={loginPassword} onChange={(e) => setLoginPassword(e.target.value)} className={`${inputClass} pl-9 pr-10`} />
                      <button type="button" onClick={() => setShowLoginPw((v) => !v)} className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
                        {showLoginPw ? <EyeOff size={15} /> : <Eye size={15} />}
                      </button>
                    </div>
                  </div>
                  {error && <p className="text-sm text-red-400">{error}</p>}
                  <button type="submit" disabled={loading} className="btn-primary w-full justify-center disabled:opacity-50">
                    {loading ? <Loader2 size={16} className="animate-spin" /> : 'Continue'}
                  </button>
                </form>
              )}

              {tab === 'register' && (
                <form onSubmit={handleRegisterSubmit} className="space-y-4">
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Username</label>
                      <div className="relative">
                        <AtSign size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                        <input type="text" required placeholder="username" value={regUsername} onChange={(e) => setRegUsername(e.target.value)} className={`${inputClass} pl-9`} />
                      </div>
                    </div>
                    <div>
                      <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Display Name</label>
                      <div className="relative">
                        <User size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                        <input type="text" required placeholder="Your Name" value={regDisplay} onChange={(e) => setRegDisplay(e.target.value)} className={`${inputClass} pl-9`} />
                      </div>
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Email</label>
                    <div className="relative">
                      <Mail size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                      <input type="email" required placeholder="you@example.com" value={regEmail} onChange={(e) => setRegEmail(e.target.value)} className={`${inputClass} pl-9`} />
                    </div>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">Password</label>
                    <div className="relative">
                      <Lock size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                      <input type={showRegPw ? 'text' : 'password'} required minLength={6} placeholder="Min. 6 characters" value={regPassword} onChange={(e) => setRegPassword(e.target.value)} className={`${inputClass} pl-9 pr-10`} />
                      <button type="button" onClick={() => setShowRegPw((v) => !v)} className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
                        {showRegPw ? <EyeOff size={15} /> : <Eye size={15} />}
                      </button>
                    </div>
                  </div>
                  {error && <p className="text-sm text-red-400">{error}</p>}
                  <button type="submit" disabled={loading} className="btn-primary w-full justify-center disabled:opacity-50">
                    {loading ? <Loader2 size={16} className="animate-spin" /> : 'Create Account'}
                  </button>
                </form>
              )}

              <p className="mt-4 text-center text-xs text-muted-foreground">
                {tab === 'login' ? "Don't have an account? " : 'Already have an account? '}
                <button onClick={() => { setTab(tab === 'login' ? 'register' : 'login'); setError(''); }} className="text-primary hover:underline">
                  {tab === 'login' ? 'Create one' : 'Sign in'}
                </button>
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
