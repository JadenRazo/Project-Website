import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { useAuth } from '../../hooks/useAuth';
import AuthModal from '../../components/auth/AuthModal';
import { useScrollTo } from '../../hooks/useScrollTo';

import { SCROLL_DELAYS } from '../../utils/scrollConfig';

// Types
interface ShortenedUrl {
  id: string;
  originalUrl: string;
  shortCode: string;
  createdAt: string;
  visits: number;
}

interface UrlResponse {
  success: boolean;
  shortCode?: string;
  shortUrl?: string;
  message?: string;
}

// Styled Components
const PageContainer = styled.div`
  padding: 2rem;
  padding-top: calc(80px + 2rem);
  max-width: 1000px;
  margin: 0 auto;
  min-height: 100vh;
  animation: fadeIn 0.4s ease-out;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 1.5rem;
    padding-top: calc(70px + 1.5rem);
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
    padding-top: calc(60px + 1rem);
  }
  
  @keyframes fadeIn {
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

const Header = styled.header`
  margin-bottom: 3rem;
  text-align: center;
  animation: slideDown 0.5s ease-out;
  animation-delay: 0.1s;
  animation-fill-mode: both;
  
  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
`;

const Title = styled.h1`
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.primary};
  font-weight: 600;
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -0.25rem;
    left: 50%;
    transform: translateX(-50%);
    width: 80px;
    height: 4px;
    background: ${({ theme }) => theme.colors.primary};
    border-radius: 2px;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 2rem;
  }
`;

const Subtitle = styled.p`
  font-size: 1.1rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const UrlForm = styled.form`
  display: flex;
  flex-direction: column;
  background: ${({ theme }) => theme.colors.card};
  padding: 2rem;
  border-radius: 12px;
  margin-bottom: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border: 1px solid ${({ theme }) => theme.colors.border};
  animation: slideUp 0.5s ease-out;
  animation-delay: 0.2s;
  animation-fill-mode: both;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }

  @media (min-width: 768px) {
    flex-direction: row;
    align-items: flex-end;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1.5rem;
  }
  
  @keyframes slideUp {
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

const FormGroup = styled.div`
  flex: 1;
  margin-bottom: 1rem;

  @media (min-width: 768px) {
    margin-bottom: 0;
    margin-right: 1rem;
  }
`;

const Label = styled.label`
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  opacity: 0.9;
`;

const Input = styled.input`
  width: 100%;
  padding: 0.875rem 1rem;
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  transition: all 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
    box-shadow: 0 0 0 3px ${({ theme }) => theme.colors.primary}20;
    transform: translateY(-1px);
  }
  
  &:hover:not(:focus) {
    border-color: ${({ theme }) => theme.colors.borderHover};
  }
  
  &::placeholder {
    color: ${({ theme }) => theme.colors.textSecondary};
    opacity: 0.7;
  }
`;

const Button = styled.button`
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.875rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
  min-height: 48px;
  
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
    animation: buttonPop 0.3s ease-out;
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

const ResultContainer = styled.div<{ success?: boolean }>`
  padding: 1rem 1.25rem;
  margin-top: 1rem;
  background: ${({ theme, success }) => 
    success ? theme.colors.success + '20' : theme.colors.error + '20'
  };
  color: ${({ theme, success }) => 
    success ? theme.colors.success : theme.colors.error
  };
  border-radius: 8px;
  border-left: 4px solid ${({ theme, success }) => 
    success ? theme.colors.success : theme.colors.error
  };
  display: flex;
  justify-content: space-between;
  align-items: center;
  animation: ${({ success }) => success ? 'slideInSuccess' : 'shake'} 0.4s ease-out;
  
  @keyframes slideInSuccess {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
  
  @keyframes shake {
    0%, 100% { transform: translateX(0); }
    10%, 30%, 50%, 70%, 90% { transform: translateX(-2px); }
    20%, 40%, 60%, 80% { transform: translateX(2px); }
  }
`;

const ShortenedUrlText = styled.a`
  font-weight: 500;
  text-decoration: none;
  color: ${({ theme }) => theme.colors.primary};

  &:hover {
    text-decoration: underline;
  }
`;

const CopyButton = styled.button`
  background: transparent;
  border: 2px solid ${({ theme }) => theme.colors.primary};
  color: ${({ theme }) => theme.colors.primary};
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
  font-size: 0.875rem;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    background: ${({ theme }) => theme.colors.primary}15;
    transform: translate(-50%, -50%);
    transition: width 0.4s ease, height 0.4s ease;
  }

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 8px ${({ theme }) => theme.colors.primary}30;
    
    &::before {
      width: 100%;
      height: 100%;
    }
  }
  
  &:active {
    transform: scale(0.95);
  }
`;

const UrlListContainer = styled.div`
  margin-top: 3rem;
  animation: fadeInUp 0.6s ease-out;
  animation-delay: 0.3s;
  animation-fill-mode: both;
  
  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
`;

