/**
 * WebSocket Service
 * Manages WebSocket connection to device backend
 * Subscribes to real-time data updates
 */

export class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts: number = 0
  private maxReconnectAttempts: number = 5
  private reconnectDelay: number = 3000

  // Message handlers
  private messageHandlers: Map<string, (data: any) => void> = new Map()

  constructor(url: string) {
    this.url = url
  }

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected')
      return
    }

    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0

        // Subscribe to variable changes
        this.subscribe()
      }

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          this.handleMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.attemptReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)

      setTimeout(() => {
        this.connect()
      }, this.reconnectDelay)
    } else {
      console.error('Max reconnect attempts reached')
    }
  }

  private subscribe() {
    // Send subscription message
    this.send({
      type: 'subscribe',
      channels: ['variable_changes', 'device_status']
    })
  }

  private send(data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  private handleMessage(message: any) {
    const { type, data } = message

    // Call registered handler
    const handler = this.messageHandlers.get(type)
    if (handler) {
      handler(data)
    }
  }

  // Register message handler
  on(type: string, handler: (data: any) => void) {
    this.messageHandlers.set(type, handler)
  }

  // Unregister message handler
  off(type: string) {
    this.messageHandlers.delete(type)
  }
}

// Create singleton instance
export const wsService = new WebSocketService('ws://localhost:8080/ws')
