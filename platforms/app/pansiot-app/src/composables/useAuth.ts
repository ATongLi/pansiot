/**
 * 认证相关的组合式函数
 */

import { computed } from 'vue';
import { useUserStore } from '@/stores/user.store';

/**
 * 使用认证
 */
export function useAuth() {
  const userStore = useUserStore();

  // 是否已登录
  const isLoggedIn = computed(() => userStore.isLoggedIn);

  // 用户信息
  const userInfo = computed(() => userStore.userInfo);

  // Token
  const token = computed(() => userStore.token);

  /**
   * 登录
   */
  const login = async (username: string, password: string) => {
    const success = await userStore.login({ username, password });
    return success;
  };

  /**
   * 登出
   */
  const logout = async () => {
    await userStore.logout();
    uni.reLaunch({
      url: '/pages/index/index',
    });
  };

  /**
   * 检查权限
   */
  const hasPermission = (permission: string): boolean => {
    return userStore.hasPermission(permission);
  };

  /**
   * 检查是否是管理员
   */
  const isAdmin = computed(() => {
    return hasPermission('*');
  });

  return {
    isLoggedIn,
    userInfo,
    token,
    login,
    logout,
    hasPermission,
    isAdmin,
  };
}
