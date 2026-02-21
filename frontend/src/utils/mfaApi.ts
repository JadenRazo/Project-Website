import axios from 'axios';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export interface TOTPSetupResponse {
  secret: string;
  qr_code: string;
  url: string;
  backup_codes: string[];
}

export interface MFAStatusResponse {
  totp_enabled: boolean;
  totp_verified: boolean;
  backup_codes_remaining: number;
  enabled_at?: string;
}

export interface MFAVerificationRequest {
  temp_token: string;
  mfa_code: string;
  is_backup_code: boolean;
}

export interface LoginResponse {
  token?: string;
  user?: {
    id: string;
    email: string;
    is_admin: boolean;
    username: string;
  };
  expires_in?: number;
  requires_mfa: boolean;
  mfa_type?: string;
  temp_token?: string;
}

const getAuthHeader = () => {
  const token = localStorage.getItem('adminToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export const mfaApi = {
  setupTOTP: async (password: string): Promise<TOTPSetupResponse> => {
    const response = await axios.post(
      `${API_BASE}/api/v1/auth/admin/mfa/totp/setup`,
      { password },
      { headers: getAuthHeader() }
    );
    return response.data;
  },

  verifyTOTP: async (token: string): Promise<void> => {
    const response = await axios.post(
      `${API_BASE}/api/v1/auth/admin/mfa/totp/verify`,
      { token },
      { headers: getAuthHeader() }
    );
    return response.data;
  },

  disableTOTP: async (password: string, token: string): Promise<void> => {
    const response = await axios.post(
      `${API_BASE}/api/v1/auth/admin/mfa/totp/disable`,
      { password, token },
      { headers: getAuthHeader() }
    );
    return response.data;
  },

  getMFAStatus: async (): Promise<MFAStatusResponse> => {
    const response = await axios.get(
      `${API_BASE}/api/v1/auth/admin/mfa/status`,
      { headers: getAuthHeader() }
    );
    return response.data;
  },

  regenerateBackupCodes: async (password: string): Promise<string[]> => {
    const response = await axios.post(
      `${API_BASE}/api/v1/auth/admin/mfa/backup/regenerate`,
      { password },
      { headers: getAuthHeader() }
    );
    return response.data.backup_codes;
  },

  verifyMFAToken: async (
    tempToken: string,
    code: string,
    isBackupCode: boolean = false
  ): Promise<LoginResponse> => {
    const response = await axios.post(
      `${API_BASE}/api/v1/auth/admin/mfa/verify`,
      {
        temp_token: tempToken,
        mfa_code: code,
        is_backup_code: isBackupCode,
      }
    );
    return response.data;
  },
};
