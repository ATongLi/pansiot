# IMP-008 代码实现映射表

**更新时间**: 2026-01-28

## 代码映射

| 功能ID | 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|--------|---------|---------|----------|---------|------|
| FE-008-01 | 项目初始化 | package.json | - | - | 1-40 | ✅ |
| FE-008-02 | 目录结构 | src/ | - | - | - | ✅ |
| FE-008-03 | TypeScript 配置 | tsconfig.json | - | - | 1-50 | ✅ |
| FE-008-04 | Vite 配置 | vite.config.ts | - | - | 1-25 | ✅ |
| FE-008-05 | ESLint 配置 | .eslintrc.js | - | - | 1-30 | ✅ |
| FE-008-06 | Prettier 配置 | .prettierrc | - | - | 1-10 | ✅ |
| FE-008-07 | 环境变量 | .env.* | - | - | 1-5 | ✅ |
| FE-008-08 | 全局类型定义 | src/types/global.d.ts | - | - | 1-60 | ✅ |
| FE-008-09 | API 类型定义 | src/api/types/api.types.ts | - | - | 1-30 | ✅ |
| FE-008-10 | HTTP Client 配置 | src/api/client/config.ts | requestConfig | - | 1-10 | ✅ |
| FE-008-11 | 请求拦截器 | src/api/client/interceptors.ts | Interceptor | request() | 18-48 | ✅ |
| FE-008-12 | 响应拦截器 | src/api/client/interceptors.ts | Interceptor | response() | 50-82 | ✅ |
| FE-008-13 | 错误拦截器 | src/api/client/interceptors.ts | Interceptor | error() | 84-107 | ✅ |
| FE-008-14 | Token 刷新 | src/api/client/interceptors.ts | Interceptor | handleTokenExpired() | 109-143 | ✅ |
| FE-008-15 | HTTP Client 封装 | src/api/client/request.ts | Request | request() | 18-37 | ✅ |
| FE-008-16 | GET 请求 | src/api/client/request.ts | Request | get() | 48-50 | ✅ |
| FE-008-17 | POST 请求 | src/api/client/request.ts | Request | post() | 52-54 | ✅ |
| FE-008-18 | PUT 请求 | src/api/client/request.ts | Request | put() | 56-58 | ✅ |
| FE-008-19 | DELETE 请求 | src/api/client/request.ts | Request | delete() | 60-62 | ✅ |
| FE-008-20 | Mock 认证 API | src/api/modules/auth.api.ts | authApi | login() | 28-52 | ✅ |
| FE-008-21 | Mock 注册 API | src/api/modules/auth.api.ts | authApi | register() | 54-78 | ✅ |
| FE-008-22 | App Store | src/stores/app.store.ts | useAppStore | setTheme() | 16-18 | ✅ |
| FE-008-23 | User Store | src/stores/user.store.ts | useUserStore | login() | 36-55 | ✅ |
| FE-008-24 | User Store | src/stores/user.store.ts | useUserStore | logout() | 57-70 | ✅ |
| FE-008-25 | Tenant Store | src/stores/tenant.store.ts | useTenantStore | switchTenant() | 20-27 | ✅ |
| FE-008-26 | 应用主入口 | src/main.ts | createApp | - | 1-25 | ✅ |
| FE-008-27 | 应用根组件 | src/App.vue | - | onLaunch() | 4-6 | ✅ |
| FE-008-28 | 全局样式 | src/styles/common.scss | - | - | 1-150 | ✅ |
| FE-008-29 | 启动页 | src/pages/index/index.vue | - | onLoad() | 20-26 | ✅ |
| FE-008-30 | 工作台页面 | src/pages/tabbar/workspace.vue | - | onLoad() | 18-20 | ✅ |
| FE-008-31 | 设备页面 | src/pages/tabbar/device.vue | - | onLoad() | 17-19 | ✅ |
| FE-008-32 | 看板页面 | src/pages/tabbar/dashboard.vue | - | onLoad() | 17-19 | ✅ |
| FE-008-33 | 消息页面 | src/pages/tabbar/message.vue | - | onLoad() | 17-19 | ✅ |
| FE-008-34 | 我的页面 | src/pages/tabbar/profile.vue | - | onLoad() | 17-19 | ✅ |
| FE-008-35 | 页面配置 | src/pages.json | - | - | 1-100 | ✅ |
| FE-008-36 | 应用配置 | src/manifest.json | - | - | 1-100 | ✅ |
| FE-008-37 | README | README.md | - | - | 1-150 | ✅ |

## 统计信息

- **总文件数**: 30+
- **总代码行数**: 2000+
- **完成功能点**: 37/37 (100%)
- **项目进度**: 基础架构完成

## 项目文件清单

### 配置文件 (10个)
- package.json
- tsconfig.json
- vite.config.ts
- .eslintrc.js
- .prettierrc
- .env.development
- .env.production
- src/pages.json
- src/manifest.json
- README.md

### 核心代码 (20+个)
- src/main.ts
- src/App.vue
- src/types/global.d.ts
- src/api/types/api.types.ts
- src/api/client/config.ts
- src/api/client/interceptors.ts
- src/api/client/request.ts
- src/api/modules/auth.api.ts
- src/stores/app.store.ts
- src/stores/user.store.ts
- src/stores/tenant.store.ts
- src/styles/common.scss
- src/pages/index/index.vue
- src/pages/tabbar/workspace.vue
- src/pages/tabbar/device.vue
- src/pages/tabbar/dashboard.vue
- src/pages/tabbar/message.vue
- src/pages/tabbar/profile.vue

## 依赖忽略记录

| ID | 依赖模块 | 忽略内容 | Mock位置 | 补齐优先级 | 状态 |
|----|---------|---------|---------|-----------|------|
| D-008-01 | FE-007 | 云平台认证 API | src/api/modules/auth.api.ts:28-78 | P0 | ⏳ |

## 实现亮点

1. **完整的 TypeScript 支持**: strict 模式,类型安全
2. **模块化架构**: 清晰的目录结构和模块划分
3. **HTTP 封装**: 完整的请求/响应拦截器,支持 Token 自动刷新
4. **状态管理**: Pinia + 持久化插件
5. **Mock 机制**: 使用依赖忽略机制与后端并行开发
6. **多端配置**: 支持 H5/微信小程序/App
7. **代码规范**: ESLint + Prettier 配置完善

## 下一步工作

1. 补充通用组件 (7个)
2. 完善各模块页面
3. 实现真实 API 对接 (等待 FE-007 完成)
4. 性能优化和测试

## 更新历史

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|---------|--------|
| 2026-01-28 | 1.0 | 初始创建,基础架构完成 | Claude Code |
