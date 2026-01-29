import React from 'react'
import Tab from './Tab'
import WindowControls from './WindowControls'
import type { Tab as TabType } from '@/store'
import './TabBar.module.css'

export interface TabBarProps {
  /** 标签列表 */
  tabs: TabType[]
  /** 当前激活的标签ID */
  activeTab: string
  /** 标签切换回调 */
  onTabChange?: (tabId: string) => void
  /** 标签关闭回调 */
  onTabClose?: (tabId: string) => void
  /** 是否显示窗口控制按钮 */
  showWindowControls?: boolean
  /** 是否显示品牌区域（Logo+标题） */
  showBrand?: boolean
}

/**
 * TabBar 应用级标签栏组件
 * 显示所有应用级标签（首页、工程等）和窗口控制按钮
 */
const TabBar: React.FC<TabBarProps> = ({
  tabs,
  activeTab,
  onTabChange,
  onTabClose,
  showWindowControls = true,
  showBrand = true,
}) => {
  return (
    <div className="tabbar">
      {/* 左侧：Logo + 应用标题（可选） */}
      {showBrand && (
        <div className="tabbar__brand">
          <span className="tabbar__logo">PanTools</span>
          <span className="tabbar__title">Scada 组态软件</span>
        </div>
      )}

      {/* 中间：标签列表 */}
      <div className="tabbar__tabs">
        {tabs.map((tab) => (
          <Tab
            key={tab.id}
            id={tab.id}
            title={tab.title}
            icon={tab.icon}
            active={tab.id === activeTab}
            closable={tab.closable}
            onClick={onTabChange}
            onClose={onTabClose}
          />
        ))}
      </div>

      {/* 右侧：窗口控制按钮 */}
      {showWindowControls && <WindowControls />}
    </div>
  )
}

export default TabBar
