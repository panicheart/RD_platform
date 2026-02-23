// Monitor types for system monitoring dashboard
// Following backend models: services/api/models/monitor.go

/**
 * System metric data point
 */
export interface SystemMetric {
  id: string;
  timestamp: string;
  cpu_usage: number;
  memory_usage: number;
  memory_total: number;
  memory_used: number;
  disk_usage: number;
  disk_total: number;
  disk_used: number;
  network_in: number;
  network_out: number;
  db_connections: number;
  api_requests: number;
  created_at: string;
}

/**
 * API metric data point
 */
export interface APIMetric {
  id: string;
  timestamp: string;
  endpoint: string;
  method: string;
  duration: number;
  status_code: number;
  user_id?: string;
  ip_address: string;
  created_at: string;
}

/**
 * Log entry data
 */
export interface LogEntry {
  id: string;
  timestamp: string;
  level: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR';
  message: string;
  source: string;
  module: string;
  user_id?: string;
  request_id?: string;
  metadata?: string;
  created_at: string;
}

/**
 * Alert rule definition
 */
export interface AlertRule {
  id: string;
  name: string;
  description?: string;
  metric: string;
  condition: '>' | '<' | '==' | '!=';
  threshold: number;
  duration: number;
  severity: 'warning' | 'critical';
  is_active: boolean;
  notify_channels: string[];
  created_at: string;
  updated_at: string;
}

/**
 * Alert history entry
 */
export interface AlertHistory {
  id: string;
  rule_id: string;
  rule_name: string;
  severity: 'warning' | 'critical';
  message: string;
  value: number;
  threshold: number;
  status: 'firing' | 'resolved';
  resolved_at?: string;
  created_at: string;
}

/**
 * Time series data point for charts
 */
export interface TimeSeriesData {
  timestamp: string;
  value: number;
  label?: string;
}

/**
 * System metrics summary for display
 */
export interface SystemMetricsSummary {
  current: {
    cpu_usage: number;
    memory_usage: number;
    disk_usage: number;
    memory_total: number;
    memory_used: number;
    disk_total: number;
    disk_used: number;
    network_in_rate: number;
    network_out_rate: number;
    db_connections: number;
  };
  history: {
    cpu: TimeSeriesData[];
    memory: TimeSeriesData[];
    disk: TimeSeriesData[];
    network: TimeSeriesData[];
  };
}

/**
 * API metrics summary for display
 */
export interface APIMetricsSummary {
  total_requests: number;
  avg_response_time: number;
  error_rate: number;
  requests_per_second: number;
  top_endpoints: {
    endpoint: string;
    method: string;
    count: number;
    avg_duration: number;
  }[];
  response_time_trend: TimeSeriesData[];
  request_count_trend: TimeSeriesData[];
  status_distribution: {
    status_code: number;
    count: number;
    percentage: number;
  }[];
}

/**
 * Log filter params
 */
export interface LogFilterParams {
  level?: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR';
  source?: string;
  module?: string;
  start_time?: string;
  end_time?: string;
  keyword?: string;
  page?: number;
  page_size?: number;
}

/**
 * Alert filter params
 */
export interface AlertFilterParams {
  severity?: 'warning' | 'critical';
  status?: 'firing' | 'resolved';
  start_time?: string;
  end_time?: string;
  page?: number;
  page_size?: number;
}

/**
 * Create alert rule request
 */
export interface CreateAlertRuleRequest {
  name: string;
  description?: string;
  metric: string;
  condition: '>' | '<' | '==' | '!=';
  threshold: number;
  duration: number;
  severity: 'warning' | 'critical';
  notify_channels: string[];
}

/**
 * Update alert rule request
 */
export interface UpdateAlertRuleRequest {
  name?: string;
  description?: string;
  metric?: string;
  condition?: '>' | '<' | '==' | '!=';
  threshold?: number;
  duration?: number;
  severity?: 'warning' | 'critical';
  is_active?: boolean;
  notify_channels?: string[];
}

/**
 * Monitor dashboard data
 */
export interface MonitorDashboardData {
  system: SystemMetricsSummary;
  api: APIMetricsSummary;
  alerts: {
    active_count: number;
    warning_count: number;
    critical_count: number;
    recent_alerts: AlertHistory[];
  };
  health: {
    status: 'healthy' | 'degraded' | 'unhealthy';
    checks: {
      name: string;
      status: 'pass' | 'fail' | 'warn';
      message?: string;
      last_check: string;
    }[];
  };
}

/**
 * Time range for monitor queries
 */
export interface MonitorTimeRange {
  start_time?: string;
  end_time?: string;
  range?: '1h' | '6h' | '24h' | '7d' | '30d';
}
