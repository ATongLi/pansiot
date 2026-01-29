/**
 * 最近工程列表项组件
 * 显示单个工程的简要信息
 */

import React, { useState, useRef, useEffect } from 'react'
import type { RecentProject } from '@/types/project'
import { formatRelativeTime } from '@/utils/dateFormat'
import './ProjectListItem.css'

interface ProjectListItemProps {
  project: RecentProject
  isActive?: boolean
  onClick: (project: RecentProject) => void
  onOpen?: (project: RecentProject) => void
  onShowInExplorer?: (project: RecentProject) => void
  onRemove?: (project: RecentProject) => void
  onCopyPath?: (project: RecentProject) => void
}

/**
 * ProjectListItem 组件
 */
const ProjectListItem: React.FC<ProjectListItemProps> = ({
  project,
  isActive = false,
  onClick,
  onOpen,
  onShowInExplorer,
  onRemove,
  onCopyPath
}) => {
  const [showContextMenu, setShowContextMenu] = useState(false)
  const contextMenuRef = useRef<HTMLDivElement>(null)

  // 处理右键菜单
  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setShowContextMenu(true)
  }

  // 关闭右键菜单
  const closeContextMenu = () => {
    setShowContextMenu(false)
  }

  // 点击外部关闭菜单
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        contextMenuRef.current &&
        !contextMenuRef.current.contains(event.target as Node)
      ) {
        closeContextMenu()
      }
    }

    if (showContextMenu) {
      document.addEventListener('mousedown', handleClickOutside)
      return () => document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [showContextMenu])

  // 处理菜单项点击
  const handleMenuAction = (action: () => void) => {
    action()
    closeContextMenu()
  }

  return (
    <div
      className={`project-list-item ${isActive ? 'project-list-item--active' : ''}`}
      onClick={() => onClick(project)}
      onContextMenu={handleContextMenu}
      tabIndex={0}
      onKeyDown={e => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          onClick(project)
        }
      }}
    >
      {/* 工程图标 */}
      <div className="project-list-item__icon">
        📁
      </div>

      {/* 工程信息 */}
      <div className="project-list-item__info">
        <div className="project-list-item__name">
          {project.name}
        </div>
        <div className="project-list-item__meta">
          {project.category && (
            <span className="project-list-item__category">
              {project.category}
            </span>
          )}
          <span className="project-list-item__time">
            {formatRelativeTime(project.lastOpenedDate)}
          </span>
        </div>
        <div className="project-list-item__path" title={project.filePath}>
          {project.filePath}
        </div>
      </div>

      {/* 加密状态图标 */}
      {project.isEncrypted && (
        <div className="project-list-item__encrypted" title="已加密">
          🔒
        </div>
      )}

      {/* 右键菜单 */}
      {showContextMenu && (
        <div
          ref={contextMenuRef}
          className="project-list-item__context-menu"
          onClick={e => e.stopPropagation()}
        >
          {onOpen && (
            <div
              className="context-menu__item"
              onClick={() => handleMenuAction(() => onOpen(project))}
            >
              打开工程
            </div>
          )}
          {onShowInExplorer && (
            <div
              className="context-menu__item"
              onClick={() => handleMenuAction(() => onShowInExplorer(project))}
            >
              在文件管理器中显示
            </div>
          )}
          {onCopyPath && (
            <div
              className="context-menu__item"
              onClick={() => handleMenuAction(() => onCopyPath(project))}
            >
              复制工程路径
            </div>
          )}
          <div className="context-menu__divider" />
          {onRemove && (
            <div
              className="context-menu__item context-menu__item--danger"
              onClick={() => handleMenuAction(() => onRemove(project))}
            >
              从列表中移除
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default ProjectListItem
