import React from 'react'
import './Tab.module.css'

export interface TabProps {
  /** 标签唯一ID */
  id: string
  /** 标签标题 */
  title: string
  /** 标签图标（可选） */
  icon?: string
  /** 是否激活 */
  active: boolean
  /** 是否可关闭 */
  closable: boolean
  /** 点击事件 */
  onClick?: (id: string) => void
  /** 关闭事件 */
  onClose?: (id: string) => void
}

/**
 * Tab 标签页组件
 * 显示单个标签，支持点击切换和关闭操作
 */
const Tab: React.FC<TabProps> = ({
  id,
  title,
  icon,
  active,
  closable,
  onClick,
  onClose,
}) => {
  /**
   * 标签点击处理
   */
  const handleClick = (): void => {
    onClick?.(id)
  }

  /**
   * 关闭按钮点击处理
   * 阻止事件冒泡，避免触发标签点击
   */
  const handleClose = (e: React.MouseEvent): void => {
    e.stopPropagation()
    onClose?.(id)
  }

  return (
    <div
      className={`tab ${active ? 'tab--active' : ''}`}
      onClick={handleClick}
      role="tab"
      tabIndex={0}
      aria-selected={active}
    >
      {/* 图标 */}
      {icon && <span className="tab__icon">{icon}</span>}

      {/* 标题 */}
      <span className="tab__title">{title}</span>

      {/* 关闭按钮 */}
      {closable && (
        <span
          className="tab__close"
          onClick={handleClose}
          role="button"
          aria-label="关闭标签"
          tabIndex={0}
        >
          ×
        </span>
      )}
    </div>
  )
}

export default Tab
