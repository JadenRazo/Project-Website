import React from 'react';
import styled from 'styled-components';
import { FaGithub, FaWindows } from 'react-icons/fa';
import { FcGoogle } from 'react-icons/fc';

interface OAuthButtonsProps {
  disabled?: boolean;
}

const OAuthButtons: React.FC<OAuthButtonsProps> = ({ disabled = false }) => {
  const handleOAuthLogin = (provider: 'google' | 'github' | 'microsoft') => {
    const redirectUrl = encodeURIComponent('/devpanel');
    window.location.href = `/api/v1/auth/admin/oauth/login/${provider}?redirect=${redirectUrl}`;
  };

  return (
    <>
      <Divider>
        <DividerLine />
        <DividerText>OR</DividerText>
        <DividerLine />
      </Divider>

      <OAuthSection>
        <OAuthButton
          type="button"
          onClick={() => handleOAuthLogin('google')}
          disabled={disabled}
          $provider="google"
        >
          <FcGoogle size={18} />
          <span>Continue with Google</span>
        </OAuthButton>

        <OAuthButton
          type="button"
          onClick={() => handleOAuthLogin('github')}
          disabled={disabled}
          $provider="github"
        >
          <FaGithub size={18} />
          <span>Continue with GitHub</span>
        </OAuthButton>

        <OAuthButton
          type="button"
          onClick={() => handleOAuthLogin('microsoft')}
          disabled={disabled}
          $provider="microsoft"
        >
          <FaWindows size={18} />
          <span>Continue with Microsoft</span>
        </OAuthButton>
      </OAuthSection>
    </>
  );
};

const Divider = styled.div`
  display: flex;
  align-items: center;
  gap: 1rem;
  margin: 1.5rem 0;
  animation: fadeIn 0.5s ease-out;
  animation-delay: 0.5s;
  animation-fill-mode: both;

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
`;

const DividerLine = styled.div`
  flex: 1;
  height: 1px;
  background: ${({ theme }) => theme.colors.border};
`;

const DividerText = styled.span`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
`;

const OAuthSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  animation: fadeIn 0.5s ease-out;
  animation-delay: 0.6s;
  animation-fill-mode: both;
`;

const OAuthButton = styled.button<{ $provider: string }>`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 0.875rem 1.5rem;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;

  svg {
    color: ${({ $provider, theme }) => {
      if ($provider === 'google') return '#EA4335';
      if ($provider === 'github') return theme.colors.text;
      if ($provider === 'microsoft') return '#00A4EF';
      return theme.colors.text;
    }};
    transition: transform 0.2s ease;
  }

  span {
    flex: 1;
    text-align: center;
  }

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    border-color: ${({ theme }) => theme.colors.borderHover || theme.colors.primary};

    svg {
      transform: scale(1.1);
    }
  }

  &:active:not(:disabled) {
    transform: scale(0.98) translateY(0);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;

    svg {
      transform: none;
    }
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.75rem 1.25rem;
    font-size: 0.9rem;

    svg {
      width: 16px;
      height: 16px;
    }
  }
`;

export default OAuthButtons;
