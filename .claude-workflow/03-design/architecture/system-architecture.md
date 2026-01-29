# Scada工程编辑器系统架构

## 概述

本文档描述Scada工程编辑器（EditorLayout）的系统架构，包括组件层次、数据流、状态管理和技术栈。

## 整体架构图

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           Scada Desktop Application                       │
│  (Electron + React + MobX)                                                      │
├──────────────────────────────────────────────────────────────────────────┤
│  TitleBar (FE-005)                                                              │
│  - Logo + Title                                                                  │
│  - WindowControls                                                               │
│  - ProjectTabs ← NEW (FE-006-01)                                               │
├──────────────────────────────────────────────────────────────────────────┤
│  Sidebar (FE-001)                                                               │
│  - NavItem: 首页, 本地, 云端, 工具, User                                        │
├──────────────────────────────────────────────────────────────────────────┤
│  MainContent (FE-001)                                                          │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │ EditorLayout (FE-006-02) - NEW                                        │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────┐ │ │
│  │  │ TopToolbar (FE-006-03) - 32px                                  │ │ │
│  │  │ 文件 | 编辑 | 视图 | 工具箱                                           │ │ │
│  │  └────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────┐ │ │
│  │  │ SubToolbar (FE-006-04) - 64px                                  │ │ │
│  │  │ 上下文相关工具 (属性/颜色/对齐等)                                 │ │ │
│  │  └────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                          │ │
│  │  ┌───────────┬─────────────────────────┬─────────────────────────┐ │ │
│  │  │           │  SubPageTabs (FE-006-12)   │                         │ │ │
│  │  │ LeftSidebar│  - 36px                  │  RightPanel (FE-006-09) │ │ │
│  │  │ (FE-006-05)├─────────────────────────┤  - 280px (可隐藏)      │ │ │
│  │  │           │                         │  - 属性 | 图层          │ │ │
│  │  │ - 工程Tab │                         │  ┌───────────────────┐ │ │
│  │  │   ↓       │                         │  │ PropertiesPanel  │ │ │
│  │  │ ┌────────┐ │                         │  │ (FE-006-10)      │ │ │
│  │  │ │Project │ │                         │  │ - 基本属性        │ │ │
│  │  │ │Panel  │ │                         │  │ - 外观属性        │ │ │
│  │  │ │(006-06)│ │                         │  │ - 数据属性        │ │ │
│  │  │ └────────┘ │                         │  └───────────────────┘ │ │
│  │  │           │                         │  ┌───────────────────┐ │ │
│  │  │ - 画面Tab │                         │  │ LayersPanel       │ │ │
│  │  │   ↓       │                         │  │ (FE-006-11)       │ │ │
│  │  │ ┌────────┐ │                         │  │ - 图层列表        │ │ │
│  │  │ │Pages   │ │                         │  │ - 锁定/隐藏       │ │ │
│  │  │ │Panel  │ │                         │  └───────────────────┘ │ │
│  │  │ │(006-07)│ │                         │                         │ │
│  │  │ └────────┘ │                         │                         │ │
│  │  │           │                         │                         │ │
│  │  │ - 组件Tab │                         │                         │ │
│  │  │   ↓       │                         │                         │ │
│  │  │ ┌────────┐ │                         │                         │ │
│  │  │ │Compo   │ │                         │                         │ │
│  │  │ │nents  │ │                         │                         │ │
│  │  │ │Panel  │ │                         │                         │ │
│  │  │ │(006-08)│ │                         │                         │ │
│  │  │ └────────┘ │                         │                         │ │
│  │  └───────────┴─────────────────────────┴─────────────────────────┘ │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────┐ │ │
│  │  │ CanvasArea (FE-006-13)                                         │ │ │
│  │  │                                                                 │ │ │
│  │  │  ┌─────────────────────────────────────────────────────────┐  │ │ │
│  │  │  │ Canvas - 网格背景                                         │  │ │ │
│  │  │  │                                                         │  │ │ │
│  │  │  │  [组件1] [组件2] [组件3] ...                            │  │ │ │
│  │  │  │                                                         │  │ │ │
│  │  │  └─────────────────────────────────────────────────────────┘  │ │ │
│  │  └────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────┐ │ │
│  │  │ StatusBar (FE-006-14) - 24px                                  │ │ │
│  │  │ 状态 | 缩放 | 坐标 | 通知                                          │ │ │
│  │  └────────────────────────────────────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────┘
```

## 组件层次结构

```
EditorLayout (FE-006-02)
│
├── TopToolbar (FE-006-03) - 32px高度
│   ├── ToolbarGroup (文件操作)
│   │   ├── ToolbarButton (新建)
│   │   ├── ToolbarButton (打开)
│   │   ├── ToolbarButton (保存)
│   │   ├── ToolbarButton (另存为)
│   │   └── ToolbarButton (导出)
│   ├── ToolbarGroup (编辑操作)
│   │   ├── ToolbarButton (撤销)
│   │   ├── ToolbarButton (重做)
│   │   ├── ToolbarButton (剪切)
│   │   ├── ToolbarButton (复制)
│   │   ├── ToolbarButton (粘贴)
│   │   └── ToolbarButton (删除)
│   ├── ToolbarGroup (视图操作)
│   │   ├── ToolbarButton (放大)
│   │   ├── ToolbarButton (缩小)
│   │   ├── ToolbarButton (适应窗口)
│   │   └── ToolbarButton (实际大小)
│   └── ToolbarGroup (工具箱)
│       ├── ToolbarButton (选择)
│       ├── ToolbarButton (画线)
│       ├── ToolbarButton (画矩形)
│       ├── ToolbarButton (画圆形)
│       ├── ToolbarButton (文本)
│       └── ToolbarButton (图片)
│
├── SubToolbar (FE-006-04) - 64px高度
│   └── ContextualControls (根据选中对象动态渲染)
│       ├── ColorPicker (填充色)
│       ├── ColorPicker (边框色)
│       ├── Select (边框宽度)
│       ├── Select (线型)
│       └── ButtonGroup (对齐方式)
│
├── LeftSidebar (FE-006-05) - 280px宽度
│   ├── SidebarTab (工程 | 画面 | 组件)
│   │
│   ├── ProjectPanel (FE-006-06) - 当activeTab='工程'
│   │   └── ProjectTree
│   │       ├── TreeNode (工程名称)
│   │       ├── TreeNode (画面1)
│   │       │   ├── TreeNode (组件1)
│   │       │   └── TreeNode (组件2)
│   │       └── TreeNode (画面2)
│   │
│   ├── PagesPanel (FE-006-07) - 当activeTab='画面'
│   │   └── PageItem[]
│   │       ├── PageThumbnail (可选)
│   │       ├── PageName
│   │       └── PageStatus (已保存/未保存)
│   │
│   └── ComponentsPanel (FE-006-08) - 当activeTab='组件'
│       └── ComponentItem[] (网格布局)
│           ├── ComponentIcon
│           ├── ComponentName
│           └── DragHandle
│
├── SubPageTabs (FE-006-12) - 36px高度
│   └── SubPageTab[]
│       ├── TabTitle (画面名称)
│       ├── CloseButton
│       └── UnsavedIndicator (小圆点)
│
├── CanvasArea (FE-006-13) - 占据剩余空间
│   └── Canvas
│       ├── GridBackground
│       ├── CanvasObject[] (组件实例)
│       │   ├── Rectangle
│       │   ├── Circle
│       │   ├── Text
│       │   └── Image
│       └── SelectionOverlay (选中框)
│
├── RightPanel (FE-006-09) - 280px宽度，默认隐藏
│   ├── SidebarTab (属性 | 图层)
│   │
│   ├── PropertiesPanel (FE-006-10) - 当activeTab='属性'
│   │   ├── PropertyGroup (基本属性)
│   │   │   ├── PropertyField (名称)
│   │   │   ├── PropertyField (位置 X, Y)
│   │   │   ├── PropertyField (尺寸 W, H)
│   │   │   ├── PropertyField (可见性)
│   │   │   └── PropertyField (锁定)
│   │   ├── PropertyGroup (外观属性)
│   │   │   ├── PropertyField (背景色)
│   │   │   ├── PropertyField (边框色)
│   │   │   ├── PropertyField (边框宽度)
│   │   │   └── PropertyField (透明度)
│   │   ├── PropertyGroup (文本属性) - 仅文本组件
│   │   │   ├── PropertyField (字体)
│   │   │   ├── PropertyField (字号)
│   │   │   └── PropertyField (颜色)
│   │   ├── PropertyGroup (数据属性)
│   │   │   ├── PropertyField (数据绑定)
│   │   │   └── PropertyField (变量关联)
│   │   └── PropertyGroup (事件属性)
│   │       ├── PropertyField (点击事件)
│   │       └── PropertyField (数据变更事件)
│   │
│   └── LayersPanel (FE-006-11) - 当activeTab='图层'
│       └── LayerItem[]
│           ├── LayerVisibilityIcon (👁️)
│           ├── LayerLockIcon (🔒)
│           ├── LayerName
│           └── DragHandle
│
└── StatusBar (FE-006-14) - 24px高度
    ├── StatusInfo (就绪 | 操作中 | 错误)
    ├── ZoomControls
    │   ├── ZoomOutButton
    │   ├── ZoomLevelDisplay (100%)
    │   ├── ZoomInButton
    │   └── FitToWindowButton
    ├── CoordinateInfo
    │   ├── MousePosition (X: 123, Y: 456)
    │   └── SelectionSize (W: 200, H: 100)
    └── Notifications (未保存提示 | 错误 | 警告)
