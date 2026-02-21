import { useEffect, useState } from 'react'
import { Table, Button, Space, Tag, Input, Select, Card, Modal, Form, DatePicker, Progress, Typography, message } from 'antd'
import { PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { api, ApiResponse, PaginatedResponse, getErrorMessage } from '@/utils/api'
import dayjs from 'dayjs'

const { Title, Text } = Typography
const { Option } = Select
const { TextArea } = Input

interface Project {
  id: string
  code: string
  name: string
  description?: string
  category: string
  status: string
  progress: number
  product_line?: string
  team?: string
  start_date?: string
  end_date?: string
  leader_id?: string
  created_at: string
}

const categoryOptions = [
  { value: 'product', label: '产品开发' },
  { value: 'research', label: '技术研究' },
  { value: 'infrastructure', label: '基础建设' },
  { value: 'tech_support', label: '技术支持' },
]

const statusOptions = [
  { value: 'draft', label: '草稿', color: 'default' },
  { value: 'planning', label: '规划中', color: 'blue' },
  { value: 'in_progress', label: '进行中', color: 'processing' },
  { value: 'review', label: '审核中', color: 'orange' },
  { value: 'completed', label: '已完成', color: 'success' },
]

export default function ProjectsPage() {
  const navigate = useNavigate()
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ page: 1, pageSize: 20, total: 0 })
  const [filters, setFilters] = useState({ search: '', status: '', category: '' })
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchProjects()
  }, [pagination.page, filters])

  const fetchProjects = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      params.append('page', pagination.page.toString())
      params.append('page_size', pagination.pageSize.toString())
      if (filters.search) params.append('search', filters.search)
      if (filters.status) params.append('status', filters.status)
      if (filters.category) params.append('category', filters.category)

      const response = await api.get<ApiResponse<PaginatedResponse<Project>>>(`/projects?${params}`)
      setProjects(response.data.data?.items || [])
      setPagination(prev => ({ ...prev, total: response.data.data?.total || 0 }))
    } catch (error) {
      message.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (values: any) => {
    try {
      const data = {
        ...values,
        start_date: values.dates?.[0]?.format('YYYY-MM-DD'),
        end_date: values.dates?.[1]?.format('YYYY-MM-DD'),
      }
      delete data.dates

      await api.post('/projects', data)
      message.success('项目创建成功')
      setCreateModalOpen(false)
      form.resetFields()
      fetchProjects()
    } catch (error) {
      message.error(getErrorMessage(error))
    }
  }

  const handleDelete = async (id: string) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此项目吗？此操作不可撤销。',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await api.delete(`/projects/${id}`)
          message.success('项目已删除')
          fetchProjects()
        } catch (error) {
          message.error(getErrorMessage(error))
        }
      },
    })
  }

  const getStatusTag = (status: string) => {
    const option = statusOptions.find(s => s.value === status)
    return <Tag color={option?.color || 'default'}>{option?.label || status}</Tag>
  }

  const columns = [
    {
      title: '项目编码',
      dataIndex: 'code',
      key: 'code',
      width: 120,
    },
    {
      title: '项目名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Project) => (
        <a onClick={() => navigate(`/projects/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '类别',
      dataIndex: 'category',
      key: 'category',
      width: 100,
      render: (category: string) => {
        const opt = categoryOptions.find(c => c.value === category)
        return opt?.label || category
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      width: 150,
      render: (progress: number) => <Progress percent={progress} size="small" />,
    },
    {
      title: '时间',
      dataIndex: 'start_date',
      key: 'dates',
      width: 180,
      render: (_: any, record: Project) => {
        if (!record.start_date) return '-'
        const start = dayjs(record.start_date).format('YYYY-MM-DD')
        const end = record.end_date ? dayjs(record.end_date).format('YYYY-MM-DD') : '进行中'
        return `${start} ~ ${end}`
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_: any, record: Project) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/projects/${record.id}`)}
          >
            查看
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => navigate(`/projects/${record.id}?edit=true`)}
          >
            编辑
          </Button>
          <Button
            type="link"
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
            <Title level={2} style={{ marginBottom: 0 }}>项目管理</Title>
            <Text type="secondary">管理和查看所有研发项目</Text>
          </div>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateModalOpen(true)}
          >
            创建项目
          </Button>
        </Space>
      </div>

      <Card>
        <Space style={{ marginBottom: 16 }} wrap>
          <Input
            placeholder="搜索项目名称或编码"
            prefix={<SearchOutlined />}
            style={{ width: 250 }}
            value={filters.search}
            onChange={e => setFilters(prev => ({ ...prev, search: e.target.value }))}
            onPressEnter={fetchProjects}
          />
          <Select
            placeholder="项目状态"
            style={{ width: 150 }}
            allowClear
            value={filters.status || undefined}
            onChange={value => setFilters(prev => ({ ...prev, status: value || '' }))}
          >
            {statusOptions.map(opt => (
              <Option key={opt.value} value={opt.value}>{opt.label}</Option>
            ))}
          </Select>
          <Select
            placeholder="项目类别"
            style={{ width: 150 }}
            allowClear
            value={filters.category || undefined}
            onChange={value => setFilters(prev => ({ ...prev, category: value || '' }))}
          >
            {categoryOptions.map(opt => (
              <Option key={opt.value} value={opt.value}>{opt.label}</Option>
            ))}
          </Select>
          <Button type="primary" icon={<SearchOutlined />} onClick={fetchProjects}>
            查询
          </Button>
        </Space>

        <Table
          columns={columns}
          dataSource={projects}
          loading={loading}
          rowKey="id"
          pagination={{
            current: pagination.page,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: (page) => setPagination(prev => ({ ...prev, page })),
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条`,
          }}
        />
      </Card>

      <Modal
        title="创建项目"
        open={createModalOpen}
        onCancel={() => setCreateModalOpen(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreate}
        >
          <Form.Item
            name="code"
            label="项目编码"
            rules={[{ required: true, message: '请输入项目编码' }]}
          >
            <Input placeholder="例如: PRD-2026-001" />
          </Form.Item>

          <Form.Item
            name="name"
            label="项目名称"
            rules={[{ required: true, message: '请输入项目名称' }]}
          >
            <Input placeholder="请输入项目名称" />
          </Form.Item>

          <Form.Item
            name="description"
            label="项目描述"
          >
            <TextArea rows={3} placeholder="请输入项目描述" />
          </Form.Item>

          <Form.Item
            name="category"
            label="项目类别"
            rules={[{ required: true, message: '请选择项目类别' }]}
          >
            <Select placeholder="请选择项目类别">
              {categoryOptions.map(opt => (
                <Option key={opt.value} value={opt.value}>{opt.label}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="dates"
            label="项目时间"
          >
            <DatePicker.RangePicker style={{ width: '100%' }} />
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
