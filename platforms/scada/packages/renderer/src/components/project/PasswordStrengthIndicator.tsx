/**
 * 密码强度指示器组件
 * 显示密码强度并给出改进建议
 */

import React from 'react'
import './PasswordStrengthIndicator.css'

/**
 * 密码强度类型
 */
export type PasswordStrength = 'weak' | 'medium' | 'strong'

interface PasswordStrengthIndicatorProps {
  strength: PasswordStrength
}

/**
 * PasswordStrengthIndicator 组件
 */
const PasswordStrengthIndicator: React.FC<PasswordStrengthIndicatorProps> = ({
  strength
}) => {
  // 获取强度配置
  const getStrengthConfig = () => {
    switch (strength) {
      case 'weak':
        return {
          label: '弱',
          color: '#f44336',
          percent: 33,
          tips: ['建议：至少8个字符', '建议：包含大小写字母', '建议：包含数字']
        }
      case 'medium':
        return {
          label: '中',
          color: '#ff9800',
          percent: 66,
          tips: ['建议：使用12个字符以上', '建议：包含特殊字符']
        }
      case 'strong':
        return {
          label: '强',
          color: '#4caf50',
          percent: 100,
          tips: ['密码强度很好']
        }
    }
  }

  const config = getStrengthConfig()

  return (
    <div className="password-strength">
      <div className="password-strength__header">
        <span className="password-strength__label">密码强度：</span>
        <span
          className="password-strength__value"
          style={{ color: config.color }}
        >
          {config.label}
        </span>
      </div>

      {/* 强度条 */}
      <div className="password-strength__bar">
        <div
          className="password-strength__bar-fill"
          style={{
            width: `${config.percent}%`,
            backgroundColor: config.color
          }}
        />
      </div>

      {/* 改进建议 */}
      {strength !== 'strong' && (
        <ul className="password-strength__tips">
          {config.tips.map((tip, index) => (
            <li key={index} className="password-strength__tip">
              {tip}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default PasswordStrengthIndicator
