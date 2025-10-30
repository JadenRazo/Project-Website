import React, { useState, useEffect, useCallback, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useScrollTo } from '../../hooks/useScrollTo';
import { ScrollableModal } from '../common/ScrollableModal';
import { useInlineFormScroll } from '../../hooks/useInlineFormScroll';
import { useCrudOperations } from '../../hooks/useCrudOperations';
import {
  Skill,
  SkillFormData,
  defaultSkillFormData,
  mapSkillToFormData,
  skillsApi,
} from '../../utils/devPanelApi';

const CATEGORIES = [
  { value: 'frontend', label: 'Frontend' },
  { value: 'backend', label: 'Backend' },
  { value: 'design', label: 'Design' },
  { value: 'database', label: 'Database' },
  { value: 'devops', label: 'DevOps' },
  { value: 'language', label: 'Language' },
  { value: 'framework', label: 'Framework' },
  { value: 'tool', label: 'Tool' }
];

const PROFICIENCY_LEVELS = [
  { value: 'beginner', label: 'Beginner' },
  { value: 'intermediate', label: 'Intermediate' },
  { value: 'advanced', label: 'Advanced' },
  { value: 'expert', label: 'Expert' }
];

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

const FilterSection = styled.div`
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  
  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

const FilterSelect = styled.select`
  padding: 0.5rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  min-width: 150px;
`;

const SkillGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
`;

const SkillCard = styled(motion.div)`
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

const SkillHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
`;

const SkillName = styled.h3`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const SkillIcon = styled.div<{ color?: string }>`
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: ${({ color }) => color || '#gray'};
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: white;
  font-weight: bold;
`;

const SkillCategory = styled.span<{ category: string }>`
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 4px;
  text-transform: uppercase;
  background: ${({ category }) => {
    const colors: { [key: string]: string } = {
      frontend: '#61DAFB20',
      backend: '#00ADD820',
      design: '#F24E1E20',
      database: '#33679120',
      devops: '#2496ED20',
      language: '#F7DF1E20',
      framework: '#00000020',
      tool: '#FCC62420'
    };
    return colors[category] || '#E5E7EB20';
  }};
  color: ${({ category }) => {
    const colors: { [key: string]: string } = {
      frontend: '#61DAFB',
      backend: '#00ADD8',
      design: '#F24E1E',
      database: '#336791',
      devops: '#2496ED',
      language: '#F7DF1E',
      framework: '#000000',
      tool: '#FCC624'
    };
    return colors[category] || '#6B7280';
  }};
`;

const SkillDescription = styled.p`
  margin: 0 0 1rem 0;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  line-height: 1.4;
`;

const ProficiencySection = styled.div`
  margin-bottom: 1rem;
`;

const ProficiencyHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
`;

const ProficiencyLevel = styled.span`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
  text-transform: capitalize;
`;

const ProficiencyValue = styled.span`
  font-weight: 600;
  color: ${({ theme }) => theme.colors.primary};
`;

const ProficiencyBar = styled.div`
  width: 100%;
  height: 8px;
  background: ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  overflow: hidden;
`;

const ProficiencyFill = styled.div<{ width: number }>`
  height: 100%;
  width: ${({ width }) => width}%;
  background: linear-gradient(90deg, #10B981, #059669);
  transition: width 0.3s ease;
`;

const SkillMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
`;

const FeaturedBadge = styled.span<{ featured: boolean }>`
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 4px;
  background: ${({ featured }) => featured ? '#10B98120' : '#E5E7EB20'};
  color: ${({ featured }) => featured ? '#10B981' : '#6B7280'};
`;

const SortOrder = styled.span`
  font-size: 0.75rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const TagContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
`;

const Tag = styled.span`
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.border};
  color: ${({ theme }) => theme.colors.textSecondary};
  border-radius: 4px;
  font-size: 0.75rem;
`;

const ActionButtons = styled.div`
  display: flex;
  gap: 0.5rem;
