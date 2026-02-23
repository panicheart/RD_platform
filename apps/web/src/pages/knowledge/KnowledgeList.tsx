import React, { useState, useEffect, useCallback } from 'react'
import {
  Card,
  Row,
  Col,
  Typography,
  Button,
  Table,
  Tag,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Badge,
  Tooltip
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  FileTextOutlined,
  TagOutlined,
  SendOutlined,
  CheckCircleOutlined,
  InboxOutlined
} from '@ant-design/icons'
import dayjs from 'dayjs'
import CategoryTree from '@/components/knowledge/CategoryTree'
import { knowledgeApi, tagApi, Knowledge, Tag as KnowledgeTag } from '@/services/knowledge'

const { Title, Text } = Typography
const { Option } = Select
const { TextArea } = Input

const KnowledgeList: React.FC = () => {
  const [knowledgeItems, setKnowledgeItems] = useState<Knowledge[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>('')
  const [selectedCategoryName, setSelectedCategoryName] = useState('All Categories')
  const [tags, setTags] = useState<KnowledgeTag[]>([])
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [currentKnowledge, setCurrentKnowledge] = useState<Knowledge | null>(null)
  const [isDetailModalVisible, setIsDetailModalVisible] = useState(false)
  const [form] = Form.useForm()

  const fetchKnowledge = useCallback(async () => {
    setLoading(true)
    try {
      const response = await knowledgeApi.list({
        page: currentPage,
        page_size: pageSize,
        category_id: selectedCategoryId || undefined
      })
      const data = response.data.data
      setKnowledgeItems(data?.items || [])
      setTotal(data?.total || 0)
    } catch (error) {
      message.error('Failed to load knowledge items')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [currentPage, pageSize, selectedCategoryId])

  const fetchTags = async () => {
    try {
      const response = await tagApi.list()
      setTags(response.data.data || [])
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  useEffect(() => {
    fetchKnowledge()
    fetchTags()
  }, [fetchKnowledge])

  const handleCategorySelect = (categoryId: string, categoryName: string) => {
    setSelectedCategoryId(categoryId)
    setSelectedCategoryName(categoryName)
    setCurrentPage(1)
  }

  const handleAdd = () => {
    setIsEditMode(false)
    setCurrentKnowledge(null)
    form.resetFields()
    if (selectedCategoryId) {
      form.setFieldsValue({ category_id: selectedCategoryId })
    }
    setIsModalVisible(true)
  }

  const handleEdit = (record: Knowledge) => {
    setIsEditMode(true)
    setCurrentKnowledge(record)
    form.setFieldsValue({
      title: record.title,
      content: record.content,
      category_id: record.category_id,
      tag_ids: record.tags?.map(t => t.id) || []
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await knowledgeApi.delete(id)
      message.success('Knowledge item deleted successfully')
      fetchKnowledge()
    } catch (error: any) {
      message.error(error.response?.data?.message || 'Failed to delete')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (isEditMode && currentKnowledge) {
        await knowledgeApi.update(currentKnowledge.id, values)
        message.success('Knowledge updated successfully')
      } else {
        await knowledgeApi.create({
          ...values,
          source: 'manual'
        })
        message.success('Knowledge created successfully')
      }
      setIsModalVisible(false)
      fetchKnowledge()
    } catch (error: any) {
      message.error(error.response?.data?.message || 'Operation failed')
    }
  }

  const handlePublish = async (id: string) => {
    try {
      await knowledgeApi.publish(id)
      message.success('Knowledge published successfully')
      fetchKnowledge()
    } catch (error: any) {
      message.error(error.response?.data?.message || 'Failed to publish')
    }
  }

  const handleArchive = async (id: string) => {
    try {
      await knowledgeApi.archive(id)
      message.success('Knowledge archived successfully')
      fetchKnowledge()
    } catch (error: any) {
      message.error(error.response?.data?.message || 'Failed to archive')
    }
  }

  const handleView = (record: Knowledge) => {
    setCurrentKnowledge(record)
    setIsDetailModalVisible(true)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'published': return 'success'
      case 'draft': return 'default'
      case 'archived': return 'warning'
      default: return 'default'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'published': return 'Published'
      case 'draft': return 'Draft'
      case 'archived': return 'Archived'
      default: return status
    }
  }

  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (text: string, record: Knowledge) => (
        <div>
          <Text strong>{text}</Text>
          {record.version > 1 && (
            <Tag color="blue" style={{ marginLeft: 8 }}>
              v{record.version}
            </Tag>
          )}
        </div>
      )
    },
    {
      title: 'Tags',
      dataIndex: 'tags',
      key: 'tags',
      width: 200,
      render: (tags: KnowledgeTag[]) => (
        <Space size={[0, 4]} wrap>
          {tags?.map(tag => (
            <Tag key={tag.id} color={tag.color || 'blue'}>
              {tag.name}
            </Tag>
          ))}
        </Space>
      )
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status)}
        </Tag>
      )
    },
    {
      title: 'Views',
      dataIndex: 'view_count',
      key: 'view_count',
      width: 80,
      render: (count: number) => (
        <Badge count={count} showZero style={{ backgroundColor: '#52c41a' }} />
      )
    },
    {
      title: 'Updated',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 150,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 200,
      render: (_: any, record: Knowledge) => (
        <Space size="small">
          <Tooltip title="View">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleView(record)}
            />
          </Tooltip>
          
          <Tooltip title="Edit">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          
          {record.status === 'draft' && (
            <Tooltip title="Publish">
              <Button
                type="text"
                icon={<CheckCircleOutlined />}
                onClick={() => handlePublish(record.id)}
              />
            </Tooltip>
          )}
          
          {record.status === 'published' && (
            <Tooltip title="Archive">
              <Button
                type="text"
                icon={<InboxOutlined />}
                onClick={() => handleArchive(record.id)}
              />
            </Tooltip>
          )}
          
          <Popconfirm
            title="Delete this knowledge item?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div style={{ padding: '24px' }}>
      <Row gutter={24}>
        <Col span={6}>
          <Card title="Categories" style={{ height: 'calc(100vh - 120px)' }}>
            <CategoryTree
              onSelect={handleCategorySelect}
              selectedCategoryId={selectedCategoryId}
              showActions={true}
            />
          </Card>
        </Col>
        
        <Col span={18}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
              <div>
                <Title level={4}>Knowledge Base</Title>
                <Text type="secondary">
                  {selectedCategoryName} • {total} items
                </Text>
              </div>
              <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                New Knowledge
              </Button>
            </div>
            
            <Table
              columns={columns}
              dataSource={knowledgeItems}
              rowKey="id"
              loading={loading}
              pagination={{
                current: currentPage,
                pageSize: pageSize,
                total: total,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `Total ${total} items`,
                onChange: (page, size) => {
                  setCurrentPage(page)
                  setPageSize(size || 20)
                }
              }}
            />
          </Card>
        </Col>
      </Row>

      <Modal
        title={isEditMode ? 'Edit Knowledge' : 'New Knowledge'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
        width={800}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="title"
            label="Title"
            rules={[{ required: true, message: 'Please enter title' }]}
          >
            <Input placeholder="Enter knowledge title" />
          </Form.Item>

          <Form.Item
            name="category_id"
            label="Category"
            rules={[{ required: true, message: 'Please select category' }]}
          >
            <Select placeholder="Select category">
              <Option value="">Select a category...</Option>
              {/* TODO: Load categories as options */}
            </Select>
          </Form.Item>

          <Form.Item name="tag_ids" label="Tags">
            <Select mode="multiple" placeholder="Select tags">
              {tags.map(tag => (
                <Option key={tag.id} value={tag.id}>
                  <Tag color={tag.color}>{tag.name}</Tag>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="Content"
            rules={[{ required: true, message: 'Please enter content' }]}
          >
            <TextArea 
              placeholder="Enter content (Markdown supported)" 
              rows={10}
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => setIsModalVisible(false)}>Cancel</Button>
              <Button type="primary" htmlType="submit">
                {isEditMode ? 'Update' : 'Create'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="Knowledge Detail"
        open={isDetailModalVisible}
        onCancel={() => setIsDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setIsDetailModalVisible(false)}>Close</Button>,
          <Button key="edit" type="primary" onClick={() => {
            setIsDetailModalVisible(false)
            currentKnowledge && handleEdit(currentKnowledge)
          }}>Edit</Button>
        ]}
        width={800}
      >
        {currentKnowledge && (
          <div>
            <div style={{ marginBottom: 16 }}>
              <Title level={3}>{currentKnowledge.title}</Title>
              <Space size="middle">
                <Tag color={getStatusColor(currentKnowledge.status)}>
                  {getStatusText(currentKnowledge.status)}
                </Tag>
                <Text type="secondary">
                  Views: {currentKnowledge.view_count}
                </Text>
                <Text type="secondary">
                  Version: {currentKnowledge.version}
                </Text>
              </Space>
            </div>
            
            <div style={{ marginBottom: 16 }}>
              {currentKnowledge.tags?.map(tag => (
                <Tag key={tag.id} color={tag.color} style={{ marginRight: 8 }}>
                  {tag.name}
                </Tag>
              ))}
            </div>
            
            <Card>
              <pre style={{ whiteSpace: 'pre-wrap', fontFamily: 'inherit' }}>
                {currentKnowledge.content}
              </pre>
            </Card>
            
            <div style={{ marginTop: 16 }}>
              <Text type="secondary">
                Created: {dayjs(currentKnowledge.created_at).format('YYYY-MM-DD HH:mm')}
              </Text>
              <br />
              <Text type="secondary">
                Updated: {dayjs(currentKnowledge.updated_at).format('YYYY-MM-DD HH:mm')}
              </Text>
            </div>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default KnowledgeList
