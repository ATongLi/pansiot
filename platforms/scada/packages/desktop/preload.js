"use strict";
/**
 * Electron Preload Script
 * 向渲染进程暴露安全的IPC API
 *
 * FE-006-25: Preload脚本更新
 */
Object.defineProperty(exports, "__esModule", { value: true });
const electron_1 = require("electron");
/**
 * 向渲染进程暴露的API
 */
const electronAPI = {
    // ==========================================
    // Dialog API - 文件对话框
    // ==========================================
    dialog: {
        /**
         * 选择保存文件的路径
         */
        selectSavePath: (options) => electron_1.ipcRenderer.invoke('dialog:selectSavePath', options),
        /**
         * 选择打开文件的路径
         */
        selectOpenPath: (options) => electron_1.ipcRenderer.invoke('dialog:selectOpenPath', options),
    },
    // ==========================================
    // Window API - 窗口控制
    // ==========================================
    window: {
        /**
         * 最小化窗口
         */
        minimize: () => electron_1.ipcRenderer.send('window:minimize'),
        /**
         * 最大化/还原窗口
         */
        maximize: () => electron_1.ipcRenderer.send('window:maximize'),
        /**
         * 判断窗口是否最大化
         */
        isMaximized: () => electron_1.ipcRenderer.invoke('window:isMaximized'),
        /**
         * 全屏/退出全屏
         */
        toggleFullScreen: () => electron_1.ipcRenderer.send('window:toggleFullScreen'),
        /**
         * 判断窗口是否全屏
         */
        isFullScreen: () => electron_1.ipcRenderer.invoke('window:isFullScreen'),
        /**
         * 关闭窗口
         */
        close: () => electron_1.ipcRenderer.send('window:close'),
        /**
         * 重载窗口
         */
        reload: () => electron_1.ipcRenderer.send('window:reload'),
        /**
         * 强制刷新
         */
        forceReload: () => electron_1.ipcRenderer.send('window:forceReload'),
        /**
         * 打开开发者工具
         */
        openDevTools: () => electron_1.ipcRenderer.send('window:openDevTools'),
        /**
         * 关闭开发者工具
         */
        closeDevTools: () => electron_1.ipcRenderer.send('window:closeDevTools'),
    },
    // ==========================================
    // File API - 文件操作
    // ==========================================
    file: {
        /**
         * 选择保存路径
         */
        selectSavePath: (options) => electron_1.ipcRenderer.invoke('file:selectSavePath', options),
        /**
         * 选择打开路径
         */
        selectOpenPath: (options) => electron_1.ipcRenderer.invoke('file:selectOpenPath', options),
        /**
         * 选择目录
         */
        selectDirectory: (options) => electron_1.ipcRenderer.invoke('file:selectDirectory', options),
        /**
         * 读取工程文件
         */
        readProject: (filePath) => electron_1.ipcRenderer.invoke('file:readProject', filePath),
        /**
         * 写入工程文件
         */
        writeProject: (filePath, project) => electron_1.ipcRenderer.invoke('file:writeProject', filePath, project),
        /**
         * 创建新工程
         */
        createProject: (name, description) => electron_1.ipcRenderer.invoke('file:createProject', name, description),
        /**
         * 删除工程文件
         */
        deleteProject: (filePath) => electron_1.ipcRenderer.invoke('file:deleteProject', filePath),
        /**
         * 检查工程文件是否存在
         */
        projectExists: (filePath) => electron_1.ipcRenderer.invoke('file:projectExists', filePath),
        /**
         * 获取工程文件信息
         */
        getProjectInfo: (filePath) => electron_1.ipcRenderer.invoke('file:getProjectInfo', filePath),
        /**
         * 复制工程文件
         */
        copyProject: (sourcePath, targetPath) => electron_1.ipcRenderer.invoke('file:copyProject', sourcePath, targetPath),
        /**
         * 重命名工程文件
         */
        renameProject: (oldPath, newPath) => electron_1.ipcRenderer.invoke('file:renameProject', oldPath, newPath),
        /**
         * 在文件管理器中显示
         */
        showInFolder: (filePath) => electron_1.ipcRenderer.invoke('file:showInFolder', filePath),
        /**
         * 获取最近的工程列表
         */
        getRecentProjects: (maxCount) => electron_1.ipcRenderer.invoke('file:getRecentProjects', maxCount),
    },
    // ==========================================
    // File System API - 底层文件系统
    // ==========================================
    fs: {
        /**
         * 检查文件是否存在
         */
        exists: (filePath) => electron_1.ipcRenderer.invoke('fs:exists', filePath),
        /**
         * 读取文件内容
         */
        readFile: (filePath) => electron_1.ipcRenderer.invoke('fs:readFile', filePath),
        /**
         * 写入文件内容
         */
        writeFile: (filePath, content) => electron_1.ipcRenderer.invoke('fs:writeFile', filePath, content),
        /**
         * 删除文件
         */
        deleteFile: (filePath) => electron_1.ipcRenderer.invoke('fs:deleteFile', filePath),
    },
    // ==========================================
    // App API - 应用信息
    // ==========================================
    app: {
        /**
         * 获取应用版本
         */
        getVersion: () => electron_1.ipcRenderer.invoke('app:getVersion'),
        /**
         * 获取应用路径
         */
        getAppPath: () => electron_1.ipcRenderer.invoke('app:getAppPath'),
        /**
         * 获取用户数据目录
         */
        getUserDataPath: () => electron_1.ipcRenderer.invoke('app:getUserDataPath'),
        /**
         * 获取应用名称
         */
        getName: () => electron_1.ipcRenderer.invoke('app:getName'),
        /**
         * 退出应用
         */
        quit: () => electron_1.ipcRenderer.send('app:quit'),
        /**
         * 重启应用
         */
        relaunch: () => electron_1.ipcRenderer.send('app:relaunch'),
    },
    // ==========================================
    // Notification API - 通知
    // ==========================================
    notification: {
        /**
         * 显示系统通知
         */
        show: (options) => electron_1.ipcRenderer.send('notification:show', options),
        /**
         * 显示信息通知
         */
        info: (title, body) => electron_1.ipcRenderer.send('notification:info', title, body),
        /**
         * 显示成功通知
         */
        success: (title, body) => electron_1.ipcRenderer.send('notification:success', title, body),
        /**
         * 显示警告通知
         */
        warning: (title, body) => electron_1.ipcRenderer.send('notification:warning', title, body),
        /**
         * 显示错误通知
         */
        error: (title, body) => electron_1.ipcRenderer.send('notification:error', title, body),
        /**
         * 获取通知历史
         */
        getHistory: (type) => electron_1.ipcRenderer.invoke('notification:getHistory', type),
        /**
         * 清空通知历史
         */
        clearHistory: () => electron_1.ipcRenderer.send('notification:clearHistory'),
        /**
         * 删除指定通知
         */
        delete: (id) => electron_1.ipcRenderer.send('notification:delete', id),
    },
    // ==========================================
    // AutoSave API - 自动保存
    // ==========================================
    autosave: {
        /**
         * 设置当前项目
         */
        setProject: (projectPath, projectData) => electron_1.ipcRenderer.send('autosave:setProject', projectPath, projectData),
        /**
         * 清除当前项目
         */
        clearProject: () => electron_1.ipcRenderer.send('autosave:clearProject'),
        /**
         * 手动触发自动保存
         */
        trigger: () => electron_1.ipcRenderer.send('autosave:trigger'),
        /**
         * 获取自动保存列表
         */
        getList: () => electron_1.ipcRenderer.invoke('autosave:getList'),
        /**
         * 恢复自动保存
         */
        restore: (filePath) => electron_1.ipcRenderer.invoke('autosave:restore', filePath),
        /**
         * 删除自动保存文件
         */
        delete: (filePath) => electron_1.ipcRenderer.invoke('autosave:delete', filePath),
        /**
         * 清空所有自动保存
         */
        clearAll: () => electron_1.ipcRenderer.invoke('autosave:clearAll'),
        /**
         * 启动自动保存
         */
        start: () => electron_1.ipcRenderer.send('autosave:start'),
        /**
         * 停止自动保存
         */
        stop: () => electron_1.ipcRenderer.send('autosave:stop'),
        /**
         * 设置自动保存间隔
         */
        setInterval: (interval) => electron_1.ipcRenderer.send('autosave:setInterval', interval),
    },
    // ==========================================
    // Utility API - 工具函数
    // ==========================================
    utility: {
        /**
         * 打开外部链接
         */
        openExternal: (url) => electron_1.ipcRenderer.send('utility:openExternal', url),
        /**
         * 在文件管理器中显示
         */
        showItemInFolder: (filePath) => electron_1.ipcRenderer.send('utility:showItemInFolder', filePath),
        /**
         * 获取系统路径
         */
        getPath: (name) => electron_1.ipcRenderer.invoke('utility:getPath', name),
        /**
         * Beep
         */
        beep: () => electron_1.ipcRenderer.send('utility:beep'),
        /**
         * 写日志到主进程
         */
        log: (...args) => electron_1.ipcRenderer.send('utility:log', ...args),
    },
    // ==========================================
    // Event Listeners - 事件监听
    // ==========================================
    on: (channel, callback) => {
        // 只允许特定的IPC通道
        const validChannels = [
            'window:maximized',
            'window:unmaximized',
            'window:fullscreen',
            'window:unfullscreen',
            'notification:info',
            'notification:success',
            'notification:warning',
            'notification:error',
            'autosave:saved',
            'autosave:error',
            'autosave:restored',
            'autosave:deleted',
            'autosave:cleared',
            'autosave:historyCleared',
            'menu:*', // 允许所有菜单事件
        ];
        // 检查通道是否有效（支持通配符）
        const isValid = validChannels.some((validChannel) => {
            if (validChannel.endsWith('*')) {
                const prefix = validChannel.slice(0, -1);
                return channel.startsWith(prefix);
            }
            return channel === validChannel;
        });
        if (isValid) {
            electron_1.ipcRenderer.on(channel, (_event, ...args) => callback(...args));
        }
    },
    /**
     * 移除事件监听器
     */
    off: (channel, callback) => {
        const validChannels = [
            'window:maximized',
            'window:unmaximized',
            'window:fullscreen',
            'window:unfullscreen',
            'notification:info',
            'notification:saved',
            'autosave:saved',
            'menu:*',
        ];
        const isValid = validChannels.some((validChannel) => {
            if (validChannel.endsWith('*')) {
                const prefix = validChannel.slice(0, -1);
                return channel.startsWith(prefix);
            }
            return channel === validChannel;
        });
        if (isValid) {
            electron_1.ipcRenderer.removeListener(channel, callback);
        }
    },
    /**
     * 一次性事件监听器
     */
    once: (channel, callback) => {
        const validChannels = [
            'window:maximized',
            'window:unmaximized',
            'autosave:saved',
        ];
        if (validChannels.includes(channel) || channel.startsWith('menu:')) {
            electron_1.ipcRenderer.once(channel, (_event, ...args) => callback(...args));
        }
    },
};
/**
 * 通过contextBridge暴露给渲染进程
 * 类型定义在 @types/electron.d.ts 中
 */
electron_1.contextBridge.exposeInMainWorld('electronAPI', electronAPI);
