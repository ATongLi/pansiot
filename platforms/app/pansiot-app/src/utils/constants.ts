/**
 * 常量定义
 */

// API 基础 URL
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'https://api.pansiot.com';

// Mock 开关
export const MOCK_ENABLED = import.meta.env.VITE_MOCK_ENABLED === 'true';

// Token 存储键
export const TOKEN_KEY = 'token';
export const REFRESH_TOKEN_KEY = 'refreshToken';
export const USER_INFO_KEY = 'userInfo';

// 缓存时间 (秒)
export const CACHE_EXPIRE_TIME = 7 * 24 * 60 * 60; // 7天

// 分页默认配置
export const DEFAULT_PAGE_SIZE = 20;
export const PAGE_SIZES = [10, 20, 50, 100];

// 设备状态
export const DEVICE_STATUS = {
  NORMAL: 'normal',
  WARNING: 'warning',
  ERROR: 'error',
  OFFLINE: 'offline',
} as const;

// 设备状态颜色
export const DEVICE_STATUS_COLOR = {
  normal: '#4cd964',
  warning: '#f0ad4e',
  error: '#dd524d',
  offline: '#8a8a8a',
} as const;

// 消息类型
export const MESSAGE_TYPE = {
  SYSTEM: 'system',
  DEVICE: 'device',
  ALERT: 'alert',
} as const;

// 消息类型颜色
export const MESSAGE_TYPE_COLOR = {
  system: '#007aff',
  device: '#4cd964',
  alert: '#dd524d',
} as const;
