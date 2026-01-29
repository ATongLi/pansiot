import React from 'react'
import { ProjectCategory } from '@/types/project'
import { HardwarePlatformConfig } from '@/types/project'
import './NewProjectDialog.css'

export interface FormData {
  name: string
  author: string
  description: string
  category: string
  platform: string
  savePath: string
  encrypted: boolean
  password: string
  confirmPassword: string
}

export interface NewProjectDialogViewProps {
  /** 是否显示对话框 */
  isOpen: boolean
  /** 表单数据 */
  formData: FormData
  /** 自定义分类 */
  customCategory: string
  /** 是否显示自定义分类输入框 */
  showCategoryInput: boolean
  /** 硬件平台列表 */
  platforms: HardwarePlatformConfig[]
  /** 平台列表加载状态 */
  platformsLoading: boolean
  /** 平台列表错误信息 */
  platformsError: string | null
  /** 提交状态 */
  isSubmitting: boolean
  /** 错误信息 */
  error: string | null
  /** 字段错误 */
  errors: Record<string, string>
  /** 关闭回调 */
  onClose: () => void
  /** 输入变化回调 */
  onInputChange: (field: keyof FormData, value: string | boolean) => void
  /** 自定义分类变化回调 */
  onCustomCategoryChange: (value: string) => void
  /** 分类变化回调 */
  onCategoryChange: (value: string) => void
  /** 选择文件夹回调 */
  onSelectFolder: () => void
  /** 提交回调 */
  onSubmit: (e: React.FormEvent) => void
  /** 额外的CSS类名 */
  className?: string
}

/**
 * NewProjectDialogView 新建工程对话框视图组件（纯展示）
 * 紧凑工业风格设计
 *
 * 功能：
 * - 工程信息表单（名称、作者、分类、平台等）
 * - 支持加密设置
 * - 表单验证显示
 * - 加载状态显示
 *
 * 设计模式：
 * - 纯展示组件，不直接访问 store 或 API
 * - 所有数据和回调通过 props 传入
 * - 可在任何上下文中重用
 */
