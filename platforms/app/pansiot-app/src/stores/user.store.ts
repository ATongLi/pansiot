/**
 * 用户状态管理
 */
import { defineStore } from 'pinia';
import type { UserInfo } from '@/types/global';
import type { LoginParams } from '@/api/types/api.types';
import { authApi } from '@/api/modules/auth.api';

interface UserState {
  token: string;
  refreshToken: string;
  userInfo: UserInfo | null;
  isLoggedIn: boolean;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: '',
    refreshToken: '',
    userInfo: null,
    isLoggedIn: false,
  }),

  getters: {
    /**
     * 检查是否有权限
     */
    hasPermission: (state) => {
      return (permission: string): boolean => {
        return (
          state.userInfo?.roles.some((role) => role.permissions.includes(permission) || role.permissions.includes('*')) ??
          false
        );
      };
    },
  },

  actions: {
    /**
     * 用户登录
     */
    async login(params: LoginParams) {
      const res = await authApi.login(params);

      if (res.code === 200 && res.data) {
        this.token = res.data.token;
        this.refreshToken = res.data.refreshToken;
        this.userInfo = res.data.user;
        this.isLoggedIn = true;

        // 保存到本地存储
        uni.setStorageSync('token', res.data.token);
        uni.setStorageSync('refreshToken', res.data.refreshToken);
        uni.setStorageSync('userInfo', res.data.user);

        return true;
      }

      return false;
    },

    /**
     * 用户登出
     */
    async logout() {
      await authApi.logout();

      this.token = '';
      this.refreshToken = '';
      this.userInfo = null;
      this.isLoggedIn = false;

      // 清除本地存储
      uni.removeStorageSync('token');
      uni.removeStorageSync('refreshToken');
      uni.removeStorageSync('userInfo');
    },

    /**
     * 从本地存储恢复用户信息
     */
    restoreUser() {
      const token = uni.getStorageSync('token');
      const refreshToken = uni.getStorageSync('refreshToken');
      const userInfo = uni.getStorageSync('userInfo');

      if (token && userInfo) {
        this.token = token;
        this.refreshToken = refreshToken;
        this.userInfo = userInfo;
        this.isLoggedIn = true;
      }
    },
  },

  // 持久化配置
  persist: {
    key: 'user-store',
    storage: {
      getItem: (key) => uni.getStorageSync(key),
      setItem: (key, value) => uni.setStorageSync(key, value),
    },
  },
});
