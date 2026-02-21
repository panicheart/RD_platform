// 用户相关类型定义

export type UserStatus = 'active' | 'inactive' | 'locked' | 'pending';
export type UserRole = 'admin' | 'manager' | 'leader' | 'designer' | 'viewer';
export type OrgType = 'department' | 'team' | 'group' | 'product_line';

export interface Role {
  id: string;
  name: string;
  code: string;
  description?: string;
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Organization {
  id: string;
  name: string;
  code: string;
  type: OrgType;
  parentId?: string;
  path: string;
  description?: string;
  memberCount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  displayName: string;
  avatar?: string;
  phone?: string;
  status: UserStatus;
  role: UserRole;
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

export interface UserFormData {
  username: string;
  email: string;
  password?: string;
  displayName: string;
  phone?: string;
  role: UserRole;
  organizationId?: string;
  title?: string;
  bio?: string;
}

export interface OrganizationFormData {
  name: string;
  code: string;
  type: OrgType;
  parentId?: string;
  description?: string;
}

export interface UserActivity {
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
