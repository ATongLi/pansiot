/**
 * FileManager
 *
 * 文件管理器 - 负责工程文件操作 (.pant格式)
 *
 * FE-006-20: 文件管理器
 */

import { dialog, BrowserWindow, ipcMain, shell } from 'electron';
import * as path from 'path';
import * as fs from 'fs/promises';

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
 * 文件过滤器接口
 */
interface FileFilter {
  name: string;
  extensions: string[];
}

/**
 * 文件对话框选项接口
 */
interface DialogOptions {
  title?: string;
  defaultPath?: string;
  filters?: FileFilter[];
}

/**
 * FileManager 类
 */
export class FileManager {
  private mainWindow: BrowserWindow | null = null;

  constructor(mainWindow: BrowserWindow | null) {
    this.mainWindow = mainWindow;
    this.setupIpcHandlers();
  }

  /**
   * 更新主窗口引用
   */
  setMainWindow(window: BrowserWindow | null): void {
    this.mainWindow = window;
  }

  // ==========================================
  // Public Methods - 文件对话框
  // ==========================================

  /**
   * 选择保存路径
   */
  async selectSavePath(options: DialogOptions = {}): Promise<string | undefined> {
    if (!this.mainWindow) {
      throw new Error('主窗口未初始化');
    }

    const result = await dialog.showSaveDialog(this.mainWindow, {
      title: options.title || '保存工程',
      defaultPath: options.defaultPath,
      filters: options.filters || this.getDefaultFilters(),
    });

    if (result.canceled || !result.filePath) {
      return undefined;
    }

    // 确保文件扩展名是 .pant
    let filePath = result.filePath;
    if (!filePath.endsWith('.pant')) {
      filePath += '.pant';
    }

    return filePath;
  }

  /**
   * 选择打开路径
   */
  async selectOpenPath(options: DialogOptions = {}): Promise<string | undefined> {
    if (!this.mainWindow) {
      throw new Error('主窗口未初始化');
    }

    const result = await dialog.showOpenDialog(this.mainWindow, {
      title: options.title || '打开工程',
      defaultPath: options.defaultPath,
      filters: options.filters || this.getDefaultFilters(),
      properties: ['openFile'],
    });

    if (result.canceled || result.filePaths.length === 0) {
      return undefined;
    }

    return result.filePaths[0];
  }

  /**
   * 选择目录
   */
  async selectDirectory(options: { title?: string; defaultPath?: string } = {}): Promise<string | undefined> {
    if (!this.mainWindow) {
      throw new Error('主窗口未初始化');
    }

    const result = await dialog.showOpenDialog(this.mainWindow, {
      title: options.title || '选择目录',
      defaultPath: options.defaultPath,
      properties: ['openDirectory', 'createDirectory'],
    });

    if (result.canceled || result.filePaths.length === 0) {
      return undefined;
    }

    return result.filePaths[0];
  }

  // ==========================================
  // Public Methods - 工程文件操作
  // ==========================================

  /**
   * 读取工程文件
   */
  async readProjectFile(filePath: string): Promise<ProjectFile> {
    try {
      // 检查文件是否存在
      await fs.access(filePath);

      // 读取文件内容
      const content = await fs.readFile(filePath, 'utf-8');

      // 解析JSON
      const project = JSON.parse(content) as ProjectFile;

      // 验证文件格式
      if (!project.version || !project.name) {
        throw new Error('无效的工程文件格式');
      }

      return project;
    } catch (error: any) {
      throw new Error(`读取工程文件失败: ${error.message}`);
    }
  }

  /**
   * 写入工程文件
   */
  async writeProjectFile(filePath: string, project: ProjectFile): Promise<void> {
    try {
      // 确保目录存在
      const dir = path.dirname(filePath);
      await fs.mkdir(dir, { recursive: true });

      // 更新时间戳
      project.updatedAt = new Date().toISOString();

      // 如果是新文件，设置创建时间
      if (!project.createdAt) {
        project.createdAt = project.updatedAt;
      }

      // 写入文件
      const content = JSON.stringify(project, null, 2);
      await fs.writeFile(filePath, content, 'utf-8');
    } catch (error: any) {
      throw new Error(`写入工程文件失败: ${error.message}`);
    }
  }

  /**
   * 创建新工程文件
   */
  async createNewProject(name: string, description?: string): Promise<ProjectFile> {
    const now = new Date().toISOString();

    return {
      version: '1.0.0',
      name,
      description,
      createdAt: now,
      updatedAt: now,
      screens: [],
      components: [],
      settings: {
        zoom: 100,
        gridSize: 10,
        snapToGrid: true,
        showGrid: true,
      },
    };
  }

  /**
   * 删除工程文件
   */
  async deleteProjectFile(filePath: string): Promise<void> {
    try {
      await fs.unlink(filePath);
    } catch (error: any) {
      throw new Error(`删除工程文件失败: ${error.message}`);
    }
  }

