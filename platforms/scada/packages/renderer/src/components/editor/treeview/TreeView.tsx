/**
 * TreeView Component
 *
 * 通用的树形视图组件
 *
 * 功能：
 * - 嵌套节点显示
 * - 展开/折叠节点
 * - 节点选中（单选/多选）
 * - 拖拽支持
 * - 图标和标签显示
 * - 上下文菜单
 */

import React from 'react';
import './TreeView.css';

/**
 * 树节点数据
 */
export interface TreeNode {
  id: string;
  name: string;
  type: string;
  icon?: string;
  expanded?: boolean;
  children?: TreeNode[];
  visible?: boolean;
  locked?: boolean;
  [key: string]: any;
}

/**
 * TreeView Props
 */
interface TreeViewProps {
  data: TreeNode[];
  selectedIds?: string[];
  onToggle?: (nodeId: string) => void;
  onSelect?: (nodeId: string, nodeType: string) => void;
  onSelectMultiple?: (nodeIds: string[]) => void;
  onRename?: (nodeId: string, newName: string) => void;
  onDelete?: (nodeId: string) => void;
  onDragStart?: (nodeId: string, nodeType: string) => void;
  onDragEnd?: () => void;
  onDrop?: (targetNodeId: string, draggedData: any) => void;
  className?: string;
}

/**
 * TreeView Component
 */
