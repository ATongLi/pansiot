import React, { useState } from 'react'
import { ACTION_ICONS } from '@/constants/icons'
import { getElectronAPI } from '@/utils/electron'
import NewProjectDialog from './NewProjectDialog'
import './ActionButtons.css'

interface ActionCard {
  id: string
  icon: string
  title: string
  description: string
  onClick: () => void
}

interface ActionButtonsProps {
  onOpenProject?: (projectName: string, projectPath: string) => void
}

const ActionButtons: React.FC<ActionButtonsProps> = ({ onOpenProject }) => {
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  const actions: ActionCard[] = [
    {
      id: 'new',
      icon: ACTION_ICONS.newProject,
      title: '新建工程',
      description: '从零开始创建新的工程配置',
      onClick: () => setIsDialogOpen(true),
    },
    {
      id: 'open',
      icon: ACTION_ICONS.openProject,
      title: '从文件打开',
      description: '打开已保存的工程文件',
      onClick: () => handleOpenProject(),
    },
    {
      id: 'copy',
      icon: ACTION_ICONS.copyProject,
      title: '复制工程',
      description: '基于现有工程创建副本',
      onClick: () => console.log('Copy project - TODO'),
    },
  ]

  /**
   * 处理打开工程文件
   */
  const handleOpenProject = async () => {
    try {
      const electronAPI = getElectronAPI()
      const filePath = await electronAPI.dialog.selectOpenPath({
        title: '选择工程文件',
        filters: [
          { name: 'PanTools工程文件', extensions: ['pant'] },
          { name: '所有文件', extensions: ['*'] },
        ],
      })

      if (filePath) {
        console.log('打开工程文件:', filePath)

        // 读取工程文件
        const project = await electronAPI.file.readProject(filePath)

        if (project && onOpenProject) {
          // 调用回调，在 App.tsx 中创建编辑器标签页
          const fileName = filePath.split(/[/\\]/).pop() || filePath
          const projectName = fileName.replace(/\.pant$/, '')
          onOpenProject(projectName || project.name, filePath)
        }
      }
    } catch (error: any) {
      console.error('打开工程失败:', error)
    }
  }

  return (
    <>
      <div className="action-cards">
        {actions.map((action) => (
          <div
            key={action.id}
            className="action-card"
            onClick={action.onClick}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                action.onClick()
              }
            }}
          >
            <div
              className="action-card__icon"
              dangerouslySetInnerHTML={{ __html: action.icon }}
            />
            <div className="action-card__content">
              <h3 className="action-card__title">{action.title}</h3>
              <p className="action-card__description">{action.description}</p>
            </div>
          </div>
        ))}
      </div>

      <NewProjectDialog
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
      />
    </>
  )
}

export default ActionButtons
