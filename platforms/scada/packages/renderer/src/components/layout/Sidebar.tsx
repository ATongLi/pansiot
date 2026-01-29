import React from 'react'
import { observer } from 'mobx-react-lite'
import NavItem from '../navigation/NavItem'
import './Sidebar.css'

/**
 * Sidebar 导航栏组件
 * 纯展示组件，接受导航项配置和激活状态
 *
 * @param navigationItems - 导航项列表
 * @param activeItemId - 当前激活的导航项ID
 * @param onItemClick - 导航项点击回调
 */
interface SidebarProps {
  navigationItems: Array<{
    id: string
    label: string
    icon: string
  }>
  activeItemId: string
  onItemClick: (itemId: string) => void
}

const Sidebar: React.FC<SidebarProps> = observer(({ navigationItems, activeItemId, onItemClick }) => {
  return (
    <div
      className="sidebar"
      style={{width: '60px', flexShrink: 0, height: '100%'}}
    >
      {navigationItems.map((item) => (
        <NavItem
          key={item.id}
          item={item}
          isActive={activeItemId === item.id}
          onClick={() => onItemClick(item.id)}
        />
      ))}
    </div>
  )
})

export default Sidebar
