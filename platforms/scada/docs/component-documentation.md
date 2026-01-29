# Scada 工程管理 - 组件文档

## 概述

本文档描述 Scada 工程管理功能的所有 React 组件，包括其 props、状态、使用方法和示例。

---

## 核心组件

### 1. NewProjectDialog

新建工程对话框组件，提供完整的工程创建表单。

**文件路径**: `src/components/project/NewProjectDialog.tsx`

**Props**:
```typescript
interface NewProjectDialogProps {
  onClose: () => void              // 关闭对话框回调
  onProjectCreated: (projectId: string) => void  // 工程创建成功回调
}
```

**状态管理**:
- 使用 MobX observer 包装，响应式更新
- 表单数据本地状态（formData）
- 验证错误本地状态（errors）
- 加载状态（isSubmitting）

**主要功能**:
- 工程基本信息输入（名称、作者、描述）
- 分类选择（预定义 + 自定义）
- 硬件平台选择
- 保存位置选择（集成 Electron 文件对话框）
- 工程加密选项
- 密码输入和强度指示
- 表单验证
- ESC 键关闭支持

**使用示例**:
```typescript
import { useState } from 'react'
import NewProjectDialog from '@/components/project/NewProjectDialog'

function App() {
  const [showDialog, setShowDialog] = useState(false)

  const handleProjectCreated = (projectId: string) => {
    console.log('工程创建成功:', projectId)
    setShowDialog(false)
    // 刷新最近工程列表
  }

  return (
    <div>
      <button onClick={() => setShowDialog(true)}>新建工程</button>
      {showDialog && (
        <NewProjectDialog
          onClose={() => setShowDialog(false)}
          onProjectCreated={handleProjectCreated}
        />
      )}
    </div>
  )
}
```

**样式文件**: `NewProjectDialog.css`

**依赖**:
- `PasswordStrengthIndicator`
- `@/api/projectApi`
- `@/utils/electron` (for `getElectronAPI`)

---

### 2. PasswordDialog

密码输入对话框，用于打开加密工程时输入密码。

**文件路径**: `src/components/project/PasswordDialog.tsx`

**Props**:
```typescript
interface PasswordDialogProps {
  isOpen: boolean                  // 是否显示对话框
  projectName: string              // 工程名称
  onSubmit: (password: string) => Promise<void>  // 提交密码回调
  onCancel: () => void             // 取消回调
}
```

**主要功能**:
- 密码输入（带可见性切换）
- 密码强度显示
- 记住密码选项（预留）
- "忘记密码?"链接（预留）
- Enter 键提交支持

**使用示例**:
```typescript
import { useState } from 'react'
import PasswordDialog from '@/components/project/PasswordDialog'
import { projectApi } from '@/api/projectApi'

function App() {
  const [showPasswordDialog, setShowPasswordDialog] = useState(false)
  const [currentProject, setCurrentProject] = useState('')

  const handlePasswordSubmit = async (password: string) => {
    try {
      const response = await projectApi.openProject({
        filePath: currentProject,
        password
      })
      if (response.success) {
        setShowPasswordDialog(false)
        // 加载工程
      }
    } catch (error) {
      console.error('打开工程失败:', error)
    }
  }

  return (
    <div>
      <button onClick={() => setShowPasswordDialog(true)}>打开加密工程</button>
      <PasswordDialog
        isOpen={showPasswordDialog}
        projectName="加密工程.pant"
        onSubmit={handlePasswordSubmit}
        onCancel={() => setShowPasswordDialog(false)}
      />
    </div>
  )
}
```

---

### 3. PasswordStrengthIndicator

密码强度指示器组件，可视化显示密码安全等级。

**文件路径**: `src/components/project/PasswordStrengthIndicator.tsx`

**Props**:
```typescript
interface PasswordStrengthIndicatorProps {
  strength: PasswordStrength        // 密码强度等级
}

type PasswordStrength = 'weak' | 'medium' | 'strong'
```

