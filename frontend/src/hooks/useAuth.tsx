import { useAuthStore } from '../stores';

export const useAuth = () => {
  const {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
  } = useAuthStore();

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
  };
}; 