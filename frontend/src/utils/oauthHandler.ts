import { api } from './apiConfig';

interface OAuthCallbackResult {
  success: boolean;
  user?: any;
  error?: string;
}

export const handleOAuthCallback = async (
  onSuccess: (userData: any) => void,
  onError?: (error: string) => void
): Promise<OAuthCallbackResult> => {
  const hash = window.location.hash.substring(1);
  const params = new URLSearchParams(hash);
  const accessToken = params.get('access_token');
  const tokenType = params.get('token_type');
  const expiresIn = params.get('expires_in');
  const errorParam = params.get('error');

  if (errorParam) {
    const errorMessage = params.get('error_description') || 'OAuth authentication failed';
    console.error('OAuth error:', errorParam, errorMessage);

    if (onError) {
      onError(errorMessage);
    }

    window.history.replaceState(null, '', window.location.pathname);

    return {
      success: false,
      error: errorMessage
    };
  }

  if (!accessToken) {
    return {
      success: false,
      error: 'No access token in callback'
    };
  }

  try {
    api.setAuthToken(accessToken);

    const userData = await api.get('/api/v1/auth/admin/validate');

    onSuccess(userData);

    window.history.replaceState(null, '', window.location.pathname);

    return {
      success: true,
      user: userData
    };
  } catch (error: any) {
    console.error('OAuth validation failed:', error);

    api.clearAuthToken();

    const errorMessage = error?.message || 'Failed to validate OAuth token';

    if (onError) {
      onError(errorMessage);
    }

    window.history.replaceState(null, '', window.location.pathname);

    return {
      success: false,
      error: errorMessage
    };
  }
};

export const checkOAuthCallback = (): boolean => {
  return window.location.hash.includes('access_token') || window.location.hash.includes('error');
};
