/**
 * 线条图标组件库
 * 用于工具栏和界面图标
 */

import React from 'react'
import './LineIcon.css'

export interface LineIconProps {
  name: string
  size?: number
  className?: string
}

/**
 * 线条图标组件
 */
export const LineIcon: React.FC<LineIconProps> = ({ name, size = 16, className = '' }) => {
  const icons: Record<string, JSX.Element> = {
    // 工程文件
    'new-project': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3h10M8 3v10M3 8h10" strokeLinecap="round" />
      </svg>
    ),
    'open-project': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M4 3h8M2 6h12M2 9h12M2 12h8" strokeLinecap="round" />
      </svg>
    ),
    'save-project': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 2v12h10V5l-3-3H3z" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'save-as': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 2v12h10V5l-3-3H3z" strokeLinecap="round" strokeLinejoin="round" />
        <path d="M10 8l3 3M13 8v6" strokeLinecap="round" />
      </svg>
    ),
    'auto-save': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 2v6M5 5h6" strokeLinecap="round" />
        <circle cx="8" cy="11" r="3" />
      </svg>
    ),
    'project-settings': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <circle cx="8" cy="8" r="3" />
        <path d="M8 1v4m0 6v4M1 8h4m6 0h4" strokeLinecap="round" />
      </svg>
    ),

    // 通用 - 剪贴板
    'copy': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <rect x="4" y="2" width="8" height="12" rx="1" />
        <path d="M6 5h4" strokeLinecap="round" />
      </svg>
    ),
    'paste': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M6 3v10M10 3v10M3 6h10M3 10h10" strokeLinecap="round" />
      </svg>
    ),
    'delete': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 4h10M5 4V3h6v1M8 7v5M6 9l2-2M10 9l-2-2M4 4v9h8V4" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'format-brush': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M13 3l-2 2M5 11l2 2M8 8l3-3M3 13l2-2" strokeLinecap="round" />
        <circle cx="6" cy="6" r="2" />
      </svg>
    ),

    // 通用 - 编辑
    'undo': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 3H5v5h5M5 8c3 0 5-2 5-5s-2-5-5-5" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'redo': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 3h3v5H6M11 8c-3 0-5-2-5-5s2-5 5-5" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'select': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3h4v4H3z M9 3h4v4H9z M3 9h4v4H3z M9 9h4v4H9z" />
      </svg>
    ),
    'multi-duplicate': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <rect x="3" y="3" width="4" height="4" rx="0.5" />
        <rect x="9" y="9" width="4" height="4" rx="0.5" />
        <path d="M7 7l2 2" strokeLinecap="round" />
      </svg>
    ),

    // 通用 - 对齐
    'align-left': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3h8M3 6h10M3 9h10M3 12h6" strokeLinecap="round" />
      </svg>
    ),
    'align-right': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M13 3H5M13 6H3M13 9H3M13 12h-6" strokeLinecap="round" />
      </svg>
    ),
    'align-center': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3h10M5 6h6M3 9h10M5 12h6" strokeLinecap="round" />
      </svg>
    ),

    // 通用 - 组合
    'group': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <rect x="3" y="5" width="5" height="6" rx="1" />
        <rect x="8" y="5" width="5" height="6" rx="1" />
      </svg>
    ),
    'ungroup': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 3v10M3 8h10" strokeLinecap="round" />
        <rect x="3" y="5" width="4" height="3" rx="0.5" />
        <rect x="9" y="8" width="4" height="3" rx="0.5" />
      </svg>
    ),

    // 运行调试
    'compile': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 2v12M3 7l5 3 5-3" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'online-compile': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <circle cx="8" cy="8" r="2" />
        <path d="M8 4v4m0 4v4M4 8h4m0 0h4" strokeLinecap="round" />
      </svg>
    ),
    'offline-compile': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 2v12M3 7l5 3 5-3" strokeLinecap="round" strokeLinejoin="round" />
        <line x1="12" y1="2" x2="12" y2="6" strokeLinecap="round" />
        <line x1="14" y1="2" x2="14" y2="6" strokeLinecap="round" />
      </svg>
    ),
    'download': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 3v9M4 8l4 4 4-4" strokeLinecap="round" strokeLinejoin="round" />
        <path d="M3 13h10" strokeLinecap="round" />
      </svg>
    ),
    'upload': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 13v-9M4 8l4-4 4 4" strokeLinecap="round" strokeLinejoin="round" />
        <path d="M3 13h10" strokeLinecap="round" />
      </svg>
    ),

    // 工具
    'text-library': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M4 3h8M8 3v10M4 8h8M4 13h8" strokeLinecap="round" />
      </svg>
    ),
    'media-library': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <rect x="2" y="3" width="12" height="10" rx="1" />
        <circle cx="6" cy="8" r="1.5" fill="currentColor" />
        <path d="M10 6l2 2-2 2" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'device-manager': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <rect x="3" y="2" width="10" height="12" rx="1" />
        <circle cx="8" cy="12" r="1" fill="currentColor" />
        <path d="M6 5h4M6 8h4" strokeLinecap="round" />
      </svg>
    ),
    'passthrough': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 8h3m5 0h5M8 3v10" strokeLinecap="round" />
        <circle cx="8" cy="8" r="1" />
      </svg>
    ),
    'debug-log': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3v10M3 13h10M8 3v6" strokeLinecap="round" />
        <circle cx="8" cy="10" r="1" fill="currentColor" />
      </svg>
    ),

    // 视图
    'layer': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M2 8l6-3 6 3M2 5l6-3 6 3M2 11l6 3 6-3" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'history': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <circle cx="8" cy="8" r="6" />
        <path d="M8 4v4l3 3" strokeLinecap="round" />
      </svg>
    ),
    'grid': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 3h3v3H3z M7 3h3v3H7z M11 3h3v3h-3z M3 7h3v3H3z M7 7h3v3H7z M11 7h3v3h-3z M3 11h3v3H3z M7 11h3v3H7z M11 11h3v3h-3z" />
      </svg>
    ),
    'snap-grid': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <circle cx="4" cy="4" r="1" />
        <circle cx="8" cy="4" r="1" />
        <circle cx="12" cy="4" r="1" />
        <circle cx="4" cy="8" r="1" />
        <circle cx="8" cy="8" r="1" />
        <circle cx="12" cy="8" r="1" />
        <circle cx="4" cy="12" r="1" />
        <circle cx="8" cy="12" r="1" />
        <circle cx="12" cy="12" r="1" />
      </svg>
    ),

    // 帮助
    'version': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 2v6M8 12h.01" strokeLinecap="round" />
        <circle cx="8" cy="8" r="6" />
      </svg>
    ),
    'changelog': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 8c0-3 2-5 5-5s5 2 5 5-2 5-5 5" strokeLinecap="round" />
        <path d="M12 5l3 3-3 3" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    ),
    'update': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M8 2v6M8 11l3-3M11 8l-3-3" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="8" cy="8" r="6" />
      </svg>
    ),
    'help-docs': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M6 6c0-2 1-3 3-3s2 1 2 3c0 2-3 1-3 4M6 13h.01" strokeLinecap="round" />
        <circle cx="8" cy="8" r="6" />
      </svg>
    ),
    'contact': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <circle cx="6" cy="6" r="2" />
        <path d="M2 14c0-2 2-3 4-3s3 1 3 3c0 1-1 2-4 2" strokeLinecap="round" />
        <path d="M10 6h4M10 10h4" strokeLinecap="round" />
      </svg>
    ),

    // 控制按钮
    'collapse': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M4 8h8" strokeLinecap="round" />
      </svg>
    ),
    'expand': (
      <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.2">
        <path d="M3 8h10M8 3v10" strokeLinecap="round" />
      </svg>
    ),
  }

  const icon = icons[name]

  if (!icon) {
    console.warn(`LineIcon: "${name}" not found`)
    return null
  }

  return (
    <span
      className={`line-icon ${className}`}
      style={{
        width: size,
        height: size,
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      {React.cloneElement(icon, {
        width: size,
        height: size,
        style: { width: '100%', height: '100%' },
        strokeWidth: 1.5,
        strokeLinecap: 'round',
        strokeLinejoin: 'round',
      })}
    </span>
  )
}

export default LineIcon
