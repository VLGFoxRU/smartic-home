import { useState } from 'react'
import { createDevice } from '../api/devices'
import { useToast } from '../contexts/ToastContext'
import './CreateHomeModal.css'

export default function CreateDeviceModal({ roomId, onClose, onCreated }) {
  const { addToast } = useToast()
  const [name, setName] = useState('')
  const [type, setType] = useState('light')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!name.trim()) return
    setLoading(true)
    setError('')
    try {
      const res = await createDevice(name, type, roomId)
      if (res.data.success) {
        addToast('Устройство успешно создано', 'success')
        onCreated(res.data.data)
        onClose()
      }
    } catch (err) {
        addToast('Ошибка создания устройства', 'error')
      setError(err.response?.data?.error || 'Ошибка создания устройства')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Новое устройство</h2>
        <form onSubmit={handleSubmit}>
          <input
            className="modal-input"
            type="text"
            placeholder="Название"
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
            required
          />
          <select
            className="modal-input"
            value={type}
            onChange={(e) => setType(e.target.value)}
          >
            <option value="light">Лампа</option>
            <option value="power_switch">Розетка</option>
            <option value="temperature_sensor">Датчик температуры</option>
            <option value="humidity_sensor">Датчик влажности</option>
          </select>
          {error && <div className="modal-error">{error}</div>}
          <div className="modal-buttons">
            <button type="button" onClick={onClose} disabled={loading}>Отмена</button>
            <button type="submit" disabled={loading}>
              {loading ? 'Создаётся...' : 'Создать'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}