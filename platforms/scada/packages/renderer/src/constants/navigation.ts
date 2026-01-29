import { NavItem } from '@/types/navigation'
import { NAV_ICONS } from './icons'  // CR-002: 使用SVG线框图标

/**
 * 导航项配置
 * CR-002: 图标从Emoji改为SVG线框风格
 */
export const NAVIGATION_ITEMS: NavItem[] = [
  { id: 'home', label: '首页', icon: NAV_ICONS.home, path: '/' },
  { id: 'local', label: '本地', icon: NAV_ICONS.local, path: '/local' },
  { id: 'cloud', label: '云端', icon: NAV_ICONS.cloud, path: '/cloud' },
  { id: 'tools', label: '工具', icon: NAV_ICONS.tools, path: '/tools' },
  { id: 'user', label: 'User', icon: NAV_ICONS.user, path: '/user' },
]

export const DEFAULT_NAV_ITEM = 'home'
