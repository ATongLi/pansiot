# PansIot Cloud Platform - Backend

云平台账号系统后端服务，基于 Golang + Gin + PostgreSQL + Redis 实现。

## 技术栈

- **语言**: Golang 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7.0+
- **认证**: JWT
- **日志**: zap
- **配置**: Viper
- **API文档**: Swagger

## 项目结构

```
platforms/cloud/backend/
├── cmd/
│   └── server/
│       └── main.go           # 主入口文件
├── internal/
│   ├── auth/                  # 认证模块
│   ├── tenant/                # 租户管理模块
│   ├── user/                  # 用户管理模块
│   ├── role/                  # 角色管理模块
│   ├── permission/            # 权限管理模块
│   ├── quota/                 # 配额管理模块
│   ├── audit/                 # 审计日志模块
│   ├── middleware/            # 中间件
│   ├── models/                # 数据模型
│   ├── config/                # 配置管理
│   └── api/                   # API处理器
├── pkg/
│   ├── database/              # 数据库连接
│   ├── logger/                # 日志封装
│   └── response/              # 统一响应
├── migrations/                # 数据库迁移
├── configs/                   # 配置文件
│   └── config.yaml           # 默认配置
├── scripts/                   # 脚本文件
│   └── init.sql              # 数据库初始化
├── logs/                      # 日志文件
├── go.mod                     # Go模块定义
├── go.sum                     # Go依赖锁定
├── Makefile                   # 构建脚本
├── Dockerfile                 # Docker镜像
└── docker-compose.yml         # Docker编排
```

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 14+
- Redis 7.0+
- Docker (可选)

### 本地开发

1. **克隆项目**
   ```bash
   cd platforms/cloud/backend
   ```

2. **安装依赖**
   ```bash
   make deps
   ```

3. **配置数据库**
   ```bash
   # 使用Docker启动PostgreSQL和Redis
   docker-compose up -d postgres redis
   ```

4. **初始化数据库**
   ```bash
   psql -h localhost -U postgres -d pansiot_cloud -f scripts/init.sql
   ```

5. **运行服务**
   ```bash
   make run
   ```

服务将在 `http://localhost:8080` 启动。

### 使用Docker

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down
```

## API文档

启动服务后，访问 Swagger 文档：
```
http://localhost:8080/swagger/index.html
```

## 主要功能模块

### 1. 认证模块 (FE-007-02)
- 用户注册（新企业/加入已有企业）
- 用户登录/登出
- Token刷新
- 密码重置

### 2. 租户管理 (FE-007-01)
- 组织CRUD
- 组织树查询
- 租户类型升级/降级

### 3. 用户管理 (FE-007-02)
- 用户CRUD
- 角色分配
- 批量操作

### 4. 角色权限管理 (FE-007-03)
- 角色CRUD
- 权限配置
- 权限验证

### 5. 功能模块与配额 (FE-007-08)
- 功能模块管理
- 配额分配
- 配额统计

### 6. 审计日志 (FE-007-09)
- 操作日志记录
- 日志查询
- 日志导出

## 核心特性

### 多租户架构
- 三层租户模型：平台超管 → 集成商 → 下游客户
- 双字段数据隔离：`tenant_id` + `managed_tenant_id`
- 租户级RBAC权限模型

### 数据隔离
- 下游客户：只能看到自己的数据
- 集成商：可以看到所有下游数据
- 跨租户严格隔离

### 动态界面
- 根据租户类型动态显示菜单
- 根据权限动态显示按钮
- 界面简化规则

## 开发规范

### 代码风格
- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化代码
- 必须编写单元测试

### 提交规范
```
feat: 添加新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 重构代码
test: 测试相关
chore: 构建/工具链相关
```

## 测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage
```

## 构建部署

```bash
# 构建二进制文件
make build

# 使用Docker构建
make docker-build
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | 服务端口 | 8080 |
| `GIN_MODE` | 运行模式 | debug |
| `DB_HOST` | 数据库地址 | localhost |
| `DB_PORT` | 数据库端口 | 5432 |
| `DB_USER` | 数据库用户 | postgres |
| `DB_PASSWORD` | 数据库密码 | postgres |
| `DB_NAME` | 数据库名称 | pansiot_cloud |
| `REDIS_HOST` | Redis地址 | localhost |
| `REDIS_PORT` | Redis端口 | 6379 |
| `JWT_SECRET` | JWT密钥 | - |

## 常见问题

### 1. 数据库连接失败
检查 PostgreSQL 是否正常运行：
```bash
docker ps | grep postgres
```

### 2. Redis连接失败
检查 Redis 是否正常运行：
```bash
docker ps | grep redis
```

### 3. 端口被占用
修改 `configs/config.yaml` 中的端口配置。

## License

Copyright © 2026 PansIot
