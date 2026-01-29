/**
 * SVG Icon Constants
 * 线框型图标定义 - 工业专业风格
 *
 * 图标规格:
 * - 描边宽度: 1.5px
 * - 描边颜色: currentColor (继承父元素文字颜色)
 * - 填充: none
 * - 圆角: stroke-linecap="round" stroke-linejoin="round"
 * - 尺寸: 24×24px viewBox
 *
 * 参考样式: 石墨文档图标风格
 */

/**
 * 导航图标集
 */
export const NAV_ICONS = {
  /**
   * 首页图标
   * 房屋轮廓线框
   */
  home: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    <polyline points="9,22 9,12 15,12 15,22" />
  </svg>`,

  /**
   * 本地/文件夹图标
   * 文件夹轮廓线框
   */
  local: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
  </svg>`,

  /**
   * 云端图标
   * 云朵轮廓线框
   */
  cloud: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z" />
  </svg>`,

  /**
   * 工具图标
   * 扳手轮廓线框
   */
  tools: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z" />
  </svg>`,

  /**
   * 用户图标
   * 用户轮廓线框
   */
  user: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
    <circle cx="12" cy="7" r="4" />
  </svg>`,
};

/**
 * 窗口控制图标集 (Windows风格)
 */
export const WINDOW_ICONS = {
  /**
   * 最小化图标
   * 横线
   */
  minimize: `<svg viewBox="0 0 10 1" fill="none" stroke="currentColor" stroke-width="1">
    <rect x="0" y="0" width="10" height="1" fill="currentColor"/>
  </svg>`,

  /**
   * 最大化图标
   * 矩形边框
   */
  maximize: `<svg viewBox="0 0 10 10" fill="none" stroke="currentColor" stroke-width="1">
    <rect x="0.5" y="0.5" width="9" height="9"/>
  </svg>`,

  /**
   * 还原图标
   * 重叠矩形
   */
  restore: `<svg viewBox="0 0 10 10" fill="none" stroke="currentColor" stroke-width="1">
    <rect x="2" y="0.5" width="7.5" height="7.5"/>
    <rect x="0.5" y="2" width="7.5" height="7.5" fill="white"/>
    <rect x="2" y="2" width="6" height="6"/>
  </svg>`,

  /**
   * 关闭图标
   * X形
   */
  close: `<svg viewBox="0 0 10 10" fill="none" stroke="currentColor" stroke-width="1">
    <path d="M1 1L9 9M9 1L1 9"/>
  </svg>`,
};

/**
 * 品牌图标集 (CR-002补充)
 */
export const BRANDING_ICONS = {
  /**
   * Logo图标
   * 简洁的品牌标识
   */
  logo: `<svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="2" y="2" width="28" height="28" rx="6" fill="#2196F3"/>
    <path d="M10 22L16 10L22 22M16 10V16" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>`,
};

/**
 * 工程操作图标集 (CR-002第三批补充)
 */
export const ACTION_ICONS = {
  /**
   * 新建工程图标
   * 48×48px大图标，包含加号
   */
  newProject: `<svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <rect x="8" y="8" width="32" height="32" rx="4"/>
    <path d="M24 16V32M16 24H32"/>
  </svg>`,

  /**
   * 从文件打开图标
   * 48×48px大图标，文件夹带勾选标记
   */
  openProject: `<svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <path d="M6 20V14C6 11.7909 7.79086 10 10 10H18L22 14H38C40.2091 14 42 15.7909 42 18V34C42 36.2091 40.2091 38 38 38H10C7.79086 38 6 36.2091 6 34V20Z"/>
    <path d="M22 26L26 30L34 22"/>
  </svg>`,

  /**
   * 复制工程图标
   * 48×48px大图标，两个重叠的文档
   */
  copyProject: `<svg viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
    <rect x="8" y="8" width="20" height="20" rx="4"/>
    <path d="M20 16H38C40.2091 16 42 17.7909 42 20V38C42 40.2091 40.2091 42 38 42H20C17.7909 42 16 40.2091 16 38V28"/>
  </svg>`,
};

/**
 * 图标映射表
 * 统一导出所有图标
 */
export const ICONS = {
  ...NAV_ICONS,
  ...WINDOW_ICONS,
  ...BRANDING_ICONS,
  ...ACTION_ICONS,
};

/**
 * 默认导出
 */
export default ICONS;
