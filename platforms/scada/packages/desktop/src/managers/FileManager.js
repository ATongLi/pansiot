"use strict";
/**
 * FileManager
 *
 * 文件管理器 - 负责工程文件操作 (.pant格式)
 *
 * FE-006-20: 文件管理器
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
exports.FileManager = void 0;
const electron_1 = require("electron");
const path = __importStar(require("path"));
const fs = __importStar(require("fs/promises"));
/**
 * FileManager 类
 */
class FileManager {
    constructor(mainWindow) {
        this.mainWindow = null;
        this.mainWindow = mainWindow;
        this.setupIpcHandlers();
    }
    /**
     * 更新主窗口引用
     */
    setMainWindow(window) {
        this.mainWindow = window;
    }
    // ==========================================
    // Public Methods - 文件对话框
    // ==========================================
    /**
     * 选择保存路径
     */
    async selectSavePath(options = {}) {
        if (!this.mainWindow) {
            throw new Error('主窗口未初始化');
        }
        const result = await electron_1.dialog.showSaveDialog(this.mainWindow, {
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
    async selectOpenPath(options = {}) {
        if (!this.mainWindow) {
            throw new Error('主窗口未初始化');
        }
        const result = await electron_1.dialog.showOpenDialog(this.mainWindow, {
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
    async selectDirectory(options = {}) {
        if (!this.mainWindow) {
            throw new Error('主窗口未初始化');
        }
        const result = await electron_1.dialog.showOpenDialog(this.mainWindow, {
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
    async readProjectFile(filePath) {
        try {
            // 检查文件是否存在
            await fs.access(filePath);
            // 读取文件内容
            const content = await fs.readFile(filePath, 'utf-8');
            // 解析JSON
            const project = JSON.parse(content);
            // 验证文件格式
            if (!project.version || !project.name) {
                throw new Error('无效的工程文件格式');
            }
            return project;
        }
        catch (error) {
            throw new Error(`读取工程文件失败: ${error.message}`);
        }
    }
    /**
     * 写入工程文件
     */
    async writeProjectFile(filePath, project) {
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
        }
        catch (error) {
            throw new Error(`写入工程文件失败: ${error.message}`);
        }
    }
    /**
     * 创建新工程文件
     */
    async createNewProject(name, description) {
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
    async deleteProjectFile(filePath) {
        try {
            await fs.unlink(filePath);
        }
        catch (error) {
            throw new Error(`删除工程文件失败: ${error.message}`);
        }
    }
    /**
     * 检查工程文件是否存在
     */
    async projectFileExists(filePath) {
        try {
            await fs.access(filePath);
            return true;
        }
        catch {
            return false;
        }
    }
    /**
     * 获取工程文件信息
     */
    async getProjectFileInfo(filePath) {
        try {
            const stats = await fs.stat(filePath);
            return {
                name: path.basename(filePath),
                path: filePath,
                size: stats.size,
                modifiedTime: stats.mtime,
            };
        }
        catch (error) {
            throw new Error(`获取文件信息失败: ${error.message}`);
        }
    }
    /**
     * 复制工程文件
     */
    async copyProjectFile(sourcePath, targetPath) {
        try {
            await fs.copyFile(sourcePath, targetPath);
        }
        catch (error) {
            throw new Error(`复制工程文件失败: ${error.message}`);
        }
    }
    /**
     * 重命名工程文件
     */
    async renameProjectFile(oldPath, newPath) {
        try {
            await fs.rename(oldPath, newPath);
        }
        catch (error) {
            throw new Error(`重命名工程文件失败: ${error.message}`);
        }
    }
    /**
     * 在文件管理器中显示文件
     */
    async showFileInFolder(filePath) {
        electron_1.shell.showItemInFolder(filePath);
    }
    /**
     * 获取最近的工程文件列表
     */
    async getRecentProjects(maxCount = 10) {
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
    setupIpcHandlers() {
        /**
         * 文件对话框：选择保存路径
         */
        electron_1.ipcMain.handle('file:selectSavePath', async (_event, options) => {
            return this.selectSavePath(options);
        });
        /**
         * 文件对话框：选择打开路径
         */
        electron_1.ipcMain.handle('file:selectOpenPath', async (_event, options) => {
            return this.selectOpenPath(options);
        });
        /**
         * 文件对话框：选择目录
         */
        electron_1.ipcMain.handle('file:selectDirectory', async (_event, options) => {
            return this.selectDirectory(options);
        });
        /**
         * 工程文件：读取
         */
        electron_1.ipcMain.handle('file:readProject', async (_event, filePath) => {
            return this.readProjectFile(filePath);
        });
        /**
         * 工程文件：写入
         */
        electron_1.ipcMain.handle('file:writeProject', async (_event, filePath, project) => {
            return this.writeProjectFile(filePath, project);
        });
        /**
         * 工程文件：创建新工程
         */
        electron_1.ipcMain.handle('file:createProject', async (_event, name, description) => {
            return this.createNewProject(name, description);
        });
        /**
         * 工程文件：删除
         */
        electron_1.ipcMain.handle('file:deleteProject', async (_event, filePath) => {
            return this.deleteProjectFile(filePath);
        });
        /**
         * 工程文件：检查是否存在
         */
        electron_1.ipcMain.handle('file:projectExists', async (_event, filePath) => {
            return this.projectFileExists(filePath);
        });
        /**
         * 工程文件：获取文件信息
         */
        electron_1.ipcMain.handle('file:getProjectInfo', async (_event, filePath) => {
            return this.getProjectFileInfo(filePath);
        });
        /**
         * 工程文件：复制
         */
        electron_1.ipcMain.handle('file:copyProject', async (_event, sourcePath, targetPath) => {
            return this.copyProjectFile(sourcePath, targetPath);
        });
        /**
         * 工程文件：重命名
         */
        electron_1.ipcMain.handle('file:renameProject', async (_event, oldPath, newPath) => {
            return this.renameProjectFile(oldPath, newPath);
        });
        /**
         * 文件系统：在文件管理器中显示
         */
        electron_1.ipcMain.handle('file:showInFolder', async (_event, filePath) => {
            return this.showFileInFolder(filePath);
        });
        /**
         * 工程文件：获取最近的工程列表
         */
        electron_1.ipcMain.handle('file:getRecentProjects', async (_event, maxCount = 10) => {
            return this.getRecentProjects(maxCount);
        });
    }
    // ==========================================
    // Utility Methods
    // ==========================================
    /**
     * 获取默认文件过滤器
     */
    getDefaultFilters() {
        return [
            { name: 'PanTools工程文件', extensions: ['pant'] },
            { name: '所有文件', extensions: ['*'] },
        ];
    }
    /**
     * 销毁文件管理器
     */
    destroy() {
        // 清理工作（如果需要）
    }
}
exports.FileManager = FileManager;
/**
 * 默认导出
 */
exports.default = FileManager;
