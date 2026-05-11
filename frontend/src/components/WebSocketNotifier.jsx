import { useEffect } from 'react'
import { useWebSocket } from '../hooks/useWebSocket'
import { useToast } from '../contexts/ToastContext'

export default function WebSocketNotifier() {
  const { lastEvent } = useWebSocket()
  const { addToast } = useToast()

  useEffect(() => {
    if (!lastEvent) return
    switch (lastEvent.type) {
      case 'anomaly.detected':
        addToast(`Обнаружена аномалия на устройстве ${lastEvent.payload?.device_id?.slice(0,8)}...`, 'warning')
        break
      case 'scene.executed':
        addToast(`Сценарий «${lastEvent.payload?.scene_name}» выполнен`, 'success')
        break
      case 'telemetry.ingested':
        // Можно не показывать тост для каждого показания, либо показывать очень кратко
        break
      default:
        break
    }
  }, [lastEvent, addToast])

  return null // компонент ничего не рендерит
}