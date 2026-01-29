# 移动端应用系统架构设计

## 架构概述

移动端应用采用 **分层架构** 和 **模块化设计**,确保代码可维护、可扩展、可测试。

## 整体架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Presentation Layer                         │
│                        (展示层 - Vue Components)                    │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │    Pages     │  │  Components  │  │     App      │              │
│  │   (页面层)    │  │   (组件层)    │  │  (应用入口)   │              │
│  │              │  │              │  │              │              │
│  │ - index/     │  │ - common/    │  │ - App.vue    │              │
│  │ - auth/      │  │ - business/  │  │ - main.ts    │              │
│  │ - device/    │  │              │  │              │              │
│  │ - workspace/ │  │              │  │              │              │
│  │ - dashboard/ │  │              │  │              │              │
│  │ - message/   │  │              │  │              │              │
│  │ - profile/   │  │              │  │              │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└────────────────────────────┬───────────────────────────────────────┘
                             │
                             ↓ (Composables)
┌─────────────────────────────────────────────────────────────────────┐
│                         Business Layer                             │
│                      (业务层 - Vue Composables)                    │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   Stores     │  │  Composables │  │   Services   │              │
│  │  (状态管理)   │  │ (组合式函数)  │  │  (业务服务)   │              │
│  │              │  │              │  │              │              │
│  │ - auth.ts    │  │ - useAuth    │  │ - AuthService│              │
│  │ - device.ts  │  │ - useDevice  │  │ - DeviceSvc  │              │
│  │ - workspace.│  │ - useReq..   │  │ - WorkspaceS │              │
│  │ - dashboard.│  │ - usePagin.. │  │              │              │
│  │ - message.ts│  │              │  │              │              │
│  │ - profile.ts│  │              │  │              │              │
│  │ - app.ts     │  │              │  │              │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│                            ↓                                        │
│                      (Pinia State Management)                      │
└────────────────────────────┬───────────────────────────────────────┘
                             │
                             ↓ (API Calls)
┌─────────────────────────────────────────────────────────────────────┐
│                           Data Layer                                │
│                      (数据层 - Utils & Types)                       │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │     API      │  │    Utils     │  │    Types     │              │
│  │  (API封装)   │  │  (工具类)     │  │  (类型定义)   │              │
│  │              │  │              │  │              │              │
│  │ - auth.ts    │  │ - request.ts │  │ - auth.d.ts  │              │
│  │ - device.ts  │  │ - storage.ts │  │ - device.d..│              │
│  │ - workspace.│  │ - validator.│  │ - workspace.│              │
│  │ - dashboard.│  │ - format.ts  │  │ - dashboard.│              │
│  │ - message.ts│  │ - date.ts    │  │ - message.d.│              │
│  │ - profile.ts│  │ - constants.│  │ - profile.d.│              │
│  │              │  │              │  │ - common.d.ts│              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└────────────────────────────┬───────────────────────────────────────┘
                             │
                             ↓ (HTTP Requests)
