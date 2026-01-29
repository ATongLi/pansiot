/**
 * ScreenPanel Component
 *
 * FE-006-07: 画面面板（画面列表）
 *
 * 功能：
 * - 显示所有画面
 * - 画面缩略图预览
 * - 画面新建/删除/重命名
 * - 拖拽排序
 * - 双击打开画面
 */

import React, { useState } from 'react';
import { observer } from 'mobx-react-lite';
import { getEditorStore } from '@/store';
import './ScreenPanel.css';

/**
 * 画面数据
 */
interface Screen {
  id: string;
  name: string;
  thumbnail?: string;
  createdAt: string;
  modifiedAt: string;
}

/**
 * ScreenPanel Props
 */
interface ScreenPanelProps {
  className?: string;
  onScreenSelect?: (screenId: string) => void;
}

/**
 * ScreenPanel Component
 */
export const ScreenPanel: React.FC<ScreenPanelProps> = observer(
  ({ className = '', onScreenSelect }) => {
    const editorStore = getEditorStore();
    const [screens, setScreens] = useState<Screen[]>([
      {
        id: 'screen-1',
        name: '画面 1',
        createdAt: new Date().toISOString(),
        modifiedAt: new Date().toISOString(),
      },
    ]);
    const [draggedScreenId, setDraggedScreenId] = useState<string | null>(null);

    // ==========================================
    // Handlers - 画面操作
    // ==========================================

    /**
     * 添加画面
     */
    const handleAddScreen = (): void => {
      const newScreen: Screen = {
        id: `screen-${Date.now()}`,
        name: `画面 ${screens.length + 1}`,
        createdAt: new Date().toISOString(),
        modifiedAt: new Date().toISOString(),
      };
      setScreens([...screens, newScreen]);
    };

    /**
     * 删除画面
     */
    const handleDeleteScreen = (screenId: string): void => {
      if (confirm('确定要删除这个画面吗？')) {
        setScreens(screens.filter((screen) => screen.id !== screenId));

        // 清除选中状态
        if (editorStore.state.selectedIds.includes(screenId)) {
          editorStore.clearSelection();
        }
      }
    };

    /**
     * 重命名画面
     */
    const handleRenameScreen = (screenId: string, newName: string): void => {
      setScreens(
        screens.map((screen) =>
          screen.id === screenId ? { ...screen, name: newName } : screen
        )
      );
    };

    /**
     * 选中画面
     */
    const handleSelectScreen = (screenId: string): void => {
      editorStore.selectOne(screenId);
      onScreenSelect?.(screenId);
    };

    /**
     * 双击打开画面
     */
    const handleDoubleClickScreen = (screenId: string): void => {
      console.log('ScreenPanel: double click to open', screenId);
      // TODO: 打开画面到画布区域
    };

    // ==========================================
    // Handlers - 拖拽排序
    // ==========================================

    /**
     * 拖拽开始
     */
    const handleDragStart = (screenId: string, e: React.DragEvent): void => {
      setDraggedScreenId(screenId);
      editorStore.startDrag('screen', { screenId });
    };

    /**
     * 拖拽结束
     */
    const handleDragEnd = (): void => {
      setDraggedScreenId(null);
      editorStore.endDrag();
    };

    /**
     * 拖拽经过
     */
    const handleDragOver = (e: React.DragEvent, targetScreenId: string): void => {
      e.preventDefault();
      if (draggedScreenId && draggedScreenId !== targetScreenId) {
        // TODO: 显示拖拽插入指示器
      }
    };

    /**
     * 拖拽放置
     */
    const handleDrop = (e: React.DragEvent, targetScreenId: string): void => {
      e.preventDefault();
      if (!draggedScreenId || draggedScreenId === targetScreenId) {
        return;
      }

      // 重新排序
      const draggedIndex = screens.findIndex((s) => s.id === draggedScreenId);
      const targetIndex = screens.findIndex((s) => s.id === targetScreenId);

      if (draggedIndex !== -1 && targetIndex !== -1) {
        const newScreens = [...screens];
        const [draggedScreen] = newScreens.splice(draggedIndex, 1);
        newScreens.splice(targetIndex, 0, draggedScreen);
        setScreens(newScreens);
      }

      setDraggedScreenId(null);
      editorStore.endDrag();
    };

    // ==========================================
    // Render - Screen Item
    // ==========================================

    const renderScreenItem = (screen: Screen, index: number) => {
      const isSelected = editorStore.state.selectedIds.includes(screen.id);
      const isDragging = draggedScreenId === screen.id;

      return (
        <div
          key={screen.id}
          className={`screen-item ${isSelected ? 'screen-item--selected' : ''} ${
            isDragging ? 'screen-item--dragging' : ''
          }`}
          draggable
          onClick={() => handleSelectScreen(screen.id)}
          onDoubleClick={() => handleDoubleClickScreen(screen.id)}
          onDragStart={(e) => handleDragStart(screen.id, e)}
          onDragEnd={handleDragEnd}
          onDragOver={(e) => handleDragOver(e, screen.id)}
          onDrop={(e) => handleDrop(e, screen.id)}
        >
          {/* 缩略图 */}
          <div className="screen-item__thumbnail">
            {screen.thumbnail ? (
              <img src={screen.thumbnail} alt={screen.name} />
            ) : (
              <div className="screen-item__thumbnail__placeholder">
                <span className="screen-item__thumbnail__icon">📄</span>
              </div>
            )}
          </div>

          {/* 信息 */}
          <div className="screen-item__info">
            <div className="screen-item__name">{screen.name}</div>
            <div className="screen-item__meta">
              {new Date(screen.modifiedAt).toLocaleDateString()}
            </div>
          </div>

          {/* 操作按钮 */}
          <div className="screen-item__actions">
            <button
              className="screen-item__action"
              onClick={(e) => {
                e.stopPropagation();
                const newName = prompt('重命名画面:', screen.name);
                if (newName && newName.trim()) {
                  handleRenameScreen(screen.id, newName.trim());
                }
              }}
              title="重命名"
            >
              ✏️
            </button>
            <button
              className="screen-item__action screen-item__action--danger"
              onClick={(e) => {
                e.stopPropagation();
                handleDeleteScreen(screen.id);
              }}
              title="删除"
            >
              🗑️
            </button>
          </div>
        </div>
      );
    };

    // ==========================================
    // Main Render
    // ==========================================

    return (
      <div className={`screen-panel ${className}`}>
        {/* 工具栏 */}
        <div className="screen-panel__toolbar">
          <button
            className="toolbar-button toolbar-button--icon-only"
            onClick={handleAddScreen}
            title="新建画面"
          >
            +
          </button>
          <div className="screen-panel__toolbar__title">画面</div>
        </div>

        {/* 画面列表 */}
        <div className="screen-panel__list editor-scrollbar">
          {screens.length === 0 ? (
            <div className="editor-empty-state">
              <div className="editor-empty-state__icon">📄</div>
              <div className="editor-empty-state__text">暂无画面</div>
              <div className="editor-empty-state__hint">点击上方 + 添加画面</div>
            </div>
          ) : (
            <div className="screen-panel__items">
              {screens.map((screen, index) => renderScreenItem(screen, index))}
            </div>
          )}
        </div>

        {/* 状态栏 */}
        <div className="screen-panel__status">
          <span className="screen-panel__status__text">
            {screens.length} 个画面
          </span>
        </div>
      </div>
    );
  }
);

/**
 * Default export
 */
export default ScreenPanel;
