import { useState, useCallback } from 'react'

export interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  custom?: (value: any) => string | null
}

export type ValidationRules<T> = Partial<Record<keyof T, ValidationRule>>

export type ValidationErrors<T> = Partial<Record<keyof T, string>>

/**
 * 表单验证 Hook
 * 提供通用的表单验证逻辑
 *
 * @param rules - 验证规则配置
 * @returns 验证函数和错误状态
 */
export function useFormValidation<T extends Record<string, any>>(
  rules: ValidationRules<T>
) {
  const [errors, setErrors] = useState<ValidationErrors<T>>({})

  /**
   * 验证单个字段
   */
  const validateField = useCallback(
    (field: keyof T, value: any): string | null => {
      const fieldRules = rules[field]
      if (!fieldRules) return null

      // 必填验证
      if (fieldRules.required && (!value || (typeof value === 'string' && !value.trim()))) {
        return '此项为必填'
      }

      // 跳过空值的其他验证
      if (!value || (typeof value === 'string' && !value.trim())) {
        return null
      }

      // 最小长度验证
      if (fieldRules.minLength && typeof value === 'string' && value.length < fieldRules.minLength) {
        return `长度不能少于 ${fieldRules.minLength} 个字符`
      }

      // 最大长度验证
      if (fieldRules.maxLength && typeof value === 'string' && value.length > fieldRules.maxLength) {
        return `长度不能超过 ${fieldRules.maxLength} 个字符`
      }

      // 正则验证
      if (fieldRules.pattern && typeof value === 'string' && !fieldRules.pattern.test(value)) {
        return '格式不正确'
      }

      // 自定义验证
      if (fieldRules.custom) {
        return fieldRules.custom(value)
      }

      return null
    },
    [rules]
  )

  /**
   * 验证整个表单
   */
  const validateForm = useCallback(
    (data: T): boolean => {
      const newErrors: ValidationErrors<T> = {}
      let hasError = false

      for (const field in rules) {
        const error = validateField(field, data[field])
        if (error) {
          newErrors[field] = error
          hasError = true
        }
      }

      setErrors(newErrors)
      return !hasError
    },
    [rules, validateField]
  )

  /**
   * 清除所有错误
   */
  const clearErrors = useCallback(() => {
    setErrors({})
  }, [])

  /**
   * 清除单个字段的错误
   */
  const clearFieldError = useCallback((field: keyof T) => {
    setErrors((prev) => {
      const newErrors = { ...prev }
      delete newErrors[field]
      return newErrors
    })
  }, [])

  return {
    errors,
    validateField,
    validateForm,
    clearErrors,
    clearFieldError,
  }
}
