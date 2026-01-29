/**
 * IPC Handlers Registry
 *
 * IPC处理器注册表 - 集中管理所有IPC处理器
 *
 * FE-006-24: IPC处理器注册
 */

import { ipcMain, dialog, shell } from 'electron';
import * as fs from 'fs/promises';

/**
 * IPC处理器注册表类
 */
export class IpcHandlersRegistry {
  private mainWindow: any;

  constructor(mainWindow: any) {
    this.mainWindow = mainWindow;
  }

  /**
   * 更新主窗口引用
   */
  setMainWindow(window: any): void {
    this.mainWindow = window;
  }

  /**
   * 注册所有IPC处理器
   */
  registerAll(): void {
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
  private registerDialogHandlers(): void {
    /**
     * 文件对话框：选择保存路径
     */
    ipcMain.handle('dialog:selectSavePath', async (_event, options) => {
      if (!this.mainWindow) return undefined;

      const result = await dialog.showSaveDialog(this.mainWindow, {
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
    ipcMain.handle('dialog:selectOpenPath', async (_event, options) => {
      if (!this.mainWindow) return undefined;

      const result = await dialog.showOpenDialog(this.mainWindow, {
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
  private registerFileSystemHandlers(): void {
    /**
     * 文件系统：检查文件是否存在
     */
    ipcMain.handle('fs:exists', async (_event, filePath: string) => {
      try {
        await fs.access(filePath);
        return true;
      } catch {
        return false;
      }
    });

    /**
     * 文件系统：读取文件内容
     */
    ipcMain.handle('fs:readFile', async (_event, filePath: string) => {
      try {
        const content = await fs.readFile(filePath, 'utf-8');
        return content;
      } catch (error: any) {
        throw new Error(`读取文件失败: ${error.message}`);
      }
    });

    /**
     * 文件系统：写入文件内容
     */
    ipcMain.handle('fs:writeFile', async (_event, filePath: string, content: string) => {
      try {
        const { dirname } = require('path');
        // 确保目录存在
        const dir = dirname(filePath);
        await fs.mkdir(dir, { recursive: true });

        // 写入文件
        await fs.writeFile(filePath, content, 'utf-8');
        return { success: true };
      } catch (error: any) {
        throw new Error(`写入文件失败: ${error.message}`);
      }
    });

    /**
     * 文件系统：删除文件
     */
    ipcMain.handle('fs:deleteFile', async (_event, filePath: string) => {
      try {
        await fs.unlink(filePath);
        return { success: true };
      } catch (error: any) {
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
  private registerAppInfoHandlers(): void {
    const { app } = require('electron');

    /**
     * 应用信息：获取版本号
     */
    ipcMain.handle('app:getVersion', () => {
      return app.getVersion();
    });

    /**
     * 应用信息：获取应用路径
     */
    ipcMain.handle('app:getAppPath', () => {
      return app.getAppPath();
    });

    /**
     * 应用信息：获取用户数据目录
     */
    ipcMain.handle('app:getUserDataPath', () => {
      return app.getPath('userData');
    });

    /**
     * 应用信息：获取应用名称
     */
    ipcMain.handle('app:getName', () => {
      return app.getName();
    });

    /**
     * 应用信息：退出应用
     */
    ipcMain.on('app:quit', () => {
      app.quit();
    });

    /**
     * 应用信息：重启应用
     */
    ipcMain.on('app:relaunch', () => {
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
  private registerUtilityHandlers(): void {
    /**
     * 工具：打开外部链接
     */
    ipcMain.on('utility:openExternal', (_event, url: string) => {
      shell.openExternal(url);
    });

    /**
     * 工具：在文件管理器中显示
     */
    ipcMain.on('utility:showItemInFolder', (_event, filePath: string) => {
      shell.showItemInFolder(filePath);
    });

    /**
     * 工具：获取路径信息
     */
    ipcMain.handle('utility:getPath', (_event, name: string) => {
      const { app } = require('electron');
      return app.getPath(name as any);
    });

    /**
     * 工具： beep
     */
    ipcMain.on('utility:beep', () => {
      shell.beep();
    });

    /**
     * 工具：写日志到主进程
     */
    ipcMain.on('utility:log', (_event, ...args: any[]) => {
      console.log('[Renderer]', ...args);
    });
  }

  /**
   * 注销所有IPC处理器
   */
  unregisterAll(): void {
    // Electron不提供直接移除handler的方法
    // 但可以通过移除channel来间接实现
    const channels = ipcMain.eventNames();
    channels.forEach((channel) => {
      ipcMain.removeAllListeners(channel as string);
    });
  }
}

/**
 * 默认导出
 */
export default IpcHandlersRegistry;
