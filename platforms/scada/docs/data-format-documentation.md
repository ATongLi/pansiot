# .pant 工程文件格式规范

## 概述

`.pant` 文件是 PanTools Scada 工程的存储格式，使用 JSON 格式存储工程的所有配置和数据。支持可选的 AES-256-GCM 加密。

**文件扩展名**: `.pant`
**MIME 类型**: `application/json`
**字符编码**: `UTF-8`
**版本**: 1.0.0

---

## 文件结构

### 未加密工程文件

```json
{
  "version": "1.0.0",
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "name": "工厂监控画面01",
    "author": "张三",
    "description": "用于监控生产线状态的HMI画面",
    "category": "分类1",
    "platform": "HMI型号1",
    "createdAt": "2026-01-21T10:00:00Z",
    "updatedAt": "2026-01-21T10:00:00Z"
  },
  "security": {
    "encrypted": false,
    "passwordHash": "",
    "deviceBinding": "",
    "fileSignature": "abc123def456..."
  },
  "canvas": {
    "width": 1920,
    "height": 1080,
    "backgroundColor": "#1a1a1a"
  },
  "components": [
    {
      "id": "comp_001",
      "type": "gauge",
      "x": 100,
      "y": 100,
      "width": 200,
      "height": 200,
      "properties": {
        "minValue": 0,
        "maxValue": 100,
        "value": 75.5,
        "unit": "°C",
        "title": "温度"
      },
      "dataBindings": [
        {
          "componentId": "comp_001",
          "property": "value",
          "source": "DV-PLC001-TEMP01"
        }
      ]
    }
  ]
}
```

### 加密工程文件

```json
{
  "version": "1.0.0",
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "name": "机密工程",
    "author": "李四",
    "description": "包含敏感数据",
    "category": "分类2",
    "platform": "HMI型号2",
    "createdAt": "2026-01-21T10:00:00Z",
    "updatedAt": "2026-01-21T10:00:00Z"
  },
  "security": {
    "encrypted": true,
    "passwordHash": "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
    "deviceBinding": "",
    "fileSignature": "xyz789abc012...",
    "kekVersion": ""
  },
  "canvas": {
    "width": 1920,
    "height": 1080,
    "backgroundColor": "#1a1a1a"
  },
  "components": [],
  "encryptedContent": "U2FsdGVkX1+xxx...（Base64编码的加密内容）"
}
```

**注意**: 加密工程文件中，`components` 数组为空，实际内容存储在 `encryptedContent` 字段中。

---

## 字段详解

### 顶层字段

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `version` | string | 是 | 文件格式版本，当前 "1.0.0" |
| `projectId` | string | 是 | 工程唯一ID（UUID v4格式） |
| `metadata` | object | 是 | 工程元数据 |
| `security` | object | 是 | 安全配置 |
| `canvas` | object | 是 | 画布配置 |
| `components` | array | 是 | 组件列表（加密工程为空） |
| `encryptedContent` | string | 否 | 加密内容（仅加密工程） |

---

### metadata（元数据）

工程的基本信息。

| 字段 | 类型 | 必填 | 说明 | 限制 |
|------|------|------|------|------|
| `name` | string | 是 | 工程名称 | 1-50字符 |
| `author` | string | 否 | 作者名称 | 最多30字符 |
| `description` | string | 否 | 工程描述 | 最多500字符 |
| `category` | string | 是 | 工程分类 | 预定义或自定义 |
| `platform` | string | 是 | 硬件平台 | 预定义枚举 |
| `createdAt` | string | 是 | 创建时间 | ISO8601格式 |
| `updatedAt` | string | 是 | 更新时间 | ISO8601格式 |

**工程分类**:
- `分类1`
- `分类2`
- 自定义分类（用户输入）

**硬件平台**:
- `HMI型号1`
- `HMI型号2`
- `网关型号1`

---

### security（安全配置）

工程的安全和加密设置。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `encrypted` | boolean | 是 | 是否启用加密 |
| `passwordHash` | string | 条件 | 密码哈希（加密时必填） |
| `deviceBinding` | string | 否 | 设备绑定标识（预留） |
| `fileSignature` | string | 是 | 文件签名（HMAC-SHA256） |
| `kekVersion` | string | 否 | KEK版本（预留，密码恢复功能） |
| `userEncrypted` | string | 否 | 用户加密的DEK（预留） |
| `officialEncrypted` | string | 否 | 官方加密的DEK（预留） |

