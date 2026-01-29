"use strict";
/**
 * WindowManager
 *
 * 窗口管理器 - 负责窗口创建、状态持久化和窗口状态管理
 *
 * FE-006-19: 窗口管理器
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
exports.WindowManager = void 0;
const electron_1 = require("electron");
const path = __importStar(require("path"));
const fs = __importStar(require("fs"));
/**
 * 默认窗口状态
 */
const DEFAULT_WINDOW_STATE = {
    width: 1200,
    height: 800,
    isMaximized: false,
    isFullScreen: false,
};
/**
 * WindowManager 类
 */
class WindowManager {
    constructor() {
        this.mainWindow = null;
        this.windowState = DEFAULT_WINDOW_STATE;
        this.stateFilePath = path.join(electron_1.app.getPath('userData'), 'window-state.json');
        this.loadWindowState();
        this.setupIpcHandlers();
    }
    // ==========================================
    // Public Methods - 窗口创建和获取
    // ==========================================
    /**
     * 创建主窗口
     */
    createMainWindow() {
        if (this.mainWindow) {
            // 如果窗口已存在，聚焦并返回
            if (!this.mainWindow.isDestroyed()) {
                this.mainWindow.focus();
                return this.mainWindow;
            }
            this.mainWindow = null;
        }
        const config = {
            width: this.windowState.width,
            height: this.windowState.height,
            minWidth: 800,
            minHeight: 600,
            frame: false, // 使用自定义标题栏
            backgroundColor: '#ffffff',
            title: 'PanTools Scada',
            show: false, // 延迟显示，等加载完成后显示
            webPreferences: {
                preload: path.join(__dirname, '../../preload.js'),
                contextIsolation: true,
                nodeIntegration: false,
                sandbox: false,
            },
        };
        // 创建窗口
        this.mainWindow = new electron_1.BrowserWindow(config);
        // 恢复窗口位置（如果存在且有效）
        if (this.windowState.x !== undefined && this.windowState.y !== undefined) {
            const { x, y } = this.windowState;
            // 确保窗口在可见区域内
            const displays = electron_1.screen.getAllDisplays();
            const isValidPosition = displays.some((display) => {
                const area = display.workArea;
                return x >= area.x && x < area.x + area.width && y >= area.y && y < area.y + area.height;
            });
            if (isValidPosition) {
                this.mainWindow.setPosition(x, y);
            }
        }
        // 恢复窗口状态
        if (this.windowState.isMaximized) {
            this.mainWindow.maximize();
        }
        if (this.windowState.isFullScreen) {
            this.mainWindow.setFullScreen(true);
        }
        // 加载应用
        this.loadApp();
        // 监听窗口事件
        this.setupWindowEvents();
        // 移除菜单栏
        this.mainWindow.setMenuBarVisibility(false);
        return this.mainWindow;
    }
    /**
     * 获取主窗口实例
     */
    getMainWindow() {
        return this.mainWindow;
    }
    /**
     * 关闭主窗口
     */
    closeMainWindow() {
        if (this.mainWindow && !this.mainWindow.isDestroyed()) {
            this.mainWindow.close();
        }
    }
    // ==========================================
    // Private Methods - 应用加载
    // ==========================================
    /**
     * 加载应用
     */
    loadApp() {
        if (!this.mainWindow)
            return;
        const isDevMode = this.isDevMode();
        console.log('WindowManager: 加载应用，开发模式:', isDevMode);
        if (isDevMode) {
            // 开发模式：加载Vite开发服务器
            const devServerUrl = 'http://localhost:5173';
            console.log('WindowManager: 加载开发服务器', devServerUrl);
            this.mainWindow.loadURL(devServerUrl);
            this.mainWindow.webContents.openDevTools();
        }
        else {
            // 生产模式：加载打包后的文件
            const indexPath = path.join(__dirname, '../renderer/index.html');
            console.log('WindowManager: 加载生产文件', indexPath);
            this.mainWindow.loadFile(indexPath);
        }
        // 监听加载完成事件
        this.mainWindow.webContents.once('did-finish-load', () => {
            console.log('WindowManager: 窗口加载完成');
        });
        // 监听加载失败事件
        this.mainWindow.webContents.once('did-fail-load', (event, errorCode, errorDescription, validatedURL) => {
            console.error('WindowManager: 窗口加载失败', errorCode, errorDescription, validatedURL);
        });
        // 窗口准备好后显示
        this.mainWindow.once('ready-to-show', () => {
            console.log('WindowManager: 窗口准备显示');
            this.mainWindow?.show();
        });
        // 阻止新窗口打开（在浏览器中打开外部链接）
        this.mainWindow.webContents.setWindowOpenHandler(({ url }) => {
            // eslint-disable-next-line @typescript-eslint/no-var-requires
            const { shell } = require('electron');
            shell.openExternal(url);
            return { action: 'deny' };
        });
    }
    // ==========================================
    // Private Methods - 窗口事件
    // ==========================================
    /**
     * 设置窗口事件监听
     */
    setupWindowEvents() {
        if (!this.mainWindow)
            return;
        // 窗口关闭事件
        this.mainWindow.on('closed', () => {
            this.mainWindow = null;
        });
        // 窗口移动事件 - 保存位置
        this.mainWindow.on('move', () => {
            if (this.mainWindow && !this.mainWindow.isMaximized() && !this.mainWindow.isFullScreen()) {
                const [x, y] = this.mainWindow.getPosition();
                this.windowState.x = x;
                this.windowState.y = y;
            }
        });
        // 窗口大小改变事件 - 保存尺寸
        this.mainWindow.on('resize', () => {
            if (this.mainWindow && !this.mainWindow.isMaximized() && !this.mainWindow.isFullScreen()) {
                const [width, height] = this.mainWindow.getSize();
                this.windowState.width = width;
                this.windowState.height = height;
            }
        });
        // 窗口最大化/还原事件
        this.mainWindow.on('maximize', () => {
            this.windowState.isMaximized = true;
            this.mainWindow?.webContents.send('window:maximized');
        });
        this.mainWindow.on('unmaximize', () => {
            this.windowState.isMaximized = false;
            this.mainWindow?.webContents.send('window:unmaximized');
        });
        // 窗口全屏/退出全屏事件
        this.mainWindow.on('enter-full-screen', () => {
            this.windowState.isFullScreen = true;
            this.mainWindow?.webContents.send('window:fullscreen');
        });
        this.mainWindow.on('leave-full-screen', () => {
            this.windowState.isFullScreen = false;
            this.mainWindow?.webContents.send('window:unfullscreen');
        });
        // 应用退出前保存窗口状态
        electron_1.app.on('before-quit', () => {
            this.saveWindowState();
        });
    }
    // ==========================================
    // Private Methods - 窗口状态持久化
    // ==========================================
    /**
     * 加载窗口状态
     */
    loadWindowState() {
        try {
            if (fs.existsSync(this.stateFilePath)) {
                const data = fs.readFileSync(this.stateFilePath, 'utf-8');
                const savedState = JSON.parse(data);
                // 合并默认状态和保存的状态
                this.windowState = { ...DEFAULT_WINDOW_STATE, ...savedState };
                console.log('WindowManager: 窗口状态已加载', this.windowState);
            }
        }
        catch (error) {
            console.error('WindowManager: 加载窗口状态失败', error);
            this.windowState = DEFAULT_WINDOW_STATE;
        }
    }
    /**
     * 保存窗口状态
     */
    saveWindowState() {
        try {
            // 确保目录存在
            const dir = path.dirname(this.stateFilePath);
            if (!fs.existsSync(dir)) {
                fs.mkdirSync(dir, { recursive: true });
            }
            // 获取当前窗口状态
            if (this.mainWindow && !this.mainWindow.isDestroyed()) {
                const [width, height] = this.mainWindow.getSize();
                const [x, y] = this.mainWindow.getPosition();
                this.windowState.width = width;
                this.windowState.height = height;
                this.windowState.x = x;
                this.windowState.y = y;
                this.windowState.isMaximized = this.mainWindow.isMaximized();
                this.windowState.isFullScreen = this.mainWindow.isFullScreen();
            }
            // 保存到文件
            fs.writeFileSync(this.stateFilePath, JSON.stringify(this.windowState, null, 2));
            console.log('WindowManager: 窗口状态已保存', this.windowState);
        }
        catch (error) {
            console.error('WindowManager: 保存窗口状态失败', error);
        }
    }
    // ==========================================
    // Private Methods - IPC Handlers
    // ==========================================
    /**
     * 设置IPC处理器
     */
    setupIpcHandlers() {
        /**
         * 窗口控制：最小化
         */
        electron_1.ipcMain.on('window:minimize', () => {
            this.mainWindow?.minimize();
        });
        /**
         * 窗口控制：最大化/还原
         */
        electron_1.ipcMain.on('window:maximize', () => {
            if (!this.mainWindow)
                return;
            if (this.mainWindow.isMaximized()) {
                this.mainWindow.unmaximize();
            }
            else {
                this.mainWindow.maximize();
            }
        });
        /**
         * 窗口控制：查询是否最大化
         */
        electron_1.ipcMain.handle('window:isMaximized', () => {
            return this.mainWindow?.isMaximized() || false;
        });
        /**
         * 窗口控制：全屏/退出全屏
         */
        electron_1.ipcMain.on('window:toggleFullScreen', () => {
            if (!this.mainWindow)
                return;
            if (this.mainWindow.isFullScreen()) {
                this.mainWindow.setFullScreen(false);
            }
            else {
                this.mainWindow.setFullScreen(true);
            }
        });
        /**
         * 窗口控制：查询是否全屏
         */
        electron_1.ipcMain.handle('window:isFullScreen', () => {
            return this.mainWindow?.isFullScreen() || false;
        });
        /**
         * 窗口控制：关闭
         */
        electron_1.ipcMain.on('window:close', () => {
            this.mainWindow?.close();
        });
        /**
         * 窗口控制：重启
         */
        electron_1.ipcMain.on('window:reload', () => {
            this.mainWindow?.reload();
        });
        /**
         * 窗口控制：强制刷新
         */
        electron_1.ipcMain.on('window:forceReload', () => {
            this.mainWindow?.webContents.reloadIgnoringCache();
        });
        /**
         * 窗口控制：打开开发者工具
         */
        electron_1.ipcMain.on('window:openDevTools', () => {
            if (this.mainWindow) {
                this.mainWindow.webContents.openDevTools();
            }
        });
        /**
         * 窗口控制：关闭开发者工具
         */
        electron_1.ipcMain.on('window:closeDevTools', () => {
            if (this.mainWindow) {
                this.mainWindow.webContents.closeDevTools();
            }
        });
    }
    // ==========================================
    // Utility Methods
    // ==========================================
    /**
     * 检查是否为开发模式
     */
    isDevMode() {
        return (process.env.NODE_ENV === 'development' ||
            process.env.DEBUG_PROD === 'true' ||
            !electron_1.app.isPackaged);
    }
    /**
     * 销毁窗口管理器
     */
    destroy() {
        this.saveWindowState();
        this.closeMainWindow();
    }
}
exports.WindowManager = WindowManager;
/**
 * 默认导出 - 单例模式
 */
exports.default = WindowManager;
