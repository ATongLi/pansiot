import { makeAutoObservable } from 'mobx'
import { NAVIGATION_ITEMS, DEFAULT_NAV_ITEM } from '@/constants/navigation'
import type { WindowState, WindowAction } from '@/types/window'

export class UIStore {
  // Navigation State
  activeNavItem: string = DEFAULT_NAV_ITEM
  navigationItems = NAVIGATION_ITEMS

  // Window State
  windowState: WindowState = {
    isMaximized: false,
    isFullscreen: false,
  }

  constructor() {
    makeAutoObservable(this)
  }

  setActiveNavItem(itemId: string) {
    const item = this.navigationItems.find(i => i.id === itemId)
    if (item) {
      this.activeNavItem = itemId
    }
  }

  async handleWindowAction(action: WindowAction) {
    // TODO(依赖): Electron Main Process
    // 需要集成 window.electronAPI.send(action)
    // 当前状态: Mock实现
    console.log('Window action:', action)

    switch (action) {
      case 'minimize':
        break
      case 'maximize':
        this.windowState.isMaximized = !this.windowState.isMaximized
        break
      case 'restore':
        this.windowState.isMaximized = false
        break
      case 'close':
        break
    }
  }
}

export const uiStore = new UIStore()
