import { useState } from 'react'
import { createRoom } from '../api/rooms'
import './CreateHomeModal.css' // используем те же стили

export default function CreateRoomModal({ homeId, onClose, onCreated }) {
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!name.trim()) return
    setLoading(true)
    setError('')
    try {
      const res = await createRoom(homeId, name)
      if (res.data.success) {
        onCreated(res.data.data)
        onClose()
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания комнаты')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Новая комната</h2>
        <form onSubmit={handleSubmit}>
          <input
            className="modal-input"
            type="text"
            placeholder="Название комнаты"
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
            required
          />
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