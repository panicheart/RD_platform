import React, { useState, useEffect } from 'react'
import { Tree, Input, Button, Space, Modal, Form, message, Dropdown, MenuProps } from 'antd'
import { 
  FolderOutlined, 
  FolderOpenOutlined, 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  MoreOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined
} from '@ant-design/icons'
import { categoryApi, Category } from '@/services/knowledge'
import type { DataNode, TreeProps } from 'antd/es/tree'

interface CategoryTreeProps {
  onSelect?: (categoryId: string, categoryName: string) => void
  selectedCategoryId?: string
  showActions?: boolean
}

const CategoryTree: React.FC<CategoryTreeProps> = ({ 
  onSelect, 
  selectedCategoryId,
  showActions = true 
}) => {
  const [treeData, setTreeData] = useState<DataNode[]>([])
  const [loading, setLoading] = useState(false)
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([])
  const [searchValue, setSearchValue] = useState('')
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [isEditMode, setIsEditMode] = useState(false)
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null)
  const [parentCategory, setParentCategory] = useState<Category | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    setLoading(true)
    try {
      const response = await categoryApi.getTree()
      const categories = response.data.data || []
      const nodes = convertToTreeNodes(categories)
      setTreeData(nodes)
    } catch (error) {
      message.error('Failed to load categories')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const convertToTreeNodes = (categories: Category[]): DataNode[] => {
    return categories.map(cat => ({
      key: cat.id,
      title: cat.name,
      children: cat.children ? convertToTreeNodes(cat.children) : undefined,
      icon: ({ expanded }: { expanded?: boolean }) => 
        expanded ? <FolderOpenOutlined style={{ color: '#faad14' }} /> : <FolderOutlined style={{ color: '#faad14' }} />,
      data: cat
    }))
  }

  const handleSelect: TreeProps['onSelect'] = (selectedKeys) => {
    if (selectedKeys.length > 0 && onSelect) {
      const selectedNode = findNodeByKey(treeData, selectedKeys[0] as string)
      if (selectedNode?.data) {
        onSelect(selectedNode.data.id, selectedNode.data.name)
      }
    }
  }

  const findNodeByKey = (nodes: DataNode[], key: string): DataNode | null => {
    for (const node of nodes) {
      if (node.key === key) return node
      if (node.children) {
        const found = findNodeByKey(node.children, key)
        if (found) return found
      }
    }
    return null
  }

  const handleAdd = (parent?: Category) => {
    setIsEditMode(false)
    setCurrentCategory(null)
    setParentCategory(parent || null)
    form.resetFields()
    if (parent) {
      form.setFieldsValue({ parent_name: parent.name })
    }
    setIsModalVisible(true)
  }

  const handleEdit = (category: Category) => {
    setIsEditMode(true)
    setCurrentCategory(category)
    form.setFieldsValue({
      name: category.name,
      description: category.description
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (category: Category) => {
    Modal.confirm({
      title: 'Delete Category',
      content: `Are you sure you want to delete "${category.name}"?`,
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await categoryApi.delete(category.id)
          message.success('Category deleted successfully')
          loadCategories()
        } catch (error: any) {
          message.error(error.response?.data?.message || 'Failed to delete category')
        }
      }
    })
  }

  const handleSubmit = async (values: { name: string; description?: string }) => {
    try {
      if (isEditMode && currentCategory) {
        await categoryApi.update(currentCategory.id, {
          name: values.name,
          description: values.description || ''
        })
        message.success('Category updated successfully')
      } else {
        await categoryApi.create({
          name: values.name,
          description: values.description,
          parent_id: parentCategory?.id
        })
        message.success('Category created successfully')
      }
      setIsModalVisible(false)
      loadCategories()
    } catch (error: any) {
      message.error(error.response?.data?.message || 'Operation failed')
    }
  }

  const getMenuItems = (node: DataNode): MenuProps['items'] => {
    const category = node.data as Category
    const items: MenuProps['items'] = []

    if (showActions) {
      if (category.level < 3) {
        items.push({
          key: 'add',
          icon: <PlusOutlined />,
          label: 'Add Subcategory',
          onClick: () => handleAdd(category)
        })
      }

      items.push(
        {
          key: 'edit',
          icon: <EditOutlined />,
          label: 'Edit',
          onClick: () => handleEdit(category)
        },
        {
          key: 'delete',
          icon: <DeleteOutlined />,
          label: 'Delete',
          danger: true,
          onClick: () => handleDelete(category)
        }
      )
    }

    return items
  }

  const renderTitle = (node: DataNode): React.ReactNode => {
    const category = node.data as Category
    const isSelected = selectedCategoryId === category.id

    return (
      <div 
        style={{ 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'space-between',
          padding: '4px 0',
          backgroundColor: isSelected ? '#e6f7ff' : 'transparent',
          borderRadius: 4
        }}
      >
        <span style={{ flex: 1, overflow: 'hidden', textOverflow: 'ellipsis' }}>
          {node.title as string}
        </span>
        {showActions && (
          <Dropdown 
            menu={{ items: getMenuItems(node) }} 
            trigger={['click']}
            onClick={(e) => e.stopPropagation()}
          >
            <Button 
              type="text" 
              size="small" 
              icon={<MoreOutlined />}
              onClick={(e) => e.stopPropagation()}
            />
          </Dropdown>
        )}
      </div>
    )
  }

  const processTreeData = (nodes: DataNode[]): DataNode[] => {
    return nodes.map(node => ({
      ...node,
      title: renderTitle(node),
      children: node.children ? processTreeData(node.children) : undefined
    }))
  }

  const filteredTreeData = React.useMemo(() => {
    if (!searchValue) return processTreeData(treeData)
    
    const filterNodes = (nodes: DataNode[]): DataNode[] => {
      return nodes.reduce((acc: DataNode[], node) => {
        const category = node.data as Category
        const match = category.name.toLowerCase().includes(searchValue.toLowerCase())
        
        const filteredChildren = node.children ? filterNodes(node.children) : []
        
        if (match || filteredChildren.length > 0) {
          acc.push({
            ...node,
            children: filteredChildren.length > 0 ? filteredChildren : undefined
          })
        }
        
        return acc
      }, [])
    }
    
    return processTreeData(filterTreeData(treeData))
  }, [treeData, searchValue, selectedCategoryId])

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Space direction="vertical" style={{ width: '100%', marginBottom: 16 }} size="small">
        <Input.Search
          placeholder="Search categories..."
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
          allowClear
        />
        {showActions && (
          <Button 
            type="primary" 
            icon={<PlusOutlined />} 
            onClick={() => handleAdd()}
            block
          >
            Add Category
          </Button>
        )}
      </Space>

      <div style={{ flex: 1, overflow: 'auto' }}>
        <Tree
          treeData={filteredTreeData}
          loadData={undefined}
          loadedKeys={[]}
          expandedKeys={expandedKeys}
          onExpand={(keys) => setExpandedKeys(keys)}
          onSelect={handleSelect}
          selectedKeys={selectedCategoryId ? [selectedCategoryId] : []}
          showIcon
          defaultExpandAll={false}
          autoExpandParent={false}
          blockNode
        />
      </div>

      <Modal
        title={isEditMode ? 'Edit Category' : (parentCategory ? `Add Subcategory to "${parentCategory.name}"` : 'Add Category')}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          {parentCategory && (
            <Form.Item name="parent_name" label="Parent Category">
              <Input disabled />
            </Form.Item>
          )}
          
          <Form.Item
            name="name"
            label="Category Name"
            rules={[{ required: true, message: 'Please enter category name' }]}
          >
            <Input placeholder="Enter category name" />
          </Form.Item>

          <Form.Item name="description" label="Description">
            <Input.TextArea 
              placeholder="Enter description (optional)" 
              rows={3}
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0, textAlign: 'right' }}>
            <Space>
              <Button onClick={() => setIsModalVisible(false)}>
                Cancel
              </Button>
              <Button type="primary" htmlType="submit">
                {isEditMode ? 'Update' : 'Create'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CategoryTree