#### passwordHash

**算法**: bcrypt
**Cost因子**: 10
**格式**: `$2a$10$...`（60字符）

**示例**:
```json
"passwordHash": "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
```

#### fileSignature

**算法**: HMAC-SHA256
**签名内容**: 原始工程JSON（不含 encryptedContent 字段）
**签名密钥**: `password + projectId`
**格式**: 十六进制字符串（64字符）

**验证**:
```go
// 伪代码
signature = HMAC-SHA256(jsonData, password + projectId)
if (signature != fileSignature) {
  return error("INVALID_SIGNATURE")
}
```

---

### canvas（画布配置）

工程画布的尺寸和样式。

| 字段 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| `width` | number | 是 | 画布宽度（像素） | 1920 |
| `height` | number | 是 | 画布高度（像素） | 1080 |
| `backgroundColor` | string | 否 | 背景色（CSS颜色） | "#ffffff" |

**常用分辨率**:
- 1920x1080 (Full HD)
- 1280x720 (HD)
- 1024x768 (XGA)
- 800x600 (SVGA)

---

### components（组件列表）

工程中的所有组件。

**未加密工程**: 完整的组件数组
**加密工程**: 空数组 `[]`，实际内容在 `encryptedContent` 中

#### 组件对象结构

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 组件唯一ID |
| `type` | string | 是 | 组件类型 |
| `x` | number | 是 | X坐标（像素） |
| `y` | number | 是 | Y坐标（像素） |
| `width` | number | 是 | 宽度（像素） |
| `height` | number | 是 | 高度（像素） |
| `properties` | object | 是 | 组件属性（键值对） |
| `dataBindings` | array | 否 | 数据绑定列表 |

#### 组件类型

**基础组件**:
- `rectangle`: 矩形
- `circle`: 圆形
- `line`: 线条
- `text`: 文本
- `image`: 图片

**工业组件**:
- `gauge`: 仪表盘
- `meter`: 仪表
- `indicator`: 指示灯
- `button`: 按钮
- `switch`: 开关
- `slider`: 滑块
- `chart`: 图表
- `table`: 表格

**容器组件**:
- `group`: 组
- `panel`: 面板
- `tab`: 标签页

#### 组件示例

**仪表盘组件**:
```json
{
  "id": "gauge_001",
  "type": "gauge",
  "x": 100,
  "y": 100,
  "width": 200,
  "height": 200,
  "properties": {
    "minValue": 0,
    "maxValue": 100,
    "value": 75.5,
    "unit": "°C",
    "title": "温度",
    "color": "#ff0000",
    "showScale": true
  },
  "dataBindings": [
    {
      "componentId": "gauge_001",
      "property": "value",
      "source": "DV-PLC001-TEMP01"
    }
  ]
}
```

**按钮组件**:
```json
{
  "id": "btn_001",
  "type": "button",
  "x": 500,
  "y": 300,
  "width": 120,
  "height": 40,
  "properties": {
    "text": "启动",
    "backgroundColor": "#4caf50",
    "fontSize": 16,
    "borderRadius": 4,
    "action": "write",
    "target": "MV-PLC001-START",
    "value": true
  },
  "dataBindings": []
}
```

---

### dataBindings（数据绑定）

组件属性与数据源的绑定关系。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `componentId` | string | 是 | 组件ID |
| `property` | string | 是 | 组件属性名 |
| `source` | string | 是 | 数据源标识符 |

**数据源格式**:
- `DV-{设备}-{变量}`: 数据变量
- `MV-{设备}-{变量}`: 内存变量
- `api:/api/path`: API端点

---

## 加密格式

### 加密流程

1. **原始数据**: 工程JSON对象（不含 `encryptedContent` 字段）
2. **序列化**: JSON序列化（缩进格式）
3. **密钥派生**: PBKDF2-HMAC-SHA256
   - 密码: 用户输入
   - 盐值: 16字节随机
   - 迭代: 100,000次
   - 输出: 32字节密钥
4. **加密**: AES-256-GCM
   - Nonce: 12字节随机
   - 认证标签: 16字节
5. **编码**: Base64
6. **格式**: `Salt(16) + Nonce(12) + Ciphertext + Tag(16)`

### encryptedContent 字段

