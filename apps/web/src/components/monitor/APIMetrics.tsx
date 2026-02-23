import React, { useEffect, useRef } from 'react';
import { Card, Row, Col, Statistic, Empty, Spin, Table, Tag } from 'antd';
import {
  ApiOutlined,
  ClockCircleOutlined,
  WarningOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons';
import * as echarts from 'echarts';
import type { APIMetricsSummary } from '@/types/monitor';

interface APIMetricsProps {
  data?: APIMetricsSummary;
  loading?: boolean;
  timeRange?: string;
}

const APIMetrics: React.FC<APIMetricsProps> = ({
  data,
  loading = false,
  timeRange = '1h',
}) => {
  const responseTimeChartRef = useRef<HTMLDivElement>(null);
  const statusChartRef = useRef<HTMLDivElement>(null);
  const responseTimeChartInstance = useRef<echarts.ECharts | null>(null);
  const statusChartInstance = useRef<echarts.ECharts | null>(null);

  // Initialize Response Time Chart
  useEffect(() => {
    if (!responseTimeChartRef.current || loading || !data?.response_time_trend) return;

    if (!responseTimeChartInstance.current) {
      responseTimeChartInstance.current = echarts.init(responseTimeChartRef.current);
    }

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const time = new Date(params[0].value[0]).toLocaleTimeString();
          const value = params[0].value[1];
          return `${time}<br/>响应时间: ${value.toFixed(0)}ms`;
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
        name: 'ms',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#666' },
      },
      series: [
        {
          name: 'Response Time',
          type: 'line',
          smooth: true,
          symbol: 'none',
          lineStyle: { width: 2, color: '#722ed1' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(114, 46, 209, 0.3)' },
              { offset: 1, color: 'rgba(114, 46, 209, 0.05)' },
            ]),
          },
          data: data.response_time_trend.map((item) => [item.timestamp, item.value]),
        },
      ],
    };

    responseTimeChartInstance.current.setOption(option);

    const handleResize = () => {
      responseTimeChartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [data?.response_time_trend, loading]);

  // Initialize Status Distribution Chart
  useEffect(() => {
    if (!statusChartRef.current || loading || !data?.status_distribution) return;

    if (!statusChartInstance.current) {
      statusChartInstance.current = echarts.init(statusChartRef.current);
    }

    const statusData = data.status_distribution.map((item) => ({
      name: `${item.status_code}`,
      value: item.count,
      itemStyle: {
        color:
          item.status_code < 300
            ? '#52c41a'
            : item.status_code < 400
              ? '#faad14'
              : '#ff4d4f',
      },
    }));

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'item',
        formatter: (params: any) => {
          const percent = data.status_distribution.find(
            (d) => d.status_code.toString() === params.name
          )?.percentage;
          return `状态码 ${params.name}<br/>数量: ${params.value}<br/>占比: ${percent?.toFixed(1)}%`;
        },
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center',
      },
      series: [
        {
          name: 'Status Code',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['40%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold',
            },
          },
          labelLine: {
            show: false,
          },
          data: statusData,
        },
      ],
    };

    statusChartInstance.current.setOption(option);

    const handleResize = () => {
      statusChartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [data?.status_distribution, loading]);

  // Cleanup charts
  useEffect(() => {
    return () => {
      responseTimeChartInstance.current?.dispose();
      responseTimeChartInstance.current = null;
      statusChartInstance.current?.dispose();
      statusChartInstance.current = null;
    };
  }, []);

  // Get status color
  const getErrorColor = (rate: number): string => {
    if (rate >= 5) return '#ff4d4f';
    if (rate >= 1) return '#faad14';
    return '#52c41a';
  };

  const getResponseTimeColor = (time: number): string => {
    if (time >= 1000) return '#ff4d4f';
    if (time >= 500) return '#faad14';
    return '#52c41a';
  };

  // Table columns for top endpoints
  const columns = [
    {
      title: '端点',
      dataIndex: 'endpoint',
      key: 'endpoint',
      ellipsis: true,
      render: (text: string, record: any) => (
        <span>
          <Tag color="blue">{record.method}</Tag> {text}
        </span>
      ),
    },
    {
      title: '请求数',
      dataIndex: 'count',
      key: 'count',
      sorter: (a: any, b: any) => a.count - b.count,
      render: (value: number) => value.toLocaleString(),
    },
    {
      title: '平均响应时间',
      dataIndex: 'avg_duration',
      key: 'avg_duration',
      sorter: (a: any, b: any) => a.avg_duration - b.avg_duration,
      render: (value: number) => (
        <span style={{ color: getResponseTimeColor(value) }}>{value.toFixed(0)} ms</span>
      ),
    },
  ];

  if (loading) {
    return (
      <Card title="API 指标">
        <div style={{ height: 400, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (!data) {
    return (
      <Card title="API 指标">
        <div style={{ height: 400, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Empty description="暂无API指标数据" />
        </div>
      </Card>
    );
  }

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <ApiOutlined />
          <span>API 指标</span>
          {timeRange && <Tag>{timeRange}</Tag>}
        </div>
      }
    >
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        {/* Total Requests */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="总请求数"
              value={data.total_requests}
              prefix={<ThunderboltOutlined />}
            />
          </Card>
        </Col>

        {/* Avg Response Time */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="平均响应时间"
              value={data.avg_response_time}
              suffix="ms"
              precision={0}
              valueStyle={{ color: getResponseTimeColor(data.avg_response_time) }}
              prefix={<ClockCircleOutlined />}
            />
          </Card>
        </Col>

        {/* Error Rate */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="错误率"
              value={data.error_rate}
              suffix="%"
              precision={2}
              valueStyle={{ color: getErrorColor(data.error_rate) }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>

        {/* Requests Per Second */}
        <Col xs={24} sm={12} lg={6}>
          <Card size="small">
            <Statistic
              title="每秒请求数"
              value={data.requests_per_second}
              suffix="req/s"
              precision={1}
              prefix={<ApiOutlined />}
            />
          </Card>
        </Col>

        {/* Response Time Trend */}
        <Col xs={24} lg={12}>
          <Card title="响应时间趋势" size="small">
            <div ref={responseTimeChartRef} style={{ height: 250 }} />
          </Card>
        </Col>

        {/* Status Distribution */}
        <Col xs={24} lg={12}>
          <Card title="状态码分布" size="small">
            <div ref={statusChartRef} style={{ height: 250 }} />
          </Card>
        </Col>

        {/* Top Endpoints */}
        <Col span={24}>
          <Card title="热门端点" size="small">
            <Table
              columns={columns}
              dataSource={data.top_endpoints || []}
              rowKey={(record) => `${record.method}-${record.endpoint}`}
              pagination={false}
              size="small"
              scroll={{ x: 'max-content' }}
            />
          </Card>
        </Col>
      </Row>
    </Card>
  );
};

export default APIMetrics;
