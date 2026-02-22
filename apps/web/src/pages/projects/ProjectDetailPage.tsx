import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Button, Space, Tag, Progress, Row, Col, Table, Avatar, List, Typography, Spin, message, Tabs } from 'antd'
import { EditOutlined, ArrowLeftOutlined, UserOutlined, CalendarOutlined, FlagOutlined } from '@ant-design/icons'
import { api, ApiResponse, getErrorMessage } from '@/utils/api'
import dayjs from 'dayjs'

const { Title, Text, Paragraph } = Typography

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
  actual_start_date?: string
  actual_end_date?: string
  leader_id?: string
  created_at: string
  members?: ProjectMember[]
}

interface ProjectMember {
  id: string
  user_id: string
  role: string
  user?: {
    id: string
    username: string
    display_name: string
    avatar_url?: string
  }
}

interface Activity {
  id: string
  name: string
  status: string
  progress: number
  start_date?: string
  end_date?: string
  assignee_id?: string
}

const categoryMap: Record<string, string> = {
  product: '产品开发',
  research: '技术研究',
  infrastructure: '基础建设',
  tech_support: '技术支持',
}

const statusMap: Record<string, { label: string; color: string }> = {
  draft: { label: '草稿', color: 'default' },
  planning: { label: '规划中', color: 'blue' },
  in_progress: { label: '进行中', color: 'processing' },
  review: { label: '审核中', color: 'orange' },
  completed: { label: '已完成', color: 'success' },
}

export default function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [project, setProject] = useState<Project | null>(null)
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('overview')

  const fetchProject = useCallback(async () => {
    try {
      const response = await api.get<ApiResponse<Project>>(`/projects/${id}`)
      setProject(response.data.data)
    } catch (error) {
      message.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }, [id])

  const fetchActivities = useCallback(async () => {
    try {
      const response = await api.get<ApiResponse<Activity[]>>(`/projects/${id}/activities`)
      setActivities(response.data.data || [])
    } catch (error) {
      console.error('Failed to fetch activities:', error)
    }
  }, [id])

  useEffect(() => {
    if (id) {
      fetchProject()
      fetchActivities()
    }
  }, [id, fetchProject, fetchActivities])

  const getStatusConfig = (status: string) => statusMap[status] || { label: status, color: 'default' }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!project) {
    return (
      <div>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/projects')}>
          返回项目列表
        </Button>
        <Card style={{ marginTop: 16 }}>
          <Text>项目不存在</Text>
        </Card>
      </div>
    )
  }

  const statusConfig = getStatusConfig(project.status)

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/projects')}>
          返回
        </Button>
        <Button type="primary" icon={<EditOutlined />} onClick={() => navigate(`/projects/${id}?edit=true`)}>
          编辑项目
        </Button>
      </Space>

      <Card>

        <Row gutter={24}>
          <Col span={16}>
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
              <div>
                <Title level={2} style={{ marginBottom: 8 }}>
                  {project.name}
                </Title>
                <Space>
                  <Tag color={statusConfig.color}>{statusConfig.label}</Tag>
                  <Text type="secondary">编码: {project.code}</Text>
                  <Text type="secondary">类别: {categoryMap[project.category] || project.category}</Text>
                </Space>
              </div>

              {project.description && (
                <Paragraph>{project.description}</Paragraph>
              )}

              <div style={{ maxWidth: 400 }}>
                <Text>进度: {project.progress}%</Text>
                <Progress percent={project.progress} status="active" />
              </div>
            </Space>
          </Col>
          <Col span={8}>
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label={<><CalendarOutlined /> 开始日期</>}>
                {project.start_date ? dayjs(project.start_date).format('YYYY-MM-DD') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label={<><CalendarOutlined /> 结束日期</>}>
                {project.end_date ? dayjs(project.end_date).format('YYYY-MM-DD') : '-'}
              </Descriptions.Item>
              <Descriptions.Item label={<><FlagOutlined /> 团队</>}>
                {project.team || '-'}
              </Descriptions.Item>
              <Descriptions.Item label={<><FlagOutlined /> 产品线</>}>
                {project.product_line || '-'}
              </Descriptions.Item>
            </Descriptions>
          </Col>
        </Row>
      </Card>

      <Card style={{ marginTop: 16 }}>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'overview',
              label: '项目概况',
              children: (
                <Row gutter={16}>
                  <Col span={12}>
                    <Card title="项目信息" size="small">
                      <Descriptions column={1} bordered size="small">
                        <Descriptions.Item label="项目编码">{project.code}</Descriptions.Item>
                        <Descriptions.Item label="项目名称">{project.name}</Descriptions.Item>
                        <Descriptions.Item label="项目类别">{categoryMap[project.category]}</Descriptions.Item>
                        <Descriptions.Item label="项目状态">{statusConfig.label}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{dayjs(project.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
                      </Descriptions>
                    </Card>
                  </Col>
                  <Col span={12}>
                    <Card title="项目成员" size="small">
                      {project.members && project.members.length > 0 ? (
                        <List
                          dataSource={project.members}
                          renderItem={(member) => (
                            <List.Item>
                              <List.Item.Meta
                                avatar={<Avatar src={member.user?.avatar_url} icon={<UserOutlined />} />}
                                title={member.user?.display_name || member.user?.username}
                                description={member.role === 'leader' ? '项目负责人' : member.role}
                              />
                            </List.Item>
                          )}
                        />
                      ) : (
                        <Text type="secondary">暂无成员</Text>
                      )}
                    </Card>
                  </Col>
                </Row>
              ),
            },
            {
              key: 'activities',
              label: '项目活动',
              children: (
                <Table
                  dataSource={activities}
                  rowKey="id"
                  size="small"
                  columns={[
                    { title: '活动名称', dataIndex: 'name', key: 'name' },
                    {
                      title: '状态',
                      dataIndex: 'status',
                      key: 'status',
                      render: (status: string) => {
                        const config = statusMap[status] || { label: status, color: 'default' }
                        return <Tag color={config.color}>{config.label}</Tag>
                      }
                    },
                    {
                      title: '进度',
                      dataIndex: 'progress',
                      key: 'progress',
                      render: (progress: number) => <Progress percent={progress} size="small" />
                    },
                    { title: '开始日期', dataIndex: 'start_date', key: 'start_date', render: (d: string) => d || '-' },
                    { title: '结束日期', dataIndex: 'end_date', key: 'end_date', render: (d: string) => d || '-' },
                  ]}
                  locale={{ emptyText: '暂无活动' }}
                />
              ),
            },
            {
              key: 'files',
              label: '项目文件',
              children: (
                <Text type="secondary">文件管理功能开发中...</Text>
              ),
            },
          ]}
        />
      </Card>
    </div>
  )
}
