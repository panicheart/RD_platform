import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './hooks/useAuth'
import MainLayout from './layouts/MainLayout'
import PortalPage from './pages/portal/PortalPage'
import WorkbenchPage from './pages/workbench/WorkbenchPage'
import LoginPage from './pages/auth/LoginPage'
import ProjectsPage from './pages/projects/ProjectsPage'
import ProjectDetailPage from './pages/projects/ProjectDetailPage'
import UsersPage from './pages/users/UsersPage'

// Protected route wrapper
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  
  return <>{children}</>
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public routes - 公开页面，无需登录 */}
          <Route path="/" element={<PortalPage />} />
          <Route path="/portal" element={<PortalPage />} />
          <Route path="/login" element={<LoginPage />} />
          
          {/* Protected routes - 需要登录的页面 */}
          <Route
            path="/workbench"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<WorkbenchPage />} />
          </Route>
          
          <Route
            path="/admin"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route path="projects" element={<ProjectsPage />} />
            <Route path="projects/:id" element={<ProjectDetailPage />} />
            <Route path="users" element={<UsersPage />} />
          </Route>
          
          {/* Redirect old routes */}
          <Route path="/projects" element={<Navigate to="/workbench" replace />} />
          <Route path="/users" element={<Navigate to="/workbench" replace />} />
          
          {/* 404 */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
