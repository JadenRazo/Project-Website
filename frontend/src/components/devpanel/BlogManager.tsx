import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import {
  BlogPost,
  BlogPostFormData,
  defaultBlogPostFormData,
  mapBlogPostToFormData,
  blogApi,
} from '../../utils/devPanelApi';

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2rem;
`;

const ErrorMessage = styled.div`
  padding: 1rem;
  margin-bottom: 1rem;
  background: ${({ theme }) => theme.colors.error}20;
  color: ${({ theme }) => theme.colors.error};
  border-radius: 8px;
  border-left: 4px solid ${({ theme }) => theme.colors.error};
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;

  @media (max-width: 768px) {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
`;

const AddButton = styled.button`
  padding: 0.6rem 1.2rem;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: opacity 0.2s;

  &:hover { opacity: 0.9; }
`;

const PostList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
`;

const PostItem = styled(motion.div)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  gap: 1rem;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: flex-start;
  }
`;

const PostInfo = styled.div`
  flex: 1;
  min-width: 0;
`;

const PostTitle = styled.div`
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  margin-bottom: 0.25rem;
`;

const PostMetaRow = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  font-size: 0.8rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const StatusBadge = styled.span<{ $status: string }>`
  padding: 0.15rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
  background: ${({ theme, $status }) =>
    $status === 'published' ? theme.colors.success + '20' :
    $status === 'archived' ? theme.colors.error + '20' :
    theme.colors.warning + '20'};
  color: ${({ theme, $status }) =>
    $status === 'published' ? theme.colors.success :
    $status === 'archived' ? theme.colors.error :
    theme.colors.warning};
`;

const Actions = styled.div`
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
`;

const ActionButton = styled.button<{ $variant?: 'edit' | 'delete' }>`
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  transition: opacity 0.2s;
  background: ${({ theme, $variant }) =>
    $variant === 'delete' ? theme.colors.error : theme.colors.primary};
  color: white;

  &:hover { opacity: 0.85; }
`;

const Form = styled(motion.form)`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
`;

const FormTitle = styled.h3`
  margin: 0 0 0.5rem;
  color: ${({ theme }) => theme.colors.text};
`;

const FormRow = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
`;

const Label = styled.label`
  font-size: 0.85rem;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const Input = styled.input`
  padding: 0.6rem 0.8rem;
  background: ${({ theme }) => theme.colors.background};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 6px;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.95rem;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const TextArea = styled.textarea`
  padding: 0.6rem 0.8rem;
  background: ${({ theme }) => theme.colors.background};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 6px;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.95rem;
  min-height: 200px;
  resize: vertical;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const Select = styled.select`
  padding: 0.6rem 0.8rem;
  background: ${({ theme }) => theme.colors.background};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 6px;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.95rem;

  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const CheckboxRow = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const FormActions = styled.div`
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  margin-top: 0.5rem;
`;

const SaveButton = styled.button`
  padding: 0.6rem 1.2rem;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;

  &:hover { opacity: 0.9; }
  &:disabled { opacity: 0.5; cursor: not-allowed; }
`;

const CancelButton = styled.button`
  padding: 0.6rem 1.2rem;
  background: transparent;
  color: ${({ theme }) => theme.colors.text};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;

  &:hover { background: ${({ theme }) => theme.colors.surface}; }
`;

const ConfirmDialog = styled(motion.div)`
  padding: 1.5rem;
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.error}40;
  border-radius: 8px;
  text-align: center;

  p {
    margin-bottom: 1rem;
    color: ${({ theme }) => theme.colors.text};
  }
`;

