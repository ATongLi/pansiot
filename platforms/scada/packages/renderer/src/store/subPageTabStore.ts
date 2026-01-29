import { makeAutoObservable } from 'mobx'

/**
 * SubPageTab 接口定义
 * 表示编辑内容页内部的标签（画面编辑、组件配置等）
 */
export interface SubPageTab {
  /** 标签唯一ID */
  id: string
  /** 标签标题 */
  title: string
  /** 标签图标（可选） */
  icon?: string
  /** 标签内容类型 */
  contentType: 'screen' | 'component' | 'property' | 'layer' | 'custom'
  /** 关联的数据ID（如画面ID、组件ID） */
  dataId?: string
  /** 是否可关闭 */
  closable: boolean
}

/**
 * SubPageTabStore - 编辑内容页内部标签状态管理
 * 负责管理工程编辑器内部的子页面标签（画面、组件、属性等）
 */
export class SubPageTabStore {
  /** 内部标签列表 */
  subPageTabs: SubPageTab[] = []

  /** 当前激活的内部标签ID */
  activeSubTab: string = ''

  constructor() {
    makeAutoObservable(this)
  }

  /**
   * 添加新的内部标签
   * @param tab 内部标签信息（不包含id）
   * @returns 新创建的标签ID
   */
  addSubTab(tab: Omit<SubPageTab, 'id'>): string {
    // 检查是否已存在相同内容的标签
    const existingTab = this.subPageTabs.find(
      t => t.contentType === tab.contentType && t.dataId === tab.dataId
    )

    // 如果已存在，直接激活该标签
    if (existingTab) {
      this.activeSubTab = existingTab.id
      return existingTab.id
    }

    // 创建新标签
    const newTab: SubPageTab = {
      ...tab,
      id: `subtab-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    }
    this.subPageTabs.push(newTab)
    this.activeSubTab = newTab.id
    return newTab.id
  }

  /**
   * 关闭内部标签
   * @param tabId 要关闭的标签ID
   */
  closeSubTab(tabId: string): void {
    const index = this.subPageTabs.findIndex(t => t.id === tabId)

    // 标签不存在
    if (index === -1) {
      console.warn(`SubPageTab ${tabId} not found`)
      return
    }

    const tabToClose = this.subPageTabs[index]

    // 不可关闭的标签
    if (!tabToClose.closable) {
      console.warn(`SubPageTab ${tabId} is not closable`)
      return
    }

    // 移除标签
    this.subPageTabs.splice(index, 1)

    // 如果关闭的是当前激活标签，切换到其他标签
    if (this.activeSubTab === tabId) {
      if (this.subPageTabs.length > 0) {
        // 优先切换到右侧标签，如果没有则切换到左侧
        const newIndex = Math.min(index, this.subPageTabs.length - 1)
        this.activeSubTab = this.subPageTabs[newIndex].id
      } else {
        // 没有标签了
        this.activeSubTab = ''
      }
    }
  }

  /**
   * 设置激活的内部标签
   * @param tabId 要激活的标签ID
   */
  setActiveSubTab(tabId: string): void {
    const tab = this.subPageTabs.find(t => t.id === tabId)
    if (tab) {
      this.activeSubTab = tabId
    } else {
      console.warn(`SubPageTab ${tabId} not found`)
    }
  }

  /**
   * 获取当前激活的内部标签对象
   */
  get activeSubTabObj(): SubPageTab | undefined {
    return this.subPageTabs.find(t => t.id === this.activeSubTab)
  }

  /**
   * 清空所有内部标签
   * 在切换工程或关闭工程时调用
   */
  clearAllSubTabs(): void {
    this.subPageTabs = []
    this.activeSubTab = ''
  }

  /**
   * 根据内容类型和数据ID查找标签
   * @param contentType 内容类型
   * @param dataId 数据ID
   */
  findTabByContent(contentType: SubPageTab['contentType'], dataId: string): SubPageTab | undefined {
    return this.subPageTabs.find(t => t.contentType === contentType && t.dataId === dataId)
  }

  /**
   * 获取特定类型的所有标签
   * @param contentType 内容类型
   */
  getTabsByType(contentType: SubPageTab['contentType']): SubPageTab[] {
    return this.subPageTabs.filter(t => t.contentType === contentType)
  }
}

// 创建单例实例
export const subPageTabStore = new SubPageTabStore()
