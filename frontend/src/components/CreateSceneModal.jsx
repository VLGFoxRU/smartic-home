import { useState, useEffect } from 'react'
import { fetchDevices } from '../api/devices'
import { createScene } from '../api/scenes'
import './CreateHomeModal.css' // используем общие стили

export default function CreateSceneModal({ homeId, onClose, onCreated }) {
  const [name, setName] = useState('')
  const [conditionType, setConditionType] = useState('temperature')
  const [conditionDevice, setConditionDevice] = useState('')
  const [operator, setOperator] = useState('>')
  const [conditionValue, setConditionValue] = useState(25)
  const [actionCommand, setActionCommand] = useState('set_power')
  const [actionDevice, setActionDevice] = useState('')
  const [actionParams, setActionParams] = useState('{"on": true}')
  const [devices, setDevices] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchDevices().then(res => setDevices(res.data.data || []))
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      let params = {}
      try { params = JSON.parse(actionParams) } catch { params = {} }
      const condition = {
        type: conditionType,
        device_id: conditionDevice,
        operator,
        value: parseFloat(conditionValue),
        cooldown_sec: 300
      }
      const action = {
        device_id: actionDevice,
        command: actionCommand,
        params
      }
      const res = await createScene({ name, home_id: homeId, condition, action })
      if (res.data.success) {
        onCreated(res.data.data)
        onClose()
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания сценария')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <h2 className="modal-title">Новый сценарий</h2>
        <form onSubmit={handleSubmit}>
          <input className="modal-input" placeholder="Название сценария" value={name} onChange={e => setName(e.target.value)} required />
          <select className="modal-input" value={conditionType} onChange={e => setConditionType(e.target.value)}>
            <option value="temperature">Температура</option>
            <option value="humidity">Влажность</option>
            <option value="power">Мощность</option>
          </select>
          <select className="modal-input" value={conditionDevice} onChange={e => setConditionDevice(e.target.value)} required>
            <option value="">Выберите датчик</option>
            {devices.filter(d => d.type.includes('sensor')).map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
          <select className="modal-input" value={operator} onChange={e => setOperator(e.target.value)}>
            <option value=">">{'>'}</option>
            <option value="<">{'<'}</option>
            <option value=">=">{'>='}</option>
            <option value="<=">{'<='}</option>
          </select>
          <input className="modal-input" type="number" value={conditionValue} onChange={e => setConditionValue(e.target.value)} required />
          <select className="modal-input" value={actionCommand} onChange={e => setActionCommand(e.target.value)}>
            <option value="set_power">Включить/Выключить</option>
          </select>
          <select className="modal-input" value={actionDevice} onChange={e => setActionDevice(e.target.value)} required>
            <option value="">Выберите устройство</option>
            {devices.filter(d => d.type === 'light' || d.type === 'power_switch').map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
          </select>
          <input className="modal-input" type="text" value={actionParams} onChange={e => setActionParams(e.target.value)} placeholder='{"on": true}' />
          {error && <div className="modal-error">{error}</div>}
          <div className="modal-buttons">
            <button type="button" onClick={onClose}>Отмена</button>
            <button type="submit" disabled={loading}>{loading ? 'Создаётся...' : 'Создать'}</button>
          </div>
        </form>
      </div>
    </div>
  )
}