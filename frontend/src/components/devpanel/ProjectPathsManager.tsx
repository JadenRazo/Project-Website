import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useScrollTo } from '../../hooks/useScrollTo';
import { ScrollableModal } from '../common/ScrollableModal';
import { useInlineFormScroll } from '../../hooks/useInlineFormScroll';
import { useCrudOperations } from '../../hooks/useCrudOperations';
import {
  ProjectPath,
  ProjectPathFormData,
  defaultProjectPathFormData,
  mapProjectPathToFormData,
  projectPathsApi,
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

const RefreshButton = styled.button<{ refreshing: boolean }>`
  padding: 0.75rem 1rem;
  background: ${({ theme }) => theme.colors.success};
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: ${({ refreshing }) => refreshing ? 'not-allowed' : 'pointer'};
  transition: all 0.2s ease;
  opacity: ${({ refreshing }) => refreshing ? 0.7 : 1};
  
  &:hover {
    background: ${({ theme, refreshing }) => 
      refreshing ? theme.colors.success : theme.colors.success + 'dd'};
    transform: ${({ refreshing }) => refreshing ? 'none' : 'translateY(-1px)'};
  }
  
  &:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }
`;

const HeaderButtons = styled.div`
  display: flex;
  gap: 0.75rem;
  align-items: center;
  
  @media (max-width: 768px) {
    width: 100%;
    justify-content: space-between;
  }
`;

const PathGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
`;

const PathCard = styled(motion.div)`
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

const PathHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
`;

const PathName = styled.h3`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.25rem;
`;

const StatusToggle = styled.button<{ active: boolean }>`
  padding: 0.25rem 0.75rem;
  border: none;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, active }) => 
    active ? theme.colors.success + '20' : theme.colors.error + '20'
  };
  color: ${({ theme, active }) => 
    active ? theme.colors.success : theme.colors.error
  };
  
  &:hover {
    opacity: 0.8;
  }
`;

const PathInfo = styled.div`
  margin-bottom: 1rem;
`;

const PathLocation = styled.div`
  font-family: 'Courier New', monospace;
  background: ${({ theme }) => theme.colors.background};
  padding: 0.5rem;
  border-radius: 4px;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: ${({ theme }) => theme.colors.text};
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const PathDescription = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 1rem;
  line-height: 1.5;
  font-size: 0.875rem;
`;

const PatternsContainer = styled.div`
  margin-bottom: 1rem;
`;

const PatternsLabel = styled.div`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
`;

const PatternsList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
`;

const Pattern = styled.span`
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.warning + '20'};
  color: ${({ theme }) => theme.colors.warning};
  border-radius: 4px;
  font-size: 0.75rem;
  font-family: 'Courier New', monospace;
`;

const PathActions = styled.div`
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
`;

const ActionButton = styled.button<{ variant: 'edit' | 'delete' }>`
  flex: 1;
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, variant }) => 
    variant === 'edit' ? theme.colors.primary : theme.colors.error
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
  min-height: 80px;
  resize: vertical;
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const CheckboxContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const Checkbox = styled.input`
  width: 18px;
  height: 18px;
  cursor: pointer;
`;

const PatternInput = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  min-height: 50px;
`;

const PatternInputField = styled.input`
  border: none;
  background: none;
  outline: none;
  color: ${({ theme }) => theme.colors.text};
  flex: 1;
  min-width: 100px;
`;

const PatternItem = styled.span`
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.warning + '20'};
  color: ${({ theme }) => theme.colors.warning};
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: 'Courier New', monospace;
`;

