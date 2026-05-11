import { useState } from 'react'
import { createHome } from '../api/homes'
import { useToast } from '../contexts/ToastContext'
import './CreateHomeModal.css'

export default function CreateHomeModal({ onClose, onCreated }) {
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const { addToast } = useToast()

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!name.trim()) return
    setLoading(true)
    setError('')
    try {
      const res = await createHome(name)
      if (res.data.success) {
        addToast('Дом успешно создан', 'success')
        onCreated(res.data.data) // передаём новый дом наверх
        onClose()
      }
    } catch (err) {
      addToast('Ошибка создания дома', 'error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2 className="modal-title">Новый дом</h2>
        <form onSubmit={handleSubmit}>
          <input
            className="modal-input"
            type="text"
            placeholder="Название дома"
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
            required
          />
          {error && <div className="modal-error">{error}</div>}
          <div className="modal-buttons">
            <button type="button" onClick={onClose} disabled={loading}>
              Отмена
            </button>
            <button type="submit" disabled={loading}>
              {loading ? 'Создаётся...' : 'Создать'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}