/**
 * 全局类型定义
 */

// API 响应类型
export interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}

// 分页参数
export interface PageParams {
  page: number;
  pageSize: number;
}

// 分页结果
export interface PageResult<T> {
  list: T[];
  total: number;
}

// 用户信息
export interface UserInfo {
  id: string;
  username: string;
  phone?: string;
  email?: string;
  avatar?: string;
  roles: Role[];
  tenantId: string;
  tenantName: string;
  createdAt: string;
}

// 角色信息
export interface Role {
  id: string;
  name: string;
  code: string;
  permissions: string[];
}

// 租户信息
export interface Tenant {
  id: string;
  name: string;
  code: string;
  type: 'platform' | 'integrator' | 'customer' | 'sub';
}

// 设备信息
export interface Device {
  id: string;
  name: string;
  type: string;
  status: DeviceStatus;
  online: boolean;
  location?: string;
  lastOnline?: string;
  properties: Record<string, any>;
  createdAt: string;
}

// 设备状态枚举
export type DeviceStatus = 'normal' | 'warning' | 'error' | 'offline';

// 消息信息
export interface Message {
  id: string;
  title: string;
  content: string;
  type: 'system' | 'device' | 'alert';
  read: boolean;
  createdAt: string;
}
