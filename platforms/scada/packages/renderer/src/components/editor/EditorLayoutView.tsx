import React from 'react'
import { StatusBar } from './statusbar/StatusBar'
import { TopToolbar } from './toolbar/TopToolbar'
import { SubToolbar } from './toolbar/SubToolbar'
import { Canvas } from './canvas/Canvas'
import { SubPageTabs } from './tabs/SubPageTabs'
import { ProjectPanel } from './sidebar/ProjectPanel'
import { ScreenPanel } from './sidebar/ScreenPanel'
import { ComponentPanel } from './sidebar/ComponentPanel'
import type { SubPageTab } from '@/store'
import { ToolbarCategory } from './toolbar/ToolbarTypes'
import './EditorLayout.css'

export interface EditorLayoutViewProps {
  // ==================== 状态 ====================
  /** 当前激活的左侧边栏标签 */
  leftSidebarActiveTab: 'project' | 'screen' | 'component'
  /** 右侧边栏是否可见 */
  rightSidebarVisible: boolean
  /** 当前激活的右侧边栏标签 */
  rightSidebarActiveTab: 'property' | 'layer'
  /** 子页面标签列表 */
  subPageTabs: SubPageTab[]
  /** 当前激活的子页面标签 ID */
  activeSubTab: string | null

  // ==================== 工具栏状态 ====================
  /** 子工具栏是否可见 */
  subToolbarVisible?: boolean
  /** 当前激活的工具栏分类 */
  activeToolbarCategory?: ToolbarCategory

  // ==================== 回调 - 文件操作 ====================
  onNewProject: () => void
  onOpenProject: () => void
  onSaveProject: () => void

  // ==================== 回调 - 侧边栏操作 ====================
  onSetLeftSidebarTab: (tab: 'project' | 'screen' | 'component') => void
  onSetRightSidebarTab: (tab: 'property' | 'layer') => void

  // ==================== 回调 - 标签页操作 ====================
  onTabChange: (tabId: string) => void
  onTabClose: (tabId: string) => void
  onTabAdd: () => void

  // ==================== 回调 - 画布操作 ====================
  onDropComponent: (component: any, x: number, y: number) => void

  // ==================== 回调 - 工具栏操作 ====================
  /** 切换子工具栏可见性 */
  onToggleSubToolbar?: () => void
  /** 设置激活的工具栏分类 */
  onSetActiveToolbarCategory?: (category: ToolbarCategory) => void

  /** 额外的CSS类名 */
  className?: string
}

/**
 * EditorLayoutView 编辑器布局视图组件（纯展示）
 * FE-009-03: 工程编辑器布局重构
 *
 * 布局结构：
 * ┌────────────────────────────────────────────────────────────┐
 * │ Top Toolbar (32px)                                         │
 * ├────────────────────────────────────────────────────────────┤
 * │ Sub Toolbar (64px) - Dynamic/Fixed mode                   │
 * ├─────────┬──────────────────────────────────┬───────────────┤
 * │         │ Sub Page Tabs (32px)             │               │
 * │ Left    ├──────────────────────────────────┤ Right         │
 * │ Sidebar │                                  │ Sidebar       │
 * │ (280px) │        Canvas Area               │ (280px)       │
 * │         │        (flex: 1)                 │               │
 * │         │                                  │               │
 * ├─────────┴──────────────────────────────────┴───────────────┤
 * │ Status Bar (24px)                                            │
 * └────────────────────────────────────────────────────────────┘
 *
 * 设计模式：
 * - 纯展示组件，不直接访问 store
 * - 所有数据和回调通过 props 传入
 * - 专注于 UI 渲染和布局
 */
