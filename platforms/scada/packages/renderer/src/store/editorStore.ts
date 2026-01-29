/**
 * Editor Store - MobX State Management
 *
 * FE-006-15: MobX状态管理架构
 * 负责编辑器UI状态管理，包括：
 * - 工具栏状态（当前工具、按钮状态）
 * - 侧边栏状态（当前Tab、面板显示/隐藏）
 * - 编辑器模式（编辑模式、预览模式、运行模式）
 * - UI交互状态（拖拽、缩放、选中）
 */

import { makeAutoObservable, runInAction } from 'mobx';

/**
 * 编辑器工具类型
 */
export enum EditorTool {
  SELECT = 'select',       // 选择工具
  RECTANGLE = 'rectangle', // 矩形工具
  CIRCLE = 'circle',       // 圆形工具
  LINE = 'line',           // 直线工具
  TEXT = 'text',           // 文本工具
  IMAGE = 'image',         // 图片工具
}

/**
 * 编辑器模式
 */
export enum EditorMode {
  EDIT = 'edit',     // 编辑模式
  PREVIEW = 'preview', // 预览模式
  RUN = 'run',       // 运行模式
}

/**
 * 左侧边栏Tab类型
 */
export enum LeftSidebarTab {
  PROJECT = 'project',   // 工程面板
  SCREEN = 'screen',     // 画面面板
  COMPONENT = 'component', // 组件面板
}

/**
 * 右侧边栏Tab类型
 */
export enum RightSidebarTab {
  PROPERTY = 'property', // 属性面板
  LAYER = 'layer',       // 图层面板
}

/**
 * Editor Store State
 */
interface EditorState {
  // 当前工具
  currentTool: EditorTool;

  // 编辑器模式
  mode: EditorMode;

  // 左侧边栏当前Tab
  leftSidebarActiveTab: LeftSidebarTab;

  // 右侧边栏当前Tab
  rightSidebarActiveTab: RightSidebarTab | null;

  // 右侧边栏是否显示
  rightSidebarVisible: boolean;

  // 是否正在拖拽
  isDragging: boolean;

  // 拖拽数据
  dragData: {
    type: string;
    data: any;
  } | null;

  // 选中元素的ID列表
  selectedIds: string[];

  // 剪贴板数据
  clipboard: {
    type: string;
    data: any;
  } | null;

  // 撤销/重做状态
  canUndo: boolean;
  canRedo: boolean;

  // 加载状态
  isLoading: boolean;

  // 错误信息
  error: string | null;
}

/**
 * Editor Store - 编辑器状态管理
 */
export class EditorStore {
  // State
  state: EditorState = {
    currentTool: EditorTool.SELECT,
    mode: EditorMode.EDIT,
    leftSidebarActiveTab: LeftSidebarTab.PROJECT,
    rightSidebarActiveTab: null,
    rightSidebarVisible: false,
    isDragging: false,
    dragData: null,
    selectedIds: [],
    clipboard: null,
    canUndo: false,
    canRedo: false,
    isLoading: false,
    error: null,
  };

  constructor() {
    makeAutoObservable(this);
  }

  // ==========================================
  // Actions - 工具操作
  // ==========================================

  /**
   * 设置当前工具
   */
  setCurrentTool = (tool: EditorTool): void => {
    this.state.currentTool = tool;
  };

  /**
   * 设置编辑器模式
   */
  setMode = (mode: EditorMode): void => {
    this.state.mode = mode;

    // 预览/运行模式下隐藏右侧边栏
    if (mode !== EditorMode.EDIT) {
      this.state.rightSidebarVisible = false;
    }
  };

  // ==========================================
  // Actions - 侧边栏操作
  // ==========================================

  /**
   * 切换左侧边栏Tab
   */
  setLeftSidebarTab = (tab: LeftSidebarTab): void => {
    this.state.leftSidebarActiveTab = tab;
  };

  /**
   * 切换右侧边栏Tab
   */
  setRightSidebarTab = (tab: RightSidebarTab | null): void => {
    this.state.rightSidebarActiveTab = tab;

    // 如果有tab被选中，显示右侧边栏
    if (tab !== null) {
      this.state.rightSidebarVisible = true;
    }
  };

  /**
   * 切换右侧边栏显示/隐藏
   */
  toggleRightSidebar = (): void => {
    this.state.rightSidebarVisible = !this.state.rightSidebarVisible;
  };

  /**
   * 隐藏右侧边栏
   */
  hideRightSidebar = (): void => {
    this.state.rightSidebarVisible = false;
  };

  // ==========================================
  // Actions - 拖拽操作
  // ==========================================

  /**
   * 开始拖拽
   */
  startDrag = (type: string, data: any): void => {
    runInAction(() => {
      this.state.isDragging = true;
      this.state.dragData = { type, data };
    });
  };

  /**
   * 结束拖拽
   */
  endDrag = (): void => {
    runInAction(() => {
      this.state.isDragging = false;
      this.state.dragData = null;
    });
  };

  // ==========================================
  // Actions - 选择操作
  // ==========================================

  /**
   * 设置选中元素
   */
  setSelectedIds = (ids: string[]): void => {
    this.state.selectedIds = ids;

    // 如果有选中元素，显示属性面板
    if (ids.length > 0 && this.state.mode === EditorMode.EDIT) {
      this.state.rightSidebarActiveTab = RightSidebarTab.PROPERTY;
      this.state.rightSidebarVisible = true;
    }
  };

