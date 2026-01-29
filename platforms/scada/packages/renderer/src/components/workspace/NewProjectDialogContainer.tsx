import React, { useState, useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { projectStore } from '@store/projectStore'
import { getElectronAPI } from '@/utils/electron'
import { ProjectCategory, HardwarePlatform } from '@/types/project'
import { HardwarePlatformConfig } from '@/types/project'
import { platformApi } from '@/api/platformApi'
import { useFormValidation } from '@/hooks/useFormValidation'
import NewProjectDialogView, { type FormData } from './NewProjectDialogView'

export interface NewProjectDialogContainerProps {
  /** 是否显示对话框 */
  isOpen: boolean
  /** 关闭回调 */
  onClose: () => void
}

/**
 * NewProjectDialogContainer 新建工程对话框容器组件
 * 使用 Container/Presenter 模式重构
 *
 * 职责：
 * - 连接 projectStore 创建工程
 * - 调用 API 获取硬件平台列表
 * - 管理表单状态
 * - 处理表单验证
 * - 处理表单提交
 *
 * 设计模式：
 * - Container 组件：负责状态管理和业务逻辑
 * - View 组件：负责纯 UI 渲染
 * - 分离关注点，提高可测试性和可重用性
 */
const NewProjectDialogContainer: React.FC<NewProjectDialogContainerProps> = observer(({ isOpen, onClose }) => {
  // 表单状态
  const [formData, setFormData] = useState<FormDataType>({
    name: '',
    author: '',
    description: '',
    category: ProjectCategory.CATEGORY_1,
    platform: HardwarePlatform.BOX1,
    savePath: '',
    encrypted: false,
    password: '',
    confirmPassword: '',
  })

  // 自定义分类状态
  const [customCategory, setCustomCategory] = useState('')
  const [showCategoryInput, setShowCategoryInput] = useState(false)

  // 硬件平台列表状态
  const [platforms, setPlatforms] = useState<HardwarePlatformConfig[]>([])
  const [platformsLoading, setPlatformsLoading] = useState(false)
  const [platformsError, setPlatformsError] = useState<string | null>(null)

  // 提交状态
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // 表单验证规则
  const validationRules = {
    name: {
      required: true,
      maxLength: 50,
      custom: (value: string) => {
        if (!value.trim()) return '工程名称不能为空'
        return null
      },
    },
    author: {
      maxLength: 30,
    },
    description: {
      maxLength: 500,
    },
    category: {
      required: true,
      custom: (value: string) => {
        if (!value) return '请选择或输入工程分类'
        return null
      },
    },
    platform: {
      required: true,
      custom: (value: string) => {
        if (!value) return '请选择或输入运行平台'
        return null
      },
    },
    savePath: {
      required: true,
      custom: (value: string) => {
        if (!value.trim()) return '请选择工程保存位置'
        return null
      },
    },
    password: {
      custom: (value: string) => {
        if (formData.encrypted) {
          if (!value) return '请输入密码'
          if (value.length < 6) return '密码长度至少6个字符'
        }
        return null
      },
    },
    confirmPassword: {
      custom: (value: string) => {
        if (formData.encrypted && value !== formData.password) {
          return '两次输入的密码不一致'
        }
        return null
      },
    },
  }

  const { errors, validateForm, clearErrors } = useFormValidation(validationRules)

  /**
   * 获取硬件平台列表
   */
  useEffect(() => {
    const fetchPlatforms = async () => {
      if (!isOpen) return

      setPlatformsLoading(true)
      setPlatformsError(null)

      try {
        const response = await platformApi.getAllPlatforms()
        if (response.success && response.data) {
          setPlatforms(response.data)
          // 设置默认平台（第一个启用的平台）
          if (response.data.length > 0) {
            setFormData((prev) => ({ ...prev, platform: response.data[0]!.id }))
          }
        }
      } catch (err: any) {
        console.error('获取硬件平台列表失败:', err)
        setPlatformsError('获取硬件平台列表失败，请刷新页面重试')
      } finally {
        setPlatformsLoading(false)
      }
    }

    fetchPlatforms()
  }, [isOpen])

  /**
   * 处理输入框变化
   */
  const handleInputChange = (field: keyof FormDataType, value: string | boolean) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  /**
   * 处理工程分类变化
   */
  const handleCategoryChange = (value: string) => {
    if (value === ProjectCategory.CUSTOM) {
      setShowCategoryInput(true)
      setFormData((prev) => ({ ...prev, category: '' }))
    } else {
      setShowCategoryInput(false)
      setCustomCategory('')
      setFormData((prev) => ({ ...prev, category: value }))
    }
  }

  /**
   * 选择保存位置文件夹
   */
  const handleSelectFolder = async () => {
    try {
      const electronAPI = getElectronAPI()
      const folderPath = await electronAPI.dialog.selectSavePath({
        title: '选择工程保存位置',
        filters: [{ name: 'PanTools工程文件', extensions: ['pant'] }],
      })

      if (folderPath) {
        setFormData((prev) => ({ ...prev, savePath: folderPath }))
      }
    } catch (err: any) {
      console.error('选择文件夹失败:', err)
    }
  }

  /**
   * 处理表单提交
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // 验证表单
    if (!validateForm(formData)) {
      return
    }

    setIsSubmitting(true)
    setError(null)

    try {
      await projectStore.createProject({
        metadata: {
          name: formData.name.trim(),
          author: formData.author.trim() || undefined,
          description: formData.description.trim() || undefined,
          category: formData.category,
          platform: formData.platform,
        },
        security: {
          encrypted: formData.encrypted,
          password: formData.encrypted ? formData.password : undefined,
        },
        savePath: formData.savePath,
      })

      // 重置表单
      resetForm()
      onClose()
    } catch (err: any) {
      setError(err.message || '创建工程失败')
    } finally {
      setIsSubmitting(false)
    }
  }

  /**
   * 重置表单
   */
  const resetForm = () => {
    setFormData({
      name: '',
      author: '',
      description: '',
      category: ProjectCategory.CATEGORY_1,
      platform: platforms.length > 0 ? platforms[0].id : HardwarePlatform.BOX1,
      savePath: '',
      encrypted: false,
      password: '',
      confirmPassword: '',
    })
    setCustomCategory('')
    setShowCategoryInput(false)
    clearErrors()
    setError(null)
  }

  const viewProps = {
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
    onInputChange: handleInputChange,
    onCustomCategoryChange: setCustomCategory,
    onCategoryChange: handleCategoryChange,
    onSelectFolder: handleSelectFolder,
    onSubmit: handleSubmit,
  }

  return <NewProjectDialogView {...viewProps} />
})

// Type alias for form data
type FormDataType = FormData

export default NewProjectDialogContainer
