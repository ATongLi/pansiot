/**
 * Electron Main Process
 * Scada Desktop Application 主进程
 *
 * FE-006-27: 主进程集成
 */

import { app, BrowserWindow } from 'electron';
import * as path from 'path';

// 导入管理器
import { WindowManager } from './src/managers/WindowManager';
import { FileManager } from './src/managers/FileManager';
import { MenuManager } from './src/managers/MenuManager';
import { NotificationManager } from './src/managers/NotificationManager';
import { AutoSaveManager } from './src/managers/AutoSaveManager';
import { IpcHandlersRegistry } from './src/ipc/handlers';

/**
 * 管理器实例
 */
let windowManager: WindowManager | null = null;
let fileManager: FileManager | null = null;
let menuManager: MenuManager | null = null;
let notificationManager: NotificationManager | null = null;
let autoSaveManager: AutoSaveManager | null = null;
let ipcHandlersRegistry: IpcHandlersRegistry | null = null;

// 调试输出
console.log('=== PanTools Scada Electron ===');
console.log('NODE_ENV:', process.env.NODE_ENV);
console.log('================================');

/**
 * 检查是否为开发模式
 */
function isDevMode(): boolean {
  return (
    process.env.NODE_ENV === 'development' ||
    process.env.DEBUG_PROD === 'true' ||
    !app.isPackaged
  );
}

/**
 * 初始化管理器
 */
function initializeManagers(mainWindow: BrowserWindow): void {
  // WindowManager 已在 app.whenReady() 中创建，跳过

  // FileManager - 文件管理器
  fileManager = new FileManager(mainWindow);

  // MenuManager - 菜单管理器
  menuManager = new MenuManager(mainWindow);

  // NotificationManager - 通知管理器
  notificationManager = new NotificationManager(mainWindow);

  // AutoSaveManager - 自动保存管理器
  autoSaveManager = new AutoSaveManager(mainWindow, {
    enabled: true,
    interval: 60000, // 60秒
    maxBackups: 5, // 保留5个备份
  });

  // IpcHandlersRegistry - IPC处理器注册表
  ipcHandlersRegistry = new IpcHandlersRegistry(mainWindow);
  ipcHandlersRegistry.registerAll();

  console.log('Main Process: 所有管理器已初始化');
}

/**
 * 清理管理器
 */
async function cleanupManagers(): Promise<void> {
  if (autoSaveManager) {
    await autoSaveManager.destroy();
    autoSaveManager = null;
  }

  if (notificationManager) {
    notificationManager.destroy();
    notificationManager = null;
  }

  if (menuManager) {
    menuManager.destroy();
    menuManager = null;
  }

  if (fileManager) {
    fileManager.destroy();
    fileManager = null;
  }

  if (windowManager) {
    windowManager.destroy();
    windowManager = null;
  }

  console.log('Main Process: 所有管理器已清理');
}

/**
 * 应用就绪事件
 */
app.whenReady().then(() => {
  // 先创建 WindowManager 实例
  if (!windowManager) {
    windowManager = new WindowManager();
  }

  // 创建主窗口
  const mainWindow = windowManager.createMainWindow();

  if (mainWindow) {
    // 初始化管理器
    initializeManagers(mainWindow);
  }

  // macOS: 点击Dock图标时重新创建窗口
  app.on('activate', () => {
    if (windowManager && BrowserWindow.getAllWindows().length === 0) {
      const newWindow = windowManager.createMainWindow();
      // 更新管理器的主窗口引用
      updateManagersWindowReference(newWindow);
    }
  });
});

/**
 * 更新管理器的主窗口引用
 */
function updateManagersWindowReference(window: BrowserWindow | null): void {
  fileManager?.setMainWindow(window);
  menuManager?.setMainWindow(window);
  notificationManager?.setMainWindow(window);
  // autoSaveManager 需要特殊处理，因为它使用 BrowserWindow 类型
  (autoSaveManager as any)?.setMainWindow(window);
  ipcHandlersRegistry?.setMainWindow(window);
}

/**
 * 所有窗口关闭事件
 */
app.on('window-all-closed', async () => {
  // 清理管理器
  await cleanupManagers();

  // macOS: 不要退出应用，只是关闭窗口
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

/**
 * 应用退出前事件
 */
app.on('before-quit', async () => {
  // 清理管理器
  await cleanupManagers();
});

/**
 * 应用退出事件
 */
app.on('will-quit', async (event) => {
  // 阻止默认退出，等待清理完成
  event.preventDefault();

  // 清理管理器
  await cleanupManagers();

  // 退出应用
  app.exit(0);
});

/**
 * 处理未捕获的异常
 */
process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error);
  // 在开发模式下显示错误通知
  if (notificationManager) {
    notificationManager.showError('未捕获的异常', error.message);
  }
});

/**
 * 处理未处理的Promise rejection
 */
process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason);
  // 在开发模式下显示错误通知
  if (notificationManager) {
    notificationManager.showError('未处理的Promise rejection', String(reason));
  }
});

/**
 * 导出管理器实例（用于测试）
 */
export {
  windowManager,
  fileManager,
  menuManager,
  notificationManager,
  autoSaveManager,
  ipcHandlersRegistry,
};
