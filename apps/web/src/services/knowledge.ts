import { api, ApiResponse } from '@/utils/api'

export interface Category {
  id: string
  name: string
  description: string
  level: number
  sort_order: number
  children: Category[]
}

export interface Tag {
  id: string
  name: string
  color: string
  count: number
}

export interface Knowledge {
  id: string
  title: string
  content: string
  category_id: string
  author_id: string
  tags: Tag[]
  status: 'draft' | 'published' | 'archived'
  version: number
  view_count: number
  source: string
  created_at: string
  updated_at: string
  published_at?: string
}

export interface KnowledgeReview {
  id: string
  knowledge_id: string
  reviewer_id: string
  status: 'pending' | 'approved' | 'rejected'
  comment: string
  created_at: string
  knowledge?: Knowledge
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface CreateCategoryRequest {
  name: string
  description?: string
  parent_id?: string
}

export interface UpdateCategoryRequest {
  name: string
  description?: string
}

export interface MoveCategoryRequest {
  parent_id?: string
}

export interface CreateKnowledgeRequest {
  title: string
  content: string
  category_id: string
  tag_ids?: string[]
  source?: string
}

export interface UpdateKnowledgeRequest {
  title: string
  content: string
  category_id?: string
  tag_ids?: string[]
}

export interface CreateTagRequest {
  name: string
  color?: string
}

export interface SubmitReviewRequest {
  reviewer_id: string
}

export interface ReviewActionRequest {
  comment?: string
}

// Category APIs
export const categoryApi = {
  getTree: () => api.get<ApiResponse<Category[]>>('/categories'),
  create: (data: CreateCategoryRequest) => api.post<ApiResponse<Category>>('/categories', data),
  update: (id: string, data: UpdateCategoryRequest) => api.put<ApiResponse<Category>>(`/categories/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/categories/${id}`),
  move: (id: string, data: MoveCategoryRequest) => api.post<ApiResponse<void>>(`/categories/${id}/move`, data),
}

// Knowledge APIs
export const knowledgeApi = {
  list: (params?: {
    page?: number
    page_size?: number
    category_id?: string
    status?: string
    tag_id?: string
    search?: string
  }) => api.get<ApiResponse<PaginatedResponse<Knowledge>>>('/knowledge', { params }),
  getById: (id: string) => api.get<ApiResponse<Knowledge>>(`/knowledge/${id}`),
  create: (data: CreateKnowledgeRequest) => api.post<ApiResponse<Knowledge>>('/knowledge', data),
  update: (id: string, data: UpdateKnowledgeRequest) => api.put<ApiResponse<Knowledge>>(`/knowledge/${id}`, data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/knowledge/${id}`),
  publish: (id: string) => api.post<ApiResponse<void>>(`/knowledge/${id}/publish`),
  archive: (id: string) => api.post<ApiResponse<void>>(`/knowledge/${id}/archive`),
  submitReview: (id: string, data: SubmitReviewRequest) => api.post<ApiResponse<KnowledgeReview>>(`/knowledge/${id}/review`, data),
}

// Tag APIs
export const tagApi = {
  list: () => api.get<ApiResponse<Tag[]>>('/tags'),
  create: (data: CreateTagRequest) => api.post<ApiResponse<Tag>>('/tags', data),
  delete: (id: string) => api.delete<ApiResponse<void>>(`/tags/${id}`),
}

// Review APIs
export const reviewApi = {
  getPending: () => api.get<ApiResponse<KnowledgeReview[]>>('/reviews/pending'),
  approve: (reviewId: string, data: ReviewActionRequest) => api.post<ApiResponse<void>>(`/reviews/${reviewId}/approve`, data),
  reject: (reviewId: string, data: ReviewActionRequest) => api.post<ApiResponse<void>>(`/reviews/${reviewId}/reject`, data),
}
