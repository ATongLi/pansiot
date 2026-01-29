// 认证相关API

import axios from 'axios';
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  SendVerificationCodeRequest,
  ChangePasswordRequest,
  ResetPasswordRequest,
} from '../types/auth';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
});

// 请求拦截器：添加token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：处理token过期
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token过期，尝试刷新
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await api.post<LoginResponse>('/auth/refresh', {
            refresh_token: refreshToken,
          });
          localStorage.setItem('access_token', response.data.access_token);
          localStorage.setItem('refresh_token', response.data.refresh_token);
          // 重试原请求
          error.config.headers.Authorization = `Bearer ${response.data.access_token}`;
          return api.request(error.config);
        } catch (refreshError) {
          // 刷新失败，清除token并跳转登录
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
          return Promise.reject(refreshError);
        }
      } else {
        // 没有refresh token，跳转登录
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// 登录
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/login', data);
  // 保存token
  localStorage.setItem('access_token', response.data.access_token);
  localStorage.setItem('refresh_token', response.data.refresh_token);
  localStorage.setItem('user', JSON.stringify(response.data.user));
  localStorage.setItem('tenant', JSON.stringify(response.data.tenant));
  return response.data;
};

// 注册（新企业）
export const register = async (data: RegisterRequest): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/register', data);
  // 保存token
  localStorage.setItem('access_token', response.data.access_token);
  localStorage.setItem('refresh_token', response.data.refresh_token);
  localStorage.setItem('user', JSON.stringify(response.data.user));
  localStorage.setItem('tenant', JSON.stringify(response.data.tenant));
  return response.data;
};

// 发送验证码
export const sendVerificationCode = async (data: SendVerificationCodeRequest): Promise<void> => {
  await api.post('/auth/send-code', data);
};

// 重置密码
export const resetPassword = async (data: ResetPasswordRequest): Promise<void> => {
  await api.post('/auth/reset-password', data);
};

// 登出
export const logout = async (): Promise<void> => {
  try {
    await api.post('/auth/logout');
  } finally {
    // 清除本地存储
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    localStorage.removeItem('tenant');
  }
};

// 获取当前用户信息
export const getCurrentUser = async () => {
  const response = await api.get('/auth/me');
  return response.data;
};

// 修改密码
export const changePassword = async (data: ChangePasswordRequest): Promise<void> => {
  await api.post('/auth/change-password', data);
};

// 刷新token
export const refreshToken = async (refreshToken: string): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>('/auth/refresh', {
    refresh_token: refreshToken,
  });
  return response.data;
};

export default api;
