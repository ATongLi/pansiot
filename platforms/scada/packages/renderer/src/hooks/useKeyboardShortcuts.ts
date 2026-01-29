/**
 * useKeyboardShortcuts Hook
 *
 * 键盘快捷键管理
 *
 * 功能：
 * - 全局快捷键监听
 * - 工具快捷键
 * - 文件操作快捷键
 * - 编辑操作快捷键
 */

import { useEffect } from 'react';
import { getEditorStore, EditorTool, EditorMode } from '@/store';

/**
 * 快捷键配置
 */
interface ShortcutConfig {
  key: string;
  ctrlKey?: boolean;
  metaKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
  handler: (e: KeyboardEvent) => void;
  description: string;
}

/**
 * useKeyboardShortcuts Hook
 */
export const useKeyboardShortcuts = (shortcuts: ShortcutConfig[]): void => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent): void => {
      for (const shortcut of shortcuts) {
        const keyMatch = e.key.toLowerCase() === shortcut.key.toLowerCase();
        const ctrlMatch = shortcut.ctrlKey === undefined || e.ctrlKey === shortcut.ctrlKey;
        const metaMatch = shortcut.metaKey === undefined || e.metaKey === shortcut.metaKey;
        const shiftMatch = shortcut.shiftKey === undefined || e.shiftKey === shortcut.shiftKey;
        const altMatch = shortcut.altKey === undefined || e.altKey === shortcut.altKey;

        if (keyMatch && ctrlMatch && metaMatch && shiftMatch && altMatch) {
          e.preventDefault();
          shortcut.handler(e);
          break;
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [shortcuts]);
};

/**
 * Editor键盘快捷键Hook
 */
export const useEditorKeyboardShortcuts = (): void => {
  const editorStore = getEditorStore();

  const shortcuts: ShortcutConfig[] = [
    // 文件操作
    {
      key: 'n',
      ctrlKey: true,
      handler: () => console.log('Shortcut: New Project'),
      description: '新建工程',
    },
    {
      key: 'o',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Open Project'),
      description: '打开工程',
    },
    {
      key: 's',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Save Project'),
      description: '保存工程',
    },

    // 撤销/重做
    {
      key: 'z',
      ctrlKey: true,
      handler: () => editorStore.undo(),
      description: '撤销',
    },
    {
      key: 'y',
      ctrlKey: true,
      handler: () => editorStore.redo(),
      description: '重做',
    },
    {
      key: 'z',
      ctrlKey: true,
      shiftKey: true,
      handler: () => editorStore.redo(),
      description: '重做',
    },

    // 删除
    {
      key: 'Delete',
      handler: () => console.log('Shortcut: Delete'),
      description: '删除选中元素',
    },
    {
      key: 'Backspace',
      handler: () => console.log('Shortcut: Delete'),
      description: '删除选中元素',
    },

    // 工具选择
    {
      key: 'v',
      handler: () => editorStore.setCurrentTool(EditorTool.SELECT),
      description: '选择工具',
    },
    {
      key: 'r',
      handler: () => editorStore.setCurrentTool(EditorTool.RECTANGLE),
      description: '矩形工具',
    },
    {
      key: 'c',
      handler: () => editorStore.setCurrentTool(EditorTool.CIRCLE),
      description: '圆形工具',
    },
    {
      key: 'l',
      handler: () => editorStore.setCurrentTool(EditorTool.LINE),
      description: '直线工具',
    },
    {
      key: 't',
      handler: () => editorStore.setCurrentTool(EditorTool.TEXT),
      description: '文本工具',
    },
    {
      key: 'i',
      handler: () => editorStore.setCurrentTool(EditorTool.IMAGE),
      description: '图片工具',
    },

    // 模式切换
    {
      key: 'F1',
      handler: () => editorStore.setMode(EditorMode.EDIT),
      description: '编辑模式',
    },
    {
      key: 'F2',
      handler: () => editorStore.setMode(EditorMode.PREVIEW),
      description: '预览模式',
    },
    {
      key: 'F5',
      handler: () => editorStore.setMode(EditorMode.RUN),
      description: '运行模式',
    },

    // 缩放
    {
      key: '=',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Zoom In'),
      description: '放大',
    },
    {
      key: '-',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Zoom Out'),
      description: '缩小',
    },
    {
      key: '0',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Reset Zoom'),
      description: '重置缩放',
    },

    // 全选
    {
      key: 'a',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Select All'),
      description: '全选',
    },

    // 复制/粘贴
    {
      key: 'c',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Copy'),
      description: '复制',
    },
    {
      key: 'x',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Cut'),
      description: '剪切',
    },
    {
      key: 'v',
      ctrlKey: true,
      handler: () => console.log('Shortcut: Paste'),
      description: '粘贴',
    },

    // 保存
    {
      key: 's',
      ctrlKey: true,
      shiftKey: true,
      handler: () => console.log('Shortcut: Save As'),
      description: '另存为',
    },
  ];

  useKeyboardShortcuts(shortcuts);
};

/**
 * Default export
 */
export default useEditorKeyboardShortcuts;
