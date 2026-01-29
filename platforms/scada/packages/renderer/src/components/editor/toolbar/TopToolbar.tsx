/**
 * TopToolbar 顶部工具栏组件
 *
 * 功能：
 * - 6个主分类按钮（工程文件、通用、运行调试、工具、视图、帮助）
 * - 展开收缩子工具栏的控制按钮
 * - 使用线条图标
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { ToolbarCategory } from './ToolbarTypes'
import { getEditorStore } from '@/store'
import './TopToolbar.css'

export interface TopToolbarProps {
  className?: string
  subToolbarVisible?: boolean
  activeCategory?: ToolbarCategory
  onToggleSubToolbar?: () => void
  onSetActiveCategory?: (category: ToolbarCategory) => void
  /** 工程文件操作回调 */
  onNewProject?: () => void
  onOpenProject?: () => void
  onSaveProject?: () => void
}

/**
 * TopToolbar 组件
 */
export const TopToolbar: React.FC<TopToolbarProps> = observer(({
  className = '',
  subToolbarVisible = true,
  activeCategory = ToolbarCategory.PROJECT,
  onToggleSubToolbar,
  onSetActiveCategory,
  onNewProject,
  onOpenProject,
  onSaveProject,
}) => {
  const editorStore = getEditorStore()

  // ==========================================
  // Handlers
  // ==========================================

  const handleCategoryClick = (category: ToolbarCategory): void => {
    onSetActiveCategory?.(category)
    // 自动展开子工具栏
    if (!subToolbarVisible) {
      onToggleSubToolbar?.()
    }
  }

  // ==========================================
  // Render - 主分类按钮
  // ==========================================

  const renderCategoryButton = (
    category: ToolbarCategory,
    icon: string,
    label: string
  ) => (
    <button
      key={category}
      className={`toolbar-category-button ${
        activeCategory === category ? 'toolbar-category-button--active' : ''
      }`}
      onClick={() => handleCategoryClick(category)}
      title={label}
    >
      <span className="toolbar-category-button__label">{label}</span>
    </button>
  )

  // ==========================================
  // Main Render
  // ==========================================

  return (
    <div className={`top-toolbar ${className}`}>
      {/* 左侧：6个主分类按钮 */}
      <div className="top-toolbar__categories">
        {renderCategoryButton(ToolbarCategory.PROJECT, 'save-project', '工程文件')}
        {renderCategoryButton(ToolbarCategory.GENERAL, 'undo', '通用')}
        {renderCategoryButton(ToolbarCategory.DEBUG, 'compile', '运行调试')}
        {renderCategoryButton(ToolbarCategory.TOOLS, 'device-manager', '工具')}
        {renderCategoryButton(ToolbarCategory.VIEW, 'grid', '视图')}
        {renderCategoryButton(ToolbarCategory.HELP, 'help-docs', '帮助')}
      </div>

      {/* 右侧：展开/收缩子工具栏按钮 */}
      <div className="top-toolbar__controls">
        <button
          className={`toolbar-toggle-button ${subToolbarVisible ? 'toolbar-toggle-button--expanded' : ''}`}
          onClick={onToggleSubToolbar}
          title={subToolbarVisible ? '收缩子工具栏' : '展开子工具栏'}
          aria-label={subToolbarVisible ? '收缩子工具栏' : '展开子工具栏'}
        >
          <span className="toolbar-toggle-button__icon">
            {subToolbarVisible ? '▲' : '▼'}
          </span>
        </button>
      </div>
    </div>
  )
})

/**
 * 默认导出
 */
export default TopToolbar