const RemovePatternButton = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.warning};
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  
  &:hover {
    background: ${({ theme }) => theme.colors.warning + '30'};
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

const EmptyState = styled.div`
  text-align: center;
  padding: 3rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

export const ProjectPathsManager: React.FC = () => {
  // Use CRUD operations hook for project paths
  const pathsCrud = useCrudOperations<ProjectPath, ProjectPathFormData>(
    {
      baseUrl: '/api/v1/code/paths',
      resourceName: 'project paths',
      requiresAuth: false,
    },
    defaultProjectPathFormData,
    {
      onItemToFormData: mapProjectPathToFormData,
      extractItemsFromResponse: (data) => data.data || [],
    }
  );

  // Local state for UI
  const [patternInput, setPatternInput] = useState('');
  const [refreshingStats, setRefreshingStats] = useState(false);
  const [pendingRefreshes, setPendingRefreshes] = useState(new Set<string>());
  
  // Scroll hooks and refs
  const { scrollToElement } = useScrollTo();
  const { formRef: errorRef } = useInlineFormScroll(!!pathsCrud.error, {
    scrollOffset: 100,
    scrollDelay: 200
  });

  useEffect(() => {
    pathsCrud.fetchItems();
  }, []);

  // Removed fetchPaths - now handled by CRUD hook

  const refreshCodeStats = async (action?: string) => {
    const requestKey = action || 'default';
    
    // Check if request is already pending
    if (pendingRefreshes.has(requestKey)) {
      console.log('Code stats refresh already in progress for:', requestKey);
      return;
    }

    try {
      // Mark request as pending
      setPendingRefreshes(prev => new Set(prev).add(requestKey));
      setRefreshingStats(true);
      
      await projectPathsApi.refreshStats(action);
      console.log('Code statistics refreshed successfully' + (action ? ` after ${action}` : ''));
    } catch (err) {
      console.warn('Error refreshing code stats:', err);
      // Don't throw error - stats refresh is non-critical
    } finally {
      // Remove from pending requests
      setPendingRefreshes(prev => {
        const newSet = new Set(prev);
        newSet.delete(requestKey);
        return newSet;
      });
      
      // Only clear refreshingStats if no other requests are pending
      if (pendingRefreshes.size <= 1) { // Will be 1 because we haven't removed it yet
        setRefreshingStats(false);
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (pathsCrud.editingItem) {
        await pathsCrud.updateItem(pathsCrud.editingItem.id, pathsCrud.formData);
        refreshCodeStats('updating project path');
      } else {
        await pathsCrud.createItem(pathsCrud.formData);
        refreshCodeStats('creating project path');
      }
      pathsCrud.setShowModal(false);
      resetForm();
    } catch (err) {
      console.error('Error saving project path:', err);
      pathsCrud.setError(err instanceof Error ? err.message : 'Failed to save project path');
    }
  };

  const handleDelete = async (pathId: string) => {
    try {
      await pathsCrud.deleteItem(pathId);
      refreshCodeStats('deleting project path');
    } catch (err) {
      console.error('Error deleting project path:', err);
      pathsCrud.setError('Failed to delete project path');
    }
  };

  const handleToggleActive = async (path: ProjectPath) => {
    try {
      const updatedData = {
        name: path.name,
        path: path.path,
        description: path.description || '',
        exclude_patterns: path.exclude_patterns,
        is_active: !path.is_active
      };
      await pathsCrud.updateItem(path.id, updatedData);
      refreshCodeStats('toggling project path status');
    } catch (err) {
      console.error('Error updating project path:', err);
      pathsCrud.setError('Failed to update project path');
    }
  };

  const handleEdit = pathsCrud.handleEdit;

  const resetForm = () => {
    pathsCrud.resetForm();
    setPatternInput('');
  };

  const handleAddPattern = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      const pattern = patternInput.trim();
      if (pattern && !pathsCrud.formData.exclude_patterns.includes(pattern)) {
        pathsCrud.setFormData(prev => ({
          ...prev,
          exclude_patterns: [...prev.exclude_patterns, pattern]
        }));
      }
      setPatternInput('');
    }
  };

  const removePattern = (patternToRemove: string) => {
    pathsCrud.setFormData(prev => ({
      ...prev,
      exclude_patterns: prev.exclude_patterns.filter(pattern => pattern !== patternToRemove)
    }));
  };

  if (pathsCrud.loading) {
    return <LoadingSpinner>Loading project paths...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <Title>Lines of Code Project Paths</Title>
        <HeaderButtons>
          <RefreshButton 
            refreshing={refreshingStats}
            onClick={() => refreshCodeStats('manual refresh')}
            disabled={refreshingStats}
          >
            {refreshingStats ? '⟳ Refreshing...' : '↻ Refresh Stats'}
          </RefreshButton>
          <AddButton onClick={pathsCrud.handleCreate}>
            Add New Path
          </AddButton>
        </HeaderButtons>
      </Header>

      {pathsCrud.error && <ErrorMessage ref={errorRef}>{pathsCrud.error}</ErrorMessage>}

      {pathsCrud.items.length === 0 ? (
        <EmptyState>
          <p>No project paths configured yet.</p>
          <p>Add your first project path to start tracking lines of code.</p>
        </EmptyState>
      ) : (
        <PathGrid>
          <AnimatePresence>
            {pathsCrud.items.map((path) => (
              <PathCard
                key={path.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                layout
              >
                <PathHeader>
                  <PathName>{path.name}</PathName>
                  <StatusToggle 
                    active={path.is_active}
                    onClick={() => handleToggleActive(path)}
                  >
                    {path.is_active ? 'Active' : 'Inactive'}
                  </StatusToggle>
                </PathHeader>

                <PathInfo>
                  <PathLocation>{path.path}</PathLocation>
                  {path.description && (
                    <PathDescription>{path.description}</PathDescription>
                  )}
                </PathInfo>

                {path.exclude_patterns && path.exclude_patterns.length > 0 && (
                  <PatternsContainer>
                    <PatternsLabel>Excluded Patterns:</PatternsLabel>
                    <PatternsList>
                      {path.exclude_patterns.map((pattern, index) => (
                        <Pattern key={index}>{pattern}</Pattern>
                      ))}
                    </PatternsList>
                  </PatternsContainer>
                )}

                <PathActions>
                  <ActionButton variant="edit" onClick={() => handleEdit(path)}>
                    Edit
                  </ActionButton>
                  <ActionButton variant="delete" onClick={() => handleDelete(path.id)}>
                    Delete
                  </ActionButton>
                </PathActions>
              </PathCard>
            ))}
          </AnimatePresence>
        </PathGrid>
      )}

      <AnimatePresence>
        <ScrollableModal 
          isOpen={pathsCrud.showModal} 
          onClose={() => pathsCrud.setShowModal(false)}
        >
          <h3 style={{ marginBottom: '1.5rem' }}>{pathsCrud.editingItem ? 'Edit Project Path' : 'Add New Project Path'}</h3>
              
              <Form onSubmit={handleSubmit}>
                <FormGroup>
                  <Label>Project Name</Label>
                  <Input
                    type="text"
                    value={pathsCrud.formData.name}
                    onChange={(e) => pathsCrud.setFormData(prev => ({ ...prev, name: e.target.value }))}
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>File Path</Label>
                  <Input
                    type="text"
                    value={pathsCrud.formData.path}
                    onChange={(e) => pathsCrud.setFormData(prev => ({ ...prev, path: e.target.value }))}
                    placeholder="/path/to/your/project"
                    required
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Description (Optional)</Label>
                  <TextArea
                    value={pathsCrud.formData.description}
                    onChange={(e) => pathsCrud.setFormData(prev => ({ ...prev, description: e.target.value }))}
                    placeholder="Brief description of this project..."
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Exclude Patterns (press Enter or comma to add)</Label>
                  <PatternInput>
                    {pathsCrud.formData.exclude_patterns.map((pattern, index) => (
                      <PatternItem key={index}>
                        {pattern}
                        <RemovePatternButton type="button" onClick={() => removePattern(pattern)}>
                          ×
                        </RemovePatternButton>
                      </PatternItem>
                    ))}
                    <PatternInputField
                      type="text"
                      value={patternInput}
                      onChange={(e) => setPatternInput(e.target.value)}
                      onKeyDown={handleAddPattern}
                      placeholder="node_modules, *.log, build..."
                    />
                  </PatternInput>
                </FormGroup>

                <FormGroup>
                  <CheckboxContainer>
                    <Checkbox
                      type="checkbox"
                      checked={pathsCrud.formData.is_active}
                      onChange={(e) => pathsCrud.setFormData(prev => ({ ...prev, is_active: e.target.checked }))}
                    />
                    <Label>Active (include in line count)</Label>
                  </CheckboxContainer>
                </FormGroup>

                <ModalActions>
                  <ModalButton 
                    type="button" 
                    variant="secondary" 
                    onClick={() => {
                      pathsCrud.setShowModal(false);
                      resetForm();
                    }}
                  >
                    Cancel
                  </ModalButton>
                  <ModalButton type="submit" variant="primary">
                    {pathsCrud.editingItem ? 'Update' : 'Create'}
                  </ModalButton>
                </ModalActions>
              </Form>
        </ScrollableModal>
      </AnimatePresence>
    </Container>
  );
};

export default ProjectPathsManager;