const NewProjectDialogView: React.FC<NewProjectDialogViewProps> = ({
  isOpen,
  formData,
  customCategory,
  showCategoryInput,
  platforms,
  platformsLoading,
  platformsError,
  isSubmitting,
  error,
  errors,
  onClose,
  onInputChange,
  onCustomCategoryChange,
  onCategoryChange,
  onSelectFolder,
  onSubmit,
  className = '',
}) => {
  if (!isOpen) {
    return null
  }

  return (
    <div className="dialog-backdrop" onClick={onClose}>
      <div className="dialog-container" onClick={(e) => e.stopPropagation()}>
        <div className="dialog-header">
          <span className="dialog-title">新建工程</span>
          <button className="dialog-close" onClick={onClose} disabled={isSubmitting}>
            ✕
          </button>
        </div>

        <form onSubmit={onSubmit} className="dialog-form">
          {/* 工程名称 - 全宽 */}
          <div className="form-row">
            <div className="form-field form-field--full">
              <label className="form-label">
                工程名称 <span className="required">*</span>
              </label>
              <input
                type="text"
                className={`form-input ${errors.name ? 'form-input--error' : ''}`}
                value={formData.name}
                onChange={(e) => onInputChange('name', e.target.value)}
                placeholder="输入工程名称"
                disabled={isSubmitting}
                autoFocus
                maxLength={50}
              />
              {errors.name && <span className="form-error">{errors.name}</span>}
            </div>
          </div>

          {/* 工程作者 + 工程分类 - 两列 */}
          <div className="form-row">
            <div className="form-field form-field--half">
              <label className="form-label">工程作者</label>
              <input
                type="text"
                className={`form-input ${errors.author ? 'form-input--error' : ''}`}
                value={formData.author}
                onChange={(e) => onInputChange('author', e.target.value)}
                placeholder="可选"
                disabled={isSubmitting}
                maxLength={30}
              />
              {errors.author && <span className="form-error">{errors.author}</span>}
            </div>

            <div className="form-field form-field--half">
              <label className="form-label">
                工程分类 <span className="required">*</span>
              </label>
              {showCategoryInput ? (
                <div className="form-input-group">
                  <input
                    type="text"
                    className={`form-input ${errors.category ? 'form-input--error' : ''}`}
                    value={customCategory}
                    onChange={(e) => onCustomCategoryChange(e.target.value)}
                    placeholder="输入新分类名称"
                    disabled={isSubmitting}
                  />
                  <button
                    type="button"
                    className="form-input-btn"
                    onClick={() => {
                      if (customCategory.trim()) {
                        onInputChange('category', customCategory.trim())
                        onCategoryChange(ProjectCategory.CATEGORY_1)
                      } else {
                        onCategoryChange(ProjectCategory.CATEGORY_1)
                      }
                    }}
                    disabled={isSubmitting}
                  >
                    ✓
                  </button>
                </div>
              ) : (
                <>
                  <select
                    className={`form-select ${errors.category ? 'form-input--error' : ''}`}
                    value={formData.category}
                    onChange={(e) => onCategoryChange(e.target.value)}
                    disabled={isSubmitting}
                  >
                    <option value={ProjectCategory.CATEGORY_1}>{ProjectCategory.CATEGORY_1}</option>
                    <option value={ProjectCategory.CATEGORY_2}>{ProjectCategory.CATEGORY_2}</option>
                    <option value={ProjectCategory.CUSTOM}>{ProjectCategory.CUSTOM}</option>
                  </select>
                  {errors.category && <span className="form-error">{errors.category}</span>}
                </>
              )}
            </div>
          </div>

          {/* 运行平台 + 保存位置 - 两列 */}
          <div className="form-row">
            <div className="form-field form-field--half">
              <label className="form-label">
                运行平台 <span className="required">*</span>
              </label>
              {platformsLoading ? (
                <select className="form-select" disabled={true}>
                  <option>加载中...</option>
                </select>
              ) : platformsError ? (
                <select className={`form-select ${errors.platform ? 'form-input--error' : ''}`} disabled={true}>
                  <option>加载失败</option>
                </select>
              ) : (
                <>
                  <select
                    className={`form-select ${errors.platform ? 'form-input--error' : ''}`}
                    value={formData.platform}
                    onChange={(e) => onInputChange('platform', e.target.value)}
                    disabled={isSubmitting}
                  >
                    {platforms.map((platform) => (
                      <option key={platform.id} value={platform.id}>
                        {platform.name}
                      </option>
                    ))}
                  </select>
                  {errors.platform && <span className="form-error">{errors.platform}</span>}
                </>
              )}
            </div>

            <div className="form-field form-field--half">
              <label className="form-label">
                保存位置 <span className="required">*</span>
              </label>
              <div className="form-input-group">
                <input
                  type="text"
                  className={`form-input ${errors.savePath ? 'form-input--error' : ''}`}
                  value={formData.savePath}
                  onChange={(e) => onInputChange('savePath', e.target.value)}
                  placeholder="选择保存位置"
                  disabled={isSubmitting}
                  readOnly
                />
                <button type="button" className="form-input-btn" onClick={onSelectFolder} disabled={isSubmitting}>
                  …
                </button>
              </div>
              {errors.savePath && <span className="form-error">{errors.savePath}</span>}
            </div>
          </div>

          {/* 工程描述 - 全宽 */}
          <div className="form-row">
            <div className="form-field form-field--full">
              <label className="form-label">工程描述</label>
              <textarea
                className={`form-textarea ${errors.description ? 'form-input--error' : ''}`}
                value={formData.description}
                onChange={(e) => onInputChange('description', e.target.value)}
                placeholder="输入工程描述（可选）"
                disabled={isSubmitting}
                rows={2}
                maxLength={500}
              />
              {errors.description && <span className="form-error">{errors.description}</span>}
            </div>
          </div>

          {/* 工程加密 - 全宽 */}
          <div className="form-row">
            <div className="form-field form-field--full">
              <label className="form-checkbox-label">
                <input
                  type="checkbox"
                  className="form-checkbox"
                  checked={formData.encrypted}
                  onChange={(e) => onInputChange('encrypted', e.target.checked)}
                  disabled={isSubmitting}
                />
                <span>工程加密</span>
              </label>
            </div>
          </div>

          {/* 密码输入 - 加密时显示，两列 */}
          {formData.encrypted && (
            <div className="form-row">
              <div className="form-field form-field--half">
                <label className="form-label">
                  密码 <span className="required">*</span>
                </label>
                <input
                  type="password"
                  className={`form-input ${errors.password ? 'form-input--error' : ''}`}
                  value={formData.password}
                  onChange={(e) => onInputChange('password', e.target.value)}
                  placeholder="至少6个字符"
                  disabled={isSubmitting}
                />
                {errors.password && <span className="form-error">{errors.password}</span>}
              </div>

              <div className="form-field form-field--half">
                <label className="form-label">
                  确认密码 <span className="required">*</span>
                </label>
                <input
                  type="password"
                  className={`form-input ${errors.confirmPassword ? 'form-input--error' : ''}`}
                  value={formData.confirmPassword}
                  onChange={(e) => onInputChange('confirmPassword', e.target.value)}
                  placeholder="再次输入"
                  disabled={isSubmitting}
                />
                {errors.confirmPassword && <span className="form-error">{errors.confirmPassword}</span>}
              </div>
            </div>
          )}

          {/* 错误提示 */}
          {error && <div className="dialog-error">{error}</div>}

          {/* 操作按钮 */}
          <div className="dialog-actions">
            <button type="button" className="btn btn--secondary" onClick={onClose} disabled={isSubmitting}>
              取消
            </button>
            <button type="submit" className="btn btn--primary" disabled={isSubmitting}>
              {isSubmitting ? '创建中...' : '创建'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default NewProjectDialogView
