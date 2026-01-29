/**
 * AutoSaveManager
 *
 * 自动保存管理器 - 负责工程文件的自动保存和恢复
 *
 * FE-006-23: 自动保存管理器
 */

import { app, ipcMain, BrowserWindow } from 'electron';
import * as path from 'path';
import * as fs from 'fs/promises';

/**
 * 自动保存配置接口
 */
interface AutoSaveConfig {
  enabled: boolean;
  interval: number; // 自动保存间隔（毫秒）
  maxBackups: number; // 最大备份数量
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
 * AutoSaveManager 类
 */
export class AutoSaveManager {
  private mainWindow: BrowserWindow | null = null;
  private config: AutoSaveConfig;
  private autoSaveTimer: NodeJS.Timeout | null = null;
  private autoSaveDir: string;
  private currentProjectPath: string | null = null;
  private currentProjectData: any = null;

  constructor(mainWindow: BrowserWindow | null, config: Partial<AutoSaveConfig> = {}) {
    this.mainWindow = mainWindow;
    this.config = {
      enabled: config.enabled ?? true,
      interval: config.interval ?? 60000, // 默认60秒
      maxBackups: config.maxBackups ?? 5, // 默认保留5个备份
    };
    this.autoSaveDir = path.join(app.getPath('userData'), 'autosave');
    this.initialize();
    this.setupIpcHandlers();
  }

  /**
   * 更新主窗口引用
   */
  setMainWindow(window: BrowserWindow | null): void {
    this.mainWindow = window;
  }

  // ==========================================
  // Public Methods - 初始化和清理
  // ==========================================

  /**
   * 初始化自动保存管理器
   */
  private async initialize(): Promise<void> {
    try {
      // 确保自动保存目录存在
      await fs.mkdir(this.autoSaveDir, { recursive: true });

      // 清理旧的自动保存文件
      await this.cleanupOldAutoSaves();

      // 如果启用自动保存，启动定时器
      if (this.config.enabled) {
        this.startAutoSave();
      }

      console.log('AutoSaveManager: 初始化完成', this.config);
    } catch (error) {
      console.error('AutoSaveManager: 初始化失败', error);
    }
  }

  /**
   * 销毁自动保存管理器
   */
  async destroy(): Promise<void> {
    this.stopAutoSave();
    await this.cleanupOldAutoSaves();
  }

  // ==========================================
  // Public Methods - 自动保存控制
  // ==========================================

  /**
   * 启动自动保存
   */
  startAutoSave(): void {
    if (this.autoSaveTimer) {
      return; // 已经启动
    }

    this.autoSaveTimer = setInterval(async () => {
      await this.performAutoSave();
    }, this.config.interval);

    console.log('AutoSaveManager: 自动保存已启动', `间隔: ${this.config.interval}ms`);
  }

  /**
   * 停止自动保存
   */
  stopAutoSave(): void {
    if (this.autoSaveTimer) {
      clearInterval(this.autoSaveTimer);
      this.autoSaveTimer = null;
      console.log('AutoSaveManager: 自动保存已停止');
    }
  }

  /**
   * 执行自动保存
   */
  async performAutoSave(): Promise<void> {
    // 如果没有当前项目数据，跳过自动保存
    if (!this.currentProjectPath || !this.currentProjectData) {
      return;
    }

    try {
      const projectName = this.currentProjectData.name || '未命名工程';
      const timestamp = Date.now();
      const fileName = `${projectName}_${timestamp}.pant`;
      const filePath = path.join(this.autoSaveDir, fileName);

      // 保存自动保存文件
      await this.saveAutoSaveFile(filePath, this.currentProjectData);

      // 清理旧文件
      await this.cleanupOldAutoSaves();

      // 通知渲染进程
      this.sendToRenderer('autosave:saved', {
        filePath,
        timestamp,
        projectName,
      });

      console.log('AutoSaveManager: 自动保存完成', filePath);
    } catch (error) {
      console.error('AutoSaveManager: 自动保存失败', error);
      this.sendToRenderer('autosave:error', {
        error: error instanceof Error ? error.message : String(error),
      });
    }
  }

  /**
   * 设置当前项目
   */
  setCurrentProject(projectPath: string, projectData: any): void {
    this.currentProjectPath = projectPath;
    this.currentProjectData = projectData;
    console.log('AutoSaveManager: 当前项目已设置', projectPath);
  }

  /**
   * 清除当前项目
   */
  clearCurrentProject(): void {
    this.currentProjectPath = null;
    this.currentProjectData = null;
    console.log('AutoSaveManager: 当前项目已清除');
  }

  // ==========================================
  // Public Methods - 自动保存文件操作
  // ==========================================

  /**
   * 获取所有自动保存文件列表
   */
  async getAutoSaveList(): Promise<AutoSaveInfo[]> {
    try {
      const files = await fs.readdir(this.autoSaveDir);
      const autoSaveFiles: AutoSaveInfo[] = [];

      for (const file of files) {
        if (file.endsWith('.pant')) {
          const filePath = path.join(this.autoSaveDir, file);
          const stats = await fs.stat(filePath);
          const match = file.match(/^(.+)_\d+\.pant$/);
          const projectName = match ? match[1] : file;

          autoSaveFiles.push({
            filePath,
            timestamp: stats.mtimeMs,
            projectName,
          });
        }
      }

      // 按时间戳降序排序
      autoSaveFiles.sort((a, b) => b.timestamp - a.timestamp);

      return autoSaveFiles;
    } catch (error) {
      console.error('AutoSaveManager: 获取自动保存列表失败', error);
      return [];
    }
  }

