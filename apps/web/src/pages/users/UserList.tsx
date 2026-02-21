import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table,
  Button,
  Space,
  Tag,
  Input,
  Select,
  Card,
  Modal,
  Avatar,
  Typography,
  message,
  Tooltip,
  Row,
  Col,
} from 'antd';
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  UserOutlined,
  EyeOutlined,
  TeamOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import { userAPI } from '@/services/user';
import type { User } from '@/types';

const { Title, Text } = Typography;
const { Option } = Select;
const { Search } = Input;

const roleOptions = [
  { value: 'admin', label: '管理员', color: 'red' },
  { value: 'manager', label: '项目经理', color: 'blue' },
  { value: 'leader', label: '负责人', color: 'green' },
  { value: 'designer', label: '工程师', color: 'default' },
  { value: 'viewer', label: '访客', color: 'default' },
];

const statusOptions = [
  { value: 'active', label: '正常', color: 'success' },
  { value: 'inactive', label: '停用', color: 'default' },
  { value: 'locked', label: '锁定', color: 'error' },
  { value: 'pending', label: '待审核', color: 'warning' },
];

export default function UserList() {
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });
  const [filters, setFilters] = useState({
    search: '',
    role: '',
    status: '',
  });
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 权限检查
  const isAdmin = currentUser?.role === 'admin' || currentUser?.roles?.some((r: Role) => r.code === 'admin');

  useEffect(() => {
    fetchUsers();
  }, [pagination.current, pagination.pageSize]);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const params: Record<string, unknown> = {
        page: pagination.current,
        page_size: pagination.pageSize,
      };
      
      if (filters.search) params.search = filters.search;
      if (filters.role) params.role = filters.role;
      if (filters.status) params.status = filters.status;

      const response = await userAPI.getUsers(params);
      setUsers(response.data || []);
      setPagination(prev => ({ ...prev, total: response.total || 0 }));
    } catch (error: any) {
      message.error(error.message || '获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setPagination(prev => ({ ...prev, current: 1 }));
    fetchUsers();
  };

  const handleReset = () => {
    setFilters({ search: '', role: '', status: '' });
    setPagination(prev => ({ ...prev, current: 1 }));
    fetchUsers();
  };

  const handleDelete = async (id: string) => {
    if (!isAdmin) {
      message.error('您没有权限删除用户');
      return;
    }

    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此用户吗？此操作不可恢复。',
      okText: '确认删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await userAPI.deleteUser(id);
          message.success('用户已删除');
          fetchUsers();
        } catch (error: any) {
          message.error(error.message || '删除用户失败');
        }
      },
    });
  };

  const handleBatchDelete = () => {
    if (!isAdmin) {
      message.error('您没有权限删除用户');
      return;
    }

    if (selectedRowKeys.length === 0) {
      message.warning('请先选择要删除的用户');
      return;
    }

    Modal.confirm({
      title: '批量删除确认',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个用户吗？`,
      okText: '确认删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await Promise.all(selectedRowKeys.map(id => userAPI.deleteUser(id as string)));
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          fetchUsers();
        } catch (error: any) {
          message.error(error.message || '批量删除失败');
        }
      },
    });
  };

  const getRoleTag = (role: string) => {
    const config = roleOptions.find(r => r.value === role);
    return <Tag color={config?.color || 'default'}>{config?.label || role}</Tag>;
  };

  const getStatusTag = (status: string) => {
    const config = statusOptions.find(s => s.value === status);
    return <Tag color={config?.color || 'default'}>{config?.label || status}</Tag>;
  };

  const columns = [
    {
      title: '用户信息',
      key: 'user',
      render: (_: unknown, record: User) => (
        <Space>
          <Avatar 
            src={record.avatar} 
            icon={!record.avatar && <UserOutlined />}
            size="large"
          />
          <div>
            <div style={{ fontWeight: 500 }}>{record.displayName}</div>
            <Text type="secondary" style={{ fontSize: 12 }}>
              @{record.username}
            </Text>
          </div>
        </Space>
      ),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 220,
      render: (email: string) => email || '-',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      width: 120,
      render: (role: string, record: User) => {
        const roleCode = record.roles?.[0]?.code || role;
        return getRoleTag(roleCode);
      },
    },
    {
      title: '部门',
      dataIndex: 'organization',
      key: 'organization',
      width: 150,
      render: (org: { name: string }) => org?.name || '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (date: string) => new Date(date).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right' as const,
      render: (_: unknown, record: User) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="text"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => navigate(`/users/${record.id}`)}
            />
          </Tooltip>
          {isAdmin && (
            <Tooltip title="编辑">
              <Button
                type="text"
                size="small"
                icon={<EditOutlined />}
                onClick={() => navigate(`/users/${record.id}/edit`)}
              />
            </Tooltip>
          )}
          {isAdmin && currentUser?.id !== record.id && (
            <Tooltip title="删除">
              <Button
                type="text"
                size="small"
                danger
                icon={<DeleteOutlined />}
                onClick={() => handleDelete(record.id)}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  const rowSelection = isAdmin
    ? {
        selectedRowKeys,
        onChange: (keys: React.Key[]) => setSelectedRowKeys(keys),
      }
    : undefined;

  return (
    <div>
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Title level={2} style={{ margin: 0 }}>
            <TeamOutlined style={{ marginRight: 12 }} />
            用户管理
          </Title>
          <Text type="secondary">管理系统用户、角色和权限</Text>
        </Col>
        <Col>
          {isAdmin && (
            <Button
              type="primary"
              icon={<PlusOutlined />}
              size="large"
              onClick={() => navigate('/register')}
            >
              添加用户
            </Button>
          )}
        </Col>
      </Row>

      <Card style={{ marginBottom: 24 }}>
        <Row gutter={16} align="middle">
          <Col flex="auto">
            <Space wrap>
              <Search
                placeholder="搜索用户名、姓名或邮箱"
                value={filters.search}
                onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
                onSearch={handleSearch}
                style={{ width: 250 }}
                allowClear
              />
              <Select
                placeholder="选择角色"
                style={{ width: 150 }}
                allowClear
                value={filters.role || undefined}
                onChange={(value) => setFilters(prev => ({ ...prev, role: value || '' }))}
              >
                {roleOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
              <Select
                placeholder="选择状态"
                style={{ width: 150 }}
                allowClear
                value={filters.status || undefined}
                onChange={(value) => setFilters(prev => ({ ...prev, status: value || '' }))}
              >
                {statusOptions.map(opt => (
                  <Option key={opt.value} value={opt.value}>{opt.label}</Option>
                ))}
              </Select>
              <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                查询
              </Button>
              <Button icon={<ReloadOutlined />} onClick={handleReset}>
                重置
              </Button>
            </Space>
          </Col>
          <Col>
            {isAdmin && selectedRowKeys.length > 0 && (
              <Button danger onClick={handleBatchDelete}>
                批量删除 ({selectedRowKeys.length})
              </Button>
            )}
          </Col>
        </Row>
      </Card>

      <Card>
        <Table
          columns={columns}
          dataSource={users}
          loading={loading}
          rowKey="id"
          rowSelection={rowSelection}
          scroll={{ x: 1200 }}
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
            onChange: (page, pageSize) => {
              setPagination(prev => ({ ...prev, current: page, pageSize: pageSize || 20 }));
            },
          }}
        />
      </Card>
    </div>
  );
}