**Base64 编码结构**:
```
[salt(16 bytes)][nonce(12 bytes)][ciphertext][authTag(16 bytes)]
```

**总长度**: 取决于原文大小，约 `原文长度 * 4/3 + 44` 字符

**解密伪代码**:
```go
// 1. Base64 解码
data := base64.Decode(encryptedContent)

// 2. 提取各部分
salt := data[0:16]
nonce := data[16:28]
ciphertext := data[28:len(data)-16]
tag := data[len(data)-16:]

// 3. 派生密钥
key := PBKDF2(password, salt, 100000)

// 4. AES-GCM 解密
plaintext := AES-GCM-Decrypt(key, nonce, ciphertext, tag)
```

---

## 版本控制

### 版本号格式

`主版本.次版本.修订版本`

**主版本**: 重大结构变更，不向后兼容
**次版本**: 新增字段，向后兼容
**修订版本**: Bug修复和小改进

### 当前版本: 1.0.0

**特性**:
- 基础工程结构
- 单加密机制
- 组件系统
- 数据绑定
- 文件签名

### 未来计划

**1.1.0**:
- 双重加密（DEK + KEK）
- 官方密码恢复
- 组件模板

**2.0.0**:
- 多页面支持
- 脚本引擎集成
- 历史记录

---

## 验证和校验

### 文件完整性

**方法**: HMAC-SHA256 签名

**验证步骤**:
1. 读取文件
2. 提取 `fileSignature` 字段
3. 移除 `encryptedContent` 字段（如果存在）
4. 重新计算签名: `HMAC-SHA256(jsonData, password + projectId)`
5. 对比签名

### 密码验证

**方法**: bcrypt 对比

**验证步骤**:
1. 读取 `passwordHash` 字段
2. 使用 bcrypt 验证: `bcrypt.CompareHashAndPassword(hash, password)`
3. 返回验证结果

### 数据验证

**必填字段检查**:
- `version`, `projectId`, `metadata`, `security`, `canvas`, `components`

**格式验证**:
- `projectId`: UUID v4 格式
- `createdAt`, `updatedAt`: ISO8601 格式
- `passwordHash`: bcrypt 格式（60字符，以 `$2a$10$` 开头）
- `fileSignature`: 64字符十六进制字符串

**业务验证**:
- `metadata.name`: 长度 1-50
- `canvas.width`, `canvas.height`: 正整数
- 组件坐标非负
- 组件尺寸正数

---

## 兼容性

### 向后兼容

**原则**: 新版本可以打开旧版本文件

**实现**:
- 旧版本文件缺少新字段时使用默认值
- 读取时忽略未知字段
- 保存时升级到当前版本

**示例**:
```json
// 1.0.0 文件
{
  "version": "1.0.0",
  "metadata": { ... }
}

// 1.1.0 读取时添加新字段
{
  "version": "1.0.0",
  "metadata": { ... },
  "kekVersion": ""  // 新增字段，默认空字符串
}
```

### 向前兼容

**原则**: 旧版本可以打开新版本文件（尽可能）

**实现**:
- 新字段使用可选
- 核心结构不变
- 提供降级转换

---

## 性能考虑

### 文件大小

**未加密工程**:
- 小工程（<10组件）: ~5 KB
- 中等工程（10-100组件）: ~50 KB
- 大工程（>100组件）: ~500 KB

**加密工程**:
- 增加 ~33% 大小（Base64编码）
- 认证标签: 16字节

### 加载时间

**未加密**:
- 小工程: < 10ms
- 中等工程: < 50ms
- 大工程: < 200ms

**加密**:
- 小工程: < 50ms
- 中等工程: < 200ms
- 大工程: < 1000ms

### 优化建议

1. **延迟加载**: 组件属性按需加载
2. **缓存**: 解密后的内容缓存在内存
3. **增量保存**: 只保存修改的部分
4. **压缩**: 大文件使用 GZIP 压缩

---

## 安全建议

### 密码管理

1. **密码强度**: 推荐 12+ 字符，混合大小写、数字、符号
2. **密码存储**: 不在日志或错误信息中泄露密码
3. **密码传输**: 仅在本地使用，不网络传输

### 文件保护

1. **权限**: .pant 文件设置适当的文件系统权限
2. **备份**: 加密工程的备份同样需要保护
3. **销毁**: 删除时安全擦除（覆盖数据）

