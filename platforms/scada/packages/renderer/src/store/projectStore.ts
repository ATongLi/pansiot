/**
 * 工程状态管理
 * 管理当前工程的状态和操作
 */

import { makeAutoObservable, runInAction } from 'mobx'
import type {
  Project,
  NewProjectFormData,
  OpenProjectRequest,
} from '@/types/project'
import { projectApi } from '@/api/projectApi'

/**
 * 工程状态类型
 */
type ProjectStatus = 'idle' | 'creating' | 'opening' | 'saving' | 'error'

/**
 * ProjectStore 工程状态管理
 */
class ProjectStore {
  // 当前工程
  currentProject: Project | null = null

  // 工程状态
  status: ProjectStatus = 'idle'

  // 错误信息
  error: string = ''

  constructor() {
    makeAutoObservable(this)
  }

  /**
   * 创建新工程
   */
  async createProject(formData: NewProjectFormData): Promise<void> {
    runInAction(() => {
      this.status = 'creating'
      this.error = ''
    })

    try {
      const response = await projectApi.createProject(formData)

      if (response.success && response.data) {
        runInAction(() => {
          this.status = 'idle'
          // 创建成功后打开工程
          this.openProject({
            filePath: response.data!.filePath,
            password: formData.password,
          })
        })
      } else {
        runInAction(() => {
          this.status = 'error'
          this.error = response.message || '创建工程失败'
        })
      }
    } catch (err) {
      runInAction(() => {
        this.status = 'error'
        this.error = err instanceof Error ? err.message : '创建工程出错'
      })
    }
  }

  /**
   * 打开工程
   */
  async openProject(request: OpenProjectRequest): Promise<void> {
    runInAction(() => {
      this.status = 'opening'
      this.error = ''
    })

    try {
      const response = await projectApi.openProject(request)

      if (response.success && response.data) {
        runInAction(() => {
          this.currentProject = response.data!.project
          this.status = 'idle'
        })
      } else {
        runInAction(() => {
          this.status = 'error'
          this.error = response.message || '打开工程失败'
        })
      }
    } catch (err) {
      runInAction(() => {
        this.status = 'error'
        this.error = err instanceof Error ? err.message : '打开工程出错'
      })
    }
  }

  /**
   * 保存工程
   */
  async saveProject(): Promise<void> {
    if (!this.currentProject) {
      return
    }

    runInAction(() => {
      this.status = 'saving'
      this.error = ''
    })

    try {
      const response = await projectApi.saveProject(this.currentProject)

      if (response.success) {
        runInAction(() => {
          this.status = 'idle'
        })
      } else {
        runInAction(() => {
          this.status = 'error'
          this.error = response.message || '保存工程失败'
        })
      }
    } catch (err) {
      runInAction(() => {
        this.status = 'error'
        this.error = err instanceof Error ? err.message : '保存工程出错'
      })
    }
  }

  /**
   * 关闭工程
   */
  closeProject() {
    runInAction(() => {
      this.currentProject = null
      this.status = 'idle'
      this.error = ''
    })
  }

  /**
   * 清除错误
   */
  clearError() {
    runInAction(() => {
      this.error = ''
    })
  }

  // Computed属性

  /**
   * 是否有打开的工程
   */
  get hasProject(): boolean {
    return this.currentProject !== null
  }

  /**
   * 工程名称
   */
  get projectName(): string {
    return this.currentProject?.metadata.name || ''
  }

  /**
   * 是否加密工程
   */
  get isEncrypted(): boolean {
    return this.currentProject?.security.encrypted || false
  }

  /**
   * 是否正在加载
   */
  get isLoading(): boolean {
    return this.status === 'creating' || this.status === 'opening' || this.status === 'saving'
  }

  /**
   * 是否出错
   */
  get hasError(): boolean {
    return this.status === 'error' && this.error !== ''
  }
}

// 导出类和单例
export { ProjectStore }
export const projectStore = new ProjectStore()
