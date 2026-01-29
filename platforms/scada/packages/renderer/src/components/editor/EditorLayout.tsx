/**
 * EditorLayout Component
 *
 * FE-009-03: 工程编辑器布局重构
 * 使用 Container/Presenter 模式
 *
 * 架构：
 * - Container (EditorLayoutContainer): 连接 stores，处理业务逻辑
 * - View (EditorLayoutView): 纯 UI 渲染，可重用
 *
 * 向后兼容：
 * - 默认导出 Container 组件，保持原有 API 不变
 * - 可单独导入 View 组件用于自定义场景
 */

// 导出 Container 组件（默认导出和命名导出，保持向后兼容）
export { default as EditorLayout, EditorLayoutContainer } from './EditorLayoutContainer'
export type { default as EditorLayoutContainerType } from './EditorLayoutContainer'

// 导出 View 组件（用于自定义场景）
export { EditorLayoutView, default as EditorLayoutViewDefault } from './EditorLayoutView'
export type { EditorLayoutViewProps } from './EditorLayoutView'

// 导出类型（保持向后兼容）
export type { EditorLayoutContainerProps } from './EditorLayoutContainer'
