const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

import { apiClient } from './api';
import type {
  SystemMetric,
  LogEntry,
  AlertRule,
  AlertHistory,
  SystemMetricsSummary,
  APIMetricsSummary,
  LogFilterParams,
  AlertFilterParams,
  CreateAlertRuleRequest,
  UpdateAlertRuleRequest,
  MonitorDashboardData,
  MonitorTimeRange,
} from '../types/monitor';

/**
 * Monitor API service for system monitoring dashboard
 */
export const monitorAPI = {
  /**
   * Get monitor dashboard overview data
   */
  getDashboard: (params?: MonitorTimeRange): Promise<MonitorDashboardData> =>
    apiClient.get('/monitor/dashboard', { params }),

  /**
   * Get system metrics summary
   */
  getSystemMetrics: (params?: MonitorTimeRange): Promise<SystemMetricsSummary> =>
    apiClient.get('/monitor/system', { params }),

  /**
   * Get latest system metric
   */
  getLatestSystemMetric: (): Promise<SystemMetric> =>
    apiClient.get('/monitor/system/latest'),

  /**
   * Get API metrics summary
   */
  getAPIMetrics: (params?: MonitorTimeRange): Promise<APIMetricsSummary> =>
    apiClient.get('/monitor/api', { params }),

  /**
   * Get logs with filtering
   */
  getLogs: (params?: LogFilterParams): Promise<{
    logs: LogEntry[];
    total: number;
    page: number;
    page_size: number;
  }> => apiClient.get('/monitor/logs', { params }),

  /**
   * Get log sources/modules for filtering
   */
  getLogSources: (): Promise<{ sources: string[]; modules: string[] }> =>
    apiClient.get('/monitor/logs/sources'),

  /**
   * Stream logs (Server-Sent Events endpoint)
   * Note: This returns an EventSource instance for real-time logs
   */
  streamLogs: (params?: LogFilterParams): string => {
    const queryParams = new URLSearchParams();
    if (params?.level) queryParams.append('level', params.level);
    if (params?.source) queryParams.append('source', params.source);
    if (params?.keyword) queryParams.append('keyword', params.keyword);
    return `${BASE_URL}/monitor/logs/stream?${queryParams.toString()}`;
  },

  /**
   * Get all alert rules
   */
  getAlertRules: (): Promise<AlertRule[]> =>
    apiClient.get('/monitor/alerts/rules'),

  /**
   * Get alert rule by ID
   */
  getAlertRule: (id: string): Promise<AlertRule> =>
    apiClient.get(`/monitor/alerts/rules/${id}`),

  /**
   * Create alert rule
   */
  createAlertRule: (data: CreateAlertRuleRequest): Promise<AlertRule> =>
    apiClient.post('/monitor/alerts/rules', data),

  /**
   * Update alert rule
   */
  updateAlertRule: (id: string, data: UpdateAlertRuleRequest): Promise<AlertRule> =>
    apiClient.put(`/monitor/alerts/rules/${id}`, data),

  /**
   * Delete alert rule
   */
  deleteAlertRule: (id: string): Promise<void> =>
    apiClient.delete(`/monitor/alerts/rules/${id}`),

  /**
   * Toggle alert rule active status
   */
  toggleAlertRule: (id: string, isActive: boolean): Promise<AlertRule> =>
    apiClient.patch(`/monitor/alerts/rules/${id}/toggle`, { is_active: isActive }),

  /**
   * Get alert history
   */
  getAlertHistory: (params?: AlertFilterParams): Promise<{
    alerts: AlertHistory[];
    total: number;
    page: number;
    page_size: number;
  }> => apiClient.get('/monitor/alerts/history', { params }),

  /**
   * Get active alerts
   */
  getActiveAlerts: (): Promise<AlertHistory[]> =>
    apiClient.get('/monitor/alerts/active'),

  /**
   * Acknowledge an alert
   */
  acknowledgeAlert: (id: string): Promise<AlertHistory> =>
    apiClient.post(`/monitor/alerts/${id}/acknowledge`),

  /**
   * Resolve an alert
   */
  resolveAlert: (id: string): Promise<AlertHistory> =>
    apiClient.post(`/monitor/alerts/${id}/resolve`),

  /**
   * Get health check status
   */
  getHealthStatus: (): Promise<{
    status: 'healthy' | 'degraded' | 'unhealthy';
    checks: {
      name: string;
      status: 'pass' | 'fail' | 'warn';
      message?: string;
      last_check: string;
    }[];
  }> => apiClient.get('/health'),
};

export default monitorAPI;
