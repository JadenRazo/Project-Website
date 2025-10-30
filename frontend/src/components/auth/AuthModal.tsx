import React, { useState } from 'react';
import styled, { keyframes } from 'styled-components';
import { X } from 'lucide-react';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import { ScrollableModal } from '../common/ScrollableModal';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialMode?: 'login' | 'register';
}

const fadeIn = keyframes`
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
`;

const slideUp = keyframes`
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
`;

const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  animation: ${slideUp} 0.3s ease-out;
`;

const Title = styled.h2`
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.textSecondary};
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.375rem;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background-color: ${({ theme }) => theme.colors.surfaceLight};
    color: ${({ theme }) => theme.colors.text};
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px ${({ theme }) => theme.colors.primary}40;
  }
`;

const ModeToggle = styled.div`
  text-align: center;
  margin-top: 1.5rem;
  animation: ${fadeIn} 0.5s ease-out;
`;

const ToggleText = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.875rem;
  margin: 0;
`;

const ToggleLink = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.primary};
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  padding: 0;
  margin-left: 0.25rem;
  text-decoration: underline;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 0.8;
  }

  &:focus {
    outline: none;
    opacity: 0.8;
  }
`;

const ModalContent = styled.div`
  max-width: 450px;
  margin: 0 auto;
  width: 100%;
`;

const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onClose, initialMode = 'login' }) => {
  const [mode, setMode] = useState<'login' | 'register'>(initialMode);

  const handleSuccess = () => {
    onClose();
  };

  const toggleMode = () => {
    setMode(mode === 'login' ? 'register' : 'login');
  };

  return (
    <ScrollableModal isOpen={isOpen} onClose={onClose}>
      <ModalContent>
        <ModalHeader>
        <Title>{mode === 'login' ? 'Sign In' : 'Create Account'}</Title>
        <CloseButton onClick={onClose} aria-label="Close modal">
          <X size={20} />
        </CloseButton>
      </ModalHeader>

      {mode === 'login' ? (
        <LoginForm onSuccess={handleSuccess} />
      ) : (
        <RegisterForm onSuccess={handleSuccess} />
      )}

      <ModeToggle>
        <ToggleText>
          {mode === 'login' ? "Don't have an account?" : 'Already have an account?'}
          <ToggleLink onClick={toggleMode}>
            {mode === 'login' ? 'Sign up' : 'Sign in'}
          </ToggleLink>
        </ToggleText>
      </ModeToggle>
      </ModalContent>
    </ScrollableModal>
  );
};

export default AuthModal;