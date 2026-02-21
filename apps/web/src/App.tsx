import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '@/hooks/useAuth';
import MainLayout from '@/layouts/MainLayout';
import PortalPage from '@/pages/portal/PortalPage';
import LoginPage from '@/pages/auth/LoginPage';
import WorkbenchPage from '@/pages/workbench/WorkbenchPage';
import ProjectsPage from '@/pages/projects/ProjectsPage';
import ProjectDetailPage from '@/pages/projects/ProjectDetailPage';
import UsersPage from '@/pages/users/UsersPage';

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Portal 和 Login 不使用 MainLayout */}
          <Route path="/" element={<Navigate to="/portal" replace />} />
          <Route path="/portal" element={<PortalPage />} />
          <Route path="/login" element={<LoginPage />} />

          {/* 需要侧边栏布局的页面 */}
          <Route element={<MainLayout />}>
            <Route path="/workbench" element={<WorkbenchPage />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="/projects/:id" element={<ProjectDetailPage />} />
            <Route path="/users" element={<UsersPage />} />
          </Route>

          <Route path="*" element={<Navigate to="/portal" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
