import React, { useState, useEffect, useCallback, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import Badge from '../common/Badge';
import { useScrollTo } from '../../hooks/useScrollTo';
import { ScrollableModal } from '../common/ScrollableModal';
import { useInlineFormScroll } from '../../hooks/useInlineFormScroll';
import { useCrudOperations } from '../../hooks/useCrudOperations';
import {
  Prompt,
  Category,
  PromptFormData,
  CategoryFormData,
  defaultPromptFormData,
  defaultCategoryFormData,
  mapPromptToFormData,
  mapCategoryToFormData,
} from '../../utils/promptApi'; // This will be created next

// Styled Components (mostly copied from CertificationsManager)
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

const ButtonGroup = styled.div`
  display: flex;
  gap: 1rem;
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

const FilterSection = styled.div`
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  
  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

const SearchInput = styled.input`
  flex: 1;
  min-width: 200px;
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const FilterSelect = styled.select`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const PromptGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  
  @media (max-width: 1024px) {
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 1.25rem;
  }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  @media (max-width: 480px) {
    gap: 0.75rem;
  }
`;

const PromptCard = styled(motion.div)`
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 12px;
  padding: 1.5rem;
  position: relative;
  transition: all 0.3s ease;
  min-height: 280px;
  display: flex;
  flex-direction: column;
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
    border-color: ${({ theme }) => theme.colors.primary}30;
  }
  
  @media (max-width: 768px) {
    padding: 1.25rem;
    min-height: 260px;
    border-radius: 10px;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    }
  }
  
  @media (max-width: 480px) {
    padding: 1rem;
    min-height: 240px;
  }
`;

const CardHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
`;

const CardTitle = styled.h3`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.1rem;
`;

const CardActions = styled.div`
  display: flex;
  gap: 0.5rem;
`;

const IconButton = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.textSecondary};
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  transition: all 0.2s ease;
  
  &:hover {
    color: ${({ theme }) => theme.colors.primary};
    background: ${({ theme }) => theme.colors.border};
  }
`;

const CardBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  flex: 1;
  
  @media (max-width: 768px) {
    gap: 0.6rem;
  }
`;

const InfoSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  
  @media (max-width: 768px) {
    gap: 0.4rem;
  }
`;

const BadgeContainer = styled.div`
  display: flex;
  gap: 0.6rem;
  margin-top: 1rem;
  padding-top: 0.75rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border}20;
  flex-wrap: wrap;
  align-items: center;
  min-height: 2rem;
  
  @media (max-width: 768px) {
    gap: 0.4rem;
    margin-top: 0.75rem;
    padding-top: 0.5rem;
  }
`;

const CardInfo = styled.div`
  font-size: 0.9rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const ModalHeader = styled.h3`
  margin: 0 0 1.5rem 0;
  color: ${({ theme }) => theme.colors.text};
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const FormRow = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
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
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const TextArea = styled.textarea`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  min-height: 100px;
  resize: vertical;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const Select = styled.select`
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const CheckboxGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const Checkbox = styled.input`
  width: 18px;
  height: 18px;
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

const LoadingSpinner = styled.div`
  text-align: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const DeleteConfirmText = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 1.5rem;
`;

