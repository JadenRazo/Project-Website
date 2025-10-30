import React, { useState, useEffect, useCallback, useRef } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import Badge from '../common/Badge';
import { useScrollTo } from '../../hooks/useScrollTo';
import { ScrollableModal } from '../common/ScrollableModal';
import { useInlineFormScroll } from '../../hooks/useInlineFormScroll';
import { useCrudOperations } from '../../hooks/useCrudOperations';
import {
  Certification,
  Category,
  CertificationFormData,
  CategoryFormData,
  defaultCertificationFormData,
  defaultCategoryFormData,
  mapCertificationToFormData,
  mapCategoryToFormData,
  categoryApi,
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

const CertificationGrid = styled.div`
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

const CertificationCard = styled(motion.div)`
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

const LinkSection = styled.div`
  margin: 0.75rem 0 0.5rem 0;
  padding: 0.75rem 0 0.5rem 0;
  border-top: 1px solid ${({ theme }) => theme.colors.border}20;
  
  @media (max-width: 768px) {
    margin: 0.5rem 0 0.25rem 0;
    padding: 0.5rem 0 0.25rem 0;
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

const CardLink = styled.a`
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.2s ease;
  
  &:hover {
    text-decoration: underline;
    color: ${({ theme }) => theme.colors.primaryHover || theme.colors.primary};
    transform: translateX(2px);
  }
  
  &::after {
    content: '↗';
    font-size: 0.8rem;
    opacity: 0.7;
  }
`;

// Badge component is now imported from common components

const ExpiryWarning = styled.span`
  color: ${({ theme }) => theme.colors.warning || '#F59E0B'};
  font-size: 0.85rem;
  font-weight: 500;
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

// Removed duplicate ErrorMessage - using the one defined above

const LoadingSpinner = styled.div`
  text-align: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

// Removed DeleteConfirmModal - using ScrollableModal instead

const DeleteConfirmText = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 1.5rem;
`;

// Component
const CertificationsManager: React.FC = () => {
  // Use CRUD operations hook for certifications
  const certificationsCrud = useCrudOperations<Certification, CertificationFormData>(
    {
      baseUrl: '/api/v1/devpanel/certifications',
      resourceName: 'certifications',
      requiresAuth: true,
    },
    defaultCertificationFormData,
    {
      onItemToFormData: mapCertificationToFormData,
      fetchQueryParams: 'include_hidden=true',
    }
  );

  // Use CRUD operations hook for categories
  const categoriesCrud = useCrudOperations<Category, CategoryFormData>(
    {
      baseUrl: '/api/v1/devpanel/certification-categories',
      resourceName: 'categories',
      requiresAuth: true,
    },
    defaultCategoryFormData,
    {
      onItemToFormData: mapCategoryToFormData,
      fetchQueryParams: 'include_hidden=true',
    }
  );

  // Local state for UI
  const [filteredCertifications, setFilteredCertifications] = useState<Certification[]>([]);
  const [showCategoryModal, setShowCategoryModal] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deletingCert, setDeletingCert] = useState<Certification | null>(null);
  
  // Filter states
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('all');
  const [visibilityFilter, setVisibilityFilter] = useState('all');
  
  // Scroll hooks and refs
  const { scrollToElement } = useScrollTo();
  const { formRef: errorRef } = useInlineFormScroll(!!(certificationsCrud.error || categoriesCrud.error), {
    scrollOffset: 100,
    scrollDelay: 200
  });

  // Fetch data on mount
  useEffect(() => {
    certificationsCrud.fetchItems();
    categoriesCrud.fetchItems();
  }, []);

  // Filter certifications
  useEffect(() => {
    let filtered = certificationsCrud.items;
    
    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(cert => 
        cert.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        cert.issuer.toLowerCase().includes(searchTerm.toLowerCase()) ||
        cert.description?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    
    // Category filter
    if (categoryFilter !== 'all') {
      filtered = filtered.filter(cert => cert.category_id === categoryFilter);
    }
    
    // Visibility filter
    if (visibilityFilter === 'visible') {
      filtered = filtered.filter(cert => cert.is_visible);
    } else if (visibilityFilter === 'hidden') {
      filtered = filtered.filter(cert => !cert.is_visible);
    }
    
    setFilteredCertifications(filtered);
  }, [certificationsCrud.items, searchTerm, categoryFilter, visibilityFilter]);

  // Handle certification form submit - now uses the CRUD hook
  const handleCertSubmit = certificationsCrud.handleSubmit;

  // Handle category form submit - now uses the CRUD hook
  const handleCategorySubmit = categoriesCrud.handleSubmit;

  // Handle delete
  const handleDelete = async () => {
    if (!deletingCert) return;
    
    try {
      await certificationsCrud.deleteItem(deletingCert.id);
      setShowDeleteConfirm(false);
      setDeletingCert(null);
    } catch (err) {
      // Error is already handled by the CRUD hook
      console.error('Delete failed:', err);
    }
  };

  // Edit certification - now uses the CRUD hook
  const handleEdit = (cert: Certification) => {
    certificationsCrud.handleEdit(cert);
  };

  // Reset forms - now use the CRUD hooks
  const resetCertForm = certificationsCrud.resetForm;
  const resetCategoryForm = categoriesCrud.resetForm;

  // Check if certification is expiring soon (within 90 days)
  const isExpiringSoon = (expiryDate: string | undefined) => {
    if (!expiryDate) return false;
    const expiry = new Date(expiryDate);
    const today = new Date();
    const daysUntilExpiry = Math.floor((expiry.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    return daysUntilExpiry <= 90 && daysUntilExpiry > 0;
  };

  // Check if certification is expired
  const isExpired = (expiryDate: string | undefined) => {
    if (!expiryDate) return false;
    const expiry = new Date(expiryDate);
    const today = new Date();
    return expiry < today;
  };

  if (certificationsCrud.loading) return <LoadingSpinner>Loading certifications...</LoadingSpinner>;

  const error = certificationsCrud.error || categoriesCrud.error;

  return (
    <Container>
      {error && (
        <ErrorMessage ref={errorRef}>
          {error}
        </ErrorMessage>
      )}
      
      <Header>
        <Title>Certifications Management</Title>
        <ButtonGroup>
          <AddButton onClick={() => {
            categoriesCrud.resetForm();
            setShowCategoryModal(true);
          }}>
            Manage Categories
          </AddButton>
          <AddButton onClick={certificationsCrud.handleCreate}>
            Add Certification
          </AddButton>
        </ButtonGroup>
      </Header>

      {error && <ErrorMessage ref={errorRef}>{error}</ErrorMessage>}

      <FilterSection>
        <SearchInput
          type="text"
          placeholder="Search certifications..."
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

      <CertificationGrid>
        <AnimatePresence>
          {filteredCertifications.map(cert => (
            <CertificationCard
              key={cert.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <CardHeader>
                <CardTitle>{cert.name}</CardTitle>
                <CardActions>
                  <IconButton onClick={() => handleEdit(cert)}>✏️</IconButton>
                  <IconButton onClick={() => {
                    setDeletingCert(cert);
                    setShowDeleteConfirm(true);
                  }}>🗑️</IconButton>
                </CardActions>
              </CardHeader>
              
              <CardBody>
                <InfoSection>
                  <CardInfo>
                    <strong>Issuer:</strong> {cert.issuer}
                  </CardInfo>
                  
                  {cert.credential_id && (
                    <CardInfo>
                      <strong>Credential ID:</strong> {cert.credential_id}
                    </CardInfo>
                  )}
                  
                  <CardInfo>
                    <strong>Issue Date:</strong> {new Date(cert.issue_date).toLocaleDateString()}
                  </CardInfo>
                  
                  {cert.expiry_date && (
                    <CardInfo>
                      <strong>Expiry Date:</strong> {new Date(cert.expiry_date).toLocaleDateString()}
                      {isExpired(cert.expiry_date) && (
                        <ExpiryWarning> (Expired)</ExpiryWarning>
                      )}
                      {isExpiringSoon(cert.expiry_date) && (
                        <ExpiryWarning> (Expiring Soon)</ExpiryWarning>
                      )}
                    </CardInfo>
                  )}
                  
                  {cert.description && (
                    <CardInfo>{cert.description}</CardInfo>
                  )}
                </InfoSection>
                
                {cert.verification_url && (
                  <LinkSection>
                    <CardLink href={cert.verification_url} target="_blank" rel="noopener noreferrer">
                      {cert.verification_text || 'View Credential'}
                    </CardLink>
                  </LinkSection>
                )}
                
                <BadgeContainer>
                  {cert.is_featured && <Badge variant="featured">Featured</Badge>}
                  {cert.category && <Badge variant="category">{cert.category.name}</Badge>}
                  {!cert.is_visible && <Badge variant="status">Hidden</Badge>}
                </BadgeContainer>
              </CardBody>
            </CertificationCard>
          ))}
        </AnimatePresence>
      </CertificationGrid>

      {/* Certification Modal */}
      <ScrollableModal 
        isOpen={certificationsCrud.showModal} 
        onClose={() => certificationsCrud.setShowModal(false)}
      >
        <ModalHeader>
          {certificationsCrud.editingItem ? 'Edit Certification' : 'Add New Certification'}
        </ModalHeader>
            
            <Form onSubmit={handleCertSubmit}>
              <FormRow>
                <FormGroup>
                  <Label>Certification Name *</Label>
                  <Input
                    type="text"
                    value={certificationsCrud.formData.name}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, name: e.target.value })}
                    required
                  />
                </FormGroup>
                
                <FormGroup>
                  <Label>Issuer *</Label>
                  <Input
                    type="text"
                    value={certificationsCrud.formData.issuer}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, issuer: e.target.value })}
                    required
                  />
                </FormGroup>
              </FormRow>
              
              <FormRow>
                <FormGroup>
                  <Label>Credential ID</Label>
                  <Input
                    type="text"
                    value={certificationsCrud.formData.credential_id}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, credential_id: e.target.value })}
                  />
                </FormGroup>
                
                <FormGroup>
                  <Label>Category</Label>
                  <Select
                    value={certificationsCrud.formData.category_id}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, category_id: e.target.value })}
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
              
              <FormRow>
                <FormGroup>
                  <Label>Issue Date *</Label>
                  <Input
                    type="date"
                    value={certificationsCrud.formData.issue_date}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, issue_date: e.target.value })}
                    required
                  />
                </FormGroup>
                
                <FormGroup>
                  <Label>Expiry Date</Label>
                  <Input
                    type="date"
                    value={certificationsCrud.formData.expiry_date}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, expiry_date: e.target.value })}
                  />
                </FormGroup>
              </FormRow>
              
              <FormGroup>
                <Label>Verification URL</Label>
                <Input
                  type="url"
                  value={certificationsCrud.formData.verification_url}
                  onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, verification_url: e.target.value })}
                  placeholder="https://example.com/verify/12345"
                />
              </FormGroup>
              
              <FormGroup>
                <Label>Verification Link Text</Label>
                <Input
                  type="text"
                  value={certificationsCrud.formData.verification_text}
                  onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, verification_text: e.target.value })}
                  placeholder="Verify on Credly"
                />
              </FormGroup>
              
              <FormGroup>
                <Label>Badge URL</Label>
                <Input
                  type="url"
                  value={certificationsCrud.formData.badge_url}
                  onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, badge_url: e.target.value })}
                  placeholder="https://example.com/badge.png"
                />
              </FormGroup>
              
              <FormGroup>
                <Label>Description</Label>
                <TextArea
                  value={certificationsCrud.formData.description}
                  onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, description: e.target.value })}
                  placeholder="Brief description of the certification..."
                />
              </FormGroup>
              
              <FormRow>
                <FormGroup>
                  <Label>Sort Order</Label>
                  <Input
                    type="number"
                    value={certificationsCrud.formData.sort_order}
                    onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, sort_order: parseInt(e.target.value) })}
                  />
                </FormGroup>
                
                <FormGroup style={{ justifyContent: 'center' }}>
                  <CheckboxGroup>
                    <Checkbox
                      type="checkbox"
                      id="featured"
                      checked={certificationsCrud.formData.is_featured}
                      onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, is_featured: e.target.checked })}
                    />
                    <Label htmlFor="featured">Featured</Label>
                  </CheckboxGroup>
                  
                  <CheckboxGroup>
                    <Checkbox
                      type="checkbox"
                      id="visible"
                      checked={certificationsCrud.formData.is_visible}
                      onChange={(e) => certificationsCrud.setFormData({ ...certificationsCrud.formData, is_visible: e.target.checked })}
                    />
                    <Label htmlFor="visible">Visible</Label>
                  </CheckboxGroup>
                </FormGroup>
              </FormRow>
              
              <ModalActions>
                <ModalButton type="button" variant="secondary" onClick={() => certificationsCrud.setShowModal(false)}>
                  Cancel
                </ModalButton>
                <ModalButton type="submit" variant="primary">
                  {certificationsCrud.editingItem ? 'Update' : 'Add'} Certification
                </ModalButton>
              </ModalActions>
            </Form>
      </ScrollableModal>

      {/* Category Modal */}
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
            
            {/* List existing categories */}
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

      {/* Delete Confirmation Modal */}
      <ScrollableModal 
        isOpen={showDeleteConfirm && !!deletingCert} 
        onClose={() => setShowDeleteConfirm(false)}
      >
        <ModalHeader>Confirm Deletion</ModalHeader>
            <DeleteConfirmText>
              Are you sure you want to delete the certification "{deletingCert?.name}"? This action cannot be undone.
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

export default CertificationsManager;