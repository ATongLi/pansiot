# Feature-to-Code Mapping: Scada主页面框架

## 映射表信息

**文档ID**: feature-to-code-map
**文档标题**: 功能到代码映射表
**关联需求**: REQ-001
**关联功能需求**: FE-001
**关联实施计划**: IMP-001
**创建日期**: 2026-01-20
**目标平台**: Scada
**状态**: ✅ 已完成

## 概述

本文档提供了从业务需求到代码实现的完整可追溯性映射。通过此文档，可以：

1. 追踪每个功能需求对应的代码文件
2. 快速定位功能的实现位置
3. 理解需求与代码之间的关系
4. 支持影响分析和变更管理

## 映射表结构

```
REQ-001 (原始需求)
  ↓
FE-001 (功能需求分解)
  ↓
FE-001-x (具体功能需求)
  ↓
Code Files (实现文件)
  ↓
Functions/Components (具体实现)
```

---

## FE-001-1: 顶部栏组件 (TopBar)

### 功能需求

**描述**: 应用顶部标题栏，显示应用标题和窗口控制

**关联文档**: FE-001.md#FE-001-1

### 代码实现

#### 1. TopBar组件

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/TopBar.tsx`

**代码行数**: ~20行

**关键实现**:
```typescript
const TopBar: React.FC = observer(() => {
  return (
    <div className="topbar">
      <div className="topbar-drag-region">
        <h1 className="topbar-title">PanTools</h1>
      </div>
      <WindowControls />
    </div>
  )
})
```

**职责**:
- 渲染顶部栏布局
- 显示"PanTools"标题
- 集成WindowControls组件
- 支持窗口拖拽区域

#### 2. TopBar样式

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/TopBar.css`

**代码行数**: ~30行

**关键样式**:
```css
.topbar {
  height: var(--topbar-height);  /* 40px */
  background: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 var(--spacing-md);
  -webkit-app-region: drag;  /* 支持拖拽 */
  user-select: none;
}

.topbar-title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
}
```

**职责**:
- 定义顶部栏高度（40px）
- 设置窗口拖拽区域
- 设置布局和间距

#### 3. WindowControls组件

**文件路径**: `platforms/scada/packages/renderer/src/components/common/WindowControls.tsx`

**代码行数**: ~50行

**关键实现**:
```typescript
const WindowControls: React.FC = observer(() => {
  const handleMinimize = () => uiStore.handleWindowAction('minimize')
  const handleMaximize = () => {
    const action = uiStore.windowState.isMaximized ? 'restore' : 'maximize'
    uiStore.handleWindowAction(action)
  }
  const handleClose = () => uiStore.handleWindowAction('close')

  return (
    <div className="window-controls">
      <button className="window-control window-control--minimize" ...>
        <span className="window-control__icon">─</span>
      </button>
      <button className="window-control window-control--maximize" ...>
        <span className="window-control__icon">
          {uiStore.windowState.isMaximized ? '❐' : '□'}
        </span>
      </button>
      <button className="window-control window-control--close" ...>
        <span className="window-control__icon">✕</span>
      </button>
    </div>
  )
})
```

**职责**:
- 渲染三个窗口控制按钮
- 处理最小化、最大化、关闭操作
- 动态显示最大化/还原图标
- 集成MobX状态（windowState）

#### 4. WindowControls样式

**文件路径**: `platforms/scada/packages/renderer/src/components/common/WindowControls.css`

**代码行数**: ~40行

**关键样式**:
```css
.window-control--minimize {
  background: #ffbd2e;  /* 黄色 */
}

.window-control--maximize {
  background: #28c940;  /* 绿色 */
}

.window-control--close {
  background: #ff5f57;  /* 红色 */
}

.window-control:hover .window-control__icon {
  opacity: 1;  /* 悬停时显示图标 */
}
```

**职责**:
- 设置交通灯颜色（黄、绿、红）
- 实现悬停效果
- 禁用窗口控制区域的拖拽

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 标题居左显示 | TopBar.tsx:11 | ✅ |
| 窗口控制按钮居右显示 | TopBar.css:6 (justify-content: space-between) | ✅ |
| 按钮有悬停效果 | WindowControls.css:20-23 | ✅ |
| 支持窗口拖拽 | TopBar.css:11 (-webkit-app-region: drag) | ✅ |

---

## FE-001-2: 侧边栏导航组件 (Sidebar)

### 功能需求

**描述**: 左侧垂直导航栏，提供主要功能入口

**关联文档**: FE-001.md#FE-001-2

### 代码实现

#### 1. Sidebar组件

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/Sidebar.tsx`

**代码行数**: ~20行

**关键实现**:
```typescript
const Sidebar: React.FC = observer(() => {
  return (
    <div className="sidebar">
      {uiStore.navigationItems.map((item) => (
        <NavItem
          key={item.id}
          item={item}
          isActive={uiStore.activeNavItem === item.id}
          onClick={() => uiStore.setActiveNavItem(item.id)}
        />
      ))}
    </div>
  )
})
```

**职责**:
- 渲染导航列表
- 从uiStore读取导航项配置
- 传递激活状态给NavItem
- 处理导航项点击事件

#### 2. Sidebar样式

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/Sidebar.css`

**代码行数**: ~15行

**关键样式**:
```css
.sidebar {
  width: var(--sidebar-width);  /* 80px */
  height: calc(100vh - var(--topbar-height));
  background: var(--color-bg-secondary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  padding: var(--spacing-md) 0;
}
```

**职责**:
- 定义侧边栏宽度（80px）
- 设置垂直布局
- 设置边框和内边距

#### 3. NavItem组件

**文件路径**: `platforms/scada/packages/renderer/src/components/navigation/NavItem.tsx`

**代码行数**: ~30行

**关键实现**:
```typescript
const NavItem: React.FC<NavItemProps> = ({ item, isActive, onClick }) => {
  return (
    <div
      className={`nav-item ${isActive ? 'nav-item--active' : ''}`}
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          onClick()
        }
      }}
    >
      <div className="nav-item__icon">{item.icon}</div>
      <div className="nav-item__label">{item.label}</div>
    </div>
  )
}
```

**职责**:
- 渲染单个导航项
- 显示图标和标签
- 处理点击和键盘事件
- 根据isActive应用激活样式

#### 4. NavItem样式

**文件路径**: `platforms/scada/packages/renderer/src/components/navigation/NavItem.css`

