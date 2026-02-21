import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, message, Typography, Select, Space } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined, ArrowLeftOutlined } from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';
import { userAPI } from '@/services/user';

const { Title, Text } = Typography;
const { Option } = Select;

interface RegisterFormValues {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  displayName: string;
  role: string;
}

const roleOptions = [
  { value: 'admin', label: '管理员', color: 'red' },
  { value: 'manager', label: '项目经理', color: 'blue' },
  { value: 'leader', label: '技术/项目负责人', color: 'green' },
  { value: 'designer', label: '设计师/工程师', color: 'default' },
  { value: 'viewer', label: '访客', color: 'default' },
];

// 仅管理员可访问
export default function Register() {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { user } = useAuth();

  // 权限检查 - 只有管理员可以创建用户
  const isAdmin = user?.role === 'admin';

  const handleSubmit = async (values: RegisterFormValues) => {
    if (!isAdmin) {
      message.error('您没有权限创建用户');
      return;
    }

    setLoading(true);
    try {
      await userAPI.createUser({
        username: values.username,
        email: values.email,
        password: values.password,
        display_name: values.displayName,
        role: values.role,
      });

      message.success('用户创建成功');
      form.resetFields();
      
      // 可选择返回用户列表
      setTimeout(() => {
        navigate('/users');
      }, 1500);
    } catch (error: any) {
      message.error(error.message || '创建用户失败');
    } finally {
      setLoading(false);
    }
  };

  const handleGoBack = () => {
    navigate(-1);
  };

  // 密码确认验证
  const validateConfirmPassword = (_: any, value: string) => {
    const password = form.getFieldValue('password');
    if (value && value !== password) {
      return Promise.reject(new Error('两次输入的密码不一致'));
    }
    return Promise.resolve();
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#f0f2f5',
      }}
    >
      <Card
        style={{ width: 480, boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)' }}
        bordered={false}
      >
        <div style={{ marginBottom: 24 }}>
          <Button 
            type="link" 
            icon={<ArrowLeftOutlined />} 
            onClick={handleGoBack}
            style={{ padding: 0 }}
          >
            返回
          </Button>
        </div>

        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Title level={3} style={{ marginBottom: 8 }}>
            创建新用户
          </Title>
          <Text type="secondary">管理员专用 - 创建系统用户账号</Text>
          {!isAdmin && (
            <div style={{ marginTop: 8 }}>
              <Text type="danger" style={{ fontSize: 12 }}>
                警告：您当前不是管理员角色
              </Text>
            </div>
          )}
        </div>

        <Form
          form={form}
          name="register"
          onFinish={handleSubmit}
          autoComplete="off"
          size="large"
          layout="vertical"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' },
              { max: 20, message: '用户名最多20个字符' },
              { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="请输入用户名"
              disabled={!isAdmin}
            />
          </Form.Item>

          <Form.Item
            name="displayName"
            label="显示名称"
            rules={[
              { required: true, message: '请输入显示名称' },
              { max: 50, message: '显示名称最多50个字符' },
            ]}
          >
            <Input placeholder="请输入显示名称，如：张三" disabled={!isAdmin} />
          </Form.Item>

          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input
              prefix={<MailOutlined />}
              placeholder="请输入邮箱"
              disabled={!isAdmin}
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6个字符' },
              { max: 32, message: '密码最多32个字符' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入密码"
              disabled={!isAdmin}
            />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            label="确认密码"
            rules={[
              { required: true, message: '请确认密码' },
              { validator: validateConfirmPassword },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请再次输入密码"
              disabled={!isAdmin}
            />
          </Form.Item>

          <Form.Item
            name="role"
            label="用户角色"
            rules={[{ required: true, message: '请选择用户角色' }]}
          >
            <Select placeholder="请选择角色" disabled={!isAdmin}>
              {roleOptions.map((role) => (
                <Option key={role.value} value={role.value}>
                  {role.label}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item style={{ marginTop: 32 }}>
            <Space style={{ width: '100%' }}>
              <Button onClick={handleGoBack} size="large" style={{ width: 120 }}>
                取消
              </Button>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                size="large"
                style={{ width: 320 }}
                disabled={!isAdmin}
              >
                创建用户
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
