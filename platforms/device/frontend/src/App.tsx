import { Leafer, Rect } from 'leafer-ui'
import { Editor } from '@leafer-ui/editor'

function App() {
  React.useEffect(() => {
    // Initialize Leafer.js canvas
    const leafer = new Leafer({
      view: 'canvas',
      type: 'runtime'
    })

    // Add a test rectangle
    const rect = new Rect({
      x: 100,
      y: 100,
      width: 200,
      height: 200,
      fill: '#32cd79',
      stroke: '#000',
      strokeWidth: 2
    })

    leafer.add(rect)

    return () => {
      leafer.destroy()
    }
  }, [])

  return (
    <div style={{ width: '100vw', height: '100vh', overflow: 'hidden' }}>
      <h1>PansIot Device Platform - HMI Runtime</h1>
      <div id="leafer-canvas" style={{ width: '100%', height: 'calc(100% - 60px)' }}></div>
    </div>
  )
}

export default App