  /**
   * 检查工程文件是否存在
   */
  async projectFileExists(filePath: string): Promise<boolean> {
    try {
      await fs.access(filePath);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * 获取工程文件信息
   */
  async getProjectFileInfo(filePath: string): Promise<{ name: string; path: string; size: number; modifiedTime: Date }> {
    try {
      const stats = await fs.stat(filePath);
      return {
        name: path.basename(filePath),
        path: filePath,
        size: stats.size,
        modifiedTime: stats.mtime,
      };
    } catch (error: any) {
      throw new Error(`获取文件信息失败: ${error.message}`);
    }
  }

  /**
   * 复制工程文件
   */
  async copyProjectFile(sourcePath: string, targetPath: string): Promise<void> {
    try {
      await fs.copyFile(sourcePath, targetPath);
    } catch (error: any) {
      throw new Error(`复制工程文件失败: ${error.message}`);
    }
  }

  /**
   * 重命名工程文件
   */
  async renameProjectFile(oldPath: string, newPath: string): Promise<void> {
    try {
      await fs.rename(oldPath, newPath);
    } catch (error: any) {
      throw new Error(`重命名工程文件失败: ${error.message}`);
    }
  }

  /**
   * 在文件管理器中显示文件
   */
  async showFileInFolder(filePath: string): Promise<void> {
    shell.showItemInFolder(filePath);
  }

  /**
   * 获取最近的工程文件列表
   */
  async getRecentProjects(maxCount = 10): Promise<ProjectFile[]> {
    // TODO: 实现最近文件列表功能
    // 可以从 app.getPath('userData')/recent-projects.json 读取
    return [];
  }

  // ==========================================
  // Private Methods - IPC Handlers
  // ==========================================

  /**
   * 设置IPC处理器
   */
  private setupIpcHandlers(): void {
    /**
     * 文件对话框：选择保存路径
     */
    ipcMain.handle('file:selectSavePath', async (_event, options: DialogOptions) => {
      return this.selectSavePath(options);
    });

    /**
     * 文件对话框：选择打开路径
     */
    ipcMain.handle('file:selectOpenPath', async (_event, options: DialogOptions) => {
      return this.selectOpenPath(options);
    });

    /**
     * 文件对话框：选择目录
     */
    ipcMain.handle('file:selectDirectory', async (_event, options: { title?: string; defaultPath?: string }) => {
      return this.selectDirectory(options);
    });

    /**
     * 工程文件：读取
     */
    ipcMain.handle('file:readProject', async (_event, filePath: string) => {
      return this.readProjectFile(filePath);
    });

    /**
     * 工程文件：写入
     */
    ipcMain.handle('file:writeProject', async (_event, filePath: string, project: ProjectFile) => {
      return this.writeProjectFile(filePath, project);
    });

    /**
     * 工程文件：创建新工程
     */
    ipcMain.handle('file:createProject', async (_event, name: string, description?: string) => {
      return this.createNewProject(name, description);
    });

    /**
     * 工程文件：删除
     */
    ipcMain.handle('file:deleteProject', async (_event, filePath: string) => {
      return this.deleteProjectFile(filePath);
    });

    /**
     * 工程文件：检查是否存在
     */
    ipcMain.handle('file:projectExists', async (_event, filePath: string) => {
      return this.projectFileExists(filePath);
    });

    /**
     * 工程文件：获取文件信息
     */
    ipcMain.handle('file:getProjectInfo', async (_event, filePath: string) => {
      return this.getProjectFileInfo(filePath);
    });

    /**
     * 工程文件：复制
     */
    ipcMain.handle('file:copyProject', async (_event, sourcePath: string, targetPath: string) => {
      return this.copyProjectFile(sourcePath, targetPath);
    });

    /**
     * 工程文件：重命名
     */
    ipcMain.handle('file:renameProject', async (_event, oldPath: string, newPath: string) => {
      return this.renameProjectFile(oldPath, newPath);
    });

    /**
     * 文件系统：在文件管理器中显示
     */
    ipcMain.handle('file:showInFolder', async (_event, filePath: string) => {
      return this.showFileInFolder(filePath);
    });

    /**
     * 工程文件：获取最近的工程列表
     */
    ipcMain.handle('file:getRecentProjects', async (_event, maxCount = 10) => {
      return this.getRecentProjects(maxCount);
    });
  }

  // ==========================================
  // Utility Methods
  // ==========================================

  /**
   * 获取默认文件过滤器
   */
  private getDefaultFilters(): FileFilter[] {
    return [
      { name: 'PanTools工程文件', extensions: ['pant'] },
      { name: '所有文件', extensions: ['*'] },
    ];
  }

  /**
   * 销毁文件管理器
   */
  destroy(): void {
    // 清理工作（如果需要）
  }
}

/**
 * 默认导出
 */
export default FileManager;
