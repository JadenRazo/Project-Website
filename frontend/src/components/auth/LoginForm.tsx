import React, { useState, useRef, useEffect } from 'react';
import styled, { keyframes } from 'styled-components';
import { useAuth } from '../../hooks/useAuth';
import { validateEmail } from '../../utils/validation';
import { useScrollTo } from '../../hooks/useScrollTo';

interface LoginFormProps {
  onSuccess: () => void;
}

const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
`;

const shake = keyframes`
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-2px); }
  20%, 40%, 60%, 80% { transform: translateX(2px); }
`;

const buttonPop = keyframes`
  0% {
    transform: scale(0.95);
  }
  40% {
    transform: scale(1.02);
  }
  100% {
    transform: scale(1);
  }
`;

const spin = keyframes`
  to { transform: rotate(360deg); }
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
  animation: ${fadeIn} 0.3s ease-out;
  animation-fill-mode: both;
  
  &:nth-child(1) { animation-delay: 0.1s; }
  &:nth-child(2) { animation-delay: 0.2s; }
  &:nth-child(3) { animation-delay: 0.3s; }
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
    border-color: ${({ theme }) => theme.colors.borderHover || theme.colors.border};
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    background: ${({ theme }) => theme.colors.surfaceLight};
  }
  
  &::placeholder {
    color: ${({ theme }) => theme.colors.textSecondary};
    opacity: 0.7;
  }
`;

const Button = styled.button`
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
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  
  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    background: rgba(255, 255, 255, 0.3);
    transform: translate(-50%, -50%);
    transition: width 0.6s cubic-bezier(0.165, 0.84, 0.44, 1), 
                height 0.6s cubic-bezier(0.165, 0.84, 0.44, 1);
  }
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px ${({ theme }) => theme.colors.primary}40;
    
    &::before {
      width: 400%;
      height: 400%;
    }
  }
  
  &:active:not(:disabled) {
    transform: scale(0.95);
    transition: transform 0.1s ease;
    animation: ${buttonPop} 0.3s ease-out;
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
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
  background: ${({ theme }) => theme.colors.error || '#dc3545'}20;
  color: ${({ theme }) => theme.colors.error || '#dc3545'};
  border-radius: 8px;
  border-left: 4px solid ${({ theme }) => theme.colors.error || '#dc3545'};
  font-size: 0.9rem;
  margin-top: 0.5rem;
  animation: ${shake} 0.4s ease-out;
`;

const LoadingSpinner = styled.div`
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: ${spin} 0.8s linear infinite;
  margin-right: 0.5rem;
`;

const ButtonContent = styled.span`
  display: inline-flex;
  align-items: center;
  justify-content: center;
`;

const LoginForm: React.FC<LoginFormProps> = ({ onSuccess }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const { login } = useAuth();
  const { scrollToElement } = useScrollTo();
  const errorRef = useRef<HTMLDivElement>(null);

  // Scroll to error message when it appears
  useEffect(() => {
    if (error && errorRef.current) {
      scrollToElement(errorRef.current, { 
        behavior: 'smooth',
        offset: 80
      });
    }
  }, [error, scrollToElement]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Client-side validation
    if (!validateEmail(email)) {
      setError('Please enter a valid email address');
      return;
    }

    if (password.length < 6) {
      setError('Password must be at least 6 characters');
      return;
    }

    setLoading(true);

    try {
      await login(email, password);
      onSuccess();
    } catch (err: any) {
      if (err.message?.includes('invalid_email')) {
        setError('Email address contains invalid characters');
      } else {
        setError(err.message || 'Invalid email or password');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      <InputGroup>
        <Label htmlFor="login-email">Email</Label>
        <Input
          id="login-email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your@email.com"
          required
          disabled={loading}
          autoComplete="email"
        />
      </InputGroup>

      <InputGroup>
        <Label htmlFor="login-password">Password</Label>
        <Input
          id="login-password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter your password"
          required
          disabled={loading}
          autoComplete="current-password"
        />
      </InputGroup>

      {error && <ErrorMessage ref={errorRef}>{error}</ErrorMessage>}

      <Button type="submit" disabled={loading}>
        <ButtonContent>
          {loading && <LoadingSpinner />}
          {loading ? 'Signing in...' : 'Sign In'}
        </ButtonContent>
      </Button>
    </Form>
  );
};

export default LoginForm;