**强度配置**:

| 等级 | 颜色 | 宽度 | 提示 |
|------|------|------|------|
| weak | #f44336 (红) | 33% | 弱: 建议12+字符，混合大小写、数字、符号 |
| medium | #ff9800 (橙) | 66% | 中: 可以更强，添加更多字符类型 |
| strong | #4caf50 (绿) | 100% | 强: 安全的密码 |

**使用示例**:
```typescript
import { useState } from 'react'
import PasswordStrengthIndicator from '@/components/project/PasswordStrengthIndicator'

function PasswordInput() {
  const [password, setPassword] = useState('')

  const calculateStrength = (): PasswordStrength => {
    if (password.length < 6) return 'weak'
    if (password.length < 12) return 'medium'
    return 'strong'
  }

  return (
    <div>
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="请输入密码"
      />
      {password && (
        <PasswordStrengthIndicator strength={calculateStrength()} />
      )}
    </div>
  )
}
```

---

### 4. RecentProjectsList

最近工程列表组件，使用虚拟滚动优化大列表性能。

**文件路径**: `src/components/workspace/RecentProjectsList.tsx`

**Props**:
```typescript
interface RecentProjectsListProps {
  projects: RecentProject[]        // 工程列表
  onProjectClick: (project: RecentProject) => void  // 点击工程回调
  onProjectRemove: (projectId: string) => void      // 删除工程回调
  emptyMessage?: string            // 空状态提示
}
```

**主要功能**:
- 虚拟滚动（react-window）
- 固定高度项目（80px）
- 自动计算列表高度（最大400px）
- 空状态提示
- 键盘导航支持（Tab）

**性能特性**:
- 只渲染可见项目
- 支持1000+工程流畅滚动
- 滚动位置保持（使用 ref）

**使用示例**:
```typescript
import { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import RecentProjectsList from '@/components/workspace/RecentProjectsList'
import { recentProjectsStore } from '@/store'

const ProjectList = observer(() => {
  useEffect(() => {
    recentProjectsStore.loadRecentProjects()
  }, [])

  const handleProjectClick = (project) => {
    console.log('打开工程:', project.name)
    // 打开工程逻辑
  }

  const handleProjectRemove = async (projectId) => {
    await recentProjectsStore.removeProject(projectId)
  }

  return (
    <RecentProjectsList
      projects={recentProjectsStore.displayProjects}
      onProjectClick={handleProjectClick}
      onProjectRemove={handleProjectRemove}
      emptyMessage="暂无最近工程"
    />
  )
})
```

**依赖**:
- `react-window` (List component)
- `ProjectListItem`

---

### 5. ProjectListItem

单个工程列表项组件，显示工程信息和上下文菜单。

**文件路径**: `src/components/workspace/ProjectListItem.tsx`

**Props**:
```typescript
interface ProjectListItemProps {
  project: RecentProject           // 工程数据
  onClick: (project: RecentProject) => void   // 点击回调
  onRemove: (projectId: string) => void       // 删除回调
}
```

**主要功能**:
- 显示工程图标、名称、时间
- 右键上下文菜单
- 悬停效果
- 键盘导航（Enter 打开，Delete 删除）
- 加密工程标识（🔒 图标）

**上下文菜单**:
- 打开工程
- 从列表中移除
- 在文件管理器中显示

**使用示例**:
```typescript
import ProjectListItem from '@/components/workspace/ProjectListItem'

function ProjectList({ projects }) {
  const handleProjectClick = (project) => {
    // 打开工程
  }

  const handleProjectRemove = (projectId) => {
    // 删除工程
  }

  return (
    <div>
      {projects.map(project => (
        <ProjectListItem
          key={project.projectId}
          project={project}
          onClick={handleProjectClick}
          onRemove={handleProjectRemove}
        />
      ))}
    </div>
  )
}
```

---

