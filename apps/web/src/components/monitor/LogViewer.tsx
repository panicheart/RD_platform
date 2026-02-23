import React, { useEffect, useRef, useState, useCallback } from 'react';
import {
  Card,
  Input,
  Select,
  Button,
  Tag,
  List,
  Empty,
  Spin,
  Space,
  Typography,
  Badge,
  Tooltip,
} from 'antd';
import {
  SearchOutlined,
  ReloadOutlined,
  CodeOutlined,
  InfoCircleOutlined,
  WarningOutlined,
  CloseCircleOutlined,
  BugOutlined,
  PauseCircleOutlined,
  PlayCircleOutlined,
} from '@ant-design/icons';
import type { LogEntry, LogFilterParams } from '@/types/monitor';


const { Text } = Typography;

interface LogViewerProps {
  logs?: LogEntry[];
  loading?: boolean;
  total?: number;
  sources?: string[];
  modules?: string[];
  onFilterChange?: (filters: LogFilterParams) => void;
  onRefresh?: () => void;
  onLoadMore?: () => void;
  autoRefresh?: boolean;
  onAutoRefreshChange?: (enabled: boolean) => void;
}

const LogViewer: React.FC<LogViewerProps> = ({
  logs = [],
  loading = false,
  total = 0,
  sources = [],
  modules = [],
  onFilterChange,
  onRefresh,
  onLoadMore,
  autoRefresh = false,
  onAutoRefreshChange,
}) => {
  const [filters, setFilters] = useState<LogFilterParams>({});
  const [searchKeyword, setSearchKeyword] = useState('');
  const listRef = useRef<HTMLDivElement>(null);
  const [isAutoScroll, setIsAutoScroll] = useState(true);

  // Get log level color
  const getLevelColor = (level: string): string => {
    switch (level) {
      case 'ERROR':
        return 'red';
      case 'WARN':
        return 'orange';
      case 'INFO':
        return 'blue';
      case 'DEBUG':
        return 'default';
      default:
        return 'default';
    }
  };

  // Get log level icon
  const getLevelIcon = (level: string) => {
    switch (level) {
      case 'ERROR':
        return <CloseCircleOutlined />;
      case 'WARN':
        return <WarningOutlined />;
      case 'INFO':
        return <InfoCircleOutlined />;
      case 'DEBUG':
        return <BugOutlined />;
      default:
        return <CodeOutlined />;
    }
  };

  // Handle filter change
  const handleFilterChange = useCallback(
    (newFilters: Partial<LogFilterParams>) => {
      const updated = { ...filters, ...newFilters, page: 1 };
      setFilters(updated);
      onFilterChange?.(updated);
    },
    [filters, onFilterChange]
  );

  // Handle search
  const handleSearch = () => {
    handleFilterChange({ keyword: searchKeyword || undefined });
  };

  // Format timestamp
  const formatTimestamp = (timestamp: string): string => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  // Auto scroll to bottom
  useEffect(() => {
    if (isAutoScroll && listRef.current && autoRefresh) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [logs, isAutoScroll, autoRefresh]);

  // Render log item
  const renderLogItem = (log: LogEntry) => {
    const levelColor = getLevelColor(log.level);

    return (
      <List.Item
        key={log.id}
        style={{
          padding: '8px 16px',
          borderBottom: '1px solid #f0f0f0',
          fontFamily: 'monospace',
          fontSize: '13px',
        }}
      >
        <div style={{ width: '100%' }}>
          <Space wrap style={{ marginBottom: 4 }}>
            <Text type="secondary">{formatTimestamp(log.timestamp)}</Text>
            <Tag color={levelColor} icon={getLevelIcon(log.level)}>{log.level}</Tag>
            <Tag>{log.source}</Tag>
            {log.module && <Tag>{log.module}</Tag>}
            {log.request_id && (
              <Tooltip title={log.request_id}>
                <Tag color="purple">请求ID</Tag>
              </Tooltip>
            )}
          </Space>
          <div style={{ marginTop: 4, wordBreak: 'break-all' }}>{log.message}</div>
          {log.metadata && (
            <Text type="secondary" style={{ fontSize: '11px', marginTop: 4, display: 'block' }}>
              {log.metadata}
            </Text>
          )}
        </div>
      </List.Item>
    );
  };

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <CodeOutlined />
          <span>系统日志</span>
          {autoRefresh && <Badge status="processing" text="实时" />}
          <Text type="secondary" style={{ marginLeft: 'auto', fontSize: 12 }}>
            共 {total.toLocaleString()} 条
          </Text>
        </div>
      }
      extra={
        <Space>
          <Tooltip title={isAutoScroll ? '自动滚动开启' : '自动滚动关闭'}>
            <Button
              icon={isAutoScroll ? <PauseCircleOutlined /> : <PlayCircleOutlined />}
              size="small"
              onClick={() => setIsAutoScroll(!isAutoScroll)}
            />
          </Tooltip>
          <Tooltip title={autoRefresh ? '暂停实时更新' : '开启实时更新'}>
            <Button
              icon={autoRefresh ? <PauseCircleOutlined /> : <PlayCircleOutlined />}
              size="small"
              type={autoRefresh ? 'primary' : 'default'}
              onClick={() => onAutoRefreshChange?.(!autoRefresh)}
            >
              {autoRefresh ? '实时' : '暂停'}
            </Button>
          </Tooltip>
          <Button icon={<ReloadOutlined />} size="small" onClick={onRefresh}>刷新</Button>
        </Space>
      }
    >
      {/* Filter Bar */}
      <Space wrap style={{ marginBottom: 16 }}>
        <Select
          placeholder="日志级别"
          allowClear
          style={{ width: 120 }}
          value={filters.level}
          onChange={(value) => handleFilterChange({ level: value })}
        >
          <Select.Option value="DEBUG">DEBUG</Select.Option>
          <Select.Option value="INFO">INFO</Select.Option>
          <Select.Option value="WARN">WARN</Select.Option>
          <Select.Option value="ERROR">ERROR</Select.Option>
        </Select>

        <Select
          placeholder="服务来源"
          allowClear
          showSearch
          style={{ width: 150 }}
          value={filters.source}
          onChange={(value) => handleFilterChange({ source: value })}
        >
          {sources.map((source) => (
            <Select.Option key={source} value={source}>{source}</Select.Option>
          ))}
        </Select>

        <Select
          placeholder="模块"
          allowClear
          showSearch
          style={{ width: 150 }}
          value={filters.module}
          onChange={(value) => handleFilterChange({ module: value })}
        >
          {modules.map((module) => (
            <Select.Option key={module} value={module}>{module}</Select.Option>
          ))}
        </Select>

        <Input.Search
          placeholder="搜索关键词"
          style={{ width: 200 }}
          value={searchKeyword}
          onChange={(e) => setSearchKeyword(e.target.value)}
          onSearch={handleSearch}
          allowClear
          enterButton={<SearchOutlined />}
        />

        <Button onClick={() => {
          setFilters({});
          setSearchKeyword('');
          onFilterChange?.({});
        }}>清空</Button>
      </Space>

      {/* Log List */}
      <div
        ref={listRef}
        style={{
          height: 400,
          overflow: 'auto',
          backgroundColor: '#fafafa',
          border: '1px solid #f0f0f0',
          borderRadius: 6,
        }}
      >
        {loading && logs.length === 0 ? (
          <div style={{ height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <Spin />
          </div>
        ) : logs.length === 0 ? (
          <div style={{ height: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <Empty description="暂无日志" />
          </div>
        ) : (
          <List
            dataSource={logs}
            renderItem={renderLogItem}
            locale={{ emptyText: '暂无日志' }}
          />
        )}

        {/* Load More */}
        {logs.length > 0 && logs.length < total && (
          <div style={{ textAlign: 'center', padding: '12px 0' }}>
            <Button onClick={onLoadMore} loading={loading}>加载更多</Button>
          </div>
        )}
      </div>
    </Card>
  );
};

export default LogViewer;