**代码行数**: ~40行

**关键样式**:
```css
.nav-item {
  height: var(--nav-item-height);  /* 70px */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color var(--transition-fast);
  margin: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius);
}

.nav-item--active {
  background: var(--color-accent-active);  /* #FF9999 */
}

.nav-item:hover {
  background: var(--color-bg-tertiary);
}
```

**职责**:
- 定义导航项高度（70px）
- 设置激活状态背景色（#FF9999）
- 实现悬停效果
- 设置圆角和过渡动画

#### 5. 导航配置常量

**文件路径**: `platforms/scada/packages/renderer/src/constants/navigation.ts`

**代码行数**: ~15行

**关键实现**:
```typescript
export const NAVIGATION_ITEMS: NavItem[] = [
  { id: 'home', label: '首页', icon: '🏠', path: '/' },
  { id: 'local', label: '本地', icon: '💾', path: '/local' },
  { id: 'cloud', label: '云端', icon: '☁️', path: '/cloud' },
  { id: 'tools', label: '工具', icon: '🔧', path: '/tools' },
  { id: 'user', label: 'User', icon: '👤', path: '/user' },
]

export const DEFAULT_NAV_ITEM = 'home'
```

**职责**:
- 定义5个导航项配置
- 设置默认激活项
- 提供导航项数据源

#### 6. 导航类型定义

**文件路径**: `platforms/scada/packages/renderer/src/types/navigation.ts`

**代码行数**: ~10行

**关键实现**:
```typescript
export interface NavItem {
  id: string
  label: string
  icon?: string
  path: string
}
```

**职责**:
- 定义NavItem接口
- 提供类型安全

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 5个导航项全部显示 | navigation.ts:4-10 | ✅ |
| 默认激活"首页" | navigation.ts:12 → uiStore.ts:8 | ✅ |
| 点击导航项切换激活状态 | Sidebar.tsx:14 (onClick) | ✅ |
| 只能有一个激活项 | uiStore.ts:18-22 (setActiveNavItem) | ✅ |
| 激活项背景色为 #FF9999 | NavItem.css:19 (var(--color-accent-active)) | ✅ |
| 悬停时背景色变浅 | NavItem.css:24-26 (:hover) | ✅ |
| 支持 Tab 键导航 | NavItem.tsx:17 (tabIndex={0}) | ✅ |
| 支持 Enter/Space 激活 | NavItem.tsx:18-22 (onKeyDown) | ✅ |

---

## FE-001-3: 主内容区组件 (MainContent)

### 功能需求

**描述**: 页面主要内容显示区域

**关联文档**: FE-001.md#FE-001-3

### 代码实现

#### 1. MainContent组件

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/MainContent.tsx`

**代码行数**: ~20行

**关键实现**:
```typescript
const MainContent: React.FC = () => {
  return (
    <div className="main-content">
      <div className="main-content__inner">
        <ActionButtons />
        <RecentProjects />
      </div>
    </div>
  )
}
```

**职责**:
- 渲染主内容区布局
- 集成ActionButtons组件
- 集成RecentProjects组件
- 提供响应式容器

#### 2. MainContent样式

**文件路径**: `platforms/scada/packages/renderer/src/components/layout/MainContent.css`

**代码行数**: ~20行

**关键样式**:
```css
.main-content {
  flex: 1;  /* 占据剩余空间 */
  height: calc(100vh - var(--topbar-height));
  background: var(--color-bg-primary);
  overflow-y: auto;  /* 支持垂直滚动 */
}

.main-content__inner {
  max-width: 1200px;  /* 最大宽度限制 */
  margin: 0 auto;  /* 居中对齐 */
  padding: var(--spacing-xl);  /* 32px内边距 */
}
```

**职责**:
- 设置flex布局（flex: 1）
- 启用垂直滚动
- 限制最大宽度（1200px）
- 居中对齐内容
- 设置内边距（32px）

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 占据剩余空间 | MainContent.css:3 (flex: 1) | ✅ |
| 内容超出时可滚动 | MainContent.css:6 (overflow-y: auto) | ✅ |
| 内容居中显示 | MainContent.css:11 (margin: 0 auto) | ✅ |
| 内边距正确（32px） | MainContent.css:12 (padding: var(--spacing-xl)) | ✅ |

---

## FE-001-4: 操作按钮组组件 (ActionButtons)

### 功能需求

**描述**: 三个主要操作按钮

**关联文档**: FE-001.md#FE-001-4

### 代码实现

#### 1. ActionButtons组件

**文件路径**: `platforms/scada/packages/renderer/src/components/workspace/ActionButtons.tsx`

**代码行数**: ~30行

**关键实现**:
```typescript
const ActionButtons: React.FC = () => {
  const buttons: ActionButton[] = [
    { id: 'open', label: '从文件打开', onClick: () => console.log('Open file') },
    { id: 'new', label: '新建工程', onClick: () => console.log('New project') },
    { id: 'copy', label: '复制工程', onClick: () => console.log('Copy project') },
  ]

  return (
    <div className="action-buttons">
      {buttons.map((button) => (
        <button key={button.id} className="action-button" onClick={button.onClick}>
          {button.label}
        </button>
      ))}
    </div>
  )
}
```

**职责**:
- 定义三个操作按钮配置
- 渲染按钮列表
- 处理按钮点击（当前为Mock）

#### 2. ActionButtons样式

**文件路径**: `platforms/scada/packages/renderer/src/components/workspace/ActionButtons.css`

**代码行数**: ~30行

**关键样式**:
```css
.action-buttons {
  display: flex;
  gap: var(--spacing-md);  /* 16px间距 */
  margin-bottom: var(--spacing-xl);
}

.action-button {
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  font-size: var(--font-size-md);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.action-button:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-accent-active);  /* #FF9999 */
}

