/**
 * ComponentGrid Component
 *
 * 组件网格显示，支持拖拽
 *
 * 功能：
 * - 网格布局显示组件
 * - 组件拖拽到画布
 * - 悬停效果
 */

import React from 'react';
import { Component } from '../sidebar/ComponentPanel';
import './ComponentGrid.css';

/**
 * ComponentGrid Props
 */
interface ComponentGridProps {
  components: Component[];
  onDragStart?: (component: Component) => void;
  onDragEnd?: () => void;
  className?: string;
}

/**
 * ComponentGrid Component
 */
export const ComponentGrid: React.FC<ComponentGridProps> = ({
  components,
  onDragStart,
  onDragEnd,
  className = '',
}) => {
  // ==========================================
  // Handlers - 拖拽
  // ==========================================

  /**
   * 拖拽开始
   */
  const handleDragStart = (component: Component, e: React.DragEvent): void => {
    onDragStart?.(component);

    // 设置拖拽数据
    e.dataTransfer.effectAllowed = 'copy';
    e.dataTransfer.setData('application/json', JSON.stringify(component));
  };

  /**
   * 拖拽结束
   */
  const handleDragEnd = (): void => {
    onDragEnd?.();
  };

  // ==========================================
  // Render - Component Item
  // ==========================================

  const renderComponentItem = (component: Component) => (
    <div
      key={component.id}
      className="component-item"
      draggable
      onDragStart={(e) => handleDragStart(component, e)}
      onDragEnd={handleDragEnd}
      title={component.name}
    >
      {/* 图标 */}
      <div className="component-item__icon">{component.icon}</div>

      {/* 标签 */}
      <div className="component-item__label">{component.name}</div>
    </div>
  );

  // ==========================================
  // Main Render
  // ==========================================

  if (!components || components.length === 0) {
    return (
      <div className={`component-grid component-grid--empty ${className}`}>
        <div className="component-grid__empty">
          <div className="component-grid__empty__icon">🔍</div>
          <div className="component-grid__empty__text">未找到组件</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`component-grid ${className}`}>
      <div className="component-grid__items">
        {components.map((component) => renderComponentItem(component))}
      </div>
    </div>
  );
};

/**
 * Default export
 */
export default ComponentGrid;
