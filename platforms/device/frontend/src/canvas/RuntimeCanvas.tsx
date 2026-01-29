import { useEffect, useRef } from 'react'
import { Leafer } from 'leafer-ui'
import { projectStore } from '../store/projectStore'

interface RuntimeCanvasProps {
  width?: number
  height?: number
}

/**
 * RuntimeCanvas Component
 * Renders the HMI runtime canvas using Leafer.js
 * Loads and displays project files
 */
export const RuntimeCanvas: React.FC<RuntimeCanvasProps> = ({
  width = 1920,
  height = 1080
}) => {
  const canvasRef = useRef<HTMLDivElement>(null)
  const leaferRef = useRef<Leafer | null>(null)

  useEffect(() => {
    if (!canvasRef.current) return

    // Initialize Leafer.js
    const leafer = new Leafer({
      view: canvasRef.current,
      type: 'runtime'
    })

    leaferRef.current = leafer

    // Load project if available
    if (projectStore.currentProjectFile) {
      // TODO: Implement project loading
      console.log('Would load project:', projectStore.currentProjectFile)
    }

    return () => {
      leafer.destroy()
    }
  }, [])

  useEffect(() => {
    // Re-render when components change
    if (!leaferRef.current) return

    // Clear existing
    leaferRef.current.clear()

    // Render components
    projectStore.components.forEach(component => {
      // TODO: Implement component rendering
      // This will create actual Leafer.js shapes based on component type
      console.log('Would render component:', component)
    })
  }, [projectStore.components])

  return (
    <div
      ref={canvasRef}
      style={{
        width: `${width}px`,
        height: `${height}px`,
        border: '1px solid #ccc',
        background: '#1a1a1a'
      }}
    />
  )
}
