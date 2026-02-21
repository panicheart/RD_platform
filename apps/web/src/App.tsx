import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from 'antd';

const { Content } = Layout;

function App() {
  return (
    <BrowserRouter>
      <Layout style={{ minHeight: '100vh' }}>
        <Content>
          <Routes>
            <Route path="/" element={<Navigate to="/portal" replace />} />
            <Route path="/portal" element={<div>Portal Page - TODO</div>} />
            <Route path="/workbench" element={<div>Workbench Page - TODO</div>} />
            <Route path="/projects" element={<div>Projects Page - TODO</div>} />
            <Route path="/profile" element={<div>Profile Page - TODO</div>} />
            <Route path="*" element={<Navigate to="/portal" replace />} />
          </Routes>
        </Content>
      </Layout>
    </BrowserRouter>
  );
}

export default App;
