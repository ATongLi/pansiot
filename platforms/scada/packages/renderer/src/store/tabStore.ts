import { makeAutoObservable, computed } from 'mobx'

/**
 * Tab 接口定义
 * 表示应用级的标签页（首页标签、工程标签）
 */
export interface Tab {
  /** 标签唯一ID */
  id: string
  /** 标签标题 */
  title: string
  /** 标签图标（可选） */
  icon?: string
  /** 标签类型：home-首页, editor-工程编辑器 */
  type: 'home' | 'editor'
  /** 是否可关闭（首页不可关闭） */
  closable: boolean
  /** 工程文件路径（仅工程标签） */
  projectPath?: string
}

/**
 * TabStore - 应用级标签状态管理
 * 负责管理应用级别的所有标签（首页、工程编辑器等）
 */
export class TabStore {
  /** 所有标签列表 */
  tabs: Tab[] = []

  /** 当前激活的标签ID */
  activeTab: string = ''

  constructor() {
    makeAutoObservable(this)

    // 初始化时添加首页标签
    this.initHomeTab()
  }

  /**
   * 初始化首页标签
   */
  private initHomeTab(): void {
    const homeTab: Tab = {
      id: 'tab-home',
      title: '首页',
      type: 'home',
      closable: false, // 首页不可关闭
    }
    this.tabs.push(homeTab)
    this.activeTab = homeTab.id
  }

  /**
   * 添加新标签
   * @param tab 标签信息（不包含id）
   * @returns 新创建的标签ID
   */
  addTab(tab: Omit<Tab, 'id'>): string {
    const newTab: Tab = {
      ...tab,
      id: `tab-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    }
    this.tabs.push(newTab)
    this.activeTab = newTab.id
    return newTab.id
  }

  /**
   * 关闭标签
   * @param tabId 要关闭的标签ID
   */
  closeTab(tabId: string): void {
    const index = this.tabs.findIndex(t => t.id === tabId)

    // 标签不存在
    if (index === -1) {
      console.warn(`Tab ${tabId} not found`)
      return
    }

    const tabToClose = this.tabs[index]

    // 不可关闭的标签
    if (!tabToClose.closable) {
      console.warn(`Tab ${tabId} is not closable`)
      return
    }

    // 移除标签
    this.tabs.splice(index, 1)

    // 如果关闭的是当前激活标签，切换到其他标签
    if (this.activeTab === tabId) {
      if (this.tabs.length > 0) {
        // 优先切换到右侧标签，如果没有则切换到左侧
        const newIndex = Math.min(index, this.tabs.length - 1)
        this.activeTab = this.tabs[newIndex].id
      } else {
        // 理论上不应该到这里（至少有首页标签）
        console.error('No tabs available after closing')
        this.activeTab = ''
      }
    }
  }

  /**
   * 设置激活标签
   * @param tabId 要激活的标签ID
   */
  setActiveTab(tabId: string): void {
    const tab = this.tabs.find(t => t.id === tabId)
    if (tab) {
      this.activeTab = tabId
    } else {
      console.warn(`Tab ${tabId} not found`)
    }
  }

  /**
   * 获取当前激活的标签对象
   */
  get activeTabObj(): Tab | undefined {
    return this.tabs.find(t => t.id === this.activeTab)
  }

  /**
   * 获取所有可关闭的标签（非首页）
   */
  get closableTabs(): Tab[] {
    return this.tabs.filter(t => t.closable)
  }

  /**
   * 根据项目路径查找标签
   * @param projectPath 工程文件路径
   */
  findTabByProjectPath(projectPath: string): Tab | undefined {
    return this.tabs.find(t => t.projectPath === projectPath)
  }
}

// 创建单例实例
export const tabStore = new TabStore()
