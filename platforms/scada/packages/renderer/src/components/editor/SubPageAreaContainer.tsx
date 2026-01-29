import React from 'react'
import { observer } from 'mobx-react-lite'
import { subPageTabStore } from '@/store'
import SubPageAreaView, { type SubPageAreaViewProps } from './SubPageAreaView'

export interface SubPageAreaContainerProps {
  /** 组件拖放回调 */
  onDropComponent?: (component: any, x: number, y: number) => void
  /** 额外的CSS类名 */
  className?: string
}

/**
 * SubPageAreaContainer 子页面内容区容器组件
 * FE-009-06: 使用 Container/Presenter 模式重构
 *
 * 职责：
 * - 连接 subPageTabStore 获取标签页状态
 * - 处理标签页操作（切换、关闭、新建）
 * - 将数据和回调传递给 View 组件
 *
 * 设计模式：
 * - Container 组件：负责状态管理和业务逻辑
 * - View 组件：负责纯 UI 渲染
 * - 分离关注点，提高可测试性和可重用性
 */
const SubPageAreaContainer: React.FC<SubPageAreaContainerProps> = observer(({
  onDropComponent,
  className = '',
}) => {
  /**
   * 处理标签页切换
   */
  const handleTabChange = (tabId: string): void => {
    subPageTabStore.setActiveSubTab(tabId)
  }

  /**
   * 处理标签页关闭
   */
  const handleTabClose = (tabId: string): void => {
    subPageTabStore.closeSubTab(tabId)
  }

  /**
   * 处理新建标签页
   */
  const handleTabAdd = (): void => {
    const tabCount = subPageTabStore.subPageTabs.length
    subPageTabStore.addSubTab({
      title: `画面 ${tabCount + 1}`,
      contentType: 'screen',
      closable: true,
    })
  }

  const viewProps: SubPageAreaViewProps = {
    tabs: subPageTabStore.subPageTabs,
    activeTab: subPageTabStore.activeSubTab,
    onTabChange: handleTabChange,
    onTabClose: handleTabClose,
    onTabAdd: handleTabAdd,
    onDropComponent,
    className,
  }

  return <SubPageAreaView {...viewProps} />
})

export default SubPageAreaContainer
