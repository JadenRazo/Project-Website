import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import SEO from '../../components/common/SEO';

interface BlogPost {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  featured_image: string;
  status: string;
  published_at: string;
  tags: string[];
  view_count: number;
  read_time_minutes: number;
  is_featured: boolean;
}

interface ListResult {
  posts: BlogPost[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

const BlogContainer = styled.div`
  max-width: var(--page-max-width, 1200px);
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

const PageHeader = styled.div`
  text-align: center;
  margin-bottom: 3rem;
`;

const PageTitle = styled.h1`
  font-size: 3rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  font-weight: 700;
  position: relative;
  display: inline-block;

  @media (max-width: 768px) {
    font-size: 2.5rem;
  }

  &::after {
    content: '';
    position: absolute;
    bottom: -10px;
    left: 50%;
    transform: translateX(-50%);
    width: 80px;
    height: 4px;
    background: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    border-radius: 2px;
  }
`;

const PageDescription = styled.p`
  max-width: 650px;
  margin: 1.5rem auto 0;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  font-size: 1.1rem;
  line-height: 1.6;

  @media (max-width: 768px) {
    font-size: 1rem;
  }
`;

const SearchBar = styled.input`
  width: 100%;
  max-width: 500px;
  margin: 2rem auto;
  display: block;
  padding: 0.75rem 1.25rem;
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'};
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  border-radius: 8px;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  font-size: 1rem;
  outline: none;
  transition: border-color 0.2s;

  &:focus {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  }

  &::placeholder {
    color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
  }
`;

const TagFilters = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 2rem;
`;

const TagPill = styled.button<{ $active: boolean }>`
  padding: 0.35rem 0.85rem;
  border-radius: 20px;
  border: 1px solid ${({ theme, $active }) =>
    $active ? theme?.colors?.primary || '#007bff' : theme?.colors?.border || 'rgba(255,255,255,0.15)'};
  background: ${({ theme, $active }) =>
    $active ? (theme?.colors?.primary || '#007bff') + '20' : 'transparent'};
  color: ${({ theme, $active }) =>
    $active ? theme?.colors?.primary || '#007bff' : theme?.colors?.textSecondary || '#aaa'};
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s;

  &:hover {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  }
`;

const FeaturedSection = styled.div`
  margin-bottom: 3rem;
`;

const FeaturedLabel = styled.h2`
  font-size: 1.4rem;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  margin-bottom: 1.5rem;
  font-weight: 600;
`;

const FeaturedCardWrapper = styled(Link)`
  text-decoration: none;
  color: inherit;
  display: block;
`;

const FeaturedCard = styled(motion.div)`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'};
  border-radius: 16px;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  overflow: hidden;
  transition: border-color 0.3s, box-shadow 0.3s;

  &:hover {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  }

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
`;

const FeaturedImage = styled.div<{ $src?: string }>`
  width: 100%;
  min-height: 280px;
  background: ${({ $src, theme }) =>
    $src ? `url(${$src}) center/cover no-repeat` : theme?.colors?.surface || '#222'};

  @media (max-width: 768px) {
    min-height: 200px;
  }
`;

const FeaturedContent = styled.div`
  padding: 2rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1rem;

  @media (max-width: 768px) {
    padding: 1.5rem;
  }
`;

const PostGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 1.5rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
`;

const PostCardLink = styled(Link)`
  text-decoration: none;
  color: inherit;
  display: block;
`;

const PostCard = styled(motion.div)`
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'};
  border-radius: 12px;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: border-color 0.3s, transform 0.3s, box-shadow 0.3s;

  &:hover {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }
`;

const PostImage = styled.div<{ $src?: string }>`
  width: 100%;
  height: 200px;
  background: ${({ $src, theme }) =>
    $src ? `url(${$src}) center/cover no-repeat` : theme?.colors?.background || '#1a1a1a'};
`;

const PostBody = styled.div`
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  flex: 1;
`;

const PostTags = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
`;

const PostTag = styled.span`
  background: ${({ theme }) => (theme?.colors?.primary || '#007bff') + '15'};
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  padding: 0.15rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
`;

const PostTitle = styled.h3`
  font-size: 1.2rem;
  font-weight: 600;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  margin: 0;
  line-height: 1.4;
`;

const PostExcerpt = styled.p`
  font-size: 0.9rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const PostMeta = styled.div`
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.8rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
  margin-top: auto;
`;

const Pagination = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  margin-top: 3rem;
`;

const PageButton = styled.button<{ $active?: boolean }>`
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid ${({ theme, $active }) =>
    $active ? theme?.colors?.primary || '#007bff' : theme?.colors?.border || 'rgba(255,255,255,0.15)'};
  background: ${({ theme, $active }) =>
    $active ? theme?.colors?.primary || '#007bff' : 'transparent'};
  color: ${({ theme, $active }) =>
    $active ? '#fff' : theme?.colors?.text || '#fff'};
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;

  &:hover:not(:disabled) {
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  }

  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
`;

const EmptyState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};

