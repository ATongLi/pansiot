# Scada 工程管理 API 文档

## 概述

本文档描述 Scada 工程管理后端 API 的所有端点、请求/响应格式和错误码。

**基础 URL**:
- 开发环境: `http://localhost:3000`
- Electron 环境: `http://localhost:3000`
- 浏览器环境: `/api` (通过 Vite 代理)

**Content-Type**: `application/json`

---

## API 端点

### 1. 健康检查

检查后端服务是否正常运行。

**端点**: `GET /health`

**请求**: 无需参数

**响应**:
```json
{
  "status": "ok",
  "message": "Scada Backend API is running"
}
```

**状态码**: `200 OK`

---

### 2. 创建工程

创建新的工程文件并保存到指定路径。

**端点**: `POST /api/projects/create`

**请求体**:
```json
{
  "metadata": {
    "name": "工程名称",
    "author": "作者名称（可选）",
    "description": "工程描述（可选）",
    "category": "分类1",
    "platform": "HMI型号1",
    "createdAt": "2026-01-21T10:00:00Z",
    "updatedAt": "2026-01-21T10:00:00Z"
  },
  "security": {
    "encrypted": false,
    "password": "密码（如果加密）"
  },
  "savePath": "C:\\Projects\\new-project.pant"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| metadata.name | string | 是 | 工程名称，最长50字符 |
| metadata.author | string | 否 | 作者名称，最长30字符 |
| metadata.description | string | 否 | 工程描述，最长500字符 |
| metadata.category | string | 是 | 工程分类 |
| metadata.platform | string | 是 | 硬件平台 |
| security.encrypted | boolean | 是 | 是否启用加密 |
| security.password | string | 条件 | 加密时必填，最少6字符 |
| savePath | string | 是 | 保存路径，必须以 `.pant` 结尾 |

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "projectId": "uuid-v4",
    "filePath": "C:\\Projects\\new-project.pant"
  }
}
```

