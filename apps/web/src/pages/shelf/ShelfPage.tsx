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
  Tooltip,
  Tabs,
  Statistic,
  Empty,
  Drawer,
  Descriptions,
  Timeline,
  Divider,
  Avatar,
  List,
} from 'antd'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  ShoppingCartOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  HistoryOutlined,
  DownloadOutlined,
  SearchOutlined,
  FilterOutlined,
  AppstoreOutlined,
  BranchesOutlined,
} from '@ant-design/icons'
import dayjs from 'dayjs'
import { 
  productApi, 
  cartApi, 
  Product, 
  ProductVersion,
  getTRLColor, 
  getTRLName 
} from '@/services/shelf'

const { Title, Text, Paragraph } = Typography
const { Option } = Select
const { TabPane } = Tabs

const ShelfPage: React.FC = () => {
  // State
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(12)
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [selectedTRL, setSelectedTRL] = useState<number | undefined>()
  const [searchQuery, setSearchQuery] = useState('')
  const [categories, setCategories] = useState<string[]>([])
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [currentProduct, setCurrentProduct] = useState<Product | null>(null)
  const [isDetailVisible, setIsDetailVisible] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [cartItemCount, setCartItemCount] = useState(0)
  const [form] = Form.useForm()

  // Load products
  const fetchProducts = useCallback(async () => {
    setLoading(true)
    try {
      const response = await productApi.list({
        page: currentPage,
        page_size: pageSize,
        category: selectedCategory || undefined,
        trl_level: selectedTRL,
        is_published: true,
        search: searchQuery || undefined,
      })
      const data = response.data.data
      setProducts(data?.items || [])
      setTotal(data?.total || 0)
    } catch (error) {
      message.error('加载产品列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }, [currentPage, pageSize, selectedCategory, selectedTRL, searchQuery])

  // Load categories
  const fetchCategories = async () => {
    try {
      const response = await productApi.getCategories()
      setCategories(response.data.data || [])
    } catch (error) {
      console.error('加载分类失败:', error)
    }
  }

  // Load cart count
  const fetchCartCount = async () => {
    try {
      const response = await cartApi.getItems()
      const items = response.data.data || []
      setCartItemCount(items.reduce((sum: number, item: any) => sum + item.quantity, 0))
    } catch (error) {
      console.error('加载购物车失败:', error)
    }
  }

  useEffect(() => {
    fetchProducts()
    fetchCategories()
    fetchCartCount()
  }, [fetchProducts])

  // Handlers
  const handleAdd = () => {
    setIsEditMode(false)
    setCurrentProduct(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (product: Product) => {
    setIsEditMode(true)
    setCurrentProduct(product)
    form.setFieldsValue({
      name: product.name,
      description: product.description,
      trl_level: product.trl_level,
      category: product.category,
      version: product.version,
    })
    setIsModalVisible(true)
  }

  const handleView = async (product: Product) => {
    try {
      const response = await productApi.getById(product.id)
      setSelectedProduct(response.data.data)
      setIsDetailVisible(true)
    } catch (error) {
      message.error('加载产品详情失败')
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await productApi.delete(id)
      message.success('产品删除成功')
      fetchProducts()
    } catch (error: any) {
      message.error(error.response?.data?.message || '删除失败')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (isEditMode && currentProduct) {
        await productApi.update(currentProduct.id, values)
        message.success('产品更新成功')
      } else {
        await productApi.create(values)
        message.success('产品创建成功')
      }
      setIsModalVisible(false)
      fetchProducts()
    } catch (error: any) {
      message.error(error.response?.data?.message || '操作失败')
    }
  }

  const handleAddToCart = async (product: Product) => {
    try {
      await cartApi.add({
        product_id: product.id,
        quantity: 1,
        notes: '',
      })
      message.success(`已将 "${product.name}" 添加到选用清单`)
      fetchCartCount()
    } catch (error: any) {
      message.error(error.response?.data?.message || '添加失败')
    }
  }

  const handlePublish = async (product: Product) => {
    try {
      if (product.is_published) {
        await productApi.unpublish(product.id)
        message.success('产品已下架')
      } else {
        await productApi.publish(product.id)
        message.success('产品已发布')
      }
      fetchProducts()
    } catch (error: any) {
      message.error(error.response?.data?.message || '操作失败')
    }
  }

  // TRL color mapping
  const getTRLTagColor = (level: number) => {
    const colors: Record<number, string> = {
      1: 'red',
      2: 'red',
      3: 'red',
      4: 'orange',
      5: 'orange',
      6: 'orange',
      7: 'green',
      8: 'green',
      9: 'green',
    }
    return colors[level] || 'default'
  }

  // Columns for table view
  const columns = [
    {
      title: '产品名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Product) => (
        <div>
          <Text strong>{text}</Text>
          {record.version && (
            <Tag color="blue" style={{ marginLeft: 8 }}>
              v{record.version}
            </Tag>
          )}
        </div>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      width: 120,
      render: (category: string) => category || '-',
    },
    {
      title: '成熟度',
      dataIndex: 'trl_level',
      key: 'trl_level',
      width: 100,
      render: (level: number) => (
        <Tag color={getTRLTagColor(level)}>
          TRL {level}
        </Tag>
      ),
    },
    {
      title: '下载',
      dataIndex: 'download_count',
      key: 'download_count',
      width: 80,
      render: (count: number) => (
        <Badge count={count} showZero style={{ backgroundColor: '#52c41a' }} />
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_published',
      key: 'is_published',
      width: 80,
      render: (isPublished: boolean) => (
        isPublished ? (
          <Tag color="success">已发布</Tag>
        ) : (
          <Tag>未发布</Tag>
        )
      ),
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 150,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, record: Product) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleView(record)}
            />
          </Tooltip>
          
          <Tooltip title="添加到选用清单">
            <Button
              type="text"
              icon={<ShoppingCartOutlined />}
              onClick={() => handleAddToCart(record)}
            />
          </Tooltip>
          
          <Tooltip title="编辑">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          
          <Tooltip title={record.is_published ? '下架' : '发布'}>
            <Button
              type="text"
              icon={record.is_published ? <CloseCircleOutlined /> : <CheckCircleOutlined />}
              onClick={() => handlePublish(record)}
            />
          </Tooltip>
          
          <Popconfirm
            title="确定删除此产品？"
            onConfirm={() => handleDelete(record.id)}
            okText="删除"
            cancelText="取消"
          >
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  // Card view for products
  const renderProductCard = (product: Product) => (
    <Card
      hoverable
      className="product-card"
      style={{ height: '100%' }}
      actions={[
        <Tooltip title="查看详情">
          <EyeOutlined onClick={() => handleView(product)} />
        </Tooltip>,
        <Tooltip title="添加到选用清单">
          <ShoppingCartOutlined onClick={() => handleAddToCart(product)} />
        </Tooltip>,
        <Tooltip title="编辑">
          <EditOutlined onClick={() => handleEdit(product)} />
        </Tooltip>,
      ]}
    >
      <Card.Meta
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Text strong>{product.name}</Text>
            {product.is_published ? (
              <Badge status="success" />
            ) : (
              <Badge status="default" />
            )}
          </div>
        }
        description={
          <div>
            <Tag color={getTRLTagColor(product.trl_level)} style={{ marginBottom: 8 }}>
              TRL {product.trl_level} - {getTRLName(product.trl_level)}
            </Tag>
            <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 8 }}>
              {product.description || '暂无描述'}
            </Paragraph>
            <Space size="small">
              {product.category && <Tag>{product.category}</Tag>}
              {product.version && <Tag color="blue">v{product.version}</Tag>}
            </Space>
            <div style={{ marginTop: 8 }}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                <DownloadOutlined /> {product.download_count} 次下载
              </Text>
            </div>
          </div>
        }
      />
    </Card>
  )

  return (
    <div style={{ padding: '24px' }}>
      {/* Header */}
      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={4}>产品货架</Title>
            <Text type="secondary">浏览、选用已成熟的模块产品</Text>
          </div>
          <Space>
            <Button 
              icon={<ShoppingCartOutlined />}
              badge={{ count: cartItemCount }}
            >
              选用清单 ({cartItemCount})
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              发布产品
            </Button>
          </Space>
        </div>
      </Card>

      {/* Filters */}
      <Card style={{ marginBottom: 24 }}>
        <Row gutter={16} align="middle">
          <Col span={8}>
            <Input.Search
              placeholder="搜索产品名称或描述"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onSearch={fetchProducts}
              allowClear
            />
          </Col>
          <Col span={6}>
            <Select
              placeholder="选择分类"
              value={selectedCategory || undefined}
              onChange={(value) => setSelectedCategory(value)}
              style={{ width: '100%' }}
              allowClear
            >
              {categories.map((cat) => (
                <Option key={cat} value={cat}>{cat}</Option>
              ))}
            </Select>
          </Col>
          <Col span={6}>
            <Select
              placeholder="成熟度等级"
              value={selectedTRL}
              onChange={(value) => setSelectedTRL(value)}
              style={{ width: '100%' }}
              allowClear
            >
              {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((level) => (
                <Option key={level} value={level}>
                  TRL {level} - {getTRLName(level)}
                </Option>
              ))}
            </Select>
          </Col>
          <Col span={4}>
            <Button onClick={() => {
              setSelectedCategory('')
              setSelectedTRL(undefined)
              setSearchQuery('')
              fetchProducts()
            }}>
              重置筛选
            </Button>
          </Col>
        </Row>
      </Card>

      {/* Stats */}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="产品总数"
              value={total}
              prefix={<AppstoreOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="已发布"
              value={products.filter(p => p.is_published).length}
              valueStyle={{ color: '#3f8600' }}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="分类数量"
              value={categories.length}
              prefix={<FilterOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="选用清单"
              value={cartItemCount}
              prefix={<ShoppingCartOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* Product List */}
      <Tabs defaultActiveKey="grid">
        <TabPane tab="网格视图" key="grid">
          <Row gutter={[16, 16]}>
            {products.map((product) => (
              <Col xs={24} sm={12} md={8} lg={6} key={product.id}>
                {renderProductCard(product)}
              </Col>
            ))}
          </Row>
          {products.length === 0 && (
            <Empty description="暂无产品" style={{ marginTop: 48 }} />
          )}
        </TabPane>
        <TabPane tab="列表视图" key="list">
          <Table
            columns={columns}
            dataSource={products}
            rowKey="id"
            loading={loading}
            pagination={{
              current: currentPage,
              pageSize: pageSize,
              total: total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total) => `共 ${total} 个产品`,
              onChange: (page, size) => {
                setCurrentPage(page)
                setPageSize(size || 12)
              },
            }}
          />
        </TabPane>
      </Tabs>

      {/* Create/Edit Modal */}
      <Modal
        title={isEditMode ? '编辑产品' : '发布新产品'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
        width={700}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="产品名称"
            rules={[{ required: true, message: '请输入产品名称' }]}
          >
            <Input placeholder="输入产品名称" />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="trl_level"
                label="成熟度等级 (TRL)"
                rules={[{ required: true, message: '请选择TRL等级' }]}
              >
                <Select placeholder="选择TRL等级">
                  {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((level) => (
                    <Option key={level} value={level}>
                      TRL {level} - {getTRLName(level)}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="category"
                label="分类"
              >
                <Select 
                  placeholder="选择或输入分类"
                  allowClear
                  showSearch
                  mode="tags"
                  maxCount={1}
                >
                  {categories.map((cat) => (
                    <Option key={cat} value={cat}>{cat}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="version"
            label="版本号"
          >
            <Input placeholder="例如：v1.0.0" />
          </Form.Item>

          <Form.Item
            name="description"
            label="产品描述"
          >
            <Input.TextArea 
              placeholder="详细描述产品功能、适用场景等"
              rows={4}
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => setIsModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit">
                {isEditMode ? '更新' : '发布'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* Product Detail Drawer */}
      <Drawer
        title="产品详情"
        placement="right"
        width={600}
        onClose={() => setIsDetailVisible(false)}
        open={isDetailVisible}
      >
        {selectedProduct && (
          <div>
            <div style={{ marginBottom: 24 }}>
              <Title level={4}>{selectedProduct.name}</Title>
              <Space size="middle">
                <Tag color={getTRLTagColor(selectedProduct.trl_level)}>
                  TRL {selectedProduct.trl_level} - {getTRLName(selectedProduct.trl_level)}
                </Tag>
                {selectedProduct.is_published ? (
                  <Tag color="success">已发布</Tag>
                ) : (
                  <Tag>未发布</Tag>
                )}
              </Space>
            </div>

            <Descriptions bordered column={1} style={{ marginBottom: 24 }}>
              <Descriptions.Item label="分类">{selectedProduct.category || '-'}</Descriptions.Item>
              <Descriptions.Item label="版本">{selectedProduct.version || '-'}</Descriptions.Item>
              <Descriptions.Item label="下载次数">{selectedProduct.download_count}</Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {dayjs(selectedProduct.created_at).format('YYYY-MM-DD HH:mm')}
              </Descriptions.Item>
              <Descriptions.Item label="更新时间">
                {dayjs(selectedProduct.updated_at).format('YYYY-MM-DD HH:mm')}
              </Descriptions.Item>
            </Descriptions>

            <Divider />

            <Paragraph>
              <Text strong>产品描述</Text>
              <div style={{ marginTop: 8 }}>
                {selectedProduct.description || '暂无描述'}
              </div>
            </Paragraph>

            <Divider />

            {selectedProduct.versions && selectedProduct.versions.length > 0 && (
              <div>
                <Text strong>版本历史</Text>
                <Timeline style={{ marginTop: 16 }}>
                  {selectedProduct.versions.map((version: ProductVersion) => (
                    <Timeline.Item key={version.id}>
                      <Text strong>v{version.version}</Text>
                      <div style={{ color: '#666', fontSize: 12 }}>
                        {dayjs(version.created_at).format('YYYY-MM-DD')}
                      </div>
                      {version.changelog && (
                        <div style={{ marginTop: 4 }}>{version.changelog}</div>
                      )}
                    </Timeline.Item>
                  ))}
                </Timeline>
              </div>
            )}

            <Divider />

            <Space>
              <Button 
                type="primary" 
                icon={<ShoppingCartOutlined />}
                onClick={() => handleAddToCart(selectedProduct)}
              >
                添加到选用清单
              </Button>
              <Button icon={<DownloadOutlined />}>下载资料</Button>
            </Space>
          </div>
        )}
      </Drawer>
    </div>
  )
}

export default ShelfPage
