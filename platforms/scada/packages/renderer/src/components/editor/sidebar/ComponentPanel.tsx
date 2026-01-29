/**
 * ComponentPanel Component
 *
 * FE-006-08: 组件面板（组件库）
 *
 * 功能：
 * - 显示可用组件库
 * - 组件分类展示
 * - 拖拽组件到画布
 * - 组件搜索过滤
 */

import React, { useState, useMemo } from 'react';
import { observer } from 'mobx-react-lite';
import { getEditorStore } from '@/store';
import { ComponentGrid } from '../componentgrid/ComponentGrid';
import './ComponentPanel.css';

/**
 * 组件分类
 */
interface ComponentCategory {
  id: string;
  name: string;
  components: Component[];
}

/**
 * 组件数据
 */
interface Component {
  id: string;
  name: string;
  type: string;
  category: string;
  icon: string;
  description?: string;
}

/**
 * 默认组件库数据
 */
const DEFAULT_COMPONENTS: ComponentCategory[] = [
  {
    id: 'basic',
    name: '基础组件',
    components: [
      { id: 'rect', name: '矩形', type: 'rectangle', category: 'basic', icon: '⬜' },
      { id: 'circle', name: '圆形', type: 'circle', category: 'basic', icon: '⚪' },
      { id: 'line', name: '直线', type: 'line', category: 'basic', icon: '📏' },
      { id: 'text', name: '文本', type: 'text', category: 'basic', icon: '📝' },
      { id: 'image', name: '图片', type: 'image', category: 'basic', icon: '🖼️' },
    ],
  },
  {
    id: 'industrial',
    name: '工业组件',
    components: [
      { id: 'button', name: '按钮', type: 'button', category: 'industrial', icon: '🔘' },
      { id: 'indicator', name: '指示灯', type: 'indicator', category: 'industrial', icon: '💡' },
      { id: 'gauge', name: '仪表盘', type: 'gauge', category: 'industrial', icon: '🎚️' },
      { id: 'slider', name: '滑动条', type: 'slider', category: 'industrial', icon: '🎚️' },
      { id: 'switch', name: '开关', type: 'switch', category: 'industrial', icon: '🔌' },
    ],
  },
  {
    id: 'chart',
    name: '图表组件',
    components: [
      { id: 'line-chart', name: '折线图', type: 'lineChart', category: 'chart', icon: '📈' },
      { id: 'bar-chart', name: '柱状图', type: 'barChart', category: 'chart', icon: '📊' },
      { id: 'pie-chart', name: '饼图', type: 'pieChart', category: 'chart', icon: '🥧' },
    ],
  },
];

/**
 * ComponentPanel Props
 */
interface ComponentPanelProps {
  className?: string;
  onComponentDragStart?: (component: Component) => void;
}

/**
 * ComponentPanel Component
 */
export const ComponentPanel: React.FC<ComponentPanelProps> = observer(
  ({ className = '', onComponentDragStart }) => {
    const editorStore = getEditorStore();
    const [searchQuery, setSearchQuery] = useState('');
    const [activeCategory, setActiveCategory] = useState<string>('basic');

    // ==========================================
    // Computed - 过滤后的组件
    // ==========================================

    const filteredComponents = useMemo(() => {
      const category = DEFAULT_COMPONENTS.find((c) => c.id === activeCategory);
      if (!category) return [];

      if (!searchQuery) return category.components;

      return category.components.filter((component) =>
        component.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }, [activeCategory, searchQuery]);

    // ==========================================
    // Handlers - 分类选择
    // ==========================================

    const handleCategoryChange = (categoryId: string): void => {
      setActiveCategory(categoryId);
      setSearchQuery(''); // 切换分类时清空搜索
    };

    // ==========================================
    // Handlers - 搜索
    // ==========================================

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
      setSearchQuery(e.target.value);
    };

    // ==========================================
    // Handlers - 拖拽
    // ==========================================

    const handleDragStart = (component: Component): void => {
      editorStore.startDrag('component', component);
      onComponentDragStart?.(component);
    };

    const handleDragEnd = (): void => {
      editorStore.endDrag();
    };

    // ==========================================
    // Render - 分类标签
    // ==========================================

    const renderCategoryTabs = () => (
      <div className="component-panel__tabs">
        {DEFAULT_COMPONENTS.map((category) => (
          <button
            key={category.id}
            className={`component-panel__tab ${
              activeCategory === category.id ? 'component-panel__tab--active' : ''
            }`}
            onClick={() => handleCategoryChange(category.id)}
          >
            {category.name}
          </button>
        ))}
      </div>
    );

    // ==========================================
    // Render - 搜索框
    // ==========================================

    const renderSearchBox = () => (
      <div className="component-panel__search">
        <input
          type="text"
          className="component-panel__search__input"
          placeholder="搜索组件..."
          value={searchQuery}
          onChange={handleSearchChange}
        />
        {searchQuery && (
          <button
            className="component-panel__search__clear"
            onClick={() => setSearchQuery('')}
            title="清除搜索"
          >
            ×
          </button>
        )}
      </div>
    );

    // ==========================================
    // Main Render
    // ==========================================

    return (
      <div className={`component-panel ${className}`}>
        {/* 分类标签 */}
        {renderCategoryTabs()}

        {/* 搜索框 */}
        {renderSearchBox()}

        {/* 组件网格 */}
        <div className="component-panel__grid editor-scrollbar">
          <ComponentGrid
            components={filteredComponents}
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
          />
        </div>

        {/* 状态栏 */}
        <div className="component-panel__status">
          <span className="component-panel__status__text">
            {filteredComponents.length} 个组件
          </span>
        </div>
      </div>
    );
  }
);

/**
 * Default export
 */
export default ComponentPanel;

/**
 * Export types
 */
export type { Component, ComponentCategory };
