"use strict";
/**
 * IPC Handlers Registry
 *
 * IPC处理器注册表 - 集中管理所有IPC处理器
 *
 * FE-006-24: IPC处理器注册
 */
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.IpcHandlersRegistry = void 0;
const electron_1 = require("electron");
const fs = __importStar(require("fs/promises"));
/**
 * IPC处理器注册表类
 */
class IpcHandlersRegistry {
    constructor(mainWindow) {
        this.mainWindow = mainWindow;
    }
    /**
     * 更新主窗口引用
     */
    setMainWindow(window) {
        this.mainWindow = window;
    }
    /**
     * 注册所有IPC处理器
     */
    registerAll() {
        this.registerDialogHandlers();
        this.registerFileSystemHandlers();
        this.registerAppInfoHandlers();
        this.registerUtilityHandlers();
    }
    // ==========================================
    // Dialog Handlers
    // ==========================================
    /**
     * 注册对话框处理器
     */
    registerDialogHandlers() {
        /**
         * 文件对话框：选择保存路径
         */
        electron_1.ipcMain.handle('dialog:selectSavePath', async (_event, options) => {
            if (!this.mainWindow)
                return undefined;
            const result = await electron_1.dialog.showSaveDialog(this.mainWindow, {
                title: options?.title || '选择保存位置',
                defaultPath: options?.defaultPath,
                filters: options?.filters || [
                    { name: 'PanTools工程文件', extensions: ['pant'] },
                    { name: '所有文件', extensions: ['*'] },
                ],
            });
            if (result.canceled || !result.filePath) {
                return undefined;
            }
            // 确保文件扩展名是.pant
            let filePath = result.filePath;
            if (!filePath.endsWith('.pant')) {
                filePath += '.pant';
            }
            return filePath;
        });
        /**
         * 文件对话框：选择打开路径
         */
        electron_1.ipcMain.handle('dialog:selectOpenPath', async (_event, options) => {
            if (!this.mainWindow)
                return undefined;
            const result = await electron_1.dialog.showOpenDialog(this.mainWindow, {
                title: options?.title || '选择工程文件',
                defaultPath: options?.defaultPath,
                filters: options?.filters || [
                    { name: 'PanTools工程文件', extensions: ['pant'] },
                    { name: '所有文件', extensions: ['*'] },
                ],
                properties: ['openFile'],
            });
            if (result.canceled || result.filePaths.length === 0) {
                return undefined;
            }
            return result.filePaths[0];
        });
    }
    // ==========================================
    // File System Handlers
    // ==========================================
    /**
     * 注册文件系统处理器
     */
    registerFileSystemHandlers() {
        /**
         * 文件系统：检查文件是否存在
         */
        electron_1.ipcMain.handle('fs:exists', async (_event, filePath) => {
            try {
                await fs.access(filePath);
                return true;
            }
            catch {
                return false;
            }
        });
        /**
         * 文件系统：读取文件内容
         */
        electron_1.ipcMain.handle('fs:readFile', async (_event, filePath) => {
            try {
                const content = await fs.readFile(filePath, 'utf-8');
                return content;
            }
            catch (error) {
                throw new Error(`读取文件失败: ${error.message}`);
            }
        });
        /**
         * 文件系统：写入文件内容
         */
        electron_1.ipcMain.handle('fs:writeFile', async (_event, filePath, content) => {
            try {
                const { dirname } = require('path');
                // 确保目录存在
                const dir = dirname(filePath);
                await fs.mkdir(dir, { recursive: true });
                // 写入文件
                await fs.writeFile(filePath, content, 'utf-8');
                return { success: true };
            }
            catch (error) {
                throw new Error(`写入文件失败: ${error.message}`);
            }
        });
        /**
         * 文件系统：删除文件
         */
        electron_1.ipcMain.handle('fs:deleteFile', async (_event, filePath) => {
            try {
                await fs.unlink(filePath);
                return { success: true };
            }
            catch (error) {
                throw new Error(`删除文件失败: ${error.message}`);
            }
        });
    }
    // ==========================================
    // App Info Handlers
    // ==========================================
    /**
     * 注册应用信息处理器
     */
    registerAppInfoHandlers() {
        const { app } = require('electron');
        /**
         * 应用信息：获取版本号
         */
        electron_1.ipcMain.handle('app:getVersion', () => {
            return app.getVersion();
        });
        /**
         * 应用信息：获取应用路径
         */
        electron_1.ipcMain.handle('app:getAppPath', () => {
            return app.getAppPath();
        });
        /**
         * 应用信息：获取用户数据目录
         */
        electron_1.ipcMain.handle('app:getUserDataPath', () => {
            return app.getPath('userData');
        });
        /**
         * 应用信息：获取应用名称
         */
        electron_1.ipcMain.handle('app:getName', () => {
            return app.getName();
        });
        /**
         * 应用信息：退出应用
         */
        electron_1.ipcMain.on('app:quit', () => {
            app.quit();
        });
        /**
         * 应用信息：重启应用
         */
        electron_1.ipcMain.on('app:relaunch', () => {
            app.relaunch();
            app.exit();
        });
    }
    // ==========================================
    // Utility Handlers
    // ==========================================
    /**
     * 注册工具处理器
     */
    registerUtilityHandlers() {
        /**
         * 工具：打开外部链接
         */
        electron_1.ipcMain.on('utility:openExternal', (_event, url) => {
            electron_1.shell.openExternal(url);
        });
        /**
         * 工具：在文件管理器中显示
         */
        electron_1.ipcMain.on('utility:showItemInFolder', (_event, filePath) => {
            electron_1.shell.showItemInFolder(filePath);
        });
        /**
         * 工具：获取路径信息
         */
        electron_1.ipcMain.handle('utility:getPath', (_event, name) => {
            const { app } = require('electron');
            return app.getPath(name);
        });
        /**
         * 工具： beep
         */
        electron_1.ipcMain.on('utility:beep', () => {
            electron_1.shell.beep();
        });
        /**
         * 工具：写日志到主进程
         */
        electron_1.ipcMain.on('utility:log', (_event, ...args) => {
            console.log('[Renderer]', ...args);
        });
    }
    /**
     * 注销所有IPC处理器
     */
    unregisterAll() {
        // Electron不提供直接移除handler的方法
        // 但可以通过移除channel来间接实现
        const channels = electron_1.ipcMain.eventNames();
        channels.forEach((channel) => {
            electron_1.ipcMain.removeAllListeners(channel);
        });
    }
}
exports.IpcHandlersRegistry = IpcHandlersRegistry;
/**
 * 默认导出
 */
exports.default = IpcHandlersRegistry;
