import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { fetchScenes, activateScene, pauseScene, deleteScene } from '../api/scenes'
import { FiPlay, FiPause, FiTrash2, FiPlus } from 'react-icons/fi'
import CreateSceneModal from '../components/CreateSceneModal'
import ConfirmDialog from '../components/ConfirmDialog'
import { useToast } from '../contexts/ToastContext'
import './ScenesPage.css'

export default function ScenesPage() {
  const { homeId } = useParams()
  const { token } = useAuth()
  const { addToast } = useToast()
  const [scenes, setScenes] = useState([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [deleteSceneId, setDeleteSceneId] = useState(null)

  useEffect(() => {
    if (!token || !homeId) return
    fetchScenes(homeId)
      .then(res => setScenes(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [homeId, token])

  const handleActivate = async (id) => {
    try {
      await activateScene(id)
      addToast('Сценарий активирован', 'success')
      setScenes(prev => prev.map(s => s.id === id ? { ...s, is_active: true } : s))
    } catch (err) {
      addToast('Не удалось активировать сценарий', 'error')
    }
  }

  const handlePause = async (id) => {
    try {
      await pauseScene(id)
      addToast('Сценарий приостановлен', 'success')
      setScenes(prev => prev.map(s => s.id === id ? { ...s, is_active: false } : s))
    } catch (err) {
      addToast('Не удалось приостановить сценарий', 'error')
    }
  }

  const handleDeleteConfirm = async () => {
    if (!deleteSceneId) return
    try {
      await deleteScene(deleteSceneId)
      addToast('Сценарий удалён', 'success')
      setScenes(prev => prev.filter(s => s.id !== deleteSceneId))
    } catch (err) {
      addToast('Ошибка удаления сценария', 'error')
    } finally {
      setDeleteSceneId(null)
    }
  }

  if (loading) return <div className="page-loading">Загрузка...</div>

  return (
    <div className="scenes-page">
      
      <div className="page-header">
        <h1>Сценарии</h1>
        <button className="create-btn" onClick={() => setShowCreate(true)}>
          <FiPlus /> Создать сценарий
        </button>
      </div>

      {scenes.length === 0 ? (
        <div className="empty-state">
          <FiPlay style={{ fontSize: '2rem', color: 'var(--color-text-muted)' }} />
          <p>Нет созданных сценариев</p>
        </div>
      ) : (
        <div className="scenes-list">
          {scenes.map(scene => (
            <div key={scene.id} className="scene-card">
              <div className="scene-info">
                <h3>{scene.name}</h3>
                <span className={`scene-status ${scene.is_active ? 'active' : ''}`}>
                  {scene.is_active ? 'Активен' : 'На паузе'}
                </span>
              </div>
              <div className="scene-actions">
                {scene.is_active ? (
                  <button onClick={() => handlePause(scene.id)} title="Приостановить">
                    <FiPause />
                  </button>
                ) : (
                  <button onClick={() => handleActivate(scene.id)} title="Активировать">
                    <FiPlay />
                  </button>
                )}
                <button onClick={() => setDeleteSceneId(scene.id)} title="Удалить">
                  <FiTrash2 />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showCreate && (
        <CreateSceneModal
          homeId={homeId}
          onClose={() => setShowCreate(false)}
          onCreated={(newScene) => setScenes(prev => [...prev, newScene])}
        />
      )}

      {deleteSceneId && (
        <ConfirmDialog
          title="Удаление сценария"
          message="Вы уверены, что хотите удалить этот сценарий? Действие необратимо."
          confirmLabel="Удалить"
          onConfirm={handleDeleteConfirm}
          onCancel={() => setDeleteSceneId(null)}
        />
      )}
    </div>
  )
}