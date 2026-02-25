// DevPanel API Types
export interface Certification {
  id: string;
  name: string;
  issuer: string;
  credential_id?: string;
  issue_date: string;
  expiry_date?: string;
  verification_url?: string;
  verification_text?: string;
  badge_url?: string;
  description?: string;
  category_id?: string;
  category?: Category;
  is_featured: boolean;
  is_visible: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  name: string;
  description: string;
  sort_order: number;
  is_visible: boolean;
  created_at: string;
  updated_at: string;
}

export interface CertificationFormData {
  name: string;
  issuer: string;
  credential_id: string;
  issue_date: string;
  expiry_date: string;
  verification_url: string;
  verification_text: string;
  badge_url: string;
  description: string;
  category_id: string;
  is_featured: boolean;
  is_visible: boolean;
  sort_order: number;
}

export interface CategoryFormData {
  name: string;
  description: string;
  sort_order: number;
  is_visible: boolean;
}

export interface Skill {
  id: string;
  name: string;
  description: string;
  category: string;
  proficiency_level: string;
  proficiency_value: number;
  is_featured: boolean;
  sort_order: number;
  icon_url?: string;
  color?: string;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface SkillFormData {
  name: string;
  description: string;
  category: string;
  proficiency_level: string;
  proficiency_value: number;
  is_featured: boolean;
  sort_order: number;
  icon_url: string;
  color: string;
  tags: string[];
}

export interface Project {
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

export interface ProjectFormData {
  name: string;
  description: string;
  repo_url: string;
  live_url: string;
  tags: string[];
  status: 'active' | 'archived' | 'draft';
}

export interface ProjectPath {
  id: string;
  name: string;
  path: string;
  description?: string;
  exclude_patterns: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProjectPathFormData {
  name: string;
  path: string;
  description: string;
  exclude_patterns: string[];
  is_active: boolean;
}

export interface BlogPost {
  id: string;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  featured_image: string;
  status: 'draft' | 'published' | 'archived';
  published_at: string;
  tags: string[];
  view_count: number;
  read_time_minutes: number;
  is_featured: boolean;
  is_visible: boolean;
  created_at: string;
  updated_at: string;
}

export interface BlogPostFormData {
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  featured_image: string;
  status: string;
  tags_input: string;
  is_featured: boolean;
  is_visible: boolean;
}

// API Helper Functions
export const getAuthHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  const token = localStorage.getItem('auth_token');
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return headers;
};

export const buildApiUrl = (endpoint: string): string => {
  const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
  return apiUrl ? `${apiUrl}${endpoint}` : endpoint;
};

// Specific API Functions for each resource
export const certificationApi = {
  getAll: async (): Promise<Certification[]> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/certifications?include_hidden=true'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch certifications');
    return response.json();
  },

  create: async (data: CertificationFormData): Promise<Certification> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/certifications'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create certification');
    return response.json();
  },

  update: async (id: string, data: CertificationFormData): Promise<Certification> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/certifications/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update certification');
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/certifications/${id}`), {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete certification');
  },
};

export const categoryApi = {
  getAll: async (): Promise<Category[]> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/certification-categories?include_hidden=true'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch categories');
    return response.json();
  },

  create: async (data: CategoryFormData): Promise<Category> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/certification-categories'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create category');
    return response.json();
  },

  update: async (id: string, data: CategoryFormData): Promise<Category> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/certification-categories/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update category');
    return response.json();
  },
};

export const skillsApi = {
  getAll: async (): Promise<{ skills: Skill[] }> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/skills'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch skills');
    return response.json();
  },

  create: async (data: SkillFormData): Promise<Skill> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/skills'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create skill');
    return response.json();
  },

  update: async (id: string, data: SkillFormData): Promise<Skill> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/skills/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update skill');
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/skills/${id}`), {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete skill');
  },

  toggleFeatured: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/skills/${id}/toggle-featured`), {
      method: 'PUT',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to toggle featured status');
  },
};

export const projectsApi = {
  getAll: async (): Promise<{ projects: Project[] }> => {
    const response = await fetch(buildApiUrl('/api/v1/projects'));
    if (!response.ok) throw new Error('Failed to fetch projects');
    return response.json();
  },

  create: async (data: ProjectFormData): Promise<Project> => {
    const response = await fetch(buildApiUrl('/api/v1/projects'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create project');
    return response.json();
  },

  update: async (id: string, data: ProjectFormData): Promise<Project> => {
    const response = await fetch(buildApiUrl(`/api/v1/projects/${id}`), {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update project');
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/projects/${id}`), {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete project');
  },
};

