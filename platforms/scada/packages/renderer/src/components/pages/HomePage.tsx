import React from 'react'
import { observer } from 'mobx-react-lite'
import Sidebar from '../layout/Sidebar'
import HomePageContent from './HomePageContent'
import { homePageStore } from '@store/homePageStore'
import './HomePage.module.css'

export interface HomePageProps {
  /** 打开工程回调 */
  onOpenProject?: (projectName: string, projectPath: string) => void
}

/**
 * HomePage 首页组件
 * 采用左右分栏布局：左侧导航栏 + 右侧内容区
 *
 * 布局结构：
 * - 左侧：Sidebar（首页导航栏，包含首页、本地、云端、工具）
 * - 右侧：HomePageContent（根据左侧导航显示对应内容）
 *
 * 使用 homePageStore 管理首页独立的导航状态
 */
const HomePage: React.FC<HomePageProps> = observer(({ onOpenProject }) => {
  /**
   * 处理导航项点击
   */
  const handleNavItemClick = (itemId: string): void => {
    homePageStore.setActiveNavItem(itemId as any)
  }

  return (
    <div className="homepage" style={{display: 'flex', width: '100%', height: '100%', flexDirection: 'row'}}>
      <Sidebar
        navigationItems={homePageStore.navigationItems}
        activeItemId={homePageStore.activeNavItem}
        onItemClick={handleNavItemClick}
      />
      <HomePageContent onOpenProject={onOpenProject} />
    </div>
  )
})

export default HomePage
