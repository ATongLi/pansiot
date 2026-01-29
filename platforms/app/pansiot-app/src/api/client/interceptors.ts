/**
 * 请求/响应拦截器
 */
import type { RequestConfig, ApiResponse } from '../types/api.types';
import { useUserStore } from '@/stores/user.store';
import { useTenantStore } from '@/stores/tenant.store';

class Interceptor {
  /**
   * 请求拦截器
   */
  request(config: RequestConfig): RequestConfig {
    // 添加 Token
    const userStore = useUserStore();
    if (userStore.token) {
      config.header = {
        ...config.header,
        Authorization: `Bearer ${userStore.token}`,
      };
    }

    // 添加租户 ID
    const tenantStore = useTenantStore();
    if (tenantStore.currentTenant?.id) {
      config.header = {
        ...config.header,
        'X-Tenant-ID': tenantStore.currentTenant.id,
      };
    }

    // 添加通用 Header
    config.header = {
      'Content-Type': 'application/json',
      ...config.header,
    };

    // 请求时间戳(防止缓存)
    if (config.method === 'GET') {
      config.data = {
        ...config.data,
        _t: Date.now(),
      };
    }

    console.log(`[Request] ${config.method} ${config.url}`, config.data);
    return config;
  }

  /**
   * 响应拦截器
   */
  response<T>(res: any, config: RequestConfig): ApiResponse<T> {
    console.log(`[Response] ${config.url}`, res);

    const { statusCode, data } = res;

    // HTTP 状态码检查
    if (statusCode !== 200) {
      return {
        code: statusCode,
        data: null as any,
        message: '请求失败',
      };
    }

    // 业务状态码检查
    if (data.code !== 200) {
      // Token 过期,自动刷新
      if (data.code === 401) {
        this.handleTokenExpired();
      }

      return {
        code: data.code,
        data: null as any,
        message: data.message || '请求失败',
      };
    }

    return {
      code: data.code,
      data: data.data,
      message: data.message,
    };
  }

  /**
   * 错误拦截器
   */
  error(err: any, config: RequestConfig): never {
    console.error(`[Error] ${config.url}`, err);

    let message = '网络错误';

    // 网络超时
    if (err.errMsg?.includes('timeout')) {
      message = '请求超时';
    }
    // 网络异常
    else if (err.errMsg?.includes('fail')) {
      message = '网络连接失败';
    }

    uni.showToast({
      title: message,
      icon: 'none',
    });

    throw new Error(message);
  }

  /**
   * 处理 Token 过期
   */
  private async handleTokenExpired() {
    const userStore = useUserStore();

    try {
      // 刷新 Token
      const res = await uni.request({
        url: `${import.meta.env.VITE_API_BASE_URL}/api/auth/refresh`,
        method: 'POST',
        data: {
          refreshToken: userStore.refreshToken,
        },
      });

      if (res.data.code === 200) {
        userStore.token = res.data.data.token;
        // 重试原请求
      } else {
        // 刷新失败,跳转登录
        userStore.logout();
        uni.navigateTo({ url: '/pages/index/login' });
      }
    } catch (error) {
      userStore.logout();
      uni.navigateTo({ url: '/pages/index/login' });
    }
  }
}

export const interceptor = new Interceptor();
