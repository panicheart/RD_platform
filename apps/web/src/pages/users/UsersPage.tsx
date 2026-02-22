import { useEffect, useState, useCallback } from 'react'
import { Table, Button, Space, Tag, Input, Select, Card, Modal, Form, Avatar, Typography, message } from 'antd'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons'
import { api, ApiResponse, PaginatedResponse, getErrorMessage } from '@/utils/api'

const { Title, Text } = Typography
const { Option } = Select

interface User {
  id: string
  username: string
  display_name: string
  email?: string
  phone?: string
  avatar_url?: string
  role: string
  team?: string
  specialty?: string
  title?: string
  is_active: boolean
  created_at: string
}

const roleOptions = [
  { value: 'admin', label: '管理员' },
  { value: 'manager', label: '项目经理' },
  { value: 'leader', label: '技术/项目负责人' },
  { value: 'designer', label: '设计师/工程师' },
]

const teamOptions = [
  { value: '微波天线', label: '微波天线' },
  { value: '射频前端', label: '射频前端' },
  { value: '数字电路', label: '数字电路' },
  { value: '结构工艺', label: '结构工艺' },
  { value: '软件', label: '软件' },
  { value: '测试', label: '测试' },
]

interface CreateUserValues {
  username: string
  display_name: string
  email?: string
  phone?: string
  role: string
  team?: string
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ page: 1, pageSize: 20, total: 0 })
  const [filters, setFilters] = useState({ search: '', role: '', team: '' })
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchUsers = useCallback(async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      params.append('page', pagination.page.toString())
      params.append('page_size', pagination.pageSize.toString())
      if (filters.search) params.append('search', filters.search)
      if (filters.role) params.append('role', filters.role)
      if (filters.team) params.append('team', filters.team)

      const response = await api.get<ApiResponse<PaginatedResponse<User>>>(`/users?${params.toString()}`)
      setUsers(response.data.data?.items || [])
      setPagination(prev => ({ ...prev, total: response.data.data?.total || 0 }))
    } catch (error) {
      message.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }, [pagination.page, pagination.pageSize, filters])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  const handleCreate = async (values: CreateUserValues) => {
    try {
      await api.post('/users', values)
      message.success('用户创建成功')
      setCreateModalOpen(false)
      form.resetFields()
      fetchUsers()
    } catch (error) {
      message.error(getErrorMessage(error))
    }
  }

  const handleDelete = async (id: string) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此用户吗？',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await api.delete(`/users/${id}`)
          message.success('用户已删除')
          fetchUsers()
        } catch (error) {
          message.error(getErrorMessage(error))
        }
      },
    })
  }

  const getRoleTag = (role: string) => {
    const roleMap: Record<string, { color: string; text: string }> = {
      admin: { color: 'red', text: '管理员' },
      manager: { color: 'blue', text: '项目经理' },
      leader: { color: 'green', text: '负责人' },
      designer: { color: 'default', text: '工程师' },
    }
    const config = roleMap[role] || { color: 'default', text: role }
    return <Tag color={config.color}>{config.text}</Tag>
  }

  const columns = [
    {
      title: '用户',
      key: 'user',
      render: (_: unknown, record: User) => (
        <Space>
          <Avatar src={record.avatar_url} icon={<UserOutlined />} />
          <div>
            <div>{record.display_name}</div>
            <Text type="secondary" style={{ fontSize: 12 }}>{record.username}</Text>
          </div>
        </Space>
      ),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      render: (email: string) => email || '-',
    },
    {
      title: '电话',
      dataIndex: 'phone',
      key: 'phone',
      render: (phone: string) => phone || '-',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      width: 100,
      render: (role: string) => getRoleTag(role),
    },
    {
      title: '团队',
      dataIndex: 'team',
      key: 'team',
      width: 100,
      render: (team: string) => team || '-',
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      width: 80,
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'default'}>
          {isActive ? '正常' : '停用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_: unknown, record: User) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />}>
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <Space style={{ width: '100%', justifyContent: 'space-between' }}>
          <div>
            <Title level={2} style={{ marginBottom: 0 }}>用户管理</Title>
            <Text type="secondary">管理系统用户和权限</Text>
          </div>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateModalOpen(true)}
          >
            添加用户
          </Button>
        </Space>
      </div>

      <Card>
        <Space style={{ marginBottom: 16 }} wrap>
          <Input
            placeholder="搜索用户名或姓名"
            prefix={<SearchOutlined />}
            style={{ width: 200 }}
            value={filters.search}
            onChange={e => setFilters(prev => ({ ...prev, search: e.target.value }))}
            onPressEnter={fetchUsers}
          />
          <Select
            placeholder="用户角色"
            style={{ width: 120 }}
            allowClear
            value={filters.role || undefined}
            onChange={value => setFilters(prev => ({ ...prev, role: value || '' }))}
          >
            {roleOptions.map(opt => (
              <Option key={opt.value} value={opt.value}>{opt.label}</Option>
            ))}
          </Select>
          <Select
            placeholder="所属团队"
            style={{ width: 120 }}
            allowClear
            value={filters.team || undefined}
            onChange={value => setFilters(prev => ({ ...prev, team: value || '' }))}
          >
            {teamOptions.map(opt => (
              <Option key={opt.value} value={opt.value}>{opt.label}</Option>
            ))}
          </Select>
          <Button type="primary" icon={<SearchOutlined />} onClick={fetchUsers}>
            查询
          </Button>
        </Space>

        <Table
          columns={columns}
          dataSource={users}
          loading={loading}
          rowKey="id"
          pagination={{
            current: pagination.page,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: (page) => setPagination(prev => ({ ...prev, page })),
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>

      <Modal
        title="添加用户"
        open={createModalOpen}
        onCancel={() => setCreateModalOpen(false)}
        footer={null}
        width={500}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreate}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          >
            <Input placeholder="请输入显示名称" />
          </Form.Item>

          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
          >
            <Input placeholder="请输入邮箱" />
          </Form.Item>

          <Form.Item
            name="phone"
            label="电话"
          >
            <Input placeholder="请输入电话" />
          </Form.Item>

          <Form.Item
            name="role"
            label="角色"
            rules={[{ required: true, message: '请选择角色' }]}
          >
            <Select placeholder="请选择角色">
              {roleOptions.map(opt => (
                <Option key={opt.value} value={opt.value}>{opt.label}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="team"
            label="团队"
          >
            <Select placeholder="请选择团队" allowClear>
              {teamOptions.map(opt => (
                <Option key={opt.value} value={opt.value}>{opt.label}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setCreateModalOpen(false)}>取消</Button>
              <Button type="primary" htmlType="submit">创建</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
