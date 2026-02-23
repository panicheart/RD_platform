import React, { useEffect, useState } from 'react';
import { Card, Statistic } from 'antd';
import type { ReactNode } from 'react';

interface StatCardProps {
  title: string;
  value: number;
  prefix?: ReactNode;
  suffix?: string;
  precision?: number;
  valueStyle?: React.CSSProperties;
  loading?: boolean;
  color?: string;
}

const StatCard: React.FC<StatCardProps> = ({
  title,
  value,
  prefix,
  suffix,
  precision = 0,
  valueStyle,
  loading = false,
  color = '#1890ff',
}) => {
  const [displayValue, setDisplayValue] = useState(0);

  useEffect(() => {
    if (loading) return;

    const duration = 1000;
    const startTime = Date.now();
    const startValue = 0;
    const endValue = value;

    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const currentValue = startValue + (endValue - startValue) * easeOutQuart;

      setDisplayValue(currentValue);

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };

    requestAnimationFrame(animate);
  }, [value, loading]);

  return (
    <Card loading={loading} bodyStyle={{ padding: '20px 24px' }}>
      <Statistic
        title={title}
        value={displayValue}
        precision={precision}
        prefix={prefix}
        suffix={suffix}
        valueStyle={{
          color,
          fontSize: '28px',
          fontWeight: 600,
          ...valueStyle,
        }}
      />
    </Card>
  );
};

export default StatCard;
