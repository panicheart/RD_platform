// 统一导出所有类型

export * from './user';

export interface Role {
  id: string;
  name: string;
  code: string;
  description?: string;
  permissions: string[];
}

export interface Organization {
  id: string;
  name: string;
  code: string;
  type: 'department' | 'team' | 'group' | 'product_line';
  parentId?: string;
  path: string;
}

// 兼容hooks/useAuth.tsx中的User类型
export interface User {
  id: string;
  username: string;
  email: string;
  // 支持两种命名风格（后端返回vs前端使用）
  displayName?: string;
  display_name?: string;
  avatar?: string;
  avatar_url?: string;
  phone?: string;
  status: 'active' | 'inactive' | 'locked' | 'pending';
  role: string;
  roles: Role[];
  organization?: Organization;
  organizationId?: string;
  title?: string;
  bio?: string;
  skills?: string[];
  isActive: boolean;
  lastLoginAt?: string;
  lastLoginIp?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Project {
  id: string;
  name: string;
  code: string;
  description?: string;
  category: ProjectCategory;
  status: ProjectStatus;
  priority: 'low' | 'medium' | 'high' | 'critical';
  classification: 'public' | 'internal' | 'confidential' | 'secret';
  ownerId: string;
  owner?: User;
  teamMembers: ProjectMember[];
  startDate?: string;
  endDate?: string;
  processTemplateId?: string;
  createdAt: string;
  updatedAt: string;
}

export type ProjectCategory =
  | 'new_product'
  | 'product_improvement'
  | 'pre_research'
  | 'tech_platform'
  | 'component_development'
  | 'process_improvement'
  | 'other';

export type ProjectStatus =
  | 'draft'
  | 'planning'
  | 'in_progress'
  | 'on_hold'
  | 'completed'
  | 'cancelled'
  | 'archived';

export interface ProjectMember {
  userId: string;
  user?: User;
  role: string;
  joinedAt: string;
}

export interface FileItem {
  id: string;
  name: string;
  path: string;
  size: number;
  mimeType: string;
  classification: 'public' | 'internal' | 'confidential' | 'secret';
  projectId?: string;
  uploaderId: string;
  status: 'pending' | 'active' | 'archived' | 'deleted';
  createdAt: string;
  updatedAt: string;
}

export interface AuditLog {
  id: string;
  userId: string;
  user?: User;
  action: 'create' | 'read' | 'update' | 'delete' | 'login' | 'logout' | 'export' | 'import';
  resource: string;
  resourceId?: string;
  details?: Record<string, unknown>;
  ipAddress: string;
  userAgent: string;
  createdAt: string;
}

export interface Notification {
  id: string;
  userId: string;
  type: 'system' | 'project' | 'workflow' | 'mention' | 'deadline';
  title: string;
  content: string;
  isRead: boolean;
  link?: string;
  createdAt: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  items?: T[];  // 兼容不同API返回格式
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface APIResponse<T = unknown> {
  code: number;
  message: string;
  data?: T;
}

export interface APIError {
  code: number;
  message: string;
}
