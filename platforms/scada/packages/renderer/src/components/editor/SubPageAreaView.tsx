import React from 'react'
import type { SubPageTab } from '@/store'
import { SubPageTabs } from './tabs/SubPageTabs'
import { Canvas } from './canvas/Canvas'
import './SubPageArea.module.css'

export interface SubPageAreaViewProps {
  /** 子页面标签列表 */
  tabs: SubPageTab[]
  /** 当前激活的标签 ID */
  activeTab: string | null
  /** 标签切换回调 */
  onTabChange: (tabId: string) => void
  /** 标签关闭回调 */
  onTabClose: (tabId: string) => void
  /** 标签添加回调 */
  onTabAdd: () => void
  /** 组件拖放回调 */
  onDropComponent?: (component: any, x: number, y: number) => void
  /** 额外的CSS类名 */
  className?: string
}

/**
 * SubPageAreaView 子页面内容区视图组件（纯展示）
 * FE-009-06: 使用 Container/Presenter 模式重构
 *
 * 结构：
 * - 顶部：SubPageTabs（内部标签栏）
 * - 主体：Canvas（画布区域）
 *
 * 设计模式：
 * - 纯展示组件，不直接访问 store
 * - 所有数据和回调通过 props 传入
 * - 可在任何上下文中重用
 */
const SubPageAreaView: React.FC<SubPageAreaViewProps> = ({
  tabs,
  activeTab,
  onTabChange,
  onTabClose,
  onTabAdd,
  onDropComponent,
  className = '',
}) => {
  return (
    <div className={`sub-page-area ${className}`.trim()}>
      {/* 子页面标签栏 */}
      <SubPageTabs
        tabs={tabs}
        activeTab={activeTab || ''}
        onTabChange={onTabChange}
        onTabClose={onTabClose}
        onTabAdd={onTabAdd}
      />

      {/* 画布区域 */}
      <div className="sub-page-area__canvas-container">
        <Canvas onDropComponent={onDropComponent} />
      </div>
    </div>
  )
}

export default SubPageAreaView
