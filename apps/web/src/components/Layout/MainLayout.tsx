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
import type { MenuProps } from 'antd';
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
  TeamOutlined,
  PartitionOutlined,
} from '@ant-design/icons';
import { useAuth } from '@/hooks/useAuth';

const { Header, Sider, Content, Footer } = Layout;
const { Text } = Typography;

type MenuItem = Required<MenuProps>['items'][number];

export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const [openKeys, setOpenKeys] = useState<string[]>(['user-management']);
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const userDropdownItems = [
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

  // 侧边栏菜单项
  const getMenuItems = (): MenuItem[] => [
    {
      key: 'portal',
      icon: <HomeOutlined />,
      label: '部门门户',
      onClick: () => navigate('/portal'),
    },
    {
      key: 'workbench',
      icon: <DashboardOutlined />,
      label: '个人工作台',
      onClick: () => navigate('/workbench'),
    },
    {
      key: 'projects',
      icon: <ProjectOutlined />,
      label: '项目管理',
      onClick: () => navigate('/projects'),
    },
    {
      key: 'user-management',
      icon: <TeamOutlined />,
      label: '用户管理',
      children: [
        {
          key: 'users',
          icon: <TeamOutlined />,
          label: '用户列表',
          onClick: () => navigate('/users'),
        },
        {
          key: 'organization',
          icon: <PartitionOutlined />,
          label: '组织架构',
          onClick: () => navigate('/organization'),
        },
      ],
    },
    {
      key: 'shelf',
      icon: <AppstoreOutlined />,
      label: '产品货架',
      onClick: () => navigate('/shelf'),
    },
    {
      key: 'knowledge',
      icon: <BookOutlined />,
      label: '知识库',
      onClick: () => navigate('/knowledge'),
    },
    {
      key: 'forum',
      icon: <MessageOutlined />,
      label: '技术论坛',
      onClick: () => navigate('/forum'),
    },
  ];

  // 获取当前选中的菜单项
  const getSelectedKeys = (): string[] => {
    const pathname = location.pathname;
    if (pathname.startsWith('/users')) return ['users'];
    if (pathname.startsWith('/organization')) return ['organization'];
    if (pathname.startsWith('/projects')) return ['projects'];
    if (pathname.startsWith('/workbench')) return ['workbench'];
    if (pathname.startsWith('/shelf')) return ['shelf'];
    if (pathname.startsWith('/knowledge')) return ['knowledge'];
    if (pathname.startsWith('/forum')) return ['forum'];
    if (pathname.startsWith('/portal')) return ['portal'];
    return [];
  };

  // 面包屑标题映射
  const getBreadcrumbTitle = (): string => {
    const pathname = location.pathname;
    if (pathname.startsWith('/users')) return '用户列表';
    if (pathname.startsWith('/organization')) return '组织架构';
    if (pathname.startsWith('/projects')) return '项目管理';
    if (pathname.startsWith('/workbench')) return '个人工作台';
    if (pathname.startsWith('/shelf')) return '产品货架';
    if (pathname.startsWith('/knowledge')) return '知识库';
    if (pathname.startsWith('/forum')) return '技术论坛';
    if (pathname.startsWith('/profile')) return '个人资料';
    if (pathname.startsWith('/settings')) return '系统设置';
    if (pathname.startsWith('/notifications')) return '消息通知';
    return '当前页面';
  };

  const selectedKeys = getSelectedKeys();

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
          openKeys={collapsed ? [] : openKeys}
          onOpenChange={setOpenKeys}
          items={getMenuItems()}
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
                { title: getBreadcrumbTitle() },
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
            <Dropdown menu={{ items: userDropdownItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar
                  src={user?.avatar || user?.avatar_url}
                  icon={!user?.avatar && !user?.avatar_url && <UserOutlined />}
                  size="small"
                />
                <Text>{user?.display_name || user?.displayName || user?.username || '用户'}</Text>
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
            微波室研发管理平台 (RDP) ©2026 版权所有
          </Text>
        </Footer>
      </Layout>
    </Layout>
  );
}

export default MainLayout;