.action-button:active {
  transform: translateY(1px);  /* 下压效果 */
}
```

**职责**:
- 设置水平布局（flex）
- 设置按钮间距（16px）
- 实现悬停高亮效果
- 实现点击下压效果

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 三个按钮水平排列 | ActionButtons.css:3 (display: flex) | ✅ |
| 按钮间距一致（16px） | ActionButtons.css:4 (gap: var(--spacing-md)) | ✅ |
| 悬停时有边框高亮（#FF9999） | ActionButtons.css:24 (border-color: var(--color-accent-active)) | ✅ |
| 点击时有下压效果 | ActionButtons.css:27 (transform: translateY(1px)) | ✅ |
| 当前为Mock实现 | ActionButtons.tsx:7-9 (console.log) | ✅ |

---

## FE-001-5: 最近工程列表组件 (RecentProjects)

### 功能需求

**描述**: 显示最近打开的工程项目

**关联文档**: FE-001.md#FE-001-5

### 代码实现

#### 1. RecentProjects组件

**文件路径**: `platforms/scada/packages/renderer/src/components/workspace/RecentProjects.tsx`

**代码行数**: ~40行

**关键实现**:
```typescript
const RecentProjects: React.FC = () => {
  const recentProjects = [
    { id: '1', name: '工程示例 1', lastOpened: '2026-01-20' },
    { id: '2', name: '工程示例 2', lastOpened: '2026-01-19' },
    { id: '3', name: '工程示例 3', lastOpened: '2026-01-18' },
  ]

  return (
    <div className="recent-projects">
      <h2 className="recent-projects__title">最近工程</h2>
      <div className="recent-projects__grid">
        {recentProjects.map((project) => (
          <div key={project.id} className="project-card">
            <div className="project-card__icon">📁</div>
            <div className="project-card__info">
              <div className="project-card__name">{project.name}</div>
              <div className="project-card__date">{project.lastOpened}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
```

**职责**:
- 定义Mock项目数据（3个）
- 渲染项目卡片列表
- 显示项目标题和网格
- 显示项目图标、名称、日期

#### 2. RecentProjects样式

**文件路径**: `platforms/scada/packages/renderer/src/components/workspace/RecentProjects.css`

**代码行数**: ~60行

**关键样式**:
```css
.recent-projects__title {
  font-size: var(--font-size-lg);
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-md);
}

.recent-projects__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));  /* 响应式网格 */
  gap: var(--spacing-md);
}

.project-card {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: var(--spacing-md);
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.project-card:hover {
  border-color: var(--color-accent-active);  /* #FF9999 */
  box-shadow: 0 2px 8px var(--color-shadow);  /* 阴影效果 */
}
```

**职责**:
- 设置标题样式
- 实现响应式网格布局
- 设置卡片样式
- 实现悬停效果（边框高亮 + 阴影）

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 标题正确显示 | RecentProjects.tsx:15 | ✅ |
| 卡片网格布局（最小宽度200px） | RecentProjects.css:17 (minmax(200px, 1fr)) | ✅ |
| 显示3个示例工程 | RecentProjects.tsx:7-11 | ✅ |
| 卡片有悬停效果 | RecentProjects.css:40-44 (:hover) | ✅ |
| 边框颜色变为 #FF9999 | RecentProjects.css:41 (border-color) | ✅ |
| 显示阴影 | RecentProjects.css:42 (box-shadow) | ✅ |

---

## FE-001-6: 状态管理 (UI Store)

### 功能需求

**描述**: MobX状态管理，管理UI状态

**关联文档**: FE-001.md#FE-001-6

### 代码实现

#### 1. UI Store

**文件路径**: `platforms/scada/packages/renderer/src/store/uiStore.ts`

**代码行数**: ~45行

**关键实现**:
```typescript
export class UIStore {
  // 导航状态
  activeNavItem: string = DEFAULT_NAV_ITEM
  navigationItems = NAVIGATION_ITEMS

  // 窗口状态
  windowState: WindowState = {
    isMaximized: false,
    isFullscreen: false,
  }

  constructor() {
    makeAutoObservable(this)  // 自动observable
  }

  setActiveNavItem(itemId: string) {
    const item = this.navigationItems.find(i => i.id === itemId)
    if (item) {
      this.activeNavItem = itemId
    }
  }

  async handleWindowAction(action: WindowAction) {
    console.log('Window action:', action)

    switch (action) {
      case 'minimize':
        break
      case 'maximize':
        this.windowState.isMaximized = !this.windowState.isMaximized
        break
      case 'restore':
        this.windowState.isMaximized = false
        break
      case 'close':
        break
    }
  }
}

export const uiStore = new UIStore()
```

**职责**:
- 管理导航状态（activeNavItem）
- 管理窗口状态（windowState）
- 提供setActiveNavItem方法
- 提供handleWindowAction方法
- 使用makeAutoObservable自动追踪

#### 2. Store导出

**文件路径**: `platforms/scada/packages/renderer/src/store/index.ts`

**代码行数**: ~5行

**关键实现**:
```typescript
export { uiStore, UIStore } from './uiStore'
```

**职责**:
- 导出uiStore单例
- 导出UIStore类型

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 使用 makeAutoObservable | uiStore.ts:16 (makeAutoObservable(this)) | ✅ |
| 状态变化触发组件重新渲染 | Sidebar.tsx:6, TopBar.tsx:6 (observer()) | ✅ |
| setActiveNavItem 方法正常工作 | uiStore.ts:18-23 | ✅ |
| handleWindowAction 处理窗口操作 | uiStore.ts:25-41 | ✅ |
| 当前为Mock实现 | uiStore.ts:28 (console.log) | ✅ |

---

## FE-001-7: 类型系统 (TypeScript Types)

### 功能需求

**描述**: TypeScript类型定义

**关联文档**: FE-001.md#FE-001-7

### 代码实现

#### 1. 导航类型

**文件路径**: `platforms/scada/packages/renderer/src/types/navigation.ts`

**代码行数**: ~10行

**关键实现**:
```typescript
export interface NavItem {
  id: string
  label: string
  icon?: string
  path: string
}
```

**职责**:
- 定义NavItem接口
- 提供导航项类型安全

#### 2. 窗口类型

**文件路径**: `platforms/scada/packages/renderer/src/types/window.ts`

**代码行数**: ~10行

**关键实现**:
```typescript
export interface WindowState {
  isMaximized: boolean
  isFullscreen: boolean
}

export type WindowAction =
  | 'minimize'
  | 'maximize'
  | 'close'
  | 'restore'
```

**职责**:
- 定义WindowState接口
- 定义WindowAction联合类型

#### 3. 项目类型

**文件路径**: `platforms/scada/packages/renderer/src/types/project.ts`

**代码行数**: ~10行

**关键实现**:
```typescript
export interface RecentProject {
  id: string
  name: string
  path: string
  lastOpened: Date
  thumbnail?: string
}
```

**职责**:
- 定义RecentProject接口
- 提供项目类型安全

### 验收标准映射

| 验收标准 | 实现位置 | 状态 |
|---------|---------|------|
| 所有类型定义完整 | navigation.ts, window.ts, project.ts | ✅ |
| 类型导出正确 | 各文件export语句 | ✅ |
| 在组件中正确使用 | NavItem.tsx:3, uiStore.ts:5, etc. | ✅ |

---

## 全局配置文件

### CSS变量系统

**文件路径**: `platforms/scada/packages/renderer/src/index.css`

**代码行数**: ~80行

**关键内容**:
```css
:root {
  /* 颜色系统 */
  --color-bg-primary: #f5f5f5;
  --color-bg-secondary: #ffffff;
  --color-accent-active: #FF9999;
  --color-text-primary: #333333;
  --color-border: #d0d0d0;

  /* 尺寸系统 */
  --topbar-height: 40px;
  --sidebar-width: 80px;
  --nav-item-height: 70px;
  --border-radius: 4px;

  /* 间距系统 (4px网格) */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;

  /* 字体系统 */
  --font-size-sm: 12px;
  --font-size-md: 14px;
  --font-size-lg: 16px;

  /* 动画系统 */
  --transition-fast: 150ms ease;
}
```

**职责**:
- 定义全局CSS变量
- 建立设计系统
- 提供主题化能力

### 应用入口

**文件路径**: `platforms/scada/packages/renderer/src/App.tsx`

**代码行数**: ~20行

**关键实现**:
```typescript
const App: React.FC = observer(() => {
  return (
    <div className="app">
      <TopBar />
      <div className="app-body">
        <Sidebar />
        <MainContent />
      </div>
    </div>
  )
})
```

**职责**:
- 组合所有布局组件
- 使用observer包装（MobX集成）
- 定义应用整体结构

**文件路径**: `platforms/scada/packages/renderer/src/main.tsx`

**代码行数**: ~10行

**关键实现**:
```typescript
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
```

**职责**:
- React应用入口
- 渲染根组件

### 工具函数

**文件路径**: `platforms/scada/packages/renderer/src/utils/electron.ts`

**代码行数**: ~35行

**关键实现**:
```typescript
export const electronAPI = {
  minimize: () => console.log('Mock: minimize window'),
  maximize: () => console.log('Mock: maximize window'),
  restore: () => console.log('Mock: restore window'),
  close: () => console.log('Mock: close window'),
  openProject: () => console.log('Mock: open project dialog'),
  saveProject: () => console.log('Mock: save project'),
}

export const isElectron = (): boolean => {
  return typeof window !== 'undefined' &&
         window.process !== undefined &&
         window.process.type === 'renderer'
}

export const getElectronAPI = () => {
  if (isElectron() && (window as any).electronAPI) {
    return (window as any).electronAPI
  }
  return electronAPI
}
```

**职责**:
- 提供Electron API封装
- 实现Mock实现用于浏览器开发
- 检测Electron环境

### 配置文件

#### Vite配置

**文件路径**: `platforms/scada/packages/renderer/vite.config.ts`

**代码行数**: ~20行

**关键内容**:
```typescript
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@store': path.resolve(__dirname, './src/store'),
      '@types': path.resolve(__dirname, './src/types'),
      '@constants': path.resolve(__dirname, './src/constants'),
      '@utils': path.resolve(__dirname, './src/utils'),
    }
  },
  base: './',  // Important for Electron
  build: {
    outDir: 'dist',
    sourcemap: true
  }
})
```

**职责**:
- 配置路径别名
- 配置React插件
- 配置Electron兼容

#### TypeScript配置

**文件路径**: `platforms/scada/packages/renderer/tsconfig.json`

**代码行数**: ~30行

**关键内容**:
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "jsx": "react-jsx",
    "strict": true,
    "experimentalDecorators": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@components/*": ["./src/components/*"],
      "@store/*": ["./src/store/*"],
      "@types/*": ["./src/types/*"],
      "@constants/*": ["./src/constants/*"],
      "@utils/*": ["./src/utils/*"]
    }
  }
}
```