const UrlListHeader = styled.h2`
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
  color: ${({ theme }) => theme.colors.text};
  font-weight: 600;
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -0.5rem;
    left: 0;
    width: 60px;
    height: 3px;
    background: ${({ theme }) => theme.colors.primary};
    border-radius: 2px;
  }
`;

const StatsTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
  background: ${({ theme }) => theme.colors.card};
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const TableHead = styled.thead`
  background: ${({ theme }) => theme.colors.background};
  border-bottom: 2px solid ${({ theme }) => theme.colors.border};
`;

const TableHeader = styled.th`
  text-align: left;
  padding: 1rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const TableBody = styled.tbody`
  tr {
    animation: tableRowFade 0.3s ease-out;
    animation-fill-mode: both;
    
    &:nth-child(1) { animation-delay: 0.1s; }
    &:nth-child(2) { animation-delay: 0.15s; }
    &:nth-child(3) { animation-delay: 0.2s; }
    &:nth-child(4) { animation-delay: 0.25s; }
    &:nth-child(5) { animation-delay: 0.3s; }
    &:nth-child(n+6) { animation-delay: 0.35s; }
  }
  
  @keyframes tableRowFade {
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

const TableRow = styled.tr`
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  transition: all 0.2s ease;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: ${({ theme }) => theme.colors.backgroundHover};
    transform: scale(1.005);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }
`;

const TableCell = styled.td`
  padding: 1rem;
  color: ${({ theme }) => theme.colors.text};
  vertical-align: middle;
`;

const OriginalUrl = styled.div`
  max-width: 250px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const ShortCode = styled.div`
  font-family: monospace;
  color: ${({ theme }) => theme.colors.primary};
`;

const VisitCount = styled.div`
  font-weight: 500;
`;

const LoadingMessage = styled.div`
  text-align: center;
  color: ${({ theme }) => theme.colors.textSecondary};
  padding: 3rem 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  
  &::before {
    content: '';
    width: 40px;
    height: 40px;
    border: 3px solid ${({ theme }) => theme.colors.border};
    border-top-color: ${({ theme }) => theme.colors.primary};
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
`;

const AuthPrompt = styled.div`
  text-align: center;
  padding: 2rem;
  background: ${({ theme }) => theme.colors.background};
  border-radius: 8px;
  border: 1px dashed ${({ theme }) => theme.colors.border};
  animation: fadeIn 0.4s ease-out;
  scroll-margin-top: 100px;
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
`;

const AuthPromptText = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 1rem;
  margin-bottom: 1.5rem;
`;

const AuthButtons = styled.div`
  display: flex;
  gap: 1rem;
  justify-content: center;
  align-items: center;
`;

const AuthButton = styled.button<{ variant?: 'primary' | 'secondary' }>`
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
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
  }
`;

// Main Component
const UrlShortener: React.FC = () => {
  // Ensure page scrolls to top when navigated to
  
  
  const [url, setUrl] = useState<string>('');
  const [customCode, setCustomCode] = useState<string>('');
  const [result, setResult] = useState<UrlResponse | null>(null);
  const [urlList, setUrlList] = useState<ShortenedUrl[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [showAuthModal, setShowAuthModal] = useState<boolean>(false);
  const [authModalMode, setAuthModalMode] = useState<'login' | 'register'>('login');

  const { user, isAuthenticated } = useAuth();
  const { scrollToElement } = useScrollTo();
  const authPromptRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchUrlList = async () => {
      // Only fetch if user is authenticated
      if (!isAuthenticated) {
        setLoading(false);
        return;
      }

      try {
        const headers: HeadersInit = {
          'Content-Type': 'application/json',
        };
        
        // Add auth token if available
        if (user?.token) {
          headers['Authorization'] = `Bearer ${user.token}`;
        }

        const response = await fetch('/api/urls/', {
          headers,
          credentials: 'include'
        });
        
        if (!response.ok) {
          console.error('Failed to fetch URLs:', response.status);
          setLoading(false);
          return;
        }
        
        const data = await response.json();
        // Handle the paginated response
        if (data.urls) {
          setUrlList(data.urls.map((item: any) => ({
            id: item.shortened_url.id,
            originalUrl: item.shortened_url.original_url,
            shortCode: item.shortened_url.short_code,
            createdAt: item.shortened_url.created_at,
            visits: item.total_clicks || 0
          })));
        } else if (Array.isArray(data)) {
          // Handle legacy array response
          setUrlList(data);
        }
      } catch (err) {
        console.error('Error fetching URL list:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchUrlList();
  }, [isAuthenticated, user?.token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url) return;
    
    try {
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };
      
      // Add auth token if available
      if (user?.token) {
        headers['Authorization'] = `Bearer ${user.token}`;
      }

      // Call the URL shortening API
      const response = await fetch('/api/urls/shorten', {
        method: 'POST',
        headers,
        body: JSON.stringify({ 
          url,
          custom_code: customCode,
          title: '',
          description: ''
        }),
        credentials: 'include'
      });
      
      if (!response.ok) {
        const errorData = await response.json();
        setResult({
          success: false,
          message: errorData.error || 'Failed to shorten URL',
        });
        return;
      }
      
      const data = await response.json();
      
      if (data.data) {
        // Handle success response from API
        setResult({
          success: true,
          shortCode: data.data.short_code,
          shortUrl: `${window.location.origin}/s/${data.data.short_code}`,
        });
      } else {
        setResult({
          success: false,
          message: data.error || 'Failed to shorten URL',
        });
      }
      
      if (data.data) {
        // Refresh URL list if authenticated
        if (isAuthenticated && user?.token) {
          const listHeaders: HeadersInit = {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${user.token}`
          };
          
          const listResponse = await fetch('/api/urls/', {
            headers: listHeaders,
            credentials: 'include'
          });
          if (listResponse.ok) {
            const listData = await listResponse.json();
            if (listData.urls) {
              setUrlList(listData.urls.map((item: any) => ({
                id: item.shortened_url.id,
                originalUrl: item.shortened_url.original_url,
                shortCode: item.shortened_url.short_code,
                createdAt: item.shortened_url.created_at,
                visits: item.total_clicks || 0
              })));
            }
          }
        }
        
        // Reset form
        setUrl('');
        setCustomCode('');
      }
    } catch (err) {
      console.error('Error shortening URL:', err);
      setResult({
        success: false,
        message: 'Failed to shorten URL. Please try again.',
      });
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
      .then(() => {
        alert('URL copied to clipboard!');
      })
      .catch(err => {
        console.error('Failed to copy:', err);
      });
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    }).format(date);
  };

  return (
    <PageContainer>
      <Header>
        <Title>URL Shortener</Title>
        <Subtitle>Create shortened links with custom aliases</Subtitle>
      </Header>

      <UrlForm onSubmit={handleSubmit}>
        <FormGroup>
          <Label htmlFor="url">Long URL</Label>
          <Input 
            id="url"
            type="url" 
            placeholder="https://example.com/your-long-url" 
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            required
          />
        </FormGroup>
        
        <FormGroup>
          <Label htmlFor="customCode">Custom Alias (Optional)</Label>
          <Input 
            id="customCode"
            type="text" 
            placeholder="e.g., my-link" 
            value={customCode}
            onChange={(e) => setCustomCode(e.target.value)}
          />
        </FormGroup>
        
        <Button type="submit">Shorten</Button>
      </UrlForm>

      {result && (
        <ResultContainer success={result.success}>
          {result.success ? (
            <>
              <ShortenedUrlText href={result.shortUrl} target="_blank" rel="noopener noreferrer">
                {result.shortUrl}
              </ShortenedUrlText>
              <CopyButton onClick={() => copyToClipboard(result.shortUrl || '')}>
                Copy
              </CopyButton>
            </>
          ) : (
            <span>{result.message}</span>
          )}
        </ResultContainer>
      )}

      <UrlListContainer>
        <UrlListHeader>Your Shortened URLs</UrlListHeader>
        
        {loading ? (
          <LoadingMessage>Loading your URLs...</LoadingMessage>
        ) : !isAuthenticated ? (
          <AuthPrompt ref={authPromptRef}>
            <AuthPromptText>Sign in to view and manage your shortened URLs</AuthPromptText>
            <AuthButtons>
              <AuthButton
                onClick={() => {
                  setAuthModalMode('login');
                  setShowAuthModal(true);
                }}
              >
                Sign In
              </AuthButton>
              <AuthButton
                variant="secondary"
                onClick={() => {
                  setAuthModalMode('register');
                  setShowAuthModal(true);
                }}
              >
                Create Account
              </AuthButton>
            </AuthButtons>
          </AuthPrompt>
        ) : urlList.length === 0 ? (
          <LoadingMessage>No shortened URLs yet. Create your first one above!</LoadingMessage>
        ) : (
          <StatsTable>
            <TableHead>
              <tr>
                <TableHeader>Original URL</TableHeader>
                <TableHeader>Short Code</TableHeader>
                <TableHeader>Created</TableHeader>
                <TableHeader>Visits</TableHeader>
                <TableHeader>Actions</TableHeader>
              </tr>
            </TableHead>
            <TableBody>
              {urlList.map(item => (
                <TableRow key={item.id}>
                  <TableCell>
                    <OriginalUrl title={item.originalUrl}>{item.originalUrl}</OriginalUrl>
                  </TableCell>
                  <TableCell>
                    <ShortCode>{item.shortCode}</ShortCode>
                  </TableCell>
                  <TableCell>{formatDate(item.createdAt)}</TableCell>
                  <TableCell>
                    <VisitCount>{item.visits}</VisitCount>
                  </TableCell>
                  <TableCell>
                    <CopyButton onClick={() => copyToClipboard(`${window.location.origin}/s/${item.shortCode}`)}>
                      Copy
                    </CopyButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </StatsTable>
        )}
      </UrlListContainer>

      <AuthModal
        isOpen={showAuthModal}
        onClose={() => setShowAuthModal(false)}
        initialMode={authModalMode}
      />
    </PageContainer>
  );
};

export default UrlShortener; 