const EmptyState = styled.div`
  text-align: center;
  padding: 3rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const BlogManager: React.FC = () => {
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState<BlogPostFormData>(defaultBlogPostFormData);
  const [saving, setSaving] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchPosts = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const result = await blogApi.getAll();
      setPosts(result.posts || []);
    } catch (err) {
      setError('Failed to load blog posts');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPosts();
  }, [fetchPosts]);

  const handleCreate = () => {
    setEditingId(null);
    setFormData(defaultBlogPostFormData);
    setShowForm(true);
  };

  const handleEdit = (post: BlogPost) => {
    setEditingId(post.id);
    setFormData(mapBlogPostToFormData(post));
    setShowForm(true);
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingId(null);
    setFormData(defaultBlogPostFormData);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError('');

    const payload: Record<string, unknown> = {
      ...formData,
      tags: formData.tags_input
        .split(',')
        .map((t) => t.trim())
        .filter(Boolean),
    };
    delete (payload as Record<string, unknown>).tags_input;

    try {
      if (editingId) {
        await blogApi.update(editingId, payload);
      } else {
        await blogApi.create(payload);
      }
      handleCancel();
      await fetchPosts();
    } catch (err) {
      setError(editingId ? 'Failed to update post' : 'Failed to create post');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await blogApi.delete(id);
      setDeletingId(null);
      await fetchPosts();
    } catch {
      setError('Failed to delete post');
    }
  };

  const updateField = (field: keyof BlogPostFormData, value: string | boolean) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <Container>
      {error && <ErrorMessage>{error}</ErrorMessage>}

      <Header>
        <h3 style={{ margin: 0 }}>
          Blog Posts {!loading && `(${posts.length})`}
        </h3>
        <AddButton onClick={handleCreate}>New Post</AddButton>
      </Header>

      <AnimatePresence>
        {showForm && (
          <Form
            key="blog-form"
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            onSubmit={handleSubmit}
          >
            <FormTitle>{editingId ? 'Edit Post' : 'Create New Post'}</FormTitle>

            <FormRow>
              <Label>Title *</Label>
              <Input
                value={formData.title}
                onChange={(e) => updateField('title', e.target.value)}
                required
              />
            </FormRow>

            <FormRow>
              <Label>Slug (auto-generated if empty)</Label>
              <Input
                value={formData.slug}
                onChange={(e) => updateField('slug', e.target.value)}
                placeholder="my-post-slug"
              />
            </FormRow>

            <FormRow>
              <Label>Excerpt</Label>
              <Input
                value={formData.excerpt}
                onChange={(e) => updateField('excerpt', e.target.value)}
              />
            </FormRow>

            <FormRow>
              <Label>Content (Markdown)</Label>
              <TextArea
                value={formData.content}
                onChange={(e) => updateField('content', e.target.value)}
              />
            </FormRow>

            <FormRow>
              <Label>Featured Image URL</Label>
              <Input
                value={formData.featured_image}
                onChange={(e) => updateField('featured_image', e.target.value)}
              />
            </FormRow>

            <FormRow>
              <Label>Tags (comma-separated)</Label>
              <Input
                value={formData.tags_input}
                onChange={(e) => updateField('tags_input', e.target.value)}
                placeholder="go, react, typescript"
              />
            </FormRow>

            <FormRow>
              <Label>Status</Label>
              <Select
                value={formData.status}
                onChange={(e) => updateField('status', e.target.value)}
              >
                <option value="draft">Draft</option>
                <option value="published">Published</option>
                <option value="archived">Archived</option>
              </Select>
            </FormRow>

            <CheckboxRow>
              <input
                type="checkbox"
                id="blog-is-featured"
                checked={formData.is_featured}
                onChange={(e) => updateField('is_featured', e.target.checked)}
              />
              <Label htmlFor="blog-is-featured">Featured</Label>
            </CheckboxRow>

            <CheckboxRow>
              <input
                type="checkbox"
                id="blog-is-visible"
                checked={formData.is_visible}
                onChange={(e) => updateField('is_visible', e.target.checked)}
              />
              <Label htmlFor="blog-is-visible">Visible</Label>
            </CheckboxRow>

            <FormActions>
              <CancelButton type="button" onClick={handleCancel}>
                Cancel
              </CancelButton>
              <SaveButton type="submit" disabled={saving || !formData.title}>
                {saving ? 'Saving...' : editingId ? 'Update' : 'Create'}
              </SaveButton>
            </FormActions>
          </Form>
        )}
      </AnimatePresence>

      {loading ? (
        <EmptyState>Loading posts...</EmptyState>
      ) : posts.length === 0 ? (
        <EmptyState>No blog posts yet. Create one to get started.</EmptyState>
      ) : (
        <PostList>
          <AnimatePresence>
            {posts.map((post) => (
              <React.Fragment key={post.id}>
                {deletingId === post.id ? (
                  <ConfirmDialog
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                  >
                    <p>Delete "{post.title}"?</p>
                    <Actions style={{ justifyContent: 'center' }}>
                      <CancelButton onClick={() => setDeletingId(null)}>
                        Cancel
                      </CancelButton>
                      <ActionButton
                        $variant="delete"
                        onClick={() => handleDelete(post.id)}
                      >
                        Delete
                      </ActionButton>
                    </Actions>
                  </ConfirmDialog>
                ) : (
                  <PostItem
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                  >
                    <PostInfo>
                      <PostTitle>{post.title}</PostTitle>
                      <PostMetaRow>
                        <StatusBadge $status={post.status}>
                          {post.status}
                        </StatusBadge>
                        <span>{post.view_count} views</span>
                        {post.published_at && (
                          <span>
                            {new Date(post.published_at).toLocaleDateString()}
                          </span>
                        )}
                        {post.is_featured && <span>Featured</span>}
                      </PostMetaRow>
                    </PostInfo>
                    <Actions>
                      <ActionButton
                        $variant="edit"
                        onClick={() => handleEdit(post)}
                      >
                        Edit
                      </ActionButton>
                      <ActionButton
                        $variant="delete"
                        onClick={() => setDeletingId(post.id)}
                      >
                        Delete
                      </ActionButton>
                    </Actions>
                  </PostItem>
                )}
              </React.Fragment>
            ))}
          </AnimatePresence>
        </PostList>
      )}
    </Container>
  );
};

export default BlogManager;
