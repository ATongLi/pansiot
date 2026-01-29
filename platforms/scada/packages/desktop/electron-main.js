"use strict";
/**
 * Electron Main Process
 * Scada Desktop Application 主进程
 *
 * FE-006-27: 主进程集成
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.ipcHandlersRegistry = exports.autoSaveManager = exports.notificationManager = exports.menuManager = exports.fileManager = exports.windowManager = void 0;
const electron_1 = require("electron");
// 导入管理器
const WindowManager_1 = require("./src/managers/WindowManager");
const FileManager_1 = require("./src/managers/FileManager");
const MenuManager_1 = require("./src/managers/MenuManager");
const NotificationManager_1 = require("./src/managers/NotificationManager");
const AutoSaveManager_1 = require("./src/managers/AutoSaveManager");
const handlers_1 = require("./src/ipc/handlers");
/**
 * 管理器实例
 */
let windowManager = null;
exports.windowManager = windowManager;
let fileManager = null;
exports.fileManager = fileManager;
let menuManager = null;
exports.menuManager = menuManager;
let notificationManager = null;
exports.notificationManager = notificationManager;
let autoSaveManager = null;
exports.autoSaveManager = autoSaveManager;
let ipcHandlersRegistry = null;
exports.ipcHandlersRegistry = ipcHandlersRegistry;
// 调试输出
console.log('=== PanTools Scada Electron ===');
console.log('NODE_ENV:', process.env.NODE_ENV);
console.log('================================');
/**
 * 检查是否为开发模式
 */
function isDevMode() {
    return (process.env.NODE_ENV === 'development' ||
        process.env.DEBUG_PROD === 'true' ||
        !electron_1.app.isPackaged);
}
/**
 * 初始化管理器
 */
function initializeManagers(mainWindow) {
    // WindowManager 已在 app.whenReady() 中创建，跳过
    // FileManager - 文件管理器
    exports.fileManager = fileManager = new FileManager_1.FileManager(mainWindow);
    // MenuManager - 菜单管理器
    exports.menuManager = menuManager = new MenuManager_1.MenuManager(mainWindow);
    // NotificationManager - 通知管理器
    exports.notificationManager = notificationManager = new NotificationManager_1.NotificationManager(mainWindow);
    // AutoSaveManager - 自动保存管理器
    exports.autoSaveManager = autoSaveManager = new AutoSaveManager_1.AutoSaveManager(mainWindow, {
        enabled: true,
        interval: 60000, // 60秒
        maxBackups: 5, // 保留5个备份
    });
    // IpcHandlersRegistry - IPC处理器注册表
    exports.ipcHandlersRegistry = ipcHandlersRegistry = new handlers_1.IpcHandlersRegistry(mainWindow);
    ipcHandlersRegistry.registerAll();
    console.log('Main Process: 所有管理器已初始化');
}
/**
 * 清理管理器
 */
async function cleanupManagers() {
    if (autoSaveManager) {
        await autoSaveManager.destroy();
        exports.autoSaveManager = autoSaveManager = null;
    }
    if (notificationManager) {
        notificationManager.destroy();
        exports.notificationManager = notificationManager = null;
    }
    if (menuManager) {
        menuManager.destroy();
        exports.menuManager = menuManager = null;
    }
    if (fileManager) {
        fileManager.destroy();
        exports.fileManager = fileManager = null;
    }
    if (windowManager) {
        windowManager.destroy();
        exports.windowManager = windowManager = null;
    }
    console.log('Main Process: 所有管理器已清理');
}
/**
 * 应用就绪事件
 */
electron_1.app.whenReady().then(() => {
    // 先创建 WindowManager 实例
    if (!windowManager) {
        exports.windowManager = windowManager = new WindowManager_1.WindowManager();
    }
    // 创建主窗口
    const mainWindow = windowManager.createMainWindow();
    if (mainWindow) {
        // 初始化管理器
        initializeManagers(mainWindow);
    }
    // macOS: 点击Dock图标时重新创建窗口
    electron_1.app.on('activate', () => {
        if (windowManager && electron_1.BrowserWindow.getAllWindows().length === 0) {
            const newWindow = windowManager.createMainWindow();
            // 更新管理器的主窗口引用
            updateManagersWindowReference(newWindow);
        }
    });
});
/**
 * 更新管理器的主窗口引用
 */
function updateManagersWindowReference(window) {
    fileManager?.setMainWindow(window);
    menuManager?.setMainWindow(window);
    notificationManager?.setMainWindow(window);
    // autoSaveManager 需要特殊处理，因为它使用 BrowserWindow 类型
    autoSaveManager?.setMainWindow(window);
    ipcHandlersRegistry?.setMainWindow(window);
}
/**
 * 所有窗口关闭事件
 */
electron_1.app.on('window-all-closed', async () => {
    // 清理管理器
    await cleanupManagers();
    // macOS: 不要退出应用，只是关闭窗口
    if (process.platform !== 'darwin') {
        electron_1.app.quit();
    }
});
/**
 * 应用退出前事件
 */
electron_1.app.on('before-quit', async () => {
    // 清理管理器
    await cleanupManagers();
});
/**
 * 应用退出事件
 */
electron_1.app.on('will-quit', async (event) => {
    // 阻止默认退出，等待清理完成
    event.preventDefault();
    // 清理管理器
    await cleanupManagers();
    // 退出应用
    electron_1.app.exit(0);
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