**职责**:
- 配置TypeScript编译选项
- 启用严格模式
- 配置装饰器支持（MobX）
- 配置路径映射

---

## 代码统计

### 文件总数

| 类别 | 数量 |
|------|------|
| **配置文件** | 3 |
| **核心应用** | 3 |
| **状态管理** | 2 |
| **类型定义** | 3 |
| **常量** | 1 |
| **布局组件** | 5 (组件) + 5 (CSS) = 10 |
| **工作区组件** | 2 (组件) + 2 (CSS) = 4 |
| **工具** | 1 |
| **总计** | **27** |

### 代码行数

| 类别 | 代码行数 | 说明 |
|------|---------|------|
| TypeScript (TSX) | ~350行 | 组件和逻辑 |
| CSS | ~250行 | 样式 |
| 配置文件 | ~50行 | JSON, TS配置 |
| 注释和空行 | ~100行 | 文档和格式 |
| **总计** | **~600行** | 含注释和空行 |

### 组件数量

| 组件类型 | 数量 |
|---------|------|
| 布局组件 | 3 (TopBar, Sidebar, MainContent) |
| 导航组件 | 1 (NavItem) |
| 工作区组件 | 2 (ActionButtons, RecentProjects) |
| 通用组件 | 1 (WindowControls) |
| **总计** | **7** |

---

## 依赖关系图

```
App.tsx (Root)
├── index.css (Global Styles)
├── TopBar.tsx
│   ├── TopBar.css
│   └── WindowControls.tsx
│       ├── WindowControls.css
│       └── uiStore.ts (MobX)
├── Sidebar.tsx
│   ├── Sidebar.css
│   ├── NavItem.tsx
│   │   ├── NavItem.css
│   │   └── navigation.ts (Constants)
│   └── uiStore.ts (MobX)
└── MainContent.tsx
    ├── MainContent.css
    ├── ActionButtons.tsx
    │   └── ActionButtons.css
    └── RecentProjects.tsx
        └── RecentProjects.css

uiStore.ts (MobX State Management)
├── navigation.ts (NAVIGATION_ITEMS)
├── types/navigation.ts (NavItem interface)
└── types/window.ts (WindowState, WindowAction)
```

---

## 影响分析

### 修改影响范围

#### 如果修改导航项配置

