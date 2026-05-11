import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useToast } from '../contexts/ToastContext'
import { fetchAnomalies, updateAnomalyStatus } from '../api/anomalies'
import { FiAlertTriangle, FiCheck, FiX, FiFlag } from 'react-icons/fi'
import './AnomaliesPage.css'

export default function AnomaliesPage() {
  const { homeId } = useParams()
  const { token } = useAuth()
  const { addToast } = useToast()
  const [anomalies, setAnomalies] = useState([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState('open')

  useEffect(() => {
    if (!token || !homeId) return
    fetchAnomalies({ status: statusFilter })
      .then(res => setAnomalies(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [homeId, token, statusFilter])

  const handleStatusChange = async (id, newStatus) => {
    try {
      await updateAnomalyStatus(id, newStatus)
      addToast('Статус аномалии обновлён', 'success')
      setAnomalies(prev => prev.map(a => 
        a.id === id ? { ...a, status: newStatus } : a
      ))
    } catch (err) {
        addToast('Ошибка изменения статуса аномалии', 'error')
      console.error('Ошибка изменения статуса аномалии', err)
    }
  }

  if (loading) return <div className="page-loading">Загрузка...</div>

  return (
    <div className="anomalies-page">      
      <div className="page-header">
        <h1>Аномалии</h1>
        <select 
          className="status-filter"
          value={statusFilter} 
          onChange={e => setStatusFilter(e.target.value)}
        >
          <option value="open">Открытые</option>
          <option value="acknowledged">Подтверждённые</option>
          <option value="resolved">Устранённые</option>
          <option value="false_positive">Ложные</option>
          <option value="">Все</option>
        </select>
      </div>

      {anomalies.length === 0 ? (
        <div className="empty-state">
          <FiAlertTriangle style={{ fontSize: '2rem', color: 'var(--color-text-muted)' }} />
          <p>Аномалий не обнаружено</p>
        </div>
      ) : (
        <div className="anomalies-list">
          {anomalies.map(anomaly => (
            <div key={anomaly.id} className={`anomaly-card ${anomaly.severity}`}>
              <div className="anomaly-main">
                <div className="anomaly-header">
                  <h3>Устройство: {anomaly.device_id?.slice(0, 8)}...</h3>
                  <span className={`severity-badge ${anomaly.severity}`}>
                    {anomaly.severity === 'critical' ? 'Критическая' : 'Предупреждение'}
                  </span>
                </div>
                <div className="anomaly-details">
                  <span>Значение: {anomaly.value}</span>
                  {anomaly.expected_value && <span>Ожидалось: {anomaly.expected_value}</span>}
                  <span>Обнаружено: {new Date(anomaly.detected_at).toLocaleString()}</span>
                </div>
              </div>
              
              {anomaly.status === 'open' || anomaly.status === 'acknowledged' ? (
                <div className="anomaly-actions">
                  {anomaly.status === 'open' && (
                    <button 
                      className="action-btn acknowledge"
                      onClick={() => handleStatusChange(anomaly.id, 'acknowledged')}
                      title="Подтвердить"
                    >
                      <FiFlag />
                    </button>
                  )}
                  {(anomaly.status === 'open' || anomaly.status === 'acknowledged') && (
                    <button 
                      className="action-btn resolve"
                      onClick={() => handleStatusChange(anomaly.id, 'resolved')}
                      title="Устранить"
                    >
                      <FiCheck />
                    </button>
                  )}
                  <button 
                    className="action-btn false-positive"
                    onClick={() => handleStatusChange(anomaly.id, 'false_positive')}
                    title="Ложное срабатывание"
                  >
                    <FiX />
                  </button>
                </div>
              ) : (
                <span className={`status-label ${anomaly.status}`}>
                  {anomaly.status === 'resolved' ? 'Устранена' : 'Ложная'}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}