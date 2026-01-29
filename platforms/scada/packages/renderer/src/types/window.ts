export interface WindowState {
  isMaximized: boolean
  isFullscreen: boolean
}

export type WindowAction = 'minimize' | 'maximize' | 'close' | 'restore'
