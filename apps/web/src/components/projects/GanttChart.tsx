import React, { useState, useEffect } from 'react';
import { Gantt, Task, ViewMode } from 'gantt-task-react';
import 'gantt-task-react/dist/index.css';
import { Card, Radio, Space, Spin, Empty } from 'antd';
import { useActivities } from '@/hooks/useActivities';
import dayjs from 'dayjs';

interface GanttChartProps {
  projectId: string;
  workflowId: string;
}

interface GanttTask extends Task {
  activityId: string;
}

const GanttChart: React.FC<GanttChartProps> = ({ projectId, workflowId }) => {
  const [viewMode, setViewMode] = useState<ViewMode>(ViewMode.Day);
  const { activities, loading } = useActivities(workflowId);

  const convertToGanttTasks = (activities: any[]): GanttTask[] => {
    return activities.map((activity, index) => {
      const start = activity.actual_start 
        ? dayjs(activity.actual_start).toDate()
        : dayjs(activity.planned_start || new Date()).toDate();
      
      const end = activity.actual_end
        ? dayjs(activity.actual_end).toDate()
        : dayjs(activity.planned_end || dayjs().add(7, 'day')).toDate();

      let type: Task['type'] = 'task';
      if (activity.type === 'milestone') type = 'milestone';
      else if (activity.progress === 100) type = 'done';

      return {
        id: activity.id,
        name: activity.name,
        start,
        end,
        progress: activity.progress || 0,
        type,
        activityId: activity.id,
        dependencies: activity.dependencies?.map((d: any) => d.depends_on_id) || [],
        styles: {
          backgroundColor: getStatusColor(activity.status),
          progressColor: '#1890ff',
        },
      };
    });
  };

  const getStatusColor = (status: string): string => {
    const colors: Record<string, string> = {
      pending: '#d9d9d9',
      ready: '#52c41a',
      running: '#1890ff',
      completed: '#52c41a',
      reviewing: '#faad14',
      approved: '#52c41a',
      rejected: '#ff4d4f',
      blocked: '#ff4d4f',
    };
    return colors[status] || '#d9d9d9';
  };

  const handleTaskChange = (task: Task) => {
    // TODO: Update activity dates via API
    console.log('Task changed:', task);
  };

  const handleTaskSelect = (task: Task) => {
    // TODO: Show activity details
    console.log('Task selected:', task);
  };

  if (loading) {
    return (
      <Card title="项目甘特图">
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Spin size="large" />
        </div>
      </Card>
    );
  }

  if (!activities || activities.length === 0) {
    return (
      <Card title="项目甘特图">
        <Empty description="暂无活动数据" />
      </Card>
    );
  }

  const tasks = convertToGanttTasks(activities);

  return (
    <Card
      title="项目甘特图"
      extra={
        <Space>
          <Radio.Group 
            value={viewMode} 
            onChange={(e) => setViewMode(e.target.value)}
            size="small"
          >
            <Radio.Button value={ViewMode.Day}>日</Radio.Button>
            <Radio.Button value={ViewMode.Week}>周</Radio.Button>
            <Radio.Button value={ViewMode.Month}>月</Radio.Button>
          </Radio.Group>
        </Space>
      }
    >
      <div style={{ height: '500px', overflow: 'auto' }}>
        <Gantt
          tasks={tasks}
          viewMode={viewMode}
          onDateChange={handleTaskChange}
          onSelect={handleTaskSelect}
          columnWidth={viewMode === ViewMode.Month ? 150 : 65}
          barBackgroundColor="#1890ff"
          barProgressColor="#52c41a"
          barBackgroundSelectedColor="#096dd9"
          barProgressSelectedColor="#389e0d"
        />
      </div>
    </Card>
  );
};

export default GanttChart;
