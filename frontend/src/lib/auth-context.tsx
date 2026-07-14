'use client';

import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import {
  getMe,
  loginUser,
  logoutUser,
  refreshAccessToken,
  registerUser,
  verifyOtp,
  forgotPassword as apiForgotPassword,
  resetPassword as apiResetPassword,
  type User,
  type AuthTokens,
} from '@/lib/api';

// ─── Types ────────────────────────────────────────────────────────────────────

interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isLoading: boolean;
}

interface AuthContextValue extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  register: (username: string, displayName: string, email: string, password: string) => Promise<{ email: string; otp_code_dev?: string }>;
  confirmOtp: (email: string, otp: string, purpose: 'LOGIN' | 'REGISTER') => Promise<void>;
  forgotPassword: (email: string) => Promise<void>;
  resetPassword: (email: string, otp: string, newPassword: string) => Promise<void>;
  logout: () => Promise<void>;
  setAuthData: (data: AuthTokens) => void;
  isAuthenticated: boolean;
}

// ─── Context ──────────────────────────────────────────────────────────────────

const AuthContext = createContext<AuthContextValue | null>(null);

const STORAGE_KEY_TOKEN = 'px_access_token';
const STORAGE_KEY_REFRESH = 'px_refresh_token';

// ─── Provider ─────────────────────────────────────────────────────────────────

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    token: null,
    refreshToken: null,
    isLoading: true,
  });

  // On mount, try to restore session
  useEffect(() => {
    const restore = async () => {
      const storedToken = localStorage.getItem(STORAGE_KEY_TOKEN);
      const storedRefresh = localStorage.getItem(STORAGE_KEY_REFRESH);

      if (!storedToken || !storedRefresh) {
        setState((s) => ({ ...s, isLoading: false }));
        return;
      }

      try {
        // Try to fetch current user with stored token
        const user = await getMe(storedToken);
        setState({ user, token: storedToken, refreshToken: storedRefresh, isLoading: false });
      } catch {
        // Token expired — try refresh
        try {
          const refreshed = await refreshAccessToken(storedRefresh);
          const user = await getMe(refreshed.access_token);
          localStorage.setItem(STORAGE_KEY_TOKEN, refreshed.access_token);
          localStorage.setItem(STORAGE_KEY_REFRESH, refreshed.refresh_token);
          setState({
            user,
            token: refreshed.access_token,
            refreshToken: refreshed.refresh_token,
            isLoading: false,
          });
        } catch {
          // Refresh failed — clear session
          localStorage.removeItem(STORAGE_KEY_TOKEN);
          localStorage.removeItem(STORAGE_KEY_REFRESH);
          setState({ user: null, token: null, refreshToken: null, isLoading: false });
        }
      }
    };

    restore();
  }, []);

  const applySession = useCallback((tokens: AuthTokens) => {
    localStorage.setItem(STORAGE_KEY_TOKEN, tokens.access_token);
    localStorage.setItem(STORAGE_KEY_REFRESH, tokens.refresh_token);
    setState({
      user: tokens.user,
      token: tokens.access_token,
      refreshToken: tokens.refresh_token,
      isLoading: false,
    });
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const tokens = await loginUser(email, password);
    applySession(tokens);
  }, [applySession]);

  const register = useCallback(async (username: string, displayName: string, email: string, password: string) => {
    const result = await registerUser(username, displayName, email, password);
    return { email: result.email, otp_code_dev: result.otp_code_dev };
  }, []);

  const confirmOtp = useCallback(
    async (email: string, otp: string, purpose: 'LOGIN' | 'REGISTER') => {
      const tokens = await verifyOtp(email, otp, purpose);
      applySession(tokens);
    },
    [applySession],
  );

  const forgotPassword = useCallback(async (email: string) => {
    await apiForgotPassword(email);
  }, []);

  const resetPassword = useCallback(async (email: string, otp: string, newPassword: string) => {
    await apiResetPassword(email, otp, newPassword);
  }, []);

  const logout = useCallback(async () => {
    const { token, refreshToken } = state;
    if (token && refreshToken) {
      try {
        await logoutUser(token, refreshToken);
      } catch { /* ignore */ }
    }
    localStorage.removeItem(STORAGE_KEY_TOKEN);
    localStorage.removeItem(STORAGE_KEY_REFRESH);
    setState({ user: null, token: null, refreshToken: null, isLoading: false });
  }, [state]);

  const setAuthData = useCallback((tokens: AuthTokens) => {
    applySession(tokens);
  }, [applySession]);

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        register,
        confirmOtp,
        forgotPassword,
        resetPassword,
        logout,
        setAuthData,
        isAuthenticated: !!state.user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// ─── Hook ─────────────────────────────────────────────────────────────────────

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