export const projectPathsApi = {
  getAll: async (): Promise<{ data: ProjectPath[] }> => {
    const response = await fetch(buildApiUrl('/api/v1/code/paths'));
    if (!response.ok) throw new Error('Failed to fetch project paths');
    return response.json();
  },

  create: async (data: ProjectPathFormData): Promise<ProjectPath> => {
    const response = await fetch(buildApiUrl('/api/v1/code/paths'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Failed to create project path');
    }
    return response.json();
  },

  update: async (id: string, data: ProjectPathFormData): Promise<ProjectPath> => {
    const response = await fetch(buildApiUrl(`/api/v1/code/paths/${id}`), {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Failed to update project path');
    }
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/code/paths/${id}`), {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete project path');
  },

  refreshStats: async (action?: string): Promise<void> => {
    const response = await fetch(buildApiUrl('/api/v1/code/stats/update'), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      console.warn('Code stats update failed:', errorData.error || 'Unknown error');
    }
  },
};

// Form data mappers for each resource
export const mapCertificationToFormData = (cert: Certification): CertificationFormData => ({
  name: cert.name,
  issuer: cert.issuer,
  credential_id: cert.credential_id || '',
  issue_date: cert.issue_date.split('T')[0],
  expiry_date: cert.expiry_date ? cert.expiry_date.split('T')[0] : '',
  verification_url: cert.verification_url || '',
  verification_text: cert.verification_text || '',
  badge_url: cert.badge_url || '',
  description: cert.description || '',
  category_id: cert.category_id || '',
  is_featured: cert.is_featured,
  is_visible: cert.is_visible,
  sort_order: cert.sort_order,
});

export const mapCategoryToFormData = (category: Category): CategoryFormData => ({
  name: category.name,
  description: category.description,
  sort_order: category.sort_order,
  is_visible: category.is_visible,
});

export const mapSkillToFormData = (skill: Skill): SkillFormData => ({
  name: skill.name,
  description: skill.description,
  category: skill.category,
  proficiency_level: skill.proficiency_level,
  proficiency_value: skill.proficiency_value,
  is_featured: skill.is_featured,
  sort_order: skill.sort_order,
  icon_url: skill.icon_url || '',
  color: skill.color || '#6B7280',
  tags: skill.tags || [],
});

export const mapProjectToFormData = (project: Project): ProjectFormData => ({
  name: project.name,
  description: project.description,
  repo_url: project.repo_url,
  live_url: project.live_url || '',
  tags: project.tags || [],
  status: project.status,
});

export const mapProjectPathToFormData = (path: ProjectPath): ProjectPathFormData => ({
  name: path.name,
  path: path.path,
  description: path.description || '',
  exclude_patterns: path.exclude_patterns || [],
  is_active: path.is_active,
});

// Default form data for each resource
export const defaultCertificationFormData: CertificationFormData = {
  name: '',
  issuer: '',
  credential_id: '',
  issue_date: '',
  expiry_date: '',
  verification_url: '',
  verification_text: '',
  badge_url: '',
  description: '',
  category_id: '',
  is_featured: false,
  is_visible: true,
  sort_order: 1000,
};

export const defaultCategoryFormData: CategoryFormData = {
  name: '',
  description: '',
  sort_order: 1000,
  is_visible: true,
};

export const defaultSkillFormData: SkillFormData = {
  name: '',
  description: '',
  category: 'frontend',
  proficiency_level: 'intermediate',
  proficiency_value: 50,
  is_featured: false,
  sort_order: 1000,
  icon_url: '',
  color: '#6B7280',
  tags: [],
};

export const defaultProjectFormData: ProjectFormData = {
  name: '',
  description: '',
  repo_url: '',
  live_url: '',
  tags: [],
  status: 'draft',
};

export const defaultProjectPathFormData: ProjectPathFormData = {
  name: '',
  path: '',
  description: '',
  exclude_patterns: [],
  is_active: true,
};

export const blogApi = {
  getAll: async (): Promise<{ posts: BlogPost[]; total: number }> => {
    const response = await fetch(buildApiUrl('/api/v1/blog/admin'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch blog posts');
    return response.json();
  },

  create: async (data: Record<string, unknown>): Promise<BlogPost> => {
    const response = await fetch(buildApiUrl('/api/v1/blog/admin'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create blog post');
    return response.json();
  },

  update: async (id: string, data: Record<string, unknown>): Promise<BlogPost> => {
    const response = await fetch(buildApiUrl(`/api/v1/blog/admin/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update blog post');
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/blog/admin/${id}`), {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete blog post');
  },
};

export const mapBlogPostToFormData = (post: BlogPost): BlogPostFormData => ({
  title: post.title,
  slug: post.slug,
  content: post.content || '',
  excerpt: post.excerpt || '',
  featured_image: post.featured_image || '',
  status: post.status,
  tags_input: (post.tags || []).join(', '),
  is_featured: post.is_featured,
  is_visible: post.is_visible,
});

export const defaultBlogPostFormData: BlogPostFormData = {
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  featured_image: '',
  status: 'draft',
  tags_input: '',
  is_featured: false,
  is_visible: true,
};