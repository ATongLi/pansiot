/**
 * 最近工程状态管理
 * 管理最近工程列表、筛选、搜索、排序
 */

import { makeAutoObservable, runInAction, computed } from 'mobx'
import type { RecentProject } from '@/types/project'
import { projectApi } from '@/api/projectApi'
import { formatRelativeTime } from '@/utils/dateFormat'

/**
 * 排序方式类型
 */
type SortBy = 'lastOpened' | 'name' | 'createdAt'
type SortOrder = 'asc' | 'desc'

/**
 * RecentProjectsStore 最近工程状态管理
 */
class RecentProjectsStore {
  // 原始数据
  projects: RecentProject[] = []

  // 分类筛选
  selectedCategory: string = 'all'

  // 搜索关键词
  searchQuery: string = ''

  // 排序方式
  sortBy: SortBy = 'lastOpened'
  sortOrder: SortOrder = 'desc'

  // 加载状态
  isLoading: boolean = false
  error: string = ''

  constructor() {
    makeAutoObservable(this)
  }

  /**
   * 加载最近工程列表
   */
  async loadRecentProjects(): Promise<void> {
    runInAction(() => {
      this.isLoading = true
      this.error = ''
    })

    try {
      const response = await projectApi.getRecentProjects()

      if (response.success && response.data) {
        runInAction(() => {
          this.projects = response.data!
          this.isLoading = false
        })
      } else {
        runInAction(() => {
          this.error = response.message || '加载最近工程失败'
          this.isLoading = false
        })
      }
    } catch (err) {
      runInAction(() => {
        this.error = err instanceof Error ? err.message : '加载最近工程出错'
        this.isLoading = false
      })
    }
  }

  /**
   * 设置分类筛选
   */
  setCategory(category: string): void {
    runInAction(() => {
      this.selectedCategory = category
    })
  }

  /**
   * 设置搜索关键词
   */
  setSearchQuery(query: string): void {
    runInAction(() => {
      this.searchQuery = query
    })
  }

  /**
   * 设置排序方式
   */
  setSortBy(sortBy: SortBy): void {
    runInAction(() => {
      this.sortBy = sortBy
    })
  }

  /**
   * 切换排序方向
   */
  toggleSortOrder(): void {
    runInAction(() => {
      this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc'
    })
  }

  /**
   * 移除工程
   */
  async removeProject(projectId: string): Promise<void> {
    try {
      const response = await projectApi.removeRecentProject(projectId)

      if (response.success) {
        runInAction(() => {
          this.projects = this.projects.filter(p => p.projectId !== projectId)
        })
      }
    } catch (err) {
      console.error('移除工程失败:', err)
    }
  }

  /**
   * 清空列表
   */
  clear(): void {
    runInAction(() => {
      this.projects = []
      this.selectedCategory = 'all'
      this.searchQuery = ''
      this.sortBy = 'lastOpened'
      this.sortOrder = 'desc'
      this.error = ''
    })
  }

  // Computed属性

  /**
   * 所有分类列表（含数量统计）
   */
  get categories(): Array<{ name: string; value: string; count: number }> {
    const categoryMap = new Map<string, number>()

    // 统计每个分类的工程数量
    this.projects.forEach(project => {
      const category = project.category || '未分类'
      categoryMap.set(category, (categoryMap.get(category) || 0) + 1)
    })

    // 转换为数组
    const result = Array.from(categoryMap.entries()).map(([name, count]) => ({
      name,
      value: name,
      count,
    }))

    // 添加"全部"选项
    return [
      { name: '全部', value: 'all', count: this.projects.length },
      ...result.sort((a, b) => a.name.localeCompare(b.name)),
    ]
  }

  /**
   * 过滤后的工程列表
   */
  get filteredProjects(): RecentProject[] {
    let result = [...this.projects]

    // 分类筛选
    if (this.selectedCategory !== 'all') {
      result = result.filter(p => p.category === this.selectedCategory)
    }

    // 搜索筛选
    if (this.searchQuery.trim()) {
      const query = this.searchQuery.toLowerCase().trim()
      result = result.filter(p => p.name.toLowerCase().includes(query))
    }

    return result
  }

  /**
   * 排序后的工程列表
   */
  get sortedProjects(): RecentProject[] {
    const result = [...this.filteredProjects]

    result.sort((a, b) => {
      let comparison = 0

      switch (this.sortBy) {
        case 'lastOpened':
          comparison = a.lastOpenedDate.getTime() - b.lastOpenedDate.getTime()
          break
        case 'name':
          comparison = a.name.localeCompare(b.name, 'zh-CN')
          break
        case 'createdAt':
          comparison = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
          break
      }

      return this.sortOrder === 'asc' ? comparison : -comparison
    })

    return result
  }

  /**
   * 显示用的工程列表（带格式化的时间）
   */
  get displayProjects(): Array<RecentProject & { lastOpened: string }> {
    return this.sortedProjects.map(project => ({
      ...project,
      lastOpened: formatRelativeTime(project.lastOpenedDate),
    }))
  }

  /**
   * 工程总数
   */
  get totalCount(): number {
    return this.projects.length
  }

  /**
   * 筛选后的数量
   */
  get filteredCount(): number {
    return this.filteredProjects.length
  }
}

// 导出类和单例
export { RecentProjectsStore }
export const recentProjectsStore = new RecentProjectsStore()
