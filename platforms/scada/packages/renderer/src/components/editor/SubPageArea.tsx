/**
 * SubPageArea 子页面内容区组件
 * FE-009-06: 使用 Container/Presenter 模式重构
 *
 * 架构：
 * - Container (SubPageAreaContainer): 连接 subPageTabStore，处理业务逻辑
 * - View (SubPageAreaView): 纯 UI 渲染，可重用
 *
 * 向后兼容：
 * - 默认导出 Container 组件，保持原有 API 不变
 * - 可单独导入 View 组件用于自定义场景
 */

// 导出 Container 组件（默认导出，保持向后兼容）
export { default } from './SubPageAreaContainer'
export type { default as SubPageAreaContainer } from './SubPageAreaContainer'

// 导出 View 组件（用于自定义场景）
export { default as SubPageAreaView } from './SubPageAreaView'
export type { SubPageAreaViewProps } from './SubPageAreaView'

// 导出类型（保持向后兼容）
export type { SubPageAreaContainerProps } from './SubPageAreaContainer'
