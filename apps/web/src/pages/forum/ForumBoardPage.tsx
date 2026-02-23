import React, { useState, useEffect, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Typography,
  Button,
  Table,
  Tag,
  Space,
  Breadcrumb,
  Badge,
  Input,
  Tooltip,
  message,
} from 'antd'
import {
  PlusOutlined,
  ArrowLeftOutlined,
  EyeOutlined,
  MessageOutlined,
  PushpinFilled,
  LockFilled,
  SearchOutlined,
} from '@ant-design/icons'
import dayjs from 'dayjs'
import { boardApi, postApi } from '@/services/forum'
import type { ForumBoard, ForumPost } from '@/types/forum'

const { Title, Text } = Typography
const { Search } = Input

const ForumBoardPage: React.FC = () => {
  const { boardId } = useParams<{ boardId: string }>()
  const navigate = useNavigate()
  const [board, setBoard] = useState<ForumBoard | null>(null)
  const [posts, setPosts] = useState<ForumPost[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [searchQuery, setSearchQuery] = useState('')

  const fetchBoard = useCallback(async () => {
    if (!boardId) return
    try {
      const response = await boardApi.getById(boardId)
      setBoard(response.data.data)
    } catch (error) {
      message.error('加载板块信息失败')
      console.error(error)
    }
  }, [boardId])

  const fetchPosts = useCallback(async () => {
    if (!boardId) return
    setLoading(true)
    try {
      const response = await postApi.listByBoard(boardId, {
        page: page,
        page_size: pageSize,
        search: searchQuery || undefined,
      })
      const data = response.data.data
      setPosts(data?.items || [])
      setTotal(data?.total || 0)
    } catch (error) {
      message.error('加载帖子列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [boardId, page, pageSize, searchQuery])

  useEffect(() => {
    fetchBoard()
    fetchPosts()
  }, [fetchBoard, fetchPosts])

  const handleCreatePost = () => {
    navigate(`/forum/posts/create?boardId=${boardId}`)
  }

  const handleViewPost = (postId: string) => {
    navigate(`/forum/posts/${postId}`)
  }

  const handleBack = () => {
    navigate('/forum')
  }

  const handleSearch = (value: string) => {
    setSearchQuery(value)
    setPage(1)
  }

  const getCategoryLabel = (category: string) => {
    const categories: Record<string, string> = {
      tech: '技术讨论',
      general: '综合讨论',
      help: '求助问答',
      announcement: '公告通知',
    }
    return categories[category] || category
  }

  const getCategoryColor = (category: string) => {
    const colors: Record<string, string> = {
      tech: 'blue',
      general: 'green',
      help: 'orange',
      announcement: 'red',
    }
    return colors[category] || 'default'
  }

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text: string, record: ForumPost) => (
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
          <Space>
            {record.is_pinned && (
              <Tooltip title="置顶">
                <PushpinFilled style={{ color: '#ff4d4f' }} />
              </Tooltip>
            )}
            {record.is_locked && (
              <Tooltip title="已锁定">
                <LockFilled style={{ color: '#faad14' }} />
              </Tooltip>
            )}
            <Text
              strong={record.is_pinned}
              style={{
                cursor: 'pointer',
                color: record.is_pinned ? '#1890ff' : 'inherit',
              }}
              onClick={() => handleViewPost(record.id)}
            >
              {text}
            </Text>
          </Space>
          {record.tags && record.tags.length > 0 && (
            <Space size="small" wrap>
                {record.tags.map((tag, index) => (
                  <Tag key={index}>{tag}</Tag>
                ))}
            </Space>
          )}
        </Space>
      ),
    },
    {
      title: '作者',
      dataIndex: 'author_name',
      key: 'author_name',
      width: 120,
    },
    {
      title: '浏览',
      dataIndex: 'view_count',
      key: 'view_count',
      width: 80,
      sorter: (a: ForumPost, b: ForumPost) => a.view_count - b.view_count,
      render: (count: number) => (
        <Space size="small">
          <EyeOutlined />
          {count}
        </Space>
      ),
    },
    {
      title: '回复',
      dataIndex: 'reply_count',
      key: 'reply_count',
      width: 80,
      sorter: (a: ForumPost, b: ForumPost) => a.reply_count - b.reply_count,
      render: (count: number) => (
        <Space size="small">
          <MessageOutlined />
          <Badge
            count={count}
            showZero
            style={{
              backgroundColor: count > 0 ? '#1890ff' : '#d9d9d9',
            }}
          />
        </Space>
      ),
    },
    {
      title: '最后回复',
      dataIndex: 'last_reply_at',
      key: 'last_reply_at',
      width: 150,
      render: (date: string, record: ForumPost) =>
        date ? dayjs(date).format('MM-DD HH:mm') : dayjs(record.created_at).format('MM-DD HH:mm'),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={handleBack}>技术论坛</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>{board?.name || '加载中...'}</Breadcrumb.Item>
      </Breadcrumb>

      <Card>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
            marginBottom: 24,
          }}
        >
          <div>
            <Space align="center">
              <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
                返回
              </Button>
              <Title level={4} style={{ margin: 0 }}>
                {board?.name || '加载中...'}
              </Title>
              {board && (
                <Tag color={getCategoryColor(board.category)}>
                  {getCategoryLabel(board.category)}
                </Tag>
              )}
            </Space>
            <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
              {board?.description}
            </Text>
            <Space style={{ marginTop: 8 }}>
              <Text type="secondary">
                主题: <strong>{board?.topic_count || 0}</strong>
              </Text>
              <Text type="secondary">
                回复: <strong>{board?.post_count || 0}</strong>
              </Text>
            </Space>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreatePost}>
            发帖
          </Button>
        </div>

        <div style={{ marginBottom: 16 }}>
          <Search
            placeholder="搜索帖子标题"
            allowClear
            enterButton={<SearchOutlined />}
            onSearch={handleSearch}
            style={{ width: 300 }}
          />
        </div>

        <Table
          columns={columns}
          dataSource={posts}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 个帖子`,
            onChange: (page, size) => {
              setPage(page)
              if (size) setPageSize(size)
            },
          }}
        />
      </Card>
    </div>
  )
}

export default ForumBoardPage
