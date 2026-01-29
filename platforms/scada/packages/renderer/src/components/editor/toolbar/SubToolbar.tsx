/**
 * SubToolbar 子工具栏组件
 *
 * 功能：
 * - 根据主分类显示对应的子工具栏内容
 * - 支持展开/收缩
 * - 使用线条图标
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { LineIcon } from '@/components/common/LineIcon'
import { ToolbarCategory, SubToolbarItem } from './ToolbarTypes'
import { getEditorStore } from '@/store'
import './SubToolbar.css'

export interface SubToolbarProps {
  className?: string
  visible?: boolean
  activeCategory?: ToolbarCategory
  /** 工程文件操作回调 */
  onNewProject?: () => void
  onOpenProject?: () => void
  onSaveProject?: () => void
}

/**
 * 子工具栏按钮配置
 */
const getSubToolbarButtons = (category: ToolbarCategory, callbacks: {
  onNewProject?: () => void
  onOpenProject?: () => void
  onSaveProject?: () => void
}): SubToolbarItem[] => {
  const editorStore = getEditorStore()

  switch (category) {
    case ToolbarCategory.PROJECT:
      return [
        { id: 'new', icon: 'new-project', label: '新建工程', shortcut: 'Ctrl+N', action: callbacks.onNewProject || (() => {}) },
        { id: 'open', icon: 'open-project', label: '打开工程', shortcut: 'Ctrl+O', action: callbacks.onOpenProject || (() => {}) },
        { id: 'save', icon: 'save-project', label: '保存工程', shortcut: 'Ctrl+S', action: callbacks.onSaveProject || (() => {}) },
        { id: 'saveAs', icon: 'save-as', label: '另存为', shortcut: 'Ctrl+Shift+S', action: () => console.log('另存为') },
        { id: 'autoSave', icon: 'auto-save', label: '自动保存', action: () => console.log('自动保存') },
        { id: 'settings', icon: 'project-settings', label: '工程设置', action: () => console.log('工程设置') },
      ]

    case ToolbarCategory.GENERAL:
      return [
        // 剪贴板组
        { id: 'copy', icon: 'copy', label: '复制', shortcut: 'Ctrl+C', action: () => editorStore.copy() },
        { id: 'paste', icon: 'paste', label: '粘贴', shortcut: 'Ctrl+V', action: () => editorStore.paste() },
        { id: 'delete', icon: 'delete', label: '删除', shortcut: 'Delete', action: () => editorStore.delete() },
        { id: 'formatBrush', icon: 'format-brush', label: '格式刷', shortcut: 'Ctrl+Shift+F', action: () => console.log('格式刷') },
        { type: 'separator' },
        // 编辑组
        { id: 'undo', icon: 'undo', label: '撤销', shortcut: 'Ctrl+Z', action: () => editorStore.undo(), disabled: !editorStore.state.canUndo },
        { id: 'redo', icon: 'redo', label: '重做', shortcut: 'Ctrl+Y', action: () => editorStore.redo(), disabled: !editorStore.state.canRedo },
        { id: 'select', icon: 'select', label: '选择', shortcut: 'V', action: () => console.log('选择') },
        { id: 'multiDuplicate', icon: 'multi-duplicate', label: '多重复制', shortcut: 'Ctrl+D', action: () => console.log('多重复制') },
        { type: 'separator' },
        // 对齐组
        { id: 'alignLeft', icon: 'align-left', label: '左对齐', action: () => console.log('左对齐') },
        { id: 'alignRight', icon: 'align-right', label: '右对齐', action: () => console.log('右对齐') },
        { id: 'alignCenter', icon: 'align-center', label: '居中对齐', action: () => console.log('居中对齐') },
        { type: 'separator' },
        // 组合组
        { id: 'group', icon: 'group', label: '组合', shortcut: 'Ctrl+G', action: () => console.log('组合') },
        { id: 'ungroup', icon: 'ungroup', label: '拆分', shortcut: 'Ctrl+Shift+G', action: () => console.log('拆分') },
      ]

    case ToolbarCategory.DEBUG:
      return [
        { id: 'compile', icon: 'compile', label: '编译', shortcut: 'F7', action: () => console.log('编译') },
        { id: 'onlineCompile', icon: 'online-compile', label: '在线编译', action: () => console.log('在线编译') },
        { id: 'offlineCompile', icon: 'offline-compile', label: '离线编译', action: () => console.log('离线编译') },
        { id: 'download', icon: 'download', label: '下载', action: () => console.log('下载') },
        { id: 'upload', icon: 'upload', label: '上传', action: () => console.log('上传') },
      ]

    case ToolbarCategory.TOOLS:
      return [
        { id: 'textLibrary', icon: 'text-library', label: '文本库', action: () => console.log('文本库') },
        { id: 'mediaLibrary', icon: 'media-library', label: '媒体库', action: () => console.log('媒体库') },
        { id: 'deviceManager', icon: 'device-manager', label: '设备管理', action: () => console.log('设备管理') },
        { id: 'passthrough', icon: 'passthrough', label: '透传', action: () => console.log('透传') },
        { id: 'debugLog', icon: 'debug-log', label: '调试日志', action: () => console.log('调试日志') },
      ]

    case ToolbarCategory.VIEW:
      return [
        // 显示组
        { id: 'layer', icon: 'layer', label: '图层', action: () => console.log('图层') },
        { id: 'history', icon: 'history', label: '历史记录', action: () => console.log('历史记录') },
        { type: 'separator' },
        // 位置组
        { id: 'grid', icon: 'grid', label: '网格', action: () => console.log('网格') },
        { id: 'snapGrid', icon: 'snap-grid', label: '对齐网格', action: () => console.log('对齐网格') },
      ]

    case ToolbarCategory.HELP:
      return [
        { id: 'version', icon: 'version', label: '版本信息', action: () => console.log('版本信息') },
        { id: 'changelog', icon: 'changelog', label: '更新日志', action: () => console.log('更新日志') },
        { id: 'update', icon: 'update', label: '检查更新', action: () => console.log('检查更新') },
        { id: 'docs', icon: 'help-docs', label: '帮助文档', shortcut: 'F1', action: () => console.log('帮助文档') },
        { id: 'contact', icon: 'contact', label: '联系我们', action: () => console.log('联系我们') },
      ]

    default:
      return []
  }
}

