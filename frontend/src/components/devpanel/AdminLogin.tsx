import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { api } from '../../utils/apiConfig';
import { handleError } from '../../utils/errorHandler';
import { useScrollTo } from '../../hooks/useScrollTo';

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
  min-height: 100vh;
  padding: 2rem;
  padding-top: calc(80px + 2rem);
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding-top: calc(70px + 1.5rem);
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding-top: calc(60px + 1rem);
    padding-left: 1rem;
    padding-right: 1rem;
  }
`;

const LoginCard = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  border: 1px solid ${({ theme }) => theme.colors.border};
  max-width: 400px;
  width: 100%;
  animation: fadeInUp 0.4s ease-out;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1.5rem;
    max-width: 100%;
  }
  
  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
`;

const LoginTitle = styled.h2`
  text-align: center;
  margin-bottom: 2rem;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.8rem;
  font-weight: 600;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -0.5rem;
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 3px;
    background: ${({ theme }) => theme.colors.primary};
    border-radius: 2px;
  }
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
`;

const InputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  animation: fadeIn 0.3s ease-out;
  animation-fill-mode: both;
  
  &:nth-child(1) { animation-delay: 0.1s; }
  &:nth-child(2) { animation-delay: 0.2s; }
  &:nth-child(3) { animation-delay: 0.3s; }
  &:nth-child(4) { animation-delay: 0.4s; }
  
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateX(-10px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
`;

const Label = styled.label`
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  opacity: 0.9;
`;

const Input = styled.input`
  padding: 0.875rem 1rem;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  font-size: 1rem;
  transition: all 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
    box-shadow: 0 0 0 3px ${({ theme }) => theme.colors.primary}20;
    transform: translateY(-1px);
  }
  
  &:hover:not(:disabled) {
    border-color: ${({ theme }) => theme.colors.borderHover};
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    background: ${({ theme }) => theme.colors.surfaceDisabled};
  }
  
  &::placeholder {
    color: ${({ theme }) => theme.colors.textSecondary};
    opacity: 0.7;
  }
`;

const Button = styled.button<{ variant?: 'primary' | 'secondary' }>`
  padding: 0.875rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 0.5rem;
  position: relative;
  overflow: hidden;
  
  background: ${({ theme, variant }) => 
    variant === 'secondary' ? theme.colors.background : theme.colors.primary};
  color: ${({ theme, variant }) => 
    variant === 'secondary' ? theme.colors.text : 'white'};
  border: ${({ theme, variant }) => 
    variant === 'secondary' ? `2px solid ${theme.colors.border}` : 'none'};
  
  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    background: ${({ theme, variant }) => 
      variant === 'secondary' 
        ? theme.colors.primary + '15'
        : 'rgba(255, 255, 255, 0.3)'};
    transform: translate(-50%, -50%);
    transition: width 0.6s cubic-bezier(0.165, 0.84, 0.44, 1), 
                height 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
  }
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px ${({ theme, variant }) => 
      variant === 'secondary' ? 'rgba(0, 0, 0, 0.1)' : theme.colors.primary + '40'};
    
    &::before {
      width: 400%;
      height: 400%;
    }
  }
  
  &:active:not(:disabled) {
    transform: scale(0.95);
    transition: transform 0.1s ease;
    animation: buttonPop 0.3s ease-out;
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }
  
  @keyframes buttonPop {
    0% {
      transform: scale(0.95);
    }
    40% {
      transform: scale(1.02);
    }
    100% {
      transform: scale(1);
    }
  }
  
  @media (hover: none) {
    &:active:not(:disabled) {
      &::before {
        width: 400%;
        height: 400%;
        transition: width 0.3s ease-out, height 0.3s ease-out;
      }
    }
  }
`;

const ErrorMessage = styled.div`
  padding: 0.875rem 1rem;
  background: ${({ theme }) => theme.colors.error}20;
  color: ${({ theme }) => theme.colors.error};
  border-radius: 8px;
  border-left: 4px solid ${({ theme }) => theme.colors.error};
  font-size: 0.9rem;
  margin-top: 0.5rem;
  animation: shake 0.4s ease-out;
  
  @keyframes shake {
    0%, 100% { transform: translateX(0); }
    10%, 30%, 50%, 70%, 90% { transform: translateX(-2px); }
    20%, 40%, 60%, 80% { transform: translateX(2px); }
  }
`;

const SuccessMessage = styled.div`
  padding: 0.875rem 1rem;
  background: ${({ theme }) => theme.colors.success}20;
  color: ${({ theme }) => theme.colors.success};
  border-radius: 8px;
  border-left: 4px solid ${({ theme }) => theme.colors.success};
  font-size: 0.9rem;
  margin-top: 0.5rem;
  animation: slideInRight 0.4s ease-out;
  
  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
`;

const SetupSection = styled.div`
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  animation: slideIn 0.4s ease-out;
  
  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
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

const LoadingContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 3rem;
  
  div {
    width: 48px;
    height: 48px;
    border: 3px solid ${({ theme }) => theme.colors.border};
    border-top-color: ${({ theme }) => theme.colors.primary};
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  p {
    color: ${({ theme }) => theme.colors.textSecondary};
    font-size: 0.9rem;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
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
  
  // Scroll hooks and refs
  const { scrollToElement } = useScrollTo();
  const errorRef = useRef<HTMLDivElement>(null);
  const setupErrorRef = useRef<HTMLDivElement>(null);
  const setupMessageRef = useRef<HTMLDivElement>(null);
  const setupSectionRef = useRef<HTMLDivElement>(null);
  const loginCardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    checkSetupStatus();
  }, []);

  // Scroll to login card on mount
  useEffect(() => {
    if (loginCardRef.current) {
      setTimeout(() => {
        scrollToElement(loginCardRef.current, { 
          behavior: 'smooth',
          offset: 100
        });
      }, 100);
    }
  }, [scrollToElement]);

  // Scroll to error messages
  useEffect(() => {
    if (error && errorRef.current) {
      scrollToElement(errorRef.current, { 
        behavior: 'smooth',
        offset: 80
      });
    }
  }, [error, scrollToElement]);

  useEffect(() => {
    if (setupError && setupErrorRef.current) {
      scrollToElement(setupErrorRef.current, { 
        behavior: 'smooth',
        offset: 80
      });
    }
  }, [setupError, scrollToElement]);

  // Scroll to success message
  useEffect(() => {
    if (setupMessage && setupMessageRef.current) {
      scrollToElement(setupMessageRef.current, { 
        behavior: 'smooth',
        offset: 80
      });
    }
  }, [setupMessage, scrollToElement]);

  // Scroll to setup section when shown
  useEffect(() => {
    if (showSetup && setupSectionRef.current) {
      setTimeout(() => {
        scrollToElement(setupSectionRef.current, { 
          behavior: 'smooth',
          offset: 80
        });
      }, 100);
    }
  }, [showSetup, scrollToElement]);

  const checkSetupStatus = async () => {
    try {
      const data = await api.get<SetupStatus>('/api/v1/auth/admin/setup/status', {
        skipAuth: true
      });
      setSetupStatus(data);
      if (!data.hasAdmin) {
        setShowSetup(true);
      }
    } catch (error) {
      handleError(error, { context: 'AdminLogin.checkSetupStatus' });
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const data = await api.post('/api/v1/auth/admin/login', 
        { email, password }, 
        { skipAuth: true }
      );
      
      api.setAuthToken(data.token);
      onLoginSuccess(data.user);
    } catch (err) {
      handleError(err, { context: 'AdminLogin.handleLogin' });
      const message = err instanceof Error && err.message.includes('error') 
        ? err.message 
        : 'Network error. Please try again.';
      setError(message);
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
      await api.post('/api/v1/auth/admin/setup/request', 
        { email: setupEmail },
        { skipAuth: true }
      );
      
      setSetupMessage('Setup email sent! Check your email for the setup token.');
      setSetupStep('complete');
    } catch (err) {
      handleError(err, { context: 'AdminLogin.handleSetupRequest' });
      const message = err instanceof Error && err.message 
        ? err.message 
        : 'Network error. Please try again.';
      setSetupError(message);
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
      await api.post('/api/v1/auth/admin/setup/complete', 
        {
          email: setupEmail,
          password: setupPassword,
          confirm_password: confirmPassword,
          setup_token: setupToken,
        },
        { skipAuth: true }
      );
      
      setSetupMessage('Admin account created successfully! You can now log in.');
      setShowSetup(false);
      checkSetupStatus();
    } catch (err) {
      handleError(err, { context: 'AdminLogin.handleSetupComplete' });
      const message = err instanceof Error && err.message 
        ? err.message 
        : 'Network error. Please try again.';
      setSetupError(message);
    } finally {
      setSetupLoading(false);
    }
  };

  return (
    <LoginContainer>
      <LoginCard ref={loginCardRef}>
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
              
              {error && <ErrorMessage ref={errorRef}>{error}</ErrorMessage>}
              
              <Button type="submit" disabled={loading}>
                {loading ? 'Signing in...' : 'Sign In'}
              </Button>
            </Form>
          </>
        ) : showSetup ? (
          <SetupSection ref={setupSectionRef}>
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
                
                {setupError && <ErrorMessage ref={setupErrorRef}>{setupError}</ErrorMessage>}
                {setupMessage && <SuccessMessage ref={setupMessageRef}>{setupMessage}</SuccessMessage>}
                
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
                
                {setupError && <ErrorMessage ref={setupErrorRef}>{setupError}</ErrorMessage>}
                {setupMessage && <SuccessMessage ref={setupMessageRef}>{setupMessage}</SuccessMessage>}
                
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
          <LoadingContainer>
            <div />
            <p>Checking authentication status...</p>
          </LoadingContainer>
        )}
      </LoginCard>
    </LoginContainer>
  );
};

export default AdminLogin;