```

## MobX Store架构

```
MobX Root Store
│
├── editorStore (编辑器全局状态)
│   ├── State (状态)
│   │   ├── projectTabs: ObservableArray<ProjectTab>
│   │   ├── activeProjectTabId: string
│   │   ├── leftSidebarActiveTab: '工程' | '画面' | '组件'
│   │   ├── rightPanelVisible: boolean
│   │   ├── rightPanelActiveTab: '属性' | '图层'
│   │   ├── selectedObjectIds: ObservableArray<string>
│   │   └── zoom: number
│   │
│   ├── Actions (操作)
│   │   ├── openProject(project: Project): void
│   │   ├── closeProject(projectId: string): void
│   │   ├── switchProject(projectId: string): void
│   │   ├── setLeftSidebarActiveTab(tab: string): void
│   │   ├── setRightPanelVisible(visible: boolean): void
│   │   ├── setRightPanelActiveTab(tab: string): void
│   │   ├── setSelectedObjects(ids: string[]): void
│   │   └── setZoom(zoom: number): void
│   │
│   └── Computed (派生状态)
│       ├── activeProjectTab: ProjectTab | undefined
│       ├── hasUnsavedChanges: boolean
│       └── canUndo: boolean
│
├── projectStore (工程数据状态)
│   ├── State
│   │   ├── currentProject: Project | null
│   │   ├── pages: ObservableArray<Page>
│   │   └── components: ObservableArray<Component>
│   │
│   ├── Actions
│   │   ├── createProject(name: string): Project
│   │   ├── openProject(file: File): Promise<Project>
│   │   ├── saveProject(): Promise<void>
│   │   ├── createPage(name: string): Page
│   │   ├── deletePage(pageId: string): void
│   │   └── renamePage(pageId: string, newName: string): void
│   │
│   └── Computed
│       ├── currentPage: Page | undefined
│       └── pageCount: number
│
└── canvasStore (画布状态)
    ├── State
    │   ├── currentPageId: string | null
    │   ├── canvasObjects: ObservableArray<CanvasObject>
    │   └── selectedObjectIds: ObservableArray<string>
    │
    ├── Actions
    │   ├── setCurrentPage(pageId: string): void
    │   ├── addComponent(component: Component, position: {x, y}): CanvasObject
    │   ├── deleteObjects(ids: string[]): void
    │   ├── selectObjects(ids: string[]): void
    │   ├── updateObjectProperty(id: string, property: string, value: any): void
    │   └── reorderLayers(ids: string[]): void
    │
    └── Computed
        ├── currentPage: Page | undefined
        ├── selectedObjects: CanvasObject[]
        └── layers: Layer[]
