import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Spin, Result, Button } from 'antd';
import { MainLayout, PortalLayout } from '@/components/Layout';
import { useAuth } from '@/hooks/useAuth';

// 懒加载页面组件
const PortalPage = lazy(() => import('./pages/portal/PortalPage'));
const WorkbenchPage = lazy(() => import('./pages/workbench/WorkbenchPage'));
const ProjectsPage = lazy(() => import('./pages/projects/ProjectsPage'));
const ProjectDetailPage = lazy(() => import('./pages/projects/ProjectDetailPage'));
const UsersPage = lazy(() => import('./pages/users/UsersPage'));
const LoginPage = lazy(() => import('./pages/auth/LoginPage'));

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
}

const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated } = useAuth();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
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
            <Route path="projects" element={<ProjectsPage />} />
            <Route path="projects/:id" element={<ProjectDetailPage />} />
            <Route path="users" element={<UsersPage />} />
            <Route path="shelf" element={<div>产品货架页面 - 待开发</div>} />
            <Route path="knowledge" element={<div>知识库页面 - 待开发</div>} />
            <Route path="forum" element={<div>技术论坛页面 - 待开发</div>} />
            <Route path="profile" element={<div>个人资料页面 - 待开发</div>} />
            <Route path="settings" element={<div>系统设置页面 - 待开发</div>} />
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
