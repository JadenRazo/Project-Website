import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import type { AuthState, User } from './types';

interface AuthActions {
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
  validateToken: (token: string) => Promise<void>;
  setLoading: (loading: boolean) => void;
  setUser: (user: User | null) => void;
  setAuthenticated: (authenticated: boolean) => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      (set, get) => ({
        user: null,
        isAuthenticated: false,
        isLoading: true,

        setLoading: (loading: boolean) =>
          set({ isLoading: loading }, false, 'setLoading'),

        setUser: (user: User | null) =>
          set({ user }, false, 'setUser'),

        setAuthenticated: (authenticated: boolean) =>
          set({ isAuthenticated: authenticated }, false, 'setAuthenticated'),

        validateToken: async (token: string) => {
          try {
            const response = await fetch('/api/v1/auth/validate', {
              headers: {
                'Authorization': `Bearer ${token}`,
              },
            });

            if (response.ok) {
              const userData = await response.json();
              set({
                user: { ...userData, token },
                isAuthenticated: true,
                isLoading: false,
              }, false, 'validateToken/success');
            } else {
              localStorage.removeItem('auth_token');
              set({
                user: null,
                isAuthenticated: false,
                isLoading: false,
              }, false, 'validateToken/invalid');
            }
          } catch (error) {
            console.error('Error validating token:', error);
            set({
              user: null,
              isAuthenticated: false,
              isLoading: false,
            }, false, 'validateToken/error');
          }
        },

        login: async (email: string, password: string) => {
          try {
            const response = await fetch('/api/v1/auth/login', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
              throw new Error('Login failed');
            }

            const userData = await response.json();
            localStorage.setItem('auth_token', userData.token);
            
            set({
              user: userData,
              isAuthenticated: true,
              isLoading: false,
            }, false, 'login/success');

            return true;
          } catch (error) {
            console.error('Login error:', error);
            return false;
          }
        },

        logout: () => {
          localStorage.removeItem('auth_token');
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
          }, false, 'logout');
        },
      }),
      {
        name: 'auth-storage',
        partialize: (state) => ({ 
          user: state.user,
          isAuthenticated: state.isAuthenticated 
        }),
      }
    ),
    { name: 'auth-store' }
  )
);

export const initializeAuth = () => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    useAuthStore.getState().validateToken(token);
  } else {
    useAuthStore.getState().setLoading(false);
  }
};