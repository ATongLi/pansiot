# Scada UI Style Guide - 官方样式规范

**版本**: 1.0.0
**最后更新**: 2026-01-28
**适用范围**: PanTools Scada 组态软件 - Renderer 进程

---

## 目录

1. [设计哲学](#1-设计哲学)
2. [颜色系统](#2-颜色系统)
3. [字体系统](#3-字体系统)
4. [间距系统](#4-间距系统)
5. [组件设计模式](#5-组件设计模式)
6. [布局标准](#6-布局标准)
7. [动画和交互](#7-动画和交互)
8. [代码示例](#8-代码示例)
9. [反模式](#9-反模式)
10. [最佳实践](#10-最佳实践)

---

## 1. 设计哲学

### 1.1 工业软件 vs 消费软件

PanTools Scada 是**工业级组态软件**，设计原则与消费类软件有本质区别：

| 特性 | 消费软件 | 工业软件 (PanTools) |
|------|---------|---------------------|
| **布局** | 宽松、留白多 | 紧凑、信息密度高 |
| **间距** | 16-24px | 8-12px |
| **字号** | 14-16px 为主 | 11-14px 为主 |
| **圆角** | 8-16px | 2-4px |
| **颜色** | 鲜艳多彩 | 灰色+蓝色系 |
| **交互** | 动画丰富 | 快速响应 |
| **目标** | 视觉吸引 | 效率优先 |

### 1.2 核心设计原则

1. **紧凑 (Compact)**: 在有限空间内展示更多信息
2. **高效 (Efficient)**: 减少操作步骤，提升工作效率
3. **专业 (Professional)**: 统一、严谨、可信赖
4. **一致 (Consistent)**: 所有组件遵循相同的设计语言

---

## 2. 颜色系统

### 2.1 颜色令牌 (Color Tokens)

所有颜色使用 CSS 变量定义，位于 `src/index.css`。

```css
/* 背景色层级 */
--bg-primary: #f5f5f5;          /* 主背景 */
--bg-secondary: #ffffff;         /* 次级背景、卡片背景 */
--bg-tertiary: #e8e8e8;          /* 三级背景 */
--bg-hover: #f0f0f0;             /* 悬停背景 */

/* 文本色层级 */
--text-primary: #212121;         /* 主要文本 */
--text-secondary: #757575;       /* 次要文本 */
--text-disabled: #bdbdbd;        /* 禁用文本 */

/* 强调色 (Material Blue) */
--color-accent: #2196F3;         /* 主强调色 */
--color-accent-light: #e3f2fd;   /* 浅强调色背景 */
--color-accent-active: #1976D2;  /* 激活状态强调色 */

/* 边框和分割线 */
--border-color: #e0e0e0;
--divider-color: #eeeeee;

/* 阴影 */
--color-shadow: rgba(0, 0, 0, 0.1);
```

### 2.2 颜色使用指南

#### 背景色使用场景

- `--bg-primary` (#f5f5f5): 应用主背景、页面背景
- `--bg-secondary` (#ffffff): 卡片背景、对话框背景、面板背景
- `--bg-tertiary` (#e8e8e8): 工具栏背景、导航栏背景
- `--bg-hover` (#f0f0f0): 按钮、列表项悬停状态

#### 文本色使用场景

- `--text-primary` (#212121): 标题、重要文本、用户输入内容
- `--text-secondary` (#757575): 描述文本、辅助信息
- `--text-disabled` (#bdbdbd): 禁用状态文本、占位符

#### 强调色使用规则

1. **主强调色** (`--color-accent`): 仅用于以下场景：
   - 主要操作按钮（保存、确认、提交）
   - 激活状态的导航项和标签页
   - 重要的交互元素（链接、开关）

2. **浅强调色** (`--color-accent-light`):
   - 激活状态的背景色
   - 标签页激活背景
   - 选中状态的边框

3. **激活强调色** (`--color-accent-active`):
   - 悬停状态的强调色
   - 焦点状态的边框
   - 按下状态的按钮

**反模式警告**: ❌ 不要将强调色用于普通文本、边框或装饰性元素。

---

## 3. 字体系统

### 3.1 字体族

```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
```

使用系统字体栈确保在各平台上的最佳显示效果。

### 3.2 字号层级

| 级别 | 大小 | 使用场景 | CSS 变量 |
|------|------|---------|----------|
| **xs** | 11px | 错误提示、辅助信息 | `--font-size-xs` |
| **sm** | 12px | 标签、次级文本 | `--font-size-sm` |
| **md** | 13px | 表单输入、按钮 | `--font-size-md` |
| **base** | 14px | 正文、常规文本 | `--font-size-base` |
| **lg** | 16px | 标题、重要文本 | `--font-size-lg` |
| **xl** | 18px | 大标题 | `--font-size-xl` |

### 3.3 字重使用规则

```css
font-weight: 400;  /* 常规：正文、描述 */
font-weight: 500;  /* 中等：小标题、标签 */
font-weight: 600;  /* 半粗：标题、重要文本 */
```

- **400 (Regular)**: 正文、描述文本、按钮文字
- **500 (Medium)**: 小标题、卡片标题、标签文本
- **600 (Semi-Bold)**: 区域标题、重要按钮文本

### 3.4 行高标准

```css
line-height: 1.4;  /* 紧凑：标题、重要文本 */
line-height: 1.5;  /* 常规：正文、表单 */
```

工业软件使用较紧凑的行高（1.4-1.5），相比消费软件的 1.6-1.8。

---

## 4. 间距系统

### 4.1 8px 网格系统

所有间距使用 8px 的倍数：

| 变量名 | 值 | 使用场景 |
|--------|----|----|----------|
| `--spacing-xs` | 4px | 最小间距、图标与文字间距 |
| `--spacing-sm` | 8px | 小间距、相关元素间距 |
| `--spacing-md` | 12px | 中等间距、表单行间距 |
| `--spacing-lg` | 16px | 大间距、组件组间距 |
| `--spacing-xl` | 24px | 超大间距、Section 间距 |

### 4.2 间距使用指南

#### 组件内部 Padding

```css
/* 紧凑按钮 */
padding: 6px 8px;

/* 标准按钮 */
padding: 12px 16px;

/* 卡片内部 */
padding: 16px;

/* 表单输入框 */
padding: 6px 8px;
```

#### 组件之间 Margin

```css
/* 相关组件 */
margin-bottom: 8px;

/* 无关组件 */
margin-bottom: 24px;

/* Section 间距 */
margin-bottom: 32px;
```

---

## 5. 组件设计模式

### 5.1 按钮模式

#### 主按钮 (Primary Button)

```css
.button-primary {
  padding: 12px 24px;
  background: var(--color-accent);
  color: #ffffff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 150ms ease;
}

.button-primary:hover {
  background: var(--color-accent-active);
}

.button-primary:active {
  transform: translateY(1px);
}
```

#### 次按钮 (Secondary Button)

```css
.button-secondary {
  padding: 12px 24px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 150ms ease;
}

.button-secondary:hover {
  background: var(--bg-hover);
  border-color: var(--color-accent-active);
}
```

### 5.2 表单控件

#### 输入框 (Input)

```css
.input {
  padding: 6px 8px;
  font-size: 13px;
  border: 1px solid var(--border-color);
  border-radius: 2px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: all 150ms ease;
}

.input:hover {
  border-color: var(--color-accent-active);
}

.input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 2px var(--color-accent-light);
}
```

#### 表单布局

```css
.form-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.form-field--half {
  flex: 1;
  min-width: 0;
}

.form-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}
```

### 5.3 对话框 (Dialog)

#### 紧凑工业风格

```css
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-container {
  width: 520px;
  max-width: 90vw;
  max-height: 90vh;
  background: var(--bg-secondary);
  border-radius: 4px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.dialog-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-body {
  padding: 16px;
  overflow-y: auto;
  max-height: calc(90vh - 45px - 64px);
}

.dialog-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
```

---

## 6. 布局标准

### 6.1 主布局结构

```
┌─────────────────────────────────────────────┐
│ TitleBar (40px)                             │  应用级标签栏
├──────────┬──────────────────────────────────┤
│ Sidebar  │ Main Content                     │
│ (60px)   │                                  │
│          │  ┌────────────────────────────┐  │
│          │  │ TopBar (32px)              │  │
│          │  ├────────────────────────────┤  │
│          │  │                           │  │
│          │  │  Content Area             │  │
│          │  │  (flex: 1)                │  │
│          │  │                           │  │
│          │  └────────────────────────────┘  │
└──────────┴──────────────────────────────────┘
```

### 6.2 尺寸标准

| 组件 | 尺寸 | 说明 |
|------|------|------|
| **TitleBar** | 40px 高度 | 应用级标题栏 |
| **TopBar** | 32px 高度 | 页面级工具栏 |
| **Sidebar** | 60px 宽度 | 应用级侧边栏（导航） |
| **LeftPanel** | 280px 宽度 | 编辑器左侧面板 |
| **RightPanel** | 280px 宽度 | 编辑器右侧面板 |
| **SubPageTab** | 32px 高度 | 内部标签栏高度 |
| **StatusBar** | 24px 高度 | 状态栏高度 |

### 6.3 对齐和定位

1. **4px 网格对齐**: 所有元素对齐到 4px 网格
2. **Flexbox 优先**: 使用 Flexbox 进行布局
3. **居中规则**:
   - 水平居中: `justify-content: center` 或 `margin: 0 auto`
   - 垂直居中: `align-items: center`

---

## 7. 动画和交互

### 7.1 过渡时长

```css
--transition-fast: 150ms;   /* 快速交互：悬停、焦点 */
--transition-normal: 250ms; /* 常规动画：展开、淡入 */
--transition-slow: 350ms;   /* 慢速动画：页面切换 */
```

### 7.2 缓动函数

```css
transition: all 150ms ease;  /* 默认缓动 */
```

使用 CSS 默认的 `ease` 缓动函数，无需自定义。

### 7.3 交互反馈

#### Hover 状态

```css
.element:hover {
  background: var(--bg-hover);
  border-color: var(--color-accent-active);
}
```

#### Active 状态

```css
.element:active {
  transform: translateY(1px);
}
```

#### Focus 状态

```css
.element:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 2px var(--color-accent-light);
}
```

---

## 8. 代码示例

### 8.1 标准按钮实现

```tsx
import React from 'react'
import './Button.css'

interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'icon'
  size?: 'small' | 'medium' | 'large'
  disabled?: boolean
  onClick?: () => void
  children: React.ReactNode
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'secondary',
  size = 'medium',
  disabled = false,
  onClick,
  children,
}) => {
  return (
    <button
      className={`button button--${variant} button--${size}`}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  )
}
```

### 8.2 紧凑表单实现

```tsx
import React from 'react'
import './Form.css'

interface FormFieldProps {
  label: string
  name: string
  type?: 'text' | 'password' | 'email'
  placeholder?: string
}

export const FormField: React.FC<FormFieldProps> = ({
  label,
  name,
  type = 'text',
  placeholder,
}) => {
  return (
    <div className="form-field">
      <label className="form-field__label" htmlFor={name}>
        {label}
      </label>
      <input
        className="form-field__input"
        type={type}
        id={name}
        name={name}
        placeholder={placeholder}
      />
    </div>
  )
}
```

---

## 9. 反模式

### ❌ 避免的做法

1. **过度使用圆角**
   ```css
   /* ❌ 错误 */
   border-radius: 12px;

   /* ✅ 正确 */
   border-radius: 4px;
   ```

2. **过大的间距**
   ```css
   /* ❌ 错误：消费软件风格 */
   padding: 24px;
   gap: 24px;

   /* ✅ 正确：工业软件紧凑风格 */
   padding: 12px;
   gap: 8px;
   ```

3. **硬编码颜色值**
   ```css
   /* ❌ 错误 */
   background: #f5f5f5;
   color: #212121;

   /* ✅ 正确 */
   background: var(--bg-primary);
   color: var(--text-primary);
   ```

4. **全屏滚动对话框**
   ```tsx
   /* ❌ 错误：整个对话框滚动 */
   <div className="dialog" style={{ overflowY: 'scroll' }}>
     <header>...</header>
     <form>...</form>
     <footer>...</footer>
   </div>

   /* ✅ 正确：只有表单区域滚动 */
   <div className="dialog">
     <header>...</header>
     <form style={{ overflowY: 'auto' }}>...</form>
     <footer>...</footer>
   </div>
   ```

---

## 10. 最佳实践

### ✅ 推荐的做法

1. **始终使用 CSS 变量**
   ```css
   color: var(--text-primary);
   margin: var(--spacing-md);
   font-size: var(--font-size-base);
   ```

2. **遵循 8px 网格系统**
   ```css
   padding: 8px;
   gap: 16px;
   margin: 24px;
   ```

3. **保持一致的 border-radius**
   ```css
   .button { border-radius: 4px; }
   .input { border-radius: 2px; }
   ```

4. **表单优先使用两列布局**
   ```tsx
   <div className="form-row">
     <FormField className="form-field--half" label="字段1" />
     <FormField className="form-field--half" label="字段2" />
   </div>
   ```

5. **使用 CSS Modules**
   ```tsx
   import styles from './Component.module.css'

   const Component = () => {
     return <div className={styles.container}>...</div>
   }
   ```

6. **组件命名遵循 BEM 规范**
   ```css
   .form-field { }
   .form-field__label { }
   .form-field__input { }
   .form-field--disabled { }
   ```

---

## 附录 A: CSS 变量速查表

```css
/* 背景色 */
--bg-primary: #f5f5f5;
--bg-secondary: #ffffff;
--bg-tertiary: #e8e8e8;
--bg-hover: #f0f0f0;

/* 文本色 */
--text-primary: #212121;
--text-secondary: #757575;
--text-disabled: #bdbdbd;

/* 强调色 */
--color-accent: #2196F3;
--color-accent-light: #e3f2fd;
--color-accent-active: #1976D2;

/* 边框和分割线 */
--border-color: #e0e0e0;
--divider-color: #eeeeee;

/* 阴影 */
--color-shadow: rgba(0, 0, 0, 0.1);

/* 间距 */
--spacing-xs: 4px;
--spacing-sm: 8px;
--spacing-md: 12px;
--spacing-lg: 16px;
--spacing-xl: 24px;

/* 字号 */
--font-size-xs: 11px;
--font-size-sm: 12px;
--font-size-md: 13px;
--font-size-base: 14px;
--font-size-lg: 16px;
--font-size-xl: 18px;

/* 过渡 */
--transition-fast: 150ms;
--transition-normal: 250ms;
--transition-slow: 350ms;

/* 圆角 */
--border-radius: 4px;
--border-radius-sm: 2px;
```

---

## 附录 B: 组件清单

### 已实现组件

✅ **应用级组件**
- TitleBar: 应用标题栏
- TabBar: 应用级标签栏
- Tab: 标签页
- WindowControls: 窗口控制按钮

✅ **页面组件**
- HomePage: 首页（左右分栏）
- HomePageSidebar: 首页左侧导航
- HomePageContent: 首页右侧内容

✅ **编辑器组件**
- EditorLayout: 编辑器主布局
- EditOptionsBar: 编辑选项栏（可配置布局）
- SubPageArea: 子页面内容区
- SubPageTabs: 内部标签栏
- Canvas: 画布区域

✅ **工作区组件**
- ActionButtons: 操作按钮卡片
- RecentProjects: 最近工程列表

---

**文档维护**: 本文档应随 UI 组件的演进持续更新。所有新增组件必须遵循本规范。

**最后审核**: 2026-01-28
**审核人**: IMP-009 实施团队
