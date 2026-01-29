# FE-008 代码映射更新

## 更新日期: 2026-01-28

## FE-008-01: 项目初始化 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 项目配置 | `package.json` | - | - | 1-40 | ✅ |
| Vite 配置 | `vite.config.ts` | - | - | 1-50 | ✅ |
| TypeScript 配置 | `tsconfig.json` | - | - | 1-50 | ✅ |
| ESLint 配置 | `.eslintrc.js` | - | - | 1-80 | ✅ |
| Prettier 配置 | `.prettierrc` | - | - | 1-20 | ✅ |
| 环境变量 | `.env.development` | - | - | 1-5 | ✅ |
| 环境变量 | `.env.production` | - | - | 1-5 | ✅ |

## FE-008-02: 基础页面框架 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 应用入口 | `src/App.vue` | App | onLaunch/onShow/onHide | 1-20 | ✅ |
| 主入口 | `src/main.ts` | - | createApp | 1-27 | ✅ |
| 路由配置 | `src/pages.json` | - | - | 1-92 | ✅ |
| 应用配置 | `src/manifest.json` | - | - | 1-100 | ✅ |
| 全局样式 | `src/styles/common.scss` | - | - | 1-139 | ✅ |

## FE-008-03: 核心模块骨架 🔄

### Auth 模块 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 用户 Store | `src/stores/user.store.ts` | useUserStore | login/logout/restoreUser | 1-105 | ✅ |
| 应用 Store | `src/stores/app.store.ts` | useAppStore | setTheme/setLanguage | 1-36 | ✅ |
| 租户 Store | `src/stores/tenant.store.ts` | useTenantStore | - | 1-30 | ✅ |
| 认证 API | `src/api/modules/auth.api.ts` | authApi | login/register/logout | 1-141 | ✅ (Mock) |
| 登录页面 | `src/pages/auth/login/index.vue` | LoginPage | handleLogin | 1-300 | ✅ |

### 其他模块 ⏳

| 模块 | Store | API | Pages | 状态 |
|------|-------|-----|-------|------|
| Device | ⏳ | ⏳ | ⏳ | 0% |
| Workspace | ⏳ | ⏳ | ⏳ | 0% |
| Dashboard | ⏳ | ⏳ | ⏳ | 0% |
| Message | ⏳ | ⏳ | ⏳ | 0% |
| Profile | ⏳ | ⏳ | ⏳ | 0% |

## FE-008-04: 通用组件库 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 自定义导航栏 | `src/components/common/CustomNavBar/index.vue` | CustomNavBar | handleBack | 1-90 | ✅ |
| 页面容器 | `src/components/common/PageContainer/index.vue` | PageContainer | - | 1-30 | ✅ |
| 加载指示器 | `src/components/common/Loading/index.vue` | Loading | - | 1-60 | ✅ |
| 空状态 | `src/components/common/EmptyState/index.vue` | EmptyState | handleAction | 1-70 | ✅ |
| 网络错误 | `src/components/common/NetworkError/index.vue` | NetworkError | handleRetry | 1-50 | ✅ |
| 下拉刷新 | `src/components/common/PullRefresh/index.vue` | PullRefresh | handleRefresh | 1-60 | ✅ |
| 上拉加载 | `src/components/common/LoadMore/index.vue` | LoadMore | handleRetry | 1-120 | ✅ |

## FE-008-05: 工具类和类型定义 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| HTTP 封装 | `src/utils/request.ts` | request | request/get/post | 1-200 | ✅ |
| 存储封装 | `src/utils/storage.ts` | storage | setStorage/getStorage | 1-80 | ✅ |
| 验证工具 | `src/utils/validator.ts` | validator | validatePhone/validateEmail | 1-100 | ✅ |
| 格式化工具 | `src/utils/format.ts` | format | formatDateTime/formatNumber | 1-180 | ✅ |
| 常量定义 | `src/utils/constants.ts` | - | - | 1-50 | ✅ |
| 全局类型 | `src/types/global.d.ts` | - | ApiResponse/UserInfo | 1-78 | ✅ |
| API 类型 | `src/api/types/api.types.ts` | - | LoginParams/LoginResult | 1-50 | ✅ |

## FE-008-06: 第一个页面实现 ✅

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 启动页 | `src/pages/index/index.vue` | IndexPage | onLoad | 1-74 | ✅ |
| 登录页 | `src/pages/auth/login/index.vue` | LoginPage | handleLogin | 1-300 | ✅ |

## FE-008-07: 开发规范 ⏳

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 |
|--------|---------|---------|----------|---------|------|
| 开发规范 | `docs/development-guide.md` | - | - | - | ⏳ |

---

## 统计信息

**总功能点**: 7 个
**已完成**: 5 个 ✅
**进行中**: 1 个 🔄
**未开始**: 1 个 ⏳

**总文件数**: 30+ 个
**已完成文件**: 28 个 ✅
**总代码行数**: 3000+ 行

**当前进度**: 75%

---

## 依赖处理

**FE-007**: 云平台账号系统

**依赖位置**: `src/api/modules/auth.api.ts`

**Mock 实现**: ✅ 已完成

**依赖标记**: ✅ TODO(依赖) 标记完整

**补齐优先级**: P0

---

## 更新历史

| 日期 | 变更内容 | 变更人 |
|------|---------|--------|
| 2026-01-28 | 初始更新,完成 75% | Claude Code |
