/**
 * HTTP 客户端封装
 */
import { interceptor } from './interceptors';
import { requestConfig } from './config';
import type { RequestConfig, ApiResponse } from '../types/api.types';

class Request {
  private config: RequestConfig;

  constructor(config: RequestConfig) {
    this.config = config;
  }

  /**
   * 通用请求方法
   */
  request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
    // 请求前拦截
    const requestConfig = interceptor.request(config);

    return new Promise((resolve, reject) => {
      uni.request({
        url: this.getFullUrl(requestConfig.url),
        method: requestConfig.method || 'GET',
        data: requestConfig.data,
        header: requestConfig.header,
        timeout: requestConfig.timeout || this.config.timeout,
        success: (res) => {
          // 响应拦截
          const response = interceptor.response<T>(res, requestConfig);
          resolve(response);
        },
        fail: (err) => {
          // 错误拦截
          const error = interceptor.error(err, requestConfig);
          reject(error);
        },
      });
    });
  }

  /**
   * GET 请求
   */
  get<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'GET', data, ...config });
  }

  /**
   * POST 请求
   */
  post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'POST', data, ...config });
  }

  /**
   * PUT 请求
   */
  put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'PUT', data, ...config });
  }

  /**
   * DELETE 请求
   */
  delete<T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'DELETE', data, ...config });
  }

  /**
   * 获取完整 URL
   */
  private getFullUrl(url: string): string {
    if (url.startsWith('http')) {
      return url;
    }
    return this.config.baseURL + url;
  }
}

export const request = new Request(requestConfig);
