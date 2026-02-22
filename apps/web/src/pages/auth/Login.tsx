import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Form, Input, Button, Card, message, Typography, Checkbox } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores/auth';
import { userAPI } from '@/services/user';

const { Title, Text } = Typography;

interface LoginFormValues {
  username: string;
  password: string;
  remember: boolean;
}

export default function Login() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { setAuth } = useAuthStore();

  // 获取登录后的跳转目标
  const from = (location.state as { from?: { pathname?: string } })?.from?.pathname || '/workbench';

  const handleSubmit = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      const response = await userAPI.login({
        username: values.username,
        password: values.password,
      });
      
      // 存储JWT Token
      const { user, access_token, refresh_token } = response;
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);
      
      // 更新全局状态
      setAuth(user, access_token);
      
      message.success('登录成功');
      
      // 跳转
      navigate(from, { replace: true });
    } catch (error: any) {
      message.error(error.message || '登录失败，请检查用户名和密码');
    } finally {
      setLoading(false);
    }
  };

  const handleGoHome = () => {
    navigate('/portal');
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #1890ff 0%, #40a9ff 100%)',
      }}
    >
      <Card
        style={{ width: 420, boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)' }}
        bordered={false}
      >
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Title level={3} style={{ marginBottom: 8 }}>
            微波室研发管理平台
          </Title>
          <Text type="secondary">Microwave R&D Management Platform</Text>
        </div>

        <Form
          name="login"
          onFinish={handleSubmit}
          autoComplete="off"
          size="large"
          initialValues={{ remember: true }}
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
              autoFocus
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
            />
          </Form.Item>

          <Form.Item>
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>记住我</Checkbox>
            </Form.Item>
            <a style={{ float: 'right' }} onClick={() => message.info('请联系管理员重置密码')}>
              忘记密码?
            </a>
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              登录
            </Button>
          </Form.Item>
        </Form>

        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            首次使用? 请联系管理员创建账号
          </Text>
        </div>

        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Button type="link" onClick={handleGoHome}>
            返回首页
          </Button>
        </div>
      </Card>
    </div>
  );
}
