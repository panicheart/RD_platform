export interface User {
  id: string;
  username: string;
  email: string;
  displayName: string;
  avatar?: string;
  status: 'active' | 'inactive' | 'locked' | 'pending';
  roles: Role[];
  organization?: Organization;
  createdAt: string;
  updatedAt: string;
}

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

export interface File {
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
