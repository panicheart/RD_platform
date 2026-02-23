import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import type { APIResponse, APIError, PaginatedResponse } from '../types';

const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

class APIClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError<APIError>) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get<APIResponse<T>>(url, config);
    return response.data.data as T;
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post<APIResponse<T>>(url, data, config);
    return response.data.data as T;
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put<APIResponse<T>>(url, data, config);
    return response.data.data as T;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete<APIResponse<T>>(url, config);
    return response.data.data as T;
  }

  async patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.patch<APIResponse<T>>(url, data, config);
    return response.data.data as T;
  }

  async getPaginated<T>(url: string, params?: Record<string, unknown>): Promise<PaginatedResponse<T>> {
    const response = await this.client.get<APIResponse<PaginatedResponse<T>>>(url, { params });
    return response.data.data as PaginatedResponse<T>;
  }
}

export const apiClient = new APIClient();

export const userAPI = {
  getCurrentUser: () => apiClient.get('/users/me'),
  getUsers: (params?: Record<string, unknown>) => apiClient.getPaginated('/users', params),
  getUser: (id: string) => apiClient.get(`/users/${id}`),
  updateUser: (id: string, data: unknown) => apiClient.put(`/users/${id}`, data),
};

export const projectAPI = {
  getProjects: (params?: Record<string, unknown>) => apiClient.getPaginated('/projects', params),
  getProject: (id: string) => apiClient.get(`/projects/${id}`),
  createProject: (data: unknown) => apiClient.post('/projects', data),
  updateProject: (id: string, data: unknown) => apiClient.put(`/projects/${id}`, data),
  deleteProject: (id: string) => apiClient.delete(`/projects/${id}`),
};

export const fileAPI = {
  upload: async (projectId: string, file: File): Promise<unknown> => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('projectId', projectId);
    
    const response = await apiClient.post('/files/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response;
  },
  download: (id: string) => `${BASE_URL}/files/${id}/download`,
  delete: (id: string) => apiClient.delete(`/files/${id}`),
};

export const auditAPI = {
  getLogs: (params?: Record<string, unknown>) => apiClient.getPaginated('/audit-logs', params),
};
