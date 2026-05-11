import { useState, useEffect, useCallback } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { fetchDevices, sendCommand } from '../api/devices'
import { getLatestTelemetry, getTelemetryHistory } from '../api/telemetry'
import { fetchAnomalies, updateAnomalyStatus } from '../api/anomalies'
import { Line } from 'react-chartjs-2'
import { FiPower, FiArrowLeft, FiRefreshCw } from 'react-icons/fi'
import { WiThermometer, WiHumidity, WiBarometer, WiLightning } from 'react-icons/wi'
import {
  Chart as ChartJS,
  LineElement,
  CategoryScale,
  LinearScale,
  PointElement,
  Tooltip,
  Legend,
} from 'chart.js'
import './DevicePage.css'

ChartJS.register(LineElement, CategoryScale, LinearScale, PointElement, Tooltip, Legend)

const periods = {
  '24ч': 24 * 60 * 60 * 1000,
  '7д': 7 * 24 * 60 * 60 * 1000,
  '30д': 30 * 24 * 60 * 60 * 1000,
}

const typeIcons = {
  temperature_sensor: <WiThermometer className="telemetry-icon temperature" />,
  humidity_sensor: <WiHumidity className="telemetry-icon humidity" />,
  electricity_meter: <WiLightning className="telemetry-icon power" />,
  water_meter: <WiBarometer className="telemetry-icon water" />,
}

export default function DevicePage() {
  const { homeId, deviceId } = useParams()
  const { token } = useAuth()
  const [device, setDevice] = useState(null)
  const [latest, setLatest] = useState(null)
  const [history, setHistory] = useState([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [period, setPeriod] = useState('24ч')
  const [anomalies, setAnomalies] = useState([])
  const [anomaliesLoading, setAnomaliesLoading] = useState(true)

  // Функция загрузки данных (вызывается при монтировании и по кнопке "Обновить")
  const fetchData = useCallback(() => {
    if (!token || !deviceId) return
    setRefreshing(true)

    // Информация об устройстве
    fetchDevices()
      .then((res) => {
        const dev = (res.data.data || []).find((d) => d.id === deviceId)
        setDevice(dev)
      })
      .catch(console.error)

    // Последние показания
    getLatestTelemetry(deviceId)
      .then((res) => setLatest(res.data.data))
      .catch(console.error)

    // История за выбранный период
    const to = new Date().toISOString()
    const from = new Date(Date.now() - periods[period]).toISOString()
    getTelemetryHistory(deviceId, from, to)
      .then((res) => setHistory(res.data.data || []))
      .catch(console.error)
      .finally(() => {
        setLoading(false)
        setRefreshing(false)
      })

    // Аномалии устройства
    setAnomaliesLoading(true)
    fetchAnomalies({ device_id: deviceId })
      .then((res) => setAnomalies(res.data.data || []))
      .catch(console.error)
      .finally(() => setAnomaliesLoading(false))
  }, [deviceId, token, period])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const togglePower = async () => {
    const currentOn = device?.state?.power === 'on'
    try {
      await sendCommand(deviceId, 'set_power', { on: !currentOn })
      setDevice((prev) => ({
        ...prev,
        state: { power: !currentOn ? 'on' : 'off' },
      }))
    } catch (err) {
      console.error('Ошибка управления устройством', err)
    }
  }

  const handleAnomalyAction = async (anomalyId, newStatus) => {
    try {
      await updateAnomalyStatus(anomalyId, newStatus)
      setAnomalies((prev) =>
        prev.map((a) =>
          a.id === anomalyId ? { ...a, status: newStatus } : a
        )
      )
    } catch (err) {
      console.error('Ошибка обновления аномалии', err)
    }
  }

  if (loading) return <div className="page-loading">Загрузка...</div>
  if (!device) return <div className="page-error">Устройство не найдено</div>

  const chartData = {
    labels: history.map((p) =>
      new Date(p.timestamp).toLocaleString([], {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      })
    ),
    datasets: [
      {
        label: device.name,
        data: history.map((p) => p.value),
        borderColor: '#C66518',
        backgroundColor: 'rgba(198, 101, 24, 0.1)',
        fill: true,
        tension: 0.3,
      },
    ],
  }

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      tooltip: {
        callbacks: {
          label: (context) => `${context.dataset.label}: ${context.raw}`,
          title: (context) => context[0].label,
        },
      },
    },
  }

  return (
    <div className="device-page">
      <div className="top-nav">
        <Link to={`/homes/${homeId}`} className="back-link">
          <FiArrowLeft /> Назад
        </Link>
      </div>

      <div className="device-header">
        <h1>{device.name}</h1>
        <span className={`status-badge ${device.status}`}>{device.status}</span>
      </div>

      {latest && (
        <div className="latest-telemetry">
          <div className="latest-icon">
            {typeIcons[device.type] || <WiBarometer className="telemetry-icon" />}
          </div>
          <span className="latest-value">{latest.value}</span>
          <span className="latest-type">{latest.type}</span>
        </div>
      )}

      {(device.type === 'light' || device.type === 'power_switch') && (
        <div className="device-control">
          <button
            className={`power-btn ${device.state?.power === 'on' ? 'on' : 'off'}`}
            onClick={togglePower}
          >
            <FiPower />
            {device.state?.power === 'on' ? ' Выключить' : ' Включить'}
          </button>
        </div>
      )}

      <div className="period-selector">
        {Object.keys(periods).map((key) => (
          <button
            key={key}
            className={`period-btn ${period === key ? 'active' : ''}`}
            onClick={() => setPeriod(key)}
          >
            {key}
          </button>
        ))}
        <button
          className="refresh-btn"
          onClick={fetchData}
          disabled={refreshing}
          title="Обновить данные"
        >
          <FiRefreshCw className={refreshing ? 'spinning' : ''} />
        </button>
      </div>

      <div className="chart-container">
        {history.length > 0 ? (
          <Line data={chartData} options={chartOptions} />
        ) : (
          <p>Нет данных телеметрии за выбранный период</p>
        )}
      </div>

      <div className="device-anomalies">
        <h3>Аномалии</h3>
        {anomaliesLoading ? (
          <p>Загрузка...</p>
        ) : anomalies.length === 0 ? (
          <p className="no-anomalies">Аномалий для этого устройства нет</p>
        ) : (
          <div className="anomalies-mini-list">
            {anomalies.slice(0, 5).map((a) => (
              <div key={a.id} className={`anomaly-mini-card ${a.severity}`}>
                <div className="anomaly-mini-info">
                  <span className="anomaly-mini-value">{a.value}</span>
                  {a.expected_value && (
                    <span className="anomaly-mini-expected">
                      Ожид: {a.expected_value}
                    </span>
                  )}
                  <span className="anomaly-mini-date">
                    {new Date(a.detected_at).toLocaleString()}
                  </span>
                </div>
                {a.status !== 'resolved' && a.status !== 'false_positive' && (
                  <div className="anomaly-mini-actions">
                    <button onClick={() => handleAnomalyAction(a.id, 'acknowledged')}>
                      ✓
                    </button>
                    <button onClick={() => handleAnomalyAction(a.id, 'resolved')}>
                      √
                    </button>
                    <button onClick={() => handleAnomalyAction(a.id, 'false_positive')}>
                      ×
                    </button>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}