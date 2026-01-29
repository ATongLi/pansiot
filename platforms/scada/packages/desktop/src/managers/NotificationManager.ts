/**
 * NotificationManager
 *
 * 通知管理器 - 负责系统通知和应用内通知
 *
 * FE-006-22: 通知管理器
 */

import { Notification, ipcMain, BrowserWindow } from 'electron';

/**
 * 通知选项接口
 */
interface NotificationOptions {
  title: string;
  body: string;
  icon?: string;
  silent?: boolean;
  urgency?: 'normal' | 'critical' | 'low';
  timeoutType?: 'default' | 'never';
}

/**
 * 通知类型
 */
type NotificationType = 'info' | 'success' | 'warning' | 'error';

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
 * NotificationManager 类
 */
export class NotificationManager {
  private mainWindow: BrowserWindow | null = null;
  private history: NotificationHistory[] = [];
  private maxHistorySize = 100;

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
  // Public Methods - 系统通知
  // ==========================================

  /**
   * 显示系统通知
   */
  showSystemNotification(options: NotificationOptions): void {
    const notification = new Notification({
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
  showInfo(title: string, body: string): void {
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
  showSuccess(title: string, body: string): void {
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
  showWarning(title: string, body: string): void {
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
  showError(title: string, body: string): void {
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
  getHistory(type?: NotificationType): NotificationHistory[] {
    if (type) {
      return this.history.filter((item) => item.type === type);
    }
    return [...this.history];
  }

  /**
   * 清空通知历史
   */
  clearHistory(): void {
    this.history = [];
    this.sendToRenderer('notification:historyCleared');
  }

  /**
   * 删除指定通知
   */
  deleteFromHistory(id: string): void {
    this.history = this.history.filter((item) => item.id !== id);
  }

  // ==========================================
  // Private Methods - 历史管理
  // ==========================================

  /**
   * 添加到历史记录
   */
  private addToHistory(type: NotificationType, title: string, body: string): void {
    const notification: NotificationHistory = {
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
  private sendToRenderer(channel: string, data?: any): void {
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
  private setupIpcHandlers(): void {
    /**
     * 显示系统通知
     */
    ipcMain.on('notification:show', (_event, options: NotificationOptions) => {
      this.showSystemNotification(options);
    });

    /**
     * 显示信息通知
     */
    ipcMain.on('notification:info', (_event, title: string, body: string) => {
      this.showInfo(title, body);
    });

    /**
     * 显示成功通知
     */
    ipcMain.on('notification:success', (_event, title: string, body: string) => {
      this.showSuccess(title, body);
    });

    /**
     * 显示警告通知
     */
    ipcMain.on('notification:warning', (_event, title: string, body: string) => {
      this.showWarning(title, body);
    });

    /**
     * 显示错误通知
     */
    ipcMain.on('notification:error', (_event, title: string, body: string) => {
      this.showError(title, body);
    });

    /**
     * 获取通知历史
     */
    ipcMain.handle('notification:getHistory', (_event, type?: NotificationType) => {
      return this.getHistory(type);
    });

    /**
     * 清空通知历史
     */
    ipcMain.on('notification:clearHistory', () => {
      this.clearHistory();
    });

    /**
     * 删除指定通知
     */
    ipcMain.on('notification:delete', (_event, id: string) => {
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
  static checkPermission(): 'granted' | 'denied' | 'default' {
    // Electron 主进程中，假设系统通知可用
    return Notification.isSupported() ? 'granted' : 'denied';
  }

  /**
   * 请求通知权限
   * 注意：Electron 主进程中不使用浏览器式的权限请求
   * 系统通知权限由操作系统自动处理
   */
  static async requestPermission(): Promise<boolean> {
    // Electron 主进程中，只需检查通知是否被支持
    return Notification.isSupported();
  }

  /**
   * 销毁通知管理器
   */
  destroy(): void {
    this.clearHistory();
  }
}

/**
 * 默认导出
 */
export default NotificationManager;
