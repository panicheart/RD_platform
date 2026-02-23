import React, { useState, useEffect, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Typography,
  Button,
  Space,
  Breadcrumb,
  Avatar,
  Divider,
  List,
  Input,
  message,
  Popconfirm,
  Tag,
} from 'antd'
import {
  ArrowLeftOutlined,
  EditOutlined,
  DeleteOutlined,
  PushpinOutlined,
  LockOutlined,
  MessageOutlined,
  EyeOutlined,
  SendOutlined,
  UserOutlined,
  CheckCircleFilled,
} from '@ant-design/icons'
import dayjs from 'dayjs'
import { postApi, replyApi } from '@/services/forum'
import type { ForumPost, ForumReply } from '@/types/forum'

const { Title, Text } = Typography

const ForumPostPage: React.FC = () => {
  const { postId } = useParams<{ postId: string }>()
  const navigate = useNavigate()
  const [post, setPost] = useState<ForumPost | null>(null)
  const [replies, setReplies] = useState<ForumReply[]>([])
  const [loadingPost, setLoadingPost] = useState(false)
  const [loadingReplies, setLoadingReplies] = useState(false)
  const [replyContent, setReplyContent] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [replyTotal, setReplyTotal] = useState(0)
  const [replyPage, setReplyPage] = useState(1)
  const [replyPageSize, setReplyPageSize] = useState(20)
  const [replyingTo, setReplyingTo] = useState<ForumReply | null>(null)
  const replyInputRef = useRef<HTMLTextAreaElement>(null)

  const fetchPost = useCallback(async () => {
    if (!postId) return
    setLoadingPost(true)
    try {
      const response = await postApi.getById(postId)
      setPost(response.data.data)
    } catch (error) {
      message.error('加载帖子失败')
      console.error(error)
    } finally {
      setLoadingPost(false)
    }
  }, [postId])

  const fetchReplies = useCallback(async () => {
    if (!postId) return
    setLoadingReplies(true)
    try {
      const response = await replyApi.listByPost(postId, {
        page: replyPage,
        page_size: replyPageSize,
      })
      const data = response.data.data
      setReplies(data?.items || [])
      setReplyTotal(data?.total || 0)
    } catch (error) {
      message.error('加载回复失败')
      console.error(error)
    } finally {
      setLoadingReplies(false)
    }
  }, [postId, replyPage, replyPageSize])

  useEffect(() => {
    fetchPost()
    fetchReplies()
  }, [fetchPost, fetchReplies])

  const handleBack = () => {
    navigate(-1)
  }

  const handleReply = async () => {
    if (!postId || !replyContent.trim()) return
    setSubmitting(true)
    try {
      await replyApi.create(postId, {
        content: replyContent,
        parent_id: replyingTo?.id || null,
      })
      message.success('回复成功')
      setReplyContent('')
      setReplyingTo(null)
      fetchReplies()
      fetchPost()
    } catch (error) {
      message.error('回复失败')
      console.error(error)
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeletePost = async () => {
    if (!postId) return
    try {
      await postApi.delete(postId)
      message.success('帖子已删除')
      navigate('/forum')
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleDeleteReply = async (replyId: string) => {
    try {
      await replyApi.delete(replyId)
      message.success('回复已删除')
      fetchReplies()
    } catch (error) {
      message.error('删除失败')
      console.error(error)
    }
  }

  const handleReplyToReply = (reply: ForumReply) => {
    setReplyingTo(reply)
    setReplyContent(`@${reply.author_name} `)
    replyInputRef.current?.focus()
  }

  const cancelReplyTo = () => {
    setReplyingTo(null)
    setReplyContent('')
  }

  const renderReply = (reply: ForumReply, isNested = false) => (
    <List.Item
      key={reply.id}
      style={{
        padding: '16px',
        backgroundColor: isNested ? '#fafafa' : 'transparent',
        borderLeft: isNested ? '3px solid #1890ff' : 'none',
        marginLeft: isNested ? 40 : 0,
        marginBottom: isNested ? 8 : 0,
        borderRadius: isNested ? '0 8px 8px 0' : 0,
      }}
    >
      <List.Item.Meta
        avatar={
          <Avatar icon={<UserOutlined />} src={undefined} />
        }
        title={
          <Space>
            <Text strong>{reply.author_name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {dayjs(reply.created_at).format('YYYY-MM-DD HH:mm')}
            </Text>
            {reply.is_best_answer && (
              <Tag color="green" icon={<CheckCircleFilled />}>
                最佳答案
              </Tag>
            )}
          </Space>
        }
        description={
          <div style={{ marginTop: 8 }}>
            <div
              style={{
                whiteSpace: 'pre-wrap',
                lineHeight: 1.6,
                fontSize: 14,
              }}
            >
              {reply.content}
            </div>
            <Space style={{ marginTop: 12 }}>
              <Button
                type="link"
                size="small"
                onClick={() => handleReplyToReply(reply)}
              >
                回复
              </Button>
              {reply.author_id === localStorage.getItem('user_id') && (
                <Popconfirm
                  title="删除回复"
                  description="确定要删除这条回复吗？"
                  onConfirm={() => handleDeleteReply(reply.id)}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button type="link" danger size="small">
                    删除
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </div>
        }
      />
    </List.Item>
  )

  return (
    <div style={{ padding: '24px' }}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={() => navigate('/forum')}>技术论坛</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>{post?.title || '加载中...'}</Breadcrumb.Item>
      </Breadcrumb>

      <Card loading={loadingPost}>
        <div style={{ marginBottom: 16 }}>
          <Space>
            <Button icon={<ArrowLeftOutlined />} onClick={handleBack}>
              返回
            </Button>
            {post?.is_pinned && (
              <Tag color="red" icon={<PushpinOutlined />}>置顶</Tag>
            )}
            {post?.is_locked && (
              <Tag color="orange" icon={<LockOutlined />}>已锁定</Tag>
            )}
          </Space>
        </div>

        <div style={{ marginBottom: 24 }}>
          <Title level={3}>{post?.title}</Title>
          <Space split={<Divider type="vertical" />}>
            <Avatar icon={<UserOutlined />} size="small" />
            <Text>{post?.author_name}</Text>
            <Text type="secondary">
              发布于 {post?.created_at && dayjs(post.created_at).format('YYYY-MM-DD HH:mm')}
            </Text>
            <Space>
              <EyeOutlined />
              <Text>{post?.view_count} 次浏览</Text>
            </Space>
            <Space>
              <MessageOutlined />
              <Text>{post?.reply_count} 条回复</Text>
            </Space>
          </Space>

          {post?.tags && post.tags.length > 0 && (
            <div style={{ marginTop: 12 }}>
              <Space wrap>
                {post.tags.map((tag, index) => (
                  <Tag key={index}>{tag}</Tag>
                ))}
              </Space>
            </div>
          )}
        </div>

        <Divider />

        <div
          className="content-area"
          style={{
            minHeight: 200,
            whiteSpace: 'pre-wrap',
            lineHeight: 1.8,
            fontSize: 14,
          }}
        >
          {post?.content}
        </div>

        <Divider />

        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Button icon={<EditOutlined />}>编辑</Button>
          </Space>
          <Popconfirm
            title="删除帖子"
            description="确定要删除这个帖子吗？此操作不可恢复。"
            onConfirm={handleDeletePost}
            okText="确定"
            cancelText="取消"
          >
            <Button danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </div>
      </Card>

      <Card
        title={<div>
          <MessageOutlined /> 回复 ({replyTotal})
        </div>}
        style={{ marginTop: 24 }}
        loading={loadingReplies}
      >
        <List
          dataSource={replies}
          renderItem={(reply) => renderReply(reply)}
          pagination={
            replyTotal > replyPageSize
              ? {
                  current: replyPage,
                  pageSize: replyPageSize,
                  total: replyTotal,
                  onChange: (page, size) => {
                    setReplyPage(page)
                    if (size) setReplyPageSize(size)
                  },
                }
              : false
          }
          locale={{
            emptyText: '暂无回复，快来抢沙发吧！',
          }}
        />

        {!post?.is_locked && (
          <div style={{ marginTop: 24 }}>
            <Divider />
            {replyingTo && (
              <div style={{ marginBottom: 12 }}>
                <Space>
                  <Text type="secondary">
                    正在回复 <strong>@{replyingTo.author_name}</strong>
                  </Text>
                  <Button type="link" size="small" onClick={cancelReplyTo}>
                    取消
                  </Button>
                </Space>
              </div>
            )}
            <Input.TextArea
              ref={replyInputRef}
              rows={4}
              placeholder="写下你的回复..."
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              disabled={submitting}
            />
            <div style={{ marginTop: 12, textAlign: 'right' }}>
              <Button
                type="primary"
                icon={<SendOutlined />}
                loading={submitting}
                disabled={!replyContent.trim()}
                onClick={handleReply}
              >
                发表回复
              </Button>
            </div>
          </div>
        )}

        {post?.is_locked && (
          <div style={{ textAlign: 'center', padding: '24px' }}>
            <LockOutlined style={{ fontSize: 24, color: '#999' }} />
            <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
              该帖子已锁定，无法回复
            </Text>
          </div>
        )}
      </Card>
    </div>
  )
}

export default ForumPostPage