┌─────────────────────────────────────────────────────────────────────┐
│                    External Services Layer                         │
│                       (外部服务层)                                  │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ Cloud API    │  │ Local Storage│  │   Uni API    │              │
│  │ (云平台接口)  │  │  (本地存储)   │  │ (UniApp接口) │              │
│  │              │  │              │  │              │              │
│  │ - Auth API   │  │ - Token      │  │ - navigateTo │              │
│  │ - Device API │  │ - UserInfo   │  │ - request    │              │
│  │ - Message API│  │ - Cache      │  │ - storage    │              │
│  │              │  │              │  │ - chooseImage│              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
```

## 分层职责

### Presentation Layer (展示层)

**职责**:
- UI 渲染和用户交互
- 页面路由和导航
- 组件复用和组合

**技术栈**:
- Vue 3 (Composition API)
- UniApp 框架
- UniUI 组件库

**主要文件**:
- `src/pages/` - 页面组件
- `src/components/` - 通用组件和业务组件
- `src/App.vue` - 应用入口组件

### Business Layer (业务层)

**职责**:
- 业务逻辑封装
- 状态管理
- 数据转换和处理

**技术栈**:
- Pinia (状态管理)
- Vue 3 Composables (组合式函数)

**主要文件**:
- `src/stores/` - Pinia Store 状态管理
- `src/composables/` - 组合式函数 (useAuth, useDevice, etc.)
- `src/services/` - 业务服务 (AuthService, DeviceService, etc.)

### Data Layer (数据层)

**职责**:
- API 请求封装
- 数据持久化
- 工具函数和类型定义

**技术栈**:
- TypeScript (类型安全)
- uni.request (HTTP 请求)
- uni.storage (本地存储)

**主要文件**:
- `src/api/` - API 接口封装
- `src/utils/` - 工具函数
- `src/types/` - TypeScript 类型定义

## 模块划分

### 1. Auth 模块 (认证模块)

```
Auth Module
├── Presentation
│   ├── pages/auth/login/index.vue          # 登录页
│   └── pages/auth/register/index.vue       # 注册页
├── Business
│   ├── stores/auth.ts                       # 认证状态
│   ├── composables/useAuth.ts              # 认证逻辑
│   └── services/AuthService.ts             # 认证服务
└── Data
    ├── api/auth.ts                          # 认证 API
    └── types/auth.d.ts                      # 认证类型
```

### 2. Device 模块 (设备模块)

```
Device Module
├── Presentation
│   ├── pages/device/list/index.vue         # 设备列表
│   ├── pages/device/detail/index.vue       # 设备详情
│   └── components/business/device/
│       ├── DeviceCard.vue                  # 设备卡片
│       └── DeviceControl.vue               # 设备控制
├── Business
│   ├── stores/device.ts                     # 设备状态
│   ├── composables/useDevice.ts            # 设备逻辑
│   └── services/DeviceService.ts           # 设备服务
└── Data
    ├── api/device.ts                        # 设备 API
    └── types/device.d.ts                    # 设备类型
```

### 3. Workspace 模块 (工作台模块)

```
Workspace Module
├── Presentation
│   ├── pages/workspace/index.vue           # 工作台
│   └── components/business/workspace/
│       ├── QuickActions.vue                # 快捷操作
│       └── ToDoList.vue                    # 待办事项
├── Business
│   ├── stores/workspace.ts                  # 工作台状态
│   └── composables/useWorkspace.ts         # 工作台逻辑
└── Data
    ├── api/workspace.ts                     # 工作台 API
    └── types/workspace.d.ts                 # 工作台类型
```

### 4. Dashboard 模块 (看板模块)

```
Dashboard Module
├── Presentation
│   ├── pages/dashboard/index.vue           # 看板
│   └── components/business/dashboard/
│       ├── ChartCard.vue                   # 图表卡片
│       └── RealTimeData.vue                # 实时数据
├── Business
│   ├── stores/dashboard.ts                  # 看板状态
│   └── composables/useDashboard.ts         # 看板逻辑
└── Data
    ├── api/dashboard.ts                     # 看板 API
    └── types/dashboard.d.ts                 # 看板类型
```

### 5. Message 模块 (消息模块)

```
Message Module
├── Presentation
│   ├── pages/message/list/index.vue        # 消息列表
│   ├── pages/message/detail/index.vue      # 消息详情
│   └── components/business/message/
│       └── MessageCard.vue                 # 消息卡片
├── Business
│   ├── stores/message.ts                    # 消息状态
│   └── composables/useMessage.ts           # 消息逻辑
└── Data
    ├── api/message.ts                       # 消息 API
    └── types/message.d.ts                   # 消息类型
