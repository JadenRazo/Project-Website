import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';

// Types
interface Project {
  id: string;
  name: string;
  description: string;
  repo_url: string;
  live_url?: string;
  tags: string[];
  status: 'active' | 'archived' | 'draft';
  created_at: string;
  updated_at: string;
}

interface ProjectFormData {
  name: string;
  description: string;
  repo_url: string;
  live_url: string;
  tags: string[];
  status: 'active' | 'archived' | 'draft';
}

// Styled Components
const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2rem;
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

const Title = styled.h2`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.5rem;
`;

const AddButton = styled.button`
  padding: 0.75rem 1.5rem;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover || theme.colors.primary};
    transform: translateY(-1px);
  }
`;

const ProjectGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
`;

const ProjectCard = styled(motion.div)`
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  padding: 1.5rem;
  transition: all 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
`;

const ProjectHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
`;

const ProjectName = styled.h3`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.25rem;
`;

const StatusBadge = styled.span<{ status: string }>`
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  background: ${({ theme, status }) => 
    status === 'active' ? theme.colors.success + '20' :
    status === 'draft' ? theme.colors.warning + '20' :
    theme.colors.error + '20'
  };
  color: ${({ theme, status }) => 
    status === 'active' ? theme.colors.success :
    status === 'draft' ? theme.colors.warning :
    theme.colors.error
  };
`;

const ProjectDescription = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 1rem;
  line-height: 1.5;
`;

const TagsContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
`;

const Tag = styled.span`
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.primary + '20'};
  color: ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
  font-size: 0.75rem;
`;

const ProjectActions = styled.div`
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
`;

const ActionButton = styled.button<{ variant: 'edit' | 'delete' | 'view' }>`
  flex: 1;
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, variant }) => 
    variant === 'edit' ? theme.colors.primary :
    variant === 'delete' ? theme.colors.error :
    theme.colors.success
  };
  color: white;
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
`;

const Modal = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 1rem;
`;

const ModalContent = styled(motion.div)`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 2rem;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const Label = styled.label`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
`;

const Input = styled.input`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const TextArea = styled.textarea`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  min-height: 100px;
  resize: vertical;
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const Select = styled.select`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const TagInput = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  min-height: 50px;
`;

const TagInputField = styled.input`
  border: none;
  background: none;
  outline: none;
  color: ${({ theme }) => theme.colors.text};
  flex: 1;
  min-width: 100px;
`;

const TagItem = styled.span`
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.primary + '20'};
  color: ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
  font-size: 0.875rem;
`;

const RemoveTagButton = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.primary};
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primary + '30'};
  }
`;

const ModalActions = styled.div`
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  justify-content: flex-end;
`;

const ModalButton = styled.button<{ variant: 'primary' | 'secondary' }>`
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, variant }) => 
    variant === 'primary' ? theme.colors.primary : theme.colors.border
  };
  color: ${({ theme, variant }) => 
    variant === 'primary' ? 'white' : theme.colors.text
  };
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
`;

const ErrorMessage = styled.div`
  padding: 1rem;
  background: ${({ theme }) => theme.colors.error + '20'};
  color: ${({ theme }) => theme.colors.error};
  border-radius: 4px;
  margin-bottom: 1rem;
