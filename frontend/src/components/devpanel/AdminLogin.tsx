import React, { useState, useEffect } from 'react';
import styled from 'styled-components';

interface AdminLoginProps {
  onLoginSuccess: (userData: any) => void;
}

interface SetupStatus {
  hasAdmin: boolean;
}

const LoginContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 2rem;
`;

const LoginCard = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  border: 1px solid ${({ theme }) => theme.colors.border};
  max-width: 400px;
  width: 100%;
`;

const LoginTitle = styled.h2`
  text-align: center;
  margin-bottom: 2rem;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.8rem;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const InputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const Label = styled.label`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.9rem;
`;

const Input = styled.input`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 6px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  font-size: 1rem;
  transition: border-color 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
`;

const Button = styled.button<{ variant?: 'primary' | 'secondary' }>`
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, variant }) => 
    variant === 'secondary' ? theme.colors.background : theme.colors.primary};
  color: ${({ theme, variant }) => 
    variant === 'secondary' ? theme.colors.text : 'white'};
  border: ${({ theme, variant }) => 
    variant === 'secondary' ? `1px solid ${theme.colors.border}` : 'none'};
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }
`;

const ErrorMessage = styled.div`
  padding: 0.75rem;
  background: ${({ theme }) => theme.colors.error}20;
  color: ${({ theme }) => theme.colors.error};
  border-radius: 6px;
  border-left: 4px solid ${({ theme }) => theme.colors.error};
  font-size: 0.9rem;
`;

const SuccessMessage = styled.div`
  padding: 0.75rem;
  background: ${({ theme }) => theme.colors.success}20;
  color: ${({ theme }) => theme.colors.success};
  border-radius: 6px;
  border-left: 4px solid ${({ theme }) => theme.colors.success};
  font-size: 0.9rem;
`;

const SetupSection = styled.div`
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border};
`;

const SetupTitle = styled.h3`
  margin-bottom: 1rem;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.2rem;
`;

const SetupDescription = styled.p`
  margin-bottom: 1rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  line-height: 1.4;
