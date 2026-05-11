import { useState, useCallback } from 'react'
import { useNavigate, Outlet } from 'react-router-dom'
import Header from './Header'
import WebSocketNotifier from './WebSocketNotifier'

export default function Layout() {
  const navigate = useNavigate()
  const [currentHomeId, setCurrentHomeId] = useState(
    localStorage.getItem('currentHomeId') || ''
  )

  const handleHomeChange = useCallback((homeId) => {
    setCurrentHomeId(homeId)
    localStorage.setItem('currentHomeId', homeId)
    if (homeId) {
      navigate(`/homes/${homeId}`)
    } else {
      navigate('/dashboard')
    }
  }, [navigate])

  return (
    <div className="app-layout">
      <Header currentHomeId={currentHomeId} onHomeChange={handleHomeChange} />
      <WebSocketNotifier />
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}