`;

const LoadingSpinner = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

// Main Component
export const ProjectManager: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [formData, setFormData] = useState<ProjectFormData>({
    name: '',
    description: '',
    repo_url: '',
    live_url: '',
    tags: [],
    status: 'draft'
  });
  const [tagInput, setTagInput] = useState('');

  useEffect(() => {
    fetchProjects();
  }, []);

  const fetchProjects = async () => {
    try {
      setLoading(true);
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/projects`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch projects');
      }
      
      const data = await response.json();
      setProjects(data.projects || []);
      setError(null);
    } catch (err) {
      console.error('Error fetching projects:', err);
      setError('Failed to load projects');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const url = editingProject 
        ? `${apiUrl}/api/v1/projects/${editingProject.id}`
        : `${apiUrl}/api/v1/projects`;
      
      const method = editingProject ? 'PUT' : 'POST';
      
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error('Failed to save project');
      }

      await fetchProjects();
      setShowModal(false);
      resetForm();
    } catch (err) {
      console.error('Error saving project:', err);
      setError('Failed to save project');
    }
  };

  const handleDelete = async (projectId: string) => {
    if (!window.confirm('Are you sure you want to delete this project?')) {
      return;
    }

    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/projects/${projectId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete project');
      }

      await fetchProjects();
    } catch (err) {
      console.error('Error deleting project:', err);
      setError('Failed to delete project');
    }
  };

  const handleEdit = (project: Project) => {
    setEditingProject(project);
    setFormData({
      name: project.name,
      description: project.description,
      repo_url: project.repo_url,
      live_url: project.live_url || '',
      tags: project.tags || [],
      status: project.status
    });
    setShowModal(true);
  };

  const resetForm = () => {
    setEditingProject(null);
    setFormData({
      name: '',
      description: '',
      repo_url: '',
      live_url: '',
      tags: [],
      status: 'draft'
    });
    setTagInput('');
  };

  const handleAddTag = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      const tag = tagInput.trim();
      if (tag && !formData.tags.includes(tag)) {
        setFormData(prev => ({
          ...prev,
          tags: [...prev.tags, tag]
        }));
      }
      setTagInput('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }));
  };

  if (loading) {
    return <LoadingSpinner>Loading projects...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <Title>Project Management</Title>
        <AddButton onClick={() => setShowModal(true)}>
          Add New Project
        </AddButton>
      </Header>

      {error && <ErrorMessage>{error}</ErrorMessage>}

      <ProjectGrid>
        <AnimatePresence>
          {projects.map((project) => (
            <ProjectCard
              key={project.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              layout
            >
              <ProjectHeader>
                <ProjectName>{project.name}</ProjectName>
                <StatusBadge status={project.status}>
                  {project.status}
                </StatusBadge>
              </ProjectHeader>

              <ProjectDescription>{project.description}</ProjectDescription>

              {project.tags && project.tags.length > 0 && (
                <TagsContainer>
                  {project.tags.map((tag) => (
                    <Tag key={tag}>{tag}</Tag>
                  ))}
                </TagsContainer>
              )}

              <ProjectActions>
                <ActionButton 
                  variant="view" 
                  onClick={() => window.open(project.live_url || project.repo_url, '_blank')}
                >
                  View
                </ActionButton>
                <ActionButton variant="edit" onClick={() => handleEdit(project)}>
                  Edit
                </ActionButton>
                <ActionButton variant="delete" onClick={() => handleDelete(project.id)}>
                  Delete
                </ActionButton>
              </ProjectActions>
            </ProjectCard>
          ))}
        </AnimatePresence>
      </ProjectGrid>

      <AnimatePresence>
        {showModal && (
          <Modal
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setShowModal(false)}
          >
            <ModalContent
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
            >
              <h3>{editingProject ? 'Edit Project' : 'Add New Project'}</h3>
              
              <Form onSubmit={handleSubmit}>
                <FormGroup>
                  <Label>Project Name</Label>
                  <Input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Description</Label>
                  <TextArea
                    value={formData.description}
                    onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Repository URL</Label>
                  <Input
                    type="url"
                    value={formData.repo_url}
                    onChange={(e) => setFormData(prev => ({ ...prev, repo_url: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Live URL (Optional)</Label>
                  <Input
                    type="url"
                    value={formData.live_url}
                    onChange={(e) => setFormData(prev => ({ ...prev, live_url: e.target.value }))}
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Status</Label>
                  <Select
                    value={formData.status}
                    onChange={(e) => setFormData(prev => ({ ...prev, status: e.target.value as any }))}
                  >
                    <option value="draft">Draft</option>
                    <option value="active">Active</option>
                    <option value="archived">Archived</option>
                  </Select>
                </FormGroup>

                <FormGroup>
                  <Label>Tags (press Enter or comma to add)</Label>
                  <TagInput>
                    {formData.tags.map((tag) => (
                      <TagItem key={tag}>
                        {tag}
                        <RemoveTagButton type="button" onClick={() => removeTag(tag)}>
                          ×
                        </RemoveTagButton>
                      </TagItem>
                    ))}
                    <TagInputField
                      type="text"
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyDown={handleAddTag}
                      placeholder="Add tags..."
                    />
                  </TagInput>
                </FormGroup>

                <ModalActions>
                  <ModalButton 
                    type="button" 
                    variant="secondary" 
                    onClick={() => {
                      setShowModal(false);
                      resetForm();
                    }}
                  >
                    Cancel
                  </ModalButton>
                  <ModalButton type="submit" variant="primary">
                    {editingProject ? 'Update' : 'Create'}
                  </ModalButton>
                </ModalActions>
              </Form>
            </ModalContent>
          </Modal>
        )}
      </AnimatePresence>
    </Container>
  );
};

export default ProjectManager;