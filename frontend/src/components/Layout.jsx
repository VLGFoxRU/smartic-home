import { useState, useEffect, useCallback } from 'react'
import { useNavigate, Outlet } from 'react-router-dom'
import Header from './Header'
import { fetchHome } from '../api/homes'

export default function Layout() {
  const navigate = useNavigate()
  const [currentHomeId, setCurrentHomeId] = useState(
    localStorage.getItem('currentHomeId') || ''
  )
  const [homeValid, setHomeValid] = useState(false)
  const [validating, setValidating] = useState(!!currentHomeId)

  useEffect(() => {
    if (!currentHomeId) {
      setValidating(false)
      return
    }
    // Проверяем, существует ли дом с таким ID
    fetchHome(currentHomeId)
      .then(() => {
        setHomeValid(true)
        navigate(`/homes/${currentHomeId}`)
      })
      .catch(() => {
        // Дом не найден – сбрасываем
        localStorage.removeItem('currentHomeId')
        setCurrentHomeId('')
        navigate('/dashboard')
      })
      .finally(() => setValidating(false))
  }, []) // запускаем один раз при монтировании

  const handleHomeChange = useCallback((homeId) => {
    setCurrentHomeId(homeId)
    localStorage.setItem('currentHomeId', homeId)
    if (homeId) {
      navigate(`/homes/${homeId}`)
    } else {
      navigate('/dashboard')
    }
  }, [navigate])

  if (validating) {
    return <div className="page-loading">Проверка данных...</div>
  }

  return (
    <div className="app-layout">
      <Header currentHomeId={currentHomeId} onHomeChange={handleHomeChange} />
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}