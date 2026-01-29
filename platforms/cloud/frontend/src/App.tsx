import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import { ConfigProvider, theme } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import EditorApp from './editor/EditorApp';
import RuntimeApp from './runtime/RuntimeApp';
import LoginForm from './components/auth/LoginForm';
import Dashboard from './pages/dashboard/Dashboard';
import UserManagement from './pages/users/UserManagement';
import RoleManagement from './pages/roles/RoleManagement';
import OrganizationManagement from './pages/organizations/OrganizationManagement';
import { authStore } from './stores/auth/AuthStore';

// 受保护的路由组件
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = observer(({ children }) => {
  if (!authStore.isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
});

// 主应用组件
const App: React.FC = observer(() => {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1890ff',
        },
      }}
    >
      <Router>
        <Routes>
          {/* 首页重定向 */}
          <Route
            path="/"
            element={
              authStore.isAuthenticated ? (
                <Navigate to="/dashboard" replace />
              ) : (
                <Navigate to="/login" replace />
              )
            }
          />

          {/* 登录页 */}
          <Route path="/login" element={<LoginForm />} />

          {/* 仪表板 */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />

          {/* 组织管理 */}
          <Route
            path="/organizations"
            element={
              <ProtectedRoute>
                <OrganizationManagement />
              </ProtectedRoute>
            }
          />

          {/* 用户管理 */}
          <Route
            path="/users"
            element={
              <ProtectedRoute>
                <UserManagement />
              </ProtectedRoute>
            }
          />

          {/* 角色管理 */}
          <Route
            path="/roles"
            element={
              <ProtectedRoute>
                <RoleManagement />
              </ProtectedRoute>
            }
          />

          {/* Web工程编辑器 */}
          <Route path="/editor" element={<EditorApp />} />
          <Route path="/runtime/:projectId" element={<RuntimeApp />} />

          {/* 404 */}
          <Route path="*" element={<div>404 - 页面未找到</div>} />
        </Routes>
      </Router>
    </ConfigProvider>
  );
});

export default App;