// Component
const PromptsManager: React.FC = () => {
  const promptsCrud = useCrudOperations<Prompt, PromptFormData>(
    {
      baseUrl: '/api/v1/devpanel/prompts',
      resourceName: 'prompts',
      requiresAuth: true,
    },
    defaultPromptFormData,
    {
      onItemToFormData: mapPromptToFormData,
      fetchQueryParams: 'include_hidden=true',
    }
  );

  const categoriesCrud = useCrudOperations<Category, CategoryFormData>(
    {
      baseUrl: '/api/v1/devpanel/prompt-categories',
      resourceName: 'categories',
      requiresAuth: true,
    },
    defaultCategoryFormData,
    {
      onItemToFormData: mapCategoryToFormData,
      fetchQueryParams: 'include_hidden=true',
    }
  );

  const [filteredPrompts, setFilteredPrompts] = useState<Prompt[]>([]);
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deletingPrompt, setDeletingPrompt] = useState<Prompt | null>(null);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [visibilityFilter, setVisibilityFilter] = useState('all');
  
  const { scrollToElement } = useScrollTo();
  const { formRef: errorRef } = useInlineFormScroll(!!(promptsCrud.error || categoriesCrud.error), {
    scrollOffset: 100,
    scrollDelay: 200
  });

  useEffect(() => {
    promptsCrud.fetchItems();
    categoriesCrud.fetchItems();
  }, []);

  useEffect(() => {
    let filtered = promptsCrud.items;
    
    if (searchTerm) {
      filtered = filtered.filter(prompt => 
        prompt.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        prompt.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        prompt.prompt.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    if (categoryFilter !== 'all') {
      filtered = filtered.filter(prompt => prompt.category_id === categoryFilter);
    }
    
    if (visibilityFilter === 'visible') {
      filtered = filtered.filter(prompt => prompt.is_visible);
    } else if (visibilityFilter === 'hidden') {
      filtered = filtered.filter(prompt => !prompt.is_visible);
    }
    
    setFilteredPrompts(filtered);
  }, [promptsCrud.items, searchTerm, categoryFilter, visibilityFilter]);

  const handlePromptSubmit = promptsCrud.handleSubmit;
  const handleCategorySubmit = categoriesCrud.handleSubmit;

  const handleDelete = async () => {
    if (!deletingPrompt) return;
    
    try {
      await promptsCrud.deleteItem(deletingPrompt.id);
      setShowDeleteConfirm(false);
      setDeletingPrompt(null);
    } catch (err) {
      console.error('Delete failed:', err);
    }
  };

  const handleEdit = (prompt: Prompt) => {
    promptsCrud.handleEdit(prompt);
  };

  if (promptsCrud.loading) return <LoadingSpinner>Loading prompts...</LoadingSpinner>;

  const error = promptsCrud.error || categoriesCrud.error;

  return (
    <Container>
      {error && (
        <ErrorMessage ref={errorRef}>
          {error}
        </ErrorMessage>
      )}
      
      <Header>
        <Title>Prompts Management</Title>
        <ButtonGroup>
          <AddButton onClick={() => {
            categoriesCrud.resetForm();
            setShowCategoryModal(true);
          }}>
            Manage Categories
          </AddButton>
          <AddButton onClick={promptsCrud.handleCreate}>
            Add Prompt
          </AddButton>
        </ButtonGroup>
      </Header>

      {error && <ErrorMessage ref={errorRef}>{error}</ErrorMessage>}

      <FilterSection>
        <SearchInput
          type="text"
          placeholder="Search prompts..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <FilterSelect
          value={categoryFilter}
          onChange={(e) => setCategoryFilter(e.target.value)}
        >
          <option value="all">All Categories</option>
          {categoriesCrud.items.map(cat => (
            <option key={cat.id} value={cat.id}>
              {cat.name}
            </option>
          ))}
        </FilterSelect>
        <FilterSelect
          value={visibilityFilter}
          onChange={(e) => setVisibilityFilter(e.target.value)}
        >
          <option value="all">All</option>
          <option value="visible">Visible Only</option>
          <option value="hidden">Hidden Only</option>
        </FilterSelect>
      </FilterSection>

      <PromptGrid>
        <AnimatePresence>
          {filteredPrompts.map(prompt => (
            <PromptCard
              key={prompt.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <CardHeader>
                <CardTitle>{prompt.name}</CardTitle>
                <CardActions>
                  <IconButton onClick={() => handleEdit(prompt)}>✏️</IconButton>
                  <IconButton onClick={() => {
                    setDeletingPrompt(prompt);
                    setShowDeleteConfirm(true);
                  }}>🗑️</IconButton>
                </CardActions>
              </CardHeader>
              
              <CardBody>
                <InfoSection>
                  {prompt.description && (
                    <CardInfo>{prompt.description}</CardInfo>
                  )}
                </InfoSection>
                
                <BadgeContainer>
                  {prompt.is_featured && <Badge variant="featured">Featured</Badge>}
                  {prompt.category && <Badge variant="category">{prompt.category.name}</Badge>}
                  {!prompt.is_visible && <Badge variant="status">Hidden</Badge>}
                </BadgeContainer>
              </CardBody>
            </PromptCard>
          ))}
        </AnimatePresence>
      </PromptGrid>

      <ScrollableModal 
        isOpen={promptsCrud.showModal} 
        onClose={() => promptsCrud.setShowModal(false)}
      >
        <ModalHeader>
          {promptsCrud.editingItem ? 'Edit Prompt' : 'Add New Prompt'}
        </ModalHeader>
            
            <Form onSubmit={handlePromptSubmit}>
              <FormRow>
                <FormGroup>
                  <Label>Prompt Name *</Label>
                  <Input
                    type="text"
                    value={promptsCrud.formData.name}
                    onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, name: e.target.value })}
                    required
                  />
                </FormGroup>
                
                <FormGroup>
                  <Label>Category</Label>
                  <Select
                    value={promptsCrud.formData.category_id}
                    onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, category_id: e.target.value })}
                  >
                    <option value="">No Category</option>
                    {categoriesCrud.items.map(cat => (
                      <option key={cat.id} value={cat.id}>
                        {cat.name}
                      </option>
                    ))}
                  </Select>
                </FormGroup>
              </FormRow>
              
              <FormGroup>
                <Label>Description</Label>
                <TextArea
                  value={promptsCrud.formData.description}
                  onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, description: e.target.value })}
                  placeholder="Brief description of the prompt..."
                />
              </FormGroup>

              <FormGroup>
                <Label>Prompt *</Label>
                <TextArea
                  value={promptsCrud.formData.prompt}
                  onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, prompt: e.target.value })}
                  placeholder="Enter the prompt here..."
                  required
                />
              </FormGroup>
              
              <FormRow>
                <FormGroup>
                  <Label>Sort Order</Label>
                  <Input
                    type="number"
                    value={promptsCrud.formData.sort_order}
                    onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, sort_order: parseInt(e.target.value) })}
                  />
                </FormGroup>
                
                <FormGroup style={{ justifyContent: 'center' }}>
                  <CheckboxGroup>
                    <Checkbox
                      type="checkbox"
                      id="featured"
                      checked={promptsCrud.formData.is_featured}
                      onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, is_featured: e.target.checked })}
                    />
                    <Label htmlFor="featured">Featured</Label>
                  </CheckboxGroup>
                  
                  <CheckboxGroup>
                    <Checkbox
                      type="checkbox"
                      id="visible"
                      checked={promptsCrud.formData.is_visible}
                      onChange={(e) => promptsCrud.setFormData({ ...promptsCrud.formData, is_visible: e.target.checked })}
                    />
                    <Label htmlFor="visible">Visible</Label>
                  </CheckboxGroup>
                </FormGroup>
              </FormRow>
              
              <ModalActions>
                <ModalButton type="button" variant="secondary" onClick={() => promptsCrud.setShowModal(false)}>
                  Cancel
                </ModalButton>
                <ModalButton type="submit" variant="primary">
                  {promptsCrud.editingItem ? 'Update' : 'Add'} Prompt
                </ModalButton>
              </ModalActions>
            </Form>
      </ScrollableModal>

      <ScrollableModal 
        isOpen={showCategoryModal} 
        onClose={() => setShowCategoryModal(false)}
      >
        <ModalHeader>Manage Categories</ModalHeader>
            
            <Form onSubmit={handleCategorySubmit}>
              <FormGroup>
                <Label>Category Name *</Label>
                <Input
                  type="text"
                  value={categoriesCrud.formData.name}
                  onChange={(e) => categoriesCrud.setFormData({ ...categoriesCrud.formData, name: e.target.value })}
                  required
                />
              </FormGroup>
              
              <FormGroup>
                <Label>Description</Label>
                <TextArea
                  value={categoriesCrud.formData.description}
                  onChange={(e) => categoriesCrud.setFormData({ ...categoriesCrud.formData, description: e.target.value })}
                />
              </FormGroup>
              
              <FormRow>
                <FormGroup>
                  <Label>Sort Order</Label>
                  <Input
                    type="number"
                    value={categoriesCrud.formData.sort_order}
                    onChange={(e) => categoriesCrud.setFormData({ ...categoriesCrud.formData, sort_order: parseInt(e.target.value) })}
                  />
                </FormGroup>
                
                <FormGroup style={{ justifyContent: 'center' }}>
                  <CheckboxGroup>
                    <Checkbox
                      type="checkbox"
                      id="catVisible"
                      checked={categoriesCrud.formData.is_visible}
                      onChange={(e) => categoriesCrud.setFormData({ ...categoriesCrud.formData, is_visible: e.target.checked })}
                    />
                    <Label htmlFor="catVisible">Visible</Label>
                  </CheckboxGroup>
                </FormGroup>
              </FormRow>
              
              <ModalActions>
                <ModalButton type="button" variant="secondary" onClick={() => setShowCategoryModal(false)}>
                  Close
                </ModalButton>
                <ModalButton type="submit" variant="primary">
                  {categoriesCrud.editingItem ? 'Update' : 'Add'} Category
                </ModalButton>
              </ModalActions>
            </Form>
            
            <div style={{ marginTop: '2rem' }}>
              <h4>Existing Categories</h4>
              {categoriesCrud.items.map(cat => (
                <div key={cat.id} style={{ padding: '0.5rem', borderBottom: '1px solid #ddd' }}>
                  <strong>{cat.name}</strong>
                  {cat.description && <p style={{ fontSize: '0.9rem', margin: '0.25rem 0' }}>{cat.description}</p>}
                  <small>Sort Order: {cat.sort_order} | {cat.is_visible ? 'Visible' : 'Hidden'}</small>
                </div>
              ))}
            </div>
      </ScrollableModal>

      <ScrollableModal 
        isOpen={showDeleteConfirm && !!deletingPrompt} 
        onClose={() => setShowDeleteConfirm(false)}
      >
        <ModalHeader>Confirm Deletion</ModalHeader>
            <DeleteConfirmText>
              Are you sure you want to delete the prompt "{deletingPrompt?.name}"? This action cannot be undone.
            </DeleteConfirmText>
            <ModalActions>
              <ModalButton variant="secondary" onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </ModalButton>
              <ModalButton variant="primary" onClick={handleDelete}>
                Delete
              </ModalButton>
            </ModalActions>
      </ScrollableModal>
    </Container>
  );
};

export default PromptsManager;
