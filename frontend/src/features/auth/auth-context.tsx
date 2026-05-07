import { createContext, useContext, useEffect, useState, useCallback, useRef, type ReactNode } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import * as authApi from '@/features/auth/api';
import type { User } from '@/features/auth/api';

type AuthContextValue = {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string, remember?: boolean) => Promise<void>;
  signup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const queryClient = useQueryClient();
  const prevUserId = useRef<number | null>(null);

  useEffect(() => {
    authApi
      .getMe()
      .then((u) => {
        if (prevUserId.current !== null && prevUserId.current !== u.id) {
          queryClient.clear();
        }
        prevUserId.current = u.id;
        setUser(u);
      })
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  const login = useCallback(async (email: string, password: string, remember?: boolean) => {
    const u = await authApi.login({ email, password, remember });
    if (prevUserId.current !== null && prevUserId.current !== u.id) {
      queryClient.clear();
    }
    prevUserId.current = u.id;
    setUser(u);
  }, [queryClient]);

  const signup = useCallback(async (name: string, email: string, password: string) => {
    const u = await authApi.signup({ name, email, password });
    prevUserId.current = u.id;
    setUser(u);
  }, []);

  const logout = useCallback(async () => {
    await authApi.logout();
    prevUserId.current = null;
    setUser(null);
    queryClient.clear();
  }, [queryClient]);

  return (
    <AuthContext.Provider value={{ user, loading, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
