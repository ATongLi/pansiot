/**
 * SubPageTabs Component
 *
 * FE-006-12: 子页面标签栏
 * IMP-009: 使用 subPageTabStore 管理状态
 *
 * 功能：
 * - 显示所有打开的画面标签
 * - 标签切换
 * - 标签关闭
 * - 新建标签
 * - 标签拖拽排序
 */

import React, { useState } from 'react';
import { observer } from 'mobx-react-lite';
import type { SubPageTab } from '@/store';
import './SubPageTabs.css';

/**
 * SubPageTabs Props
 * IMP-009: tabs 和 activeTab 从父组件传入（来自 subPageTabStore）
 */
interface SubPageTabsProps {
  className?: string;
  tabs: SubPageTab[];
  activeTab: string;
  onTabChange?: (tabId: string) => void;
  onTabClose?: (tabId: string) => void;
  onTabAdd?: () => void;
}

/**
 * SubPageTabs Component
 * IMP-009: 纯展示组件，状态由 subPageTabStore 管理
 */
export const SubPageTabs: React.FC<SubPageTabsProps> = observer(
  ({
    className = '',
    tabs,
    activeTab,
    onTabChange,
    onTabClose,
    onTabAdd,
  }) => {
    const [draggedTabId, setDraggedTabId] = useState<string | null>(null);

    // ==========================================
    // Handlers - 标签操作
    // ==========================================

    /**
     * 切换标签
     */
    const handleTabChange = (tabId: string): void => {
      onTabChange?.(tabId);
      console.log('SubPageTabs: tab changed to', tabId);
    };

    /**
     * 关闭标签
     */
    const handleTabClose = (tabId: string, e: React.MouseEvent): void => {
      e.stopPropagation();

      const tab = tabs.find((t) => t.id === tabId);
      if (tab && !tab.closable) {
        return; // 不允许关闭
      }

      // IMP-009: 状态管理由 subPageTabStore 处理，这里只调用回调
      onTabClose?.(tabId);
    };

    /**
     * 添加新标签
     * IMP-009: 只调用回调，状态管理由 subPageTabStore 处理
     */
    const handleAddTab = (): void => {
      onTabAdd?.();
    };

    // ==========================================
    // Handlers - 拖拽排序
    // ==========================================

    /**
     * 拖拽开始
     */
    const handleDragStart = (tabId: string, e: React.DragEvent): void => {
      setDraggedTabId(tabId);
      e.dataTransfer.effectAllowed = 'move';
    };

    /**
     * 拖拽结束
     */
    const handleDragEnd = (): void => {
      setDraggedTabId(null);
    };

    /**
     * 拖拽经过
     */
    const handleDragOver = (e: React.DragEvent): void => {
      e.preventDefault();
      e.dataTransfer.dropEffect = 'move';
    };

    /**
     * 拖拽放置
     * IMP-009: TODO - 需要在 subPageTabStore 中实现重新排序方法
     */
    const handleDrop = (targetTabId: string, e: React.DragEvent): void => {
      e.preventDefault();
      e.stopPropagation();

      if (!draggedTabId || draggedTabId === targetTabId) {
        return;
      }

      // TODO: IMP-009 - 调用 subPageTabStore.reorderTabs(draggedTabId, targetTabId)
      console.log('SubPageTabs: drag and drop reorder', draggedTabId, '→', targetTabId);

      setDraggedTabId(null);
    };

    // ==========================================
    // Render - Tab
    // ==========================================

    const renderTab = (tab: SubPageTab) => {
      const isActive = activeTab === tab.id;
      const isDragging = draggedTabId === tab.id;

      return (
        <div
          key={tab.id}
          className={`sub-page-tab ${isActive ? 'sub-page-tab--active' : ''} ${
            isDragging ? 'sub-page-tab--dragging' : ''
          }`}
          draggable={tab.closable}
          onClick={() => handleTabChange(tab.id)}
          onDragStart={(e) => handleDragStart(tab.id, e)}
          onDragEnd={handleDragEnd}
          onDragOver={handleDragOver}
          onDrop={(e) => handleDrop(tab.id, e)}
        >
          {/* 标签文本 */}
          <span className="sub-page-tab__label">
            {tab.title}
            {/* TODO: IMP-009 - 添加 modified 状态支持 */}
          </span>

          {/* 关闭按钮 */}
          {tab.closable && (
            <span
              className="sub-page-tab__close"
              onClick={(e) => handleTabClose(tab.id, e)}
              title="关闭标签"
            >
              ×
            </span>
          )}
        </div>
      );
    };

    // ==========================================
    // Main Render
    // ==========================================

    return (
      <div className={`sub-page-tabs ${className}`}>
        <div className="sub-page-tabs__list">
          {tabs.map((tab) => renderTab(tab))}
        </div>

        {/* 添加标签按钮 */}
        <button
          className="sub-page-tabs__add"
          onClick={handleAddTab}
          title="新建画面"
        >
          +
        </button>
      </div>
    );
  }
);

/**
 * Default export
 */
export default SubPageTabs;
