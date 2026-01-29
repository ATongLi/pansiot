# FE-002 Bug修复：URL配置错误

## 问题描述

前端发起的API请求变成了 `/api/api/platforms`（双重 `/api`），导致404错误。

**错误日志：**
```
12:21:02 | 404 | 0s | 127.0.0.1 | GET | /api/api/platforms | Cannot GET /api/api/platforms
```

## 根本原因

### URL构建逻辑错误

**原有配置：**
```typescript
// electron.ts
export const getApiBaseUrl = (): string => {
  if (isElectron()) {
    return 'http://localhost:3000'  // ✅ 正确
  }
  return '/api'  // ❌ 错误：会重复前缀
}

// platformApi.ts
async getAllPlatforms() {
  return this.request('/api/platforms', {  // endpoint包含/api
    method: 'GET'
  })
}

// 实际URL构建
const url = `${this.baseURL}${endpoint}`
// 开发环境: '/api' + '/api/platforms' = '/api/api/platforms' ❌
// Electron: 'http://localhost:3000' + '/api/platforms' = 'http://localhost:3000/api/platforms' ✅
```

**Vite代理配置：**
```typescript
// vite.config.ts
proxy: {
  '/api': {  // 拦截所有 /api 开头的请求
    target: 'http://localhost:3000',
    changeOrigin: true,
  }
}
```

**问题分析：**
1. endpoint 已包含 `/api` 前缀（如 `/api/platforms`）
2. baseURL 在开发环境返回 `/api`
3. 最终URL = `/api` + `/api/platforms` = `/api/api/platforms` ❌

## 修复方案

### 修改 `getApiBaseUrl()` 函数

**修复后的配置：**
```typescript
// electron.ts
export const getApiBaseUrl = (): string => {
  if (isElectron()) {
    // Electron环境：使用完整URL
    return 'http://localhost:3000'
  }

  // 开发环境：返回空字符串
  // 因为endpoint已包含/api，Vite代理会自动处理
  return ''
}
```

**修复后的URL构建：**
```typescript
const url = `${this.baseURL}${endpoint}`

// 开发环境:
// baseURL = ''
// endpoint = '/api/platforms'
// url = '/api/platforms'
// → Vite拦截并代理到 → http://localhost:3000/api/platforms ✅

// Electron环境:
// baseURL = 'http://localhost:3000'
// endpoint = '/api/platforms'
// url = 'http://localhost:3000/api/platforms' ✅
```

## 工作原理

### 开发环境（Vite）

1. 前端发起请求到 `/api/platforms`
2. Vite开发服务器拦截 `/api/*` 请求
3. Vite代理转发到 `http://localhost:3000/api/platforms`
4. 后端Fiber服务器处理请求

**流程图：**
```
浏览器 → /api/platforms
        ↓
    Vite Proxy (拦截 /api/*)
        ↓
    http://localhost:3000/api/platforms
        ↓
    后端Fiber服务器 → 返回JSON响应
```

### Electron环境

1. 前端发起请求到 `http://localhost:3000/api/platforms`
2. 直接到达后端Fiber服务器
3. 后端处理并返回响应

**流程图：**
```
Electron Renderer → http://localhost:3000/api/platforms
                        ↓
                    后端Fiber服务器
                        ↓
                    返回JSON响应
```

## 验证步骤

### 1. 刷新前端页面

由于Vite的热更新（HMR），修改应该会自动生效。如果没有自动刷新，手动刷新浏览器页面。

### 2. 打开浏览器开发者工具

按 F12 打开开发者工具，切换到 Network 标签。

### 3. 触发API请求

点击"新建工程"按钮，观察Network面板中的请求。

### 4. 验证请求URL

**期望看到：**
```
GET /api/platforms
Status: 200 OK
```

**不应该看到：**
```
❌ GET /api/api/platforms
❌ Status: 404 Not Found
```

### 5. 验证响应数据

点击请求，查看Response标签页，应该包含：
```json
{
  "success": true,
  "data": [
    {"id": "box1", "name": "BOX1", ...},
    {"id": "hmi01", "name": "HMI01", ...},
    {"id": "tbox1", "name": "TBOX1", ...}
  ]
}
```

### 6. 验证下拉框

"运行平台"下拉框应该显示3个选项：
- BOX1
- HMI01
- TBOX1

## 影响范围

### 修改的文件

| 文件 | 修改内容 |
|------|---------|
| `src/utils/electron.ts` | `getApiBaseUrl()` 开发环境返回空字符串 |

### 受影响的API客户端

所有使用 `getApiBaseUrl()` 的API客户端都会受益：
- ✅ `projectApi.ts` - 工程管理API
- ✅ `platformApi.ts` - 硬件平台API
- ✅ 未来新增的任何API客户端

### 向后兼容性

- ✅ Electron环境：无影响，URL构建逻辑保持不变
- ✅ 开发环境：修复了双重 `/api` 问题
- ✅ 生产环境：Electron打包后使用完整URL，不受影响

## 总结

**问题：** 开发环境中URL前缀重复，导致404错误
**原因：** `getApiBaseUrl()` 返回 `/api`，但endpoint也包含 `/api`
**修复：** 开发环境 `getApiBaseUrl()` 返回空字符串
**验证：** 刷新页面，检查Network标签，确认请求URL正确

## 修复时间
2026-01-22

## 修复人员
Claude Code
