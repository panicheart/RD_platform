import { Outlet } from 'react-router-dom';
import { Layout, Menu, Avatar, Badge, Dropdown, Button, Typography, Space } from 'antd';
import {
  HomeOutlined,
  DashboardOutlined,
  ProjectOutlined,
  AppstoreOutlined,
  BookOutlined,
  MessageOutlined,
  UserOutlined,
  BellOutlined,
  LogoutOutlined,
  SettingOutlined,
  DownOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const { Header, Content, Footer } = Layout;
const { Title, Text } = Typography;

interface NavItem {
  key: string;
  icon: React.ReactNode;
  label: string;
  path: string;
}

const navItems: NavItem[] = [
  { key: 'portal', icon: <HomeOutlined />, label: '部门门户', path: '/portal' },
  { key: 'workbench', icon: <DashboardOutlined />, label: '工作台', path: '/workbench' },
  { key: 'projects', icon: <ProjectOutlined />, label: '项目管理', path: '/projects' },
  { key: 'shelf', icon: <AppstoreOutlined />, label: '产品货架', path: '/shelf' },
  { key: 'knowledge', icon: <BookOutlined />, label: '知识库', path: '/knowledge' },
  { key: 'forum', icon: <MessageOutlined />, label: '技术论坛', path: '/forum' },
];

export function PortalLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();

  const handleNavClick = (key: string) => {
    const item = navItems.find((i) => i.key === key);
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

  const selectedKeys = navItems
    .filter((item) => location.pathname.startsWith(item.path))
    .map((item) => item.key);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {/* 顶部导航栏 */}
      <Header
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          zIndex: 1000,
          background: '#fff',
          boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
          padding: '0 48px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          height: 64,
        }}
      >
        {/* Logo区域 */}
        <Space size={16}>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <div
              style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <span style={{ color: 'white', fontWeight: 'bold', fontSize: 16 }}>M</span>
            </div>
            <Title level={4} style={{ margin: '0 0 0 12px', color: '#1890ff' }}>
              微波室研发管理平台
            </Title>
          </div>
        </Space>

        {/* 导航菜单 */}
        <Menu
          mode="horizontal"
          selectedKeys={selectedKeys}
          items={navItems.map((item) => ({
            key: item.key,
            icon: item.icon,
            label: item.label,
          }))}
          onClick={({ key }) => handleNavClick(key)}
          style={{
            flex: 1,
            justifyContent: 'center',
            borderBottom: 'none',
            background: 'transparent',
            maxWidth: 700,
          }}
        />

        {/* 右侧功能区 */}
        <Space size={20}>
          <Button
            type="text"
            icon={<SearchOutlined style={{ fontSize: 18 }} />}
            onClick={() => navigate('/search')}
          />
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
                size="default"
              />
              <Text>{user?.display_name || user?.username || '用户'}</Text>
              <DownOutlined style={{ fontSize: 12 }} />
            </Space>
          </Dropdown>
        </Space>
      </Header>

      {/* 主内容区域 - 全宽设计 */}
      <Content style={{ marginTop: 64 }}>
        <Outlet />
      </Content>

      {/* 页脚 */}
      <Footer
        style={{
          background: '#001529',
          color: '#fff',
          padding: '48px 48px 24px',
        }}
      >
        <div
          style={{
            maxWidth: 1200,
            margin: '0 auto',
            display: 'flex',
            justifyContent: 'space-between',
            flexWrap: 'wrap',
            gap: 24,
          }}
        >
          <div>
            <Title level={4} style={{ color: '#fff', marginBottom: 16 }}>
              微波室研发管理平台
            </Title>
            <Text style={{ color: 'rgba(255,255,255,0.65)' }}>
              统筹项目全生命周期 · 沉淀技术资产 · 系统化知识管理
            </Text>
          </div>
          <div style={{ textAlign: 'right' }}>
            <Text style={{ color: 'rgba(255,255,255,0.45)', fontSize: 12 }}>
              © 2026 微波室 版权所有
            </Text>
            <br />
            <Text style={{ color: 'rgba(255,255,255,0.45)', fontSize: 12 }}>
              RDP - R&D Platform for Microwave Engineering Department
            </Text>
          </div>
        </div>
      </Footer>
    </Layout>
  );
}

export default PortalLayout;
