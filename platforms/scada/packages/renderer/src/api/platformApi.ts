/**
 * 硬件平台配置 API 客户端
 * 提供与 Scada 后端硬件平台接口的通信
 */

import type { HardwarePlatformConfig } from '@/types/project'
import type { ApiResponse } from '@/types/project'
import { getApiBaseUrl } from '@/utils/electron'

/**
 * API 错误类
 */
export class ApiError extends Error {
  constructor(
    public code: string,
    public message: string,
    public statusCode?: number
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/**
 * 硬件平台配置 API 客户端类
 */
class PlatformApiClient {
  private baseURL: string

  constructor() {
    // 动态获取API基础URL（Electron vs 开发环境）
    this.baseURL = getApiBaseUrl()
  }

  /**
   * 通用请求方法
   */
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`

    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      }
    }

    try {
      const response = await fetch(url, config)
      const data = await response.json()

      if (!response.ok) {
        throw new ApiError(
          data.error || 'UNKNOWN_ERROR',
          data.message || '请求失败',
          response.status
        )
      }

      return data
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      throw new ApiError('NETWORK_ERROR', '网络请求失败')
    }
  }

  /**
   * 获取所有启用的硬件平台
   * @returns 硬件平台列表
   */
  async getAllPlatforms(): Promise<ApiResponse<HardwarePlatformConfig[]>> {
    return this.request<HardwarePlatformConfig[]>('/api/platforms', {
      method: 'GET'
    })
  }
}

/**
 * 导出单例实例
 */
export const platformApi = new PlatformApiClient()

/**
 * 导出类型
 */
export type { ApiError }