```

### 6. Profile 模块 (个人中心模块)

```
Profile Module
├── Presentation
│   ├── pages/profile/center/index.vue      # 个人中心
│   ├── pages/profile/settings/index.vue    # 设置
│   └── components/business/profile/
│       ├── UserInfo.vue                    # 用户信息
│       └── SettingItem.vue                 # 设置项
├── Business
│   ├── stores/profile.ts                    # 个人中心状态
│   └── composables/useProfile.ts           # 个人中心逻辑
└── Data
    ├── api/profile.ts                       # 个人中心 API
    └── types/profile.d.ts                   # 个人中心类型
```

## 数据流

### 1. 用户登录流程

```
┌──────────┐
│   User   │ 用户输入用户名密码
└────┬─────┘
     │
     ↓
┌─────────────────────────┐
│ Login Page Component    │ 调用 useAuth composable
│ (pages/auth/login)      │
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ useAuth Composable      │ 调用 AuthService
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ AuthService             │ 调用 auth API
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ auth API                │ 发送 HTTP 请求
│ (api/auth.ts)           │
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ Cloud API               │ 返回响应
│ /api/auth/login         │
└────┬────────────────────┘
     │
     ↓ (Token + UserInfo)
┌─────────────────────────┐
│ auth Store              │ 保存状态
│ (stores/auth.ts)        │
└────┬────────────────────┘
     │
     ↓ (响应式状态更新)
┌─────────────────────────┐
│ Login Page Component    │ 页面跳转
└─────────────────────────┘
```

### 2. 设备列表加载流程

```
┌──────────┐
│   User   │ 打开设备列表页
└────┬─────┘
     │
     ↓
┌─────────────────────────┐
│ Device List Component   │ 调用 useDevice composable
│ (pages/device/list)     │
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ useDevice Composable    │ 调用 DeviceService
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ DeviceService           │ 调用 device API
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ device API              │ 发送 HTTP 请求
│ (api/device.ts)         │
└────┬────────────────────┘
     │
     ↓
┌─────────────────────────┐
│ Cloud API               │ 返回设备列表
│ /api/devices            │
└────┬────────────────────┘
     │
     ↓ (Device[])
┌─────────────────────────┐
│ device Store            │ 保存设备列表
│ (stores/device.ts)      │
└────┬────────────────────┘
     │
     ↓ (响应式状态更新)
┌─────────────────────────┐
│ Device List Component   │ 渲染设备列表
└─────────────────────────┘
```

## 状态管理架构

### Pinia Store 模块

```
Pinia Store
├── app.ts                 # 应用全局状态
│   - theme: string
│   - loading: boolean
│   - networkStatus: boolean
│
├── auth.ts                # 认证状态
│   - token: string
│   - userInfo: UserInfo
│   - isLoggedIn: boolean
│   - login()
│   - logout()
│
├── device.ts              # 设备状态
│   - devices: Device[]
│   - currentDevice: Device
│   - deviceStatus: Record
│   - fetchDevices()
│   - setCurrentDevice()
│
├── workspace.ts           # 工作台状态
│   - quickActions: QuickAction[]
│   - todoList: TodoItem[]
│   - notifications: Notification[]
│
├── dashboard.ts           # 看板状态
│   - charts: ChartData[]
│   - realTimeData: RealTimeData[]
│
├── message.ts             # 消息状态
│   - messages: Message[]
│   - unreadCount: number
│   - fetchMessages()
│   - markAsRead()
│
└── profile.ts             # 个人中心状态
    - userInfo: UserInfo
    - settings: UserSettings
