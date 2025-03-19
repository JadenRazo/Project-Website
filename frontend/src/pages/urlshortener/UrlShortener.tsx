import React, { useState, useEffect } from 'react';
import styled from 'styled-components';

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
  max-width: 1000px;
  margin: 0 auto;
`;

const Header = styled.header`
  margin-bottom: 2rem;
  text-align: center;
`;

const Title = styled.h1`
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.primary};
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
  border-radius: 8px;
  margin-bottom: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  @media (min-width: 768px) {
    flex-direction: row;
    align-items: flex-end;
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
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
`;

const Input = styled.input`
  width: 100%;
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.input};
  color: ${({ theme }) => theme.colors.text};
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const Button = styled.button`
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
  }

  @media (min-width: 768px) {
    height: 42px;
  }
`;

const ResultContainer = styled.div<{ success?: boolean }>`
  padding: 1rem;
  margin-top: 1rem;
  background: ${({ theme, success }) => 
    success ? theme.colors.success + '20' : theme.colors.error + '20'
  };
  color: ${({ theme, success }) => 
    success ? theme.colors.success : theme.colors.error
  };
  border-radius: 4px;
  border-left: 4px solid ${({ theme, success }) => 
    success ? theme.colors.success : theme.colors.error
  };
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  border: 1px solid ${({ theme }) => theme.colors.primary};
  color: ${({ theme }) => theme.colors.primary};
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: ${({ theme }) => theme.colors.primary}20;
  }
`;

const UrlListContainer = styled.div`
  margin-top: 3rem;
`;

const UrlListHeader = styled.h2`
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: ${({ theme }) => theme.colors.text};
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  padding-bottom: 0.5rem;
`;

const StatsTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
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

const TableBody = styled.tbody``;

const TableRow = styled.tr`
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  
  &:hover {
    background: ${({ theme }) => theme.colors.backgroundHover};
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

const LoadingMessage = styled.p`
  text-align: center;
  color: ${({ theme }) => theme.colors.textSecondary};
  padding: 2rem 0;
`;

// Main Component
const UrlShortener: React.FC = () => {
  const [url, setUrl] = useState<string>('');
  const [customCode, setCustomCode] = useState<string>('');
  const [result, setResult] = useState<UrlResponse | null>(null);
  const [urlList, setUrlList] = useState<ShortenedUrl[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchUrlList = async () => {
      try {
        // Replace with your actual API endpoint
        const response = await fetch('/urlshortener/api/urls');
        
        // If the API is not ready, use mock data
        if (!response.ok) {
          // Simulate API response with mock data
          setTimeout(() => {
            setUrlList([
              {
                id: '1',
                originalUrl: 'https://www.example.com/some/very/long/path/that/needs/shortening',
                shortCode: 'abc123',
                createdAt: new Date().toISOString(),
                visits: 42
              },
              {
                id: '2',
                originalUrl: 'https://github.com/JadenRazo/Project-Website',
                shortCode: 'github',
                createdAt: new Date(Date.now() - 86400000).toISOString(),
                visits: 127
              },
              {
                id: '3',
                originalUrl: 'https://jadenrazo.dev/contact',
                shortCode: 'contact',
                createdAt: new Date(Date.now() - 172800000).toISOString(),
                visits: 15
              }
            ]);
            setLoading(false);
          }, 1000);
          return;
        }
        
        const data = await response.json();
        setUrlList(data);
      } catch (err) {
        console.error('Error fetching URL list:', err);
        
        // Fallback to mock data
        setUrlList([
          {
            id: '1',
            originalUrl: 'https://www.example.com/some/very/long/path/that/needs/shortening',
            shortCode: 'abc123',
            createdAt: new Date().toISOString(),
            visits: 42
          },
          {
            id: '2',
            originalUrl: 'https://github.com/JadenRazo/Project-Website',
            shortCode: 'github',
            createdAt: new Date(Date.now() - 86400000).toISOString(),
            visits: 127
          },
          {
            id: '3',
            originalUrl: 'https://jadenrazo.dev/contact',
            shortCode: 'contact',
            createdAt: new Date(Date.now() - 172800000).toISOString(),
            visits: 15
          }
        ]);
      } finally {
        setLoading(false);
      }
    };

    fetchUrlList();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!url) return;
    
    try {
      // Replace with your actual API endpoint
      const response = await fetch('/urlshortener/api/shorten', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url, customCode }),
      });
      
      // If the API is not ready, simulate a response
      if (!response.ok) {
        // Simulate API response
        const mockShortCode = customCode || Math.random().toString(36).substring(2, 8);
        setResult({
          success: true,
          shortCode: mockShortCode,
          shortUrl: `${window.location.origin}/s/${mockShortCode}`,
        });
        
        // Add to list
        setUrlList(prev => [
          {
            id: Math.random().toString(),
            originalUrl: url,
            shortCode: mockShortCode,
            createdAt: new Date().toISOString(),
            visits: 0
          },
          ...prev
        ]);
        
        // Reset form
        setUrl('');
        setCustomCode('');
        return;
      }
      
      const data = await response.json();
      setResult(data);
      
      if (data.success) {
        // Refresh URL list
        const listResponse = await fetch('/urlshortener/api/urls');
        const listData = await listResponse.json();
        setUrlList(listData);
        
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
    </PageContainer>
  );
};

export default UrlShortener; 