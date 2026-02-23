import React, { useState } from 'react';
import {
  Card,
  Table,
  Tag,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  Switch,
  Tabs,
  Badge,
  Tooltip,
  Popconfirm,
  message,
  Empty,
} from 'antd';
import {
  BellOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  AlertOutlined,
  WarningOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type {
  AlertRule,
  AlertHistory,
  CreateAlertRuleRequest,
  UpdateAlertRuleRequest,

} from '@/types/monitor';

const { TabPane } = Tabs;
const { TextArea } = Input;

interface AlertManagerProps {
  rules?: AlertRule[];
  history?: AlertHistory[];
  activeAlerts?: AlertHistory[];
  loading?: boolean;
  onCreateRule?: (data: CreateAlertRuleRequest) => Promise<void>;
  onUpdateRule?: (id: string, data: UpdateAlertRuleRequest) => Promise<void>;
  onDeleteRule?: (id: string) => Promise<void>;
  onToggleRule?: (id: string, isActive: boolean) => Promise<void>;
  onAcknowledgeAlert?: (id: string) => Promise<void>;
  onResolveAlert?: (id: string) => Promise<void>;
  onRefresh?: () => void;
}

const AlertManager: React.FC<AlertManagerProps> = ({
  rules = [],
  history = [],
  activeAlerts = [],
  loading = false,
  onCreateRule,
  onUpdateRule,
  onDeleteRule,
  onToggleRule,

  onResolveAlert,
  onRefresh,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<AlertRule | null>(null);
  const [form] = Form.useForm();
  const [activeTab, setActiveTab] = useState('active');

  // Available metrics for alerting
  const availableMetrics = [
    { value: 'cpu_usage', label: 'CPU 使用率 (%)' },
    { value: 'memory_usage', label: '内存使用率 (%)' },
    { value: 'disk_usage', label: '磁盘使用率 (%)' },
    { value: 'api_response_time', label: 'API 响应时间 (ms)' },
    { value: 'api_error_rate', label: 'API 错误率 (%)' },
    { value: 'db_connections', label: '数据库连接数' },
  ];

  // Condition options
  const conditionOptions = [
    { value: '>', label: '> (大于)' },
    { value: '<', label: '< (小于)' },
    { value: '==', label: '== (等于)' },
    { value: '!=', label: '!= (不等于)' },
  ];

  // Severity options
  const severityOptions = [
    { value: 'warning', label: '警告', color: 'orange' },
    { value: 'critical', label: '严重', color: 'red' },
  ];

  // Handle create/edit rule
  const handleSaveRule = async (values: any) => {
    try {
      const data: CreateAlertRuleRequest = {
        name: values.name,
        description: values.description,
        metric: values.metric,
        condition: values.condition,
        threshold: values.threshold,
        duration: values.duration,
        severity: values.severity,
        notify_channels: values.notify_channels || [],
      };

      if (editingRule) {
        await onUpdateRule?.(editingRule.id, data);
        message.success('告警规则已更新');
      } else {
        await onCreateRule?.(data);
        message.success('告警规则已创建');
      }

      setIsModalVisible(false);
      setEditingRule(null);
      form.resetFields();
      onRefresh?.();
    } catch (error) {
      message.error('操作失败');
    }
  };

  // Open edit modal
  const handleEdit = (rule: AlertRule) => {
    setEditingRule(rule);
    form.setFieldsValue({
      name: rule.name,
      description: rule.description,
      metric: rule.metric,
      condition: rule.condition,
      threshold: rule.threshold,
      duration: rule.duration,
      severity: rule.severity,
      notify_channels: rule.notify_channels,
    });
    setIsModalVisible(true);
  };

  // Open create modal
  const handleCreate = () => {
    setEditingRule(null);
    form.resetFields();
    setIsModalVisible(true);
  };

  // Get severity color
  const getSeverityColor = (severity: string): string => {
    return severity === 'critical' ? 'red' : 'orange';
  };

  // Get severity icon
  const getSeverityIcon = (severity: string) => {
    return severity === 'critical' ? <ExclamationCircleOutlined /> : <WarningOutlined />;
  };

  // Rules table columns
  const ruleColumns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: AlertRule) => (
        <div>
          <div style={{ fontWeight: 500 }}>{text}</div>
          {record.description && (
            <div style={{ fontSize: 12, color: '#999' }}>{record.description}</div>
          )}
        </div>
      ),
    },
    {
      title: '规则',
      key: 'rule',
      render: (_: any, record: AlertRule) => {
        const metric = availableMetrics.find((m) => m.value === record.metric);
        return (
          <span>
            {metric?.label || record.metric} {record.condition} {record.threshold}
          </span>
        );
      },
    },
    {
      title: '持续时间',
      dataIndex: 'duration',
      key: 'duration',
      render: (value: number) => `${value} 分钟`,
    },
    {
      title: '级别',
      dataIndex: 'severity',
      key: 'severity',
      render: (value: string) => (
        <Tag color={getSeverityColor(value)} icon={getSeverityIcon(value)}>
          {value === 'critical' ? '严重' : '警告'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (value: boolean, record: AlertRule) => (
        <Switch
          checked={value}
          onChange={(checked) => onToggleRule?.(record.id, checked)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: AlertRule) => (
        <Space>
          <Tooltip title="编辑">
            <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(record)} />
          </Tooltip>
          <Popconfirm
            title="确认删除"
            description="确定要删除此告警规则吗？"
            onConfirm={() => onDeleteRule?.(record.id)}
            okText="删除"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button icon={<DeleteOutlined />} size="small" danger />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // History table columns
  const historyColumns = [
    {
      title: '告警',
      key: 'alert',
      render: (_: any, record: AlertHistory) => (
        <div>
          <Space>
            <Tag color={getSeverityColor(record.severity)}>{record.rule_name}</Tag>
            <span>{record.message}</span>
          </Space>
          <div style={{ fontSize: 12, color: '#999', marginTop: 4 }}>
            {new Date(record.created_at).toLocaleString()}
          </div>
        </div>
      ),
    },
    {
      title: '值 / 阈值',
      key: 'value',
      render: (_: any, record: AlertHistory) => (
        <span>
          {record.value.toFixed(2)} / {record.threshold.toFixed(2)}
        </span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (value: string) => (
        <Tag color={value === 'resolved' ? 'green' : 'red'}>
          {value === 'resolved' ? '已解决' : '告警中'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: AlertHistory) => (
        <Space>
          {record.status === 'firing' && (
            <>
              <Button
                size="small"
                icon={<CheckCircleOutlined />}
                onClick={() => onResolveAlert?.(record.id)}
              >
                解决
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ];

  // Active alerts columns
  const activeAlertColumns = [
    {
      title: '告警',
      key: 'alert',
      render: (_: any, record: AlertHistory) => (
        <div>
          <Space>
            <Badge status="error" />
            <Tag color={getSeverityColor(record.severity)}>{record.rule_name}</Tag>
          </Space>
          <div style={{ marginTop: 4 }}>{record.message}</div>
          <div style={{ fontSize: 12, color: '#999', marginTop: 4 }}>
            {new Date(record.created_at).toLocaleString()}
          </div>
        </div>
      ),
    },
    {
      title: '当前值',
      dataIndex: 'value',
      key: 'value',
      render: (value: number) => <span style={{ color: '#ff4d4f', fontWeight: 500 }}>{value.toFixed(2)}</span>,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: AlertHistory) => (
        <Button
          type="primary"
          size="small"
          icon={<CheckCircleOutlined />}
          onClick={() => onResolveAlert?.(record.id)}
        >
          确认解决
        </Button>
      ),
    },
  ];

  const activeCount = activeAlerts.length;
  const warningCount = activeAlerts.filter((a) => a.severity === 'warning').length;
  const criticalCount = activeAlerts.filter((a) => a.severity === 'critical').length;

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <BellOutlined />
          <span>告警管理</span>
          {activeCount > 0 && (
            <Badge count={activeCount} style={{ backgroundColor: '#ff4d4f' }} />
          )}
        </div>
      }
      extra={
        <Space>
          <Button icon={<ReloadOutlined />} size="small" onClick={onRefresh}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} size="small" onClick={handleCreate}>
            新建规则
          </Button>
        </Space>
      }
    >
      {/* Alert Summary */}
      {activeCount > 0 && (
        <div style={{ marginBottom: 16, padding: 12, backgroundColor: '#fff2f0', borderRadius: 6 }}>
          <Space size="large">
            <span>
              <Badge status="error" /> 活跃告警: <strong>{activeCount}</strong>
            </span>
            {criticalCount > 0 && (
              <span style={{ color: '#ff4d4f' }}>
                <ExclamationCircleOutlined /> 严重: {criticalCount}
              </span>
            )}
            {warningCount > 0 && (
              <span style={{ color: '#faad14' }}>
                <WarningOutlined /> 警告: {warningCount}
              </span>
            )}
          </Space>
        </div>
      )}

      <Tabs activeKey={activeTab} onChange={setActiveTab}>
        <TabPane
          tab={
            <span>
              <AlertOutlined /> 活跃告警
              {activeCount > 0 && <Badge count={activeCount} style={{ marginLeft: 4 }} />}
            </span>
          }
          key="active"
        >
          <Table
            columns={activeAlertColumns}
            dataSource={activeAlerts}
            rowKey="id"
            loading={loading}
            pagination={false}
            locale={{ emptyText: <Empty description="暂无活跃告警" /> }}
          />
        </TabPane>

        <TabPane
          tab={
            <span>
              <BellOutlined /> 告警规则 ({rules.length})
            </span>
          }
          key="rules"
        >
          <Table
            columns={ruleColumns}
            dataSource={rules}
            rowKey="id"
            loading={loading}
            pagination={false}
            scroll={{ x: 'max-content' }}
          />
        </TabPane>

        <TabPane
          tab={
            <span>
              <ExclamationCircleOutlined /> 历史记录 ({history.length})
            </span>
          }
          key="history"
        >
          <Table
            columns={historyColumns}
            dataSource={history}
            rowKey="id"
            loading={loading}
            pagination={{ pageSize: 10 }}
          />
        </TabPane>
      </Tabs>

      {/* Create/Edit Modal */}
      <Modal
        title={editingRule ? '编辑告警规则' : '新建告警规则'}
        open={isModalVisible}
        onOk={() => form.submit()}
        onCancel={() => {
          setIsModalVisible(false);
          setEditingRule(null);
          form.resetFields();
        }}
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSaveRule}
          initialValues={{
            condition: '>',
            duration: 5,
            severity: 'warning',
            is_active: true,
          }}
        >
          <Form.Item
            name="name"
            label="规则名称"
            rules={[{ required: true, message: '请输入规则名称' }]}
          >
            <Input placeholder="例如：CPU 使用率过高" />
          </Form.Item>

          <Form.Item name="description" label="描述">
            <TextArea rows={2} placeholder="告警规则描述（可选）" />
          </Form.Item>

          <Form.Item
            name="metric"
            label="监控指标"
            rules={[{ required: true, message: '请选择监控指标' }]}
          >
            <Select placeholder="选择要监控的指标">
              {availableMetrics.map((metric) => (
                <Select.Option key={metric.value} value={metric.value}>{metric.label}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Space style={{ display: 'flex' }} align="baseline">
            <Form.Item
              name="condition"
              label="条件"
              rules={[{ required: true }]}
              style={{ width: 120 }}
            >
              <Select>
                {conditionOptions.map((opt) => (
                  <Select.Option key={opt.value} value={opt.value}>{opt.label}</Select.Option>
                ))}
              </Select>
            </Form.Item>

            <Form.Item
              name="threshold"
              label="阈值"
              rules={[{ required: true, message: '请输入阈值' }]}
              style={{ width: 150 }}
            >
              <InputNumber style={{ width: '100%' }} placeholder="例如：80" />
            </Form.Item>

            <Form.Item
              name="duration"
              label="持续时间（分钟）"
              rules={[{ required: true }]}
              style={{ width: 150 }}
            >
              <InputNumber min={1} max={60} style={{ width: '100%' }} />
            </Form.Item>
          </Space>

          <Form.Item
            name="severity"
            label="告警级别"
            rules={[{ required: true }]}
          >
            <Select>
              {severityOptions.map((opt) => (
                <Select.Option key={opt.value} value={opt.value}>
                  <Tag color={opt.color}>{opt.label}</Tag>
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="notify_channels" label="通知渠道">
            <Select mode="multiple" placeholder="选择通知渠道">
              <Select.Option value="notification">站内通知</Select.Option>
              <Select.Option value="email">邮件</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default AlertManager;
