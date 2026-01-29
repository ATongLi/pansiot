# IMP-002 集成测试指南

## 概述
本文档提供 Scada 工程管理功能的完整集成测试指南，包括端到端测试流程、验证标准和问题排查。

## 测试环境准备

### 1. 启动后端服务

```bash
cd platforms/scada/backend

# 安装依赖
go mod tidy

# 运行后端服务（开发模式）
go run main.go

# 预期输出：
# ➜ Starting Scada Backend API...
# ➜ Database initialized: ~/.pansiot/pantool.db
# ➜ Server running on http://localhost:3000
# ➜ Health check: http://localhost:3000/health
```

### 2. 启动前端开发服务器

```bash
cd platforms/scada/packages/renderer

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 预期输出：
# ➜ VITE v5.4.2 ready in XXX ms
# ➜ ➜ Local: http://localhost:5173/
# ➜ ➜ Network: use --host to expose
```

### 3. 启动 Electron 桌面应用（可选）

```bash
cd platforms/scada/packages/desktop

# 安装依赖
pnpm install

# 编译 TypeScript
pnpm compile

# 启动 Electron
pnpm dev

# 预期：Electron 窗口打开，加载 http://localhost:5173
```

## 测试场景

### 场景 1：创建新工程（未加密）

**前置条件**:
- 后端服务运行中
- 前端应用已加载

**测试步骤**:

1. **打开新建工程对话框**
   - 点击"新建工程"按钮
   - 验证：对话框显示，包含所有表单字段

2. **填写基本信息**
   - 工程名称: `测试工程001`
   - 工程作者: `测试用户`
   - 工程描述: `这是一个测试工程`
   - 工程分类: `分类1`
   - 硬件平台: `HMI型号1`

3. **选择保存位置**
   - 点击"浏览..."按钮
   - **浏览器环境**: 验证控制台输出 `[Mock] selectSavePath`
   - **Electron环境**: 验证文件对话框打开
   - 选择保存路径（或确认Mock路径）

4. **提交表单**
   - 点击"确定"按钮
   - 验证：显示"创建中..."状态
   - 验证：控制台无错误信息

5. **验证结果**
   - 后端日志显示: `POST /api/projects/create 200 OK`
   - 数据库记录: 检查 `~/.pansiot/pantool.db` 的 `recent_projects` 表
   - 文件创建: 验证 `.pant` 文件在指定路径创建

**预期结果**:
- ✅ 工程创建成功
- ✅ 文件保存到指定位置
- ✅ 最近工程列表更新
- ✅ 无错误提示

---

### 场景 2：创建加密工程

**测试步骤**:

1. **打开新建工程对话框**
2. **填写基本信息**（同场景1）
3. **启用加密**
   - 勾选"启用工程加密"复选框
   - 验证：密码字段显示

4. **设置密码**
   - 密码: `Test@123456`
   - 确认密码: `Test@123456`
   - 验证：密码强度指示器显示"强"

5. **提交表单**
   - 点击"确定"按钮

6. **验证加密**
   ```bash
   # 检查生成的 .pant 文件
   cat "测试工程001.pant" | grep encryptedContent
   # 应该看到 Base64 编码的加密内容
   ```

**预期结果**:
- ✅ 工程文件包含 `encryptedContent` 字段
- ✅ `passwordHash` 字段存在（bcrypt hash）
- ✅ `fileSignature` 字段存在（HMAC-SHA256）
- ✅ 原始内容不在文件中（只有加密后的内容）

---

### 场景 3：打开未加密工程

**测试步骤**:

1. **打开工程对话框**
   - 点击"打开工程"按钮
   - 或从最近工程列表点击工程卡片

2. **选择工程文件**
   - **浏览器环境**: Mock路径自动填入
   - **Electron环境**: 文件对话框选择 `.pant` 文件

3. **打开工程**
   - 点击"确定"按钮

4. **验证加载**
   - 后端日志: `POST /api/projects/open 200 OK`
   - 前端状态: `projectStore.currentProject` 不为 null
   - UI显示: 工程信息正确显示

