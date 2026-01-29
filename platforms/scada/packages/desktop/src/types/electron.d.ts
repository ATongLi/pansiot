/**
 * Electron API 类型定义
 * 为渲染进程提供 window.electronAPI 的类型提示
 *
 * FE-006-26: 类型定义更新
 */

/**
 * 通知类型
 */
type NotificationType = 'info' | 'success' | 'warning' | 'error';

/**
 * 文件过滤器
 */
interface FileFilter {
  name: string;
  extensions: string[];
}

/**
 * 工程文件接口
 */
interface ProjectFile {
  version: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
  screens: any[];
  components: any[];
  settings: Record<string, any>;
}

/**
 * 自动保存信息接口
 */
interface AutoSaveInfo {
  filePath: string;
  timestamp: number;
  projectName: string;
}

/**
 * 通知历史记录接口
 */
interface NotificationHistory {
  id: string;
  type: NotificationType;
  title: string;
  body: string;
  timestamp: number;
}

/**
 * Electron API 接口
 */
export interface ElectronAPI {
  // ==========================================
  // Dialog API - 文件对话框
  // ==========================================
  dialog: {
    selectSavePath: (options?: {
      title?: string;
      defaultPath?: string;
      filters?: FileFilter[];
    }) => Promise<string | undefined>;

    selectOpenPath: (options?: {
      title?: string;
      defaultPath?: string;
      filters?: FileFilter[];
    }) => Promise<string | undefined>;
  };

  // ==========================================
  // Window API - 窗口控制
  // ==========================================
  window: {
    minimize: () => void;
    maximize: () => void;
    isMaximized: () => Promise<boolean>;
    toggleFullScreen: () => void;
    isFullScreen: () => Promise<boolean>;
    close: () => void;
    reload: () => void;
    forceReload: () => void;
    openDevTools: () => void;
    closeDevTools: () => void;
  };

  // ==========================================
  // File API - 文件操作
  // ==========================================
  file: {
    selectSavePath: (options?: {
      title?: string;
      defaultPath?: string;
      filters?: FileFilter[];
    }) => Promise<string | undefined>;

    selectOpenPath: (options?: {
      title?: string;
      defaultPath?: string;
      filters?: FileFilter[];
    }) => Promise<string | undefined>;

    selectDirectory: (options?: {
      title?: string;
      defaultPath?: string;
    }) => Promise<string | undefined>;

    readProject: (filePath: string) => Promise<ProjectFile>;
    writeProject: (filePath: string, project: ProjectFile) => Promise<void>;
    createProject: (name: string, description?: string) => Promise<ProjectFile>;
    deleteProject: (filePath: string) => Promise<void>;
    projectExists: (filePath: string) => Promise<boolean>;
    getProjectInfo: (filePath: string) => Promise<{
      name: string;
      path: string;
      size: number;
      modifiedTime: Date;
    }>;
    copyProject: (sourcePath: string, targetPath: string) => Promise<void>;
    renameProject: (oldPath: string, newPath: string) => Promise<void>;
    showInFolder: (filePath: string) => Promise<void>;
    getRecentProjects: (maxCount?: number) => Promise<ProjectFile[]>;
  };

  // ==========================================
  // File System API - 底层文件系统
  // ==========================================
  fs: {
    exists: (filePath: string) => Promise<boolean>;
    readFile: (filePath: string) => Promise<string>;
    writeFile: (filePath: string, content: string) => Promise<{ success: boolean }>;
    deleteFile: (filePath: string) => Promise<{ success: boolean }>;
  };

  // ==========================================
  // App API - 应用信息
  // ==========================================
  app: {
    getVersion: () => Promise<string>;
    getAppPath: () => Promise<string>;
    getUserDataPath: () => Promise<string>;
    getName: () => Promise<string>;
    quit: () => void;
    relaunch: () => void;
  };

  // ==========================================
  // Notification API - 通知
  // ==========================================
  notification: {
    show: (options: {
      title: string;
      body: string;
      icon?: string;
      silent?: boolean;
      urgency?: 'normal' | 'critical' | 'low';
    }) => void;

    info: (title: string, body: string) => void;
    success: (title: string, body: string) => void;
    warning: (title: string, body: string) => void;
    error: (title: string, body: string) => void;

    getHistory: (type?: NotificationType) => Promise<NotificationHistory[]>;
    clearHistory: () => void;
    delete: (id: string) => void;
  };

  // ==========================================
  // AutoSave API - 自动保存
  // ==========================================
  autosave: {
    setProject: (projectPath: string, projectData: any) => void;
    clearProject: () => void;
    trigger: () => void;

    getList: () => Promise<AutoSaveInfo[]>;
    restore: (filePath: string) => Promise<any>;
    delete: (filePath: string) => Promise<void>;
    clearAll: () => Promise<void>;

    start: () => void;
    stop: () => void;
    setInterval: (interval: number) => void;
  };

  // ==========================================
  // Utility API - 工具函数
  // ==========================================
  utility: {
    openExternal: (url: string) => void;
    showItemInFolder: (filePath: string) => void;
    getPath: (name: string) => Promise<string>;
    beep: () => void;
    log: (...args: any[]) => void;
  };

  // ==========================================
  // Event Listeners - 事件监听
  // ==========================================
  on: (channel: string, callback: (...args: any[]) => void) => void;
  off: (channel: string, callback: (...args: any[]) => void) => void;
  once: (channel: string, callback: (...args: any[]) => void) => void;
}

/**
 * 扩展Window接口
 */
declare global {
  interface Window {
    electronAPI: ElectronAPI;
  }
}

export {};
