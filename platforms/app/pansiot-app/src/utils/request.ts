/**
 * HTTP 请求封装
 */

import type { ApiResponse } from '@/types/global';
import { TOKEN_KEY, API_BASE_URL } from './constants';

/**
 * 请求配置
 */
interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  header?: Record<string, string>;
  timeout?: number;
}

/**
 * 请求拦截器
 */
const requestInterceptor = (config: RequestConfig): RequestConfig => {
  // 添加 Token
  const token = uni.getStorageSync(TOKEN_KEY);
  if (token) {
    config.header = {
      ...config.header,
      Authorization: `Bearer ${token}`,
    };
  }

  // 添加基础 URL
  if (!config.url.startsWith('http')) {
    config.url = API_BASE_URL + config.url;
  }

  console.log('[Request]', config.method, config.url, config.data);

  return config;
};

/**
 * 响应拦截器
 */
const responseInterceptor = (response: any): ApiResponse => {
  console.log('[Response]', response);

  const { statusCode, data } = response;

  // HTTP 状态码检查
  if (statusCode !== 200) {
    uni.showToast({
      title: `请求失败 (${statusCode})`,
      icon: 'none',
    });
    return {
      code: statusCode,
      data: null,
      message: '请求失败',
    };
  }

  // 业务状态码检查
  if (data.code !== 200) {
    // Token 过期,跳转登录
    if (data.code === 401) {
      uni.removeStorageSync(TOKEN_KEY);
      uni.reLaunch({
        url: '/pages/index/index',
      });
    }

    uni.showToast({
      title: data.message || '请求失败',
      icon: 'none',
    });

    return data;
  }

  return data;
};

/**
 * 统一请求方法
 */
export const request = <T = any>(config: RequestConfig): Promise<ApiResponse<T>> => {
  return new Promise((resolve, reject) => {
    // 请求拦截
    const interceptedConfig = requestInterceptor(config);

    uni.request({
      url: interceptedConfig.url,
      method: interceptedConfig.method || 'GET',
      data: interceptedConfig.data,
      header: interceptedConfig.header,
      timeout: interceptedConfig.timeout || 30000,
      success: (res) => {
        const result = responseInterceptor(res);
        resolve(result as ApiResponse<T>);
      },
      fail: (err) => {
        console.error('[Request Error]', err);
        uni.showToast({
          title: '网络请求失败',
          icon: 'none',
        });
        reject(err);
      },
    });
  });
};

/**
 * GET 请求
 */
export const get = <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return request<T>({
    url,
    method: 'GET',
    data,
  });
};

/**
 * POST 请求
 */
export const post = <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return request<T>({
    url,
    method: 'POST',
    data,
  });
};

/**
 * PUT 请求
 */
export const put = <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return request<T>({
    url,
    method: 'PUT',
    data,
  });
};

/**
 * DELETE 请求
 */
export const del = <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return request<T>({
    url,
    method: 'DELETE',
    data,
  });
};

export default {
  request,
  get,
  post,
  put,
  delete: del,
};
