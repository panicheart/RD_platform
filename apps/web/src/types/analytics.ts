// Analytics types for data analysis dashboard
// Following backend models: services/api/models/analytics.go

/**
 * Dashboard overview data
 */
export interface DashboardOverview {
  totalProjects: number;
  activeUsers: number;
  completedTasks: number;
  weeklyNewItems: number;
  avgProgress: number;
  onTimeRate: number;
}

/**
 * Project statistics
 */
export interface ProjectStats {
  id: string;
  date: string;
  totalProjects: number;
  activeProjects: number;
  completedProjects: number;
  delayedProjects: number;
  avgProgress: number;
}

/**
 * Project status distribution for charts
 */
export interface ProjectStatusDistribution {
  status: string;
  count: number;
  percentage: number;
}

/**
 * User statistics
 */
export interface UserStats {
  id: string;
  userId: string;
  date: string;
  tasksCompleted: number;
  tasksCreated: number;
  projectsJoined: number;
  reviewsDone: number;
  contribution: number;
  workHours: number;
}

/**
 * User contribution heatmap data
 * Format: [date, value] where date is YYYY-MM-DD
 */
export interface ContributionHeatmapData {
  date: string;
  count: number;
}

/**
 * Trend data for time series charts
 */
export interface TrendData {
  date: string;
  value: number;
  label?: string;
}

/**
 * Shelf statistics
 */
export interface ShelfStats {
  totalProducts: number;
  adoptedProducts: number;
  adoptionRate: number;
  avgReuses: number;
  trends: TrendData[];
}

/**
 * Knowledge statistics
 */
export interface KnowledgeStats {
  totalArticles: number;
  totalViews: number;
  totalLikes: number;
  topArticles: {
    id: string;
    title: string;
    views: number;
    likes: number;
  }[];
}

/**
 * Time range filter params
 */
export interface AnalyticsTimeRange {
  startDate?: string;
  endDate?: string;
}

/**
 * Dashboard data response
 */
export interface DashboardData {
  overview: DashboardOverview;
  projectDistribution: ProjectStatusDistribution[];
  projectTrends: TrendData[];
  userHeatmap: ContributionHeatmapData[];
  userTrends: TrendData[];
  shelfStats: ShelfStats;
  knowledgeStats: KnowledgeStats;
}

/**
 * Report template
 */
export interface ReportTemplate {
  id: string;
  name: string;
  description?: string;
  type: 'project' | 'user' | 'system';
  format: 'pdf' | 'excel';
  template: string;
  isActive: boolean;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Export request params
 */
export interface ExportRequest {
  templateId?: string;
  format: 'pdf' | 'excel';
  startDate?: string;
  endDate?: string;
  filters?: Record<string, unknown>;
}