`;

const ActionButton = styled.button<{ variant: 'edit' | 'delete' | 'featured' }>`
  padding: 0.5rem 0.75rem;
  border: none;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  
  background: ${({ theme, variant }) => {
    const colors = {
      edit: theme.colors.primary,
      delete: theme.colors.error || '#EF4444',
      featured: '#10B981'
    };
    return colors[variant];
  }};
  color: white;
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
`;

// Modal styles - Removed in favor of ScrollableModal component

const ModalHeader = styled.h3`
  margin: 0 0 1.5rem 0;
  color: ${({ theme }) => theme.colors.text};
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

const RangeGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const RangeInput = styled.input`
  width: 100%;
`;

const RangeValue = styled.span`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.primary};
`;

const TagInputGroup = styled.div`
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
`;

const TagInput = styled.input`
  flex: 1;
  min-width: 200px;
  padding: 0.5rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
`;

const AddTagButton = styled.button`
  padding: 0.5rem 1rem;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  
  &:hover {
    opacity: 0.9;
  }
`;

const TagList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
`;

const EditableTag = styled.span`
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  background: ${({ theme }) => theme.colors.border};
  color: ${({ theme }) => theme.colors.textSecondary};
  border-radius: 4px;
  font-size: 0.75rem;
`;

const RemoveTagButton = styled.button`
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.error || '#EF4444'};
  cursor: pointer;
  padding: 0;
  font-size: 0.75rem;
  
  &:hover {
    opacity: 0.7;
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
export const SkillsManager: React.FC = () => {
  // Use CRUD operations hook for skills
  const skillsCrud = useCrudOperations<Skill, SkillFormData>(
    {
      baseUrl: '/api/v1/devpanel/skills',
      resourceName: 'skills',
      requiresAuth: true,
    },
    defaultSkillFormData,
    {
      onItemToFormData: mapSkillToFormData,
      extractItemsFromResponse: (data) => data.skills || [],
    }
  );

  // Local state for UI
  const [filteredSkills, setFilteredSkills] = useState<Skill[]>([]);
  const [categoryFilter, setCategoryFilter] = useState('');
  const [proficiencyFilter, setProficiencyFilter] = useState('');
  const [featuredFilter, setFeaturedFilter] = useState('');
  const [tagInput, setTagInput] = useState('');
  
  // Scroll hooks and refs
  const { scrollToElement } = useScrollTo();
  const { formRef: errorRef } = useInlineFormScroll(!!skillsCrud.error, {
    scrollOffset: 100,
    scrollDelay: 200
  });

  useEffect(() => {
    skillsCrud.fetchItems();
  }, []);

  const filterSkills = useCallback(() => {
    let filtered = [...skillsCrud.items];

    if (categoryFilter) {
      filtered = filtered.filter(skill => skill.category === categoryFilter);
    }

    if (proficiencyFilter) {
      filtered = filtered.filter(skill => skill.proficiency_level === proficiencyFilter);
    }

    if (featuredFilter === 'featured') {
      filtered = filtered.filter(skill => skill.is_featured);
    } else if (featuredFilter === 'not-featured') {
      filtered = filtered.filter(skill => !skill.is_featured);
    }

    // Sort by sort_order, then by name
    filtered.sort((a, b) => {
      if (a.sort_order !== b.sort_order) {
        return a.sort_order - b.sort_order;
      }
      return a.name.localeCompare(b.name);
    });

    setFilteredSkills(filtered);
  }, [skillsCrud.items, categoryFilter, proficiencyFilter, featuredFilter]);

  useEffect(() => {
    filterSkills();
  }, [filterSkills]);

  const handleSubmit = skillsCrud.handleSubmit;

  const handleDelete = skillsCrud.handleDelete;

  const handleToggleFeatured = async (skillId: string) => {
    try {
      await skillsApi.toggleFeatured(skillId);
      await skillsCrud.fetchItems();
    } catch (err) {
      console.error('Error toggling featured status:', err);
      skillsCrud.setError('Failed to toggle featured status');
    }
  };

  const handleEdit = skillsCrud.handleEdit;
  const handleAddNew = skillsCrud.handleCreate;
  const resetForm = () => {
    skillsCrud.resetForm();
    setTagInput('');
  };

  const handleInputChange = (field: keyof SkillFormData, value: any) => {
    skillsCrud.setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const addTag = () => {
    if (tagInput.trim() && !skillsCrud.formData.tags.includes(tagInput.trim())) {
      skillsCrud.setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, tagInput.trim()]
      }));
      setTagInput('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    skillsCrud.setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }));
  };

  if (skillsCrud.loading) {
    return <LoadingSpinner>Loading skills...</LoadingSpinner>;
  }

  return (
    <Container>
      <Header>
        <Title>Skills Management</Title>
        <AddButton onClick={handleAddNew}>
          Add New Skill
        </AddButton>
      </Header>

      {skillsCrud.error && <ErrorMessage ref={errorRef}>{skillsCrud.error}</ErrorMessage>}

      <FilterSection>
        <FilterSelect 
          value={categoryFilter} 
          onChange={(e) => setCategoryFilter(e.target.value)}
        >
          <option value="">All Categories</option>
          {CATEGORIES.map(cat => (
            <option key={cat.value} value={cat.value}>{cat.label}</option>
          ))}
        </FilterSelect>

        <FilterSelect 
          value={proficiencyFilter} 
          onChange={(e) => setProficiencyFilter(e.target.value)}
        >
          <option value="">All Proficiency Levels</option>
          {PROFICIENCY_LEVELS.map(level => (
            <option key={level.value} value={level.value}>{level.label}</option>
          ))}
        </FilterSelect>

        <FilterSelect 
          value={featuredFilter} 
          onChange={(e) => setFeaturedFilter(e.target.value)}
        >
          <option value="">All Skills</option>
          <option value="featured">Featured Only</option>
          <option value="not-featured">Not Featured</option>
        </FilterSelect>
      </FilterSection>

      <SkillGrid>
        <AnimatePresence>
          {filteredSkills.map(skill => (
            <SkillCard
              key={skill.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.2 }}
            >
              <SkillHeader>
                <SkillName>
                  {skill.color && (
                    <SkillIcon color={skill.color}>
                      {skill.name.substring(0, 2).toUpperCase()}
                    </SkillIcon>
                  )}
                  {skill.name}
                </SkillName>
                <SkillCategory category={skill.category}>
                  {CATEGORIES.find(c => c.value === skill.category)?.label}
                </SkillCategory>
              </SkillHeader>

              <SkillDescription>{skill.description}</SkillDescription>

              <ProficiencySection>
                <ProficiencyHeader>
                  <ProficiencyLevel>{skill.proficiency_level}</ProficiencyLevel>
                  <ProficiencyValue>{skill.proficiency_value}%</ProficiencyValue>
                </ProficiencyHeader>
                <ProficiencyBar>
                  <ProficiencyFill width={skill.proficiency_value} />
                </ProficiencyBar>
              </ProficiencySection>

              <SkillMeta>
                <FeaturedBadge featured={skill.is_featured}>
                  {skill.is_featured ? 'Featured' : 'Not Featured'}
                </FeaturedBadge>
                <SortOrder>Order: {skill.sort_order}</SortOrder>
              </SkillMeta>

              {skill.tags && skill.tags.length > 0 && (
                <TagContainer>
                  {skill.tags.map(tag => (
                    <Tag key={tag}>{tag}</Tag>
                  ))}
                </TagContainer>
              )}

              <ActionButtons>
                <ActionButton variant="edit" onClick={() => handleEdit(skill)}>
                  Edit
                </ActionButton>
                <ActionButton 
                  variant="featured" 
                  onClick={() => handleToggleFeatured(skill.id)}
                >
                  {skill.is_featured ? 'Unfeature' : 'Feature'}
                </ActionButton>
                <ActionButton variant="delete" onClick={() => handleDelete(skill.id)}>
                  Delete
                </ActionButton>
              </ActionButtons>
            </SkillCard>
          ))}
        </AnimatePresence>
      </SkillGrid>

      <ScrollableModal 
        isOpen={skillsCrud.showModal} 
        onClose={() => skillsCrud.setShowModal(false)}
      >
        <ModalHeader>
          {skillsCrud.editingItem ? 'Edit Skill' : 'Add New Skill'}
        </ModalHeader>

            <Form onSubmit={handleSubmit}>
              <FormGroup>
                <Label>Name</Label>
                <Input
                  type="text"
                  value={skillsCrud.formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  required
                />
              </FormGroup>

              <FormGroup>
                <Label>Description</Label>
                <TextArea
                  value={skillsCrud.formData.description}
                  onChange={(e) => handleInputChange('description', e.target.value)}
                  placeholder="Brief description of the skill..."
                />
              </FormGroup>

              <FormGroup>
                <Label>Category</Label>
                <Select
                  value={skillsCrud.formData.category}
                  onChange={(e) => handleInputChange('category', e.target.value)}
                  required
                >
                  {CATEGORIES.map(cat => (
                    <option key={cat.value} value={cat.value}>{cat.label}</option>
                  ))}
                </Select>
              </FormGroup>

              <FormGroup>
                <Label>Proficiency Level</Label>
                <Select
                  value={skillsCrud.formData.proficiency_level}
                  onChange={(e) => handleInputChange('proficiency_level', e.target.value)}
                  required
                >
                  {PROFICIENCY_LEVELS.map(level => (
                    <option key={level.value} value={level.value}>{level.label}</option>
                  ))}
                </Select>
              </FormGroup>

              <FormGroup>
                <RangeGroup>
                  <Label>Proficiency Value: <RangeValue>{skillsCrud.formData.proficiency_value}%</RangeValue></Label>
                  <RangeInput
                    type="range"
                    min="0"
                    max="100"
                    value={skillsCrud.formData.proficiency_value}
                    onChange={(e) => handleInputChange('proficiency_value', parseInt(e.target.value))}
                  />
                </RangeGroup>
              </FormGroup>

              <FormGroup>
                <Label>Sort Order</Label>
                <Input
                  type="number"
                  value={skillsCrud.formData.sort_order}
                  onChange={(e) => handleInputChange('sort_order', parseInt(e.target.value))}
                  min="0"
                />
              </FormGroup>

              <FormGroup>
                <Label>Color (Hex)</Label>
                <Input
                  type="color"
                  value={skillsCrud.formData.color}
                  onChange={(e) => handleInputChange('color', e.target.value)}
                />
              </FormGroup>

              <FormGroup>
                <Label>Icon URL (optional)</Label>
                <Input
                  type="url"
                  value={skillsCrud.formData.icon_url}
                  onChange={(e) => handleInputChange('icon_url', e.target.value)}
                  placeholder="https://example.com/icon.png"
                />
              </FormGroup>

              <FormGroup>
                <CheckboxGroup>
                  <Checkbox
                    type="checkbox"
                    checked={skillsCrud.formData.is_featured}
                    onChange={(e) => handleInputChange('is_featured', e.target.checked)}
                  />
                  <Label>Featured Skill</Label>
                </CheckboxGroup>
              </FormGroup>

              <FormGroup>
                <Label>Tags</Label>
                <TagInputGroup>
                  <TagInput
                    type="text"
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    placeholder="Add a tag..."
                    onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
                  />
                  <AddTagButton type="button" onClick={addTag}>
                    Add Tag
                  </AddTagButton>
                </TagInputGroup>
                <TagList>
                  {skillsCrud.formData.tags.map(tag => (
                    <EditableTag key={tag}>
                      {tag}
                      <RemoveTagButton type="button" onClick={() => removeTag(tag)}>
                        ×
                      </RemoveTagButton>
                    </EditableTag>
                  ))}
                </TagList>
              </FormGroup>

              <ModalActions>
                <ModalButton type="button" variant="secondary" onClick={() => skillsCrud.setShowModal(false)}>
                  Cancel
                </ModalButton>
                <ModalButton type="submit" variant="primary">
                  {skillsCrud.editingItem ? 'Update Skill' : 'Create Skill'}
                </ModalButton>
              </ModalActions>
            </Form>
      </ScrollableModal>
    </Container>
  );
};

export default SkillsManager;