**错误响应** (`500 Internal Server Error`):
```json
{
  "success": false,
  "error": "CREATE_FAILED",
  "message": "创建工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 创建成功
- `400 Bad Request`: 请求格式错误
- `500 Internal Server Error`: 服务器内部错误

---

### 3. 打开工程

打开已存在的工程文件。

**端点**: `POST /api/projects/open`

**请求体**:
```json
{
  "filePath": "C:\\Projects\\existing-project.pant",
  "password": "密码（如果工程加密）"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| filePath | string | 是 | 工程文件路径 |
| password | string | 条件 | 加密工程必填 |

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "project": {
      "version": "1.0.0",
      "projectId": "uuid-v4",
      "metadata": { ... },
      "security": { ... },
      "canvas": { ... },
      "components": [ ... ]
    }
  }
}
```

**错误响应**:

#### 密码错误 (`401 Unauthorized`)
```json
{
  "success": false,
  "error": "INVALID_PASSWORD",
  "message": "密码错误"
}
```

#### 签名验证失败 (`400 Bad Request`)
```json
{
  "success": false,
  "error": "INVALID_SIGNATURE",
  "message": "工程文件签名验证失败，文件可能已被篡改"
}
```

#### 其他错误 (`500 Internal Server Error`)
```json
{
  "success": false,
  "error": "OPEN_FAILED",
  "message": "打开工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 打开成功
- `400 Bad Request`: 签名验证失败
- `401 Unauthorized`: 密码错误
- `500 Internal Server Error`: 文件不存在或其他服务器错误

---

### 4. 保存工程

保存当前工程的修改。

**端点**: `POST /api/projects/save`

**请求体**:
```json
{
  "project": {
    "version": "1.0.0",
    "projectId": "uuid-v4",
    "metadata": { ... },
    "security": { ... },
    "canvas": { ... },
    "components": [ ... ],
    "encryptedContent": "加密内容（如果加密）"
  }
}
```

**字段说明**: 完整的工程对象（见"数据模型"章节）

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "filePath": "C:\\Projects\\existing-project.pant"
  }
}
```

**错误响应** (`500 Internal Server Error`):
```json
{
  "success": false,
  "error": "SAVE_FAILED",
  "message": "保存工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 保存成功
- `400 Bad Request`: 请求格式错误
- `500 Internal Server Error`: 服务器内部错误

---

### 5. 验证密码

验证工程密码是否正确（不解密整个工程）。

**端点**: `POST /api/projects/validate-password`

**请求体**:
```json
{
  "filePath": "C:\\Projects\\encrypted-project.pant",
  "password": "用户输入的密码"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| filePath | string | 是 | 工程文件路径 |
| password | string | 是 | 待验证的密码 |

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "valid": true
  }
}
```

**密码错误** (`200 OK` 但 `valid: false`):
```json
{
  "success": true,
  "data": {
    "valid": false
  }
}
```

**错误响应** (`500 Internal Server Error`):
```json
{
  "success": false,
  "error": "VALIDATE_FAILED",
  "message": "验证密码失败的具体原因"
}
```

**状态码**:
- `200 OK`: 验证完成（检查 `data.valid` 字段）
- `500 Internal Server Error`: 服务器内部错误

**使用场景**:
- 用户打开加密工程前先验证密码
- 避免多次错误输入导致的性能问题

---

### 6. 获取最近工程列表

获取最近打开的工程列表。

**端点**: `GET /api/projects/recent`

**请求**: 无需参数

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": [
    {
      "projectId": "uuid-v4",
      "name": "工程名称",
      "category": "分类1",
      "filePath": "C:\\Projects\\project.pant",
      "lastOpened": "2026-01-21T10:00:00Z",
      "isEncrypted": false,
      "createdAt": "2026-01-20T15:30:00Z"
    },
    ...
  ]
}
```

**字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| projectId | string | 工程唯一ID |
| name | string | 工程名称 |
| category | string | 工程分类（可能为空） |
| filePath | string | 文件路径 |
| lastOpened | string (ISO8601) | 最后打开时间 |
| isEncrypted | boolean | 是否加密 |
| createdAt | string (ISO8601) | 创建时间 |

**错误响应** (`500 Internal Server Error`):
```json
{
  "success": false,
  "error": "GET_RECENT_FAILED",
  "message": "获取最近工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 获取成功
- `500 Internal Server Error`: 服务器内部错误

**排序**: 返回结果按 `lastOpened` 降序排列（最近打开的在前）

---

### 7. 添加或更新最近工程

添加新工程到最近列表，或更新现有工程的最后打开时间。

**端点**: `POST /api/projects/recent`

**请求体**:
```json
{
  "projectId": "uuid-v4",
  "name": "工程名称",
  "category": "分类1",
  "filePath": "C:\\Projects\\project.pant",
  "lastOpened": "2026-01-21T10:00:00Z",
  "isEncrypted": false,
  "createdAt": "2026-01-20T15:30:00Z"
}
```

**字段说明**: 与"获取最近工程列表"响应中的字段相同

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "success": true
  }
}
```

**错误响应** (`500 Internal Server Error`):
```json
{
  "success": false,
  "error": "ADD_RECENT_FAILED",
  "message": "添加最近工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 操作成功
- `400 Bad Request`: 请求格式错误
- `500 Internal Server Error`: 服务器内部错误

**自动调用**:
- 创建工程成功后自动调用
- 打开工程成功后自动调用（更新 lastOpened 时间）

---

### 8. 删除最近工程

从最近工程列表中删除指定工程。

**端点**: `DELETE /api/projects/recent/:projectId`

**URL 参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| projectId | string | 是 | 工程ID |

**请求**: 无需请求体

**成功响应** (`200 OK`):
```json
{
  "success": true,
  "data": {
    "success": true
  }
}
```

**错误响应**:

#### 项目ID为空 (`400 Bad Request`)
```json
{
  "success": false,
  "error": "INVALID_PROJECT_ID",
  "message": "工程ID不能为空"
}
```

#### 服务器错误 (`500 Internal Server Error`)
```json
{
  "success": false,
  "error": "REMOVE_RECENT_FAILED",
  "message": "删除最近工程失败的具体原因"
}
```

**状态码**:
- `200 OK`: 删除成功
- `400 Bad Request`: 项目ID无效
- `500 Internal Server Error`: 服务器内部错误

**注意**: 此操作只从数据库删除记录，**不删除**实际工程文件

---

## 数据模型

### Project（工程对象）

完整的工程数据结构。

```typescript
interface Project {
  version: string                  // 版本号，当前 "1.0.0"
  projectId: string                // 工程唯一ID（UUID v4）
  metadata: ProjectMetadata        // 元数据
  security: ProjectSecurity        // 安全配置
  canvas: CanvasConfig             // 画布配置
  components: Component[]          // 组件列表
  encryptedContent?: string        // 加密内容（如果启用加密）
}

interface ProjectMetadata {
  name: string                     // 工程名称
  author?: string                  // 作者
  description?: string             // 描述
  category: string                 // 分类
  platform: HardwarePlatform       // 硬件平台
  createdAt: string                // 创建时间（ISO8601）
  updatedAt: string                // 更新时间（ISO8601）
}

interface ProjectSecurity {
  encrypted: boolean               // 是否加密
  passwordHash?: string            // 密码哈希（bcrypt）
  deviceBinding?: string           // 设备绑定（预留）
  fileSignature: string            // 文件签名（HMAC-SHA256）
  kekVersion?: string              // KEK版本（预留）
  userEncrypted?: string           // 用户加密的DEK（预留）
  officialEncrypted?: string       // 官方加密的DEK（预留）
}

interface CanvasConfig {
  width: number                    // 画布宽度
  height: number                   // 画布高度
  backgroundColor?: string         // 背景色
}

interface Component {
  id: string                       // 组件ID
  type: string                     // 组件类型
  x: number                        // X坐标
  y: number                        // Y坐标
  width: number                    // 宽度
  height: number                   // 高度
  properties: Record<string, any>  // 属性字典
  dataBindings?: DataBinding[]     // 数据绑定
}

interface DataBinding {
  componentId: string              // 组件ID
  property: string                 // 属性名
  source: string                   // 数据源
}
```

### RecentProject（最近工程对象）

最近工程列表项。

```typescript
interface RecentProject {
  projectId: string                // 工程ID
  name: string                     // 工程名称
  category?: string                // 分类
  filePath: string                 // 文件路径
  lastOpened: string               // 最后打开时间（ISO8601）
  isEncrypted: boolean             // 是否加密
  createdAt: string                // 创建时间（ISO8601）
}
```

---

## 错误码参考

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `INVALID_REQUEST` | 400 | 请求格式错误 |
| `INVALID_PASSWORD` | 401 | 密码错误 |
| `INVALID_SIGNATURE` | 400 | 文件签名验证失败 |
| `INVALID_PROJECT_ID` | 400 | 工程ID无效 |
| `CREATE_FAILED` | 500 | 创建工程失败 |
| `OPEN_FAILED` | 500 | 打开工程失败 |
| `SAVE_FAILED` | 500 | 保存工程失败 |
| `VALIDATE_FAILED` | 500 | 验证密码失败 |
| `GET_RECENT_FAILED` | 500 | 获取最近工程失败 |
| `ADD_RECENT_FAILED` | 500 | 添加最近工程失败 |
| `REMOVE_RECENT_FAILED` | 500 | 删除最近工程失败 |
| `NETWORK_ERROR` | - | 网络请求失败（前端错误） |

---

## 安全机制

### 1. 加密流程

**创建加密工程**:
1. 用户输入密码 → PBKDF2 派生密钥（100,000次迭代）
2. 生成随机盐值 → 存储在加密内容中
3. AES-256-GCM 加密工程内容 → Base64 编码
4. bcrypt 哈希密码（cost 10）→ 存储 `passwordHash`
5. HMAC-SHA256 签名原始内容 → 存储 `fileSignature`
6. 保存 `.pant` 文件

**打开加密工程**:
1. 读取 `.pant` 文件
2. 验证 `passwordHash`（bcrypt）
3. 解密 `encryptedContent`（AES-256-GCM）
4. 验证 `fileSignature`（HMAC-SHA256）
5. 返回工程对象

### 2. 密码强度

**推荐密码要求**:
- 最少 6 个字符（强制）
- 推荐 12+ 个字符
- 包含大小写字母、数字、特殊字符

**强度计算**:
- 弱: ≤ 2 项（长度、大小写、数字、特殊字符）
- 中: 3 项
- 强: 4-5 项

### 3. 文件完整性

- **签名算法**: HMAC-SHA256
- **签名密钥**: `password + projectId`
- **验证时机**: 每次打开加密工程
- **篡改检测**: 签名不匹配时拒绝打开

---

## CORS 配置

后端配置了宽松的 CORS 策略用于开发：

```go
app.Use(cors.New(cors.Config{
  AllowOrigins:     "*",
  AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
  AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
  ExposeHeaders:    "Content-Length",
  AllowCredentials: true,
}))
```

**生产环境建议**: 限制 `AllowOrigins` 为特定域名。

---

## 性能考虑

### 响应时间目标

| 操作 | 目标时间 | 说明 |
|------|---------|------|
| 健康检查 | < 10ms | 简单状态返回 |
| 创建工程 | < 100ms | 不包含文件写入 |
| 打开工程 | < 200ms | 包含解密和验证 |
| 保存工程 | < 100ms | 不包含文件写入 |
| 获取最近工程 | < 50ms | 数据库查询 |

### 优化建议

1. **数据库索引**:
   - `recent_projects.project_id`: 唯一索引
   - `recent_projects.last_opened`: 普通索引（排序优化）

2. **缓存**:
   - 最近工程列表可缓存（TTL: 60s）
   - 密码哈希结果可缓存（验证接口）

3. **分页**:
   - 最近工程列表支持分页（预留功能）
   - 默认返回全部，未来可添加 `?limit=20&offset=0`

---

## 使用示例

### JavaScript/TypeScript

```typescript
import { projectApi } from '@/api/projectApi'

// 创建工程
const createResponse = await projectApi.createProject({
  metadata: {
    name: '测试工程',
    category: '分类1',
    platform: 'HMI型号1',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  security: {
    encrypted: false
  },
  savePath: 'C:\\Projects\\test.pant'
})

if (createResponse.success) {
  console.log('工程ID:', createResponse.data.projectId)
}

// 打开工程
const openResponse = await projectApi.openProject({
  filePath: 'C:\\Projects\\test.pant'
})

if (openResponse.success) {
  const project = openResponse.data.project
  console.log('工程名称:', project.metadata.name)
}

// 获取最近工程
const recentResponse = await projectApi.getRecentProjects()
if (recentResponse.success) {
  const projects = recentResponse.data
  console.log('最近工程数量:', projects.length)
}
```

### cURL

```bash
# 健康检查
curl http://localhost:3000/health

# 创建工程
curl -X POST http://localhost:3000/api/projects/create \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "name": "测试工程",
      "category": "分类1",
      "platform": "HMI型号1",
      "createdAt": "2026-01-21T10:00:00Z",
      "updatedAt": "2026-01-21T10:00:00Z"
    },
    "security": {
      "encrypted": false
    },
    "savePath": "test.pant"
  }'

# 打开工程
curl -X POST http://localhost:3000/api/projects/open \
  -H "Content-Type: application/json" \
  -d '{
    "filePath": "test.pant"
  }'

# 获取最近工程
curl http://localhost:3000/api/projects/recent
```

---

## 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0.0 | 2026-01-21 | 初始版本，实现基础 CRUD 功能 |

---

## 支持

如有问题或建议，请联系开发团队或提交 Issue。
