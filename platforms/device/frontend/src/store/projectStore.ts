import { makeAutoObservable } from 'mobx'

/**
 * HMI Project Store
 * Manages project loading, canvas state, and data bindings
 */
export class ProjectStore {
  // Project data
  projectName: string = ''
  projectType: 'hmi' = 'hmi'
  canvasWidth: number = 1920
  canvasHeight: number = 1080
  components: any[] = []
  dataBindings: any[] = []

  // Runtime state
  isLoading: boolean = false
  isPlaying: boolean = true
  currentProjectFile: string = ''

  // Real-time data cache
  dataCache: Map<string, any> = new Map()

  constructor() {
    makeAutoObservable(this)
  }

  // Load project from file
  async loadProject(projectPath: string) {
    this.isLoading = true
    try {
      // TODO: Implement actual project file loading
      // This will read the .工程文件 and parse it
      console.log('Loading project from:', projectPath)

      // Simulated load
      this.currentProjectFile = projectPath
      this.projectName = 'Default HMI Project'
    } catch (error) {
      console.error('Failed to load project:', error)
    } finally {
      this.isLoading = false
    }
  }

  // Update component value from real-time data
  updateComponentValue(componentId: string, value: any) {
    const component = this.components.find(c => c.id === componentId)
    if (component) {
      // Update component properties
      if (component.properties) {
        Object.assign(component.properties, value)
      }
    }
  }

  // Update data cache (from WebSocket)
  updateDataCache(variable: string, value: any) {
    this.dataCache.set(variable, value)

    // Find and update bound components
    const bindings = this.dataBindings.filter(
      b => b.variable === variable
    )

    bindings.forEach(binding => {
      this.updateComponentValue(binding.componentId, {
        [binding.property]: value
      })
    })
  }

  // Play/Pause runtime
  togglePlay() {
    this.isPlaying = !this.isPlaying
  }

  // Reset to initial state
  reset() {
    this.projectName = ''
    this.components = []
    this.dataBindings = []
    this.dataCache.clear()
    this.isPlaying = true
  }
}

// Create singleton instance
export const projectStore = new ProjectStore()
