import { useEffect, useState } from 'react'
import { Row, Col, Card, List, Avatar, Tag, Badge, Button, Space, Typography, Progress, Empty, Skeleton } from 'antd'
import {
  ProjectOutlined,
  FileTextOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  PlusOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { api, ApiResponse, PaginatedResponse } from '@/utils/api'
import dayjs from 'dayjs'

const { Title, Text } = Typography

interface Project {
  id: string
  name: string
  code: string
  status: string
  progress: number
  leader_id?: string
}

interface Activity {
  id: string
  name: string
  status: string
  project_id: string
  project_name?: string
  end_date?: string
}

interface Notification {
  id: string
  title: string
  type: string
  is_read: boolean
  created_at: string
}

interface DashboardStats {
  myProjectsCount: number
  pendingTasksCount: number
  notificationsCount: number
}

export default function WorkbenchPage() {
  const navigate = useNavigate()
  const [projects, setProjects] = useState<Project[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [stats, setStats] = useState<DashboardStats>({
    myProjectsCount: 0,
    pendingTasksCount: 0,
    notificationsCount: 0,
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [projectsRes, activitiesRes, notificationsRes] = await Promise.all([
        api.get<ApiResponse<PaginatedResponse<Project>>>('/users/me/projects'),
        api.get<ApiResponse<Activity[]>>('/activities?status=pending'),
        api.get<ApiResponse<PaginatedResponse<Notification>>>('/notifications?unread_only=true&page_size=5'),
      ])

      setProjects(projectsRes.data.data?.items || [])
      setActivities(activitiesRes.data.data || [])
      setNotifications(notificationsRes.data.data?.items || [])
      setStats({
        myProjectsCount: projectsRes.data.data?.total || 0,
        pendingTasksCount: activitiesRes.data.data?.length || 0,
        notificationsCount: notificationsRes.data.data?.total || 0,
      })
    } catch (error) {
      console.error('Failed to fetch workbench data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      draft: { color: 'default', text: '草稿' },
      planning: { color: 'blue', text: '规划中' },
      in_progress: { color: 'processing', text: '进行中' },
      review: { color: 'orange', text: '审核中' },
      completed: { color: 'success', text: '已完成' },
      pending: { color: 'warning', text: '待处理' },
      done: { color: 'success', text: '已完成' },
    }
    const config = statusMap[status] || { color: 'default', text: status }
    return <Tag color={config.color}>{config.text}</Tag>
  }

  const formatDate = (dateStr: string) => {
    return dayjs(dateStr).format('YYYY-MM-DD')
  }

  if (loading) {
    return (
      <div>
        <div className="page-header">
          <Title level={2}>个人工作台</Title>
          <Text type="secondary">欢迎回来，这里是您的工作概览</Text>
        </div>
        <Row gutter={16} style={{ marginBottom: 24 }}>
          {[1, 2, 3].map((i) => (
            <Col span={8} key={i}>
              <Card><Skeleton active avatar paragraph={{ rows: 1 }} /></Card>
            </Col>
          ))}
        </Row>
        <Row gutter={24}>
          <Col span={12}><Card><Skeleton active paragraph={{ rows: 5 }} /></Card></Col>
          <Col span={12}><Card><Skeleton active paragraph={{ rows: 5 }} /></Card></Col>
        </Row>
      </div>
    )
  }

  return (
    <div>
      <div className="page-header">
        <Title level={2}>个人工作台</Title>
        <Text type="secondary">欢迎回来，这里是您的工作概览</Text>
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={8}>
          <Card
            hoverable
            onClick={() => navigate('/projects')}
            style={{ transition: 'transform 0.2s, box-shadow 0.2s', cursor: 'pointer' }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = 'translateY(-2px)'; e.currentTarget.style.boxShadow = '0 4px 12px rgba(0,0,0,0.1)' }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = ''; e.currentTarget.style.boxShadow = '' }}
          >
            <Card.Meta
              avatar={<Avatar size={48} icon={<ProjectOutlined />} style={{ backgroundColor: '#1890ff' }} />}
              title="我的项目"
              description={`${stats.myProjectsCount} 个项目`}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card
            hoverable
            style={{ transition: 'transform 0.2s, box-shadow 0.2s', cursor: 'pointer' }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = 'translateY(-2px)'; e.currentTarget.style.boxShadow = '0 4px 12px rgba(0,0,0,0.1)' }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = ''; e.currentTarget.style.boxShadow = '' }}
          >
            <Card.Meta
              avatar={<Avatar size={48} icon={<ClockCircleOutlined />} style={{ backgroundColor: '#faad14' }} />}
              title="待办任务"
              description={`${stats.pendingTasksCount} 个任务`}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card
            hoverable
            style={{ transition: 'transform 0.2s, box-shadow 0.2s', cursor: 'pointer' }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = 'translateY(-2px)'; e.currentTarget.style.boxShadow = '0 4px 12px rgba(0,0,0,0.1)' }}
            onMouseLeave={(e) => { e.currentTarget.style.transform = ''; e.currentTarget.style.boxShadow = '' }}
          >
            <Card.Meta
              avatar={<Badge count={stats.notificationsCount} size="small">
                <Avatar size={48} icon={<FileTextOutlined />} style={{ backgroundColor: '#52c41a' }} />
              </Badge>}
              title="未读消息"
              description={`${stats.notificationsCount} 条消息`}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={24}>
        {/* My Projects */}
        <Col span={12}>
          <Card
            title="我的项目"
            extra={
              <Button type="link" onClick={() => navigate('/projects')}>
                查看全部 <ArrowRightOutlined />
              </Button>
            }
          >
            {projects.length > 0 ? (
              <List
                loading={loading}
                dataSource={projects.slice(0, 5)}
                renderItem={(item) => (
                  <List.Item
                    key={item.id}
                    actions={[
                      <Button
                        key="view"
                        type="link"
                        onClick={() => navigate(`/projects/${item.id}`)}
                      >
                        查看
                      </Button>,
                    ]}
                  >
                    <List.Item.Meta
                      avatar={<Avatar icon={<ProjectOutlined />} />}
                      title={item.name}
                      description={
                        <Space direction="vertical" size={0}>
                          <Text type="secondary">编码: {item.code}</Text>
                          <Progress
                            percent={item.progress}
                            size="small"
                            style={{ width: 150 }}
                          />
                        </Space>
                      }
                    />
                    {getStatusTag(item.status)}
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无项目">
                <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/projects')}>
                  创建项目
                </Button>
              </Empty>
            )}
          </Card>
        </Col>

        {/* Pending Activities */}
        <Col span={12}>
          <Card
            title="待办任务"
            extra={
              <Button type="link">
                查看全部 <ArrowRightOutlined />
              </Button>
            }
          >
            {activities.length > 0 ? (
              <List
                loading={loading}
                dataSource={activities.slice(0, 5)}
                renderItem={(item) => (
                  <List.Item
                    key={item.id}
                    actions={[
                      <Button key="start" type="link" size="small">
                        开始
                      </Button>,
                    ]}
                  >
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          icon={<CheckCircleOutlined />}
                          style={{
                            backgroundColor:
                              item.status === 'pending' ? '#faad14' : '#52c41a',
                          }}
                        />
                      }
                      title={item.name}
                      description={
                        <Space>
                          <Text type="secondary">{item.project_name}</Text>
                          {item.end_date && (
                            <Text type="secondary">截止: {formatDate(item.end_date)}</Text>
                          )}
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无待办任务" />
            )}
          </Card>
        </Col>
      </Row>

      {/* Recent Notifications */}
      <Card title="最新通知" style={{ marginTop: 24 }}>
        {notifications.length > 0 ? (
          <List
            loading={loading}
            dataSource={notifications}
            renderItem={(item) => (
              <List.Item key={item.id}>
                <List.Item.Meta
                  avatar={
                    <Badge dot={!item.is_read}>
                      <Avatar icon={<FileTextOutlined />} />
                    </Badge>
                  }
                  title={item.title}
                  description={formatDate(item.created_at)}
                />
                <Tag>{item.type}</Tag>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无新通知" />
        )}
      </Card>
    </div>
  )
}
