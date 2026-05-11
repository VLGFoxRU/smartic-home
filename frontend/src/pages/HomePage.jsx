import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useToast } from '../contexts/ToastContext'
import { fetchRooms } from '../api/rooms'
import { fetchDevices, sendCommand, deleteDevice } from '../api/devices'
import { FiPlus, FiPower, FiSettings, FiAlertTriangle, FiArrowLeft, FiTrash2 } from 'react-icons/fi'
import { getLatestTelemetry } from '../api/telemetry'
import CreateRoomModal from '../components/CreateRoomModal'
import CreateDeviceModal from '../components/CreateDeviceModal'
import ConfirmDialog from '../components/ConfirmDialog'
import './HomePage.css'

export default function HomePage() {
  const { homeId } = useParams()
  const { token } = useAuth()
  const { addToast } = useToast()
  const [rooms, setRooms] = useState([])
  const [devices, setDevices] = useState([])
  const [loading, setLoading] = useState(true)
  const [showRoomModal, setShowRoomModal] = useState(false)
  const [showDeviceModal, setShowDeviceModal] = useState(false)
  const [selectedRoomId, setSelectedRoomId] = useState(null)
  const [telemetryData, setTelemetryData] = useState({})
  const [deleteConfirm, setDeleteConfirm] = useState({ open: false, deviceId: null, deviceName: '' })

  useEffect(() => {
    if (!token || !homeId) return
    Promise.all([
      fetchRooms(homeId),
      fetchDevices()
    ])
      .then(([roomsRes, devicesRes]) => {
        const fetchedDevices = devicesRes.data.data || []
        setRooms(roomsRes.data.data || [])
        setDevices(fetchedDevices)

        // Загружаем последнюю телеметрию для всех датчиков
        const sensorDevices = fetchedDevices.filter(d => d.type.includes('sensor'))
        return Promise.all(sensorDevices.map(d =>
            getLatestTelemetry(d.id)
                .then(res => ({ id: d.id, data: res.data.data }))
                .catch(() => ({ id: d.id, data: null }))
        ))
      })
      .then(telemetryResults => {
        const newTelemetry = {}
        telemetryResults.forEach(t => {
            if (t.data) newTelemetry[t.id] = t.data
        })
        setTelemetryData(newTelemetry)
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [homeId, token])

  const handleRoomCreated = (newRoom) => {
    setRooms(prev => [...prev, newRoom])
  }

  const handleDeviceCreated = (newDevice) => {
    setDevices(prev => [...prev, newDevice])
  }

  const togglePower = async (deviceId, currentState) => {
    const newOn = !(currentState?.power === 'on')
    try {
      await sendCommand(deviceId, 'set_power', { on: newOn })
      addToast(`Лампа ${newOn ? 'включена' : 'выключена'}`, 'success')
      // Оптимистично обновляем состояние в UI
      setDevices(prev => prev.map(d =>
        d.id === deviceId
          ? { ...d, state: { ...d.state, power: newOn ? 'on' : 'off' } }
          : d
      ))
    } catch (err) {
        addToast('Ошибка управления устройством', 'error')
        console.error('Ошибка управления устройством', err)
    }
  }

  // Группируем устройства по комнатам
  const roomsWithDevices = rooms.map(room => ({
    ...room,
    devices: devices.filter(d => d.room_id === room.id)
  }))

  if (loading) return <div className="home-loading">Загрузка...</div>

  return (
    <div className="home-page">
      <Link to="/dashboard" className="back-link">
        <FiArrowLeft />Назад
      </Link>

      <div className="home-actions">
        <button className="create-btn" onClick={() => setShowRoomModal(true)}>
            <FiPlus /> Добавить комнату
        </button>
      </div>

      <div className="rooms-grid">
        {roomsWithDevices.map(room => (
          <div key={room.id} className="room-card">
            <div className="room-header">
              <h2 className="room-name">{room.name}</h2>
              <button
                className="create-btn small"
                onClick={() => {
                  setSelectedRoomId(room.id)
                  setShowDeviceModal(true)
                }}
              >
                <FiPlus />
              </button>
            </div>
            <div className="devices-list">
              {room.devices.map(device => (
                <div key={device.id} className="device-item">
                    <button
                        className="delete-device-btn"
                        onClick={() => setDeleteConfirm({ open: true, deviceId: device.id, deviceName: device.name })}
                        title="Удалить устройство">
                        <FiTrash2 />
                    </button>
                  <div className="device-info">
                    <span className={`device-status ${device.status}`} />
                    <Link to={`/homes/${homeId}/devices/${device.id}`} className="device-name">
                      {device.name}
                    </Link>
                  </div>
                  <div className="device-meta">
                    
                    {(device.type === 'light' || device.type === 'power_switch') && (
                      <button
                        className={`power-btn ${device.state?.power === 'on' ? 'on' : 'off'}`}
                        onClick={() => togglePower(device.id, device.state)}
                      >
                        <FiPower />
                      </button>
                    )}
                    {device.type === 'temperature_sensor' && (
                        <span className="sensor-value">
                            {telemetryData[device.id]?.value?.toFixed(1) ?? '--'}°
                        </span>
                    )}
                    {device.type === 'humidity_sensor' && (
                        <span className="sensor-value">
                            {telemetryData[device.id]?.value?.toFixed(0) ?? '--'}%
                        </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {showRoomModal && (
        <CreateRoomModal
          homeId={homeId}
          onClose={() => setShowRoomModal(false)}
          onCreated={handleRoomCreated}
        />
      )}
      {showDeviceModal && selectedRoomId && (
        <CreateDeviceModal
          roomId={selectedRoomId}
          onClose={() => setShowDeviceModal(false)}
          onCreated={handleDeviceCreated}
        />
      )}
      {deleteConfirm.open && (
        <ConfirmDialog
            title="Удаление устройства"
            message={`Вы действительно хотите удалить «${deleteConfirm.deviceName}»? Это действие нельзя отменить.`}
            confirmLabel="Удалить"
            onConfirm={async () => {
            try {
                await deleteDevice(deleteConfirm.deviceId)
                addToast('Устройство удалено', 'success')
                setDevices(prev => prev.filter(d => d.id !== deleteConfirm.deviceId))
            } catch (err) {
                addToast('Ошибка удаления устройства', 'error')
            }
            setDeleteConfirm({ open: false })
            }}
            onCancel={() => setDeleteConfirm({ open: false })}
        />
      )}
    </div>
  )
}