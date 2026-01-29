/**
 * EditOptionsBar 编辑选项栏组件
 * FE-009-05: 使用 Container/Presenter 模式重构
 *
 * 架构：
 * - Container (EditOptionsBarContainer): 连接 layoutStore，处理业务逻辑
 * - View (EditOptionsBarView): 纯 UI 渲染，可重用
 *
 * 向后兼容：
 * - 默认导出 Container 组件，保持原有 API 不变
 * - 可单独导入 View 组件用于自定义场景
 */

// 导出 Container 组件（默认导出，保持向后兼容）
export { default } from './EditOptionsBarContainer'
export type { default as EditOptionsBarContainer } from './EditOptionsBarContainer'

// 导出 View 组件（用于自定义场景）
export { default as EditOptionsBarView } from './EditOptionsBarView'
export type { EditOptionsBarViewProps } from './EditOptionsBarView'

// 导出类型（保持向后兼容）
export type { EditOptionsBarContainerProps } from './EditOptionsBarContainer'