### 6. CategoryFilter

分类筛选组件，支持横向滚动的分类标签。

**文件路径**: `src/components/workspace/CategoryFilter.tsx`

**Props**:
```typescript
interface CategoryFilterProps {
  categories: Array<{             // 分类列表
    name: string                  // 显示名称
    value: string                 // 值
    count: number                 // 工程数量
  }>
  selectedCategory: string        // 当前选中的分类
  onCategoryChange: (category: string) => void  // 分类变更回调
}
```

**主要功能**:
- 横向滚动（分类较多时）
- 数量徽章显示
- 选中状态高亮
- 点击切换分类

**样式特性**:
- Flexbox 布局
- 隐藏滚动条但保持可滚动
- 平滑过渡动画

**使用示例**:
```typescript
import { observer } from 'mobx-react-lite'
import CategoryFilter from '@/components/workspace/CategoryFilter'
import { recentProjectsStore } from '@/store'

const CategoryFilterWrapper = observer(() => {
  return (
    <CategoryFilter
      categories={recentProjectsStore.categories}
      selectedCategory={recentProjectsStore.selectedCategory}
      onCategoryChange={(category) => recentProjectsStore.setCategory(category)}
    />
  )
})
```

---

### 7. SearchBox

搜索框组件，支持防抖和键盘快捷键。

**文件路径**: `src/components/workspace/SearchBox.tsx`

**Props**:
```typescript
interface SearchBoxProps {
  onSearch: (query: string) => void  // 搜索回调
  debounceMs?: number                // 防抖延迟（默认100ms）
  placeholder?: string               // 占位符文本
}
```

**主要功能**:
- 实时搜索输入
- 防抖优化（默认100ms）
- 清除按钮
- 键盘快捷键（Ctrl/Cmd + F 聚焦）
- ESC 清空

**使用示例**:
```typescript
import { observer } from 'mobx-react-lite'
import SearchBox from '@/components/workspace/SearchBox'
import { recentProjectsStore } from '@/store'

const SearchBoxWrapper = observer(() => {
  const handleSearch = (query) => {
    recentProjectsStore.setSearchQuery(query)
  }

  return (
    <SearchBox
      onSearch={handleSearch}
      placeholder="搜索工程名称..."
      debounceMs={100}
    />
  )
})
```

---

## MobX Stores

### 1. ProjectStore

当前工程状态管理。

**文件路径**: `src/store/projectStore.ts`

**状态**:
```typescript
class ProjectStore {
  currentProject: Project | null     // 当前工程
  status: ProjectStatus              // 操作状态
  error: string                      // 错误信息
}

type ProjectStatus = 'idle' | 'creating' | 'opening' | 'saving' | 'error'
```

**Computed**:
- `hasProject`: 是否有打开的工程
- `projectName`: 工程名称
- `isEncrypted`: 是否加密工程
- `isLoading`: 是否正在加载
- `hasError`: 是否有错误

**Actions**:
- `createProject(formData)`: 创建新工程
- `openProject(request)`: 打开工程
- `saveProject()`: 保存当前工程
- `closeProject()`: 关闭工程
- `clearError()`: 清除错误

**使用示例**:
```typescript
import { observer } from 'mobx-react-lite'
import { projectStore } from '@/store'

const ProjectHeader = observer(() => {
  if (!projectStore.hasProject) {
    return <div>未打开工程</div>
  }

  return (
    <div>
      <h1>{projectStore.projectName}</h1>
      {projectStore.isEncrypted && <span>🔒 加密工程</span>}
      {projectStore.isLoading && <span>加载中...</span>}
      {projectStore.hasError && <div>{projectStore.error}</div>}
    </div>
  )
})
```

---

### 2. RecentProjectsStore

最近工程列表状态管理。

**文件路径**: `src/store/recentProjectsStore.ts`