```

## 数据流

### 用户操作流

```
用户点击左侧导航栏"画面"Tab
        ↓
LeftSidebar.tsx
  触发 onTabChange('画面')
        ↓
editorStore.setLeftSidebarActiveTab('画面')
        ↓
MobX响应式更新
  leftSidebarActiveTab = '画面'
        ↓
LeftSidebar重新渲染
  显示PagesPanel组件
```

### 对象选中流

```
用户在Canvas中点击矩形组件
        ↓
Canvas.tsx
  触发 onSelect(['rect-001'])
        ↓
canvasStore.selectObjects(['rect-001'])
        ↓
editorStore.setSelectedObjects(['rect-001'])
        ↓
MobX响应式更新
  selectedObjectIds = ['rect-001']
  rightPanelVisible = true
        ↓
多个组件同时响应:
  - Canvas: 显示选中框
  - RightPanel: 从隐藏变为显示
  - PropertiesPanel: 显示rect-001的属性
  - LayersPanel: 高亮rect-001图层
```

### 属性编辑流

```
用户在PropertiesPanel修改矩形背景色
        ↓
PropertiesPanel.tsx
  触发 onPropertyChange('rect-001', 'backgroundColor', '#FF0000')
        ↓
canvasStore.updateObjectProperty('rect-001', 'backgroundColor', '#FF0000')
        ↓
