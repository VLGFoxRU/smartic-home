import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { fetchHomes } from '../api/homes'
import { FiHome, FiLogOut, FiPlus } from 'react-icons/fi'
import CreateHomeModal from '../components/CreateHomeModal'
import CopyButton from '../components/CopyButton'
import './Dashboard.css'

export default function DashboardPage() {
  const { token, logout } = useAuth()
  const [homes, setHomes] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) return
    fetchHomes()
      .then((res) => {
        if (res.data.success) {
          setHomes(res.data.data)
        }
      })
      .catch((err) => {
        console.error('Ошибка загрузки домов:', err)
        setError('Не удалось загрузить дома')
      })
      .finally(() => setLoading(false))
  }, [token])

  const handleHomeCreated = (newHome) => {
    setHomes((prev) => [...prev, newHome])
  }

  if (!token) return <Navigate to="/login" replace />

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <button
            className="create-home-btn"
            onClick={() => setShowCreateModal(true)}
          >
            <FiPlus />
            Создать дом
          </button>
        </div>
      </header>

      {loading && (
        <div className="dashboard-loading">
          <div className="spinner" />
        </div>
      )}

      {error && <div className="dashboard-error">{error}</div>}

      {!loading && !error && homes.length === 0 && (
        <div className="dashboard-empty">
          <FiHome style={{ fontSize: '3rem', color: 'var(--color-text-muted)' }} />
          <p>У вас пока нет домов. Создайте первый дом через приложение.</p>
          <button
            className="create-home-btn"
            onClick={() => setShowCreateModal(true)}
          >
            <FiPlus />
            Создать дом
          </button>
        </div>
      )}

      {!loading && homes.length > 0 && (
        <section className="homes-grid">
          {homes.map((home) => (
            <div key={home.id} className="home-card" onClick={() => navigate(`/homes/${home.id}`)}>
              <h2 className="home-name">{home.name}</h2>
              <p className="home-id">
                <CopyButton text={home.id} small />
                ID: {home.id.slice(0, 8)}...
              </p>
            </div>
          ))}
        </section>
      )}

      {showCreateModal && (
        <CreateHomeModal
          onClose={() => setShowCreateModal(false)}
          onCreated={handleHomeCreated}
        />
      )}
    </div>
  )
}