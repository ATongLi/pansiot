/**
 * API 类型定义
 */

import type { ApiResponse } from '@/types/global';

// 请求配置
export interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  header?: Record<string, string>;
  timeout?: number;
}

// 登录参数
export interface LoginParams {
  username: string;
  password: string;
  tenantId?: string;
}

// 登录结果
export interface LoginResult {
  token: string;
  refreshToken: string;
  expiresIn: number;
  user: import('@/types/global').UserInfo;
}

// 注册参数
export interface RegisterParams {
  username: string;
  password: string;
  phone: string;
  code: string;
  tenantId?: string;
}
