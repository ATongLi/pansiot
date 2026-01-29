/**
 * 搜索框组件
 * 提供实时搜索功能，带防抖处理
 */

import React, { useState, useEffect, useRef } from 'react'
import './SearchBox.css'

interface SearchBoxProps {
  placeholder?: string
  onSearch: (query: string) => void
  debounceMs?: number
}

/**
 * SearchBox 组件
 */
const SearchBox: React.FC<SearchBoxProps> = ({
  placeholder = '搜索工程名称...',
  onSearch,
  debounceMs = 100
}) => {
  const [query, setQuery] = useState('')
  const [isFocused, setIsFocused] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const debounceTimerRef = useRef<NodeJS.Timeout>()

  // 处理输入变化（带防抖）
  const handleChange = (value: string) => {
    setQuery(value)

    // 清除之前的定时器
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current)
    }

    // 设置新的防抖定时器
    debounceTimerRef.current = setTimeout(() => {
      onSearch(value)
    }, debounceMs)
  }

  // 清除搜索
  const handleClear = () => {
    setQuery('')
    onSearch('')
    inputRef.current?.focus()
  }

  // 组件卸载时清除定时器
  useEffect(() => {
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
      }
    }
  }, [])

  return (
    <div
      className={`search-box ${isFocused ? 'search-box--focused' : ''}`}
      onFocus={() => setIsFocused(true)}
      onBlur={() => setIsFocused(false)}
    >
      {/* 搜索图标 */}
      <div className="search-box__icon">🔍</div>

      {/* 输入框 */}
      <input
        ref={inputRef}
        type="text"
        className="search-box__input"
        value={query}
        onChange={e => handleChange(e.target.value)}
        placeholder={placeholder}
      />

      {/* 清除按钮 */}
      {query && (
        <button
          type="button"
          className="search-box__clear"
          onClick={handleClear}
          aria-label="清除搜索"
        >
          ✕
        </button>
      )}
    </div>
  )
}

export default SearchBox
