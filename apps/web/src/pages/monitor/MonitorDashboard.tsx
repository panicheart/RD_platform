import React, { useEffect, useState, useCallback } from 'react';
import { Typography, Row, Col, DatePicker, Button, Space, message, Tabs, Tag, Badge } from 'antd';
import {
  DesktopOutlined,
  ApiOutlined,
  CodeOutlined,
  BellOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import SystemMetrics from '@/components/monitor/SystemMetrics';
import APIMetrics from '@/components/monitor/APIMetrics';
import LogViewer from '@/components/monitor/LogViewer';
import AlertManager from '@/components/monitor/AlertManager';
import { monitorAPI } from '@/services/monitor';
import type {
  SystemMetricsSummary,
  APIMetricsSummary,
  LogEntry,
  AlertRule,
  AlertHistory,
  MonitorTimeRange,
  LogFilterParams,
  CreateAlertRuleRequest,
  UpdateAlertRuleRequest,
} from '@/types/monitor';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;
const { TabPane } = Tabs;

const MonitorDashboard: React.FC = () => {
  const [loading, setLoading] = useState({
    dashboard: false,
    system: false,
    api: false,
    logs: false,
    alerts: false,
  });

  const [systemMetrics, setSystemMetrics] = useState<SystemMetricsSummary | undefined>();
  const [apiMetrics, setApiMetrics] = useState<APIMetricsSummary | undefined>();
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [logTotal, setLogTotal] = useState(0);
  const [logSources, setLogSources] = useState<string[]>([]);
  const [logModules, setLogModules] = useState<string[]>([]);
  const [alertRules, setAlertRules] = useState<AlertRule[]>([]);
  const [alertHistory, setAlertHistory] = useState<AlertHistory[]>([]);
  const [activeAlerts, setActiveAlerts] = useState<AlertHistory[]>([]);
  const [healthStatus, setHealthStatus] = useState<{
    status: 'healthy' | 'degraded' | 'unhealthy';
    checks: any[];
  } | null>(null);

  const [timeRange, setTimeRange] = useState<MonitorTimeRange>({ range: '1h' });
  const [logFilters, setLogFilters] = useState<LogFilterParams>({ page: 1, page_size: 50 });
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

  const fetchDashboard = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, dashboard: true }));
      const data = await monitorAPI.getDashboard(timeRange);
      setSystemMetrics(data.system);
      setApiMetrics(data.api);
      setActiveAlerts(data.alerts.recent_alerts || []);
      setHealthStatus(data.health);
      setLastUpdate(new Date());
    } catch (error) {
      console.error('Failed to fetch dashboard:', error);
    } finally {
      setLoading((prev) => ({ ...prev, dashboard: false }));
    }
  }, [timeRange]);

  const fetchSystemMetrics = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, system: true }));
      const data = await monitorAPI.getSystemMetrics(timeRange);
      setSystemMetrics(data);
    } catch (error) {
      console.error('Failed to fetch system metrics:', error);
    } finally {
      setLoading((prev) => ({ ...prev, system: false }));
    }
  }, [timeRange]);

  const fetchAPIMetrics = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, api: true }));
      const data = await monitorAPI.getAPIMetrics(timeRange);
      setApiMetrics(data);
    } catch (error) {
      console.error('Failed to fetch API metrics:', error);
    } finally {
      setLoading((prev) => ({ ...prev, api: false }));
    }
  }, [timeRange]);

  const fetchLogs = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, logs: true }));
      const data = await monitorAPI.getLogs(logFilters);
      if (logFilters.page === 1) {
        setLogs(data.logs);
      } else {
        setLogs((prev) => [...prev, ...data.logs]);
      }
      setLogTotal(data.total);
    } catch (error) {
      console.error('Failed to fetch logs:', error);
      void message.error('获取日志失败');
    } finally {
      setLoading((prev) => ({ ...prev, logs: false }));
    }
  }, [logFilters]);

  const fetchLogSources = useCallback(async () => {
    try {
      const data = await monitorAPI.getLogSources();
      setLogSources(data.sources);
      setLogModules(data.modules);
    } catch (error) {
      console.error('Failed to fetch log sources:', error);
    }
  }, []);

  const fetchAlertRules = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, alerts: true }));
      const rules = await monitorAPI.getAlertRules();
      setAlertRules(rules);
    } catch (error) {
      console.error('Failed to fetch alert rules:', error);
    }
  }, []);

  const fetchAlertHistory = useCallback(async () => {
    try {
      setLoading((prev) => ({ ...prev, alerts: true }));
      const data = await monitorAPI.getAlertHistory();
      setAlertHistory(data.alerts);
      const active = await monitorAPI.getActiveAlerts();
      setActiveAlerts(active);
    } catch (error) {
      console.error('Failed to fetch alert history:', error);
    } finally {
      setLoading((prev) => ({ ...prev, alerts: false }));
    }
  }, []);

  useEffect(() => {
    void fetchDashboard();
    void fetchLogs();
    void fetchLogSources();
    void fetchAlertRules();
    void fetchAlertHistory();
  }, []);

  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      void fetchDashboard();
    }, 5000);

    return () => clearInterval(interval);
  }, [autoRefresh, fetchDashboard]);

  const handleTimeRangeChange = (dates: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      const newRange: MonitorTimeRange = {
        start_time: dates[0].toISOString(),
        end_time: dates[1].toISOString(),
      };
      setTimeRange(newRange);
      void fetchDashboard();
      void fetchSystemMetrics();
      void fetchAPIMetrics();
    }
  };

  const handleQuickRange = (range: '1h' | '6h' | '24h' | '7d') => {
    setTimeRange({ range });
    void fetchDashboard();
    void fetchSystemMetrics();
    void fetchAPIMetrics();
  };

  const handleLogFilterChange = (filters: LogFilterParams) => {
    setLogFilters({ ...filters, page: 1 });
    setTimeout(() => void fetchLogs(), 0);
  };

  const handleLoadMoreLogs = () => {
    setLogFilters((prev) => ({ ...prev, page: (prev.page || 1) + 1 }));
    setTimeout(() => void fetchLogs(), 0);
  };

  const handleCreateRule = async (data: CreateAlertRuleRequest) => {
    await monitorAPI.createAlertRule(data);
    void fetchAlertRules();
  };

  const handleUpdateRule = async (id: string, data: UpdateAlertRuleRequest) => {
    await monitorAPI.updateAlertRule(id, data);
    void fetchAlertRules();
  };

  const handleDeleteRule = async (id: string) => {
    await monitorAPI.deleteAlertRule(id);
    void fetchAlertRules();
  };

  const handleToggleRule = async (id: string, isActive: boolean) => {
    await monitorAPI.toggleAlertRule(id, isActive);
    void fetchAlertRules();
  };

  const handleResolveAlert = async (id: string) => {
    await monitorAPI.resolveAlert(id);
    void fetchAlertHistory();
    void fetchDashboard();
  };

  const getHealthColor = (status: string): string => {
    switch (status) {
      case 'healthy':
        return '#52c41a';
      case 'degraded':
        return '#faad14';
      case 'unhealthy':
        return '#ff4d4f';
      default:
        return '#d9d9d9';
    }
  };

  const getHealthText = (status: string): string => {
    switch (status) {
      case 'healthy':
        return '健康';
      case 'degraded':
        return '降级';
      case 'unhealthy':
        return '异常';
      default:
        return '未知';
    }
  };

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Title level={4}>系统监控</Title>
          <Space direction="vertical" size={0}>
            <Text type="secondary">实时监控系统运行状态</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              最后更新: {lastUpdate.toLocaleTimeString()}
            </Text>
          </Space>
        </Col>
        <Col>
          <Space>
            {healthStatus && (
              <Tag
                icon={<CheckCircleOutlined />}
                color={getHealthColor(healthStatus.status)}
              >
                {getHealthText(healthStatus.status)}
              </Tag>
            )}
            <Space.Compact>
              {(['1h', '6h', '24h', '7d'] as const).map((range) => (
                <Button
                  key={range}
                  size="small"
                  type={timeRange.range === range ? 'primary' : 'default'}
                  onClick={() => handleQuickRange(range)}
                >
                  {range === '1h' && '1小时'}
                  {range === '6h' && '6小时'}
                  {range === '24h' && '24小时'}
                  {range === '7d' && '7天'}
                </Button>
              ))}
            </Space.Compact>

            <RangePicker
              showTime
              style={{ width: 320 }}
              placeholder={['开始时间', '结束时间']}
              onChange={handleTimeRangeChange}
            />

            <Button
              icon={<ReloadOutlined spin={loading.dashboard} />}
              onClick={() => void fetchDashboard()}
            >
              刷新
            </Button>
          </Space>
        </Col>
      </Row>
      {activeAlerts.length > 0 && (
        <div
          style={{
            marginBottom: 24,
            padding: 16,
            backgroundColor: '#fff2f0',
            border: '1px solid #ffccc7',
            borderRadius: 6,
          }}
        >
          <Row align="middle" gutter={16}>
            <Col>
              <ExclamationCircleOutlined style={{ fontSize: 24, color: '#ff4d4f' }} />
            </Col>
            <Col flex="auto">
              <div style={{ fontWeight: 500, fontSize: 16 }}>
                当前有 {activeAlerts.length} 个活跃告警
              </div>
              <Space size="large" style={{ marginTop: 4 }}>
                {activeAlerts.filter((a) => a.severity === 'critical').length > 0 && (
                  <Badge
                    count={activeAlerts.filter((a) => a.severity === 'critical').length}
                    style={{ backgroundColor: '#ff4d4f' }}
                  >
                    <span style={{ color: '#ff4d4f' }}>严重告警</span>
                  </Badge>
                )}
                {activeAlerts.filter((a) => a.severity === 'warning').length > 0 && (
                  <Badge
                    count={activeAlerts.filter((a) => a.severity === 'warning').length}
                    style={{ backgroundColor: '#faad14' }}
                  >
                    <span style={{ color: '#faad14' }}>警告告警</span>
                  </Badge>
                )}
              </Space>
            </Col>
            <Col>
              <Button type="primary" danger onClick={() => setActiveTab('alerts')}>
                查看详情
              </Button>
            </Col>
          </Row>
        </div>
      )}
      <Tabs activeKey={activeTab} onChange={setActiveTab} type="card">
        <TabPane
          tab={
            <span>
              <DesktopOutlined /> 系统监控
            </span>
          }
          key="overview"
        >
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <SystemMetrics
                data={systemMetrics}
                loading={loading.system || loading.dashboard}
                autoRefresh={autoRefresh}
                refreshInterval={5000}
                onRefresh={fetchSystemMetrics}
              />
            </Col>
          </Row>
        </TabPane>

        <TabPane
          tab={
            <span>
              <ApiOutlined /> API监控
            </span>
          }
          key="api"
        >
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <APIMetrics
                data={apiMetrics}
                loading={loading.api || loading.dashboard}
                timeRange={timeRange.range}
              />
            </Col>
          </Row>
        </TabPane>

        <TabPane
          tab={
            <span>
              <CodeOutlined /> 日志查看
              {autoRefresh && <Badge status="processing" style={{ marginLeft: 4 }} />}
            </span>
          }
          key="logs"
        >
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <LogViewer
                logs={logs}
                loading={loading.logs}
                total={logTotal}
                sources={logSources}
                modules={logModules}
                onFilterChange={handleLogFilterChange}
                onRefresh={fetchLogs}
                onLoadMore={handleLoadMoreLogs}
                autoRefresh={autoRefresh}
                onAutoRefreshChange={setAutoRefresh}
              />
            </Col>
          </Row>
        </TabPane>

        <TabPane
          tab={
            <span>
              <BellOutlined /> 告警管理
              {activeAlerts.length > 0 && (
                <Badge count={activeAlerts.length} style={{ marginLeft: 4 }} />
              )}
            </span>
          }
          key="alerts"
        >
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <AlertManager
                rules={alertRules}
                history={alertHistory}
                activeAlerts={activeAlerts}
                loading={loading.alerts}
                onCreateRule={handleCreateRule}
                onUpdateRule={handleUpdateRule}
                onDeleteRule={handleDeleteRule}
                onToggleRule={handleToggleRule}
                onResolveAlert={handleResolveAlert}
                onRefresh={fetchAlertHistory}
              />
            </Col>
          </Row>
        </TabPane>
      </Tabs>
    </div>
  );
};

export default MonitorDashboard;
