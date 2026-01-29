import { makeAutoObservable, runInAction } from 'mobx';
import {
  login,
  logout,
  getCurrentUser,
  changePassword,
  sendVerificationCode,
} from '../../api/auth';
import type {
  User,
  Tenant,
  LoginRequest,
  RegisterRequest,
  ChangePasswordRequest,
} from '../../types/auth';

class AuthStore {
  // 状态
  user: User | null = null;
  tenant: Tenant | null = null;
  isAuthenticated = false;
  isLoading = false;
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
    // 从localStorage恢复状态
    this.loadFromStorage();
  }

  // 加载本地存储的状态
  loadFromStorage() {
    try {
      const userStr = localStorage.getItem('user');
      const tenantStr = localStorage.getItem('tenant');
      const token = localStorage.getItem('access_token');

      if (userStr && tenantStr && token) {
        this.user = JSON.parse(userStr);
        this.tenant = JSON.parse(tenantStr);
        this.isAuthenticated = true;
      }
    } catch (error) {
      console.error('Failed to load auth state from storage:', error);
    }
  }

  // 登录
  login = async (data: LoginRequest) => {
    this.isLoading = true;
    this.error = null;
    try {
      const response = await login(data);
      runInAction(() => {
        this.user = response.user;
        this.tenant = response.tenant;
        this.isAuthenticated = true;
        this.isLoading = false;
      });
    } catch (error: any) {
      runInAction(() => {
        this.error = error.response?.data?.message || '登录失败';
        this.isLoading = false;
      });
      throw error;
    }
  };

  // 注册
  register = async (data: RegisterRequest) => {
    this.isLoading = true;
    this.error = null;
    try {
      const response = await login(data); // 注册后自动登录
      runInAction(() => {
        this.user = response.user;
        this.tenant = response.tenant;
        this.isAuthenticated = true;
        this.isLoading = false;
      });
    } catch (error: any) {
      runInAction(() => {
        this.error = error.response?.data?.message || '注册失败';
        this.isLoading = false;
      });
      throw error;
    }
  };

  // 登出
  logout = async () => {
    this.isLoading = true;
    try {
      await logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      runInAction(() => {
        this.user = null;
        this.tenant = null;
        this.isAuthenticated = false;
        this.isLoading = false;
      });
    }
  };

  // 获取当前用户信息
  fetchCurrentUser = async () => {
    this.isLoading = true;
    try {
      const data = await getCurrentUser();
      runInAction(() => {
        this.user = data.user;
        this.tenant = data.tenant;
        this.isAuthenticated = true;
        this.isLoading = false;
      });
    } catch (error: any) {
      runInAction(() => {
        this.error = error.response?.data?.message || '获取用户信息失败';
        this.isLoading = false;
      });
      throw error;
    }
  };

  // 修改密码
  updatePassword = async (data: ChangePasswordRequest) => {
    this.isLoading = true;
    this.error = null;
    try {
      await changePassword(data);
      runInAction(() => {
        this.isLoading = false;
      });
    } catch (error: any) {
      runInAction(() => {
        this.error = error.response?.data?.message || '修改密码失败';
        this.isLoading = false;
      });
      throw error;
    }
  };

  // 发送验证码
  sendCode = async (email: string, phone?: string) => {
    this.isLoading = true;
    this.error = null;
    try {
      await sendVerificationCode({ email, phone });
      runInAction(() => {
        this.isLoading = false;
      });
    } catch (error: any) {
      runInAction(() => {
        this.error = error.response?.data?.message || '发送验证码失败';
        this.isLoading = false;
      });
      throw error;
    }
  };

  // 清除错误
  clearError = () => {
    this.error = null;
  };

  // 检查是否是集成商
  get isIntegrator() {
    return this.tenant?.tenant_type === 'INTEGRATOR';
  }

  // 检查是否是终端租户
  get isTerminal() {
    return this.tenant?.tenant_type === 'TERMINAL';
  }
}

export const authStore = new AuthStore();
