import { makeAutoObservable } from 'mobx'

/**
 * 编辑选项栏布局模式
 * - left: 左侧竖向布局（默认）
 * - top: 顶部横向布局
 */
export type EditOptionLayout = 'left' | 'top'

/**
 * LayoutStore - 布局配置管理
 * 负责管理编辑器的布局配置（编辑选项栏布局、侧边栏折叠等）
 */
export class LayoutStore {
  /** 编辑选项栏布局模式 */
  editOptionLayout: EditOptionLayout = 'left'

  /** 侧边栏是否折叠 */
  sidebarCollapsed: boolean = false

  constructor() {
    makeAutoObservable(this)

    // 从 localStorage 加载用户偏好
    this.loadSettings()
  }

  /**
   * 设置编辑选项栏布局模式
   * @param layout 布局模式
   */
  setEditOptionLayout(layout: EditOptionLayout): void {
    this.editOptionLayout = layout

    // 持久化到 localStorage
    this.saveSettings()
  }

  /**
   * 切换编辑选项栏布局模式
   */
  toggleEditOptionLayout(): void {
    const newLayout: EditOptionLayout = this.editOptionLayout === 'left' ? 'top' : 'left'
    this.setEditOptionLayout(newLayout)
  }

  /**
   * 切换侧边栏折叠状态
   */
  toggleSidebar(): void {
    this.sidebarCollapsed = !this.sidebarCollapsed

    // 持久化到 localStorage
    this.saveSettings()
  }

  /**
   * 设置侧边栏折叠状态
   * @param collapsed 是否折叠
   */
  setSidebarCollapsed(collapsed: boolean): void {
    this.sidebarCollapsed = collapsed

    // 持久化到 localStorage
    this.saveSettings()
  }

  /**
   * 保存设置到 localStorage
   */
  private saveSettings(): void {
    try {
      localStorage.setItem('editOptionLayout', this.editOptionLayout)
      localStorage.setItem('sidebarCollapsed', String(this.sidebarCollapsed))
    } catch (error) {
      console.error('Failed to save layout settings:', error)
    }
  }

  /**
   * 从 localStorage 加载设置
   */
  private loadSettings(): void {
    try {
      // 加载编辑选项栏布局
      const savedLayout = localStorage.getItem('editOptionLayout') as EditOptionLayout | null
      if (savedLayout && (savedLayout === 'left' || savedLayout === 'top')) {
        this.editOptionLayout = savedLayout
      }

      // 加载侧边栏折叠状态
      const savedCollapsed = localStorage.getItem('sidebarCollapsed')
      if (savedCollapsed !== null) {
        this.sidebarCollapsed = savedCollapsed === 'true'
      }
    } catch (error) {
      console.error('Failed to load layout settings:', error)
    }
  }

  /**
   * 重置为默认设置
   */
  resetSettings(): void {
    this.editOptionLayout = 'left'
    this.sidebarCollapsed = false
    this.saveSettings()
  }

  /**
   * 获取当前布局配置摘要
   */
  get layoutSummary(): {
    editOptionLayout: EditOptionLayout
    sidebarCollapsed: boolean
  } {
    return {
      editOptionLayout: this.editOptionLayout,
      sidebarCollapsed: this.sidebarCollapsed,
    }
  }
}

// 创建单例实例
export const layoutStore = new LayoutStore()