  /**
   * 读取自动保存文件
   */
  async readAutoSaveFile(filePath: string): Promise<any> {
    try {
      const content = await fs.readFile(filePath, 'utf-8');
      return JSON.parse(content);
    } catch (error) {
      console.error('AutoSaveManager: 读取自动保存文件失败', error);
      throw new Error('读取自动保存文件失败');
    }
  }

  /**
   * 恢复自动保存文件
   */
  async restoreAutoSave(filePath: string): Promise<any> {
    try {
      const projectData = await this.readAutoSaveFile(filePath);

      // 通知渲染进程恢复项目
      this.sendToRenderer('autosave:restored', {
        filePath,
        projectData,
      });

      return projectData;
    } catch (error) {
      console.error('AutoSaveManager: 恢复自动保存失败', error);
      throw new Error('恢复自动保存失败');
    }
  }

  /**
   * 删除自动保存文件
   */
  async deleteAutoSave(filePath: string): Promise<void> {
    try {
      await fs.unlink(filePath);
      this.sendToRenderer('autosave:deleted', { filePath });
      console.log('AutoSaveManager: 自动保存文件已删除', filePath);
    } catch (error) {
      console.error('AutoSaveManager: 删除自动保存文件失败', error);
      throw new Error('删除自动保存文件失败');
    }
  }

  /**
   * 清空所有自动保存文件
   */
  async clearAllAutoSaves(): Promise<void> {
    try {
      const files = await fs.readdir(this.autoSaveDir);
      await Promise.all(
        files.map((file) => fs.unlink(path.join(this.autoSaveDir, file)))
      );
      this.sendToRenderer('autosave:cleared');
      console.log('AutoSaveManager: 所有自动保存文件已清空');
    } catch (error) {
      console.error('AutoSaveManager: 清空自动保存文件失败', error);
      throw new Error('清空自动保存文件失败');
    }
  }

  // ==========================================
  // Private Methods - 自动保存文件管理
  // ==========================================

  /**
   * 保存自动保存文件
   */
  private async saveAutoSaveFile(filePath: string, projectData: any): Promise<void> {
    const content = JSON.stringify(projectData, null, 2);
    await fs.writeFile(filePath, content, 'utf-8');
  }

  /**
   * 清理旧的自动保存文件
   */
  private async cleanupOldAutoSaves(): Promise<void> {
    try {
      const files = await fs.readdir(this.autoSaveDir);
      const pantFiles = files.filter((f) => f.endsWith('.pant'));

      // 如果文件数量超过最大备份数量，删除最旧的文件
      if (pantFiles.length > this.config.maxBackups) {
        const fileStats = await Promise.all(
          pantFiles.map(async (file) => {
            const filePath = path.join(this.autoSaveDir, file);
            const stats = await fs.stat(filePath);
            return { filePath, mtimeMs: stats.mtimeMs };
          })
        );

        // 按修改时间排序，保留最新的文件
        fileStats.sort((a, b) => b.mtimeMs - a.mtimeMs);
        const filesToDelete = fileStats.slice(this.config.maxBackups);

        await Promise.all(filesToDelete.map((item) => fs.unlink(item.filePath)));
        console.log('AutoSaveManager: 已清理旧的自动保存文件', filesToDelete.length);
      }
    } catch (error) {
      console.error('AutoSaveManager: 清理旧文件失败', error);
    }
  }

  // ==========================================
  // Private Methods - IPC Handlers
  // ==========================================

  /**
   * 设置IPC处理器
   */
  private setupIpcHandlers(): void {
    /**
     * 设置当前项目
     */
    ipcMain.on('autosave:setProject', (_event, projectPath: string, projectData: any) => {
      this.setCurrentProject(projectPath, projectData);
    });

    /**
     * 清除当前项目
     */
    ipcMain.on('autosave:clearProject', () => {
      this.clearCurrentProject();
    });

    /**
     * 手动触发自动保存
     */
    ipcMain.on('autosave:trigger', async () => {
      await this.performAutoSave();
    });

    /**
     * 获取自动保存列表
     */
    ipcMain.handle('autosave:getList', async () => {
      return this.getAutoSaveList();
    });

    /**
     * 恢复自动保存
     */
    ipcMain.handle('autosave:restore', async (_event, filePath: string) => {
      return this.restoreAutoSave(filePath);
    });

    /**
     * 删除自动保存文件
     */
    ipcMain.handle('autosave:delete', async (_event, filePath: string) => {
      return this.deleteAutoSave(filePath);
    });

    /**
     * 清空所有自动保存
     */
    ipcMain.handle('autosave:clearAll', async () => {
      return this.clearAllAutoSaves();
    });

    /**
     * 启动自动保存
     */
    ipcMain.on('autosave:start', () => {
      this.startAutoSave();
    });

    /**
     * 停止自动保存
     */
    ipcMain.on('autosave:stop', () => {
      this.stopAutoSave();
    });

    /**
     * 设置自动保存间隔
     */
    ipcMain.on('autosave:setInterval', (_event, interval: number) => {
      this.config.interval = interval;
      // 重启定时器以应用新间隔
      if (this.autoSaveTimer) {
        this.stopAutoSave();
        this.startAutoSave();
      }
    });
  }

  // ==========================================
  // Private Methods - 渲染进程通信
  // ==========================================

  /**
   * 向渲染进程发送消息
   */
  private sendToRenderer(channel: string, data?: any): void {
    if (this.mainWindow && !this.mainWindow.isDestroyed()) {
      this.mainWindow.webContents.send(channel, data);
    }
  }
}

/**
 * 默认导出
 */
export default AutoSaveManager;