export const EditorLayoutView: React.FC<EditorLayoutViewProps> = ({
  leftSidebarActiveTab,
  rightSidebarVisible,
  rightSidebarActiveTab,
  subPageTabs,
  activeSubTab,
  subToolbarVisible = true,
  activeToolbarCategory = ToolbarCategory.PROJECT,
  onNewProject,
  onOpenProject,
  onSaveProject,
  onSetLeftSidebarTab,
  onSetRightSidebarTab,
  onTabChange,
  onTabClose,
  onTabAdd,
  onDropComponent,
  onToggleSubToolbar,
  onSetActiveToolbarCategory,
  className = '',
}) => {
  // ==========================================
  // Render - Left Sidebar Tab Buttons
  // ==========================================

  const renderLeftSidebarTabs = () => (
    <div className="sidebar-tabs">
      <button
        className={`sidebar-tab ${leftSidebarActiveTab === 'project' ? 'sidebar-tab--active' : ''}`}
        onClick={() => onSetLeftSidebarTab('project')}
      >
        工程
      </button>
      <button
        className={`sidebar-tab ${leftSidebarActiveTab === 'screen' ? 'sidebar-tab--active' : ''}`}
        onClick={() => onSetLeftSidebarTab('screen')}
      >
        画面
      </button>
      <button
        className={`sidebar-tab ${leftSidebarActiveTab === 'component' ? 'sidebar-tab--active' : ''}`}
        onClick={() => onSetLeftSidebarTab('component')}
      >
        组件
      </button>
    </div>
  )

  // ==========================================
  // Render - Right Sidebar
  // ==========================================

  const renderRightSidebar = () => (
    <div
      className={`editor__right-sidebar ${
        !rightSidebarVisible ? 'editor__right-sidebar--hidden' : ''
      }`}
    >
      {/* 右侧边栏Tab */}
      <div className="sidebar-tabs">
        <button
          className={`sidebar-tab ${
            rightSidebarActiveTab === 'property' ? 'sidebar-tab--active' : ''
          }`}
          onClick={() => onSetRightSidebarTab('property')}
        >
          属性
        </button>
        <button
          className={`sidebar-tab ${
            rightSidebarActiveTab === 'layer' ? 'sidebar-tab--active' : ''
          }`}
          onClick={() => onSetRightSidebarTab('layer')}
        >
          图层
        </button>
      </div>

      {/* 属性面板 */}
      <div
        className={`sidebar-panel ${
          rightSidebarActiveTab !== 'property' ? 'sidebar-panel--hidden' : ''
        }`}
      >
        <div className="sidebar-panel__header">属性</div>
        <div className="sidebar-panel__content">
          {/* TODO: 集成 PropertyPanel (阶段4已跳过) */}
          <div className="editor-empty-state">
            <div className="editor-empty-state__icon">⚙️</div>
            <div className="editor-empty-state__text">未选中元素</div>
            <div className="editor-empty-state__hint">选中画布中的元素以编辑属性</div>
          </div>
        </div>
      </div>

      {/* 图层面板 */}
      <div
        className={`sidebar-panel ${
          rightSidebarActiveTab !== 'layer' ? 'sidebar-panel--hidden' : ''
        }`}
      >
        <div className="sidebar-panel__header">图层</div>
        <div className="sidebar-panel__content">
          {/* TODO: 集成 LayerPanel (阶段4已跳过) */}
          <div className="editor-empty-state">
            <div className="editor-empty-state__icon">📑</div>
            <div className="editor-empty-state__text">暂无图层</div>
            <div className="editor-empty-state__hint">添加组件后图层将显示在此</div>
          </div>
        </div>
      </div>
    </div>
  )

  // ==========================================
  // Main Render
  // ==========================================

  return (
    <div className={`editor-layout ${className}`}>
      {/* Top Toolbar */}
      <TopToolbar
        subToolbarVisible={subToolbarVisible}
        activeCategory={activeToolbarCategory}
        onToggleSubToolbar={onToggleSubToolbar}
        onSetActiveCategory={onSetActiveToolbarCategory}
        onNewProject={onNewProject}
        onOpenProject={onOpenProject}
        onSaveProject={onSaveProject}
      />

      {/* Sub Toolbar */}
      <SubToolbar
        visible={subToolbarVisible}
        activeCategory={activeToolbarCategory}
        onNewProject={onNewProject}
        onOpenProject={onOpenProject}
        onSaveProject={onSaveProject}
      />

      {/* Main Content Area */}
      <div className="editor__main-content">
        {/* Left Sidebar */}
        <div className="editor__left-sidebar">
          {/* Tab Buttons */}
          {renderLeftSidebarTabs()}

          {/* Panel Content */}
          <div className="sidebar-content">
            {leftSidebarActiveTab === 'project' && <ProjectPanel />}
            {leftSidebarActiveTab === 'screen' && <ScreenPanel />}
            {leftSidebarActiveTab === 'component' && <ComponentPanel />}
          </div>
        </div>

        {/* Canvas Area */}
        <div className="editor__canvas-area">
          {/* Sub Page Tabs */}
          <SubPageTabs
            tabs={subPageTabs}
            activeTab={activeSubTab || ''}
            onTabChange={onTabChange}
            onTabClose={onTabClose}
            onTabAdd={onTabAdd}
          />

          {/* Canvas */}
          <Canvas onDropComponent={onDropComponent} />
        </div>

        {/* Right Sidebar */}
        {renderRightSidebar()}
      </div>

      {/* Status Bar */}
      <StatusBar />
    </div>
  )
}

export default EditorLayoutView
