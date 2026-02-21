import React, { useState, useEffect } from 'react';
import { mfaApi, LoginResponse } from '../../utils/mfaApi';

interface MFAVerificationProps {
  tempToken: string;
  onSuccess: (response: LoginResponse) => void;
  onCancel: () => void;
}

export const MFAVerification: React.FC<MFAVerificationProps> = ({
  tempToken,
  onSuccess,
  onCancel,
}) => {
  const [mfaCode, setMfaCode] = useState('');
  const [useBackupCode, setUseBackupCode] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [failedAttempts, setFailedAttempts] = useState(0);

  useEffect(() => {
    const codeInput = document.getElementById('mfa-code');
    if (codeInput) {
      codeInput.focus();
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await mfaApi.verifyMFAToken(tempToken, mfaCode, useBackupCode);
      onSuccess(response);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Verification failed';
      setError(errorMessage);
      setFailedAttempts((prev) => prev + 1);
      setMfaCode('');
    } finally {
      setLoading(false);
    }
  };

  const handleCodeChange = (value: string) => {
    if (useBackupCode) {
      const formatted = value.replace(/[^0-9-]/g, '').toUpperCase();
      setMfaCode(formatted);
    } else {
      const digits = value.replace(/\D/g, '').slice(0, 6);
      setMfaCode(digits);
    }
  };

  const isCodeComplete = useBackupCode
    ? mfaCode.length >= 10
    : mfaCode.length === 6;

  return (
    <div className="mfa-verification">
      <div className="mfa-verification-content">
        <h2>Two-Factor Authentication</h2>
        <p>
          {useBackupCode
            ? 'Enter one of your backup codes to continue'
            : 'Enter the 6-digit code from your authenticator app'}
        </p>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="mfa-code">
              {useBackupCode ? 'Backup Code' : 'Verification Code'}
            </label>
            <input
              type="text"
              id="mfa-code"
              value={mfaCode}
              onChange={(e) => handleCodeChange(e.target.value)}
              required
              disabled={loading}
              placeholder={useBackupCode ? '00-00-00-00-00-00' : '000000'}
              autoComplete="off"
              className={failedAttempts > 0 ? 'error-shake' : ''}
            />
          </div>

          {error && (
            <div className="error-message">
              {error}
              {failedAttempts >= 3 && (
                <div className="warning-text">
                  Too many failed attempts. Your account may be temporarily locked.
                </div>
              )}
            </div>
          )}

          {failedAttempts > 0 && !useBackupCode && (
            <div className="info-message">
              Having trouble? You can use a backup code instead.
            </div>
          )}

          <div className="button-group">
            <button
              type="button"
              onClick={() => {
                setUseBackupCode(!useBackupCode);
                setMfaCode('');
                setError('');
              }}
              className="link-button"
              disabled={loading}
            >
              {useBackupCode ? 'Use Authenticator Code' : 'Use Backup Code'}
            </button>
          </div>

          <div className="button-group">
            <button
              type="button"
              onClick={onCancel}
              disabled={loading}
              className="secondary"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || !isCodeComplete}
              className="primary"
            >
              {loading ? 'Verifying...' : 'Verify'}
            </button>
          </div>
        </form>

        {failedAttempts >= 5 && (
          <div className="lockout-warning">
            Your account will be temporarily locked after 10 failed attempts.
            Remaining attempts: {10 - failedAttempts}
          </div>
        )}
      </div>
    </div>
  );
};
