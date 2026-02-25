import React, { useState, useEffect, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import SEO from '../../components/common/SEO';

interface BlogPostData {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  featured_image: string;
  published_at: string;
  tags: string[];
  view_count: number;
  read_time_minutes: number;
}

const PostContainer = styled.div`
  max-width: 800px;
  width: 100%;
  margin: 0 auto;
  padding: calc(2rem + 60px) 2rem 4rem 2rem;
  min-height: 100vh;
  background: ${({ theme }) => theme?.colors?.background || '#111'};
  color: ${({ theme }) => theme?.colors?.text || '#fff'};

  @media (max-width: 768px) {
    padding: calc(1.5rem + 60px) 1rem 3rem 1rem;
  }
`;

const BackLink = styled(Link)`
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  text-decoration: none;
  font-size: 0.9rem;
  margin-bottom: 2rem;
  transition: opacity 0.2s;

  &:hover {
    opacity: 0.8;
  }
`;

const PostHeader = styled(motion.div)`
  margin-bottom: 2.5rem;
`;

const Tags = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-bottom: 1rem;
`;

const Tag = styled.span`
  background: ${({ theme }) => (theme?.colors?.primary || '#007bff') + '15'};
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  padding: 0.2rem 0.6rem;
  border-radius: 10px;
  font-size: 0.8rem;
  font-weight: 500;
`;

const Title = styled.h1`
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 700;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  line-height: 1.3;
  margin-bottom: 1rem;
`;

const Meta = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
  font-size: 0.9rem;
`;

const FeaturedImage = styled.div<{ $src: string }>`
  width: 100%;
  aspect-ratio: 21 / 9;
  background: url(${({ $src }) => $src}) center/cover no-repeat;
  border-radius: 12px;
  margin-bottom: 2.5rem;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
`;

const MarkdownContent = styled.div`
  line-height: 1.8;
  font-size: 1.05rem;
  color: ${({ theme }) => theme?.colors?.text || '#e0e0e0'};

  h1, h2, h3, h4, h5, h6 {
    color: ${({ theme }) => theme?.colors?.text || '#fff'};
    margin: 2rem 0 1rem;
    line-height: 1.3;
    font-weight: 600;
  }

  h1 { font-size: 2rem; }
  h2 { font-size: 1.6rem; }
  h3 { font-size: 1.3rem; }

  p {
    margin: 1rem 0;
  }

  a {
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    text-decoration: none;
    &:hover { text-decoration: underline; }
  }

  ul, ol {
    margin: 1rem 0;
    padding-left: 1.5rem;
  }

  li {
    margin: 0.4rem 0;
  }

  blockquote {
    border-left: 4px solid ${({ theme }) => theme?.colors?.primary || '#007bff'};
    margin: 1.5rem 0;
    padding: 0.75rem 1.5rem;
    background: ${({ theme }) => (theme?.colors?.surface || 'rgba(255,255,255,0.03)')};
    border-radius: 0 8px 8px 0;
    font-style: italic;
    color: ${({ theme }) => theme?.colors?.textSecondary || '#bbb'};
  }

  code {
    background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.08)'};
    color: ${({ theme }) => theme?.colors?.primary || '#00ff88'};
    padding: 0.15rem 0.4rem;
    border-radius: 4px;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.9em;
  }

  pre {
    background: ${({ theme }) => theme?.colors?.surface || 'rgba(0,0,0,0.4)'};
    border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
    border-radius: 8px;
    padding: 1.25rem;
    overflow-x: auto;
    margin: 1.5rem 0;

    code {
      background: none;
      color: inherit;
      padding: 0;
      font-size: 0.9rem;
    }
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
    font-size: 0.95rem;

    th, td {
      border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.15)'};
      padding: 0.6rem 1rem;
      text-align: left;
    }

    th {
      background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'};
      font-weight: 600;
    }
  }

  img {
    max-width: 100%;
    border-radius: 8px;
    margin: 1rem 0;
  }

  hr {
    border: none;
    border-top: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
    margin: 2rem 0;
  }
`;

const PostFooter = styled.div`
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
`;

const ViewCount = styled.span`
  font-size: 0.9rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
`;

const ShareButton = styled.button`
  padding: 0.5rem 1rem;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.15)'};
  border-radius: 6px;
  background: transparent;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;

  &:hover {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  }
`;

const LoadingState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
  font-size: 1.1rem;
`;

const ErrorState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};

  h2 {
    color: ${({ theme }) => theme?.colors?.text || '#fff'};
    margin-bottom: 0.5rem;
  }
`;

const formatDate = (dateStr: string) => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const BlogPost: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const [post, setPost] = useState<BlogPostData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [copied, setCopied] = useState(false);

  const fetchPost = useCallback(async () => {
    if (!slug) return;
    setLoading(true);
    setError(false);
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
      const endpoint = apiUrl
        ? `${apiUrl}/api/v1/blog/${slug}`
        : `/api/v1/blog/${slug}`;
      const res = await fetch(endpoint);
      if (!res.ok) throw new Error('Not found');
      const data: BlogPostData = await res.json();
      setPost(data);
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  }, [slug]);

  useEffect(() => {
    fetchPost();
  }, [fetchPost]);

  useEffect(() => {
    if (!slug) return;
    const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
    const endpoint = apiUrl
      ? `${apiUrl}/api/v1/blog/${slug}/view`
      : `/api/v1/blog/${slug}/view`;
    fetch(endpoint, { method: 'POST' })
      .then((res) => {
        if (res.ok) {
          setPost((prev) => prev ? { ...prev, view_count: prev.view_count + 1 } : prev);
        }
      })
      .catch(() => {});
  }, [slug]);

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  if (loading) {
    return (
      <PostContainer>
        <LoadingState>Loading post...</LoadingState>
      </PostContainer>
    );
  }

  if (error || !post) {
    return (
      <PostContainer>
        <BackLink to="/blog">← Back to Blog</BackLink>
        <ErrorState>
          <h2>Post not found</h2>
          <p>The post you're looking for doesn't exist or has been removed.</p>
        </ErrorState>
      </PostContainer>
    );
  }

  const jsonLd = {
    '@context': 'https://schema.org',
    '@type': 'BlogPosting',
    headline: post.title,
    description: post.excerpt,
    image: post.featured_image || undefined,
    datePublished: post.published_at,
    author: {
      '@type': 'Person',
      name: 'Jaden Razo',
    },
    url: window.location.href,
  };

  return (
    <>
      <SEO
        title={`${post.title} | Jaden Razo Blog`}
        description={post.excerpt}
        path={`/blog/${post.slug}`}
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
      />
      <PostContainer>
        <BackLink to="/blog">← Back to Blog</BackLink>

        <PostHeader
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          {post.tags?.length > 0 && (
            <Tags>
              {post.tags.map((tag) => (
                <Tag key={tag}>{tag}</Tag>
              ))}
            </Tags>
          )}
          <Title>{post.title}</Title>
          <Meta>
            <span>{formatDate(post.published_at)}</span>
            <span>{post.read_time_minutes} min read</span>
            <span>{post.view_count} views</span>
          </Meta>
        </PostHeader>

        {post.featured_image && (
          <FeaturedImage $src={post.featured_image} />
        )}

        <MarkdownContent>
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeHighlight]}
          >
            {post.content}
          </ReactMarkdown>
        </MarkdownContent>

        <PostFooter>
          <ViewCount>{post.view_count} views</ViewCount>
          <ShareButton onClick={handleShare}>
            {copied ? 'Copied!' : 'Copy Link'}
          </ShareButton>
        </PostFooter>
      </PostContainer>
    </>
  );
};

export default BlogPost;