  /**
   * 清空选中
   */
  clearSelection = (): void => {
    this.state.selectedIds = [];
  };

  /**
   * 选中单个元素
   */
  selectOne = (id: string): void => {
    this.state.selectedIds = [id];
  };

  /**
   * 追加选中元素
   */
  addToSelection = (id: string): void => {
    if (!this.state.selectedIds.includes(id)) {
      this.state.selectedIds = [...this.state.selectedIds, id];
    }
  };

  /**
   * 从选中中移除
   */
  removeFromSelection = (id: string): void => {
    this.state.selectedIds = this.state.selectedIds.filter(
      (selectedId) => selectedId !== id
    );
  };

  /**
   * 切换选中状态
   */
  toggleSelection = (id: string): void => {
    if (this.state.selectedIds.includes(id)) {
      this.removeFromSelection(id);
    } else {
      this.addToSelection(id);
    }
  };

  /**
   * 全选
   */
  selectAll = (ids: string[]): void => {
    this.state.selectedIds = ids;
  };

  // ==========================================
  // Actions - 剪贴板操作
  // ==========================================

  /**
   * 复制
   */
  copy = (type: string, data: any): void => {
    this.state.clipboard = { type, data };
  };

  /**
   * 剪切
   */
  cut = (type: string, data: any): void => {
    this.state.clipboard = { type, data };
  };

  /**
   * 粘贴
   */
  paste = (): { type: string; data: any } | null => {
    return this.state.clipboard;
  };

  /**
   * 清空剪贴板
   */
  clearClipboard = (): void => {
    this.state.clipboard = null;
  };

  // ==========================================
  // Actions - 撤销/重做
  // ==========================================

  /**
   * 设置撤销状态
   */
  setUndoRedoState = (canUndo: boolean, canRedo: boolean): void => {
    this.state.canUndo = canUndo;
    this.state.canRedo = canRedo;
  };

  /**
   * 撤销
   */
  undo = (): void => {
    // 由CanvasStore处理实际撤销逻辑
    console.log('EditorStore: undo triggered');
  };

  /**
   * 重做
   */
  redo = (): void => {
    // 由CanvasStore处理实际重做逻辑
    console.log('EditorStore: redo triggered');
  };

  // ==========================================
  // Actions - 加载和错误状态
  // ==========================================

  /**
   * 设置加载状态
   */
  setLoading = (loading: boolean): void => {
    this.state.isLoading = loading;
  };

  /**
   * 设置错误信息
   */
  setError = (error: string | null): void => {
    this.state.error = error;
  };

  /**
   * 清空错误
   */
  clearError = (): void => {
    this.state.error = null;
  };

  // ==========================================
  // Computed - 派生状态
  // ==========================================

  /**
   * 是否有选中元素
   */
  get hasSelection(): boolean {
    return this.state.selectedIds.length > 0;
  }

  /**
   * 选中元素数量
   */
  get selectionCount(): number {
    return this.state.selectedIds.length;
  }

  /**
   * 是否可以复制
   */
  get canCopy(): boolean {
    return this.state.mode === EditorMode.EDIT && this.hasSelection;
  }

  /**
   * 是否可以粘贴
   */
  get canPaste(): boolean {
    return this.state.mode === EditorMode.EDIT && this.state.clipboard !== null;
  }

  /**
   * 是否可以删除
   */
  get canDelete(): boolean {
    return this.state.mode === EditorMode.EDIT && this.hasSelection;
  }

  /**
   * 当前工具的显示名称
   */
  get currentToolLabel(): string {
    const labels: Record<EditorTool, string> = {
      [EditorTool.SELECT]: '选择',
      [EditorTool.RECTANGLE]: '矩形',
      [EditorTool.CIRCLE]: '圆形',
      [EditorTool.LINE]: '直线',
      [EditorTool.TEXT]: '文本',
      [EditorTool.IMAGE]: '图片',
    };
    return labels[this.state.currentTool];
  }

  /**
   * 当前模式的显示名称
   */
  get currentModeLabel(): string {
    const labels: Record<EditorMode, string> = {
      [EditorMode.EDIT]: '编辑',
      [EditorMode.PREVIEW]: '预览',
      [EditorMode.RUN]: '运行',
    };
    return labels[this.state.mode];
  }

  /**
   * 左侧边栏当前Tab的显示名称
   */
  get leftSidebarTabLabel(): string {
    const labels: Record<LeftSidebarTab, string> = {
      [LeftSidebarTab.PROJECT]: '工程',
      [LeftSidebarTab.SCREEN]: '画面',
      [LeftSidebarTab.COMPONENT]: '组件',
    };
    return labels[this.state.leftSidebarActiveTab];
  }

  /**
   * 右侧边栏当前Tab的显示名称
   */
  get rightSidebarTabLabel(): string | null {
    if (this.state.rightSidebarActiveTab === null) {
      return null;
    }
    const labels: Record<RightSidebarTab, string> = {
      [RightSidebarTab.PROPERTY]: '属性',
      [RightSidebarTab.LAYER]: '图层',
    };
    return labels[this.state.rightSidebarActiveTab];
  }
}

// 创建全局单例
let editorStoreInstance: EditorStore | null = null;

/**
 * 获取EditorStore单例
 */
export const getEditorStore = (): EditorStore => {
  if (!editorStoreInstance) {
    editorStoreInstance = new EditorStore();
  }
  return editorStoreInstance;
};

/**
 * 重置EditorStore（用于测试）
 */
export const resetEditorStore = (): void => {
  editorStoreInstance = null;
};
