import React, { useState, useEffect } from 'react';
import { mfaApi, MFAStatusResponse } from '../../utils/mfaApi';
import { MFASetup } from './MFASetup';

interface DisableModalState {
  show: boolean;
  password: string;
  token: string;
  error: string;
  loading: boolean;
}

interface RegenerateModalState {
  show: boolean;
  password: string;
  codes: string[] | null;
  error: string;
  loading: boolean;
}

export const MFAManagement: React.FC = () => {
  const [status, setStatus] = useState<MFAStatusResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [showSetup, setShowSetup] = useState(false);
  const [disableModal, setDisableModal] = useState<DisableModalState>({
    show: false,
    password: '',
    token: '',
    error: '',
    loading: false,
  });
  const [regenerateModal, setRegenerateModal] = useState<RegenerateModalState>({
    show: false,
    password: '',
    codes: null,
    error: '',
    loading: false,
  });

  const loadStatus = async () => {
    setLoading(true);
    try {
      const data = await mfaApi.getMFAStatus();
      setStatus(data);
    } catch (err) {
      console.error('Failed to load MFA status:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadStatus();
  }, []);

  const handleDisableMFA = async (e: React.FormEvent) => {
    e.preventDefault();
    setDisableModal((prev) => ({ ...prev, error: '', loading: true }));

    try {
      await mfaApi.disableTOTP(disableModal.password, disableModal.token);
      await loadStatus();
      setDisableModal({
        show: false,
        password: '',
        token: '',
        error: '',
        loading: false,
      });
    } catch (err: any) {
      setDisableModal((prev) => ({
        ...prev,
        error: err.response?.data?.error || 'Failed to disable MFA',
        loading: false,
      }));
    }
  };

  const handleRegenerateBackupCodes = async (e: React.FormEvent) => {
    e.preventDefault();
    setRegenerateModal((prev) => ({ ...prev, error: '', loading: true }));

    try {
      const codes = await mfaApi.regenerateBackupCodes(regenerateModal.password);
      setRegenerateModal((prev) => ({
        ...prev,
        codes,
        error: '',
        loading: false,
      }));
      await loadStatus();
    } catch (err: any) {
      setRegenerateModal((prev) => ({
        ...prev,
        error: err.response?.data?.error || 'Failed to regenerate backup codes',
        loading: false,
      }));
    }
  };

  const handleDownloadBackupCodes = () => {
    if (regenerateModal.codes) {
      const codesText = regenerateModal.codes.join('\n');
      const blob = new Blob([codesText], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'mfa-backup-codes-new.txt';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }
  };

  if (loading) {
    return <div className="mfa-management loading">Loading MFA settings...</div>;
  }

  if (showSetup) {
    return (
      <MFASetup
        onComplete={() => {
          setShowSetup(false);
          loadStatus();
        }}
        onCancel={() => setShowSetup(false)}
      />
    );
  }

  return (
    <div className="mfa-management">
      <h2>Two-Factor Authentication</h2>

      <div className="mfa-status-card">
        <div className="status-header">
          <h3>TOTP Status</h3>
          <span className={`status-badge ${status?.totp_enabled ? 'enabled' : 'disabled'}`}>
            {status?.totp_enabled ? 'Enabled' : 'Disabled'}
          </span>
        </div>

        {status?.totp_enabled && (
          <div className="mfa-details">
            <p>
              <strong>Verified:</strong> {status.totp_verified ? 'Yes' : 'No'}
            </p>
            <p>
              <strong>Backup Codes Remaining:</strong> {status.backup_codes_remaining}
            </p>
            {status.enabled_at && (
              <p>
                <strong>Enabled At:</strong>{' '}
                {new Date(status.enabled_at).toLocaleString()}
              </p>
            )}
          </div>
        )}

        <div className="action-buttons">
          {status?.totp_enabled ? (
            <>
              <button
                onClick={() => setRegenerateModal({ ...regenerateModal, show: true })}
                className="secondary"
              >
                Regenerate Backup Codes
              </button>
              <button
                onClick={() => setDisableModal({ ...disableModal, show: true })}
                className="danger"
              >
                Disable MFA
              </button>
            </>
          ) : (
            <button onClick={() => setShowSetup(true)} className="primary">
              Enable MFA
            </button>
          )}
        </div>
      </div>

      {disableModal.show && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>Disable Two-Factor Authentication</h3>
            <p>Enter your password and current TOTP code to disable MFA:</p>

            <form onSubmit={handleDisableMFA}>
              <div className="form-group">
                <label htmlFor="disable-password">Password</label>
                <input
                  type="password"
                  id="disable-password"
                  value={disableModal.password}
                  onChange={(e) =>
                    setDisableModal({ ...disableModal, password: e.target.value })
                  }
                  required
                  disabled={disableModal.loading}
                />
              </div>

              <div className="form-group">
                <label htmlFor="disable-token">TOTP Code</label>
                <input
                  type="text"
                  id="disable-token"
                  value={disableModal.token}
                  onChange={(e) =>
                    setDisableModal({
                      ...disableModal,
                      token: e.target.value.replace(/\D/g, '').slice(0, 6),
                    })
                  }
                  maxLength={6}
                  required
                  disabled={disableModal.loading}
                />
              </div>

              {disableModal.error && (
                <div className="error-message">{disableModal.error}</div>
              )}

              <div className="button-group">
                <button
                  type="button"
                  onClick={() =>
                    setDisableModal({
                      show: false,
                      password: '',
                      token: '',
                      error: '',
                      loading: false,
                    })
                  }
                  disabled={disableModal.loading}
                >
                  Cancel
                </button>
                <button type="submit" disabled={disableModal.loading} className="danger">
                  {disableModal.loading ? 'Disabling...' : 'Disable MFA'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {regenerateModal.show && (
        <div className="modal-overlay">
          <div className="modal-content">
            <h3>Regenerate Backup Codes</h3>

            {!regenerateModal.codes ? (
              <>
                <p>This will invalidate your existing backup codes. Enter your password to continue:</p>

                <form onSubmit={handleRegenerateBackupCodes}>
                  <div className="form-group">
                    <label htmlFor="regenerate-password">Password</label>
                    <input
                      type="password"
                      id="regenerate-password"
                      value={regenerateModal.password}
                      onChange={(e) =>
                        setRegenerateModal({ ...regenerateModal, password: e.target.value })
                      }
                      required
                      disabled={regenerateModal.loading}
                    />
                  </div>

                  {regenerateModal.error && (
                    <div className="error-message">{regenerateModal.error}</div>
                  )}

                  <div className="button-group">
                    <button
                      type="button"
                      onClick={() =>
                        setRegenerateModal({
                          show: false,
                          password: '',
                          codes: null,
                          error: '',
                          loading: false,
                        })
                      }
                      disabled={regenerateModal.loading}
                    >
                      Cancel
                    </button>
                    <button type="submit" disabled={regenerateModal.loading}>
                      {regenerateModal.loading ? 'Generating...' : 'Regenerate'}
                    </button>
                  </div>
                </form>
              </>
            ) : (
              <>
                <p>Your new backup codes:</p>

                <div className="backup-codes-list">
                  {regenerateModal.codes.map((code, index) => (
                    <div key={index} className="backup-code">
                      {code}
                    </div>
                  ))}
                </div>

                <div className="button-group">
                  <button onClick={handleDownloadBackupCodes} className="secondary">
                    Download Codes
                  </button>
                  <button
                    onClick={() =>
                      setRegenerateModal({
                        show: false,
                        password: '',
                        codes: null,
                        error: '',
                        loading: false,
                      })
                    }
                    className="primary"
                  >
                    Done
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