**状态**:
```typescript
class RecentProjectsStore {
  projects: RecentProject[]         // 原始数据
  selectedCategory: string          // 选中的分类
  searchQuery: string               // 搜索关键词
  sortBy: SortBy                    // 排序字段
  sortOrder: SortOrder              // 排序方向
  isLoading: boolean                // 加载状态
  error: string                     // 错误信息
}

type SortBy = 'lastOpened' | 'name' | 'createdAt'
type SortOrder = 'asc' | 'desc'
```

**Computed**:
- `categories`: 所有分类（含数量统计）
- `filteredProjects`: 过滤后的工程列表
- `sortedProjects`: 排序后的工程列表
- `displayProjects`: 显示用的工程列表（带格式化时间）
- `totalCount`: 工程总数
- `filteredCount`: 筛选后的数量

**Actions**:
- `loadRecentProjects()`: 加载最近工程列表
- `setCategory(category)`: 设置分类筛选
- `setSearchQuery(query)`: 设置搜索关键词
- `setSortBy(sortBy)`: 设置排序字段
- `toggleSortOrder()`: 切换排序方向
- `removeProject(projectId)`: 移除工程
- `clear()`: 清空列表

**使用示例**:
```typescript
import { observer } from 'mobx-react-lite'
import { recentProjectsStore } from '@/store'

const RecentProjects = observer(() => {
  useEffect(() => {
    recentProjectsStore.loadRecentProjects()
  }, [])

  return (
    <div>
      <div>总计: {recentProjectsStore.totalCount}</div>
      <div>筛选: {recentProjectsStore.filteredCount}</div>

      <CategoryFilter
        categories={recentProjectsStore.categories}
        selectedCategory={recentProjectsStore.selectedCategory}
        onCategoryChange={(cat) => recentProjectsStore.setCategory(cat)}
      />

      <SearchBox
        onSearch={(query) => recentProjectsStore.setSearchQuery(query)}
      />

      <RecentProjectsList
        projects={recentProjectsStore.displayProjects}
        onProjectClick={(p) => openProject(p.filePath)}
        onProjectRemove={(id) => recentProjectsStore.removeProject(id)}
      />
    </div>
  )
})
```

---

## 工具函数

### 1. Electron 工具

**文件路径**: `src/utils/electron.ts`

**函数**:

#### `isElectron()`
检测是否运行在 Electron 环境。

```typescript
const isElectron = (): boolean
```

**返回**: `true` 如果在 Electron 中运行

#### `getElectronAPI()`
获取 Electron API 或 Mock 实现。

```typescript
const getElectronAPI = (): ElectronAPI
```

**返回**:
- Electron 环境: 真实的 `window.electronAPI`
- 浏览器环境: Mock 实现

#### `getApiBaseUrl()`
获取后端 API 基础 URL。

```typescript
const getApiBaseUrl = (): string
```

**返回**:
- Electron: `http://localhost:3000`
- 浏览器: `/api` (Vite 代理)

---

### 2. 日期格式化

**文件路径**: `src/utils/dateFormat.ts`

**函数**:

#### `formatRelativeTime()`
格式化相对时间（如 "2小时前"）。

```typescript
const formatRelativeTime = (date: Date | string): string
```

**示例**:
```typescript
formatRelativeTime(new Date()) // "刚刚"
formatRelativeTime(new Date(Date.now() - 3600000)) // "1小时前"
formatRelativeTime(new Date(Date.now() - 86400000)) // "1天前"
formatRelativeTime("2026-01-20T10:00:00Z") // "X天前"
```

---

## 样式约定

### CSS 变量

所有组件使用统一的 CSS 变量系统：

```css
:root {
  /* Colors */
  --color-bg-primary: #f5f5f5;
  --color-bg-secondary: #ffffff;
  --color-accent-active: #FF9999;
  --color-text-primary: #333333;
  --color-border: #d0d0d0;

  /* Spacing */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;

  /* Typography */
  --font-size-sm: 12px;
  --font-size-md: 14px;
  --font-size-lg: 16px;

  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 250ms ease;
}
```

