import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Descriptions,
  Tag,
  Button,
  Space,
  Avatar,
  Typography,
  Skeleton,
  message,
  Tabs,
  Timeline,
  Empty,
  Row,
  Col,
  Form,
  Input,
  Select,
  Switch,
  Divider,
} from 'antd';
import {
  ArrowLeftOutlined,
  EditOutlined,
  UserOutlined,
  MailOutlined,
  PhoneOutlined,
  TeamOutlined,
  SafetyOutlined,
  HistoryOutlined,
  SaveOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import { userAPI } from '@/services/user';
import type { User, Role } from '@/types';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;
const { TextArea } = Input;

interface UserActivity {
  id: string;
  action: string;
  resource: string;
  createdAt: string;
  details?: string;
}

export default function UserDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [activities, setActivities] = useState<UserActivity[]>([]);
  const [form] = Form.useForm();

  // 权限检查
  const isAdmin = currentUser?.role === 'admin' || currentUser?.roles?.some((r: Role) => r.code === 'admin');
  const isSelf = currentUser?.id === id;
  const canEdit = isAdmin || isSelf;

  useEffect(() => {
    if (id) {
      fetchUserDetail();
      fetchUserActivities();
    }
  }, [id]);

  const fetchUserDetail = async () => {
    setLoading(true);
    try {
      const data = await userAPI.getUser(id!);
      setUser(data);
      form.setFieldsValue({
        ...data,
        createdAt: data.createdAt ? dayjs(data.createdAt) : undefined,
      });
    } catch (error: any) {
      message.error(error.message || '获取用户信息失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchUserActivities = async () => {
    try {
      // 实际应从 API 获取用户活动日志
      const mockActivities: UserActivity[] = [
        {
          id: '1',
          action: '登录',
          resource: '系统',
          createdAt: new Date().toISOString(),
          details: 'IP: 192.168.1.100',
        },
        {
          id: '2',
          action: '更新',
          resource: '个人资料',
          createdAt: new Date(Date.now() - 86400000).toISOString(),
        },
      ];
      setActivities(mockActivities);
    } catch (error) {
      // 忽略错误
    }
  };

  const handleSave = async (values: any) => {
    if (!canEdit) {
      message.error('您没有权限编辑此用户');
      return;
    }

    try {
      await userAPI.updateUser(id!, {
        ...values,
        createdAt: values.createdAt?.toISOString(),
      });
      message.success('用户信息已更新');
      setEditing(false);
      fetchUserDetail();
    } catch (error: any) {
      message.error(error.message || '更新用户信息失败');
    }
  };

  const getRoleTag = (role: string) => {
    const roleMap: Record<string, { color: string; text: string }> = {
      admin: { color: 'red', text: '管理员' },
      manager: { color: 'blue', text: '项目经理' },
      leader: { color: 'green', text: '负责人' },
      designer: { color: 'default', text: '工程师' },
      viewer: { color: 'default', text: '访客' },
    };
    const config = roleMap[role] || { color: 'default', text: role };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      active: { color: 'success', text: '正常' },
      inactive: { color: 'default', text: '停用' },
      locked: { color: 'error', text: '锁定' },
      pending: { color: 'warning', text: '待审核' },
    };
    const config = statusMap[status] || { color: 'default', text: status };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  if (loading) {
    return (
      <Card>
        <Skeleton active avatar paragraph={{ rows: 6 }} />
      </Card>
    );
  }

  if (!user) {
    return (
      <Card>
        <Empty description="用户不存在或已被删除" />
        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Button onClick={() => navigate('/users')}>返回用户列表</Button>
        </div>
      </Card>
    );
  }

  return (
    <div>
      {/* 页面头部 */}
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/users')}
            style={{ marginRight: 16 }}
          >
            返回
          </Button>
          <Title level={2} style={{ display: 'inline', margin: 0 }}>
            用户详情
          </Title>
        </Col>
        <Col>
          {canEdit && !editing && (
            <Button
              type="primary"
              icon={<EditOutlined />}
              onClick={() => setEditing(true)}
            >
              编辑
            </Button>
          )}
          {editing && (
            <Space>
              <Button icon={<CloseOutlined />} onClick={() => setEditing(false)}>
                取消
              </Button>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={() => form.submit()}
              >
                保存
              </Button>
            </Space>
          )}
        </Col>
      </Row>

      <Tabs defaultActiveKey="basic">
        <TabPane
          tab={
            <span>
              <UserOutlined />
              基本信息
            </span>
          }
          key="basic"
        >
          <Row gutter={24}>
            {/* 左侧：头像和快速信息 */}
            <Col span={6}>
              <Card style={{ textAlign: 'center' }}>
                <Avatar
                  src={user.avatar}
                  icon={!user.avatar && <UserOutlined />}
                  size={120}
                  style={{ marginBottom: 16 }}
                />
                <Title level={4} style={{ margin: 0 }}>{user.displayName}</Title>
                <Text type="secondary">@{user.username}</Text>
                <div style={{ marginTop: 16 }}>
                  {getRoleTag(user.roles?.[0]?.code || 'designer')}
                  {getStatusTag(user.status)}
                </div>
              </Card>
            </Col>

            {/* 右侧：详细信息 */}
            <Col span={18}>
              <Card>
                {editing ? (
                  <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSave}
                    initialValues={user}
                  >
                    <Row gutter={16}>
                      <Col span={12}>
                        <Form.Item
                          name="username"
                          label="用户名"
                          rules={[{ required: true }]}
                        >
                          <Input disabled />
                        </Form.Item>
                      </Col>
                      <Col span={12}>
                        <Form.Item
                          name="displayName"
                          label="显示名称"
                          rules={[{ required: true }]}
                        >
                          <Input />
                        </Form.Item>
                      </Col>
                    </Row>

                    <Row gutter={16}>
                      <Col span={12}>
                        <Form.Item
                          name="email"
                          label="邮箱"
                          rules={[{ type: 'email' }]}
                        >
                          <Input prefix={<MailOutlined />} />
                        </Form.Item>
                      </Col>
                      <Col span={12}>
                        <Form.Item name="phone" label="电话">
                          <Input prefix={<PhoneOutlined />} />
                        </Form.Item>
                      </Col>
                    </Row>

                    {isAdmin && (
                      <>
                        <Divider />
                        <Row gutter={16}>
                          <Col span={12}>
                            <Form.Item name="role" label="角色">
                              <Select>
                                <Option value="admin">管理员</Option>
                                <Option value="manager">项目经理</Option>
                                <Option value="leader">负责人</Option>
                                <Option value="designer">工程师</Option>
                                <Option value="viewer">访客</Option>
                              </Select>
                            </Form.Item>
                          </Col>
                          <Col span={12}>
                            <Form.Item name="status" label="状态">
                              <Select>
                                <Option value="active">正常</Option>
                                <Option value="inactive">停用</Option>
                                <Option value="locked">锁定</Option>
                                <Option value="pending">待审核</Option>
                              </Select>
                            </Form.Item>
                          </Col>
                        </Row>
                        <Row gutter={16}>
                          <Col span={24}>
                            <Form.Item name="isActive" label="账号启用" valuePropName="checked">
                              <Switch />
                            </Form.Item>
                          </Col>
                        </Row>
                      </>
                    )}

                    <Form.Item name="bio" label="个人简介">
                      <TextArea rows={4} placeholder="请输入个人简介" />
                    </Form.Item>
                  </Form>
                ) : (
                  <Descriptions bordered column={2}>
                    <Descriptions.Item label="用户名">
                      {user.username}
                    </Descriptions.Item>
                    <Descriptions.Item label="显示名称">
                      {user.displayName}
                    </Descriptions.Item>
                    <Descriptions.Item label="邮箱">
                      <Space>
                        <MailOutlined />
                        {user.email || '-'}
                      </Space>
                    </Descriptions.Item>
                    <Descriptions.Item label="电话">
                      <Space>
                        <PhoneOutlined />
                        {user.phone || '-'}
                      </Space>
                    </Descriptions.Item>
                    <Descriptions.Item label="角色">
                      {getRoleTag(user.roles?.[0]?.code || 'designer')}
                    </Descriptions.Item>
                    <Descriptions.Item label="状态">
                      {getStatusTag(user.status)}
                    </Descriptions.Item>
                    <Descriptions.Item label="部门" span={2}>
                      <Space>
                        <TeamOutlined />
                        {user.organization?.name || '-'}
                      </Space>
                    </Descriptions.Item>
                    <Descriptions.Item label="创建时间">
                      {user.createdAt ? new Date(user.createdAt).toLocaleString('zh-CN') : '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="更新时间">
                      {user.updatedAt ? new Date(user.updatedAt).toLocaleString('zh-CN') : '-'}
                    </Descriptions.Item>
                  </Descriptions>
                )}
              </Card>
            </Col>
          </Row>
        </TabPane>

        <TabPane
          tab={
            <span>
              <SafetyOutlined />
              权限信息
            </span>
          }
          key="permissions"
        >
          <Card>
            <Title level={5}>角色权限</Title>
            {user.roles?.map((role) => (
              <Card key={role.id} size="small" style={{ marginBottom: 16 }}>
                <div style={{ marginBottom: 8 }}>
                  <Text strong>{role.name}</Text>
                  <Tag style={{ marginLeft: 8 }}>{role.code}</Tag>
                </div>
                <div>
                  {role.permissions?.map((perm) => (
                    <Tag key={perm} color="blue" style={{ margin: 4 }}>
                      {perm}
                    </Tag>
                  ))}
                </div>
              </Card>
            )) || <Empty description="暂无权限信息" />}
          </Card>
        </TabPane>

        <TabPane
          tab={
            <span>
              <HistoryOutlined />
              操作记录
            </span>
          }
          key="activities"
        >
          <Card>
            <Timeline mode="left">
              {activities.map((activity) => (
                <Timeline.Item key={activity.id}>
                  <p>
                    <Text strong>{activity.action}</Text>
                    <Text type="secondary" style={{ marginLeft: 8 }}>
                      {activity.resource}
                    </Text>
                  </p>
                  <p>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {new Date(activity.createdAt).toLocaleString('zh-CN')}
                    </Text>
                  </p>
                  {activity.details && (
                    <p>
                      <Text type="secondary">{activity.details}</Text>
                    </p>
                  )}
                </Timeline.Item>
              ))}
            </Timeline>
          </Card>
        </TabPane>
      </Tabs>
    </div>
  );
}
