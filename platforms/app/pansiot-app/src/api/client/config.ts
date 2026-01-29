/**
 * HTTP 客户端配置
 */
import type { RequestConfig } from '../types/api.types';

export const requestConfig: RequestConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 30000,
  header: {
    'Content-Type': 'application/json',
  },
};
