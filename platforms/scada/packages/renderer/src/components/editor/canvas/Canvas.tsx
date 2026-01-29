/**
 * Canvas Component
 *
 * FE-006-13: 画布区域
 *
 * 功能：
 * - 显示画布内容
 * - 网格背景
 * - 缩放和平移
 * - 拖拽放置组件
 * - 选中元素渲染
 * - 画布事件处理
 */

import React, { useRef, useState, useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { getEditorStore } from '@/store';
import './Canvas.css';

/**
 * Canvas Props
 */
interface CanvasProps {
  className?: string;
  onDropComponent?: (component: any, x: number, y: number) => void;
}

/**
 * Canvas Component
 */
export const Canvas: React.FC<CanvasProps> = observer(({ className = '', onDropComponent }) => {
  const editorStore = getEditorStore();
  const canvasRef = useRef<HTMLDivElement>(null);
  const viewportRef = useRef<HTMLDivElement>(null);

  // 缩放和平移状态
  const [zoom, setZoom] = useState(100);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  // 画布尺寸
  const [canvasSize, setCanvasSize] = useState({ width: 1920, height: 1080 });

  // ==========================================
  // Handlers - 缩放和平移
  // ==========================================

  /**
   * 处理滚轮缩放
   */
  const handleWheel = (e: React.WheelEvent): void => {
    if (e.ctrlKey || e.metaKey) {
      // 缩放
      e.preventDefault();
      const delta = e.deltaY > 0 ? -10 : 10;
      const newZoom = Math.max(25, Math.min(200, zoom + delta));
      setZoom(newZoom);
    } else {
      // 平移
      // TODO: 实现平移逻辑
    }
  };

  /**
   * 处理鼠标按下（开始拖拽画布）
   */
  const handleMouseDown = (e: React.MouseEvent): void => {
    if (e.button === 1 || (e.button === 0 && e.altKey)) {
      // 中键或Alt+左键：平移画布
      setIsDragging(true);
      setDragStart({ x: e.clientX - pan.x, y: e.clientY - pan.y });
    }
  };

  /**
   * 处理鼠标移动（拖拽画布）
   */
  const handleMouseMove = (e: React.MouseEvent): void => {
    if (isDragging) {
      const newPan = {
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      };
      setPan(newPan);
    }
  };

  /**
   * 处理鼠标抬起（结束拖拽）
   */
  const handleMouseUp = (): void => {
    setIsDragging(false);
  };

  // ==========================================
  // Handlers - 拖拽放置
  // ==========================================

  /**
   * 拖拽经过
   */
  const handleDragOver = (e: React.DragEvent): void => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
  };

  /**
   * 拖拽放置
   */
  const handleDrop = (e: React.DragEvent): void => {
    e.preventDefault();

    try {
      const component = JSON.parse(e.dataTransfer.getData('application/json'));
      const rect = canvasRef.current?.getBoundingClientRect();
      if (rect) {
        const x = (e.clientX - rect.left - pan.x) / (zoom / 100);
        const y = (e.clientY - rect.top - pan.y) / (zoom / 100);
        onDropComponent?.(component, x, y);
      }
    } catch (error) {
      console.error('Canvas: drop error', error);
    }
  };

  // ==========================================
  // Effects
  // ==========================================

  /**
   * 监听键盘事件
   */
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent): void => {
      // 空格+拖拽平移画布
      if (e.code === 'Space' && !isDragging) {
        e.preventDefault();
      }
    };

    const handleKeyUp = (e: KeyboardEvent): void => {
      if (e.code === 'Space') {
        setIsDragging(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('keyup', handleKeyUp);
    };
  }, [isDragging]);

  // ==========================================
  // Render
  // ==========================================

  return (
    <div
      ref={canvasRef}
      className={`canvas ${className} ${isDragging ? 'canvas--grabbing' : ''}`}
      onWheel={handleWheel}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
    >
      <div
        ref={viewportRef}
        className="canvas__viewport"
        style={{
          transform: `scale(${zoom / 100}) translate(${pan.x}px, ${pan.y}px)`,
          transformOrigin: '0 0',
          width: canvasSize.width,
          height: canvasSize.height,
        }}
      >
        {/* 画布内容 */}
        <div className="canvas__content">
          {/* TODO: 渲染画布元素 */}
          <div className="canvas-placeholder">
            <div className="canvas-placeholder__text">
              拖拽组件到此处
            </div>
          </div>
        </div>

        {/* 选中框 */}
        {editorStore.hasSelection && (
          <div className="canvas-selection">
            {/* TODO: 渲染选中框 */}
          </div>
        )}
      </div>

      {/* 缩放提示 */}
      {zoom !== 100 && (
        <div className="canvas-zoom-indicator">
          {zoom}%
        </div>
      )}
    </div>
  );
});

/**
 * Default export
 */
export default Canvas;