export const TreeView: React.FC<TreeViewProps> = ({
  data,
  selectedIds = [],
  onToggle,
  onSelect,
  onSelectMultiple,
  onRename,
  onDelete,
  onDragStart,
  onDragEnd,
  onDrop,
  className = '',
}) => {
  const [draggedNodeId, setDraggedNodeId] = React.useState<string | null>(null);
  const [editingNodeId, setEditingNodeId] = React.useState<string | null>(null);
  const [editingName, setEditingName] = React.useState('');

  // ==========================================
  // Handlers - 节点操作
  // ==========================================

  /**
   * 切换节点展开/折叠
   */
  const handleToggle = (nodeId: string, e: React.MouseEvent): void => {
    e.stopPropagation();
    onToggle?.(nodeId);
  };

  /**
   * 选中节点
   */
  const handleSelect = (
    nodeId: string,
    nodeType: string,
    e: React.MouseEvent
  ): void => {
    e.stopPropagation();

    if (e.ctrlKey || e.metaKey) {
      // 多选
      const newSelectedIds = selectedIds.includes(nodeId)
        ? selectedIds.filter((id) => id !== nodeId)
        : [...selectedIds, nodeId];
      onSelectMultiple?.(newSelectedIds);
    } else {
      // 单选
      onSelect?.(nodeId, nodeType);
    }
  };

  /**
   * 开始重命名
   */
  const handleStartRename = (node: TreeNode, e: React.MouseEvent): void => {
    e.stopPropagation();
    setEditingNodeId(node.id);
    setEditingName(node.name);
  };

  /**
   * 提交重命名
   */
  const handleCommitRename = (): void => {
    if (editingNodeId && editingName.trim()) {
      onRename?.(editingNodeId, editingName.trim());
    }
    setEditingNodeId(null);
    setEditingName('');
  };

  /**
   * 取消重命名
   */
  const handleCancelRename = (): void => {
    setEditingNodeId(null);
    setEditingName('');
  };

  /**
   * 删除节点
   */
  const handleDelete = (nodeId: string, e: React.MouseEvent): void => {
    e.stopPropagation();
    if (confirm('确定要删除此节点吗？')) {
      onDelete?.(nodeId);
    }
  };

  // ==========================================
  // Handlers - 拖拽
  // ==========================================

  /**
   * 拖拽开始
   */
  const handleDragStart = (node: TreeNode, e: React.DragEvent): void => {
    setDraggedNodeId(node.id);
    onDragStart?.(node.id, node.type);

    // 设置拖拽数据
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('application/json', JSON.stringify(node));
  };

  /**
   * 拖拽结束
   */
  const handleDragEnd = (): void => {
    setDraggedNodeId(null);
    onDragEnd?.();
  };

  /**
   * 拖拽经过
   */
  const handleDragOver = (nodeId: string, e: React.DragEvent): void => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  /**
   * 拖拽放置
   */
  const handleDrop = (targetNodeId: string, e: React.DragEvent): void => {
    e.preventDefault();
    e.stopPropagation();

    if (!draggedNodeId || draggedNodeId === targetNodeId) {
      return;
    }

    try {
      const draggedData = JSON.parse(e.dataTransfer.getData('application/json'));
      onDrop?.(targetNodeId, draggedData);
    } catch (error) {
      console.error('TreeView: drop error', error);
    }

    setDraggedNodeId(null);
  };

  // ==========================================
  // Render - Tree Node
  // ==========================================

  const renderNode = (node: TreeNode, level: number = 0): React.ReactNode => {
    const hasChildren = node.children && node.children.length > 0;
    const isSelected = selectedIds.includes(node.id);
    const isDragging = draggedNodeId === node.id;
    const isEditing = editingNodeId === node.id;

    return (
      <div key={node.id} className="tree-node__container">
        {/* 节点行 */}
        <div
          className={`tree-node ${isSelected ? 'tree-node--selected' : ''} ${
            isDragging ? 'tree-node--dragging' : ''
          }`}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          draggable
          onClick={(e) => handleSelect(node.id, node.type, e)}
          onDragStart={(e) => handleDragStart(node, e)}
          onDragEnd={handleDragEnd}
          onDragOver={(e) => handleDragOver(node.id, e)}
          onDrop={(e) => handleDrop(node.id, e)}
        >
          {/* 展开/折叠按钮 */}
          {hasChildren ? (
            <span
              className={`tree-node__toggle ${node.expanded ? 'tree-node__toggle--expanded' : ''}`}
              onClick={(e) => handleToggle(node.id, e)}
            >
              ▶
            </span>
          ) : (
            <span className="tree-node__toggle tree-node__toggle--empty" />
          )}

          {/* 节点图标 */}
          <span className="tree-node__icon">{node.icon || '📄'}</span>

          {/* 节点标签 */}
          {isEditing ? (
            <input
              type="text"
              className="tree-node__input"
              value={editingName}
              onChange={(e) => setEditingName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleCommitRename();
                } else if (e.key === 'Escape') {
                  handleCancelRename();
                }
              }}
              onBlur={handleCommitRename}
              autoFocus
              onClick={(e) => e.stopPropagation()}
            />
          ) : (
            <span className="tree-node__label">{node.name}</span>
          )}

          {/* 可见性/锁定图标 */}
          {node.visible !== undefined && (
            <span
              className={`tree-node__visibility ${
                !node.visible ? 'tree-node__visibility--hidden' : ''
              }`}
              title={node.visible ? '可见' : '隐藏'}
            >
              👁
            </span>
          )}

          {node.locked && (
            <span className="tree-node__lock" title="已锁定">
              🔒
            </span>
          )}

          {/* 操作按钮 */}
          <div className="tree-node__actions">
            <button
              className="tree-node__action"
              onClick={(e) => handleStartRename(node, e)}
              title="重命名"
            >
              ✏️
            </button>
            <button
              className="tree-node__action tree-node__action--danger"
              onClick={(e) => handleDelete(node.id, e)}
              title="删除"
            >
              🗑️
            </button>
          </div>
        </div>

        {/* 子节点 */}
        {hasChildren && node.expanded && (
          <div className="tree-node__children">
            {node.children!.map((child) => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  // ==========================================
  // Main Render
  // ==========================================

  if (!data || data.length === 0) {
    return (
      <div className={`tree-view tree-view--empty ${className}`}>
        <div className="tree-view__empty">
          <div className="tree-view__empty__icon">📁</div>
          <div className="tree-view__empty__text">暂无数据</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`tree-view ${className}`}>
      <div className="tree-view__nodes">{data.map((node) => renderNode(node))}</div>
    </div>
  );
};

/**
 * Default export
 */
export default TreeView;
