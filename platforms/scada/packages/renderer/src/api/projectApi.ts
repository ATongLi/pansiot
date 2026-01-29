/**
 * 工程管理 API 客户端
 * 提供与 Scada 后端通信的接口
 */

import type {
  Project,
  NewProjectFormData,
  OpenProjectRequest,
  ApiResponse,
  CreateProjectResponse,
  OpenProjectResponse,
  RecentProject
} from '@/types/project'
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
 * 工程管理 API 客户端类
 */
class ProjectApiClient {
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
   * 创建新工程
   * @param formData 工程表单数据
   * @returns 创建的工程信息
   */
  async createProject(
    formData: NewProjectFormData
  ): Promise<ApiResponse<CreateProjectResponse>> {
    return this.request<CreateProjectResponse>('/api/projects/create', {
      method: 'POST',
      body: JSON.stringify(formData)
    })
  }

  /**
   * 打开工程
   * @param requestData 打开工程请求数据
   * @returns 工程数据
   */
  async openProject(
    requestData: OpenProjectRequest
  ): Promise<ApiResponse<OpenProjectResponse>> {
    return this.request<OpenProjectResponse>('/api/projects/open', {
      method: 'POST',
      body: JSON.stringify(requestData)
    })
  }

  /**
   * 保存工程
   * @param project 工程数据
   * @returns 保存结果
   */
  async saveProject(
    project: Project
  ): Promise<ApiResponse<{ filePath: string }>> {
    return this.request<{ filePath: string }>('/api/projects/save', {
      method: 'POST',
      body: JSON.stringify(project)
    })
  }

  /**
   * 验证工程密码
   * @param filePath 工程文件路径
   * @param password 密码
   * @returns 验证结果
   */
  async validatePassword(
    filePath: string,
    password: string
  ): Promise<ApiResponse<{ valid: boolean }>> {
    return this.request<{ valid: boolean }>('/api/projects/validate-password', {
      method: 'POST',
      body: JSON.stringify({ filePath, password })
    })
  }

  /**
   * 获取最近工程列表
   * @returns 最近工程列表
   */
  async getRecentProjects(): Promise<ApiResponse<RecentProject[]>> {
    return this.request<RecentProject[]>('/api/projects/recent', {
      method: 'GET'
    })
  }

  /**
   * 添加或更新最近工程
   * @param project 工程数据
   * @returns 操作结果
   */
  async addOrUpdateRecentProject(
    project: RecentProject
  ): Promise<ApiResponse<{ success: boolean }>> {
    return this.request<{ success: boolean }>('/api/projects/recent', {
      method: 'POST',
      body: JSON.stringify(project)
    })
  }

  /**
   * 从最近工程列表中移除
   * @param projectId 工程ID
   * @returns 操作结果
   */
  async removeRecentProject(
    projectId: string
  ): Promise<ApiResponse<{ success: boolean }>> {
    return this.request<{ success: boolean }>(`/api/projects/recent/${projectId}`, {
      method: 'DELETE'
    })
  }
}

/**
 * 导出单例实例
 */
export const projectApi = new ProjectApiClient()

/**
 * 导出类型
 */
export type { ApiError }