`;

const AdminLogin: React.FC<AdminLoginProps> = ({ onLoginSuccess }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [setupEmail, setSetupEmail] = useState('');
  const [setupPassword, setSetupPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [setupToken, setSetupToken] = useState('');
  const [loading, setLoading] = useState(false);
  const [setupLoading, setSetupLoading] = useState(false);
  const [error, setError] = useState('');
  const [setupError, setSetupError] = useState('');
  const [setupMessage, setSetupMessage] = useState('');
  const [setupStatus, setSetupStatus] = useState<SetupStatus | null>(null);
  const [showSetup, setShowSetup] = useState(false);
  const [setupStep, setSetupStep] = useState<'request' | 'complete'>('request');

  useEffect(() => {
    checkSetupStatus();
  }, []);

  const checkSetupStatus = async () => {
    try {
      const response = await fetch('/api/v1/auth/admin/setup/status');
      if (response.ok) {
        const data = await response.json();
        setSetupStatus(data);
        if (!data.hasAdmin) {
          setShowSetup(true);
        }
      }
    } catch (error) {
      console.error('Error checking setup status:', error);
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/v1/auth/admin/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (response.ok) {
        localStorage.setItem('auth_token', data.token);
        onLoginSuccess(data.user);
      } else {
        setError(data.error || 'Login failed');
      }
    } catch (err) {
      setError('Network error. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleSetupRequest = async (e: React.FormEvent) => {
    e.preventDefault();
    setSetupLoading(true);
    setSetupError('');
    setSetupMessage('');

    try {
      const response = await fetch('/api/v1/auth/admin/setup/request', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: setupEmail }),
      });

      const data = await response.json();

      if (response.ok) {
        setSetupMessage('Setup email sent! Check your email for the setup token.');
        setSetupStep('complete');
      } else {
        setSetupError(data.error || 'Setup request failed');
      }
    } catch (err) {
      setSetupError('Network error. Please try again.');
    } finally {
      setSetupLoading(false);
    }
  };

  const handleSetupComplete = async (e: React.FormEvent) => {
    e.preventDefault();
    setSetupLoading(true);
    setSetupError('');
    setSetupMessage('');

    if (setupPassword !== confirmPassword) {
      setSetupError('Passwords do not match');
      setSetupLoading(false);
      return;
    }

    try {
      const response = await fetch('/api/v1/auth/admin/setup/complete', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: setupEmail,
          password: setupPassword,
          confirm_password: confirmPassword,
          setup_token: setupToken,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        setSetupMessage('Admin account created successfully! You can now log in.');
        setShowSetup(false);
        checkSetupStatus();
      } else {
        setSetupError(data.error || 'Setup completion failed');
      }
    } catch (err) {
      setSetupError('Network error. Please try again.');
    } finally {
      setSetupLoading(false);
    }
  };

  return (
    <LoginContainer>
      <LoginCard>
        <LoginTitle>DevPanel Admin Access</LoginTitle>
        
        {setupStatus?.hasAdmin ? (
          <>
            <Form onSubmit={handleLogin}>
              <InputGroup>
                <Label htmlFor="email">Admin Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                  required
                  placeholder="support@jadenrazo.dev"
                />
              </InputGroup>
              
              <InputGroup>
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={loading}
                  required
                  placeholder="Enter your password"
                />
              </InputGroup>
              
              {error && <ErrorMessage>{error}</ErrorMessage>}
              
              <Button type="submit" disabled={loading}>
                {loading ? 'Signing in...' : 'Sign In'}
              </Button>
            </Form>
          </>
        ) : showSetup ? (
          <SetupSection>
            <SetupTitle>Admin Setup Required</SetupTitle>
            <SetupDescription>
              No admin account exists. Create one using an authorized email address (@jadenrazo.dev).
            </SetupDescription>
            
            {setupStep === 'request' ? (
              <Form onSubmit={handleSetupRequest}>
                <InputGroup>
                  <Label htmlFor="setupEmail">Admin Email</Label>
                  <Input
                    id="setupEmail"
                    type="email"
                    value={setupEmail}
                    onChange={(e) => setSetupEmail(e.target.value)}
                    disabled={setupLoading}
                    required
                    placeholder="support@jadenrazo.dev"
                  />
                </InputGroup>
                
                {setupError && <ErrorMessage>{setupError}</ErrorMessage>}
                {setupMessage && <SuccessMessage>{setupMessage}</SuccessMessage>}
                
                <Button type="submit" disabled={setupLoading}>
                  {setupLoading ? 'Sending...' : 'Request Setup'}
                </Button>
              </Form>
            ) : (
              <Form onSubmit={handleSetupComplete}>
                <InputGroup>
                  <Label htmlFor="setupToken">Setup Token</Label>
                  <Input
                    id="setupToken"
                    type="text"
                    value={setupToken}
                    onChange={(e) => setSetupToken(e.target.value)}
                    disabled={setupLoading}
                    required
                    placeholder="Token from email"
                  />
                </InputGroup>
                
                <InputGroup>
                  <Label htmlFor="setupPassword">Password</Label>
                  <Input
                    id="setupPassword"
                    type="password"
                    value={setupPassword}
                    onChange={(e) => setSetupPassword(e.target.value)}
                    disabled={setupLoading}
                    required
                    placeholder="Choose a strong password"
                    minLength={8}
                  />
                </InputGroup>
                
                <InputGroup>
                  <Label htmlFor="confirmPassword">Confirm Password</Label>
                  <Input
                    id="confirmPassword"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    disabled={setupLoading}
                    required
                    placeholder="Confirm your password"
                    minLength={8}
                  />
                </InputGroup>
                
                {setupError && <ErrorMessage>{setupError}</ErrorMessage>}
                {setupMessage && <SuccessMessage>{setupMessage}</SuccessMessage>}
                
                <Button type="submit" disabled={setupLoading}>
                  {setupLoading ? 'Creating Account...' : 'Complete Setup'}
                </Button>
                
                <Button 
                  type="button" 
                  variant="secondary" 
                  onClick={() => setSetupStep('request')}
                  disabled={setupLoading}
                >
                  Back to Request
                </Button>
              </Form>
            )}
          </SetupSection>
        ) : (
          <div>Loading...</div>
        )}
      </LoginCard>
    </LoginContainer>
  );
};

export default AdminLogin;
