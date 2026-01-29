"use strict";
/**
 * NotificationManager
 *
 * 通知管理器 - 负责系统通知和应用内通知
 *
 * FE-006-22: 通知管理器
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.NotificationManager = void 0;
const electron_1 = require("electron");
/**
 * NotificationManager 类
 */
class NotificationManager {
    constructor(mainWindow) {
        this.mainWindow = null;
        this.history = [];
        this.maxHistorySize = 100;
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
    // Public Methods - 系统通知
    // ==========================================
    /**
     * 显示系统通知
     */
    showSystemNotification(options) {
        const notification = new electron_1.Notification({
            title: options.title,
            body: options.body,
            icon: options.icon,
            silent: options.silent ?? false,
            urgency: options.urgency ?? 'normal',
            timeoutType: options.timeoutType ?? 'default',
        });
        notification.on('click', () => {
            // 通知被点击时聚焦主窗口
            if (this.mainWindow && !this.mainWindow.isDestroyed()) {
                if (this.mainWindow.isMinimized()) {
                    this.mainWindow.restore();
                }
                this.mainWindow.focus();
            }
        });
        notification.on('close', () => {
            // 通知关闭时的处理（如果需要）
        });
        notification.show();
    }
    /**
     * 显示信息通知
     */
    showInfo(title, body) {
        this.showSystemNotification({
            title,
            body,
            urgency: 'normal',
        });
        this.addToHistory('info', title, body);
        this.sendToRenderer('notification:info', { title, body });
    }
    /**
     * 显示成功通知
     */
    showSuccess(title, body) {
        this.showSystemNotification({
            title,
            body,
            urgency: 'normal',
        });
        this.addToHistory('success', title, body);
        this.sendToRenderer('notification:success', { title, body });
    }
    /**
     * 显示警告通知
     */
    showWarning(title, body) {
        this.showSystemNotification({
            title,
            body,
            urgency: 'normal',
        });
        this.addToHistory('warning', title, body);
        this.sendToRenderer('notification:warning', { title, body });
    }
    /**
     * 显示错误通知
     */
    showError(title, body) {
        this.showSystemNotification({
            title,
            body,
            urgency: 'critical',
        });
        this.addToHistory('error', title, body);
        this.sendToRenderer('notification:error', { title, body });
    }
    // ==========================================
    // Public Methods - 通知历史
    // ==========================================
    /**
     * 获取通知历史
     */
    getHistory(type) {
        if (type) {
            return this.history.filter((item) => item.type === type);
        }
        return [...this.history];
    }
    /**
     * 清空通知历史
     */
    clearHistory() {
        this.history = [];
        this.sendToRenderer('notification:historyCleared');
    }
    /**
     * 删除指定通知
     */
    deleteFromHistory(id) {
        this.history = this.history.filter((item) => item.id !== id);
    }
    // ==========================================
    // Private Methods - 历史管理
    // ==========================================
    /**
     * 添加到历史记录
     */
    addToHistory(type, title, body) {
        const notification = {
            id: `notification-${Date.now()}-${Math.random()}`,
            type,
            title,
            body,
            timestamp: Date.now(),
        };
        this.history.push(notification);
        // 限制历史记录大小
        if (this.history.length > this.maxHistorySize) {
            this.history = this.history.slice(-this.maxHistorySize);
        }
    }
    // ==========================================
    // Private Methods - 渲染进程通信
    // ==========================================
    /**
     * 向渲染进程发送通知
     */
    sendToRenderer(channel, data) {
        if (this.mainWindow && !this.mainWindow.isDestroyed()) {
            this.mainWindow.webContents.send(channel, data);
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
         * 显示系统通知
         */
        electron_1.ipcMain.on('notification:show', (_event, options) => {
            this.showSystemNotification(options);
        });
        /**
         * 显示信息通知
         */
        electron_1.ipcMain.on('notification:info', (_event, title, body) => {
            this.showInfo(title, body);
        });
        /**
         * 显示成功通知
         */
        electron_1.ipcMain.on('notification:success', (_event, title, body) => {
            this.showSuccess(title, body);
        });
        /**
         * 显示警告通知
         */
        electron_1.ipcMain.on('notification:warning', (_event, title, body) => {
            this.showWarning(title, body);
        });
        /**
         * 显示错误通知
         */
        electron_1.ipcMain.on('notification:error', (_event, title, body) => {
            this.showError(title, body);
        });
        /**
         * 获取通知历史
         */
        electron_1.ipcMain.handle('notification:getHistory', (_event, type) => {
            return this.getHistory(type);
        });
        /**
         * 清空通知历史
         */
        electron_1.ipcMain.on('notification:clearHistory', () => {
            this.clearHistory();
        });
        /**
         * 删除指定通知
         */
        electron_1.ipcMain.on('notification:delete', (_event, id) => {
            this.deleteFromHistory(id);
        });
    }
    // ==========================================
    // Utility Methods
    // ==========================================
    /**
     * 检查通知权限
     * 注意：Electron 主进程中 Notification.permission 不可用
     * 系统通知权限由操作系统自动管理
     */
    static checkPermission() {
        // Electron 主进程中，假设系统通知可用
        return electron_1.Notification.isSupported() ? 'granted' : 'denied';
    }
    /**
     * 请求通知权限
     * 注意：Electron 主进程中不使用浏览器式的权限请求
     * 系统通知权限由操作系统自动处理
     */
    static async requestPermission() {
        // Electron 主进程中，只需检查通知是否被支持
        return electron_1.Notification.isSupported();
    }
    /**
     * 销毁通知管理器
     */
    destroy() {
        this.clearHistory();
    }
}
exports.NotificationManager = NotificationManager;
/**
 * 默认导出
 */
exports.default = NotificationManager;
