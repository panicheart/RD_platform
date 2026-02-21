import { apiClient } from './api';
import type { User, PaginatedResponse } from '../types';

// 登录/认证相关
interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  display_name: string;
  role: string;
  organization_id?: string;
}

interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// 组织架构相关
interface OrganizationNode {
  id: string;
  name: string;
  code: string;
  type: 'department' | 'team' | 'group' | 'product_line';
  parentId?: string;
  description?: string;
  memberCount?: number;
}

interface CreateOrganizationRequest {
  name: string;
  code: string;
  type: 'department' | 'team' | 'group' | 'product_line';
  parentId?: string;
  description?: string;
}

interface UpdateOrganizationRequest {
  name?: string;
  type?: 'department' | 'team' | 'group' | 'product_line';
  parentId?: string;
  description?: string;
}

// 用户 API
export const userAPI = {
  // 认证相关
  login: (data: LoginRequest): Promise<LoginResponse> =>
    apiClient.post('/auth/login', data),

  logout: (): Promise<void> =>
    apiClient.post('/auth/logout'),

  refreshToken: (refreshToken: string): Promise<{ access_token: string }> =>
    apiClient.post('/auth/refresh', { refresh_token: refreshToken }),

  // 当前用户
  getCurrentUser: (): Promise<User> =>
    apiClient.get('/users/me'),

  updateCurrentUser: (data: Partial<User>): Promise<User> =>
    apiClient.put('/users/me', data),

  changePassword: (data: ChangePasswordRequest): Promise<void> =>
    apiClient.post('/users/me/password', data),

  uploadAvatar: (file: File): Promise<{ url: string }> => {
    const formData = new FormData();
    formData.append('file', file);
    return apiClient.post('/users/me/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },

  // 用户管理（需要管理员权限）
  getUsers: (params?: Record<string, unknown>): Promise<PaginatedResponse<User>> =>
    apiClient.getPaginated('/users', params),

  getUser: (id: string): Promise<User> =>
    apiClient.get(`/users/${id}`),

  createUser: (data: RegisterRequest): Promise<User> =>
    apiClient.post('/users', data),

  updateUser: (id: string, data: Partial<User>): Promise<User> =>
    apiClient.put(`/users/${id}`, data),

  deleteUser: (id: string): Promise<void> =>
    apiClient.delete(`/users/${id}`),

  updateUserStatus: (id: string, status: string): Promise<User> =>
    apiClient.put(`/users/${id}/status`, { status }),

  updateUserRole: (id: string, roleId: string): Promise<User> =>
    apiClient.put(`/users/${id}/role`, { role_id: roleId }),

  // 组织架构
  getOrganizationTree: (): Promise<OrganizationNode[]> =>
    apiClient.get('/organizations/tree'),

  getOrganization: (id: string): Promise<OrganizationNode> =>
    apiClient.get(`/organizations/${id}`),

  getOrganizationMembers: (orgId: string): Promise<User[]> =>
    apiClient.get(`/organizations/${orgId}/members`),

  createOrganization: (data: CreateOrganizationRequest): Promise<OrganizationNode> =>
    apiClient.post('/organizations', data),

  updateOrganization: (id: string, data: UpdateOrganizationRequest): Promise<OrganizationNode> =>
    apiClient.put(`/organizations/${id}`, data),

  deleteOrganization: (id: string): Promise<void> =>
    apiClient.delete(`/organizations/${id}`),

  moveUserToOrg: (userId: string, orgId: string): Promise<User> =>
    apiClient.put(`/users/${userId}/organization`, { organization_id: orgId }),

  // 技能标签
  updateUserSkills: (userId: string, skills: string[]): Promise<User> =>
    apiClient.put(`/users/${userId}/skills`, { skills }),

  // 搜索用户
  searchUsers: (query: string): Promise<User[]> =>
    apiClient.get('/users/search', { params: { q: query } }),

  // 获取用户活动日志
  getUserActivities: (userId: string, params?: Record<string, unknown>): Promise<PaginatedResponse<{
    id: string;
    action: string;
    resource: string;
    createdAt: string;
    details?: string;
  }>> => apiClient.getPaginated(`/users/${userId}/activities`, params),
};

export default userAPI;