### BEM 命名

所有 CSS 类名使用 BEM (Block Element Modifier) 约定：

```css
.block { }
.block__element { }
.block__element--modifier { }
```

**示例**:
```css
.new-project-dialog { }
.new-project-dialog__header { }
.new-project-dialog__title { }
.new-project-dialog__input--error { }
```

---

## 性能优化

### 1. 虚拟滚动

`RecentProjectsList` 使用 `react-window` 实现虚拟滚动：

- 只渲染可见项目
- 固定项目高度（80px）
- 支持1000+工程流畅滚动
- 自动计算列表高度

### 2. MobX 优化

- 使用 `observer` 包装组件，精确追踪依赖
- Computed 属性缓存计算结果
- `runInAction` 批量更新状态

### 3. 防抖

`SearchBox` 使用100ms防抖：

- 减少搜索请求次数
- 提升输入体验
- 可配置延迟时间

### 4. 懒加载

对话框组件按需加载：

```typescript
const NewProjectDialog = lazy(() => import('@/components/project/NewProjectDialog'))
```

---

## 测试指南

### 单元测试示例

```typescript
import { render, screen } from '@testing-library/react'
import NewProjectDialog from '@/components/project/NewProjectDialog'

describe('NewProjectDialog', () => {
  it('should render dialog title', () => {
    render(
      <NewProjectDialog
        onClose={() => {}}
        onProjectCreated={() => {}}
      />
    )
    expect(screen.getByText('新建工程')).toBeInTheDocument()
  })

  it('should call onClose on cancel', () => {
    const handleClose = jest.fn()
    render(
      <NewProjectDialog
        onClose={handleClose}
        onProjectCreated={() => {}}
      />
    )
    fireEvent.click(screen.getByText('取消'))
    expect(handleClose).toHaveBeenCalled()
  })
})
```

### 集成测试示例

```typescript
import { render, screen, waitFor } from '@testing-library/react'
import { projectApi } from '@/api/projectApi'

jest.mock('@/api/projectApi')

describe('Create Project Flow', () => {
  it('should create project successfully', async () => {
    projectApi.createProject.mockResolvedValue({
      success: true,
      data: { projectId: 'test-id', filePath: 'test.pant' }
    })

    // ... render and interact with component

    await waitFor(() => {
      expect(projectApi.createProject).toHaveBeenCalled()
    })
  })
})
```

---

## 可访问性

### ARIA 属性

所有交互组件包含适当的 ARIA 属性：

```typescript
<button
  aria-label="关闭"
  onClick={onClose}
>
  ✕
</button>

<div
  role="button"
  tabIndex={0}
  onKeyDown={(e) => {
    if (e.key === 'Enter') onClick()
  }}
>
  点击区域
</div>
```

### 键盘导航

- `Tab`: 焦点移动
- `Enter` / `Space`: 激活按钮
- `Escape`: 关闭对话框
- `Arrow Keys`: 列表导航

---

## 故障排查

### 常见问题

**问题**: `window.electronAPI is undefined`

**原因**: 不在 Electron 环境中运行

**解决**: 使用 `getElectronAPI()` 获取 API，会自动回退到 Mock 实现

**问题**: MobX 状态更新不触发重渲染

**原因**: 没有使用 `observer()` 包装组件

**解决**: 确保组件使用 `observer()` 包装

**问题**: 虚拟滚动列表不滚动

**原因**: 父容器高度未设置或 `overflow: hidden`

**解决**: 设置父容器固定高度和 `overflow: hidden`

---

## 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0.0 | 2026-01-21 | 初始版本，实现所有核心组件 |

---

## 支持

如有组件使用问题，请参考：
- API 文档: `api-documentation.md`
- 集成测试指南: `integration-test-guide.md`
- 数据格式文档: `data-format-documentation.md`
