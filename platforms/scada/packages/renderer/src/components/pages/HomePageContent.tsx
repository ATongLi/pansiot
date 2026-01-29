import React from 'react'
import { observer } from 'mobx-react-lite'
import ActionButtons from '../workspace/ActionButtons'
import RecentProjects from '../workspace/RecentProjects'
import { homePageStore } from '@store/homePageStore'
import './HomePageContent.css'

export interface HomePageContentProps {
  /** 打开工程回调 */
  onOpenProject?: (projectName: string, projectPath: string) => void
}

/**
 * HomePageContent 首页右侧内容区
 * 根据 homePageStore.activeNavItem 显示不同的内容
 */
const HomePageContent: React.FC<HomePageContentProps> = observer(({ onOpenProject }) => {
  // 直接在组件体中访问 activeNav，建立响应式依赖
  const activeNav = homePageStore.activeNavItem

  /**
   * 渲染对应导航的内容
   */
  const renderContent = (): React.ReactNode => {
    switch (activeNav) {
      case 'home':
        return (
          <>
            <h2 className="homepage-content__section-title">开始</h2>
            <ActionButtons onOpenProject={onOpenProject} />
            <h2 className="homepage-content__section-title">最近工程</h2>
            <RecentProjects onOpenProject={onOpenProject} />
          </>
        )

      case 'local':
        return (
          <>
            <h2 className="homepage-content__section-title">本地工程</h2>
            <RecentProjects onOpenProject={onOpenProject} />
          </>
        )

      case 'cloud':
        return (
          <div className="homepage-content__placeholder">
            <div className="homepage-content__placeholder-icon">☁️</div>
            <h3>云端工程</h3>
            <p>云端工程功能开发中...</p>
          </div>
        )

      case 'tools':
        return (
          <div className="homepage-content__placeholder">
            <div className="homepage-content__placeholder-icon">🔧</div>
            <h3>工具箱</h3>
            <p>工具箱功能开发中...</p>
          </div>
        )

      default:
        return null
    }
  }

  return (
    <div className="homepage-content" style={{flex: 1, overflow: 'auto', background: 'var(--bg-primary)'}}>
      <div className="homepage-content__inner">{renderContent()}</div>
    </div>
  )
})

export default HomePageContent
