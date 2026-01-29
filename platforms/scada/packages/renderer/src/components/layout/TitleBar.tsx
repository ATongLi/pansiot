/**
 * 自定义标题栏组件
 * 替代Electron默认标题栏
 * 布局：[Logo] [标题] [标签区域...] [扩展按钮] [窗口控制]
 */

import React from 'react'
import WindowControls from '@/components/common/WindowControls'
import type { Tab } from '@/store'
import './TitleBar.css'

interface TitleBarProps {
  /**
   * 当前激活的标签页ID
   */
  activeTab?: string

  /**
   * 标签页列表
   */
  tabs?: Tab[]

  /**
   * 标签页切换回调
   */
  onTabChange?: (tabId: string) => void

  /**
   * 标签页关闭回调
   */
  onTabClose?: (tabId: string) => void
}

/**
 * TitleBar 标题栏组件
 */
const TitleBar: React.FC<TitleBarProps> = ({
  activeTab = 'home',
  tabs = [
    { id: 'home', title: '首页', type: 'home', closable: false },
  ],
  onTabChange,
  onTabClose,
}) => {
  return (
    <div className="title-bar">
      {/* Logo图标区域 */}
      <div className="title-bar__logo">
        <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          {/* PanTools Logo - 六边形图标 */}
          <path
            d="M12 2L2 7V17L12 22L22 17V7L12 2Z"
            fill="url(#gradient)"
            stroke="currentColor"
            strokeWidth="1.5"
          />
          <path
            d="M12 6L6 9V15L12 18L18 15V9L12 6Z"
            fill="currentColor"
            opacity="0.3"
          />
          <defs>
            <linearGradient id="gradient" x1="2" y1="2" x2="22" y2="22">
              <stop offset="0%" stopColor="#2196F3" />
              <stop offset="100%" stopColor="#64B5F6" />
            </linearGradient>
          </defs>
        </svg>
      </div>

      {/* 应用标题 */}
      <div className="title-bar__app-title">PanTools Scada</div>

      {/* 标签页区域 */}
      <div className="title-bar__tabs-section">
        <div className="title-bar__tabs">
          {tabs.map((tab) => (
            <div
              key={tab.id}
              className={`title-bar__tab ${activeTab === tab.id ? 'title-bar__tab--active' : ''}`}
              onClick={() => onTabChange?.(tab.id)}
            >
              {/* 只在首页标签显示图标 */}
              {tab.type === 'home' && (
                <svg className="title-bar__tab-icon title-bar__tab-icon--home" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                  {/* 圆润可爱的房子图标 */}
                  {/* 房子外框 - 更圆润的屋顶 */}
                  <path
                    d="M3 9L9 4L15 9V14H3V9Z"
                    stroke="currentColor"
                    strokeWidth="1.3"
                    fill="none"
                    strokeLinejoin="round"
                    strokeLinecap="round"
                  />
                  {/* 门 - 一竖，加粗 */}
                  <path
                    d="M9 14V10"
                    stroke="currentColor"
                    strokeWidth="1.3"
                    strokeLinecap="round"
                  />
                </svg>
              )}

              <span className="title-bar__tab-title">{tab.title}</span>

              {tab.closable && (
                <span
                  className="title-bar__tab-close"
                  onClick={(e) => {
                    e.stopPropagation()
                    onTabClose?.(tab.id)
                  }}
                >
                  <svg viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
                    {/* 极简叉 - 两条细线 */}
                    <path
                      d="M3 3L9 9M9 3L3 9"
                      stroke="currentColor"
                      strokeWidth="1"
                      strokeLinecap="round"
                    />
                  </svg>
                </span>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* 扩展功能区域 */}
      <div className="title-bar__extensions">
        <button
          className="title-bar__ext-button title-bar__ext-button--settings"
          title="设置"
        >
          <span className="title-bar__ext-icon">⚙️</span>
        </button>
      </div>

      {/* 窗口控制按钮 */}
      <WindowControls />
    </div>
  )
}

export default TitleBar
