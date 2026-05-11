import { useState, useEffect } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useToast } from '../contexts/ToastContext'
import { fetchMe } from '../api/auth'
import apiClient from '../api/client'
import './ProfilePage.css'

export default function ProfilePage() {
  const { token } = useAuth()
  const { addToast } = useToast()
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  // Поля для смены пароля
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!token) return
    fetchMe()
      .then(res => setUser(res.data.data))
      .catch(() => addToast('Не удалось загрузить профиль', 'error'))
      .finally(() => setLoading(false))
  }, [token])

  const handleChangePassword = async (e) => {
    e.preventDefault()
    if (!oldPassword || !newPassword) {
      addToast('Заполните оба поля', 'error')
      return
    }
    setSaving(true)
    try {
      await apiClient.post('/auth/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      })
      addToast('Пароль успешно изменён', 'success')
      setOldPassword('')
      setNewPassword('')
    } catch (err) {
      addToast(err.response?.data?.error || 'Ошибка смены пароля', 'error')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="page-loading">Загрузка профиля...</div>
  }

  return (
    <div className="profile-page">
      <h1 className="page-title">Профиль</h1>

      {user && (
        <div className="profile-info">
          <div className="info-row">
            <span className="info-label">Имя пользователя</span>
            <span className="info-value">{user.username}</span>
          </div>
          <div className="info-row">
            <span className="info-label">Email</span>
            <span className="info-value">{user.email}</span>
          </div>
          <div className="info-row">
            <span className="info-label">Роль</span>
            <span className="info-value role-badge">{user.role}</span>
          </div>
        </div>
      )}

      <div className="password-section">
        <h2>Сменить пароль</h2>
        <form onSubmit={handleChangePassword} className="password-form">
          <input
            type="password"
            placeholder="Старый пароль"
            value={oldPassword}
            onChange={e => setOldPassword(e.target.value)}
            autoComplete="current-password"
            required
          />
          <input
            type="password"
            placeholder="Новый пароль"
            value={newPassword}
            onChange={e => setNewPassword(e.target.value)}
            autoComplete="new-password"
            minLength={6}
            required
          />
          <button type="submit" disabled={saving} className="save-btn">
            {saving ? 'Сохранение...' : 'Изменить пароль'}
          </button>
        </form>
      </div>
    </div>
  )
}