/**
 * 新建工程对话框
 * 提供工程创建表单和验证功能
 */

import React, { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { ProjectCategory, HardwarePlatform } from '@/types/project'
import { projectApi } from '@/api/projectApi'
import { getElectronAPI } from '@/utils/electron'
import PasswordStrengthIndicator from './PasswordStrengthIndicator'
import './NewProjectDialog.css'

interface NewProjectDialogProps {
  onClose: () => void
  onProjectCreated: (projectId: string) => void
}

/**
 * 密码强度类型
 */
type PasswordStrength = 'weak' | 'medium' | 'strong'

/**
 * 表单验证错误类型
 */
interface FormErrors {
  name?: string
  author?: string
  description?: string
  category?: string
  platform?: string
  savePath?: string
  password?: string
  confirmPassword?: string
}

/**
 * NewProjectDialog 组件
 */
const NewProjectDialog: React.FC<NewProjectDialogProps> = observer(
  ({ onClose, onProjectCreated }) => {
    // 表单数据状态
    const [formData, setFormData] = useState({
      name: '',
      author: '',
      description: '',
      category: ProjectCategory.CATEGORY_1,
      platform: HardwarePlatform.HMI_MODEL_1,
      encrypted: false,
      password: '',
      confirmPassword: '',
      savePath: ''
    })

    // 验证错误状态
    const [errors, setErrors] = useState<FormErrors>({})

    // 自定义分类输入状态
    const [customCategory, setCustomCategory] = useState('')
    const [showCustomCategory, setShowCustomCategory] = useState(false)

    // 密码可见性状态
    const [showPassword, setShowPassword] = useState(false)
    const [showConfirmPassword, setShowConfirmPassword] = useState(false)

    // 提交中状态
    const [isSubmitting, setIsSubmitting] = useState(false)

    // 处理表单字段变化
    const handleFieldChange = (
      field: string,
      value: string | boolean | ProjectCategory | HardwarePlatform
    ) => {
      setFormData(prev => ({ ...prev, [field]: value }))

      // 清除该字段的错误提示
      if (errors[field as keyof FormErrors]) {
        setErrors(prev => ({ ...prev, [field]: undefined }))
      }
    }

    // 验证表单
    const validateForm = (): boolean => {
      const newErrors: FormErrors = {}

      // 工程名称验证
      if (!formData.name.trim()) {
        newErrors.name = '工程名称不能为空'
      } else if (formData.name.length > 50) {
        newErrors.name = '工程名称不能超过50个字符'
      }

      // 工程作者验证
      if (formData.author && formData.author.length > 30) {
        newErrors.author = '作者名称不能超过30个字符'
      }

      // 工程描述验证
      if (formData.description && formData.description.length > 500) {
        newErrors.description = '工程描述不能超过500个字符'
      }

      // 保存位置验证
      if (!formData.savePath) {
        newErrors.savePath = '请选择工程保存位置'
      }

      // 加密工程密码验证
      if (formData.encrypted) {
        if (!formData.password) {
          newErrors.password = '请设置密码'
        } else if (formData.password.length < 6) {
          newErrors.password = '密码至少6个字符'
        }

        if (formData.password !== formData.confirmPassword) {
          newErrors.confirmPassword = '两次输入的密码不一致'
        }
      }

      setErrors(newErrors)
      return Object.keys(newErrors).length === 0
    }

    // 处理文件夹选择
    const handleSelectFolder = async () => {
      try {
        const electronAPI = getElectronAPI()
        const filePath = await electronAPI.dialog.selectSavePath({
          title: '选择工程保存位置',
          defaultPath: formData.name ? `${formData.name}.pant` : undefined,
          filters: [
            { name: 'PanTools工程文件', extensions: ['pant'] },
            { name: '所有文件', extensions: ['*'] },
          ],
        })

        if (filePath) {
          handleFieldChange('savePath', filePath)
        }
      } catch (error) {
        console.error('选择文件夹失败:', error)
      }
    }

    // 处理自定义分类
    const handleCustomCategory = () => {
      if (customCategory.trim()) {
        setFormData(prev => ({ ...prev, category: customCategory.trim() }))
        setShowCustomCategory(false)
        setCustomCategory('')
      }
    }

    // 计算密码强度
    const calculatePasswordStrength = (): PasswordStrength => {
      const password = formData.password
      if (!password) return 'weak'

      let strength = 0
      if (password.length >= 8) strength++
      if (password.length >= 12) strength++
      if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
      if (/\d/.test(password)) strength++
      if (/[^a-zA-Z0-9]/.test(password)) strength++

      if (strength <= 2) return 'weak'
      if (strength <= 3) return 'medium'
      return 'strong'
    }

    // 处理表单提交
    const handleSubmit = async (e: React.FormEvent) => {
      e.preventDefault()

      if (!validateForm()) {
        return
      }

      setIsSubmitting(true)

      try {
        // 调用 API 创建工程
        const response = await projectApi.createProject(formData)

        if (response.success && response.data) {
          // 工程创建成功
          onProjectCreated(response.data.projectId)
          onClose()
        } else {
          // 显示错误信息
          console.error('创建工程失败:', response.message)
        }
      } catch (error) {
        console.error('创建工程出错:', error)
      } finally {
        setIsSubmitting(false)
      }
    }

    // 处理取消
    const handleCancel = () => {
      onClose()
    }

    // 处理 ESC 键
    useEffect(() => {
      const handleEsc = (e: KeyboardEvent) => {
        if (e.key === 'Escape') {
          handleCancel()
        }
      }

      window.addEventListener('keydown', handleEsc)
      return () => window.removeEventListener('keydown', handleEsc)
    }, [])

    return (
      <div className="new-project-dialog-overlay" onClick={handleCancel}>
        <div className="new-project-dialog" onClick={e => e.stopPropagation()}>
          {/* 头部 */}
          <div className="new-project-dialog__header">
            <h2 className="new-project-dialog__title">新建工程</h2>
            <button
              className="new-project-dialog__close"
              onClick={handleCancel}
              aria-label="关闭"
            >
              ✕
            </button>
          </div>

          {/* 表单内容 */}
          <form onSubmit={handleSubmit} className="new-project-dialog__body">
            {/* 工程名称 */}
            <div className="form-group">
              <label className="form-group__label">
                工程名称 <span className="form-group__required">*</span>
              </label>
              <input
                type="text"
                className={`form-group__input ${errors.name ? 'form-group__input--error' : ''}`}
                value={formData.name}
                onChange={e => handleFieldChange('name', e.target.value)}
                placeholder="请输入工程名称"
                maxLength={50}
              />
              {errors.name && <span className="form-group__error">{errors.name}</span>}
            </div>

            {/* 工程作者 */}
            <div className="form-group">
              <label className="form-group__label">工程作者</label>
              <input
                type="text"
                className={`form-group__input ${errors.author ? 'form-group__input--error' : ''}`}
                value={formData.author}
                onChange={e => handleFieldChange('author', e.target.value)}
                placeholder="请输入作者名称（可选）"
                maxLength={30}
              />
              {errors.author && <span className="form-group__error">{errors.author}</span>}
            </div>

            {/* 工程描述 */}
            <div className="form-group">
              <label className="form-group__label">工程描述</label>
              <textarea
                className={`form-group__textarea ${errors.description ? 'form-group__input--error' : ''}`}
                value={formData.description}
                onChange={e => handleFieldChange('description', e.target.value)}
                placeholder="请输入工程描述（可选）"
                rows={3}
                maxLength={500}
              />
              {errors.description && <span className="form-group__error">{errors.description}</span>}
            </div>

            {/* 工程分类 */}
            <div className="form-group">
              <label className="form-group__label">工程分类</label>
              <div className="form-group__select-wrapper">
                <select
                  className="form-group__select"
                  value={formData.category}
                  onChange={e => {
                    if (e.target.value === 'custom') {
                      setShowCustomCategory(true)
                    } else {
                      handleFieldChange('category', e.target.value as ProjectCategory)
                    }
                  }}
                >
                  <option value={ProjectCategory.CATEGORY_1}>分类1</option>
                  <option value={ProjectCategory.CATEGORY_2}>分类2</option>
                  <option value="custom">自定义分类...</option>
                </select>
              </div>
              {showCustomCategory && (
                <div className="form-group__custom-category">
                  <input
                    type="text"
                    className="form-group__input"
                    value={customCategory}
                    onChange={e => setCustomCategory(e.target.value)}
                    placeholder="请输入自定义分类名称"
                    onBlur={handleCustomCategory}
                    autoFocus
                  />
                  <button
                    type="button"
                    className="form-group__custom-category-btn"
                    onClick={handleCustomCategory}
                  >
                    确定
                  </button>
                </div>
              )}
            </div>

            {/* 硬件平台 */}
            <div className="form-group">
              <label className="form-group__label">
                硬件平台 <span className="form-group__required">*</span>
              </label>
              <select
                className="form-group__select"
                value={formData.platform}
                onChange={e => handleFieldChange('platform', e.target.value as HardwarePlatform)}
              >
                <option value={HardwarePlatform.HMI_MODEL_1}>HMI型号1</option>
                <option value={HardwarePlatform.HMI_MODEL_2}>HMI型号2</option>
                <option value={HardwarePlatform.GATEWAY_MODEL_1}>网关型号1</option>
              </select>
            </div>

            {/* 保存位置 */}
            <div className="form-group">
              <label className="form-group__label">
                保存位置 <span className="form-group__required">*</span>
              </label>
              <div className="form-group__path-selector">
                <input
                  type="text"
                  className={`form-group__input ${errors.savePath ? 'form-group__input--error' : ''}`}
                  value={formData.savePath}
                  onChange={e => handleFieldChange('savePath', e.target.value)}
                  placeholder="请选择工程保存位置"
                  readOnly
                />
                <button
                  type="button"
                  className="form-group__browse-btn"
                  onClick={handleSelectFolder}
                >
                  浏览...
                </button>
              </div>
              {errors.savePath && <span className="form-group__error">{errors.savePath}</span>}
            </div>

            {/* 工程加密 */}
            <div className="form-group">
              <label className="form-group__checkbox">
                <input
                  type="checkbox"
                  checked={formData.encrypted}
                  onChange={e => handleFieldChange('encrypted', e.target.checked)}
                />
                <span>启用工程加密</span>
              </label>
            </div>

            {/* 密码输入（加密工程） */}
            {formData.encrypted && (
              <>
                <div className="form-group">
                  <label className="form-group__label">
                    密码 <span className="form-group__required">*</span>
                  </label>
                  <div className="form-group__password-input">
                    <input
                      type={showPassword ? 'text' : 'password'}
                      className={`form-group__input ${errors.password ? 'form-group__input--error' : ''}`}
                      value={formData.password}
                      onChange={e => handleFieldChange('password', e.target.value)}
                      placeholder="请设置密码（至少6个字符）"
                    />
                    <button
                      type="button"
                      className="form-group__password-toggle"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? '👁️' : '👁️‍🗨️'}
                    </button>
                  </div>
                  {errors.password && <span className="form-group__error">{errors.password}</span>}
                  {formData.password && (
                    <PasswordStrengthIndicator strength={calculatePasswordStrength()} />
                  )}
                </div>

                <div className="form-group">
                  <label className="form-group__label">
                    确认密码 <span className="form-group__required">*</span>
                  </label>
                  <div className="form-group__password-input">
                    <input
                      type={showConfirmPassword ? 'text' : 'password'}
                      className={`form-group__input ${errors.confirmPassword ? 'form-group__input--error' : ''}`}
                      value={formData.confirmPassword}
                      onChange={e => handleFieldChange('confirmPassword', e.target.value)}
                      placeholder="请再次输入密码"
                    />
                    <button
                      type="button"
                      className="form-group__password-toggle"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    >
                      {showConfirmPassword ? '👁️' : '👁️‍🗨️'}
                    </button>
                  </div>
                  {errors.confirmPassword && (
                    <span className="form-group__error">{errors.confirmPassword}</span>
                  )}
                </div>
              </>
            )}
          </form>

          {/* 底部按钮 */}
          <div className="new-project-dialog__footer">
            <button
              type="button"
              className="btn btn--secondary"
              onClick={handleCancel}
              disabled={isSubmitting}
            >
              取消
            </button>
            <button
              type="button"
              className="btn btn--primary"
              onClick={handleSubmit}
              disabled={isSubmitting}
            >
              {isSubmitting ? '创建中...' : '确定'}
            </button>
          </div>
        </div>
      </div>
    )
  }
)

export default NewProjectDialog