```

## 路由设计

### TabBar 配置

```json
{
  "tabBar": {
    "color": "#999999",
    "selectedColor": "#007AFF",
    "backgroundColor": "#FFFFFF",
    "borderStyle": "black",
    "list": [
      {
        "pagePath": "pages/workspace/index",
        "text": "工作台",
        "iconPath": "static/icons/workspace.png",
        "selectedIconPath": "static/icons/workspace-active.png"
      },
      {
        "pagePath": "pages/device/list/index",
        "text": "设备",
        "iconPath": "static/icons/device.png",
        "selectedIconPath": "static/icons/device-active.png"
      },
      {
        "pagePath": "pages/dashboard/index",
        "text": "看板",
        "iconPath": "static/icons/dashboard.png",
        "selectedIconPath": "static/icons/dashboard-active.png"
      },
      {
        "pagePath": "pages/message/list/index",
        "text": "消息",
        "iconPath": "static/icons/message.png",
        "selectedIconPath": "static/icons/message-active.png"
      },
      {
        "pagePath": "pages/profile/center/index",
        "text": "我的",
        "iconPath": "static/icons/profile.png",
        "selectedIconPath": "static/icons/profile-active.png"
      }
    ]
  }
}
```

### 页面路由

```json
{
  "pages": [
    {
      "path": "pages/index/index",
      "style": { "navigationBarTitleText": "启动页" }
    },
    {
      "path": "pages/auth/login/index",
      "style": { "navigationBarTitleText": "登录" }
    },
    {
      "path": "pages/workspace/index",
      "style": { "navigationBarTitleText": "工作台" }
    },
    {
      "path": "pages/device/list/index",
      "style": { "navigationBarTitleText": "设备列表" }
    },
    {
      "path": "pages/device/detail/index",
      "style": { "navigationBarTitleText": "设备详情" }
    },
    {
      "path": "pages/dashboard/index",
      "style": { "navigationBarTitleText": "看板" }
    },
    {
      "path": "pages/message/list/index",
      "style": { "navigationBarTitleText": "消息" }
    },
    {
      "path": "pages/message/detail/index",
      "style": { "navigationBarTitleText": "消息详情" }
    },
    {
      "path": "pages/profile/center/index",
      "style": { "navigationBarTitleText": "个人中心" }
    },
    {
      "path": "pages/profile/settings/index",
      "style": { "navigationBarTitleText": "设置" }
    }
  ]
}
```

## 安全架构

### 认证流程

```
┌──────────────┐
│  User Login  │
└──────┬───────┘
       │
       ↓
┌────────────────────┐
│ Send Credentials   │ POST /api/auth/login
└──────┬─────────────┘
       │
       ↓
┌────────────────────┐
│ Receive Token      │ JWT Token + Refresh Token
└──────┬─────────────┘
       │
       ↓
┌────────────────────┐
│ Store Token        │ uni.storage (加密存储)
└──────┬─────────────┘
       │
       ↓
┌────────────────────┐
│ Add to Headers     │ Authorization: Bearer {token}
└──────┬─────────────┘
       │
       ↓
┌────────────────────┐
│ Auto Refresh       │ Token expired -> refresh
└────────────────────┘
```

### 数据加密

```
Sensitive Data
├── Token
│   └── uni.setStorageSync('token', token) # 加密存储
├── RefreshToken
│   └── uni.setStorageSync('refreshToken', refreshToken) # 加密存储
└── UserSecret
    └── uni.setStorageSync('userSecret', secret) # 加密存储
```

## 技术栈总结

| 层级 | 技术选型 | 说明 |
|------|---------|------|
| **展示层** | Vue 3 + UniApp | 页面渲染、路由、组件 |
| **业务层** | Pinia + Composables | 状态管理、业务逻辑 |
| **数据层** | TypeScript + uni.request | 类型安全、HTTP 请求 |
| **工具层** | ESLint + Prettier + Vite | 代码规范、构建工具 |
| **UI 组件** | UniUI | 官方 UI 组件库 |

## 设计原则

1. **单一职责**: 每个模块、组件只负责一项功能
2. **松耦合**: 模块间依赖最小化,通过接口通信
3. **高内聚**: 相关功能聚合在同一模块内
4. **可测试**: 分层架构便于单元测试和集成测试
5. **可扩展**: 模块化设计便于功能扩展
6. **可维护**: 清晰的目录结构和代码组织

---

**文档版本**: 1.0
**最后更新**: 2026-01-28
**维护人**: Claude Code
