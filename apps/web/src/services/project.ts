import { apiClient } from './api';
import type { Project, ProjectMember, User, PaginatedResponse } from '../types';

// 项目相关类型定义
export interface CreateProjectRequest {
  name: string;
  code: string;
  description?: string;
  category: Project['category'];
  priority: Project['priority'];
  classification: Project['classification'];
  processTemplateId?: string;
  ownerId: string;
  startDate?: string;
  endDate?: string;
  teamMembers?: { userId: string; role: string }[];
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  category?: Project['category'];
  priority?: Project['priority'];
  classification?: Project['classification'];
  status?: Project['status'];
  ownerId?: string;
  startDate?: string;
  endDate?: string;
}

export interface ProjectFilterParams {
  page?: number;
  pageSize?: number;
  search?: string;
  status?: string;
  category?: string;
  ownerId?: string;
  priority?: string;
  startDateFrom?: string;
  startDateTo?: string;
}

export interface ProjectActivity {
  id: string;
  projectId: string;
  name: string;
  description?: string;
  type: 'task' | 'milestone' | 'review' | 'approval';
  status: 'pending' | 'ready' | 'running' | 'completed' | 'reviewing' | 'approved' | 'rejected' | 'blocked';
  progress: number;
  plannedStart?: string;
  plannedEnd?: string;
  actualStart?: string;
  actualEnd?: string;
  assigneeId?: string;
  assignee?: User;
  dependencies?: string[];
  deliverables?: ProjectDeliverable[];
  createdAt: string;
  updatedAt: string;
}

export interface ProjectDeliverable {
  id: string;
  activityId: string;
  name: string;
  description?: string;
  type: 'document' | 'design' | 'code' | 'test_report' | 'review_report' | 'other';
  status: 'pending' | 'in_progress' | 'submitted' | 'reviewing' | 'approved' | 'rejected';
  templateUrl?: string;
  submittedAt?: string;
  reviewedAt?: string;
  reviewedBy?: string;
  reviewComment?: string;
  createdAt: string;
}

export interface GanttTask {
  id: string;
  name: string;
  start: string;
  end: string;
  progress: number;
  type: 'task' | 'milestone';
  status: string;
  assignee?: User;
  dependencies: string[];
}

export interface ProcessTemplate {
  id: string;
  name: string;
  description?: string;
  category: string;
  version: string;
  activitiesCount: number;
  estimatedDuration: number;
  isActive: boolean;
}

export interface AddMemberRequest {
  userId: string;
  role: string;
}

// 项目API服务
export const projectAPI = {
  // 项目CRUD
  getProjects: (params?: ProjectFilterParams): Promise<PaginatedResponse<Project>> =>
    apiClient.getPaginated('/projects', params),

  getProject: (id: string): Promise<Project> =>
    apiClient.get(`/projects/${id}`),

  createProject: (data: CreateProjectRequest): Promise<Project> =>
    apiClient.post('/projects', data),

  updateProject: (id: string, data: UpdateProjectRequest): Promise<Project> =>
    apiClient.put(`/projects/${id}`, data),

  deleteProject: (id: string): Promise<void> =>
    apiClient.delete(`/projects/${id}`),

  // 项目状态管理
  updateProjectStatus: (id: string, status: Project['status']): Promise<Project> =>
    apiClient.put(`/projects/${id}/status`, { status }),

  startProject: (id: string): Promise<Project> =>
    apiClient.post(`/projects/${id}/start`),

  completeProject: (id: string): Promise<Project> =>
    apiClient.post(`/projects/${id}/complete`),

  archiveProject: (id: string): Promise<Project> =>
    apiClient.post(`/projects/${id}/archive`),

  // 项目成员管理
  getProjectMembers: (projectId: string): Promise<ProjectMember[]> =>
    apiClient.get(`/projects/${projectId}/members`),

  addProjectMember: (projectId: string, data: AddMemberRequest): Promise<ProjectMember> =>
    apiClient.post(`/projects/${projectId}/members`, data),

  updateMemberRole: (projectId: string, userId: string, role: string): Promise<ProjectMember> =>
    apiClient.put(`/projects/${projectId}/members/${userId}`, { role }),

  removeProjectMember: (projectId: string, userId: string): Promise<void> =>
    apiClient.delete(`/projects/${projectId}/members/${userId}`),

  // 项目活动
  getProjectActivities: (projectId: string): Promise<ProjectActivity[]> =>
    apiClient.get(`/projects/${projectId}/activities`),

  getActivity: (projectId: string, activityId: string): Promise<ProjectActivity> =>
    apiClient.get(`/projects/${projectId}/activities/${activityId}`),

  updateActivityProgress: (projectId: string, activityId: string, progress: number): Promise<ProjectActivity> =>
    apiClient.put(`/projects/${projectId}/activities/${activityId}/progress`, { progress }),

  updateActivityDates: (projectId: string, activityId: string, dates: { start?: string; end?: string }): Promise<ProjectActivity> =>
    apiClient.put(`/projects/${projectId}/activities/${activityId}/dates`, dates),

  // 甘特图数据
  getGanttData: (projectId: string): Promise<GanttTask[]> =>
    apiClient.get(`/projects/${projectId}/gantt`),

  updateGanttTask: (projectId: string, taskId: string, data: Partial<GanttTask>): Promise<GanttTask> =>
    apiClient.put(`/projects/${projectId}/gantt/${taskId}`, data),

  // 流程模板
  getProcessTemplates: (category?: string): Promise<ProcessTemplate[]> =>
    apiClient.get('/process-templates', { params: { category } }),

  getProcessTemplate: (id: string): Promise<ProcessTemplate> =>
    apiClient.get(`/process-templates/${id}`),

  // 项目文件
  getProjectFiles: (projectId: string): Promise<{ id: string; name: string; size: number; uploadedAt: string; uploader: User }[]> =>
    apiClient.get(`/projects/${projectId}/files`),

  // 项目统计
  getProjectStats: (projectId: string): Promise<{
    totalActivities: number;
    completedActivities: number;
    pendingActivities: number;
    progress: number;
    daysRemaining?: number;
    delayDays?: number;
  }> => apiClient.get(`/projects/${projectId}/stats`),

  // 批量操作
  batchUpdateProjects: (ids: string[], data: Partial<UpdateProjectRequest>): Promise<void> =>
    apiClient.post('/projects/batch-update', { ids, data }),

  batchDeleteProjects: (ids: string[]): Promise<void> =>
    apiClient.post('/projects/batch-delete', { ids }),
};

export default projectAPI;
