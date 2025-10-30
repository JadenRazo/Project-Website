// DevPanel API Types
export interface Prompt {
  id: string;
  name: string;
  description: string;
  prompt: string;
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

export interface PromptFormData {
  name: string;
  description: string;
  prompt: string;
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
export const promptApi = {
  getAll: async (): Promise<Prompt[]> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/prompts?include_hidden=true'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch prompts');
    return response.json();
  },

  create: async (data: PromptFormData): Promise<Prompt> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/prompts'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create prompt');
    return response.json();
  },

  update: async (id: string, data: PromptFormData): Promise<Prompt> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/prompts/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update prompt');
    return response.json();
  },

  delete: async (id: string): Promise<void> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/prompts/${id}`), {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete prompt');
  },
};

export const categoryApi = {
  getAll: async (): Promise<Category[]> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/prompt-categories?include_hidden=true'), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch categories');
    return response.json();
  },

  create: async (data: CategoryFormData): Promise<Category> => {
    const response = await fetch(buildApiUrl('/api/v1/devpanel/prompt-categories'), {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to create category');
    return response.json();
  },

  update: async (id: string, data: CategoryFormData): Promise<Category> => {
    const response = await fetch(buildApiUrl(`/api/v1/devpanel/prompt-categories/${id}`), {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update category');
    return response.json();
  },
};

// Form data mappers for each resource
export const mapPromptToFormData = (prompt: Prompt): PromptFormData => ({
  name: prompt.name,
  description: prompt.description || '',
  prompt: prompt.prompt || '',
  category_id: prompt.category_id || '',
  is_featured: prompt.is_featured,
  is_visible: prompt.is_visible,
  sort_order: prompt.sort_order,
});

export const mapCategoryToFormData = (category: Category): CategoryFormData => ({
  name: category.name,
  description: category.description,
  sort_order: category.sort_order,
  is_visible: category.is_visible,
});

// Default form data for each resource
export const defaultPromptFormData: PromptFormData = {
  name: '',
  description: '',
  prompt: '',
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
