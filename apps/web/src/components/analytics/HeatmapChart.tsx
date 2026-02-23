import React, { useEffect, useRef } from 'react';
import * as echarts from 'echarts';
import { Card, Empty, Spin } from 'antd';

interface HeatmapChartProps {
  title?: string;
  data: { date: string; count: number }[];
  loading?: boolean;
  year?: number;
}

const HeatmapChart: React.FC<HeatmapChartProps> = ({
  title = '贡献热力图',
  data,
  loading = false,
  year = new Date().getFullYear(),
}) => {
  const chartRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  useEffect(() => {
    if (!chartRef.current || loading) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(chartRef.current);
    }

    const getVirtualData = () => {
      const dateList: [string, number][] = [];
      const startDate = new Date(`${year}-01-01`);
      const endDate = new Date(`${year}-12-31`);

      for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
        const dateStr = d.toISOString().split('T')[0] || '';
        const item = data.find((item) => item.date === dateStr);
        dateList.push([dateStr, item?.count || 0]);
      }

      return dateList;
    };

    const chartData = getVirtualData();
    const maxValue = Math.max(...chartData.map((item) => item[1]), 1);

    const option: echarts.EChartsOption = {
      tooltip: {
        position: 'top',
        formatter: (params: unknown) => {
          const p = params as { data: [string, number] };
          return `${p.data[0]}: ${p.data[1]} 次贡献`;
        },
      },
      visualMap: {
        min: 0,
        max: maxValue,
        calculable: false,
        orient: 'horizontal',
        left: 'center',
        bottom: 0,
        inRange: {
          color: ['#ebedf0', '#c6e48b', '#7bc96f', '#239a3b', '#196127'],
        },
        show: false,
      },
      calendar: {
        top: 30,
        left: 30,
        right: 30,
        cellSize: ['auto', 13],
        range: `${year}`,
        itemStyle: {
          borderWidth: 2,
          borderColor: '#fff',
          borderRadius: 2,
        },
        splitLine: {
          show: false,
        },
        dayLabel: {
          nameMap: 'cn',
          color: '#666',
        },
        monthLabel: {
          nameMap: 'cn',
          color: '#666',
        },
        yearLabel: {
          show: false,
        },
      },
      series: [
        {
          type: 'heatmap',
          coordinateSystem: 'calendar',
          data: chartData,
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
  }, [data, loading, year]);

  useEffect(() => {
    return () => {
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  if (loading) {
    return (
      <Card title={title}>
        <div style={{ height: 200, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Spin />
        </div>
      </Card>
    );
  }

  if (!data || data.length === 0) {
    return (
      <Card title={title}>
        <div style={{ height: 200, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Empty description="暂无数据" />
        </div>
      </Card>
    );
  }

  return (
    <Card title={title}>
      <div ref={chartRef} style={{ height: 200 }} />
    </Card>
  );
};

export default HeatmapChart;
