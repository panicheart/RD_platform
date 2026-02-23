import { api, ApiResponse } from '@/utils/api'

export interface Product {
  id: string
  name: string
  description: string
  trl_level: number
  category: string
  version: string
  source_project_id?: string
  owner_id?: string
  is_published: boolean
  published_at?: string
  download_count: number
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
  created_by?: string
  owner?: {
    id: string
    username: string
    display_name: string
  }
  source_project?: {
    id: string
    name: string
  }
  versions?: ProductVersion[]
}

export interface ProductVersion {
  id: string
  product_id: string
  version: string
  parent_version_id?: string
  changelog: string
  files?: Record<string, any>
  created_at: string
  created_by?: string
}

export interface Technology {
  id: string
  name: string
  description: string
  trl_level: number
  category: string
  parent_id?: string
  owner_id?: string
  is_published: boolean
  created_at: string
  updated_at: string
  created_by?: string
  owner?: {
    id: string
    username: string
    display_name: string
  }
  children?: Technology[]
}

export interface CartItem {
  id: string
  user_id: string
  product_id: string
  project_id?: string
  quantity: number
  notes: string
  created_at: string
  product?: Product
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface CreateProductRequest {
  name: string
  description?: string
  trl_level: number
  category?: string
  version?: string
  source_project_id?: string
  owner_id?: string
  metadata?: Record<string, any>
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  trl_level?: number
  category?: string
  version?: string
  owner_id?: string
  metadata?: Record<string, any>
}

export interface CreateProductVersionRequest {
  version: string
  parent_version_id?: string
  changelog?: string
  files?: Record<string, any>
}

export interface CreateTechnologyRequest {
  name: string
  description?: string
  trl_level: number
  category?: string
  parent_id?: string
  owner_id?: string
}

export interface UpdateTechnologyRequest {
  name?: string
  description?: string
  trl_level?: number
  category?: string
  owner_id?: string
  is_published?: boolean
}

export interface AddToCartRequest {
  product_id: string
  quantity?: number
  notes?: string
}

export interface UpdateCartItemRequest {
  quantity: number
  notes?: string
}

export interface TRLLevel {
  level: number
  name: string
  color: string
}

// Product APIs
export const productApi = {
  list: (params?: {
    page?: number
    page_size?: number
    category?: string
    trl_level?: number
    trl_min?: number
    trl_max?: number
    is_published?: boolean
    search?: string
  }) => api.get<ApiResponse<PaginatedResponse<Product>>>('/products', { params }),
  
  getById: (id: string) => api.get<ApiResponse<Product>>(`/products/${id}`),
  
  create: (data: CreateProductRequest) => api.post<ApiResponse<Product>>('/products', data),
  
  update: (id: string, data: UpdateProductRequest) => api.put<ApiResponse<Product>>(`/products/${id}`, data),
  
  delete: (id: string) => api.delete<ApiResponse<void>>(`/products/${id}`),
  
  publish: (id: string) => api.post<ApiResponse<Product>>(`/products/${id}/publish`),
  
  unpublish: (id: string) => api.post<ApiResponse<Product>>(`/products/${id}/unpublish`),
  
  getCategories: () => api.get<ApiResponse<string[]>>('/products/categories'),
  
  // Versions
  getVersions: (productId: string) => api.get<ApiResponse<ProductVersion[]>>(`/products/${productId}/versions`),
  
  createVersion: (productId: string, data: CreateProductVersionRequest) => 
    api.post<ApiResponse<ProductVersion>>(`/products/${productId}/versions`, data),
  
  deleteVersion: (versionId: string) => api.delete<ApiResponse<void>>(`/products/versions/${versionId}`),
}

// Technology APIs
export const technologyApi = {
  list: (params?: {
    page?: number
    page_size?: number
    category?: string
    trl_level?: number
    is_published?: boolean
    search?: string
  }) => api.get<ApiResponse<PaginatedResponse<Technology>>>('/technologies', { params }),
  
  getById: (id: string) => api.get<ApiResponse<Technology>>(`/technologies/${id}`),
  
  create: (data: CreateTechnologyRequest) => api.post<ApiResponse<Technology>>('/technologies', data),
  
  update: (id: string, data: UpdateTechnologyRequest) => api.put<ApiResponse<Technology>>(`/technologies/${id}`, data),
  
  delete: (id: string) => api.delete<ApiResponse<void>>(`/technologies/${id}`),
}

// Cart APIs
export const cartApi = {
  getItems: () => api.get<ApiResponse<CartItem[]>>('/cart'),
  
  add: (data: AddToCartRequest) => api.post<ApiResponse<CartItem>>('/cart', data),
  
  update: (itemId: string, data: UpdateCartItemRequest) => api.put<ApiResponse<CartItem>>(`/cart/${itemId}`, data),
  
  remove: (itemId: string) => api.delete<ApiResponse<void>>(`/cart/${itemId}`),
  
  clear: () => api.delete<ApiResponse<void>>('/cart'),
}

// TRL Level APIs
export const trlApi = {
  getLevels: () => api.get<ApiResponse<TRLLevel[]>>('/trl-levels'),
}

// Helper functions
export const getTRLColor = (level: number): string => {
  if (level <= 3) return 'red'
  if (level <= 6) return 'orange'
  return 'green'
}

export const getTRLName = (level: number): string => {
  const names: Record<number, string> = {
    1: '基本原理发现',
    2: '技术概念形成',
    3: '概念验证',
    4: '实验室验证',
    5: '相关环境验证',
    6: '系统/子系统验证',
    7: '系统原型验证',
    8: '系统完成验证',
    9: '实际应用验证',
  }
  return names[level] || `TRL ${level}`
}
