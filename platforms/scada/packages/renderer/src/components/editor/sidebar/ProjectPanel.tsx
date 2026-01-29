/**
 * ProjectPanel Component
 *
 * FE-006-06: 工程面板（工程树形结构）
 *
 * 功能：
 * - 显示工程的树形结构
 * - 工程节点展开/折叠
 * - 节点选中状态
 * - 画面列表管理
 * - 拖拽支持
 */

import React, { useState } from 'react';
import { observer } from 'mobx-react-lite';
import { getEditorStore } from '@/store';
import { TreeView } from '../treeview/TreeView';
import './ProjectPanel.css';

/**
 * 画面节点数据
 */
interface ScreenNode {
  id: string;
  name: string;
  type: 'screen';
  children: SceneNode[];
  expanded: boolean;
}

/**
 * 场景节点（画面中的元素）
 */
interface SceneNode {
  id: string;
  name: string;
  type: 'component' | 'group' | 'layer';
  componentType?: string;
  visible: boolean;
  locked: boolean;
  children: SceneNode[];
  expanded: boolean;
}

/**
 * ProjectPanel Props
 */
interface ProjectPanelProps {
  className?: string;
}

/**
 * ProjectPanel Component
 */
export const ProjectPanel: React.FC<ProjectPanelProps> = observer(({ className = '' }) => {
  const editorStore = getEditorStore();
  const [screens, setScreens] = useState<ScreenNode[]>([
    {
      id: 'screen-1',
      name: '画面 1',
      type: 'screen',
      expanded: true,
      children: [
        {
          id: 'comp-1',
          name: '矩形按钮',
          type: 'component',
          componentType: 'button',
          visible: true,
          locked: false,
          children: [],
          expanded: false,
        },
        {
          id: 'comp-2',
          name: '文本标签',
          type: 'component',
          componentType: 'text',
          visible: true,
          locked: false,
          children: [],
          expanded: false,
        },
      ],
    },
  ]);

  // ==========================================
  // Handlers - 画面操作
  // ==========================================

  /**
   * 添加画面
   */
  const handleAddScreen = (): void => {
    const newScreen: ScreenNode = {
      id: `screen-${Date.now()}`,
      name: `画面 ${screens.length + 1}`,
      type: 'screen',
      expanded: true,
      children: [],
    };
    setScreens([...screens, newScreen]);
  };

  /**
   * 删除画面
   */
  const handleDeleteScreen = (screenId: string): void => {
    setScreens(screens.filter((screen) => screen.id !== screenId));

    // 如果删除的是当前选中的画面，清空选择
    if (editorStore.state.selectedIds.includes(screenId)) {
      editorStore.clearSelection();
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
   * 切换画面展开状态
   */
  const handleToggleScreen = (screenId: string): void => {
    setScreens(
      screens.map((screen) =>
        screen.id === screenId ? { ...screen, expanded: !screen.expanded } : screen
      )
    );
  };

  // ==========================================
  // Handlers - 节点选择
  // ==========================================

  /**
   * 处理节点选中
   */
  const handleSelectNode = (nodeId: string, nodeType: string): void => {
    editorStore.selectOne(nodeId);

    // 如果是组件节点，显示属性面板
    if (nodeType === 'component' || nodeType === 'group') {
      editorStore.setRightSidebarTab('property');
    }
  };

  /**
   * 处理节点多选
   */
  const handleSelectMultiple = (nodeIds: string[]): void => {
    editorStore.setSelectedIds(nodeIds);
  };

  // ==========================================
  // Handlers - 拖拽
  // ==========================================

  /**
   * 处理拖拽开始
   */
  const handleDragStart = (nodeId: string, nodeType: string): void => {
    editorStore.startDrag(nodeType, { nodeId, nodeType });
  };

  /**
   * 处理拖拽结束
   */
  const handleDragEnd = (): void => {
    editorStore.endDrag();
  };

  /**
   * 处理拖拽放置
   */
  const handleDrop = (targetNodeId: string, draggedData: any): void => {
    console.log('ProjectPanel: drop', { targetNodeId, draggedData });
    // TODO: 实现拖拽放置逻辑
  };

  // ==========================================
  // Render
  // ==========================================

  return (
    <div className={`project-panel ${className}`}>
      {/* 工具栏 */}
      <div className="project-panel__toolbar">
        <button
          className="toolbar-button toolbar-button--icon-only"
          onClick={handleAddScreen}
          title="新建画面"
        >
          +
        </button>
        <div className="project-panel__toolbar__title">工程</div>
      </div>

      {/* 工程树 */}
      <div className="project-panel__tree editor-scrollbar">
        {screens.length === 0 ? (
          <div className="editor-empty-state">
            <div className="editor-empty-state__icon">📁</div>
            <div className="editor-empty-state__text">暂无画面</div>
            <div className="editor-empty-state__hint">点击上方 + 添加画面</div>
          </div>
        ) : (
          <TreeView
            data={screens}
            selectedIds={editorStore.state.selectedIds}
            onToggle={handleToggleScreen}
            onSelect={handleSelectNode}
            onSelectMultiple={handleSelectMultiple}
            onRename={handleRenameScreen}
            onDelete={handleDeleteScreen}
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
            onDrop={handleDrop}
          />
        )}
      </div>

      {/* 状态栏 */}
      <div className="project-panel__status">
        <span className="project-panel__status__text">
          {screens.length} 个画面
        </span>
      </div>
    </div>
  );
});

/**
 * Default export
 */
export default ProjectPanel;
