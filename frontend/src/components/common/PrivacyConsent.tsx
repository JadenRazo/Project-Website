import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { api } from '../../utils/apiConfig';

interface ConsentProps {
  onConsentUpdate?: (consents: ConsentStatus) => void;
}

interface ConsentStatus {
  analytics: boolean;
  functional: boolean;
  marketing: boolean;
}

const ConsentBanner = styled.div<{ show: boolean }>`
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: ${({ theme }) => theme.colors.card};
  border-top: 2px solid ${({ theme }) => theme.colors.primary};
  padding: 1.5rem;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
  transform: translateY(${({ show }) => show ? '0' : '100%'});
  transition: transform 0.3s ease;
  z-index: 1000;

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
  }
`;

const ConsentContent = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const ConsentHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: start;
  gap: 2rem;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    flex-direction: column;
    gap: 1rem;
  }
`;

const ConsentText = styled.div`
  flex: 1;

  h3 {
    font-size: 1.25rem;
    margin-bottom: 0.5rem;
    color: ${({ theme }) => theme.colors.text};
  }

  p {
    font-size: 0.95rem;
    color: ${({ theme }) => theme.colors.textSecondary};
    line-height: 1.6;
  }
`;

const ConsentOptions = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin: 1rem 0;
`;

const ConsentOption = styled.label`
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background 0.2s ease;

  &:hover {
    background: ${({ theme }) => theme.colors.background};
  }

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  span {
    font-size: 0.95rem;
    color: ${({ theme }) => theme.colors.text};
  }

  small {
    color: ${({ theme }) => theme.colors.textSecondary};
    font-size: 0.85rem;
    margin-left: auto;
  }
`;

const ConsentActions = styled.div`
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  flex-wrap: wrap;

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    justify-content: stretch;

    button {
      flex: 1;
    }
  }
`;

const Button = styled.button<{ variant?: 'primary' | 'secondary' | 'text' }>`
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  ${({ variant, theme }) => {
    switch (variant) {
      case 'primary':
        return `
          background: ${theme.colors.primary};
          color: white;
          &:hover {
            background: ${theme.colors.primaryHover};
            transform: translateY(-1px);
          }
        `;
      case 'secondary':
        return `
          background: transparent;
          color: ${theme.colors.primary};
          border: 1px solid ${theme.colors.primary};
          &:hover {
            background: ${theme.colors.primary}10;
          }
        `;
      default:
        return `
          background: transparent;
          color: ${theme.colors.textSecondary};
          &:hover {
            color: ${theme.colors.text};
          }
        `;
    }
  }}
`;

const PrivacyLink = styled.a`
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
`;

const generateSessionHash = (): string => {
  const timestamp = Date.now().toString();
  const random = Math.random().toString(36).substring(2);
  return btoa(`${timestamp}-${random}`).replace(/[^a-zA-Z0-9]/g, '');
};

const PrivacyConsent: React.FC<ConsentProps> = ({ onConsentUpdate }) => {
  const [showBanner, setShowBanner] = useState(false);
  const [consents, setConsents] = useState<ConsentStatus>({
    analytics: false,
    functional: true, // Functional cookies are usually necessary
    marketing: false,
  });
  const [sessionHash] = useState(() => generateSessionHash());

  useEffect(() => {
    // Check if consent has already been given
    const consentCookie = document.cookie
      .split('; ')
      .find(row => row.startsWith('privacy_consent='));
    
    if (!consentCookie) {
      // Show banner after a short delay
      setTimeout(() => setShowBanner(true), 1000);
    } else {
      // Load existing consent status
      loadConsentStatus();
    }
  }, []);

  const loadConsentStatus = async () => {
    try {
      const response = await api.get<ConsentStatus>(`/api/v1/privacy/consent/${sessionHash}`);
      setConsents(response);
      if (onConsentUpdate) {
        onConsentUpdate(response);
      }
    } catch (error) {
      // If no consent found, use defaults
    }
  };

  const handleAcceptAll = async () => {
    const allConsents = {
      analytics: true,
      functional: true,
      marketing: true,
    };
    
    await saveConsents(allConsents);
  };

  const handleAcceptSelected = async () => {
    await saveConsents(consents);
  };

  const handleRejectAll = async () => {
    const minimalConsents = {
      analytics: false,
      functional: true, // Keep functional as they're necessary
      marketing: false,
    };
    
    await saveConsents(minimalConsents);
  };

  const saveConsents = async (consentData: ConsentStatus) => {
    try {
      await api.post('/api/v1/privacy/consent', {
        sessionHash,
        consents: consentData,
      });

      // Set cookie to remember consent was given
      document.cookie = `privacy_consent=granted; path=/; max-age=${365 * 24 * 60 * 60}; SameSite=Strict`;
      
      setConsents(consentData);
      setShowBanner(false);
      
      if (onConsentUpdate) {
        onConsentUpdate(consentData);
      }
    } catch (error) {
      console.error('Failed to save consent preferences:', error);
    }
  };

  const toggleConsent = (type: keyof ConsentStatus) => {
    setConsents(prev => ({
      ...prev,
      [type]: !prev[type],
    }));
  };

  return (
    <ConsentBanner show={showBanner}>
      <ConsentContent>
        <ConsentHeader>
          <ConsentText>
            <h3>Your Privacy Matters</h3>
            <p>
              We use cookies and similar technologies to enhance your experience, analyze site traffic,
              and provide personalized content. You can choose which types of cookies you allow.
              Read our <PrivacyLink href="/privacy" target="_blank">Privacy Policy</PrivacyLink> to learn more.
            </p>
          </ConsentText>
        </ConsentHeader>

        <ConsentOptions>
          <ConsentOption>
            <input
              type="checkbox"
              checked={consents.functional}
              disabled
              readOnly
            />
            <span>Necessary Cookies</span>
            <small>Always active</small>
          </ConsentOption>
          
          <ConsentOption>
            <input
              type="checkbox"
              checked={consents.analytics}
              onChange={() => toggleConsent('analytics')}
            />
            <span>Analytics Cookies</span>
            <small>Help us improve</small>
          </ConsentOption>
          
          <ConsentOption>
            <input
              type="checkbox"
              checked={consents.marketing}
              onChange={() => toggleConsent('marketing')}
            />
            <span>Marketing Cookies</span>
            <small>Personalized ads</small>
          </ConsentOption>
        </ConsentOptions>

        <ConsentActions>
          <Button variant="text" onClick={handleRejectAll}>
            Reject All
          </Button>
          <Button variant="secondary" onClick={handleAcceptSelected}>
            Accept Selected
          </Button>
          <Button variant="primary" onClick={handleAcceptAll}>
            Accept All
          </Button>
        </ConsentActions>
      </ConsentContent>
    </ConsentBanner>
  );
};

export default PrivacyConsent;