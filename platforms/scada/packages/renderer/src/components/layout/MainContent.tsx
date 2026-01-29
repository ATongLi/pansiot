import React from 'react'
import { HomePage } from '../pages'
import './MainContent.css'

/**
 * MainContent 主内容区组件
 * IMP-009: 重构为使用新的 HomePage 组件
 *
 * 保留原组件以兼容 App.tsx 中的引用
 * 现在 HomePage 组件内部包含左右分栏布局和所有首页逻辑
 */
interface MainContentProps {
  onOpenProject?: (projectName: string, projectPath: string) => void
}

const MainContent: React.FC<MainContentProps> = ({ onOpenProject }) => {
  return (
    <div className="main-content">
      <HomePage onOpenProject={onOpenProject} />
    </div>
  )
}

export default MainContent