/**
 * SubToolbar 组件
 */
export const SubToolbar: React.FC<SubToolbarProps> = observer(({
  className = '',
  visible = true,
  activeCategory = ToolbarCategory.PROJECT,
  onNewProject,
  onOpenProject,
  onSaveProject,
}) => {
  const editorStore = getEditorStore()
  const buttons = getSubToolbarButtons(activeCategory, {
    onNewProject,
    onOpenProject,
    onSaveProject,
  })

  // ==========================================
  // Render
  // ==========================================

  if (!visible) {
    return null
  }

  return (
    <div className={`sub-toolbar ${className}`}>
      {buttons.map((button, index) => {
        // Handle separator
        if ('type' in button && button.type === 'separator') {
          return <div key={`separator-${index}`} className="toolbar-separator" />
        }

        // Handle regular button
        return (
          <button
            key={button.id}
            className={`toolbar-button ${button.disabled ? 'toolbar-button--disabled' : ''}`}
            onClick={button.action}
            disabled={button.disabled}
            title={`${button.label}${button.shortcut ? ` (${button.shortcut})` : ''}`}
            data-shortcut={button.shortcut || undefined}
          >
            <LineIcon name={button.icon} size={15} />
            <span className="toolbar-button__label">{button.label}</span>
          </button>
        )
      })}
    </div>
  )
})

/**
 * 默认导出
 */
export default SubToolbar
