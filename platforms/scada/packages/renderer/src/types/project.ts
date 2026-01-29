/**
 * 工程类型定义
 * 定义 Scada 工程管理功能的所有数据结构
 */

/**
 * 工程分类枚举（内置分类 + 自定义分类）
 */
export enum ProjectCategory {
  CATEGORY_1 = '分类1',
  CATEGORY_2 = '分类2',
  CUSTOM = '新建分类...'
}

/**
 * 运行平台枚举（设备型号）
 * 注意：平台列表从后端API获取，此枚举仅用于TypeScript类型检查
 */
export enum HardwarePlatform {
  BOX1 = 'BOX1',     // 标准BOX型号
  HMI01 = 'HMI01',   // HMI触摸屏
  TBOX1 = 'TBOX1',   // 网关盒子
  // 未来可能添加的平台
  // BOX2 = 'BOX2',
  // HMI02 = 'HMI02',
}

/**
 * 硬件平台接口（从后端获取）
 */
export interface HardwarePlatformConfig {
  id: string              // 平台唯一标识
  name: string            // 平台显示名称
  type: PlatformType      // 平台类型
  resolution?: string     // 屏幕分辨率
  enabled: boolean        // 是否启用
}

/**
 * 平台类型枚举
 */
export enum PlatformType {
  BOX = 'box',
  HMI = 'hmi',
  GATEWAY = 'gateway'
}

/**
 * 工程元数据
 */
export interface ProjectMetadata {
  /** 工程名称 */
  name: string
  /** 工程作者 */
  author?: string
  /** 工程描述 */
  description?: string
  /** 工程分类 */
  category: string
  /** 运行平台（设备型号） */
  platform: string
  /** 创建时间 (ISO 8601) */
  createdAt: string
  /** 修改时间 (ISO 8601) */
  updatedAt: string
}

/**
 * 工程安全配置
 */
export interface ProjectSecurity {
  /** 是否加密 */
  encrypted: boolean
  /** 密码哈希 (bcrypt, 仅加密工程有) */
  passwordHash?: string
  /** 设备指纹 (可选) */
  deviceBinding?: string
  /** 文件签名 (HMAC-SHA256) */
  fileSignature: string
  /** KEK版本 (双重加密机制) */
  kekVersion?: string
  /** 数据加密密钥 (双重加密机制) */
  dek?: {
    /** 用户密码加密的DEK */
    userEncrypted: string
    /** 官方KEK加密的DEK备份 */
    officialEncrypted: string
  }
}

/**
 * 画布配置
 */
export interface CanvasConfig {
  /** 画布宽度 */
  width: number
  /** 画布高度 */
  height: number
  /** 背景颜色 */
  backgroundColor?: string
}

/**
 * 组件定义
 */
export interface Component {
  /** 组件ID */
  id: string
  /** 组件类型 */
  type: string
  /** X坐标 */
  x: number
  /** Y坐标 */
  y: number
  /** 宽度 */
  width: number
  /** 高度 */
  height: number
  /** 组件属性 */
  properties: Record<string, any>
  /** 数据绑定 */
  dataBindings?: DataBinding[]
}

/**
 * 数据绑定
 */
export interface DataBinding {
  /** 组件ID */
  componentId: string
  /** 属性名 */
  property: string
  /** 数据源 */
  source: string
}

/**
 * 工程主数据结构
 */
export interface Project {
  /** 文件格式版本 */
  version: string
  /** 工程唯一标识 (UUID v4) */
  projectId: string
  /** 工程元数据 */
  metadata: ProjectMetadata
  /** 安全配置 */
  security: ProjectSecurity
  /** 画布配置 */
  canvas: CanvasConfig
  /** 组件列表 */
  components: Component[]
  /** 加密内容 (仅加密工程有) */
  encryptedContent?: string
}

/**
 * 最近工程数据
 */
export interface RecentProject {
  /** 工程唯一标识 */
  projectId: string
  /** 工程名称 */
  name: string
  /** 工程分类 */
  category?: string
  /** 工程文件路径 */
  filePath: string
  /** 最后打开时间 (相对时间字符串, 如 "2小时前") */
  lastOpened: string
  /** 最后打开时间 (Date对象) */
  lastOpenedDate: Date
  /** 是否加密 */
  isEncrypted: boolean
  /** 创建时间 */
  createdAt: string
}

/**
 * 新建工程表单数据
 */
export interface NewProjectFormData {
  /** 工程名称 */
  name: string
  /** 工程作者 */
  author?: string
  /** 工程描述 */
  description?: string
  /** 工程分类 */
  category: string
  /** 运行平台 */
  platform: string
  /** 工程加密 */
  encrypted: boolean
  /** 密码 (加密工程需要) */
  password?: string
  /** 确认密码 */
  confirmPassword?: string
  /** 保存位置 */
  savePath: string
}

/**
 * 打开工程请求数据
 */
export interface OpenProjectRequest {
  /** 工程文件路径 */
  filePath: string
  /** 密码 (加密工程需要) */
  password?: string
}

/**
 * API响应基础结构
 */
export interface ApiResponse<T = any> {
  /** 是否成功 */
  success: boolean
  /** 数据 */
  data?: T
  /** 错误码 */
  error?: string
  /** 错误信息 */
  message?: string
}

/**
 * 创建工程响应
 */
export interface CreateProjectResponse {
  /** 工程ID */
  projectId: string
  /** 工程文件路径 */
  filePath: string
}

/**
 * 打开工程响应
 */
export interface OpenProjectResponse {
  /** 工程数据 */
  project: Project
}
