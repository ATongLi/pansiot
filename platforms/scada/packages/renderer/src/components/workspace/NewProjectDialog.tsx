/**
 * NewProjectDialog 新建工程对话框组件
 * 使用 Container/Presenter 模式重构
 *
 * 架构：
 * - Container (NewProjectDialogContainer): 连接 store，处理业务逻辑
 * - View (NewProjectDialogView): 纯 UI 渲染，可重用
 *
 * 向后兼容：
 * - 默认导出 Container 组件，保持原有 API 不变
 * - 可单独导入 View 组件用于自定义场景
 */

// 导出 Container 组件（默认导出，保持向后兼容）
export { default } from './NewProjectDialogContainer'
export type { default as NewProjectDialogContainer } from './NewProjectDialogContainer'

// 导出 View 组件（用于自定义场景）
export { default as NewProjectDialogView } from './NewProjectDialogView'
export type { NewProjectDialogViewProps, FormData } from './NewProjectDialogView'

// 导出类型（保持向后兼容）
export type { NewProjectDialogContainerProps } from './NewProjectDialogContainer'
