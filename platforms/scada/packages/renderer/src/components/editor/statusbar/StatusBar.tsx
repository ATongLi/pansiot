/**
 * StatusBar Component
 *
 * FE-006-14: 底部状态栏
 *
 * 功能：
 * - 显示编辑器状态信息
 * - 画布缩放控制
 * - 鼠标坐标显示
 * - 通知和消息提示
 *
 * Layout:
 * ┌─────────────────────────────────────────────────────────┐
 * │ [状态信息]              [缩放控制] [坐标] [通知]        │
 * └─────────────────────────────────────────────────────────┘
 */

import React, { useState } from 'react';
import { observer } from 'mobx-react-lite';
import { getEditorStore } from '@/store';
import './StatusBar.css';

/**
 * StatusBar Props
 */
interface StatusBarProps {
  className?: string;
}

/**
 * 缩放级别预设
 */
const ZOOM_LEVELS = [25, 50, 75, 100, 125, 150, 200];

/**
 * StatusBar Component
 */
export const StatusBar: React.FC<StatusBarProps> = observer(({ className = '' }) => {
  const editorStore = getEditorStore();
  const [zoom, setZoom] = useState(100);

  // ==========================================
  // Handlers - 缩放控制
  // ==========================================

  /**
   * 放大
   */
  const handleZoomIn = (): void => {
    const currentZoomIndex = ZOOM_LEVELS.indexOf(zoom);
    if (currentZoomIndex < ZOOM_LEVELS.length - 1) {
      setZoom(ZOOM_LEVELS[currentZoomIndex + 1]);
    }
  };

  /**
   * 缩小
   */
  const handleZoomOut = (): void => {
    const currentZoomIndex = ZOOM_LEVELS.indexOf(zoom);
    if (currentZoomIndex > 0) {
      setZoom(ZOOM_LEVELS[currentZoomIndex - 1]);
    }
  };

  /**
   * 重置缩放
   */
  const handleZoomReset = (): void => {
    setZoom(100);
  };

  // ==========================================
  // Render
  // ==========================================

  return (
    <div className={`statusbar ${className}`}>
      {/* 左侧 - 状态信息 */}
      <div className="statusbar__section statusbar__section--left">
        {/* 当前工具 */}
        <div className="statusbar__item">
          <span className="statusbar__item__label">工具:</span>
          <span className="statusbar__item__value">
            {editorStore.currentToolLabel}
          </span>
        </div>

        {/* 当前模式 */}
        <div className="statusbar__item">
          <span className="statusbar__item__label">模式:</span>
          <span className="statusbar__item__value">
            {editorStore.currentModeLabel}
          </span>
        </div>

        {/* 修改状态 */}
        {editorStore.state.mode === 'edit' && (
          <div className="statusbar__item">
            <span className="statusbar__item__value statusbar__item__value--modified">
              未保存
            </span>
          </div>
        )}
      </div>

      {/* 中间 - 缩放控制 */}
      <div className="statusbar__section statusbar__section--center">
        <div className="statusbar__zoom-control">
          {/* 缩小按钮 */}
          <button
            className="statusbar__zoom-button statusbar__zoom-button--icon"
            onClick={handleZoomOut}
            disabled={zoom <= ZOOM_LEVELS[0]}
            title="缩小 (Ctrl+-)"
            aria-label="缩小"
          >
            −
          </button>

          {/* 缩放标签 - 可点击重置 */}
          <span
            className="statusbar__zoom-label statusbar__item--interactive"
            onClick={handleZoomReset}
            title="点击重置缩放 (Ctrl+0)"
          >
            {zoom}%
          </span>

          {/* 放大按钮 */}
          <button
            className="statusbar__zoom-button statusbar__zoom-button--icon"
            onClick={handleZoomIn}
            disabled={zoom >= ZOOM_LEVELS[ZOOM_LEVELS.length - 1]}
            title="放大 (Ctrl++)"
            aria-label="放大"
          >
            +
          </button>
        </div>
      </div>

      {/* 右侧 - 坐标和通知 */}
      <div className="statusbar__section statusbar__section--right">
        {/* 坐标信息 */}
        <div className="statusbar__item">
          <span className="statusbar__item__label">X:</span>
          <span className="statusbar__item__value">0</span>
          <span className="statusbar__item__label statusbar__item__label--spacer">Y:</span>
          <span className="statusbar__item__value">0</span>
        </div>

        {/* 选中元素数量 */}
        {editorStore.hasSelection && (
          <div className="statusbar__item">
            <span className="statusbar__item__label">已选:</span>
            <span className="statusbar__item__value">
              {editorStore.selectionCount}
            </span>
          </div>
        )}

        {/* 撤销/重做状态 */}
        <div className="statusbar__item">
          {editorStore.state.canUndo && (
            <span
              className="statusbar__item__value statusbar__item--interactive"
              title="撤销 (Ctrl+Z)"
            >
              ↶ 撤销
            </span>
          )}
          {editorStore.state.canRedo && (
            <span
              className="statusbar__item__value statusbar__item--interactive"
              title="重做 (Ctrl+Y)"
            >
              ↷ 重做
            </span>
          )}
        </div>

        {/* 通知图标 */}
        {editorStore.state.error && (
          <div className="statusbar__item statusbar__item--error" title={editorStore.state.error}>
            ⚠
          </div>
        )}
      </div>
    </div>
  );
});

/**
 * Default export
 */
export default StatusBar;