### 密钥管理

1. **盐值**: 每个文件使用唯一随机盐值
2. **Nonce**: 每次加密使用唯一随机 nonce
3. **密钥派生**: 使用足够的迭代次数（100,000+）

---

## 工具支持

### 官方工具

**PanTools Scada**:
- 创建工程
- 打开工程
- 保存工程
- 导出/导入

**密码恢复工具** (Phase 7):
- 官方移除密码
- 批量处理
- 所有权验证

### 第三方工具

**JSON 验证器**:
- [JSONLint](https://jsonlint.com/)
- VS Code 插件

**加密测试**:
- OpenSSL: `openssl enc -aes-256-gcm`
- Python: `cryptography` 库

---

## 示例文件

### 完整示例

**文件**: `factory-monitor.pant`

```json
{
  "version": "1.0.0",
  "projectId": "550e8400-e29b-41d4-a716-446655440000",
  "metadata": {
    "name": "工厂生产线监控",
    "author": "工程师001",
    "description": "监控3条生产线的温度、压力和运行状态",
    "category": "分类1",
    "platform": "HMI型号1",
    "createdAt": "2026-01-21T08:00:00Z",
    "updatedAt": "2026-01-21T10:30:00Z"
  },
  "security": {
    "encrypted": false,
    "passwordHash": "",
    "deviceBinding": "",
    "fileSignature": "a1b2c3d4e5f6..."
  },
  "canvas": {
    "width": 1920,
    "height": 1080,
    "backgroundColor": "#1a1a1a"
  },
  "components": [
    {
      "id": "gauge_temp_01",
      "type": "gauge",
      "x": 50,
      "y": 50,
      "width": 200,
      "height": 200,
      "properties": {
        "minValue": 0,
        "maxValue": 150,
        "value": 85.5,
        "unit": "°C",
        "title": "生产线1温度"
      },
      "dataBindings": [
        {
          "componentId": "gauge_temp_01",
          "property": "value",
          "source": "DV-PLC001-TEMP01"
        }
      ]
    },
    {
      "id": "gauge_press_01",
      "type": "gauge",
      "x": 300,
      "y": 50,
      "width": 200,
      "height": 200,
      "properties": {
        "minValue": 0,
        "maxValue": 10,
        "value": 5.2,
        "unit": "MPa",
        "title": "生产线1压力"
      },
      "dataBindings": [
        {
          "componentId": "gauge_press_01",
          "property": "value",
          "source": "DV-PLC001-PRESS01"
        }
      ]
    },
    {
      "id": "indicator_status_01",
      "type": "indicator",
      "x": 550,
      "y": 100,
      "width": 50,
      "height": 50,
      "properties": {
        "state": "running",
        "color": "#4caf50"
      },
      "dataBindings": [
        {
          "componentId": "indicator_status_01",
          "property": "state",
          "source": "DV-PLC001-STATUS"
        }
      ]
    }
  ]
}
```

---

## 故障排查

### 常见问题

**问题**: 文件无法打开

**排查步骤**:
1. 检查文件签名: `fileSignature` 是否匹配
2. 检查版本号: 是否支持
3. 检查必填字段: 是否完整
4. 检查JSON格式: 是否有效

**问题**: 密码验证失败

**排查步骤**:
1. 确认密码正确
2. 检查 `passwordHash` 格式
3. 验证 bcrypt cost 因子
4. 检查文件是否被篡改

**问题**: 解密失败

**排查步骤**:
1. 检查 `encryptedContent` Base64 格式
2. 验证盐值和 nonce 长度
3. 检查 AES-GCM 认证标签
4. 确认密钥派生参数

---

## 更新日志

### v1.0.0 (2026-01-21)

**新增**:
- 初始文件格式定义
- 基础组件系统
- 单加密机制
- 文件签名验证
- 数据绑定

---

## 参考资料

- [JSON 规范](https://www.json.org/)
- [AES-GCM 加密](https://tools.ietf.org/html/rfc5116)
- [HMAC-SHA256](https://tools.ietf.org/html/rfc2104)
- [PBKDF2](https://tools.ietf.org/html/rfc2898)
- [bcrypt](https://github.com/patrickfav/bcrypt)
- [UUID v4](https://tools.ietf.org/html/rfc4122)

---

## 联系方式

如有格式问题或建议，请联系开发团队。
