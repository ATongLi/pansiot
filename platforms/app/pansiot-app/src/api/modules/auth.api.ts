/**
 * 认证 API
 */

// TODO(依赖): FE-007 - 云平台账号系统
// 说明: 此处需要调用云平台认证 API 进行用户登录
// 当前状态: 使用 Mock 实现
// 依赖模块: IMP-007
// 补齐优先级: P0
// 预计补齐日期: 2026-02-15
// 负责人: 开发团队

import type { LoginParams, RegisterParams, LoginResult } from '../types/api.types';
import type { ApiResponse } from '@/types/global';

// === 依赖忽略开始 ===
export const authApi = {
  /**
   * 用户登录 (Mock)
   */
  login: (params: LoginParams): Promise<ApiResponse<LoginResult>> => {
    console.log('[Mock API] login:', params);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 200,
          data: {
            token: 'mock-token-' + Date.now(),
            refreshToken: 'mock-refresh-token',
            expiresIn: 7200,
            user: {
              id: 'mock-user-001',
              username: params.username,
              phone: params.username,
              avatar: '',
              roles: [
                {
                  id: 'role-001',
                  name: '系统管理员',
                  code: 'admin',
                  permissions: ['*'],
                },
              ],
              tenantId: 'mock-tenant-001',
              tenantName: '测试企业',
              createdAt: '2026-01-28',
            },
          },
          message: '登录成功',
        });
      }, 500); // 模拟网络延迟
    });
  },

  /**
   * 用户注册 (Mock)
   */
  register: (params: RegisterParams): Promise<ApiResponse<LoginResult>> => {
    console.log('[Mock API] register:', params);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 200,
          data: {
            token: 'mock-token-' + Date.now(),
            refreshToken: 'mock-refresh-token',
            expiresIn: 7200,
            user: {
              id: 'mock-user-' + Date.now(),
              username: params.username,
              phone: params.phone,
              avatar: '',
              roles: [],
              tenantId: params.tenantId || 'mock-tenant-001',
              tenantName: '测试企业',
              createdAt: new Date().toISOString(),
            },
          },
          message: '注册成功',
        });
      }, 500);
    });
  },

  /**
   * 刷新 Token (Mock)
   */
  refreshToken: (refreshToken: string): Promise<ApiResponse<{ token: string }>> => {
    console.log('[Mock API] refreshToken:', refreshToken);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 200,
          data: {
            token: 'mock-new-token-' + Date.now(),
          },
          message: '刷新成功',
        });
      }, 300);
    });
  },

  /**
   * 登出 (Mock)
   */
  logout: (): Promise<ApiResponse<void>> => {
    console.log('[Mock API] logout');

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 200,
          data: null,
          message: '登出成功',
        });
      }, 300);
    });
  },

  /**
   * 发送验证码 (Mock)
   */
  sendCode: (phone: string): Promise<ApiResponse<void>> => {
    console.log('[Mock API] sendCode:', phone);

    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          code: 200,
          data: null,
          message: '验证码已发送',
        });
      }, 300);
    });
  },
};
// === 依赖忽略结束 ===
