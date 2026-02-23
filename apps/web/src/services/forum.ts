import { api, ApiResponse, PaginatedResponse } from '@/utils/api'
import type {
  ForumBoard,
  ForumPost,
  ForumReply,
  ForumTag,
  CreateBoardRequest,
  UpdateBoardRequest,
  CreatePostRequest,
  UpdatePostRequest,
  CreateReplyRequest,
  UpdateReplyRequest,
  ListBoardsQuery,
  ListPostsQuery,
  ListRepliesQuery,
} from '@/types/forum'

// Board APIs
export const boardApi = {
  list: (params?: ListBoardsQuery) =>
    api.get<ApiResponse<PaginatedResponse<ForumBoard>>>('/boards', { params }),
  getById: (id: string) =>
    api.get<ApiResponse<ForumBoard>>(`/boards/${id}`),
  create: (data: CreateBoardRequest) =>
    api.post<ApiResponse<ForumBoard>>('/boards', data),
  update: (id: string, data: UpdateBoardRequest) =>
    api.put<ApiResponse<ForumBoard>>(`/boards/${id}`, data),
  delete: (id: string) =>
    api.delete<ApiResponse<void>>(`/boards/${id}`),
}

// Post APIs
export const postApi = {
  list: (params?: ListPostsQuery) =>
    api.get<ApiResponse<PaginatedResponse<ForumPost>>>('/posts', { params }),
  listByBoard: (boardId: string, params?: Omit<ListPostsQuery, 'board_id'>) =>
    api.get<ApiResponse<PaginatedResponse<ForumPost>>>('/posts', { params: { ...params, board_id: boardId } }),
  getById: (id: string) =>
    api.get<ApiResponse<ForumPost>>(`/posts/${id}`),
  create: (data: CreatePostRequest) =>
    api.post<ApiResponse<ForumPost>>('/posts', data),
  update: (id: string, data: UpdatePostRequest) =>
    api.put<ApiResponse<ForumPost>>(`/posts/${id}`, data),
  delete: (id: string) =>
    api.delete<ApiResponse<void>>(`/posts/${id}`),
  togglePin: (id: string) =>
    api.post<ApiResponse<void>>(`/posts/${id}/pin`),
  toggleLock: (id: string) =>
    api.post<ApiResponse<void>>(`/posts/${id}/lock`),
}

// Reply APIs
export const replyApi = {
  listByPost: (postId: string, params?: ListRepliesQuery) =>
    api.get<ApiResponse<PaginatedResponse<ForumReply>>>(`/posts/${postId}/replies`, { params }),
  create: (postId: string, data: CreateReplyRequest) =>
    api.post<ApiResponse<ForumReply>>(`/posts/${postId}/replies`, data),
  update: (id: string, data: UpdateReplyRequest) =>
    api.put<ApiResponse<ForumReply>>(`/replies/${id}`, data),
  delete: (id: string) =>
    api.delete<ApiResponse<void>>(`/replies/${id}`),
}

// Tag APIs
export const forumTagApi = {
  list: () =>
    api.get<ApiResponse<ForumTag[]>>('/forum-tags'),
  create: (data: { name: string; color?: string }) =>
    api.post<ApiResponse<ForumTag>>('/forum-tags', data),
  delete: (id: string) =>
    api.delete<ApiResponse<void>>(`/forum-tags/${id}`),
}
