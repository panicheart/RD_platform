import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { Form, Input, Button, Card, Checkbox, message, Typography, Divider } from 'antd'
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons'
import { useAuth } from '@/hooks/useAuth'
import { getErrorMessage } from '@/utils/api'

const { Title, Text } = Typography

export default function LoginPage() {
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const from = (location.state as { from?: string })?.from || '/workbench'

  const handleSubmit = async (values: { username: string; password: string; remember: boolean }) => {
    setLoading(true)
    try {
      await login(values.username, values.password)
      if (values.remember) {
        localStorage.setItem('rdp_remember_user', values.username)
      } else {
        localStorage.removeItem('rdp_remember_user')
      }
      message.success('登录成功')
      navigate(from)
    } catch (error) {
      message.error(getErrorMessage(error))
    } finally {
      setLoading(false)
    }
  }

  const handleGoHome = () => {
    navigate('/portal')
  }

  const rememberedUser = localStorage.getItem('rdp_remember_user') || ''

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #0f172a 0%, #1e3a5f 50%, #0f172a 100%)',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      {/* Background decoration */}
      <div style={{
        position: 'absolute',
        top: '15%',
        left: '10%',
        fontSize: '160px',
        opacity: 0.03,
        color: 'white',
        fontFamily: 'serif',
        transform: 'rotate(-12deg)',
        pointerEvents: 'none',
      }}>
        ∇ × E
      </div>
      <div style={{
        position: 'absolute',
        bottom: '10%',
        right: '10%',
        fontSize: '120px',
        opacity: 0.03,
        color: 'white',
        fontFamily: 'serif',
        transform: 'rotate(8deg)',
        pointerEvents: 'none',
      }}>
        ∇ × B
      </div>

      <Card
        style={{
          width: 420,
          borderRadius: 16,
          boxShadow: '0 20px 60px rgba(0, 0, 0, 0.3)',
          border: 'none',
        }}
        styles={{ body: { padding: '40px 36px 32px' } }}
      >
        {/* Logo */}
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <div style={{
            width: 56,
            height: 56,
            borderRadius: '50%',
            background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            marginBottom: 16,
          }}>
            <SafetyOutlined style={{ color: 'white', fontSize: 28 }} />
          </div>
          <Title level={3} style={{ marginBottom: 4 }}>
            微波室研发管理平台
          </Title>
          <Text type="secondary" style={{ fontSize: 13 }}>Microwave R&D Management Platform</Text>
        </div>

        <Form
          name="login"
          onFinish={handleSubmit}
          autoComplete="off"
          size="large"
          initialValues={{ remember: !!rememberedUser, username: rememberedUser }}
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined style={{ color: '#bfbfbf' }} />}
              placeholder="用户名"
              style={{ borderRadius: 8 }}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6位' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#bfbfbf' }} />}
              placeholder="密码"
              style={{ borderRadius: 8 }}
            />
          </Form.Item>

          <Form.Item name="remember" valuePropName="checked" style={{ marginBottom: 16 }}>
            <Checkbox>记住用户名</Checkbox>
          </Form.Item>

          <Form.Item style={{ marginBottom: 16 }}>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              style={{
                height: 44,
                borderRadius: 8,
                fontSize: 16,
                background: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
                border: 'none',
              }}
            >
              登录
            </Button>
          </Form.Item>
        </Form>

        <div style={{ textAlign: 'center' }}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            默认账号: admin / Admin@123
          </Text>
        </div>

        <Divider style={{ margin: '16px 0' }} />

        <div style={{ textAlign: 'center' }}>
          <a onClick={handleGoHome} style={{ cursor: 'pointer', color: '#3b82f6', fontSize: 13 }}>
            返回首页
          </a>
        </div>
      </Card>
    </div>
  )
}