MobX响应式更新
  canvasObjects[0].style.backgroundColor = '#FF0000'
        ↓
多个组件同时响应:
  - Canvas: 矩形组件实时更新背景色
  - PropertiesPanel: 保持当前值
  - editorStore.hasUnsavedChanges = true
```

## 技术栈

### 前端框架
- **React 18.3.1**: UI组件化框架
- **TypeScript 5.x**: 类型安全
- **Vite 5.4.2**: 构建工具

### 状态管理
- **MobX 6.12.0**: 响应式状态管理
  - observable: 响应式数据
  - computed: 派生状态
  - actions: 状态更新
  - makeAutoObservable: 自动装饰器

### 样式方案
- **Plain CSS**: 原生CSS
- **CSS Variables**: 动态主题
- **BEM命名**: 块-元素-修饰符

### 图标方案
- **SVG线框图标**: 可缩放矢量图标
- **1.5px描边**: 统一描边宽度
- **round caps/joins**: 圆角端点

### 构建工具
- **Vite 5.4.2**: 快速开发服务器
- **ESBuild**: 高性能打包
- **PostCSS**: CSS后处理

## 组件通信模式

### 1. Props传递 (父 → 子)
```typescript
// EditorLayout → LeftSidebar
<LeftSidebar
  activeTab={editorStore.leftSidebarActiveTab}
  onTabChange={(tab) => editorStore.setLeftSidebarActiveTab(tab)}
  project={projectStore.currentProject}
  pages={projectStore.pages}
  components={projectStore.components}
/>
```

### 2. 回调函数 (子 → 父)
```typescript
// LeftSidebar → EditorLayout
const handleTabChange = (tab: '工程' | '画面' | '组件') => {
  editorStore.setLeftSidebarActiveTab(tab);
};

<SidebarTab activeTab={activeTab} onChange={handleTabChange} />
```

### 3. MobX Store (跨组件)
```typescript
// Canvas组件更新状态
import { useCanvasStore } from '@/store';

