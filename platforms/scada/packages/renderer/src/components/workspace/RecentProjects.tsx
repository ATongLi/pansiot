import React from 'react'
import './RecentProjects.css'

interface RecentProjectsProps {
  onOpenProject?: (projectName: string, projectPath: string) => void
}

const RecentProjects: React.FC<RecentProjectsProps> = ({ onOpenProject }) => {
  const recentProjects = [
    { id: '1', name: '工程示例 1', lastOpened: '2026-01-20', path: '/path/to/project1.pant' },
    { id: '2', name: '工程示例 2', lastOpened: '2026-01-19', path: '/path/to/project2.pant' },
    { id: '3', name: '工程示例 3', lastOpened: '2026-01-18', path: '/path/to/project3.pant' },
  ]

  const handleProjectClick = (project: typeof recentProjects[0]) => {
    if (onOpenProject) {
      onOpenProject(project.name, project.path)
    }
  }

  return (
    <div className="recent-projects">
      <div className="recent-projects__grid">
        {recentProjects.map((project) => (
          <div
            key={project.id}
            className="project-card"
            onClick={() => handleProjectClick(project)}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault()
                handleProjectClick(project)
              }
            }}
          >
            <div className="project-card__icon">📁</div>
            <div className="project-card__info">
              <div className="project-card__name">{project.name}</div>
              <div className="project-card__date">{project.lastOpened}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default RecentProjects
