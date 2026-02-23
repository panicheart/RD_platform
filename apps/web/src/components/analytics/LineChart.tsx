import React, { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import { Card, Empty, Spin } from 'antd';

interface LineChartProps {
  title?: string;
  data: { date: string; value: number; label?: string }[];
  loading?: boolean;
  color?: string;
  yAxisName?: string;
  showArea?: boolean;
}

const LineChart: React.FC<LineChartProps> = ({
  title = '趋势分析',
  data,
  loading = false,
  color = '#1890ff',
  yAxisName,
  showArea = true,
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current || loading) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const dates = data.map((item) => item.date);
    const values = data.map((item) => item.value);

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: {
            backgroundColor: '#6a7985',
          },
        },
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: [
        {
          type: 'category',
          boundaryGap: false,
          data: dates,
          axisLine: {
            lineStyle: {
              color: '#d9d9d9',
            },
          },
          axisLabel: {
            color: '#666',
          },
        },
      ],
      yAxis: [
        {
          type: 'value',
          name: yAxisName,
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              color: '#f0f0f0',
            },
          },
          axisLabel: {
            color: '#666',
          },
        },
      ],
      series: [
        {
          name: '数值',
          type: 'line',
          stack: 'Total',
          smooth: true,
          lineStyle: {
            width: 3,
            color,
          },
          showSymbol: false,
          areaStyle: showArea
            ? {
                opacity: 0.3,
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  {
                    offset: 0,
                    color: color,
                  },
                  {
                    offset: 1,
                    color: 'rgba(255, 255, 255, 0)',
                  },
                ]),
              }
            : undefined,
          emphasis: {
            focus: 'series',
          },
          data: values,
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
  }, [data, loading, color, yAxisName, showArea]);

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

export default LineChart;
