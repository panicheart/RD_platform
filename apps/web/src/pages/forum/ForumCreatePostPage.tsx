import React, { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import {
  Card,
  Typography,
  Button,
  Form,
  Input,
  Select,
  Space,
  Breadcrumb,
  message,
} from 'antd'
import {
  ArrowLeftOutlined,
  SendOutlined,
} from '@ant-design/icons'
import { boardApi, postApi } from '@/services/forum'
import type { ForumBoard } from '@/types/forum'

const { Title } = Typography
const { Option } = Select
const { TextArea } = Input

interface PostFormData {
  title: string
  content: string
  board_id: string
  tags?: string[]
}

const ForumCreatePostPage: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const initialBoardId = searchParams.get('boardId') || ''
  
  const [form] = Form.useForm<PostFormData>()
  const [boards, setBoards] = useState<ForumBoard[]>([])
  const [submitting, setSubmitting] = useState(false)
  const [boardLoading, setBoardLoading] = useState(false)

  useEffect(() => {
    fetchBoards()
  }, [])

  useEffect(() => {
    if (initialBoardId && boards.length > 0) {
      form.setFieldsValue({ board_id: initialBoardId })
    }
  }, [initialBoardId, boards, form])

  const fetchBoards = async () => {
    setBoardLoading(true)
    try {
      const response = await boardApi.list({ page_size: 100 })
      const activeBoards = (response.data.data?.items || []).filter(
        (board) => board.is_active
      )
      setBoards(activeBoards)
    } catch (error) {
      message.error('加载板块列表失败')
      console.error(error)
    } finally {
      setBoardLoading(false)
    }
  }

  const handleSubmit = async (values: PostFormData) => {
    setSubmitting(true)
    try {
      const response = await postApi.create({
        title: values.title,
        content: values.content,
        board_id: values.board_id,
        tags: values.tags || [],
        knowledge_id: null,
      })
      message.success('帖子发布成功')
      const newPost = response.data.data
      navigate(`/forum/posts/${newPost.id}`)
    } catch (error) {
      message.error('发布失败')
      console.error(error)
    } finally {
      setSubmitting(false)
    }
  }

  const handleCancel = () => {
    navigate('/forum')
  }

  const tagOptions = [
    '问题', '讨论', '分享', '教程', '求助', '公告', '活动'
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Breadcrumb style={{ marginBottom: 16 }}>
        <Breadcrumb.Item>
          <a onClick={() => navigate('/forum')}>技术论坛</a>
        </Breadcrumb.Item>
        <Breadcrumb.Item>发布新帖</Breadcrumb.Item>
      </Breadcrumb>

      <Card>
        <div style={{ marginBottom: 24 }}>
          <Space>
            <Button icon={<ArrowLeftOutlined />} onClick={handleCancel}>
              返回
            </Button>
            <Title level={4} style={{ margin: 0 }}>发布新帖</Title>
          </Space>
        </div>

        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          autoComplete="off"
        >
          <Form.Item
            name="board_id"
            label="选择板块"
            rules={[{ required: true, message: '请选择一个板块' }]}
          >
            <Select
              placeholder="请选择要发布的板块"
              loading={boardLoading}
              showSearch
              optionFilterProp="children"
            >
              {boards.map((board) => (
                <Option key={board.id} value={board.id}>
                  <div>
                    <strong>{board.name}</strong>
                    <div style={{ fontSize: 12, color: '#999' }}>
                      {board.description}
                    </div>
                  </div>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="title"
            label="标题"
            rules={[
              { required: true, message: '请输入标题' },
              { max: 100, message: '标题最多100个字符' },
            ]}
          >
            <Input
              placeholder="请输入帖子标题，简洁明了"
              maxLength={100}
              showCount
            />
          </Form.Item>

          <Form.Item
            name="tags"
            label="标签"
          >
            <Select
              mode="tags"
              placeholder="添加标签（可选）"
              allowClear
              style={{ width: '100%' }}
            >
              {tagOptions.map((tag) => (
                <Option key={tag} value={tag}>{tag}</Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[
              { required: true, message: '请输入内容' },
              { min: 10, message: '内容至少10个字符' },
            ]}
          >
            <TextArea
              placeholder="请输入帖子内容...&#10;支持 Markdown 格式&#10;&#10;- 使用 # 表示标题&#10;- 使用 **粗体** 或 *斜体*&#10;- 使用 ``` 包裹代码块"
              rows={15}
              showCount
              maxLength={10000}
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button
                type="primary"
                htmlType="submit"
                icon={<SendOutlined />}
                loading={submitting}
                size="large"
              >
                发布帖子
              </Button>
              <Button
                onClick={handleCancel}
                size="large"
              >
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>

        <div style={{ marginTop: 24, padding: 16, backgroundColor: '#f5f5f5', borderRadius: 8 }}>
          <Title level={5}>发帖须知</Title>
          <ul style={{ color: '#666', lineHeight: '1.8' }}>
            <li>请遵守社区规范，友善交流</li>
            <li>提问前请先搜索，避免重复问题</li>
            <li>提供清晰的标题和详细的描述</li>
            <li>代码请使用 Markdown 代码块格式</li>
            <li>选择合适的板块和标签</li>
          </ul>
        </div>
      </Card>
    </div>
  )
}

export default ForumCreatePostPage
