import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useScrollTo } from '../../hooks/useScrollTo';
import { ScrollableModal } from '../common/ScrollableModal';
import { useInlineFormScroll } from '../../hooks/useInlineFormScroll';
import { useCrudOperations } from '../../hooks/useCrudOperations';
import {
  Project,
  ProjectFormData,
  defaultProjectFormData,
  mapProjectToFormData,
} from '../../utils/devPanelApi';

// Styled Components
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
  animation: slideIn 0.3s ease-out;
  scroll-margin-top: 100px;
  
  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
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

// Modal styles - Removed in favor of ScrollableModal component

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

// Removed duplicate ErrorMessage - using the one defined above

const LoadingSpinner = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

// Main Component
export const ProjectManager: React.FC = () => {
  // Use CRUD operations hook for projects
  const projectsCrud = useCrudOperations<Project, ProjectFormData>(
    {
      baseUrl: '/api/v1/projects',
      resourceName: 'projects',
      requiresAuth: false,
    },
    defaultProjectFormData,
    {
      onItemToFormData: mapProjectToFormData,
      extractItemsFromResponse: (data) => data.projects || [],
    }
  );

  // Local state for UI
  const [tagInput, setTagInput] = useState('');
  const { scrollToElement } = useScrollTo();
  const { formRef: errorRef } = useInlineFormScroll(!!projectsCrud.error, {
    scrollOffset: 100,
    scrollDelay: 200
  });

  useEffect(() => {
    projectsCrud.fetchItems();
  }, []);

  const handleSubmit = projectsCrud.handleSubmit;

  const handleDelete = projectsCrud.handleDelete;

  const handleEdit = projectsCrud.handleEdit;
  
  const resetForm = () => {
    projectsCrud.resetForm();
    setTagInput('');
  };

  const handleAddTag = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      const tag = tagInput.trim();
      if (tag && !projectsCrud.formData.tags.includes(tag)) {
        projectsCrud.setFormData(prev => ({
          ...prev,
          tags: [...prev.tags, tag]
        }));
      }
      setTagInput('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    projectsCrud.setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }));
  };

  if (projectsCrud.loading) {
    return <LoadingSpinner>Loading projects...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <Title>Project Management</Title>
        <AddButton onClick={projectsCrud.handleCreate}>
          Add New Project
        </AddButton>
      </Header>

      {projectsCrud.error && <ErrorMessage ref={errorRef}>{projectsCrud.error}</ErrorMessage>}

      <ProjectGrid>
        <AnimatePresence>
          {projectsCrud.items.map((project) => (
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
        <ScrollableModal 
          isOpen={projectsCrud.showModal} 
          onClose={() => projectsCrud.setShowModal(false)}
        >
          <h3 style={{ marginBottom: '1.5rem' }}>{projectsCrud.editingItem ? 'Edit Project' : 'Add New Project'}</h3>
              
              <Form onSubmit={handleSubmit}>
                <FormGroup>
                  <Label>Project Name</Label>
                  <Input
                    type="text"
                    value={projectsCrud.formData.name}
                    onChange={(e) => projectsCrud.setFormData(prev => ({ ...prev, name: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Description</Label>
                  <TextArea
                    value={projectsCrud.formData.description}
                    onChange={(e) => projectsCrud.setFormData(prev => ({ ...prev, description: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Repository URL</Label>
                  <Input
                    type="url"
                    value={projectsCrud.formData.repo_url}
                    onChange={(e) => projectsCrud.setFormData(prev => ({ ...prev, repo_url: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Live URL (Optional)</Label>
                  <Input
                    type="url"
                    value={projectsCrud.formData.live_url}
                    onChange={(e) => projectsCrud.setFormData(prev => ({ ...prev, live_url: e.target.value }))}
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Status</Label>
                  <Select
                    value={projectsCrud.formData.status}
                    onChange={(e) => projectsCrud.setFormData(prev => ({ ...prev, status: e.target.value as any }))}
                  >
                    <option value="draft">Draft</option>
                    <option value="active">Active</option>
                    <option value="archived">Archived</option>
                  </Select>
                </FormGroup>

                <FormGroup>
                  <Label>Tags (press Enter or comma to add)</Label>
                  <TagInput>
                    {projectsCrud.formData.tags.map((tag) => (
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
                      projectsCrud.setShowModal(false);
                      resetForm();
                    }}
                  >
                    Cancel
                  </ModalButton>
                  <ModalButton type="submit" variant="primary">
                    {projectsCrud.editingItem ? 'Update' : 'Create'}
                  </ModalButton>
                </ModalActions>
              </Form>
        </ScrollableModal>
      </AnimatePresence>
    </Container>
  );
};

export default ProjectManager;