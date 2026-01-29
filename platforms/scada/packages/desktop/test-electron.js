/**
 * 简单的Electron测试脚本
 * 用于诊断Electron是否能正常工作
 */

const { app, BrowserWindow } = require('electron')

let mainWindow = null

app.on('ready', () => {
  console.log('=== Electron Ready ===')
  console.log('Node version:', process.versions.node)
  console.log('Chrome version:', process.versions.chrome)
  console.log('Electron version:', process.versions.electron)

  mainWindow = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true
    }
  })

  // 测试1: 加载本地HTML
  console.log('Test 1: Loading local HTML...')
  mainWindow.loadURL('data:text/html,<h1>Electron Test</h1><p>If you see this, Electron works!</p>')

  mainWindow.webContents.on('did-finish-load', () => {
    console.log('✅ Test 1 PASSED: Page loaded successfully')
    console.log('Current URL:', mainWindow.webContents.getURL())

    // 测试2: 延迟2秒后加载Vite
    setTimeout(() => {
      console.log('Test 2: Loading Vite server at http://localhost:5175')
      mainWindow.loadURL('http://localhost:5175')
    }, 2000)
  })

  mainWindow.webContents.on('did-fail-load', (event, errorCode, errorDescription) => {
    console.error('❌ FAILED TO LOAD:', errorCode, errorDescription)
  })

  // 打开开发者工具
  mainWindow.webContents.openDevTools()

  mainWindow.on('closed', () => {
    console.log('Window closed')
    mainWindow = null
  })
})

app.on('window-all-closed', () => {
  console.log('All windows closed, quitting app')
  app.quit()
})

process.on('uncaughtException', (error) => {
  console.error('=== UNCAUGHT EXCEPTION ===')
  console.error(error)
})
