/**
 * 分类筛选组件
 * 提供横向滚动的分类标签筛选
 */

import React from 'react'
import './CategoryFilter.css'

export interface CategoryItem {
  name: string
  count: number
  value: string
}

interface CategoryFilterProps {
  categories: CategoryItem[]
  selectedCategory: string
  onCategoryChange: (category: string) => void
}

/**
 * CategoryFilter 组件
 */
const CategoryFilter: React.FC<CategoryFilterProps> = ({
  categories,
  selectedCategory,
  onCategoryChange
}) => {
  // 处理分类点击
  const handleCategoryClick = (category: CategoryItem) => {
    onCategoryChange(category.value)
  }

  return (
    <div className="category-filter">
      <div className="category-filter__scroll">
        {categories.map(category => (
          <button
            key={category.value}
            className={`category-filter__item ${
              selectedCategory === category.value
                ? 'category-filter__item--active'
                : ''
            }`}
            onClick={() => handleCategoryClick(category)}
          >
            <span className="category-filter__item-name">
              {category.name}
            </span>
            <span className="category-filter__item-count">
              {category.count}
            </span>
          </button>
        ))}
      </div>
    </div>
  )
}

export default CategoryFilter
