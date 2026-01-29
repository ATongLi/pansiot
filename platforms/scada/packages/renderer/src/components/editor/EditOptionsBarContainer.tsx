import React from 'react'
import { observer } from 'mobx-react-lite'
import { layoutStore } from '@/store'
import EditOptionsBarView, { type EditOptionsBarViewProps } from './EditOptionsBarView'

export interface EditOptionsBarContainerProps {
  /** 子组件内容（左侧或顶部布局的内容） */
  children: React.ReactNode
  /** 额外的CSS类名 */
  className?: string
}

/**
 * EditOptionsBarContainer 编辑选项栏容器组件
 * FE-009-05: 使用 Container/Presenter 模式重构
 *
 * 职责：
 * - 连接 layoutStore 获取布局状态
 * - 处理布局切换逻辑
 * - 将数据和回调传递给 View 组件
 *
 * 设计模式：
 * - Container 组件：负责状态管理和业务逻辑
 * - View 组件：负责纯 UI 渲染
 * - 分离关注点，提高可测试性和可重用性
 */
const EditOptionsBarContainer: React.FC<EditOptionsBarContainerProps> = observer(({ children, className = '' }) => {
  /**
   * 切换布局方向
   */
  const handleToggleLayout = (): void => {
    const newLayout: EditOptionsBarViewProps['layout'] = layoutStore.editOptionLayout === 'left' ? 'top' : 'left'
    layoutStore.setEditOptionLayout(newLayout)
  }

  const viewProps: EditOptionsBarViewProps = {
    layout: layoutStore.editOptionLayout,
    onToggleLayout: handleToggleLayout,
    children,
    className,
  }

  return <EditOptionsBarView {...viewProps} />
})

export default EditOptionsBarContainer
