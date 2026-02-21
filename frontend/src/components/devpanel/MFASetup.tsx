import React, { useState } from 'react';
import { mfaApi, TOTPSetupResponse } from '../../utils/mfaApi';

interface MFASetupProps {
  onComplete: () => void;
  onCancel: () => void;
}

export const MFASetup: React.FC<MFASetupProps> = ({ onComplete, onCancel }) => {
  const [step, setStep] = useState<'password' | 'qr' | 'verify' | 'complete'>('password');
  const [password, setPassword] = useState('');
  const [setupData, setSetupData] = useState<TOTPSetupResponse | null>(null);
  const [verificationCode, setVerificationCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [copiedBackupCodes, setCopiedBackupCodes] = useState(false);
  const [downloadedBackupCodes, setDownloadedBackupCodes] = useState(false);

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const data = await mfaApi.setupTOTP(password);
      setSetupData(data);
      setStep('qr');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to initiate MFA setup');
    } finally {
      setLoading(false);
    }
  };

  const handleCopyBackupCodes = () => {
    if (setupData) {
      const codesText = setupData.backup_codes.join('\n');
      navigator.clipboard.writeText(codesText);
      setCopiedBackupCodes(true);
      setTimeout(() => setCopiedBackupCodes(false), 3000);
    }
  };

  const handleDownloadBackupCodes = () => {
    if (setupData) {
      const codesText = setupData.backup_codes.join('\n');
      const blob = new Blob([codesText], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'mfa-backup-codes.txt';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      setDownloadedBackupCodes(true);
    }
  };

  const handleVerifyCode = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await mfaApi.verifyTOTP(verificationCode);
      setStep('complete');
      setTimeout(() => onComplete(), 2000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid verification code');
    } finally {
      setLoading(false);
    }
  };

  if (step === 'password') {
    return (
      <div className="mfa-setup-modal">
        <div className="mfa-setup-content">
          <h2>Enable Two-Factor Authentication</h2>
          <p>Enter your password to continue</p>

          <form onSubmit={handlePasswordSubmit}>
            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={loading}
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <div className="button-group">
              <button type="button" onClick={onCancel} disabled={loading}>
                Cancel
              </button>
              <button type="submit" disabled={loading}>
                {loading ? 'Processing...' : 'Continue'}
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  }

  if (step === 'qr') {
    return (
      <div className="mfa-setup-modal">
        <div className="mfa-setup-content">
          <h2>Scan QR Code</h2>
          <p>Use an authenticator app like Google Authenticator, Authy, 1Password, or Tofu to scan this QR code:</p>

          {setupData && (
            <>
              <div className="qr-code-container">
                <img src={setupData.qr_code} alt="TOTP QR Code" />
              </div>

              <div className="manual-entry">
                <p>Can't scan? Enter this code manually:</p>
                <code>{setupData.secret}</code>
              </div>

              <div className="backup-codes-section">
                <h3>Save Your Backup Codes</h3>
                <p>These codes can be used to access your account if you lose your authenticator device. Save them securely!</p>

                <div className="backup-codes-list">
                  {setupData.backup_codes.map((code, index) => (
                    <div key={index} className="backup-code">{code}</div>
                  ))}
                </div>

                <div className="button-group">
                  <button
                    type="button"
                    onClick={handleCopyBackupCodes}
                    className="secondary"
                  >
                    {copiedBackupCodes ? 'Copied!' : 'Copy Codes'}
                  </button>
                  <button
                    type="button"
                    onClick={handleDownloadBackupCodes}
                    className="secondary"
                  >
                    {downloadedBackupCodes ? 'Downloaded!' : 'Download Codes'}
                  </button>
                </div>
              </div>

              <button
                type="button"
                onClick={() => setStep('verify')}
                className="primary"
              >
                I've Saved My Codes - Continue
              </button>
            </>
          )}
        </div>
      </div>
    );
  }

  if (step === 'verify') {
    return (
      <div className="mfa-setup-modal">
        <div className="mfa-setup-content">
          <h2>Verify Your Setup</h2>
          <p>Enter the 6-digit code from your authenticator app to complete setup:</p>

          <form onSubmit={handleVerifyCode}>
            <div className="form-group">
              <label htmlFor="verification-code">Verification Code</label>
              <input
                type="text"
                id="verification-code"
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                maxLength={6}
                required
                disabled={loading}
                pattern="[0-9]{6}"
                autoComplete="off"
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <div className="button-group">
              <button type="button" onClick={() => setStep('qr')} disabled={loading}>
                Back
              </button>
              <button type="submit" disabled={loading || verificationCode.length !== 6}>
                {loading ? 'Verifying...' : 'Verify and Enable'}
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  }

  if (step === 'complete') {
    return (
      <div className="mfa-setup-modal">
        <div className="mfa-setup-content success">
          <h2>MFA Enabled Successfully!</h2>
          <p>Your account is now protected with two-factor authentication.</p>
          <div className="success-icon">✓</div>
        </div>
      </div>
    );
  }

  return null;
};
