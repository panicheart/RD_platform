import { useState, useMemo } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Badge, Space, Breadcrumb, theme } from 'antd'
import {
  HomeOutlined,
  ProjectOutlined,
  AppstoreOutlined,
  TeamOutlined,
  UserOutlined,
  BellOutlined,
  SettingOutlined,
  LogoutOutlined,
  BookOutlined,
  MessageOutlined,
  BarChartOutlined,
  MonitorOutlined,
  ShoppingOutlined,
} from '@ant-design/icons'
import { useAuth } from '@/hooks/useAuth'

const { Header, Sider, Content } = Layout

const menuItems = [
  {
    key: '/portal',
    icon: <HomeOutlined />,
    label: '门户首页',
  },
  {
    key: '/workbench',
    icon: <AppstoreOutlined />,
    label: '个人工作台',
  },
  {
    key: '/projects',
    icon: <ProjectOutlined />,
    label: '项目管理',
  },
  {
    key: '/users',
    icon: <TeamOutlined />,
    label: '用户管理',
  },
  {
    key: '/shelf',
    icon: <ShoppingOutlined />,
    label: '产品货架',
  },
  {
    key: '/knowledge',
    icon: <BookOutlined />,
    label: '知识库',
  },
  {
    key: '/forum',
    icon: <MessageOutlined />,
    label: '技术论坛',
  },
  {
    key: '/analytics',
    icon: <BarChartOutlined />,
    label: '数据分析',
  },
  {
    key: '/monitor',
    icon: <MonitorOutlined />,
    label: '系统监控',
  },
]

const breadcrumbNameMap: Record<string, string> = {
  '/workbench': '个人工作台',
  '/projects': '项目管理',
  '/users': '用户管理',
  '/portal': '门户首页',
  '/shelf': '产品货架',
  '/knowledge': '知识库',
  '/forum': '技术论坛',
  '/analytics': '数据分析',
  '/monitor': '系统监控',
}

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const { token } = theme.useToken()

  const breadcrumbItems = useMemo(() => {
    const pathSnippets = location.pathname.split('/').filter((i) => i)
    const items = [
      {
        title: <a onClick={() => navigate('/workbench')}>首页</a>,
      },
    ]

    pathSnippets.forEach((_, index) => {
      const url = `/${pathSnippets.slice(0, index + 1).join('/')}`
      const name = breadcrumbNameMap[url]
      if (name) {
        items.push({
          title: index === pathSnippets.length - 1
            ? <span>{name}</span>
            : <a onClick={() => navigate(url)}>{name}</a>,
        })
      } else if (index > 0) {
        items.push({ title: <span>详情</span> })
      }
    })

    return items
  }, [location.pathname, navigate])

  const handleMenuClick = (key: string) => {
    navigate(key)
  }

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const selectedKeys = useMemo(() => {
    const match = location.pathname.match(/^\/[^/]+/)
    return match ? [match[0]] : []
  }, [location.pathname])

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '系统设置',
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
      onClick: handleLogout,
    },
  ]

  return (
    <Layout className="main-layout">
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        theme="light"
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          zIndex: 100,
          borderRight: `1px solid ${token.colorBorderSecondary}`,
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
          }}
        >
          {collapsed ? (
            <div style={{
              width: 36,
              height: 36,
              borderRadius: '50%',
              background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}>
              <span style={{ color: 'white', fontWeight: 'bold', fontSize: 16 }}>M</span>
            </div>
          ) : (
            <Space>
              <div style={{
                width: 32,
                height: 32,
                borderRadius: '50%',
                background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <span style={{ color: 'white', fontWeight: 'bold', fontSize: 14 }}>M</span>
              </div>
              <span style={{ fontSize: 16, fontWeight: 600 }}>研发管理平台</span>
            </Space>
          )}
        </div>
        <Menu
          mode="inline"
          selectedKeys={selectedKeys}
          items={menuItems}
          onClick={({ key }) => handleMenuClick(key)}
          style={{ borderRight: 0 }}
        />
      </Sider>
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'margin-left 0.2s' }}>
        <Header
          className="main-header"
          style={{
            padding: '0 24px',
            background: token.colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            position: 'sticky',
            top: 0,
            zIndex: 99,
            width: '100%',
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
          }}
        >
          <Breadcrumb items={breadcrumbItems} />
          <Space size="middle">
            <Badge count={3} size="small">
              <BellOutlined style={{ fontSize: 18, cursor: 'pointer' }} />
            </Badge>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar
                  src={user?.avatar_url}
                  icon={!user?.avatar_url && <UserOutlined />}
                  style={{ backgroundColor: token.colorPrimary }}
                />
                <span>{user?.display_name || user?.username}</span>
              </Space>
            </Dropdown>
          </Space>
        </Header>
        <Content className="main-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
