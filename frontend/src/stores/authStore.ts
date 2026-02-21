import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import type { AuthState, User } from './types';

interface AuthActions {
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  validateToken: (token: string) => Promise<void>;
  setLoading: (loading: boolean) => void;
  setUser: (user: User | null) => void;
  setAuthenticated: (authenticated: boolean) => void;
  setAuthModalOpen: (open: boolean) => void;
  setAuthModalMode: (mode: 'login' | 'register') => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      (set, get) => ({
        user: null,
        isAuthenticated: false,
        isLoading: true,
        authModalOpen: false,
        authModalMode: 'login',

        setLoading: (loading: boolean) =>
          set({ isLoading: loading }, false, 'setLoading'),

        setUser: (user: User | null) =>
          set({ user }, false, 'setUser'),

        setAuthenticated: (authenticated: boolean) =>
          set({ isAuthenticated: authenticated }, false, 'setAuthenticated'),

        setAuthModalOpen: (open: boolean) =>
          set({ authModalOpen: open }, false, 'setAuthModalOpen'),

        setAuthModalMode: (mode: 'login' | 'register') =>
          set({ authModalMode: mode }, false, 'setAuthModalMode'),

        validateToken: async (token: string) => {
          try {
            const response = await fetch('/api/v1/auth/admin/validate', {
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
            const response = await fetch('/api/auth/login', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ email_or_username: email, password }),
            });

            if (!response.ok) {
              const errorData = await response.json().catch(() => ({}));
              throw new Error(errorData.error || 'Invalid email or password');
            }

            const userData = await response.json();
            localStorage.setItem('auth_token', userData.token);
            
            set({
              user: userData,
              isAuthenticated: true,
              isLoading: false,
              authModalOpen: false,
            }, false, 'login/success');
          } catch (error) {
            throw error;
          }
        },

        register: async (email: string, password: string) => {
          try {
            const response = await fetch('/api/auth/register', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
              const errorData = await response.json().catch(() => ({}));
              throw new Error(errorData.error || 'Failed to create account');
            }

            const userData = await response.json();
            localStorage.setItem('auth_token', userData.token);
            
            set({
              user: userData,
              isAuthenticated: true,
              isLoading: false,
              authModalOpen: false,
            }, false, 'register/success');
          } catch (error) {
            throw error;
          }
        },

        logout: () => {
          localStorage.removeItem('auth_token');
          localStorage.removeItem('auth-storage');
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