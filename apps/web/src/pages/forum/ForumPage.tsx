import React, { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Card,
  Typography,
  Button,
  Table,
  Tag,
  Space,
  Tabs,
  Select,
  Badge,
  message,
} from 'antd'
import {
  PlusOutlined,
  FolderOutlined,
  CommentOutlined,
  EyeOutlined,
} from '@ant-design/icons'
import dayjs from 'dayjs'
import { boardApi, postApi } from '@/services/forum'
import type { ForumBoard, ForumPost } from '@/types/forum'

const { Title, Text } = Typography
const { TabPane } = Tabs
const { Option } = Select

const ForumPage: React.FC = () => {
  const navigate = useNavigate()
  const [boards, setBoards] = useState<ForumBoard[]>([])
  const [posts, setPosts] = useState<ForumPost[]>([])
  const [loadingBoards, setLoadingBoards] = useState(false)
  const [loadingPosts, setLoadingPosts] = useState(false)
  const [boardsTotal, setBoardsTotal] = useState(0)
  const [postsTotal, setPostsTotal] = useState(0)
  const [boardsPage, setBoardsPage] = useState(1)
  const [postsPage, setPostsPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [selectedCategory, setSelectedCategory] = useState<string>('')

  const fetchBoards = useCallback(async () => {
    setLoadingBoards(true)
    try {
      const response = await boardApi.list({
        page: boardsPage,
        page_size: pageSize,
        category: selectedCategory || undefined,
      })
      const data = response.data.data
      setBoards(data?.items || [])
      setBoardsTotal(data?.total || 0)
    } catch (error) {
      message.error('加载板块列表失败')
      console.error(error)
    } finally {
      setLoadingBoards(false)
    }
  }, [boardsPage, pageSize, selectedCategory])

  const fetchLatestPosts = useCallback(async () => {
    setLoadingPosts(true)
    try {
      const response = await postApi.list({
        page: postsPage,
        page_size: pageSize,
      })
      const data = response.data.data
      setPosts(data?.items || [])
      setPostsTotal(data?.total || 0)
    } catch (error) {
      message.error('加载帖子列表失败')
      console.error(error)
    } finally {
      setLoadingPosts(false)
    }
  }, [postsPage, pageSize])

  useEffect(() => {
    fetchBoards()
    fetchLatestPosts()
  }, [fetchBoards, fetchLatestPosts])

  const handleCategoryChange = (value: string) => {
    setSelectedCategory(value)
    setBoardsPage(1)
  }

  const handleCreatePost = () => {
    navigate('/forum/posts/create')
  }

  const handleEnterBoard = (boardId: string) => {
    navigate(`/forum/boards/${boardId}`)
  }

  const handleViewPost = (postId: string) => {
    navigate(`/forum/posts/${postId}`)
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

  const boardColumns = [
    {
      title: '板块',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: ForumBoard) => (
        <Space>
          <FolderOutlined />
          <Text strong>{text}</Text>
          <Tag color={getCategoryColor(record.category)}>
            {getCategoryLabel(record.category)}
          </Tag>
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '主题',
      dataIndex: 'topic_count',
      key: 'topic_count',
      width: 100,
      render: (count: number) => (
        <Badge count={count} showZero style={{ backgroundColor: '#52c41a' }} />
      ),
    },
    {
      title: '回复',
      dataIndex: 'post_count',
      key: 'post_count',
      width: 100,
      render: (count: number) => (
        <Badge count={count} showZero style={{ backgroundColor: '#1890ff' }} />
      ),
    },
    {
      title: '最后发帖',
      dataIndex: 'last_post_at',
      key: 'last_post_at',
      width: 150,
      render: (date: string) =>
        date ? dayjs(date).format('MM-DD HH:mm') : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_: unknown, record: ForumBoard) => (
        <Button type="link" onClick={() => handleEnterBoard(record.id)}>
          进入板块
        </Button>
      ),
    },
  ]

  const postColumns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text: string, record: ForumPost) => (
        <Space>
          {record.is_pinned && (
            <Tag color="red">置顶</Tag>
          )}
          <Text
            strong={record.is_pinned}
            style={{ cursor: 'pointer' }}
            onClick={() => handleViewPost(record.id)}
          >
            {text}
          </Text>
          {record.reply_count > 0 && (
            <Badge
              count={record.reply_count}
              style={{ backgroundColor: '#1890ff' }}
            />
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
      render: (count: number) => (
        <Space size="small">
          <EyeOutlined />
          {count}
        </Space>
      ),
    },
    {
      title: '发布时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150,
      render: (date: string) => dayjs(date).format('MM-DD HH:mm'),
    },
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Card>
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: 16,
          }}
        >
          <div>
            <Title level={4}>技术论坛</Title>
            <Text type="secondary">交流讨论、问题解答、技术分享</Text>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreatePost}>
            发帖
          </Button>
        </div>

        <Tabs defaultActiveKey="boards">
          <TabPane
            tab={
              <span>
                <FolderOutlined />
                板块列表 ({boardsTotal})
              </span>
            }
            key="boards"
          >
            <div style={{ marginBottom: 16 }}>
              <Select
                placeholder="选择分类"
                allowClear
                style={{ width: 200 }}
                value={selectedCategory || undefined}
                onChange={handleCategoryChange}
              >
                <Option value="tech">技术讨论</Option>
                <Option value="general">综合讨论</Option>
                <Option value="help">求助问答</Option>
                <Option value="announcement">公告通知</Option>
              </Select>
            </div>
            <Table
              columns={boardColumns}
              dataSource={boards}
              rowKey="id"
              loading={loadingBoards}
              pagination={{
                current: boardsPage,
                pageSize: pageSize,
                total: boardsTotal,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 个板块`,
                onChange: (page, size) => {
                  setBoardsPage(page)
                  if (size) setPageSize(size)
                },
              }}
            />
          </TabPane>

          <TabPane
            tab={
              <span>
                <CommentOutlined />
                最新帖子 ({postsTotal})
              </span>
            }
            key="posts"
          >
            <Table
              columns={postColumns}
              dataSource={posts}
              rowKey="id"
              loading={loadingPosts}
              pagination={{
                current: postsPage,
                pageSize: pageSize,
                total: postsTotal,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 个帖子`,
                onChange: (page, size) => {
                  setPostsPage(page)
                  if (size) setPageSize(size)
                },
              }}
            />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  )
}

export default ForumPage
