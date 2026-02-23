import React, { useEffect, useRef } from 'react';
import { Card, Row, Col, Statistic, Empty, Spin, Badge } from 'antd';
import {
  DesktopOutlined,
  DatabaseOutlined,
  CloudOutlined,
  SwapOutlined,
} from '@ant-design/icons';
import * as echarts from 'echarts';
import type { SystemMetricsSummary } from '@/types/monitor';

interface SystemMetricsProps {
  data?: SystemMetricsSummary;
  loading?: boolean;
  autoRefresh?: boolean;
  refreshInterval?: number;
  onRefresh?: () => void;
}

const SystemMetrics: React.FC<SystemMetricsProps> = ({
  data,
  loading = false,
  autoRefresh = true,
  refreshInterval = 5000,
  onRefresh,
}) => {
  const cpuChartRef = useRef<HTMLDivElement>(null);
  const memoryChartRef = useRef<HTMLDivElement>(null);
  const cpuChartInstance = useRef<echarts.ECharts | null>(null);
  const memoryChartInstance = useRef<echarts.ECharts | null>(null);

  // Format bytes to human-readable
  const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Get status color based on usage
  const getStatusColor = (usage: number): string => {
    if (usage >= 90) return '#ff4d4f';
    if (usage >= 70) return '#faad14';
    return '#52c41a';
  };

  // Initialize CPU chart
  useEffect(() => {
    if (!cpuChartRef.current || loading || !data?.history.cpu) return;

    if (!cpuChartInstance.current) {
      cpuChartInstance.current = echarts.init(cpuChartRef.current);
    }

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const time = new Date(params[0].value[0]).toLocaleTimeString();
          const value = params[0].value[1];
          return `${time}<br/>CPU: ${value.toFixed(1)}%`;
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: '10%',
        containLabel: true,
      },
      xAxis: {
        type: 'time',
        axisLine: { lineStyle: { color: '#d9d9d9' } },
        axisLabel: {
          color: '#666',
          formatter: (value: number) => {
            return new Date(value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
          },
        },
        splitLine: { show: false },
      },
      yAxis: {
        type: 'value',
        min: 0,
        max: 100,
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#666', formatter: '{value}%' },
      },
      series: [
        {
          name: 'CPU',
          type: 'line',
          smooth: true,
          symbol: 'none',
          lineStyle: { width: 2, color: '#1890ff' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
              { offset: 1, color: 'rgba(24, 144, 255, 0.05)' },
            ]),
          },
          data: data.history.cpu.map((item) => [item.timestamp, item.value]),
        },
      ],
    };

    cpuChartInstance.current.setOption(option);

    const handleResize = () => {
      cpuChartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [data?.history.cpu, loading]);

  // Initialize Memory chart
  useEffect(() => {
    if (!memoryChartRef.current || loading || !data?.history.memory) return;

    if (!memoryChartInstance.current) {
      memoryChartInstance.current = echarts.init(memoryChartRef.current);
    }

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const time = new Date(params[0].value[0]).toLocaleTimeString();
          const value = params[0].value[1];
          return `${time}<br/>内存: ${value.toFixed(1)}%`;
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: '10%',
        containLabel: true,
      },
      xAxis: {
        type: 'time',
        axisLine: { lineStyle: { color: '#d9d9d9' } },
        axisLabel: {
          color: '#666',
          formatter: (value: number) => {
            return new Date(value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
          },
        },
        splitLine: { show: false },
      },
      yAxis: {
        type: 'value',
        min: 0,
        max: 100,
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#666', formatter: '{value}%' },
      },
      series: [
        {
          name: 'Memory',
          type: 'line',
          smooth: true,
          symbol: 'none',
          lineStyle: { width: 2, color: '#52c41a' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(82, 196, 26, 0.3)' },
              { offset: 1, color: 'rgba(82, 196, 26, 0.05)' },
            ]),
          },
          data: data.history.memory.map((item) => [item.timestamp, item.value]),
        },
      ],
    };

    memoryChartInstance.current.setOption(option);

    const handleResize = () => {
      memoryChartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [data?.history.memory, loading]);

  // Auto refresh
  useEffect(() => {
    if (!autoRefresh || !onRefresh) return;

    const interval = setInterval(() => {
      onRefresh();
    }, refreshInterval);

    return () => clearInterval(interval);
  }, [autoRefresh, refreshInterval, onRefresh]);

  // Cleanup charts
  useEffect(() => {
    return () => {
      cpuChartInstance.current?.dispose();
      cpuChartInstance.current = null;
      memoryChartInstance.current?.dispose();
      memoryChartInstance.current = null;
    };
  }, []);

  const current = data?.current;

  if (loading) {
    return (
      <Card title="系统指标" extra={<Badge status="processing" text="加载中..." />}>
        <div style={{ height: 400, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (!data || !current) {
    return (
      <Card title="系统指标">
        <div style={{ height: 400, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Empty description="暂无系统指标数据" />
        </div>
      </Card>
    );
  }

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <DesktopOutlined />
          <span>系统指标</span>
          {autoRefresh && <Badge status="success" text="实时" />}
        </div>
      }
    >
      <Row gutter={[16, 16]}>
        {/* CPU */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="CPU 使用率"
              value={current.cpu_usage}
              suffix="%"
              precision={1}
              valueStyle={{ color: getStatusColor(current.cpu_usage) }}
              prefix={<DesktopOutlined />}
            />
          </Card>
        </Col>

        {/* Memory */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="内存使用率"
              value={current.memory_usage}
              suffix="%"
              precision={1}
              valueStyle={{ color: getStatusColor(current.memory_usage) }}
              prefix={<DatabaseOutlined />}
            />
            <div style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
              {formatBytes(current.memory_used)} / {formatBytes(current.memory_total)}
            </div>
          </Card>
        </Col>

        {/* Disk */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="磁盘使用率"
              value={current.disk_usage}
              suffix="%"
              precision={1}
              valueStyle={{ color: getStatusColor(current.disk_usage) }}
              prefix={<CloudOutlined />}
            />
            <div style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
              {formatBytes(current.disk_used)} / {formatBytes(current.disk_total)}
            </div>
          </Card>
        </Col>

        {/* Network */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="数据库连接"
              value={current.db_connections}
              prefix={<SwapOutlined />}
            />
            <div style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
              入: {formatBytes(current.network_in_rate || 0)}/s | 出: {formatBytes(current.network_out_rate || 0)}/s
            </div>
          </Card>
        </Col>

        {/* CPU Chart */}
        <Col xs={24} lg={12}>
          <Card title="CPU 使用趋势" size="small">
            <div ref={cpuChartRef} style={{ height: 250 }} />
          </Card>
        </Col>

        {/* Memory Chart */}
        <Col xs={24} lg={12}>
          <Card title="内存使用趋势" size="small">
            <div ref={memoryChartRef} style={{ height: 250 }} />
          </Card>
        </Col>
      </Row>
    </Card>
  );
};

export default SystemMetrics;
