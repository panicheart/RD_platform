import { useState } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  Layout,
  Menu,
  Avatar,
  Badge,
  Dropdown,
  Button,
  theme,
  Breadcrumb,
  Space,
  Typography,
} from 'antd';
import {
  HomeOutlined,
  DashboardOutlined,
  ProjectOutlined,
  AppstoreOutlined,
  BookOutlined,
  MessageOutlined,
  SettingOutlined,
  UserOutlined,
  BellOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  DownOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';

const { Header, Sider, Content, Footer } = Layout;
const { Text } = Typography;

interface MenuItem {
  key: string;
  icon: React.ReactNode;
  label: string;
  path: string;
}

const menuItems: MenuItem[] = [
  { key: 'portal', icon: <HomeOutlined />, label: '部门门户', path: '/portal' },
  { key: 'workbench', icon: <DashboardOutlined />, label: '个人工作台', path: '/workbench' },
  { key: 'projects', icon: <ProjectOutlined />, label: '项目管理', path: '/projects' },
  { key: 'shelf', icon: <AppstoreOutlined />, label: '产品货架', path: '/shelf' },
  { key: 'knowledge', icon: <BookOutlined />, label: '知识库', path: '/knowledge' },
  { key: 'forum', icon: <MessageOutlined />, label: '技术论坛', path: '/forum' },
];

export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const handleMenuClick = (key: string) => {
    const item = menuItems.find((i) => i.key === key);
    if (item) {
      navigate(item.path);
    }
  };

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '系统设置',
      onClick: () => navigate('/settings'),
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
  ];

  const selectedKeys = menuItems
    .filter((item) => location.pathname.startsWith(item.path))
    .map((item) => item.key);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        style={{
          boxShadow: '2px 0 8px rgba(0,0,0,0.06)',
          zIndex: 10,
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          {collapsed ? (
            <span style={{ fontSize: 20, fontWeight: 600, color: '#1890ff' }}>RDP</span>
          ) : (
            <span style={{ fontSize: 16, fontWeight: 600 }}>研发管理平台</span>
          )}
        </div>
        <Menu
          mode="inline"
          selectedKeys={selectedKeys}
          items={menuItems.map((item) => ({
            key: item.key,
            icon: item.icon,
            label: item.label,
          }))}
          onClick={({ key }) => handleMenuClick(key)}
          style={{ borderRight: 0 }}
        />
      </Sider>
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'margin-left 0.2s' }}>
        <Header
          style={{
            padding: '0 24px',
            background: colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
            zIndex: 9,
            position: 'sticky',
            top: 0,
          }}
        >
          <Space>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{ fontSize: 16 }}
            />
            <Breadcrumb
              items={[
                { title: '首页' },
                {
                  title:
                    menuItems.find((i) => location.pathname.startsWith(i.path))
                      ?.label || '当前页面',
                },
              ]}
            />
          </Space>
          <Space size={24}>
            <Badge count={5} size="small">
              <Button
                type="text"
                icon={<BellOutlined style={{ fontSize: 18 }} />}
                onClick={() => navigate('/notifications')}
              />
            </Badge>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar
                  src={user?.avatar_url}
                  icon={!user?.avatar_url && <UserOutlined />}
                  size="small"
                />
                <Text>{user?.display_name || user?.username || '用户'}</Text>
                <DownOutlined style={{ fontSize: 12 }} />
              </Space>
            </Dropdown>
          </Space>
        </Header>
        <Content
          style={{
            margin: 24,
            padding: 24,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
            minHeight: 280,
            overflow: 'auto',
          }}
        >
          <Outlet />
        </Content>
        <Footer style={{ textAlign: 'center', padding: '12px 50px' }}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            微波工程部研发管理平台 (RDP) ©2026 版权所有
          </Text>
        </Footer>
      </Layout>
    </Layout>
  );
}

export default MainLayout;
