/**
 * 工具栏分类类型定义
 */

/**
 * 主工具栏分类
 */
export enum ToolbarCategory {
  PROJECT = 'project',      // 工程文件
  GENERAL = 'general',      // 通用
  DEBUG = 'debug',          // 运行调试
  TOOLS = 'tools',          // 工具
  VIEW = 'view',            // 视图
  HELP = 'help',            // 帮助
}

/**
 * 子工具栏按钮类型
 */
export interface SubToolbarButton {
  id: string
  icon: React.ReactNode
  label: string
  shortcut?: string
  disabled?: boolean
  action: () => void
}

/**
 * 子工具栏分隔符类型
 */
export interface SubToolbarSeparator {
  type: 'separator'
}

/**
 * 子工具栏项类型（按钮或分隔符）
 */
export type SubToolbarItem = SubToolbarButton | SubToolbarSeparator

/**
 * 主工具栏分类配置
 */
export interface ToolbarCategoryConfig {
  id: ToolbarCategory
  label: string
  icon: React.ReactNode
  buttons: SubToolbarButton[]
}
