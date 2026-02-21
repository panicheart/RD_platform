import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  Avatar,
  Typography,
  Tabs,
  Form,
  Input,
  Button,
  Space,
  Row,
  Col,
  Upload,
  message,
  Divider,
  Tag,
  List,
  Modal,
  Descriptions,
} from 'antd';
import {
  UserOutlined,
  MailOutlined,
  PhoneOutlined,
  SafetyOutlined,
  EditOutlined,
  SaveOutlined,
  LockOutlined,
  HistoryOutlined,
  TagsOutlined,
  ProjectOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import { userAPI } from '@/services/user';
import type { User } from '@/types';

const { Title, Text, Paragraph } = Typography;
const { TabPane } = Tabs;
const { TextArea } = Input;

interface PasswordFormValues {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
}

export default function Profile() {
  const navigate = useNavigate();
  const { user: currentUser, refreshUser } = useAuth();
  
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [editing, setEditing] = useState(false);
  const [activeTab, setActiveTab] = useState('basic');
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [skills, setSkills] = useState<string[]>([]);
  const [newSkill, setNewSkill] = useState('');
  const [form] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    if (currentUser?.id) {
      fetchUserProfile();
    }
  }, [currentUser?.id]);

  const fetchUserProfile = async () => {
    setLoading(true);
    try {
      const data = await userAPI.getUser(currentUser!.id);
      setUser(data);
      setSkills(data.skills || []);
      form.setFieldsValue(data);
    } catch (error: any) {
      message.error(error.message || '获取用户信息失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (values: any) => {
    try {
      await userAPI.updateUser(currentUser!.id, {
        ...values,
        skills,
      });
      message.success('个人资料已更新');
      setEditing(false);
      refreshUser();
      fetchUserProfile();
    } catch (error: any) {
      message.error(error.message || '更新失败');
    }
  };

  const handleChangePassword = async (values: PasswordFormValues) => {
    try {
      await userAPI.changePassword({
        old_password: values.oldPassword,
        new_password: values.newPassword,
      });
      message.success('密码已修改，请重新登录');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
      // 可选：登出并跳转登录页
    } catch (error: any) {
      message.error(error.message || '密码修改失败');
    }
  };

  const handleAddSkill = () => {
    if (!newSkill.trim()) return;
    if (skills.includes(newSkill.trim())) {
      message.warning('该技能已存在');
      return;
    }
    setSkills([...skills, newSkill.trim()]);
    setNewSkill('');
  };

  const handleRemoveSkill = (skill: string) => {
    setSkills(skills.filter((s) => s !== skill));
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

  const uploadProps = {
    name: 'file',
    action: '/api/v1/users/avatar',
    headers: {
      authorization: `Bearer ${localStorage.getItem('access_token')}`,
    },
    showUploadList: false,
    onChange(info: any) {
      if (info.file.status === 'done') {
        message.success('头像上传成功');
        fetchUserProfile();
      } else if (info.file.status === 'error') {
        message.error('头像上传失败');
      }
    },
  };

  if (!user) {
    return (
      <Card loading={true}>
        <div style={{ height: 400 }} />
      </Card>
    );
  }

  return (
    <div>
      {/* 页面头部 */}
      <Title level={2} style={{ marginBottom: 24 }}>
        <UserOutlined style={{ marginRight: 12 }} />
        个人资料
      </Title>

      <Row gutter={24}>
        {/* 左侧：头像和基本信息 */}
        <Col span={8}>
          <Card style={{ textAlign: 'center' }} loading={loading}>
            <Upload {...uploadProps} accept="image/*">
              <div style={{ position: 'relative', display: 'inline-block' }}>
                <Avatar
                  src={user.avatar}
                  icon={!user.avatar && <UserOutlined />}
                  size={120}
                  style={{ cursor: 'pointer' }}
                />
                <div
                  style={{
                    position: 'absolute',
                    bottom: 0,
                    right: 0,
                    background: '#1890ff',
                    borderRadius: '50%',
                    padding: 8,
                    cursor: 'pointer',
                  }}
                >
                  <UploadOutlined style={{ color: '#fff', fontSize: 14 }} />
                </div>
              </div>
            </Upload>

            <Title level={4} style={{ marginTop: 16, marginBottom: 4 }}>
              {user.displayName}
            </Title>
            <Text type="secondary">@{user.username}</Text>

            <div style={{ marginTop: 16 }}>
              {getRoleTag(user.roles?.[0]?.code || 'designer')}
              {getStatusTag(user.status)}
            </div>

            <Divider />

            <Descriptions column={1} size="small">
              <Descriptions.Item label="部门">
                {user.organization?.name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="邮箱">
                {user.email || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="电话">
                {user.phone || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="入职时间">
                {user.createdAt
                  ? new Date(user.createdAt).toLocaleDateString('zh-CN')
                  : '-'}
              </Descriptions.Item>
            </Descriptions>

            <Divider />

            <Space direction="vertical" style={{ width: '100%' }}>
              <Button
                type="primary"
                icon={<EditOutlined />}
                block
                onClick={() => setEditing(true)}
              >
                编辑资料
              </Button>
              <Button
                icon={<LockOutlined />}
                block
                onClick={() => setPasswordModalVisible(true)}
              >
                修改密码
              </Button>
            </Space>
          </Card>
        </Col>

        {/* 右侧：详细信息 */}
        <Col span={16}>
          <Tabs activeKey={activeTab} onChange={setActiveTab}>
            <TabPane
              tab={
                <span>
                  <UserOutlined />
                  基本信息
                </span>
              }
              key="basic"
            >
              <Card loading={loading}>
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
                          name="displayName"
                          label="显示名称"
                          rules={[{ required: true, message: '请输入显示名称' }]}
                        >
                          <Input prefix={<UserOutlined />} />
                        </Form.Item>
                      </Col>
                      <Col span={12}>
                        <Form.Item name="phone" label="电话">
                          <Input prefix={<PhoneOutlined />} />
                        </Form.Item>
                      </Col>
                    </Row>

                    <Form.Item
                      name="email"
                      label="邮箱"
                      rules={[
                        { type: 'email', message: '请输入有效的邮箱地址' },
                      ]}
                    >
                      <Input prefix={<MailOutlined />} />
                    </Form.Item>

                    <Form.Item name="bio" label="个人简介">
                      <TextArea
                        rows={4}
                        placeholder="介绍一下你自己..."
                      />
                    </Form.Item>

                    <Form.Item>
                      <Space>
                        <Button type="primary" htmlType="submit" icon={<SaveOutlined />}>
                          保存
                        </Button>
                        <Button onClick={() => setEditing(false)}>取消</Button>
                      </Space>
                    </Form.Item>
                  </Form>
                ) : (
                  <Descriptions column={2} bordered>
                    <Descriptions.Item label="用户名" span={2}>
                      {user.username}
                    </Descriptions.Item>
                    <Descriptions.Item label="显示名称">
                      {user.displayName}
                    </Descriptions.Item>
                    <Descriptions.Item label="角色">
                      {getRoleTag(user.roles?.[0]?.code || 'designer')}
                    </Descriptions.Item>
                    <Descriptions.Item label="邮箱" span={2}>
                      {user.email || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="电话">
                      {user.phone || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="部门">
                      {user.organization?.name || '-'}
                    </Descriptions.Item>
                    <Descriptions.Item label="个人简介" span={2}>
                      <Paragraph>
                        {user.bio || '暂无个人简介'}
                      </Paragraph>
                    </Descriptions.Item>
                  </Descriptions>
                )}
              </Card>
            </TabPane>

            <TabPane
              tab={
                <span>
                  <TagsOutlined />
                  技能标签
                </span>
              }
              key="skills"
            >
              <Card>
                <Title level={5}>我的技能</Title>
                <Paragraph type="secondary">
                  添加您的专业技能，便于项目匹配和团队协作
                </Paragraph>

                <div style={{ marginBottom: 24 }}>
                  <Space wrap>
                    {skills.map((skill) => (
                      <Tag
                        key={skill}
                        closable
                        onClose={() => handleRemoveSkill(skill)}
                        style={{ fontSize: 14, padding: '4px 12px' }}
                      >
                        {skill}
                      </Tag>
                    ))}
                  </Space>
                </div>

                <Space.Compact style={{ width: '100%' }}>
                  <Input
                    placeholder="输入技能名称，如：React、微波设计、射频电路..."
                    value={newSkill}
                    onChange={(e) => setNewSkill(e.target.value)}
                    onPressEnter={handleAddSkill}
                  />
                  <Button type="primary" onClick={handleAddSkill}>
                    添加
                  </Button>
                </Space.Compact>

                <Divider />

                <Title level={5}>推荐技能</Title>
                <Space wrap>
                  {['微波设计', '射频电路', '天线设计', 'PCB设计', 'FPGA', '嵌入式', 'Python', 'Java'].map(
                    (skill) => (
                      <Tag
                        key={skill}
                        style={{ cursor: 'pointer' }}
                        onClick={() => {
                          if (!skills.includes(skill)) {
                            setSkills([...skills, skill]);
                          }
                        }}
                      >
                        + {skill}
                      </Tag>
                    )
                  )}
                </Space>
              </Card>
            </TabPane>

            <TabPane
              tab={
                <span>
                  <ProjectOutlined />
                  参与项目
                </span>
              }
              key="projects"
            >
              <Card>
                <List
                  itemLayout="horizontal"
                  dataSource={[]}
                  renderItem={(item) => (
                    <List.Item
                      actions={[
                        <Button type="link" onClick={() => navigate(`/projects/${item}`)}>
                          查看
                        </Button>,
                      ]}
                    >
                      <List.Item.Meta
                        title={<Text strong>项目名称</Text>}
                        description="项目描述..."
                      />
                    </List.Item>
                  )}
                />
                <div style={{ textAlign: 'center', padding: 40 }}>
                  <Text type="secondary">暂无项目数据</Text>
                </div>
              </Card>
            </TabPane>

            <TabPane
              tab={
                <span>
                  <HistoryOutlined />
                  登录记录
                </span>
              }
              key="history"
            >
              <Card>
                <List
                  itemLayout="horizontal"
                  dataSource={[
                    {
                      time: new Date().toLocaleString('zh-CN'),
                      ip: '192.168.1.100',
                      device: 'Chrome on Windows',
                    },
                  ]}
                  renderItem={(item) => (
                    <List.Item>
                      <List.Item.Meta
                        title={
                          <Space>
                            <Text strong>登录成功</Text>
                            <Tag color="success">当前会话</Tag>
                          </Space>
                        }
                        description={
                          <Space direction="vertical" size={0}>
                            <Text type="secondary">时间: {item.time}</Text>
                            <Text type="secondary">IP: {item.ip}</Text>
                            <Text type="secondary">设备: {item.device}</Text>
                          </Space>
                        }
                      />
                    </List.Item>
                  )}
                />
              </Card>
            </TabPane>
          </Tabs>
        </Col>
      </Row>

      {/* 修改密码弹窗 */}
      <Modal
        title={
          <Space>
            <SafetyOutlined />
            修改密码
          </Space>
        }
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false);
          passwordForm.resetFields();
        }}
        footer={null}
        width={400}
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleChangePassword}
        >
          <Form.Item
            name="oldPassword"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="当前密码" />
          </Form.Item>

          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="新密码" />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认新密码"
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="确认新密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
