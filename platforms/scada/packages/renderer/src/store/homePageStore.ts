import { makeAutoObservable } from 'mobx'
import { NAV_ICONS } from '@/constants/icons'

/**
 * 首页导航项类型
 */
export type HomePageNavItem = 'home' | 'local' | 'cloud' | 'tools'

/**
 * 首页导航项配置
 * 与 NavItem 类型保持兼容
 */
export interface HomePageNavItemConfig {
  id: HomePageNavItem
  label: string
  icon: string
  path?: string // 可选的 path 属性
}

/**
 * 首页导航项列表
 * 使用 NAV_ICONS 中的 SVG 图标
 */
export const HOME_PAGE_NAV_ITEMS: HomePageNavItemConfig[] = [
  { id: 'home', label: '首页', icon: NAV_ICONS.home },
  { id: 'local', label: '本地', icon: NAV_ICONS.local },
  { id: 'cloud', label: '云端', icon: NAV_ICONS.cloud },
  { id: 'tools', label: '工具', icon: NAV_ICONS.tools },
]

/**
 * HomePageStore 首页状态管理
 * 管理首页内的导航状态（首页/本地/云端/工具）
 * 独立于 uiStore，避免状态冲突
 */
export class HomePageStore {
  /**
   * 当前激活的导航项
   */
  activeNavItem: HomePageNavItem = 'home'

  /**
   * 导航项列表
   */
  navigationItems = HOME_PAGE_NAV_ITEMS

  constructor() {
    makeAutoObservable(this)
  }

  /**
   * 设置激活的导航项
   */
  setActiveNavItem(itemId: HomePageNavItem) {
    const item = this.navigationItems.find(i => i.id === itemId)
    if (item) {
      this.activeNavItem = itemId
    }
  }
}

export const homePageStore = new HomePageStore()
