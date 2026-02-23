import React, { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import { Card, Empty, Spin } from 'antd';

interface PieChartProps {
  title?: string;
  data: { name: string; value: number }[];
  loading?: boolean;
  colors?: string[];
}

const DEFAULT_COLORS = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'];

const PieChart: React.FC<PieChartProps> = ({
  title = '项目状态分布',
  data,
  loading = false,
  colors = DEFAULT_COLORS,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current || loading) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const option: echarts.EChartsOption = {
      color: colors,
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        top: 'center',
      },
      series: [
        {
          name: '项目状态',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
            position: 'center',
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 18,
              fontWeight: 'bold',
            },
          },
          labelLine: {
            show: false,
          },
          data,
        },
      ],
    };

    chartInstance.current.setOption(option);

    const handleResize = () => {
      chartInstance.current?.resize();
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [data, loading, colors]);

  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  if (loading) {
    return (
      <Card title={title}>
        <div style={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Spin />
        </div>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card title={title}>
        <div style={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Empty description="暂无数据" />
        </div>
      </Card>
    );
  }

  return (
    <Card title={title}>
      <div ref={chartRef} style={{ height: 300 }} />
    </Card>
  );
};

export default PieChart;