**影响文件**:
- `src/constants/navigation.ts` (直接修改)
- `src/components/layout/Sidebar.tsx` (使用配置)
- `src/store/uiStore.ts` (导入配置)

**影响范围**: 中等

**修改建议**:
1. 修改navigation.ts中的NAVIGATION_ITEMS
2. 无需修改组件代码
3. 验证导航项显示正常

#### 如果修改颜色主题

**影响文件**:
- `src/index.css` (CSS变量定义)
- 所有CSS文件 (使用CSS变量)

**影响范围**: 广泛

**修改建议**:
1. 修改index.css中的CSS变量
2. 全局自动生效
3. 验证所有组件颜色

#### 如果修改布局尺寸

**影响文件**:
- `src/index.css` (尺寸CSS变量)
- 相应组件CSS文件

**影响范围**: 局部

**修改建议**:
1. 修改index.css中的尺寸变量
2. 特定组件自动生效
3. 验证布局无破坏

---

## 可追溯性矩阵

### REQ-001 → FE-001 → Code Files

| REQ-001需求 | FE-001功能需求 | 实现文件 | 状态 |
|------------|---------------|---------|------|
| 整体风格要求 | - | index.css (CSS变量系统) | ✅ |
| 顶部栏 (40px) | FE-001-1 | TopBar.tsx/css, WindowControls.tsx/css | ✅ |
| 侧边栏 (80px) | FE-001-2 | Sidebar.tsx/css, NavItem.tsx/css | ✅ |
| 主内容区 | FE-001-3 | MainContent.tsx/css | ✅ |
| 操作按钮组 | FE-001-4 | ActionButtons.tsx/css | ✅ |
| 最近工程列表 | FE-001-5 | RecentProjects.tsx/css | ✅ |
| 状态管理 | FE-001-6 | uiStore.ts, store/index.ts | ✅ |
| 类型系统 | FE-001-7 | types/*.ts | ✅ |

---

## 实施验证

### 功能验证清单

- [x] FE-001-1: TopBar + WindowControls正常工作
- [x] FE-001-2: Sidebar + NavItem导航正常
- [x] FE-001-3: MainContent布局正确
- [x] FE-001-4: ActionButtons可点击
- [x] FE-001-5: RecentProjects显示正常
- [x] FE-001-6: MobX状态管理正常
- [x] FE-001-7: TypeScript无类型错误

### 视觉验证清单

- [x] 颜色符合设计规范（#FF9999等）
- [x] 尺寸符合设计规范（40px, 80px等）
- [x] 间距符合4px网格系统
- [x] 字体清晰，大小合适
- [x] 圆角统一为4px

### 技术验证清单

- [x] TypeScript编译无错误
- [x] Vite开发服务器正常启动
- [x] MobX响应式正常工作
- [x] CSS变量正确应用
- [x] 路径别名正确解析

---

## 后续扩展

### 待实现功能映射

| 功能 | 依赖文件 | 状态 |
|------|---------|------|
| Electron窗口控制 | utils/electron.ts, Electron Main Process | 🔄 待实现 |
| 真实项目数据 | components/workspace/RecentProjects.tsx | 🔄 待实现 |
| 路由系统 | App.tsx, React Router | 🔄 待实现 |
| 国际化 | 所有组件, i18next | 🔄 待实现 |
| 主题切换 | index.css, 主题Provider | 🔄 待实现 |

### 技术债务

| 项目 | 文件 | 描述 | 优先级 |
|------|------|------|--------|
| Mock实现 | utils/electron.ts | 需要集成真实Electron API | P1 |
| 硬编码数据 | RecentProjects.tsx | 需要从文件系统加载 | P2 |
| 无错误边界 | App.tsx | 需要添加ErrorBoundary | P2 |
| 无路由 | App.tsx | 需要集成React Router | P3 |

---

## 变更历史

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| 1.0 | 2026-01-20 | 初始版本，完成功能到代码映射 | Claude Code |

## 参考资料

### 相关文档

- REQ-001: 原始需求
- FE-001: 功能需求
- US-001: 用户故事
- SOL-001: 技术方案
- ADR-001: 架构决策
- IMP-001: 实施计划

### 代码位置

- 根目录: `platforms/scada/packages/renderer/`
- 源代码: `src/`
- 组件: `src/components/`
- 状态: `src/store/`
- 类型: `src/types/`
- 常量: `src/constants/`

---

# Feature-to-Code Mapping: 云平台账号系统

## 映射表信息

**文档ID**: feature-to-code-map (Cloud)
**文档标题**: 功能到代码映射表
**关联需求**: REQ-007
**关联功能需求**: FE-007-01 ~ FE-007-09
**关联实施计划**: IMP-007
**创建日期**: 2026-01-28
**目标平台**: Cloud
**状态**: 🔄 进行中 (Phase 1 - 60%)

## 概述

云平台账号系统的功能到代码映射，包括多租户组织管理、RBAC权限模型、用户认证注册、前端动态界面、功能模块与配额管理、系统审计日志等9个核心功能模块。

---

## Phase 1: 基础架构搭建 ✅

### 基础设施代码映射

| 功能ID | 功能点 | 文件路径 | 类/函数 | 状态 |
|--------|--------|---------|---------|------|
| - | 配置管理 | internal/config/config.go | Config.LoadConfig() | ✅ |
| - | 日志系统 | pkg/logger/logger.go | InitLogger(), GetLogger() | ✅ |
| - | 数据库连接 | pkg/database/postgres.go | InitPostgres(), GetDB() | ✅ |
| - | Redis连接 | pkg/database/redis.go | InitRedis(), GetRedis() | ✅ |
| - | 统一响应 | pkg/response/response.go | Success(), Error() | ✅ |
| - | JWT认证 | internal/middleware/auth.go | GenerateToken(), ParseToken() | ✅ |
| - | CORS中间件 | internal/middleware/cors.go | CORS() | ✅ |
| - | 租户隔离 | internal/middleware/tenant.go | TenantScope() | ⏳ |
| - | 主入口 | cmd/server/main.go | main() | ✅ |
| - | 构建脚本 | Makefile | all, run, build | ✅ |
| - | Docker配置 | Dockerfile, docker-compose.yml | - | ✅ |
| - | 数据库脚本 | scripts/init.sql | generate_serial_number() | ✅ |

### 核心数据模型映射

| 功能ID | 功能点 | 文件路径 | 结构体 | 状态 |
|--------|--------|---------|--------|------|
| FE-007-01 | 租户模型 | internal/models/tenant.go | Tenant | ✅ |
| FE-007-02 | 用户模型 | internal/models/user.go | User | ✅ |
| FE-007-03 | 角色模型 | internal/models/user.go | Role | ✅ |
| FE-007-03 | 用户角色关联 | internal/models/user.go | UserRole | ✅ |
| FE-007-03 | 权限模型 | internal/models/permission.go | Permission | ✅ |
| FE-007-03 | 角色权限关联 | internal/models/permission.go | RolePermission | ✅ |
| FE-007-09 | 审计日志表 | scripts/init.sql | audit_logs | ✅ |

---

## FE-007-01: 多租户组织管理

### 功能需求

**描述**: 三层租户架构（平台超管 → 集成商 → 下游客户），支持组织树管理

**关联文档**: FE-007-01.md

### 代码实现

#### 1. Tenant数据模型

**文件路径**: `internal/models/tenant.go`

**代码行数**: ~70行

**关键字段**:
```go
type Tenant struct {
    ID               int64      `json:"id"`
    SerialNumber     string     `json:"serial_number"`     // 8位企业序列号
    Name             string     `json:"name"`               // 企业名称
    TenantType       string     `json:"tenant_type"`        // INTEGRATOR/TERMINAL
    Industry         string     `json:"industry"`           // 所属行业
    ParentTenantID   *int64     `json:"parent_tenant_id"`   // 上级租户ID
    Status           string     `json:"status"`             // ACTIVE/SUSPENDED/DELETED
    MaxSubTenants    int        `json:"max_sub_tenants"`    // 最大子租户数
    MaxUsers         int        `json:"max_users"`          // 最大用户数
    MaxDevices       int        `json:"max_devices"`        // 最大设备数
    MaxStorageGB     int        `json:"max_storage_gb"`     // 最大存储空间
}
```

**职责**:
- 定义租户数据结构
- 支持租户层级（parent_tenant_id）
- 支持配额限制
- 支持租户类型（集成商/下游客户）

#### 2. 企业序列号生成

**文件路径**: `scripts/init.sql`

**代码行数**: ~50行

**关键SQL**:
```sql
CREATE OR REPLACE FUNCTION generate_serial_number()
RETURNS VARCHAR(8) AS $$
DECLARE
  prefix VARCHAR(4);
  suffix INT;
  serial_number VARCHAR(8);
BEGIN
  -- 1. 生成4位随机字符
  prefix := upper(substring(encode(gen_random_bytes(3), 'base64'), 1, 4));
  prefix := regexp_replace(prefix, '[^A-Z0-9]', '', 'g');

  -- 2. 获取下一个自增ID
  suffix := nextval('serial_number_seq');

  -- 3. 拼接（4位随机+4位自增）
  serial_number := prefix || LPAD(suffix::TEXT, 4, '0');

  RETURN serial_number;
END;
$$ LANGUAGE plpgsql;
```

**职责**:
- 生成8位企业序列号
- 格式：4位随机字符 + 4位自增ID
- 示例：A3F20001, X7K10001

---

## FE-007-02: 用户注册与登录

### 功能需求

**描述**: 支持邮箱/手机号注册，新企业注册/加入已有企业，JWT认证

**关联文档**: FE-007-02.md

### 代码实现

#### 1. User数据模型

**文件路径**: `internal/models/user.go`

**代码行数**: ~40行

**关键字段**:
```go
type User struct {
    ID               int64      `json:"id"`
    TenantID         int64      `json:"tenant_id"`          // 归属租户
    Username         string     `json:"username"`
    Email            string     `json:"email"`               // 邮箱
    Phone            string     `json:"phone"`               // 手机号
    PhoneCountryCode string     `json:"phone_country_code"` // 国家码
    PasswordHash     string     `json:"-"`                  // 密码哈希（不返回）
    RealName         string     `json:"real_name"`          // 真实姓名
    Status           string     `json:"status"`              // ACTIVE/SUSPENDED
    LastLoginAt      *time.Time `json:"last_login_at"`
}
```

**职责**:
- 定义用户数据结构
- 支持邮箱和手机号登录
- 支持国际化手机号
- 记录最后登录信息

#### 2. JWT认证中间件

**文件路径**: `internal/middleware/auth.go`

**代码行数**: ~100行

**关键函数**:
```go
// 生成JWT Token
func GenerateToken(userID int64, username string, tenantID int64, expireTime int) (string, error)

// 解析JWT Token
func ParseToken(tokenString string) (*Claims, error)

// JWT认证中间件
func Auth() gin.HandlerFunc

// 从上下文获取用户ID
func GetUserID(c *gin.Context) int64

// 从上下文获取租户ID
func GetTenantID(c *gin.Context) int64
```

**职责**:
- JWT Token生成和解析
- 认证中间件
- 用户信息上下文管理

---

## FE-007-03: RBAC权限模型

### 功能需求

**描述**: 三级权限架构（系统角色 → 功能权限 → 操作权限）

**关联文档**: FE-007-03.md

### 代码实现

#### 1. Role数据模型

**文件路径**: `internal/models/user.go`

**代码行数**: ~30行

**关键字段**:
```go
type Role struct {
    ID          int64  `json:"id"`
    TenantID    int64  `json:"tenant_id"`
    RoleCode    string `json:"role_code"`    // SYSTEM_ADMIN/ORG_ADMIN/NORMAL_USER
    RoleName    string `json:"role_name"`
    Description string `json:"description"`
    IsSystem    bool   `json:"is_system"`     // 是否系统角色
    IsDeletable bool   `json:"is_deletable"` // 是否可删除
}
```

#### 2. Permission数据模型

**文件路径**: `internal/models/permission.go`

**代码行数**: ~60行

**关键常量**:
```go
// 功能权限
const (
    FeatureSystemConfig       = "SYSTEM_CONFIG"
    FeatureOrganizationMgmt   = "ORGANIZATION_MANAGEMENT"
    FeatureUserMgmt           = "USER_MANAGEMENT"
    FeatureRoleMgmt           = "ROLE_MANAGEMENT"
    FeatureDeviceMgmt         = "DEVICE_MANAGEMENT"
    FeatureDataView           = "DATA_VIEW"
    FeatureAlertMgmt          = "ALERT_MANAGEMENT"
    FeatureQuotaMgmt          = "QUOTA_MANAGEMENT"
    FeatureAuditLogView       = "AUDIT_LOG_VIEW"
)

// 操作权限
const (
    ActionView   = "VIEW"
    ActionCreate = "CREATE"
    ActionEdit   = "EDIT"
    ActionDelete = "DELETE"
    ActionExport = "EXPORT"
    ActionImport = "IMPORT"
)
```

**职责**:
- 定义9个功能模块权限
- 定义6个操作权限
- 三级权限架构基础

---

## FE-007-04: 租户数据隔离

### 功能需求

**描述**: 双字段数据隔离（tenant_id + managed_tenant_id）

**关联文档**: FE-007-04.md

### 代码实现

#### 1. 租户隔离中间件

**文件路径**: `internal/middleware/tenant.go`

**代码行数**: ~50行

**关键函数**:
```go
// 租户隔离Scope
func TenantScope(db *gorm.DB) func(*gorm.DB) *gorm.DB

// 租户隔离中间件
func TenantIsolation() gin.HandlerFunc

// 判断是否为集成商
func IsIntegrator(c *gin.Context) bool

// 判断是否为下游客户
func IsTerminal(c *gin.Context) bool
```

**职责**:
- 自动应用租户隔离查询Scope
- 集成商可查看所有下游数据
- 下游客户仅查看自身数据

---

## 代码统计

### 文件总数

| 类别 | 数量 |
|------|------|
| **配置文件** | 3 (config.yaml, Dockerfile, docker-compose.yml) |
| **基础设施** | 7 (config, logger, database, response) |
| **中间件** | 3 (auth, cors, tenant) |
| **数据模型** | 3 (tenant, user, permission) |
| **主入口** | 1 (main.go) |
| **脚本** | 2 (init.sql, Makefile) |
| **文档** | 1 (README.md) |
| **总计** | **20** |

### 代码行数

| 类别 | 代码行数 | 说明 |
|------|---------|------|
| Go (后端) | ~1200行 | 模型、中间件、工具 |
| SQL | ~50行 | 数据库脚本 |
| YAML | ~50行 | 配置文件 |
| Markdown | ~200行 | README.md |
| **总计** | **~1500行** | 含注释和空行 |

---

## 待实现功能映射

| 功能模块 | 依赖文件 | 状态 |
|---------|---------|------|
| 注册/登录API | internal/auth/ | 📋 待实现 (Phase 2) |
| 组织管理API | internal/tenant/ | 📋 待实现 (Phase 2) |
| 用户管理API | internal/user/ | 📋 待实现 (Phase 2) |
| 角色权限API | internal/role/, internal/permission/ | 📋 待实现 (Phase 3) |
| 权限验证中间件 | internal/middleware/permission.go | 📋 待实现 (Phase 3) |
| 审计日志中间件 | internal/middleware/audit.go | 📋 待实现 (Phase 4) |
| 配额管理中间件 | internal/middleware/quota.go | 📋 待实现 (Phase 4) |
| 前端项目 | platforms/cloud/frontend/ | 📋 待实现 (Phase 1.4) |

---

## 变更历史

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| 1.0 | 2026-01-28 | 初始版本，Phase 1基础架构映射 | Claude Code |

---

# Feature-to-Code Mapping: 移动端项目初始化

## 映射表信息

**文档ID**: feature-to-code-map (APP)
**文档标题**: 功能到代码映射表
**关联需求**: REQ-008
**关联功能需求**: FE-008
**关联实施计划**: IMP-008
**创建日期**: 2026-01-28
**目标平台**: APP (移动端)
**状态**: ⏳ 待实现

## 概述

移动端应用 (UniApp + Vue 3 + TypeScript) 的功能到代码映射，包括项目初始化、基础框架、核心模块、通用组件、工具类等。

---

## FE-008: 移动端项目初始化

### FE-008-01: 项目初始化 (4小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 项目创建 | `package.json` | - | - | 1-50 | ⏳ | 依赖配置 |
| 目录结构 | `src/` | - | - | - | ⏳ | 所有目录 |
| TS 配置 | `tsconfig.json` | - | - | 1-50 | ⏳ | Strict 模式 |
| Vite 配置 | `vite.config.ts` | - | - | 1-80 | ⏳ | 路径别名 |
| ESLint | `.eslintrc.js` | - | - | 1-50 | ⏳ | 代码规范 |
| Prettier | `.prettierrc` | - | - | 1-30 | ⏳ | 格式化 |

### FE-008-02: 基础页面框架 (3小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| App 入口 | `src/App.vue` | App | onLaunch | 1-50 | ⏳ | 生命周期 |
| 路由配置 | `src/pages.json` | - | - | 1-150 | ⏳ | TabBar 配置 |
| 应用配置 | `src/manifest.json` | - | - | 1-100 | ⏳ | UniApp 配置 |
| 全局样式 | `src/styles/index.scss` | - | - | 1-20 | ⏳ | 样式入口 |
| 主题变量 | `src/styles/variables.scss` | - | - | 1-100 | ⏳ | 颜色、间距 |
| 样式混入 | `src/styles/mixins.scss` | - | - | 1-50 | ⏳ | 常用混入 |
| 样式重置 | `src/styles/reset.scss` | - | - | 1-30 | ⏳ | CSS 重置 |

### FE-008-03: 核心模块骨架 (8.5小时)

#### Auth 模块 (2小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 登录页 | `src/pages/auth/login/index.vue` | LoginPage | handleLogin | 1-200 | ⏳ | 表单验证 |
| 注册页 | `src/pages/auth/register/index.vue` | RegisterPage | handleRegister | 1-150 | ⏳ | 注册表单 |
| Auth Store | `src/stores/auth.ts` | AuthStore | login/logout | 1-100 | ⏳ | 状态管理 |
| Auth API | `src/api/auth.ts` | authApi | login | 1-100 | ⏳ | Mock 实现 |
| useAuth | `src/composables/useAuth.ts` | useAuth | login | 1-80 | ⏳ | 业务逻辑 |
| Auth 类型 | `src/types/auth.d.ts` | - | - | 1-100 | ⏳ | TypeScript |

#### Device 模块 (1.5小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 设备列表 | `src/pages/device/list/index.vue` | DeviceList | onLoad | 1-100 | ⏳ | 列表展示 |
| 设备详情 | `src/pages/device/detail/index.vue` | DeviceDetail | loadDetail | 1-150 | ⏳ | 详情展示 |
| Device Store | `src/stores/device.ts` | DeviceStore | fetchDevices | 1-100 | ⏳ | 状态管理 |
| Device API | `src/api/device.ts` | deviceApi | getList | 1-100 | ⏳ | API 封装 |
| useDevice | `src/composables/useDevice.ts` | useDevice | fetchList | 1-80 | ⏳ | 业务逻辑 |
| Device 类型 | `src/types/device.d.ts` | - | - | 1-80 | ⏳ | TypeScript |

#### Workspace 模块 (1小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 工作台页 | `src/pages/workspace/index.vue` | WorkspacePage | onLoad | 1-80 | ⏳ | 主页 |
| Workspace Store | `src/stores/workspace.ts` | WorkspaceStore | - | 1-80 | ⏳ | 状态管理 |
| Workspace API | `src/api/workspace.ts` | workspaceApi | - | 1-80 | ⏳ | API 封装 |
| Workspace 类型 | `src/types/workspace.d.ts` | - | - | 1-60 | ⏳ | TypeScript |

#### Dashboard 模块 (1小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 看板页 | `src/pages/dashboard/index.vue` | DashboardPage | onLoad | 1-80 | ⏳ | 主页 |
| Dashboard Store | `src/stores/dashboard.ts` | DashboardStore | - | 1-80 | ⏳ | 状态管理 |
| Dashboard API | `src/api/dashboard.ts` | dashboardApi | - | 1-80 | ⏳ | API 封装 |
| Dashboard 类型 | `src/types/dashboard.d.ts` | - | - | 1-60 | ⏳ | TypeScript |

#### Message 模块 (1小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 消息列表 | `src/pages/message/list/index.vue` | MessageList | onLoad | 1-80 | ⏳ | 列表展示 |
| 消息详情 | `src/pages/message/detail/index.vue` | MessageDetail | loadDetail | 1-100 | ⏳ | 详情展示 |
| Message Store | `src/stores/message.ts` | MessageStore | fetchMessages | 1-80 | ⏳ | 状态管理 |
| Message API | `src/api/message.ts` | messageApi | getList | 1-80 | ⏳ | API 封装 |
| Message 类型 | `src/types/message.d.ts` | - | - | 1-60 | ⏳ | TypeScript |

#### Profile 模块 (1小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 个人中心 | `src/pages/profile/center/index.vue` | ProfileCenter | onLoad | 1-80 | ⏳ | 主页 |
| 设置页 | `src/pages/profile/settings/index.vue` | SettingsPage | - | 1-100 | ⏳ | 设置 |
| Profile Store | `src/stores/profile.ts` | ProfileStore | - | 1-80 | ⏳ | 状态管理 |
| Profile API | `src/api/profile.ts` | profileApi | - | 1-80 | ⏳ | API 封装 |
| Profile 类型 | `src/types/profile.d.ts` | - | - | 1-60 | ⏳ | TypeScript |

#### App 全局 Store (0.5小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| App Store | `src/stores/app.ts` | AppStore | setTheme | 1-60 | ⏳ | 全局状态 |

### FE-008-04: 通用组件库 (4.5小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| CustomNavBar | `src/components/common/CustomNavBar/index.vue` | CustomNavBar | - | 1-120 | ⏳ | 自定义导航栏 |
| PageContainer | `src/components/common/PageContainer/index.vue` | PageContainer | - | 1-60 | ⏳ | 页面容器 |
| Loading | `src/components/common/Loading/index.vue` | Loading | - | 1-80 | ⏳ | 加载指示器 |
| EmptyState | `src/components/common/EmptyState/index.vue` | EmptyState | - | 1-80 | ⏳ | 空状态提示 |
| NetworkError | `src/components/common/NetworkError/index.vue` | NetworkError | - | 1-80 | ⏳ | 网络异常 |
| PullRefresh | `src/components/common/PullRefresh/index.vue` | PullRefresh | - | 1-100 | ⏳ | 下拉刷新 |
| LoadMore | `src/components/common/LoadMore/index.vue` | LoadMore | - | 1-100 | ⏳ | 上拉加载 |

### FE-008-05: 工具类和类型定义 (3.5小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| HTTP 封装 | `src/utils/request.ts` | request | request | 1-200 | ⏳ | 请求封装 |
| Storage 封装 | `src/utils/storage.ts` | storage | set/get | 1-80 | ⏳ | 存储封装 |
| 验证工具 | `src/utils/validator.ts` | validator | validate | 1-100 | ⏳ | 表单验证 |
| 格式化工具 | `src/utils/format.ts` | format | formatDate | 1-80 | ⏳ | 格式化 |
| 常量定义 | `src/utils/constants.ts` | - | - | 1-50 | ⏳ | 常量 |
| API 类型 | `src/types/api.d.ts` | - | - | 1-50 | ⏳ | API 响应 |
| 通用类型 | `src/types/common.d.ts` | - | - | 1-100 | ⏳ | 通用类型 |

### FE-008-06: 第一个页面实现 (4小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 启动页 | `src/pages/index/index.vue` | IndexPage | onLaunch | 1-100 | ⏳ | Logo 展示 |
| 登录页 | `src/pages/auth/login/index.vue` | LoginPage | handleLogin | 1-200 | ⏳ | Mock 登录 |

### FE-008-07: 开发规范文档 (1小时)

| 功能点 | 文件路径 | 类/组件 | 函数/方法 | 行号范围 | 状态 | 备注 |
|--------|---------|---------|----------|---------|------|------|
| 开发规范 | `docs/development-guide.md` | - | - | 1-500 | ⏳ | 规范文档 |

---

## 统计信息

**总功能点**: 7 个 (FE-008-01 ~ FE-008-07)
**总文件数**: 50+ 个
**总代码行数**: 预估 5000+ 行
**当前进度**: 0% (全部待实现)

---

## 依赖关系

### 弱依赖

**FE-007**: 云平台账号系统 (多租户认证 API)

**依赖处理**: 使用 Mock 实现

**Mock 位置**: `src/api/auth.ts`, `src/utils/mock.ts`

**补齐优先级**: P0

**补齐时机**: IMP-007 完成后

---

## 更新历史

| 日期 | 版本 | 变更内容 | 变更人 |
|------|------|---------|--------|
| 2026-01-28 | 1.0 | 初始创建,基于 IMP-008 | Claude Code |