**预期结果**:
- ✅ 工程成功加载
- ✅ 元数据显示正确
- ✅ 组件列表显示
- ✅ 最近工程列表更新（lastOpened 时间更新）

---

### 场景 4：打开加密工程（正确密码）

**测试步骤**:

1. **创建加密工程**（使用场景2）
2. **重新打开工程**
   - 选择该加密工程文件
   - 验证：密码对话框显示

3. **输入正确密码**
   - 输入: `Test@123456`
   - 点击"确定"

4. **验证解密**
   - 后端日志显示解密成功
   - 工程内容正确显示
   - 无错误提示

**预期结果**:
- ✅ 密码验证通过
- ✅ 内容成功解密
- ✅ 签名验证通过
- ✅ 工程正常加载

---

### 场景 5：打开加密工程（错误密码）

**测试步骤**:

1. **尝试打开加密工程**
2. **输入错误密码**
   - 输入: `WrongPassword`
   - 点击"确定"

3. **验证错误处理**
   - 后端返回: `401 Unauthorized`
   - 错误信息: `"error": "INVALID_PASSWORD"`
   - 前端显示: "密码错误"

**预期结果**:
- ✅ 密码错误提示
- ✅ 工程未打开
- ✅ 可以重新输入密码

---

### 场景 6：最近工程列表（筛选、搜索、排序）

**测试步骤**:

1. **加载最近工程列表**
   - 页面加载时自动调用 `GET /api/projects/recent`
   - 验证：列表显示

2. **分类筛选**
   - 点击"分类1"标签
   - 验证：只显示该分类的工程
   - 数量徽章正确

3. **搜索**
   - 在搜索框输入: `测试`
   - 验证：只显示名称包含"测试"的工程
   - 搜索有100ms防抖

4. **排序**
   - 点击"最后打开"列头
   - 验证：按 lastOpenedDate 降序排列
   - 再次点击：切换为升序

5. **组合操作**
   - 选择分类 + 输入搜索词
   - 验证：结果同时满足两个条件

**预期结果**:
- ✅ 分类筛选正常
- ✅ 搜索功能正常
- ✅ 排序功能正常
- ✅ 组合操作正常
- ✅ 虚拟滚动流畅（1000+ 工程）

---

### 场景 7：错误处理

#### 7.1 文件不存在

```bash
# 尝试打开不存在的文件
curl -X POST http://localhost:3000/api/projects/open \
  -H "Content-Type: application/json" \
  -d '{"filePath": "/nonexistent/file.pant"}'

# 预期: 500 Internal Server Error
# message: "读取文件失败"
```

#### 7.2 文件签名篡改

```bash
# 创建工程后，手动修改文件
echo "modified" >> test-project.pant

# 尝试打开
curl -X POST http://localhost:3000/api/projects/open \
  -H "Content-Type: application/json" \
  -d '{"filePath": "test-project.pant"}'

# 预期: 400 Bad Request
# error: "INVALID_SIGNATURE"
# message: "工程文件签名验证失败，文件可能已被篡改"
```

#### 7.3 网络错误

**测试步骤**:
1. 停止后端服务
2. 尝试创建/打开工程
3. 验证前端错误提示: "网络请求失败"

---

## 性能测试

### 虚拟滚动性能

```typescript
// 在浏览器控制台执行
// 创建1000个测试工程
const tests = []
for (let i = 0; i < 1000; i++) {
  tests.push({
    projectId: `test-${i}`,
    name: `测试工程${i}`,
    lastOpenedDate: new Date(),
    filePath: `/test/${i}.pant`,
    isEncrypted: false,
    createdAt: new Date().toISOString()
  })
}

// 模拟加载
recentProjectsStore.projects = tests

// 测试滚动性能
console.time('scroll')
// 快速滚动列表
// 预期: 无卡顿，保持 60 FPS
console.timeEnd('scroll')
```

### API 响应时间

```bash
# 测试创建工程响应时间
time curl -X POST http://localhost:3000/api/projects/create \
  -H "Content-Type: application/json" \
  -d @test-project.json

# 预期: < 100ms（不包含文件写入）
```