  h3 {
    font-size: 1.4rem;
    margin-bottom: 0.5rem;
    color: ${({ theme }) => theme?.colors?.text || '#fff'};
  }
`;

const LoadingState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#888'};
  font-size: 1.1rem;
`;

const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: { delay: i * 0.08, duration: 0.4, ease: 'easeOut' },
  }),
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const Blog: React.FC = () => {
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [featured, setFeatured] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');
  const [activeTag, setActiveTag] = useState('');
  const [searchDebounce, setSearchDebounce] = useState('');

  useEffect(() => {
    const timer = setTimeout(() => setSearchDebounce(search), 300);
    return () => clearTimeout(timer);
  }, [search]);

  const fetchPosts = useCallback(async () => {
    setLoading(true);
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
      const params = new URLSearchParams({
        page: String(page),
        page_size: '9',
      });
      if (searchDebounce) params.set('search', searchDebounce);
      if (activeTag) params.set('tag', activeTag);

      const endpoint = apiUrl
        ? `${apiUrl}/api/v1/blog?${params}`
        : `/api/v1/blog?${params}`;
      const res = await fetch(endpoint);
      if (!res.ok) throw new Error('Failed to fetch');
      const data: ListResult = await res.json();
      setPosts(data.posts || []);
      setTotalPages(data.total_pages || 1);
    } catch {
      setPosts([]);
    } finally {
      setLoading(false);
    }
  }, [page, searchDebounce, activeTag]);

  const fetchFeatured = useCallback(async () => {
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
      const endpoint = apiUrl
        ? `${apiUrl}/api/v1/blog/featured`
        : '/api/v1/blog/featured';
      const res = await fetch(endpoint);
      if (!res.ok) throw new Error('Failed to fetch');
      const data: BlogPost[] = await res.json();
      setFeatured(data || []);
    } catch {
      setFeatured([]);
    }
  }, []);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  useEffect(() => {
    fetchFeatured();
  }, [fetchFeatured]);

  useEffect(() => {
    setPage(1);
  }, [searchDebounce, activeTag]);

  const allTags = useMemo(() => {
    const tagSet = new Set<string>();
    posts.forEach((p) => p.tags?.forEach((t) => tagSet.add(t)));
    featured.forEach((p) => p.tags?.forEach((t) => tagSet.add(t)));
    return Array.from(tagSet).sort();
  }, [posts, featured]);

  const mainFeatured = featured[0];

  return (
    <>
      <SEO
        title="Blog | Jaden Razo - Software Engineering & Web Development"
        description="Read about software engineering, web development, and technology insights from Jaden Razo."
        path="/blog"
      />
      <BlogContainer>
        <PageHeader>
          <PageTitle>Blog</PageTitle>
          <PageDescription>
            Thoughts on software engineering, web development, and the technologies I work with day to day.
          </PageDescription>
        </PageHeader>

        <SearchBar
          type="text"
          placeholder="Search posts..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          aria-label="Search blog posts"
        />

        {allTags.length > 0 && (
          <TagFilters>
            <TagPill
              $active={activeTag === ''}
              onClick={() => setActiveTag('')}
            >
              All
            </TagPill>
            {allTags.map((tag) => (
              <TagPill
                key={tag}
                $active={activeTag === tag}
                onClick={() => setActiveTag(activeTag === tag ? '' : tag)}
              >
                {tag}
              </TagPill>
            ))}
          </TagFilters>
        )}

        {loading ? (
          <LoadingState>Loading posts...</LoadingState>
        ) : (
          <>
            {mainFeatured && !searchDebounce && !activeTag && page === 1 && (
              <FeaturedSection>
                <FeaturedLabel>Featured</FeaturedLabel>
                <FeaturedCardWrapper to={`/blog/${mainFeatured.slug}`}>
                  <FeaturedCard
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5 }}
                  >
                    <FeaturedImage $src={mainFeatured.featured_image} />
                    <FeaturedContent>
                      <PostTags>
                        {mainFeatured.tags?.map((tag) => (
                          <PostTag key={tag}>{tag}</PostTag>
                        ))}
                      </PostTags>
                      <PostTitle style={{ fontSize: '1.6rem' }}>
                        {mainFeatured.title}
                      </PostTitle>
                      <PostExcerpt>{mainFeatured.excerpt}</PostExcerpt>
                      <PostMeta>
                        <span>{formatDate(mainFeatured.published_at)}</span>
                        <span>{mainFeatured.read_time_minutes} min read</span>
                      </PostMeta>
                    </FeaturedContent>
                  </FeaturedCard>
                </FeaturedCardWrapper>
              </FeaturedSection>
            )}

            {posts.length > 0 ? (
              <AnimatePresence mode="wait">
                <PostGrid key={`${page}-${activeTag}-${searchDebounce}`}>
                  {posts.map((post, i) => (
                    <PostCardLink key={post.id} to={`/blog/${post.slug}`}>
                      <PostCard
                        custom={i}
                        variants={cardVariants}
                        initial="hidden"
                        animate="visible"
                      >
                        {post.featured_image && (
                          <PostImage $src={post.featured_image} />
                        )}
                        <PostBody>
                          {post.tags?.length > 0 && (
                            <PostTags>
                              {post.tags.map((tag) => (
                                <PostTag key={tag}>{tag}</PostTag>
                              ))}
                            </PostTags>
                          )}
                          <PostTitle>{post.title}</PostTitle>
                          <PostExcerpt>{post.excerpt}</PostExcerpt>
                          <PostMeta>
                            <span>{formatDate(post.published_at)}</span>
                            <span>{post.read_time_minutes} min read</span>
                          </PostMeta>
                        </PostBody>
                      </PostCard>
                    </PostCardLink>
                  ))}
                </PostGrid>
              </AnimatePresence>
            ) : (
              <EmptyState>
                <h3>No posts found</h3>
                <p>
                  {searchDebounce || activeTag
                    ? 'Try adjusting your search or filters.'
                    : 'Check back soon for new content.'}
                </p>
              </EmptyState>
            )}

            {totalPages > 1 && (
              <Pagination>
                <PageButton
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
                >
                  Previous
                </PageButton>
                {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                  (p) => (
                    <PageButton
                      key={p}
                      $active={p === page}
                      onClick={() => setPage(p)}
                    >
                      {p}
                    </PageButton>
                  )
                )}
                <PageButton
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </PageButton>
              </Pagination>
            )}
          </>
        )}
      </BlogContainer>
    </>
  );
};

export default Blog;
