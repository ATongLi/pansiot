import React, { useState } from 'react';
import { Form, Input, Button, Card, Tabs, message, Checkbox } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import { authStore } from '../../stores/auth/AuthStore';
import type { LoginRequest } from '../../types/auth';
import type { RegisterRequest } from '../../types/auth';
import './auth.css';

const LoginForm: React.FC = observer(() => {
  const [activeTab, setActiveTab] = useState('login');
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const handleLogin = async (values: LoginRequest) => {
    try {
      await authStore.login(values);
      message.success('登录成功');
      navigate('/dashboard');
    } catch (error) {
      message.error(authStore.error || '登录失败');
    }
  };

  const handleRegister = async (values: RegisterRequest) => {
    try {
      await authStore.register(values);
      message.success('注册成功');
      navigate('/dashboard');
    } catch (error) {
      message.error(authStore.error || '注册失败');
    }
  };

  const sendCode = async () => {
    const email = form.getFieldValue('email');
    if (!email) {
      message.warning('请先输入邮箱');
      return;
    }
    try {
      await authStore.sendCode(email);
      message.success('验证码已发送');
    } catch (error) {
      message.error('发送验证码失败');
    }
  };

  return (
    <div className="login-container">
      <Card className="login-card">
        <div className="login-header">
          <h1>PansIot Cloud Platform</h1>
          <p>多租户组织管理与RBAC权限控制</p>
        </div>

        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          centered
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form
                  form={form}
                  name="login"
                  onFinish={handleLogin}
                  autoComplete="off"
                  layout="vertical"
                >
                  <Form.Item
                    name="username"
                    rules={[{ required: true, message: '请输入用户名' }]}
                  >
                    <Input
                      prefix={<UserOutlined />}
                      placeholder="用户名"
                      size="large"
                    />
                  </Form.Item>

                  <Form.Item
                    name="password"
                    rules={[{ required: true, message: '请输入密码' }]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="密码"
                      size="large"
                    />
                  </Form.Item>

                  <Form.Item>
                    <Form.Item name="remember" valuePropName="checked" noStyle>
                      <Checkbox>记住我</Checkbox>
                    </Form.Item>
                    <a className="login-form-forgot" href="">
                      忘记密码
                    </a>
                  </Form.Item>

                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      size="large"
                      loading={authStore.isLoading}
                      block
                    >
                      登录
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
            {
              key: 'register',
              label: '注册',
              children: (
                <Form
                  form={form}
                  name="register"
                  onFinish={handleRegister}
                  autoComplete="off"
                  layout="vertical"
                >
                  <Form.Item
                    name="tenantName"
                    label="企业名称"
                    rules={[{ required: true, message: '请输入企业名称' }]}
                  >
                    <Input placeholder="企业名称" size="large" />
                  </Form.Item>

                  <Form.Item
                    name="industry"
                    label="所属行业"
                    rules={[{ required: true, message: '请选择所属行业' }]}
                  >
                    <select className="ant-input" style={{ height: '40px' }}>
                      <option value="">请选择</option>
                      <option value="制造业">制造业</option>
                      <option value="能源">能源</option>
                      <option value="交通">交通</option>
                      <option value="建筑">建筑</option>
                      <option value="农业">农业</option>
                      <option value="医疗">医疗</option>
                      <option value="教育">教育</option>
                      <option value="金融">金融</option>
                      <option value="零售">零售</option>
                      <option value="物流">物流</option>
                      <option value="其他">其他</option>
                    </select>
                  </Form.Item>

                  <Form.Item
                    name="username"
                    label="用户名"
                    rules={[
                      { required: true, message: '请输入用户名' },
                      { min: 3, message: '用户名至少3位' },
                    ]}
                  >
                    <Input
                      prefix={<UserOutlined />}
                      placeholder="用户名"
                      size="large"
                    />
                  </Form.Item>

                  <Form.Item
                    name="email"
                    label="邮箱"
                    rules={[
                      { required: true, message: '请输入邮箱' },
                      { type: 'email', message: '邮箱格式不正确' },
                    ]}
                  >
                    <Input
                      prefix={<MailOutlined />}
                      placeholder="邮箱"
                      size="large"
                    />
                  </Form.Item>

                  <Form.Item
                    name="password"
                    label="密码"
                    rules={[
                      { required: true, message: '请输入密码' },
                      { min: 8, message: '密码至少8位' },
                    ]}
                  >
                    <Input.Password
                      prefix={<LockOutlined />}
                      placeholder="密码（至少8位）"
                      size="large"
                    />
                  </Form.Item>

                  <Form.Item
                    name="verificationCode"
                    label="验证码"
                    rules={[{ required: true, message: '请输入验证码' }]}
                  >
                    <Input
                      placeholder="验证码"
                      size="large"
                      suffix={
                        <Button type="link" onClick={sendCode}>
                          发送验证码
                        </Button>
                      }
                    />
                  </Form.Item>

                  <Form.Item>
                    <Button
                      type="primary"
                      htmlType="submit"
                      size="large"
                      loading={authStore.isLoading}
                      block
                    >
                      注册
                    </Button>
                  </Form.Item>
                </Form>
              ),
            },
          ]}
        />
      </Card>
    </div>
  );
});

export default LoginForm;
