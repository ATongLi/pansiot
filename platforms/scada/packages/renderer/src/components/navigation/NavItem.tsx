import React from 'react'
import type { NavItem as NavItemType } from '@/types/navigation'
import './NavItem.css'

interface NavItemProps {
  item: NavItemType
  isActive: boolean
  onClick: () => void
}

const NavItem: React.FC<NavItemProps> = ({ item, isActive, onClick }) => {
  return (
    <div
      className={`nav-item ${isActive ? 'nav-item--active' : ''}`}
      onClick={onClick}
      role="button"
      tabIndex={0}
      title={item.label}  /* 添加 tooltip 提示 */
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          onClick()
        }
      }}
    >
      {/* 纯图标设计 */}
      <div
        className="nav-item__icon"
        dangerouslySetInnerHTML={{ __html: item.icon || '' }}
      />
      {/* label 移到 tooltip 中，不直接显示 */}
    </div>
  )
}

export default NavItem
