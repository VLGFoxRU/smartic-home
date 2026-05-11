import { useEffect, useRef, useState, useCallback } from 'react'
import { useAuth } from '../contexts/AuthContext'

export function useWebSocket() {
  const { token } = useAuth()
  const [lastEvent, setLastEvent] = useState(null)
  const wsRef = useRef(null)
  const reconnectTimeoutRef = useRef(null)

  const connect = useCallback(() => {
    if (!token) return
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    const ws = new WebSocket(`ws://${window.location.host}/ws?token=${token}`)
    wsRef.current = ws

    ws.onopen = () => {
      console.log('WebSocket подключён')
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        setLastEvent(data)
      } catch (err) {
        console.error('Ошибка парсинга WebSocket сообщения', err)
      }
    }

    ws.onclose = (event) => {
      console.log('WebSocket закрыт', event.code, event.reason)
      // Переподключение через 5 секунд, если не было намеренного закрытия
      if (event.code !== 1000) {
        reconnectTimeoutRef.current = setTimeout(() => {
          connect()
        }, 5000)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket ошибка', error)
    }
  }, [token])

  useEffect(() => {
    connect()
    return () => {
      clearTimeout(reconnectTimeoutRef.current)
      wsRef.current?.close(1000, 'Компонент размонтирован')
    }
  }, [connect])

  return { lastEvent }
}