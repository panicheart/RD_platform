import React, { useEffect, useState } from 'react';
import { Typography, Row, Col, DatePicker, Button, Space, message } from 'antd';
import {
  TeamOutlined,
  ProjectOutlined,
  CheckCircleOutlined,
  RiseOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import StatCard from '@/components/analytics/StatCard';
import HeatmapChart from '@/components/analytics/HeatmapChart';
import PieChart from '@/components/analytics/PieChart';
import LineChart from '@/components/analytics/LineChart';
import { analyticsAPI } from '@/services/analytics';
import type { DashboardData, AnalyticsTimeRange } from '@/types/analytics';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

const AnalyticsDashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<DashboardData | null>(null);
  const [dateRange, setDateRange] = useState<AnalyticsTimeRange>({
    startDate: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
    endDate: dayjs().format('YYYY-MM-DD'),
  });

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      const response = await analyticsAPI.getDashboard();
      setData(response);
    } catch (error) {
      void message.error('获取仪表盘数据失败');
      console.error('Failed to fetch dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchDashboardData();
  }, []);

  const handleDateRangeChange = (dates: [dayjs.Dayjs | null, dayjs.Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setDateRange({
        startDate: dates[0].format('YYYY-MM-DD'),
        endDate: dates[1].format('YYYY-MM-DD'),
      });
    }
  };

  const handleGenerateReport = () => {
    void message.info('报表生成功能开发中');
  };

  const getProjectStatusData = () => {
    if (!data?.projectDistribution) return [];
    return data.projectDistribution.map((item) => ({
      name: getStatusLabel(item.status),
      value: item.count,
    }));
  };

  const getStatusLabel = (status: string): string => {
    const labels: Record<string, string> = {
      planning: '规划中',
      in_progress: '进行中',
      completed: '已完成',
      on_hold: '已暂停',
      cancelled: '已取消',
    };
    return labels[status] || status;
  };

  return (
    <div style={{ padding: '24px' }}>
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Title level={4}>数据分析仪表盘</Title>
          <Text type="secondary">查看项目和人员统计数据</Text>
        </Col>
        <Col>
          <Space>
            <RangePicker
              value={[
                dateRange.startDate ? dayjs(dateRange.startDate) : null,
                dateRange.endDate ? dayjs(dateRange.endDate) : null,
              ]}
              onChange={handleDateRangeChange}
            />
            <Button type="primary" onClick={handleGenerateReport}>
              生成报表
            </Button>
            <Button onClick={() => void fetchDashboardData()}>刷新</Button>
          </Space>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="项目总数"
            value={data?.overview.totalProjects || 0}
            prefix={<ProjectOutlined />}
            color="#1890ff"
            loading={loading}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="活跃用户"
            value={data?.overview.activeUsers || 0}
            prefix={<TeamOutlined />}
            color="#52c41a"
            loading={loading}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="已完成任务"
            value={data?.overview.completedTasks || 0}
            prefix={<CheckCircleOutlined />}
            color="#faad14"
            loading={loading}
          />
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <StatCard
            title="本周新增"
            value={data?.overview.weeklyNewItems || 0}
            prefix={<RiseOutlined />}
            color="#722ed1"
            loading={loading}
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} lg={12}>
          <PieChart
            title="项目状态分布"
            data={getProjectStatusData()}
            loading={loading}
          />
        </Col>
        <Col xs={24} lg={12}>
          <LineChart
            title="项目增长趋势"
            data={data?.projectTrends || []}
            loading={loading}
            yAxisName="项目数"
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <LineChart
            title="用户活跃度趋势"
            data={data?.userTrends || []}
            loading={loading}
            color="#52c41a"
            yAxisName="活跃用户数"
          />
        </Col>
        <Col xs={24} lg={12}>
          <HeatmapChart
            title="用户贡献热力图"
            data={data?.userHeatmap || []}
            loading={loading}
          />
        </Col>
      </Row>
    </div>
  );
};

export default AnalyticsDashboard;