const canvasStore = useCanvasStore();
const handleSelectObject = (id: string) => {
  canvasStore.selectObjects([id]);
  // editorStore自动响应，RightPanel显示
};
```

### 4. MobX Reactions (自动响应)
```typescript
// PropertiesPanel自动响应选中对象变化
const PropertiesPanel = observer(() => {
  const canvasStore = useCanvasStore();
  const selectedObjects = canvasStore.selectedObjects; // 自动追踪

  return (
    <div>
      {selectedObjects.map(obj => (
        <PropertyField
          key={obj.id}
          object={obj}
          onChange={(value) => canvasStore.updateObjectProperty(obj.id, 'backgroundColor', value)}
        />
      ))}
    </div>
  );
});
```

## 布局策略

### Flexbox布局

```css
.EditorLayout {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.EditorLayout__center {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.EditorLayout__main {
  display: flex;
  flex: 1;
  overflow: hidden;
}
```

### 响应式布局

```css
/* 窗口宽度 < 1200px */
@media (max-width: 1199px) {
  .EditorLayout__main {
    flex-direction: column;
  }

  .LeftSidebar,
  .RightPanel {
    width: 240px; /* 缩小侧边栏宽度 */
  }
}

/* 窗口宽度 < 768px */
@media (max-width: 767px) {
  .LeftSidebar {
    position: absolute;
    z-index: 100;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .LeftSidebar--open {
    transform: translateX(0);
  }
}
```

### 动态显示/隐藏

```css
.RightPanel {
  width: 280px;
  transition: margin-right 0.3s ease;
}

.RightPanel--hidden {
  margin-right: -280px;
}
```

## 性能优化策略

### 1. 组件级优化

```typescript
// 使用React.memo避免不必要的重渲染
export const ToolbarButton = React.memo<ToolbarButtonProps>(({ icon, onClick, active }) => {
  return <button onClick={onClick}>{icon}</button>;
});

// 使用useMemo缓存计算结果
const expensiveValue = useMemo(() => {
  return computeExpensiveValue(props.data);
}, [props.data]);

// 使用useCallback缓存回调函数
const handleTabChange = useCallback((tab: string) => {
  editorStore.setLeftSidebarActiveTab(tab);
}, []);
```

### 2. MobX优化

```typescript
// 使用computed派生状态
get hasUnsavedChanges(): boolean {
  return this.projectTabs.some(tab => tab.isDirty);
}

// 精确追踪observable
const selectedObjects = computed(() => {
  return canvasStore.selectedObjectIds
    .map(id => canvasStore.canvasObjects.find(obj => obj.id === id))
    .filter(Boolean);
});
```

### 3. 虚拟滚动（可选，P2）

```typescript
// 对于大量组件，使用虚拟滚动
import { FixedSizeList } from 'react-window';

<FixedSizeList
  height={600}
  itemCount={components.length}
  itemSize={80}
  width="100%"
>
  {({ index, style }) => (
    <ComponentItem
      key={components[index].id}
      component={components[index]}
      style={style}
    />
  )}
</FixedSizeList>
```

## 文件组织结构

```
platforms/scada/packages/renderer/src/
├── components/
│   └── editor/                    # 工程编辑器组件
│       ├── EditorLayout.tsx       # 主布局容器
│       ├── EditorLayout.css
│       ├── toolbar/                # 工具栏组件
│       │   ├── TopToolbar.tsx
│       │   ├── TopToolbar.css
│       │   ├── SubToolbar.tsx
│       │   ├── SubToolbar.css
│       │   ├── ToolbarButton.tsx
│       │   ├── ToolbarButton.css
│       │   └── ToolbarGroup.tsx
│       ├── sidebar/                # 左侧导航栏组件
│       │   ├── LeftSidebar.tsx
│       │   ├── LeftSidebar.css
│       │   ├── SidebarTab.tsx
│       │   ├── SidebarTab.css
│       │   ├── ProjectPanel.tsx
│       │   ├── ProjectPanel.css
│       │   ├── PagesPanel.tsx
│       │   ├── PagesPanel.css
│       │   ├── ComponentsPanel.tsx
│       │   └── ComponentsPanel.css
│       ├── rightpanel/             # 右侧属性栏组件
│       │   ├── RightPanel.tsx
│       │   ├── RightPanel.css
│       │   ├── PropertiesPanel.tsx
│       │   ├── PropertiesPanel.css
│       │   ├── LayersPanel.tsx
│       │   └── LayersPanel.css
│       ├── tabs/                   # 标签页组件
│       │   ├── ProjectTab.tsx
│       │   ├── ProjectTab.css
│       │   ├── SubPageTabs.tsx
│       │   ├── SubPageTabs.css
│       │   ├── SubPageTab.tsx
│       │   └── SubPageTab.css
│       ├── canvas/                 # 画布组件
│       │   ├── CanvasArea.tsx
│       │   ├── CanvasArea.css
│       │   ├── Canvas.tsx
│       │   └── Canvas.css
│       └── statusbar/              # 状态栏组件
│           ├── StatusBar.tsx
│           ├── StatusBar.css
│           ├── ZoomControls.tsx
│           └── ZoomControls.css
├── store/
│   ├── editorStore.ts             # 编辑器全局状态
│   ├── projectStore.ts            # 工程数据状态
│   ├── canvasStore.ts             # 画布状态
│   └── index.ts
├── types/
│   ├── editor.ts                  # 编辑器类型定义
│   ├── project.ts                 # 工程类型定义
│   ├── canvas.ts                  # 画布类型定义
│   └── index.ts
├── styles/
│   └── editor.css                 # 编辑器样式变量
└── constants/
    ├── toolbarItems.ts           # 工具栏配置
    └── componentLibrary.ts       # 组件库定义
```

## 集成点

### 与TitleBar集成 (FE-005)

```typescript
// TitleBar.tsx
import { ProjectTabs } from './editor/tabs/ProjectTabs';

export const TitleBar: React.FC = () => {
  return (
    <div className="title-bar">
      <LogoSection />
      <ProjectTabs /> {/* NEW: FE-006-01 */}
      <WindowControls />
    </div>
  );
};
```

### 与MainContent集成 (FE-001)

```typescript
// MainContent.tsx
import { EditorLayout } from './editor/EditorLayout';
import { useEditorStore } from '@/store';

export const MainContent: React.FC = () => {
  const editorStore = useEditorStore();
  const activeProjectTab = editorStore.activeProjectTab;

  if (activeProjectTab === 'home') {
    return <HomePageContent />; // 原有首页内容
  }

  return <EditorLayout />; // NEW: FE-006
};
```

## 扩展点

### 1. 自定义工具栏按钮

```typescript
// constants/toolbarItems.ts
export const toolbarItems = [
  {
    group: '文件',
    items: [
      { id: 'new', icon: 'icon-new', label: '新建', action: 'createProject' },
      { id: 'open', icon: 'icon-open', label: '打开', action: 'openProject' },
      // 可以添加更多按钮
    ]
  },
  // ...
];
```

### 2. 自定义组件库

```typescript
// constants/componentLibrary.ts
export const componentLibrary = [
  {
    id: 'button',
    category: 'basic',
    name: '按钮',
    icon: 'ButtonIcon',
    defaultSize: { width: 80, height: 32 },
    // 可以扩展更多组件
  },
  // ...
];
```

### 3. 自定义属性编辑器

```typescript
// components/editor/rightpanel/PropertiesPanel.tsx
// 可以根据组件类型动态渲染不同的属性编辑器
const propertyEditors = {
  rectangle: RectanglePropertyEditor,
  circle: CirclePropertyEditor,
  text: TextPropertyEditor,
  // 可以扩展更多编辑器
};
```

## 总结

Scada工程编辑器采用模块化、组件化的架构设计，通过MobX实现响应式状态管理，使用Flexbox实现灵活布局。整体架构清晰，组件职责单一，便于后续功能扩展和维护。
