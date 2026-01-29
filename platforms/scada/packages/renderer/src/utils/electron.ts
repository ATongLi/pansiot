/**
 * Electron API Utilities
 * 提供Electron环境检测和API访问
 */

import type { ElectronAPI } from '@pansiot/scada-desktop/src/types/electron'

/**
 * 检测是否运行在Electron环境中
 */
export const isElectron = (): boolean => {
  // 检查 window.electronAPI 是否存在（通过 preload 脚本暴露）
  // 这是检测 Electron 环境最可靠的方法
  return typeof window !== 'undefined' && !!(window as any).electronAPI
}

/**
 * 获取Electron API
 * 如果在Electron环境中返回真实的API，否则返回mock实现
 */
export const getElectronAPI = (): ElectronAPI => {
  if (isElectron() && window.electronAPI) {
    return window.electronAPI
  }

  // Mock实现用于浏览器开发
  return {
    dialog: {
      selectSavePath: async (options) => {
        console.log('[Mock] selectSavePath:', options)
        // 模拟用户选择路径
        const mockPath = options?.defaultPath || 'C:\\Projects\\new-project.pant'
        return Promise.resolve(mockPath)
      },
      selectOpenPath: async (options) => {
        console.log('[Mock] selectOpenPath:', options)
        const mockPath = 'C:\\Projects\\existing-project.pant'
        return Promise.resolve(mockPath)
      },
    },
    file: {
      selectSavePath: async (options) => {
        console.log('[Mock] file:selectSavePath:', options)
        const mockPath = options?.defaultPath || 'C:\\Projects\\new-project.pant'
        return Promise.resolve(mockPath)
      },
      selectOpenPath: async (options) => {
        console.log('[Mock] file:selectOpenPath:', options)
        const mockPath = 'C:\\Projects\\existing-project.pant'
        return Promise.resolve(mockPath)
      },
      selectDirectory: async (options) => {
        console.log('[Mock] file:selectDirectory:', options)
        return Promise.resolve('C:\\Projects')
      },
      readProject: async (filePath: string) => {
        console.log('[Mock] file:readProject:', filePath)
        // 返回模拟的工程数据
        return Promise.resolve({
          version: '1.0.0',
          name: 'Mock Project',
          description: 'Mock project description',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          screens: [],
          components: [],
          settings: {},
        })
      },
      writeProject: async (filePath: string, project: any) => {
        console.log('[Mock] file:writeProject:', filePath, project)
        return Promise.resolve()
      },
      createProject: async (name: string, description?: string) => {
        console.log('[Mock] file:createProject:', name, description)
        return Promise.resolve({
          version: '1.0.0',
          name,
          description,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          screens: [],
          components: [],
          settings: {},
        })
      },
      deleteProject: async (filePath: string) => {
        console.log('[Mock] file:deleteProject:', filePath)
        return Promise.resolve()
      },
      projectExists: async (filePath: string) => {
        console.log('[Mock] file:projectExists:', filePath)
        return Promise.resolve(false)
      },
      getProjectInfo: async (filePath: string) => {
        console.log('[Mock] file:getProjectInfo:', filePath)
        return Promise.resolve({
          name: 'Mock Project',
          path: filePath,
          size: 1024,
          modifiedTime: new Date(),
        })
      },
      copyProject: async (sourcePath: string, targetPath: string) => {
        console.log('[Mock] file:copyProject:', sourcePath, targetPath)
        return Promise.resolve()
      },
      renameProject: async (oldPath: string, newPath: string) => {
        console.log('[Mock] file:renameProject:', oldPath, newPath)
        return Promise.resolve()
      },
      showInFolder: async (filePath: string) => {
        console.log('[Mock] file:showInFolder:', filePath)
        return Promise.resolve()
      },
      getRecentProjects: async (maxCount = 10) => {
        console.log('[Mock] file:getRecentProjects:', maxCount)
        return Promise.resolve([])
      },
    },
    window: {
      minimize: () => console.log('[Mock] minimize window'),
      maximize: () => console.log('[Mock] maximize window'),
      close: () => console.log('[Mock] close window'),
      isMaximized: async () => {
        console.log('[Mock] isMaximized window')
        return false
      },
    },
    fs: {
      exists: async (path) => {
        console.log('[Mock] fs:exists:', path)
        return Promise.resolve(false)
      },
      readFile: async (path) => {
        console.log('[Mock] fs:readFile:', path)
        return Promise.resolve('{}')
      },
      writeFile: async (path, content) => {
        console.log('[Mock] fs:writeFile:', path)
        return Promise.resolve({ success: true })
      },
      deleteFile: async (path) => {
        console.log('[Mock] fs:deleteFile:', path)
        return Promise.resolve({ success: true })
      },
    },
    app: {
      getVersion: async () => Promise.resolve('1.0.0-mock'),
      getAppPath: async () => Promise.resolve('/mock/app/path'),
    },
    notification: {
      show: (options) => console.log('[Mock] notification:show', options),
      info: (title, body) => console.log('[Mock] notification:info', title, body),
      success: (title, body) => console.log('[Mock] notification:success', title, body),
      warning: (title, body) => console.log('[Mock] notification:warning', title, body),
      error: (title, body) => console.log('[Mock] notification:error', title, body),
      getHistory: async () => Promise.resolve([]),
      clearHistory: () => console.log('[Mock] notification:clearHistory'),
      delete: (id) => console.log('[Mock] notification:delete', id),
    },
    autosave: {
      setProject: (projectPath, projectData) => console.log('[Mock] autosave:setProject', projectPath),
      clearProject: () => console.log('[Mock] autosave:clearProject'),
      trigger: () => console.log('[Mock] autosave:trigger'),
      getList: async () => Promise.resolve([]),
      restore: async (filePath) => Promise.resolve({}),
      delete: (filePath) => console.log('[Mock] autosave:delete', filePath),
      clearAll: () => console.log('[Mock] autosave:clearAll'),
      start: () => console.log('[Mock] autosave:start'),
      stop: () => console.log('[Mock] autosave:stop'),
      setInterval: (interval) => console.log('[Mock] autosave:setInterval', interval),
    },
    utility: {
      openExternal: (url) => console.log('[Mock] utility:openExternal', url),
      showItemInFolder: (filePath) => console.log('[Mock] utility:showItemInFolder', filePath),
      getPath: async (name) => Promise.resolve(`/mock/${name}`),
      beep: () => console.log('[Mock] utility:beep'),
      log: (...args) => console.log('[Mock] utility:log', ...args),
    },
    on: (channel, callback) => console.log('[Mock] on:', channel),
    off: (channel, callback) => console.log('[Mock] off:', channel),
  }
}

/**
 * 获取后端API基础URL
 * 在Electron环境中返回本地服务器地址，在开发环境返回空（使用Vite代理）
 */
export const getApiBaseUrl = (): string => {
  if (isElectron()) {
    // Electron环境：使用本地Golang后端的完整URL
    return 'http://localhost:3000'
  }

  // 开发环境：返回空字符串，因为endpoint已包含/api前缀，
  // Vite代理会自动将/api/*请求转发到http://localhost:3000/api/*
  return ''
}

/**
 * 获取应用环境信息
 */
export const getAppInfo = () => {
  return {
    isElectron: isElectron(),
    platform: typeof window !== 'undefined' ? window.navigator.platform : 'unknown',
    userAgent: typeof window !== 'undefined' ? window.navigator.userAgent : 'unknown',
  }
}
