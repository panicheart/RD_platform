import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Spin, Result, Button } from 'antd';
import { MainLayout, PortalLayout } from '@/components/Layout';
import { useAuth } from '@/hooks/useAuth';
import type { Role } from '@/types';

// 懒加载页面组件
const PortalPage = lazy(() => import('./pages/portal/PortalPage'));
const WorkbenchPage = lazy(() => import('./pages/workbench/WorkbenchPage'));
const ProjectsPage = lazy(() => import('./pages/projects/ProjectsPage'));
const ProjectDetailPage = lazy(() => import('./pages/projects/ProjectDetailPage'));
const LoginPage = lazy(() => import('./pages/auth/LoginPage'));

// 用户管理页面
const UserList = lazy(() => import('./pages/users/UserList'));
const UserDetail = lazy(() => import('./pages/users/UserDetail'));
const OrgChart = lazy(() => import('./pages/users/OrgChart'));
const Profile = lazy(() => import('./pages/users/Profile'));
const Login = lazy(() => import('./pages/auth/Login'));
const Register = lazy(() => import('./pages/auth/Register'));

// 加载中组件
const PageLoading = () => (
  <div style={{ 
    display: 'flex', 
    justifyContent: 'center', 
    alignItems: 'center', 
    height: '100vh' 
  }}>
    <Spin size="large" tip="页面加载中..." />
  </div>
);

// 路由守卫组件
interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAdmin?: boolean;
}

const ProtectedRoute = ({ children, requireAdmin = false }: ProtectedRouteProps) => {
  const { isAuthenticated, user } = useAuth();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  if (requireAdmin) {
    const isAdmin = user?.role === 'admin' || user?.roles?.some((r: Role) => r.code === 'admin');
    if (!isAdmin) {
      return (
        <Result
          status="403"
          title="403"
          subTitle="抱歉，您没有权限访问此页面"
          extra={<Button type="primary" href="/">返回首页</Button>}
        />
      );
    }
  }
  
  return <>{children}</>;
};

function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<PageLoading />}>
        <Routes>
          {/* 登录页面 - 无需认证 */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/login-new" element={<Login />} />
          
          {/* 注册页面 - 需要管理员权限 */}
          <Route
            path="/register"
            element={
              <ProtectedRoute requireAdmin>
                <Register />
              </ProtectedRoute>
            }
          />
          
          {/* 门户页面 - 使用 PortalLayout，全宽设计 */}
          <Route
            path="/"
            element={
              <PortalLayout />
            }
          >
            <Route index element={<Navigate to="/portal" replace />} />
            <Route path="portal" element={<PortalPage />} />
          </Route>
          
          {/* 业务页面 - 使用 MainLayout，侧边菜单 */}
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route path="workbench" element={<WorkbenchPage />} />
            
            {/* 项目管理 */}
            <Route path="projects" element={<ProjectsPage />} />
            <Route path="projects/:id" element={<ProjectDetailPage />} />
            
            {/* 用户管理 */}
            <Route path="users" element={<UserList />} />
            <Route path="users/:id" element={<UserDetail />} />
            <Route path="users/:id/edit" element={<UserDetail />} />
            <Route path="organization" element={<OrgChart />} />
            
            {/* 产品货架 */}
            <Route path="shelf" element={<div>产品货架页面 - 待开发</div>} />
            
            {/* 知识库 */}
            <Route path="knowledge" element={<div>知识库页面 - 待开发</div>} />
            
            {/* 技术论坛 */}
            <Route path="forum" element={<div>技术论坛页面 - 待开发</div>} />
            
            {/* 个人资料 */}
            <Route path="profile" element={<Profile />} />
            
            {/* 系统设置 */}
            <Route path="settings" element={<div>系统设置页面 - 待开发</div>} />
            
            {/* 消息通知 */}
            <Route path="notifications" element={<div>消息通知页面 - 待开发</div>} />
          </Route>
          
          {/* 404页面 */}
          <Route
            path="*"
            element={
              <Result
                status="404"
                title="404"
                subTitle="抱歉，您访问的页面不存在"
                extra={<Button type="primary" href="/">返回首页</Button>}
              />
            }
          />
        </Routes>
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
