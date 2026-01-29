import React from 'react'
import type { EditOptionLayout } from '@/store'
import './EditOptionsBar.module.css'

export interface EditOptionsBarViewProps {
  /** 当前布局方向 */
  layout: EditOptionLayout
  /** 布局切换回调 */
  onToggleLayout: () => void
  /** 子组件内容（左侧或顶部布局的内容） */
  children: React.ReactNode
  /** 额外的CSS类名 */
  className?: string
}

/**
 * EditOptionsBarView 编辑选项栏视图组件（纯展示）
 * FE-009-05: 可配置的编辑选项栏（支持左侧/顶部布局）
 *
 * 功能：
 * - 根据布局方向动态切换布局方向
 * - 左侧布局：垂直排列，固定宽度
 * - 顶部布局：水平排列，固定高度
 * - 通过切换按钮实时改变布局方向
 *
 * 设计模式：
 * - 纯展示组件，不直接访问 store
 * - 所有数据和回调通过 props 传入
 * - 可在任何上下文中重用
 */
const EditOptionsBarView: React.FC<EditOptionsBarViewProps> = ({
  layout,
  onToggleLayout,
  children,
  className = '',
}) => {
  return (
    <div
      className={`edit-options-bar edit-options-bar--${layout} ${className}`.trim()}
    >
      {/* 布局切换按钮 */}
      <button
        className="edit-options-bar__toggle"
        onClick={onToggleLayout}
        title={`切换至${layout === 'left' ? '顶部' : '左侧'}布局`}
        aria-label={`切换至${layout === 'left' ? '顶部' : '左侧'}布局`}
      >
        {layout === 'left' ? '⬅️' : '⬆️'}
      </button>

      {/* 子组件内容区域 */}
      <div className="edit-options-bar__content">
        {children}
      </div>
    </div>
  )
}

export default EditOptionsBarView
