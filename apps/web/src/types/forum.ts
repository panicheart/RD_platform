/**
 * Forum type definitions
 * RDP - Forum Module Types
 */

export interface ForumBoard {
  id: string;
  name: string;
  description: string;
  category: string;
  icon: string;
  sort_order: number;
  topic_count: number;
  post_count: number;
  last_post_at: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ForumPost {
  id: string;
  board_id: string;
  title: string;
  content: string;
  author_id: string;
  author_name: string;
  view_count: number;
  reply_count: number;
  is_pinned: boolean;
  is_locked: boolean;
  is_best_answer: boolean;
  tags: string[];
  knowledge_id?: string;
  last_reply_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface ForumReply {
  id: string;
  post_id: string;
  parent_id: string | null;
  content: string;
  author_id: string;
  author_name: string;
  is_best_answer: boolean;
  mentions: string[];
  created_at: string;
  updated_at: string;
}

export interface ForumTag {
  id: string;
  name: string;
  color: string;
  count: number;
  created_at: string;
}

// Request types
export interface CreateBoardRequest {
  name: string;
  description: string;
  category: string;
  icon?: string;
  sort_order?: number;
}

export interface UpdateBoardRequest {
  name?: string;
  description?: string;
  category?: string;
  icon?: string;
  sort_order?: number;
  is_active?: boolean;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  board_id: string;
  tags?: string[];
  knowledge_id?: string | null;
}

export interface UpdatePostRequest {
  title?: string;
  content?: string;
  tags?: string[];
}

export interface CreateReplyRequest {
  content: string;
  parent_id?: string | null;
}

export interface UpdateReplyRequest {
  content: string;
}

// Query types
export interface ListBoardsQuery {
  category?: string;
  page?: number;
  page_size?: number;
}

export interface ListPostsQuery {
  board_id?: string;
  author_id?: string;
  search?: string;
  is_pinned?: boolean;
  page?: number;
  page_size?: number;
}

export interface ListRepliesQuery {
  parent_id?: string;
  page?: number;
  page_size?: number;
}

// Paginated response
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

// API Response wrapper
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// Board categories
export const BOARD_CATEGORIES = [
  { value: 'tech', label: '技术讨论', icon: 'code' },
  { value: 'general', label: '综合讨论', icon: 'message' },
  { value: 'help', label: '求助问答', icon: 'question-circle' },
  { value: 'announcement', label: '公告通知', icon: 'notification' },
] as const;
