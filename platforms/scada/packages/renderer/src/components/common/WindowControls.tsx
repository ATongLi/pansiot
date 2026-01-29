import React, { useEffect, useState } from 'react'
import { getElectronAPI } from '@/utils/electron'
import './WindowControls.module.css'

/**
 * WindowControls 窗口控制按钮组件
 * 提供最小化、最大化、关闭窗口的控制按钮
 */
const WindowControls: React.FC = () => {
  const [isMaximized, setIsMaximized] = useState(false)

  useEffect(() => {
    // 监听窗口最大化状态变化
    const electronAPI = getElectronAPI()

    const handleMaximized = () => setIsMaximized(true)
    const handleUnmaximized = () => setIsMaximized(false)

    // 获取初始状态
    electronAPI?.window.isMaximized?.().then(setIsMaximized).catch(() => {})

    // 监听状态变化
    electronAPI?.on('window:maximized', handleMaximized)
    electronAPI?.on('window:unmaximized', handleUnmaximized)

    return () => {
      electronAPI?.off('window:maximized', handleMaximized)
      electronAPI?.off('window:unmaximized', handleUnmaximized)
    }
  }, [])

  /**
   * 最小化窗口
   */
  const handleMinimize = (): void => {
    const electronAPI = getElectronAPI()
    electronAPI?.window.minimize()
  }

  /**
   * 最大化/还原窗口
   */
  const handleMaximize = (): void => {
    const electronAPI = getElectronAPI()
    if (isMaximized) {
      electronAPI?.window.unmaximize()
    } else {
      electronAPI?.window.maximize()
    }
  }

  /**
   * 关闭窗口
   */
  const handleClose = (): void => {
    const electronAPI = getElectronAPI()
    electronAPI?.window.close()
  }

  return (
    <div className="window-controls">
      <button
        className="window-controls__button window-controls__button--minimize"
        onClick={handleMinimize}
        aria-label="最小化"
        title="最小化"
      >
        <span className="window-controls__icon">−</span>
      </button>
      <button
        className="window-controls__button window-controls__button--maximize"
        onClick={handleMaximize}
        aria-label={isMaximized ? '向下还原' : '最大化'}
        title={isMaximized ? '向下还原' : '最大化'}
      >
        {isMaximized ? (
          <span className="window-controls__icon window-controls__icon--restore">❐</span>
        ) : (
          <span className="window-controls__icon">□</span>
        )}
      </button>
      <button
        className="window-controls__button window-controls__button--close"
        onClick={handleClose}
        aria-label="关闭"
        title="关闭"
      >
        <span className="window-controls__icon">×</span>
      </button>
    </div>
  )
}

export default WindowControls
