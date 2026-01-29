import { makeAutoObservable } from 'mobx'

/**
 * 主工具栏选项类型
 */
export type ToolOptionType = 'file' | 'edit' | 'view' | 'tools' | 'help'

/**
 * 子工具栏显示模式
 * - fixed: 固定模式，子工具栏一直显示
 * - dynamic: 动态模式，选中时显示，失去焦点后自动隐藏
 */
export type SubToolbarMode = 'fixed' | 'dynamic'

/**
 * ToolbarStore - 工具栏状态管理
 * 负责管理主工具栏和子工具栏的状态
 */
export class ToolbarStore {
  /** 当前选中的主工具栏选项 */
  selectedToolOption: ToolOptionType = 'file'

  /** 子工具栏显示模式 */
  subToolbarMode: SubToolbarMode = 'dynamic'

  /** 子工具栏是否可见 */
  subToolbarVisible: boolean = false

  /** 子工具栏是否获得焦点（用于动态模式） */
  private isSubToolbarFocused: boolean = false

  /** 失去焦点延迟定时器 */
  private hideTimer: ReturnType<typeof setTimeout> | null = null

  constructor() {
    makeAutoObservable(this)

    // 从 localStorage 加载用户偏好
    this.loadSettings()
  }

  /**
   * 选择主工具栏选项
   * @param optionId 工具选项ID
   */
  selectToolOption(optionId: ToolOptionType): void {
    this.selectedToolOption = optionId

    // 选中工具选项后，显示子工具栏
    this.subToolbarVisible = true

    // 清除之前的隐藏定时器
    this.clearHideTimer()
  }

  /**
   * 设置子工具栏显示模式
   * @param mode 显示模式
   */
  setSubToolbarMode(mode: SubToolbarMode): void {
    this.subToolbarMode = mode

    // 固定模式下，子工具栏一直显示
    if (mode === 'fixed') {
      this.subToolbarVisible = true
    }

    // 保存到 localStorage
    this.saveSettings()
  }

  /**
   * 切换子工具栏显示模式
   */
  toggleSubToolbarMode(): void {
    const newMode: SubToolbarMode = this.subToolbarMode === 'fixed' ? 'dynamic' : 'fixed'
    this.setSubToolbarMode(newMode)
  }

  /**
   * 子工具栏获得焦点
   */
  focusSubToolbar(): void {
    this.isSubToolbarFocused = true
    this.clearHideTimer()
  }

  /**
   * 子工具栏失去焦点
   */
  blurSubToolbar(): void {
    this.isSubToolbarFocused = false

    // 动态模式下，延迟隐藏子工具栏
    if (this.subToolbarMode === 'dynamic') {
      this.clearHideTimer()
      this.hideTimer = setTimeout(() => {
        // 确认仍然失去焦点后才隐藏
        if (!this.isSubToolbarFocused) {
          this.subToolbarVisible = false
        }
      }, 200) // 200ms 延迟，避免误触
    }
  }

  /**
   * 手动隐藏子工具栏
   */
  hideSubToolbar(): void {
    if (this.subToolbarMode === 'dynamic') {
      this.subToolbarVisible = false
    }
  }

  /**
   * 显示子工具栏
   */
  showSubToolbar(): void {
    this.subToolbarVisible = true
  }

  /**
   * 清除隐藏定时器
   */
  private clearHideTimer(): void {
    if (this.hideTimer) {
      clearTimeout(this.hideTimer)
      this.hideTimer = null
    }
  }

  /**
   * 保存设置到 localStorage
   */
  private saveSettings(): void {
    try {
      localStorage.setItem('subToolbarMode', this.subToolbarMode)
    } catch (error) {
      console.error('Failed to save toolbar settings:', error)
    }
  }

  /**
   * 从 localStorage 加载设置
   */
  private loadSettings(): void {
    try {
      const savedMode = localStorage.getItem('subToolbarMode') as SubToolbarMode | null
      if (savedMode && (savedMode === 'fixed' || savedMode === 'dynamic')) {
        this.subToolbarMode = savedMode
      }
    } catch (error) {
      console.error('Failed to load toolbar settings:', error)
    }
  }

  /**
   * 重置为默认设置
   */
  resetSettings(): void {
    this.subToolbarMode = 'dynamic'
    this.selectedToolOption = 'file'
    this.subToolbarVisible = false
    this.saveSettings()
  }
}

// 创建单例实例
export const toolbarStore = new ToolbarStore()
