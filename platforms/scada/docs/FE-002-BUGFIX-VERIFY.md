# FE-002 Bug修复验证指南

## 问题描述
前端在获取硬件平台列表时出现网络请求失败错误：
```
获取硬件平台列表失败: ApiError: 网络请求失败
```

## 根本原因
后端API路由注册存在冲突：
1. `project_api.go` 创建了 `/api` 路由组
2. `platform_api.go` 也创建了独立的 `/api` 路由组
3. 导致路由冲突，平台API无法正确注册

## 修复内容

### 1. 重构路由注册方式

**修改前：**
```go
// project_api.go
func (api *ProjectAPI) RegisterRoutes(app *fiber.App) {
    app.Use(middleware...)
    apiGroup := app.Group("/api")
    // ...
}

// platform_api.go
func (api *PlatformAPI) RegisterRoutes(app *fiber.App) {
    apiGroup := app.Group("/api")  // ❌ 冲突！
    // ...
}

// main.go
projectAPI.RegisterRoutes(app)
platformAPI.RegisterRoutes(app)  // ❌ 创建了两个 /api 组
```

**修改后：**
```go
// project_api.go
func (api *ProjectAPI) RegisterRoutes(apiGroup fiber.Router) {
    projects := apiGroup.Group("/projects")
    // ...
}

func (api *ProjectAPI) RegisterGlobalMiddleware(app *fiber.App) {
    app.Use(middleware...)
    app.Get("/health", ...)
}

// platform_api.go
func (api *PlatformAPI) RegisterRoutes(apiGroup fiber.Router) {
    apiGroup.Get("/platforms", api.getAllPlatforms)  // ✅ 直接注册到已有路由组
}

// main.go
projectAPI.RegisterGlobalMiddleware(app)
apiGroup := app.Group("/api")  // ✅ 只创建一个 /api 组
projectAPI.RegisterRoutes(apiGroup)
platformAPI.RegisterRoutes(apiGroup)
```

### 2. 修改文件列表

- ✅ `backend/internal/api/project_api.go`: 分离路由注册和全局中间件
- ✅ `backend/internal/api/platform_api.go`: 接受路由组参数，直接注册路由
- ✅ `backend/main.go`: 统一路由组管理

## 验证步骤

### 1. 启动后端服务器

```bash
cd D:\Project\pansiot\platforms\scada\backend
go run main.go
```

期望输出：
```
Database initialized successfully
Database migration completed
Starting PanTools Scada API Server on :3000
```

### 2. 验证健康检查端点

```bash
curl http://localhost:3000/health
```

期望响应：
```json
{
  "status": "ok",
  "message": "Scada Backend API is running"
}
```

### 3. 验证平台API端点

```bash
curl http://localhost:3000/api/platforms
```

期望响应：
```json
{
  "success": true,
  "data": [
    {
      "id": "box1",
      "name": "BOX1",
      "type": "box",
      "resolution": "1920x1080",
      "enabled": true
    },
    {
      "id": "hmi01",
      "name": "HMI01",
      "type": "hmi",
      "resolution": "1280x800",
      "enabled": true
    },
    {
      "id": "tbox1",
      "name": "TBOX1",
      "type": "gateway",
      "resolution": "1024x600",
      "enabled": true
    }
  ]
}
```

### 4. 验证前端集成

#### 方法A：在开发环境（Vite开发服务器）

1. 启动后端服务器（步骤1）
2. 启动前端开发服务器：
   ```bash
   cd D:\Project\pansiot\platforms\scada\packages\renderer
   pnpm dev
   ```
3. 打开浏览器访问 `http://localhost:5173`
4. 点击"新建工程"按钮
5. 检查"运行平台"下拉框是否显示：BOX1, HMI01, TBOX1

#### 方法B：在Electron环境

1. 启动后端服务器（步骤1）
2. 启动Electron应用：
   ```bash
   cd D:\Project\pansiot\platforms\scada\packages\desktop
   pnpm dev
   ```
3. 在Electron窗口中点击"新建工程"按钮
4. 检查"运行平台"下拉框是否显示：BOX1, HMI01, TBOX1

### 5. 检查浏览器控制台

打开浏览器开发者工具（F12），检查：
- ✅ 没有"获取硬件平台列表失败"错误
- ✅ 可以看到成功的API请求：`GET /api/platforms 200`
- ✅ 响应数据包含3个平台

## 常见问题排查

### 问题1：后端无法启动

**检查端口占用：**
```bash
# Windows
netstat -ano | findstr :3000

# 如果被占用，终止进程或修改端口
```

**检查数据库连接：**
确保 SQLite 数据库文件路径正确，应用有读写权限。

### 问题2：API返回404

**验证路由是否正确注册：**
启动后端时，检查日志中是否有错误信息。

**手动测试：**
```bash
# 测试根路径
curl http://localhost:3000/

# 测试健康检查
curl http://localhost:3000/health

# 测试平台API
curl http://localhost:3000/api/platforms
```

### 问题3：前端仍报网络错误

**检查Vite代理配置：**
确认 `vite.config.ts` 中有正确的代理配置：
```typescript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:3000',
      changeOrigin: true,
    }
  }
}
```

**检查API基础URL：**
在浏览器控制台执行：
```javascript
console.log(window.location.origin)
```
- 开发环境应该输出 `http://localhost:5173`
- API请求应该被代理到 `http://localhost:3000`

### 问题4：CORS错误

如果遇到CORS（跨域）错误，确认后端CORS中间件配置正确：
```go
app.Use(cors.New(cors.Config{
    AllowOrigins:     "*",
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
    ExposeHeaders:    "Content-Length",
    AllowCredentials: true,
}))
```

## 成功标准

修复成功的标志：
- ✅ 后端服务器正常启动，监听3000端口
- ✅ `/api/platforms` 端点返回200状态码和正确的JSON数据
- ✅ 前端"新建工程"对话框的"运行平台"下拉框显示3个选项
- ✅ 浏览器控制台没有网络错误
- ✅ 可以正常选择平台并创建工程

## 修复时间
2026-01-22

## 修复人员
Claude Code