---

## 验证标准

### ✅ 完整通过的标志

- [ ] 所有8个场景测试通过
- [ ] 无控制台错误
- [ ] 数据库记录正确
- [ ] 文件系统操作正确
- [ ] UI 响应流畅（60 FPS）
- [ ] 虚拟滚动性能良好
- [ ] 错误处理健壮

### ⚠️ 部分通过的标志

- [ ] 基本功能正常，但部分高级特性未实现
- [ ] 性能可接受但不够优化
- [ ] 错误处理基本完善，但边界情况未覆盖

### ❌ 未通过的标志

- [ ] 核心功能缺失
- [ ] 频繁崩溃或严重错误
- [ ] 性能不可接受（< 30 FPS）
- [ ] 数据丢失或损坏

---

## 常见问题排查

### 问题 1: 后端服务无法启动

**症状**: `panic: failed to initialize database`

**解决**:
```bash
# 检查数据库权限
ls -la ~/.pansiot/pantool.db

# 删除旧数据库重新初始化
rm ~/.pansiot/pantool.db
cd platforms/scada/backend && go run main.go
```

### 问题 2: 前端无法连接后端

**症状**: `Network request failed`

**解决**:
```bash
# 检查后端是否运行
curl http://localhost:3000/health

# 检查 Vite 代理配置
# vite.config.ts 应该有:
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:3000',
      changeOrigin: true,
    }
  }
}
```

### 问题 3: Electron 文件对话框不工作

**症状**: 点击"浏览..."按钮无反应

**解决**:
```bash
# 检查 preload.js 是否正确编译
ls platforms/scada/packages/desktop/preload.js

# 检查 electron-main.ts 是否正确加载
# 应该看到日志: "Starting Scada Desktop Application..."

# 检查 contextBridge 配置
# 在渲染进程控制台执行:
console.log(window.electronAPI) // 应该不是 undefined
```

---

## 测试清单

### 功能测试

- [ ] 创建未加密工程
- [ ] 创建加密工程
- [ ] 打开未加密工程
- [ ] 打开加密工程（正确密码）
- [ ] 打开加密工程（错误密码）
- [ ] 保存工程
- [ ] 最近工程列表加载
- [ ] 分类筛选
- [ ] 搜索功能
- [ ] 排序功能
- [ ] 虚拟滚动（大列表）
- [ ] 工程删除

### 集成测试

- [ ] 前端 ↔ 后端 API 通信
- [ ] Electron 文件对话框
- [ ] Electron 窗口控制
- [ ] 数据库 CRUD 操作
- [ ] 文件系统读写操作

### 错误处理测试

- [ ] 文件不存在
- [ ] 密码错误
- [ ] 签名篡改
- [ ] 网络错误
- [ ] 数据库连接失败

### 性能测试

- [ ] 虚拟滚动性能（1000+ 工程）
- [ ] 搜索防抖
- [ ] API 响应时间
- [ ] MobX 响应式更新
- [ ] 内存使用情况

---

## 测试报告模板

```markdown
# IMP-002 集成测试报告

**测试日期**: 2026-01-XX
**测试人员**: XXX
**环境**:
- OS: Windows 10 / macOS 14 / Ubuntu 22.04
- Node.js: v20.x.x
- Go: 1.21.x
- Electron: 28.x.x

## 测试结果汇总

- 通过场景: X / 8
- 失败场景: Y / 8
- 阻塞问题: Z 个

## 详细测试结果

### 场景 1: 创建新工程（未加密）
- 状态: ✅ 通过 / ⚠️ 部分通过 / ❌ 失败
- 备注: ...

### 场景 2-8: ...

## 问题列表

1. **问题描述**
   - 严重程度: 高 / 中 / 低
   - 复现步骤: ...
   - 预期结果: ...
   - 实际结果: ...

## 性能数据

- 虚拟滚动 FPS: XX
- API 平均响应时间: XX ms
- 内存占用: XX MB

## 结论

- [ ] 可以发布
- [ ] 需要修复问题
- [ ] 需要进一步优化
```
