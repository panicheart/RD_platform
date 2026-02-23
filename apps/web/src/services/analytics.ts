import { apiClient } from './api';
import type {
  DashboardData,
  ProjectStats,
  UserStats,
  ShelfStats,
  KnowledgeStats,
  AnalyticsTimeRange,
  ReportTemplate,
  ExportRequest,
  TrendData,
  ProjectStatusDistribution,
} from '../types/analytics';

/**
 * Analytics API service for data analysis dashboard
 */
export const analyticsAPI = {
  /**
   * Get dashboard overview data
   */
  getDashboard: (): Promise<DashboardData> =>
    apiClient.get('/analytics/dashboard'),

  /**
   * Get project statistics with optional time range filter
   */
  getProjectStats: (params?: AnalyticsTimeRange): Promise<{
    stats: ProjectStats[];
    distribution: ProjectStatusDistribution[];
    trends: TrendData[];
  }> => apiClient.get('/analytics/projects', { params }),

  /**
   * Get user statistics with optional time range filter
   */
  getUserStats: (params?: AnalyticsTimeRange): Promise<{
    stats: UserStats[];
    trends: TrendData[];
    heatmap: { date: string; count: number }[];
  }> => apiClient.get('/analytics/users', { params }),

  /**
   * Get shelf/product statistics with optional time range filter
   */
  getShelfStats: (params?: AnalyticsTimeRange): Promise<ShelfStats> =>
    apiClient.get('/analytics/shelf', { params }),

  /**
   * Get knowledge base statistics with optional time range filter
   */
  getKnowledgeStats: (params?: AnalyticsTimeRange): Promise<KnowledgeStats> =>
    apiClient.get('/analytics/knowledge', { params }),

  /**
   * Get report templates
   */
  getReportTemplates: (): Promise<ReportTemplate[]> =>
    apiClient.get('/analytics/report-templates'),

  /**
   * Export report
   */
  exportReport: (data: ExportRequest): Promise<Blob> =>
    apiClient.post('/analytics/exports', data),
};

export default analyticsAPI;
