/**
 * Editor Components Export
 * IMP-009: 工程编辑器组件导出（使用 Container/Presenter 模式）
 */

// EditOptionsBar - 导出 Container（默认），同时导出类型
export { default as EditOptionsBar } from './EditOptionsBar'
export type { EditOptionsBarContainerProps as EditOptionsBarProps } from './EditOptionsBarContainer'

// SubPageArea - 导出 Container（默认），同时导出类型
export { default as SubPageArea } from './SubPageArea'
export type { SubPageAreaContainerProps as SubPageAreaProps } from './SubPageAreaContainer'

// EditorLayout - 导出 Container（命名导出），同时导出类型
export { EditorLayout } from './EditorLayout'
export type { EditorLayoutContainerProps as EditorLayoutProps } from './EditorLayoutContainer'
