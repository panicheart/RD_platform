import { Suspense, lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Spin, Result, Button } from 'antd';
import { AuthProvider } from '@/hooks/useAuth';
import MainLayout from '@/layouts/MainLayout';
import PortalLayout from '@/layouts/PortalLayout';

// Lazy load pages
const PortalPage = lazy(() => import('./pages/portal/PortalPage'));
const WorkbenchPage = lazy(() => import('./pages/workbench/WorkbenchPage'));
const ProjectsPage = lazy(() => import('./pages/projects/ProjectsPage'));
const ProjectDetailPage = lazy(() => import('./pages/projects/ProjectDetailPage'));
const LoginPage = lazy(() => import('./pages/auth/LoginPage'));
const UsersPage = lazy(() => import('./pages/users/UsersPage'));
const ShelfPage = lazy(() => import('./pages/shelf/ShelfPage'));
const KnowledgeList = lazy(() => import('./pages/knowledge/KnowledgeList'));
const ForumPage = lazy(() => import('./pages/forum/ForumPage'));
const ForumBoardPage = lazy(() => import('./pages/forum/ForumBoardPage'));
const ForumPostPage = lazy(() => import('./pages/forum/ForumPostPage'));
const ForumCreatePostPage = lazy(() => import('./pages/forum/ForumCreatePostPage'));
const AnalyticsDashboard = lazy(() => import('./pages/analytics/AnalyticsDashboard'));
const MonitorDashboard = lazy(() => import('./pages/monitor/MonitorDashboard'));

// Page loading component
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

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Suspense fallback={<PageLoading />}>
          <Routes>
            {/* Public routes */}
            <Route path="/" element={<Navigate to="/portal" replace />} />
            <Route path="/portal" element={<PortalPage />} />
            <Route path="/login" element={<LoginPage />} />

            {/* Protected routes with MainLayout */}
            <Route element={<MainLayout />}>
              <Route path="/workbench" element={<WorkbenchPage />} />
              <Route path="/projects" element={<ProjectsPage />} />
              <Route path="/projects/:id" element={<ProjectDetailPage />} />
              <Route path="/users" element={<UsersPage />} />
              <Route path="/shelf" element={<ShelfPage />} />
              <Route path="/knowledge" element={<KnowledgeList />} />
              <Route path="/forum" element={<ForumPage />} />
              <Route path="/forum/boards/:boardId" element={<ForumBoardPage />} />
              <Route path="/forum/posts/create" element={<ForumCreatePostPage />} />
              <Route path="/forum/posts/:postId" element={<ForumPostPage />} />
              <Route path="/analytics" element={<AnalyticsDashboard />} />
              <Route path="/monitor" element={<MonitorDashboard />} />
              <Route path="/profile" element={<div>个人资料页面 - 待开发</div>} />
              <Route path="/settings" element={<div>系统设置页面 - 待开发</div>} />
              <Route path="/notifications" element={<div>消息通知页面 - 待开发</div>} />
            </Route>

            {/* 404 Page */}
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
    </AuthProvider>
  );
}

export default App;
