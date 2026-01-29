/**
 * 最近工程列表组件
 * 使用虚拟滚动显示最近工程
 */

import React, { useState, useEffect, useRef } from 'react'
import { FixedSizeList as List, ListChildComponentProps } from 'react-window'
import type { RecentProject } from '@/types/project'
import ProjectListItem from './ProjectListItem'
import './RecentProjectsList.css'

interface RecentProjectsListProps {
  projects: RecentProject[]
  selectedProjectId?: string
  onProjectClick: (project: RecentProject) => void
  onProjectOpen?: (project: RecentProject) => void
  onShowInExplorer?: (project: RecentProject) => void
  onRemoveProject?: (project: RecentProject) => void
  onCopyPath?: (project: RecentProject) => void
}

/**
 * RecentProjectsList 组件
 */
const RecentProjectsList: React.FC<RecentProjectsListProps> = ({
  projects,
  selectedProjectId,
  onProjectClick,
  onProjectOpen,
  onShowInExplorer,
  onRemoveProject,
  onCopyPath
}) => {
  const listRef = useRef<List>(null)
  const [focusedIndex, setFocusedIndex] = useState<number>(-1)

  // 列表项高度
  const ITEM_HEIGHT = 80

  // 处理键盘导航
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (projects.length === 0) return

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setFocusedIndex(prev => {
          const next = prev + 1
          if (next < projects.length) {
            // 滚动到可见区域
            listRef.current?.scrollToItem(next, 'smart')
            return next
          }
          return prev
        })
        break
      case 'ArrowUp':
        e.preventDefault()
        setFocusedIndex(prev => {
          const next = prev - 1
          if (next >= 0) {
            listRef.current?.scrollToItem(next, 'smart')
            return next
          }
          return prev
        })
        break
      case 'Enter':
      case ' ':
        e.preventDefault()
        if (focusedIndex >= 0 && focusedIndex < projects.length) {
          onProjectClick(projects[focusedIndex])
        }
        break
      case 'Escape':
        setFocusedIndex(-1)
        break
    }
  }

  // 渲染列表项
  const Row = ({ index, style }: ListChildComponentProps) => {
    const project = projects[index]

    return (
      <div style={style}>
        <ProjectListItem
          key={project.projectId}
          project={project}
          isActive={selectedProjectId === project.projectId || focusedIndex === index}
          onClick={onProjectClick}
          onOpen={onProjectOpen}
          onShowInExplorer={onShowInExplorer}
          onRemove={onRemoveProject}
          onCopyPath={onCopyPath}
        />
      </div>
    )
  }

  // 空状态
  if (projects.length === 0) {
    return (
      <div className="recent-projects-list recent-projects-list--empty">
        <div className="recent-projects-list__empty-icon">📁</div>
        <div className="recent-projects-list__empty-text">
          暂无最近工程
        </div>
        <div className="recent-projects-list__empty-hint">
          创建或打开工程后将显示在此处
        </div>
      </div>
    )
  }

  return (
    <div
      className="recent-projects-list"
      onKeyDown={handleKeyDown}
      tabIndex={0}
    >
      <List
        ref={listRef}
        height={Math.min(projects.length * ITEM_HEIGHT, 400)}
        itemCount={projects.length}
        itemSize={ITEM_HEIGHT}
        width="100%"
        overscanCount={5}
      >
        {Row}
      </List>
    </div>
  )
}

export default RecentProjectsList
