import React from 'react';
import { Layout, Menu, Avatar, Dropdown } from 'antd';
import {
  UserOutlined,
  LogoutOutlined,
  DashboardOutlined,
  TeamOutlined,
  UserSwitchOutlined,
  SafetyCertificateOutlined,
  FileTextOutlined,
} from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { useNavigate } from 'react-router-dom';
import { authStore } from '../../stores/auth/AuthStore';
import './Dashboard.css';

const { Header, Sider, Content } = Layout;

const Dashboard: React.FC = observer(() => {
  const navigate = useNavigate();

  const handleLogout = async () => {
    await authStore.logout();
    navigate('/login');
  };

  const menuItems = [
    {
      key: 'dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: 'organizations',
      icon: <TeamOutlined />,
      label: '组织管理',
      onClick: () => navigate('/organizations'),
    },
    {
      key: 'users',
      icon: <UserSwitchOutlined />,
      label: '用户管理',
      onClick: () => navigate('/users'),
    },
    {
      key: 'roles',
      icon: <SafetyCertificateOutlined />,
      label: '角色管理',
      onClick: () => navigate('/roles'),
    },
    {
      key: 'audit-logs',
      icon: <FileTextOutlined />,
      label: '审计日志',
      disabled: true,
    },
  ];

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人设置',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={200} theme="light">
        <div className="logo">
          <h2>PansIot Cloud</h2>
        </div>
        <Menu
          mode="inline"
          defaultSelectedKeys={['dashboard']}
          style={{ height: '100%', borderRight: 0 }}
          items={menuItems}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h2 style={{ margin: 0 }}>仪表板</h2>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span>
              {authStore.tenant?.name} ({authStore.tenant?.tenant_type === 'INTEGRATOR' ? '集成商' : '终端租户'})
            </span>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Avatar icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
            </Dropdown>
          </div>
        </Header>
        <Content style={{ margin: '24px', padding: 24, background: '#fff', minHeight: 280 }}>
          <div className="dashboard-content">
            <h1>欢迎, {authStore.user?.real_name || authStore.user?.username}!</h1>

            <div style={{ marginTop: 24 }}>
              <h3>组织信息</h3>
              <p><strong>组织名称:</strong> {authStore.tenant?.name}</p>
              <p><strong>组织类型:</strong> {authStore.tenant?.tenant_type === 'INTEGRATOR' ? '集成商' : '终端租户'}</p>
              <p><strong>所属行业:</strong> {authStore.tenant?.industry}</p>
              <p><strong>企业序列号:</strong> {authStore.tenant?.serial_number}</p>
            </div>

            <div style={{ marginTop: 24 }}>
              <h3>用户信息</h3>
              <p><strong>用户名:</strong> {authStore.user?.username}</p>
              <p><strong>邮箱:</strong> {authStore.user?.email}</p>
              <p><strong>手机号:</strong> {authStore.user?.phone || '未设置'}</p>
            </div>

            <div style={{ marginTop: 24 }}>
              <h3>功能模块</h3>
              <p>当前可用的功能模块正在开发中...</p>
            </div>
          </div>
        </Content>
      </Layout>
    </Layout>
  );
});

export default Dashboard;
