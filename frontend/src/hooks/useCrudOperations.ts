import { useState, useCallback } from 'react';

export interface CrudConfig {
  baseUrl: string;
  resourceName: string;
  requiresAuth?: boolean;
}

export interface CrudOperations<T, TFormData> {
  items: T[];
  loading: boolean;
  error: string | null;
  showModal: boolean;
  editingItem: T | null;
  formData: TFormData;
  
  // Actions
  setItems: React.Dispatch<React.SetStateAction<T[]>>;
  setLoading: React.Dispatch<React.SetStateAction<boolean>>;
  setError: React.Dispatch<React.SetStateAction<string | null>>;
  setShowModal: React.Dispatch<React.SetStateAction<boolean>>;
  setEditingItem: React.Dispatch<React.SetStateAction<T | null>>;
  setFormData: React.Dispatch<React.SetStateAction<TFormData>>;
  
  // CRUD Operations
  fetchItems: () => Promise<void>;
  createItem: (data: TFormData) => Promise<void>;
  updateItem: (id: string, data: TFormData) => Promise<void>;
  deleteItem: (id: string) => Promise<void>;
  
  // UI Helpers
  handleEdit: (item: T) => void;
  handleCreate: () => void;
  handleSubmit: (e: React.FormEvent) => Promise<void>;
  handleDelete: (id: string) => Promise<void>;
  resetForm: () => void;
  clearError: () => void;
}

export function useCrudOperations<T extends { id: string }, TFormData>(
  config: CrudConfig,
  initialFormData: TFormData,
  options: {
    onItemToFormData?: (item: T) => TFormData;
    onSuccess?: (action: 'create' | 'update' | 'delete') => void;
    customDelete?: boolean;
    fetchQueryParams?: string;
    extractItemsFromResponse?: (data: any) => T[];
  } = {}
): CrudOperations<T, TFormData> {
  const [items, setItems] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState<T | null>(null);
  const [formData, setFormData] = useState<TFormData>(initialFormData);

  const getHeaders = useCallback(() => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    
    if (config.requiresAuth) {
      const token = localStorage.getItem('auth_token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }
    
    return headers;
  }, [config.requiresAuth]);

  const buildUrl = useCallback((path = '') => {
    const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
    const baseUrl = apiUrl ? `${apiUrl}${config.baseUrl}` : config.baseUrl;
    return path ? `${baseUrl}/${path}` : baseUrl;
  }, [config.baseUrl]);

  const fetchItems = useCallback(async () => {
    try {
      setLoading(true);
      const url = buildUrl() + (options.fetchQueryParams ? `?${options.fetchQueryParams}` : '');
      const response = await fetch(url, {
        headers: getHeaders(),
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch ${config.resourceName}`);
      }

      const data = await response.json();
      const extractedItems = options.extractItemsFromResponse 
        ? options.extractItemsFromResponse(data)
        : (data[config.resourceName] || data.data || data);
      
      setItems(extractedItems);
      setError(null);
    } catch (err) {
      console.error(`Error fetching ${config.resourceName}:`, err);
      setError(err instanceof Error ? err.message : `Failed to load ${config.resourceName}`);
    } finally {
      setLoading(false);
    }
  }, [buildUrl, getHeaders, config.resourceName, options.fetchQueryParams, options.extractItemsFromResponse]);

  const createItem = useCallback(async (data: TFormData) => {
    const response = await fetch(buildUrl(), {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || errorData.message || `Failed to create ${config.resourceName}`);
    }

    await fetchItems();
    options.onSuccess?.('create');
  }, [buildUrl, getHeaders, fetchItems, config.resourceName, options.onSuccess]);

  const updateItem = useCallback(async (id: string, data: TFormData) => {
    const response = await fetch(buildUrl(id), {
      method: 'PUT',
      headers: getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || errorData.message || `Failed to update ${config.resourceName}`);
    }

    await fetchItems();
    options.onSuccess?.('update');
  }, [buildUrl, getHeaders, fetchItems, config.resourceName, options.onSuccess]);

  const deleteItem = useCallback(async (id: string) => {
    if (!options.customDelete && !window.confirm(`Are you sure you want to delete this ${config.resourceName}?`)) {
      return;
    }

    const response = await fetch(buildUrl(id), {
      method: 'DELETE',
      headers: getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to delete ${config.resourceName}`);
    }

    await fetchItems();
    options.onSuccess?.('delete');
  }, [buildUrl, getHeaders, fetchItems, config.resourceName, options.customDelete, options.onSuccess]);

  const handleEdit = useCallback((item: T) => {
    setEditingItem(item);
    const newFormData = options.onItemToFormData ? options.onItemToFormData(item) : item as unknown as TFormData;
    setFormData(newFormData);
    setShowModal(true);
  }, [options.onItemToFormData]);

  const handleCreate = useCallback(() => {
    setEditingItem(null);
    setFormData(initialFormData);
    setShowModal(true);
  }, [initialFormData]);

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      if (editingItem) {
        await updateItem(editingItem.id, formData);
      } else {
        await createItem(formData);
      }
      setShowModal(false);
      setFormData(initialFormData);
      setEditingItem(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to save ${config.resourceName}`);
    }
  }, [editingItem, formData, updateItem, createItem, initialFormData, config.resourceName]);

  const handleDelete = useCallback(async (id: string) => {
    try {
      await deleteItem(id);
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to delete ${config.resourceName}`);
    }
  }, [deleteItem, config.resourceName]);

  const resetForm = useCallback(() => {
    setEditingItem(null);
    setFormData(initialFormData);
    setError(null);
  }, [initialFormData]);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    items,
    loading,
    error,
    showModal,
    editingItem,
    formData,
    
    setItems,
    setLoading,
    setError,
    setShowModal,
    setEditingItem,
    setFormData,
    
    fetchItems,
    createItem,
    updateItem,
    deleteItem,
    
    handleEdit,
    handleCreate,
    handleSubmit,
    handleDelete,
    resetForm,
    clearError,
  };
}