import React from 'react'
import { observer } from 'mobx-react-lite'
import TitleBar from './components/layout/TitleBar'
import Sidebar from './components/layout/Sidebar'
import MainContent from './components/layout/MainContent'
import { EditorLayout } from './components/editor/EditorLayout'
import { tabStore, uiStore } from './store'
import type { Tab } from './store'
import './index.css'

/**
 * App 应用根组件
 * 负责管理应用级标签栏和内容渲染
 */
const App: React.FC = observer(() => {
  /**
   * 打开工程并创建编辑器标签页
   * @param projectName 工程名称
   * @param projectPath 工程文件路径
   */
  const openProjectInEditor = (projectName: string, projectPath: string): void => {
    // 使用 TabStore 添加新的工程标签
    tabStore.addTab({
      title: projectName,
      type: 'editor',
      closable: true,
      projectPath,
    })

    console.log('App: 打开工程', projectName, projectPath)
  }

  /**
   * 标签页切换处理
   * @param tabId 要切换到的标签ID
   */
  const handleTabChange = (tabId: string): void => {
    tabStore.setActiveTab(tabId)
  }

  /**
   * 标签页关闭处理
   * @param tabId 要关闭的标签ID
   */
  const handleTabClose = (tabId: string): void => {
    // TabStore 会自动处理首页不可关闭和标签切换逻辑
    tabStore.closeTab(tabId)
  }

  /**
   * 处理应用级导航项点击
   * @param itemId 导航项ID
   */
  const handleAppNavItemClick = (itemId: string): void => {
    uiStore.setActiveNavItem(itemId)
  }

  /**
   * 渲染当前活动标签的内容
   * 根据标签类型渲染不同的组件：
   * - 'home': 首页内容（MainContent，包含自己的左侧导航）
   * - 'editor': 工程编辑器（EditorLayout）
   */
  const renderTabContent = (): React.ReactNode => {
    const activeTabObj = tabStore.activeTabObj

    // 没有激活标签时，显示首页
    if (!activeTabObj) {
      return <MainContent onOpenProject={openProjectInEditor} />
    }

    // 工程编辑器标签
    if (activeTabObj.type === 'editor') {
      return <EditorLayout />
    }

    // 首页标签（默认）
    return <MainContent onOpenProject={openProjectInEditor} />
  }

  /**
   * 判断是否显示应用级 Sidebar
   * 首页和编辑器都有自己的内部侧边栏，不需要应用级 Sidebar
   */
  const shouldShowSidebar = (): boolean => {
    // 始终返回 false，因为：
    // - 首页有 HomePage 内部的 Sidebar
    // - 编辑器有 EditorLayout 内部的左右侧边栏
    return false
  }

  return (
    <div className="app">
      {/* 自定义标题栏（包含应用级标签栏） */}
      <TitleBar
        activeTab={tabStore.activeTab}
        tabs={tabStore.tabs as Tab[]}
        onTabChange={handleTabChange}
        onTabClose={handleTabClose}
      />
      <div className="app-body">
        {/* 应用级 Sidebar：只在首页显示，编辑器模式隐藏以提供更大编辑空间 */}
        {shouldShowSidebar() && (
          <Sidebar
            navigationItems={uiStore.navigationItems}
            activeItemId={uiStore.activeNavItem}
            onItemClick={handleAppNavItemClick}
          />
        )}
        {renderTabContent()}
      </div>
    </div>
  )
})

export default App
