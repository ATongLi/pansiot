/**
 * EditorLayoutContainer 编辑器布局容器组件
 * FE-009-03: 使用 Container/Presenter 模式重构
 *
 * 职责：
 * - 连接多个 stores (getEditorStore, subPageTabStore)
 * - 处理文件操作 (新建、打开、保存工程)
 * - 处理键盘快捷键
 * - 管理侧边栏和画布状态
 * - 监听 Electron 事件
 *
 * 设计模式：
 * - Container 组件：负责状态管理和业务逻辑
 * - View 组件：负责纯 UI 渲染
 * - 分离关注点，提高可测试性和可重用性
 */

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { getEditorStore } from '@/store'
import { subPageTabStore } from '@/store'
import { useEditorKeyboardShortcuts } from '@/hooks/useKeyboardShortcuts'
import EditorLayoutView, { type EditorLayoutViewProps } from './EditorLayoutView'
import { ToolbarCategory } from './toolbar/ToolbarTypes'

export interface EditorLayoutContainerProps {
  className?: string
}

/**
 * EditorLayoutContainer 组件
 */
export const EditorLayoutContainer: React.FC<EditorLayoutContainerProps> = observer(({ className = '' }) => {
  const editorStore = getEditorStore()

  // 启用键盘快捷键
  useEditorKeyboardShortcuts()

  // ==========================================
  // State - 工具栏
  // ==========================================

  /** 子工具栏是否可见 */
  const [subToolbarVisible, setSubToolbarVisible] = React.useState(true)
  /** 当前激活的工具栏分类 */
  const [activeToolbarCategory, setActiveToolbarCategory] = React.useState<ToolbarCategory>(ToolbarCategory.PROJECT)

  // ==========================================
  // Handlers - 文件操作
  // ==========================================

  /**
   * 添加工程标签页
   * IMP-009: 使用 subPageTabStore 添加标签
   */
  const addProjectTab = (projectName: string, projectPath?: string): void => {
    // 使用 SubPageTabStore 添加标签页
    subPageTabStore.addSubTab({
      title: projectName,
      contentType: 'screen',
      dataId: projectPath,
      closable: true,
    })
    console.log('EditorLayout: 添加工程标签', projectName, projectPath || '(未保存)')
  }

  /**
   * 处理新建工程
   */
  const handleNewProject = async (): Promise<void> => {
    try {
      if (typeof window !== 'undefined' && window.electronAPI) {
        const project = await window.electronAPI.file.createProject('新工程')
        if (project) {
          addProjectTab(project.name)
          // 新建工程后弹出保存对话框
          const savePath = await window.electronAPI.file.selectSavePath({
            title: '保存新工程',
            defaultPath: `${project.name}.pant`,
          })

          if (savePath) {
            // 保存工程文件
            await window.electronAPI.file.writeProject(savePath, project)
            console.log('EditorLayout: 新工程已保存到', savePath)
            // 更新标签名称（从路径提取文件名）
            const fileName = savePath.split(/[/\\]/).pop() || savePath
            const projectName = fileName.replace(/\.pant$/, '')
            // IMP-009: 更新当前激活标签的标题
            if (subPageTabStore.activeSubTab) {
              const activeTab = subPageTabStore.subPageTabs.find((t) => t.id === subPageTabStore.activeSubTab)
              if (activeTab) {
                activeTab.title = projectName
              }
            }
          } else {
            console.log('EditorLayout: 用户取消保存')
          }
        }
      }
    } catch (error) {
      console.error('EditorLayout: 新建工程失败', error)
    }
  }

  /**
   * 处理打开工程
   */
  const handleOpenProject = async (): Promise<void> => {
    console.log('EditorLayout: handleOpenProject 开始执行')
    try {
      if (typeof window !== 'undefined' && window.electronAPI) {
        console.log('EditorLayout: electronAPI 可用，调用 selectOpenPath')
        const result = await window.electronAPI.file.selectOpenPath()
        console.log('EditorLayout: selectOpenPath 返回', result)

        if (result && !result.canceled && result.filePath) {
          console.log('EditorLayout: 用户选择了文件', result.filePath)
          const project = await window.electronAPI.file.readProject(result.filePath)
          console.log('EditorLayout: 读取到工程数据', project)

          if (project) {
            console.log('EditorLayout: 准备添加标签页', project.name)
            addProjectTab(project.name, result.filePath)
            console.log('EditorLayout: 标签页已添加')
          } else {
            console.error('EditorLayout: 工程数据为空')
          }
        } else {
          console.log('EditorLayout: 用户取消选择或无有效路径')
        }
      } else {
        console.error('EditorLayout: electronAPI 不可用')
      }
    } catch (error) {
      console.error('EditorLayout: 打开工程失败', error)
    }
  }

  /**
   * 处理保存工程
   */
  const handleSaveProject = async (): Promise<void> => {
    try {
      if (typeof window !== 'undefined' && window.electronAPI) {
        // TODO: 保存当前活动标签的工程数据
        console.log('EditorLayout: 保存工程', subPageTabStore.activeSubTab)
        await window.electronAPI.notification.success('保存成功', '工程已保存')
      }
    } catch (error) {
      console.error('EditorLayout: 保存工程失败', error)
    }
  }

  // ==========================================
  // Handlers - 画布操作
  // ==========================================

  /**
   * 处理组件放置
   */
  const handleDropComponent = (component: any, x: number, y: number): void => {
    console.log('EditorLayout: component dropped', component, 'at', x, y)
    // TODO: 添加组件到画布
  }

  // ==========================================
  // Handlers - 侧边栏操作
  // ==========================================

  const handleSetLeftSidebarTab = (tab: 'project' | 'screen' | 'component'): void => {
    editorStore.setLeftSidebarTab(tab)
  }

  const handleSetRightSidebarTab = (tab: 'property' | 'layer'): void => {
    editorStore.setRightSidebarTab(tab)
  }

  // ==========================================
  // Handlers - 工具栏操作
  // ==========================================

  /**
   * 切换子工具栏可见性
   */
  const handleToggleSubToolbar = (): void => {
    setSubToolbarVisible(!subToolbarVisible)
  }

  /**
   * 设置激活的工具栏分类
   */
  const handleSetActiveToolbarCategory = (category: ToolbarCategory): void => {
    setActiveToolbarCategory(category)
  }

  // ==========================================
  // Handlers - 标签页操作
  // ==========================================

  /**
   * 处理标签页切换
   * IMP-009: 使用 subPageTabStore 切换标签
   */
  const handleTabChange = (tabId: string): void => {
    subPageTabStore.setActiveSubTab(tabId)
    console.log('EditorLayout: tab changed to', tabId)
  }

  /**
   * 处理标签页关闭
   * IMP-009: 使用 subPageTabStore 关闭标签
   */
  const handleTabClose = (tabId: string): void => {
    subPageTabStore.closeSubTab(tabId)
    console.log('EditorLayout: tab closed', tabId)
  }

  /**
   * 处理新建标签页
   * IMP-009: 使用 subPageTabStore 添加新标签
   */
  const handleTabAdd = (): void => {
    const tabCount = subPageTabStore.subPageTabs.length
    subPageTabStore.addSubTab({
      title: `画面 ${tabCount + 1}`,
      contentType: 'screen',
      closable: true,
    })
    console.log('EditorLayout: new tab added')
  }

  // ==========================================
  // Effects - 初始化
  // ==========================================

  useEffect(() => {
    console.log('EditorLayout: mounted')

    // 检查 Electron API 是否可用
    if (typeof window !== 'undefined' && window.electronAPI) {
      console.log('EditorLayout: electronAPI available')

      // 监听窗口事件
      const handleMaximized = () => console.log('Window: maximized')
      const handleUnmaximized = () => console.log('Window: unmaximized')

      // 监听文件操作事件（来自菜单）
      const handleFileNew = () => {
        console.log('EditorLayout: received file:new event')
        handleNewProject()
      }
      const handleFileOpen = () => {
        console.log('EditorLayout: received file:open event')
        handleOpenProject()
      }
      const handleFileSave = () => {
        console.log('EditorLayout: received file:save event')
        handleSaveProject()
      }

      window.electronAPI.on('window:maximized', handleMaximized)
      window.electronAPI.on('window:unmaximized', handleUnmaximized)
      window.electronAPI.on('file:new', handleFileNew)
      window.electronAPI.on('file:open', handleFileOpen)
      window.electronAPI.on('file:save', handleFileSave)

      return () => {
        window.electronAPI.off('window:maximized', handleMaximized)
        window.electronAPI.off('window:unmaximized', handleUnmaximized)
        window.electronAPI.off('file:new', handleFileNew)
        window.electronAPI.off('file:open', handleFileOpen)
        window.electronAPI.off('file:save', handleFileSave)
      }
    }
  }, [])

  // ==========================================
  // 准备 View Props
  // ==========================================

  const viewProps: EditorLayoutViewProps = {
    // 状态
    leftSidebarActiveTab: editorStore.state.leftSidebarActiveTab,
    rightSidebarVisible: editorStore.state.rightSidebarVisible,
    rightSidebarActiveTab: editorStore.state.rightSidebarActiveTab,
    subPageTabs: subPageTabStore.subPageTabs,
    activeSubTab: subPageTabStore.activeSubTab,

    // 工具栏状态
    subToolbarVisible,
    activeToolbarCategory,

    // 回调 - 文件操作
    onNewProject: handleNewProject,
    onOpenProject: handleOpenProject,
    onSaveProject: handleSaveProject,

    // 回调 - 侧边栏操作
    onSetLeftSidebarTab: handleSetLeftSidebarTab,
    onSetRightSidebarTab: handleSetRightSidebarTab,

    // 回调 - 标签页操作
    onTabChange: handleTabChange,
    onTabClose: handleTabClose,
    onTabAdd: handleTabAdd,

    // 回调 - 画布操作
    onDropComponent: handleDropComponent,

    // 回调 - 工具栏操作
    onToggleSubToolbar: handleToggleSubToolbar,
    onSetActiveToolbarCategory: handleSetActiveToolbarCategory,

    className,
  }

  return <EditorLayoutView {...viewProps} />
})

/**
 * Default export
 */
export default EditorLayoutContainer
