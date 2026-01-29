/**
 * 本地存储工具
 */

import type { Device } from '@/types/global';

/**
 * 存储数据
 */
export const setStorage = <T>(key: string, value: T): void => {
  try {
    const data = JSON.stringify(value);
    uni.setStorageSync(key, data);
  } catch (error) {
    console.error('Storage set error:', error);
  }
};

/**
 * 获取数据
 */
export const getStorage = <T>(key: string): T | null => {
  try {
    const data = uni.getStorageSync(key);
    return data ? JSON.parse(data) : null;
  } catch (error) {
    console.error('Storage get error:', error);
    return null;
  }
};

/**
 * 移除数据
 */
export const removeStorage = (key: string): void => {
  try {
    uni.removeStorageSync(key);
  } catch (error) {
    console.error('Storage remove error:', error);
  }
};

/**
 * 清空所有数据
 */
export const clearStorage = (): void => {
  try {
    uni.clearStorageSync();
  } catch (error) {
    console.error('Storage clear error:', error);
  }
};

/**
 * 获取存储信息
 */
export const getStorageInfo = (): {
  keys: string[];
  currentSize: number;
  limitSize: number;
} => {
  try {
    return uni.getStorageInfoSync();
  } catch (error) {
    console.error('Storage info error:', error);
    return {
      keys: [],
      currentSize: 0,
      limitSize: 0,
    };
  }
};
