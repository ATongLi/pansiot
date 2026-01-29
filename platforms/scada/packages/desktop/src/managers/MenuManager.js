"use strict";
/**
 * MenuManager
 *
 * 菜单管理器 - 负责应用程序菜单管理
 *
 * FE-006-21: 菜单管理器
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.MenuManager = void 0;
const electron_1 = require("electron");
/**
 * MenuManager 类
 */
class MenuManager {
    constructor(mainWindow) {
        this.mainWindow = null;
        this.mainWindow = mainWindow;
        this.buildMenu();
    }
    /**
     * 更新主窗口引用
     */
    setMainWindow(window) {
        this.mainWindow = window;
    }
    // ==========================================
    // Public Methods - 菜单构建
    // ==========================================
    /**
     * 构建应用程序菜单
     */
    buildMenu() {
        const template = [];
        // macOS: 应用程序菜单
        if (process.platform === 'darwin') {
            template.push(this.createMacOSAppMenu());
        }
        // 文件菜单
        template.push(this.createFileMenu());
        // 编辑菜单
        template.push(this.createEditMenu());
        // 视图菜单
        template.push(this.createViewMenu());
        // 工具菜单
        template.push(this.createToolsMenu());
        // 帮助菜单
        template.push(this.createHelpMenu());
        // 构建菜单
        const menu = electron_1.Menu.buildFromTemplate(template);
        electron_1.Menu.setApplicationMenu(menu);
    }
    /**
     * 创建macOS应用程序菜单
     */
    createMacOSAppMenu() {
        return {
            label: electron_1.app.getName(),
            submenu: [
                { role: 'about', label: '关于 PanTools' },
                { type: 'separator' },
                { role: 'services', label: '服务' },
                { type: 'separator' },
                { role: 'hide', label: '隐藏 PanTools' },
                { role: 'hideOthers', label: '隐藏其他' },
                { role: 'unhide', label: '显示全部' },
                { type: 'separator' },
                { role: 'quit', label: '退出 PanTools' },
            ],
        };
    }
    /**
     * 创建文件菜单
     */
    createFileMenu() {
        return {
            label: '文件',
            submenu: [
                {
                    label: '新建工程',
                    accelerator: 'CmdOrCtrl+N',
                    click: () => this.sendToMainWindow('file:new'),
                },
                {
                    label: '打开工程...',
                    accelerator: 'CmdOrCtrl+O',
                    click: () => this.sendToMainWindow('file:open'),
                },
                {
                    label: '打开最近的',
                    submenu: [
                        {
                            label: '清除最近文件列表',
                            click: () => this.sendToMainWindow('file:clearRecent'),
                        },
                    ],
                },
                { type: 'separator' },
                {
                    label: '保存',
                    accelerator: 'CmdOrCtrl+S',
                    click: () => this.sendToMainWindow('file:save'),
                },
                {
                    label: '另存为...',
                    accelerator: 'CmdOrCtrl+Shift+S',
                    click: () => this.sendToMainWindow('file:saveAs'),
                },
                { type: 'separator' },
                {
                    label: '导入',
                    submenu: [
                        {
                            label: '从文件导入...',
                            click: () => this.sendToMainWindow('file:import'),
                        },
                    ],
                },
                {
                    label: '导出',
                    submenu: [
                        {
                            label: '导出为文件...',
                            click: () => this.sendToMainWindow('file:export'),
                        },
                        {
                            label: '导出为图片...',
                            click: () => this.sendToMainWindow('file:exportImage'),
                        },
                    ],
                },
                { type: 'separator' },
                { type: 'separator' },
                {
                    label: '页面设置...',
                    click: () => this.sendToMainWindow('file:pageSetup'),
                },
                {
                    label: '打印...',
                    accelerator: 'CmdOrCtrl+P',
                    click: () => this.sendToMainWindow('file:print'),
                },
                { type: 'separator' },
                {
                    label: '退出',
                    accelerator: process.platform === 'darwin' ? 'CmdOrCtrl+Q' : 'Alt+F4',
                    click: () => {
                        electron_1.app.quit();
                    },
                },
            ],
        };
    }
    /**
     * 创建编辑菜单
     */
    createEditMenu() {
        return {
            label: '编辑',
            submenu: [
                { role: 'undo', label: '撤销' },
                { role: 'redo', label: '重做' },
                { type: 'separator' },
                { role: 'cut', label: '剪切' },
                { role: 'copy', label: '复制' },
                { role: 'paste', label: '粘贴' },
                { role: 'pasteAndMatchStyle', label: '粘贴并匹配样式' },
                { role: 'delete', label: '删除' },
                { role: 'selectAll', label: '全选' },
                { type: 'separator' },
                {
                    label: '查找...',
                    accelerator: 'CmdOrCtrl+F',
                    click: () => this.sendToMainWindow('edit:find'),
                },
                {
                    label: '替换...',
                    accelerator: 'CmdOrCtrl+H',
                    click: () => this.sendToMainWindow('edit:replace'),
                },
            ],
        };
    }
    /**
     * 创建视图菜单
     */
    createViewMenu() {
        return {
            label: '视图',
            submenu: [
                {
                    label: '工具栏',
                    submenu: [
                        {
                            label: '显示顶部工具栏',
                            type: 'checkbox',
                            checked: true,
                            click: () => this.sendToMainWindow('view:toggleTopToolbar'),
                        },
                        {
                            label: '显示左侧工具栏',
                            type: 'checkbox',
                            checked: true,
                            click: () => this.sendToMainWindow('view:toggleLeftToolbar'),
                        },
                        {
                            label: '显示右侧面板',
                            type: 'checkbox',
                            checked: false,
                            click: () => this.sendToMainWindow('view:toggleRightPanel'),
                        },
                        {
                            label: '显示状态栏',
                            type: 'checkbox',
                            checked: true,
                            click: () => this.sendToMainWindow('view:toggleStatusBar'),
                        },
                    ],
                },
                { type: 'separator' },
                {
                    label: '缩放',
                    submenu: [
                        {
                            label: '放大',
                            accelerator: 'CmdOrCtrl+Plus',
                            click: () => this.sendToMainWindow('view:zoomIn'),
                        },
                        {
                            label: '缩小',
                            accelerator: 'CmdOrCtrl+-',
                            click: () => this.sendToMainWindow('view:zoomOut'),
                        },
                        {
                            label: '重置缩放',
                            accelerator: 'CmdOrCtrl+0',
                            click: () => this.sendToMainWindow('view:resetZoom'),
                        },
                        { type: 'separator' },
                        {
                            label: '25%',
                            click: () => this.sendToMainWindow('view:setZoom', 25),
                        },
                        {
                            label: '50%',
                            click: () => this.sendToMainWindow('view:setZoom', 50),
                        },
                        {
                            label: '75%',
                            click: () => this.sendToMainWindow('view:setZoom', 75),
                        },
                        {
                            label: '100%',
                            click: () => this.sendToMainWindow('view:setZoom', 100),
                        },
                        {
                            label: '125%',
                            click: () => this.sendToMainWindow('view:setZoom', 125),
                        },
                        {
                            label: '150%',
                            click: () => this.sendToMainWindow('view:setZoom', 150),
                        },
                        {
                            label: '200%',
                            click: () => this.sendToMainWindow('view:setZoom', 200),
                        },
                    ],
                },
                { type: 'separator' },
                {
                    label: '网格',
                    submenu: [
                        {
                            label: '显示网格',
                            type: 'checkbox',
                            checked: true,
                            click: () => this.sendToMainWindow('view:toggleGrid'),
                        },
                        {
                            label: '吸附网格',
                            type: 'checkbox',
                            checked: true,
                            click: () => this.sendToMainWindow('view:toggleSnap'),
                        },
                        { type: 'separator' },
                        {
                            label: '网格大小',
                            submenu: [
                                {
                                    label: '5px',
                                    click: () => this.sendToMainWindow('view:setGridSize', 5),
                                },
                                {
                                    label: '10px',
                                    click: () => this.sendToMainWindow('view:setGridSize', 10),
                                },
                                {
                                    label: '20px',
                                    click: () => this.sendToMainWindow('view:setGridSize', 20),
                                },
                            ],
                        },
                    ],
                },
                { type: 'separator' },
                {
                    label: '全屏',
                    accelerator: process.platform === 'darwin' ? 'Ctrl+Command+F' : 'F11',
                    role: 'togglefullscreen',
                },
                {
                    label: '开发者工具',
                    accelerator: 'CmdOrCtrl+Shift+I',
                    click: () => {
                        if (this.mainWindow) {
                            this.mainWindow.webContents.toggleDevTools();
                        }
                    },
                },
                {
                    label: '重载',
                    accelerator: 'CmdOrCtrl+R',
                    click: () => {
                        this.mainWindow?.reload();
                    },
                },
                {
                    label: '强制重载',
                    accelerator: 'CmdOrCtrl+Shift+R',
                    click: () => {
                        this.mainWindow?.webContents.reloadIgnoringCache();
                    },
                },
            ],
        };
    }
    /**
     * 创建工具菜单
     */
    createToolsMenu() {
        return {
            label: '工具',
            submenu: [
                {
                    label: '选择工具',
                    accelerator: 'V',
                    click: () => this.sendToMainWindow('tool:select'),
                },
                {
                    label: '矩形工具',
                    accelerator: 'R',
                    click: () => this.sendToMainWindow('tool:rectangle'),
                },
                {
                    label: '圆形工具',
                    accelerator: 'C',
                    click: () => this.sendToMainWindow('tool:circle'),
                },
                {
                    label: '直线工具',
                    accelerator: 'L',
                    click: () => this.sendToMainWindow('tool:line'),
                },
                {
                    label: '文本工具',
                    accelerator: 'T',
                    click: () => this.sendToMainWindow('tool:text'),
                },
                {
                    label: '图片工具',
                    accelerator: 'I',
                    click: () => this.sendToMainWindow('tool:image'),
                },
                { type: 'separator' },
                {
                    label: '编辑模式',
                    accelerator: 'F1',
                    click: () => this.sendToMainWindow('mode:edit'),
                },
                {
                    label: '预览模式',
                    accelerator: 'F2',
                    click: () => this.sendToMainWindow('mode:preview'),
                },
                {
                    label: '运行模式',
                    accelerator: 'F5',
                    click: () => this.sendToMainWindow('mode:run'),
                },
                { type: 'separator' },
                {
                    label: '设置...',
                    click: () => this.sendToMainWindow('tools:settings'),
                },
            ],
        };
    }
    /**
     * 创建帮助菜单
     */
    createHelpMenu() {
        return {
            role: 'help',
            submenu: [
                {
                    label: '文档',
                    click: () => this.sendToMainWindow('help:documentation'),
                },
                {
                    label: '报告问题',
                    click: () => this.sendToMainWindow('help:reportIssue'),
                },
                {
                    label: '快捷键',
                    accelerator: 'CmdOrCtrl+/',
                    click: () => this.sendToMainWindow('help:shortcuts'),
                },
                { type: 'separator' },
                {
                    label: '关于 PanTools',
                    click: () => {
                        if (process.platform === 'darwin') {
                            electron_1.app.showAboutPanel();
                        }
                        else {
                            this.sendToMainWindow('help:about');
                        }
                    },
                },
                {
                    label: '检查更新',
                    click: () => this.sendToMainWindow('help:checkUpdates'),
                },
            ],
        };
    }
    // ==========================================
    // Utility Methods
    // ==========================================
    /**
     * 向主窗口发送消息
     */
    sendToMainWindow(channel, ...args) {
        if (this.mainWindow && !this.mainWindow.isDestroyed()) {
            this.mainWindow.webContents.send(`menu:${channel}`, ...args);
        }
    }
    /**
     * 销毁菜单管理器
     */
    destroy() {
        // 清理工作（如果需要）
    }
}
exports.MenuManager = MenuManager;
/**
 * 默认导出
 */
exports.default = MenuManager;
