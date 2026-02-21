import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Tree,
  Typography,
  Space,
  Avatar,
  Tag,
  Row,
  Col,
  List,
  Badge,
  Empty,
  message,
  Button,
  Modal,
  Form,
  Input,
  Select,
} from 'antd';
import {
  TeamOutlined,
  UserOutlined,
  PartitionOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ApartmentOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import { userAPI } from '@/services/user';
import type { Role } from '@/types';

const { Title, Text } = Typography;
const { DirectoryTree } = Tree;
const { Option } = Select;

interface OrgNode {
  id: string;
  key?: string;
  title?: string;
  name: string;
  code: string;
  type: 'department' | 'team' | 'group' | 'product_line';
  parentId?: string;
  children?: OrgNode[];
  memberCount?: number;
  description?: string;
}

interface OrgMember {
  id: string;
  username: string;
  displayName: string;
  avatar?: string;
  role: string;
  title?: string;
  email?: string;
}

export default function OrgChart() {
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  
  const [orgTree, setOrgTree] = useState<OrgNode[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedOrg, setSelectedOrg] = useState<OrgNode | null>(null);
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [membersLoading, setMembersLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [editingNode, setEditingNode] = useState<OrgNode | null>(null);

  // 权限检查
  const isAdmin = currentUser?.role === 'admin' || currentUser?.roles?.some((r: Role) => r.code === 'admin');

  useEffect(() => {
    fetchOrgTree();
  }, []);

  useEffect(() => {
    if (selectedOrg) {
      fetchOrgMembers(selectedOrg.id);
    }
  }, [selectedOrg]);

  const fetchOrgTree = async () => {
    setLoading(true);
    try {
      const data = await userAPI.getOrganizationTree();
      const treeData = buildTree(data);
      setOrgTree(treeData);
      // 默认选中第一个部门
      if (treeData.length > 0 && !selectedOrg) {
        setSelectedOrg(treeData[0] || null);
      }
    } catch (error: any) {
      message.error(error.message || '获取组织架构失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchOrgMembers = async (orgId: string) => {
    setMembersLoading(true);
    try {
      const data = await userAPI.getOrganizationMembers(orgId);
      setMembers(data as any);
    } catch (error: any) {
      message.error(error.message || '获取成员列表失败');
    } finally {
      setMembersLoading(false);
    }
  };

  // 构建树形结构
  const buildTree = (nodes: OrgNode[]): OrgNode[] => {
    const nodeMap = new Map<string, OrgNode>();
    const roots: OrgNode[] = [];

    // 首先创建所有节点的映射
    nodes.forEach(node => {
      nodeMap.set(node.id, { ...node, key: node.id, children: [] });
    });

    // 然后构建树形结构
    nodes.forEach(node => {
      const currentNode = nodeMap.get(node.id)!;
      if (node.parentId && nodeMap.has(node.parentId)) {
        const parent = nodeMap.get(node.parentId)!;
        if (!parent.children) parent.children = [];
        parent.children.push(currentNode);
      } else {
        roots.push(currentNode);
      }
    });

    return roots;
  };

  const getTypeTag = (type: string) => {
    const typeMap: Record<string, { color: string; text: string; icon: React.ReactNode }> = {
      department: { color: 'blue', text: '部门', icon: <ApartmentOutlined /> },
      team: { color: 'green', text: '团队', icon: <TeamOutlined /> },
      group: { color: 'orange', text: '小组', icon: <UserOutlined /> },
      product_line: { color: 'purple', text: '产品线', icon: <PartitionOutlined /> },
    };
    const config = typeMap[type] || { color: 'default', text: type, icon: null };
    return (
      <Tag color={config.color} icon={config.icon}>
        {config.text}
      </Tag>
    );
  };

  const getRoleTag = (role: string) => {
    const roleMap: Record<string, { color: string; text: string }> = {
      admin: { color: 'red', text: '管理员' },
      manager: { color: 'blue', text: '经理' },
      leader: { color: 'green', text: '负责人' },
      designer: { color: 'default', text: '工程师' },
    };
    const config = roleMap[role] || { color: 'default', text: role };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const handleSelect = (_: React.Key[], { node }: { node: any }) => {
    setSelectedOrg(node as OrgNode);
  };

  const handleAdd = () => {
    setEditingNode(null);
    form.resetFields();
    if (selectedOrg) {
      form.setFieldValue('parentId', selectedOrg.id);
    }
    setModalVisible(true);
  };

  const handleEdit = (node: OrgNode) => {
    setEditingNode(node);
    form.setFieldsValue(node);
    setModalVisible(true);
  };

  const handleDelete = (node: OrgNode) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除 "${node.name}" 吗？该部门下的所有成员将被移至未分配部门。`,
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await userAPI.deleteOrganization(node.id);
          message.success('部门已删除');
          fetchOrgTree();
        } catch (error: any) {
          message.error(error.message || '删除失败');
        }
      },
    });
  };

  const handleSave = async (values: any) => {
    try {
      if (editingNode) {
        await userAPI.updateOrganization(editingNode.id, values);
        message.success('部门已更新');
      } else {
        await userAPI.createOrganization(values);
        message.success('部门已创建');
      }
      setModalVisible(false);
      fetchOrgTree();
    } catch (error: any) {
      message.error(error.message || '保存失败');
    }
  };

  const titleRender = (node: OrgNode) => (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', width: '100%' }}>
      <Space>
        {getTypeTag(node.type)}
        <span>{node.name}</span>
        {node.memberCount !== undefined && (
          <Badge count={node.memberCount} style={{ backgroundColor: '#1890ff' }} />
        )}
      </Space>
      {isAdmin && (
        <Space size="small" style={{ marginLeft: 16 }}>
          <Button
            type="text"
            size="small"
            icon={<PlusOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              setSelectedOrg(node);
              handleAdd();
            }}
          />
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              handleEdit(node);
            }}
          />
          {!node.children?.length && (
            <Button
              type="text"
              size="small"
              danger
              icon={<DeleteOutlined />}
              onClick={(e) => {
                e.stopPropagation();
                handleDelete(node);
              }}
            />
          )}
        </Space>
      )}
    </div>
  );

  return (
    <div>
      {/* 页面头部 */}
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Title level={2} style={{ margin: 0 }}>
            <PartitionOutlined style={{ marginRight: 12 }} />
            组织架构
          </Title>
          <Text type="secondary">管理部门、团队和人员组织关系</Text>
        </Col>
        <Col>
          {isAdmin && (
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添加部门
            </Button>
          )}
        </Col>
      </Row>

      <Row gutter={24}>
        {/* 左侧：组织树 */}
        <Col span={8}>
          <Card title="组织层级" loading={loading}>
            {orgTree.length > 0 ? (
              <DirectoryTree
                treeData={orgTree as any}
                titleRender={titleRender as any}
                onSelect={handleSelect as any}
                selectedKeys={selectedOrg?.key ? [selectedOrg.key] : []}
                defaultExpandAll
              />
            ) : (
              <Empty description="暂无组织架构" />
            )}
          </Card>
        </Col>

        {/* 右侧：部门详情和成员 */}
        <Col span={16}>
          {selectedOrg ? (
            <>
              <Card style={{ marginBottom: 24 }}>
                <Row align="middle">
                  <Col flex="auto">
                    <Title level={4} style={{ margin: 0 }}>
                      {selectedOrg.name}
                    </Title>
                    <Space style={{ marginTop: 8 }}>
                      {getTypeTag(selectedOrg.type)}
                      <Text type="secondary">代码: {selectedOrg.code}</Text>
                    </Space>
                    {selectedOrg.description && (
                      <p style={{ marginTop: 8, color: '#666' }}>
                        {selectedOrg.description}
                      </p>
                    )}
                  </Col>
                  <Col>
                    <StatisticCard
                      title="成员数"
                      value={members.length}
                      icon={<TeamOutlined />}
                    />
                  </Col>
                </Row>
              </Card>

              <Card title="部门成员" loading={membersLoading}>
                <List
                  dataSource={members}
                  renderItem={(member) => (
                    <List.Item
                      actions={[
                        <Button
                          type="link"
                          onClick={() => navigate(`/users/${member.id}`)}
                        >
                          查看详情
                        </Button>,
                      ]}
                    >
                      <List.Item.Meta
                        avatar={<Avatar src={member.avatar} icon={<UserOutlined />} />}
                        title={
                          <Space>
                            <span>{member.displayName}</span>
                            {getRoleTag(member.role)}
                          </Space>
                        }
                        description={
                          <Space direction="vertical" size={0}>
                            <Text type="secondary">@{member.username}</Text>
                            {member.title && (
                              <Text type="secondary">{member.title}</Text>
                            )}
                            {member.email && (
                              <Text type="secondary">{member.email}</Text>
                            )}
                          </Space>
                        }
                      />
                    </List.Item>
                  )}
                />
              </Card>
            </>
          ) : (
            <Card>
              <Empty description="请选择左侧部门查看详情" />
            </Card>
          )}
        </Col>
      </Row>

      {/* 添加/编辑部门弹窗 */}
      <Modal
        title={editingNode ? '编辑部门' : '添加部门'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={500}
      >
        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item
            name="name"
            label="部门名称"
            rules={[{ required: true, message: '请输入部门名称' }]}
          >
            <Input placeholder="请输入部门名称" />
          </Form.Item>

          <Form.Item
            name="code"
            label="部门代码"
            rules={[
              { required: true, message: '请输入部门代码' },
              { pattern: /^[a-zA-Z0-9_]+$/, message: '代码只能包含字母、数字和下划线' },
            ]}
          >
            <Input placeholder="如: software_dept" disabled={!!editingNode} />
          </Form.Item>

          <Form.Item
            name="type"
            label="组织类型"
            rules={[{ required: true }]}
          >
            <Select placeholder="请选择类型">
              <Option value="department">部门</Option>
              <Option value="team">团队</Option>
              <Option value="group">小组</Option>
              <Option value="product_line">产品线</Option>
            </Select>
          </Form.Item>

          <Form.Item name="parentId" label="上级部门">
            <Select placeholder="不选则为顶级部门" allowClear>
              {flattenTree(orgTree).map((node) => (
                <Option key={node.id} value={node.id}>
                  {node.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入部门描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

// 辅助组件：统计卡片
function StatisticCard({ title, value, icon }: { title: string; value: number; icon: React.ReactNode }) {
  return (
    <div style={{ textAlign: 'center', padding: '0 24px' }}>
      <div style={{ fontSize: 24, color: '#1890ff', marginBottom: 8 }}>{icon}</div>
      <div style={{ fontSize: 24, fontWeight: 'bold' }}>{value}</div>
      <div style={{ fontSize: 12, color: '#666' }}>{title}</div>
    </div>
  );
}

// 辅助函数：扁平化树
function flattenTree(nodes: OrgNode[]): OrgNode[] {
  const result: OrgNode[] = [];
  const traverse = (nodeList: OrgNode[]) => {
    nodeList.forEach((node) => {
      result.push(node);
      if (node.children) {
        traverse(node.children);
      }
    });
  };
  traverse(nodes);
  return result;
}
