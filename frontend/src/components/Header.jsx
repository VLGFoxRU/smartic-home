import { useState, useEffect, useRef } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useTheme } from '../contexts/ThemeContext'
import { fetchHomes } from '../api/homes'
import { FiHome, FiSettings, FiAlertTriangle, FiUser, FiLogOut, FiChevronDown, FiCheck, FiSun, FiMoon } from 'react-icons/fi'
import './Header.css'

export default function Header({ currentHomeId, onHomeChange }) {
  const { token, logout } = useAuth()
  const location = useLocation()
  const [homes, setHomes] = useState([])
  const [dropdownOpen, setDropdownOpen] = useState(false)
  const dropdownRef = useRef(null)
  const { dark, toggle } = useTheme()

  // закрытие дропдауна при клике вне
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        setDropdownOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    if (!token) return
    fetchHomes()
      .then(res => setHomes(res.data.data || []))
      .catch(console.error)
  }, [token])

  const currentHome = homes.find(h => h.id === currentHomeId)

  const isActive = (path) => location.pathname.startsWith(path)

  return (
    <header className="app-header">
      <Link to="/dashboard" className="header-logo">
        SmarticHome
      </Link>

      <nav className="header-nav">
        <Link
          to="/dashboard"
          className={`nav-link ${isActive('/dashboard') ? 'active' : ''}`}
        >
          <FiHome /> Дашборд
        </Link>
        {currentHomeId && (
          <>
            <Link
              to={`/homes/${currentHomeId}/scenes`}
              className={`nav-link ${isActive(`/homes/${currentHomeId}/scenes`) ? 'active' : ''}`}
            >
              <FiSettings /> Сценарии
            </Link>
            <Link
              to={`/homes/${currentHomeId}/anomalies`}
              className={`nav-link ${isActive(`/homes/${currentHomeId}/anomalies`) ? 'active' : ''}`}
            >
              <FiAlertTriangle /> Аномалии
            </Link>
          </>
        )}
      </nav>

      <div className="header-right">
        {homes.length > 0 && (
          <div className="home-dropdown" ref={dropdownRef}>
            <button
              className="home-dropdown-trigger"
              onClick={() => setDropdownOpen(!dropdownOpen)}
            >
              <FiHome className="trigger-icon" />
              <span>{currentHome ? currentHome.name : 'Выберите дом'}</span>
              <FiChevronDown className={`chevron ${dropdownOpen ? 'open' : ''}`} />
            </button>
            {dropdownOpen && (
              <div className="home-dropdown-menu">
                {homes.map(home => (
                  <button
                    key={home.id}
                    className={`home-dropdown-item ${home.id === currentHomeId ? 'active' : ''}`}
                    onClick={() => {
                      onHomeChange(home.id)
                      setDropdownOpen(false)
                    }}
                  >
                    <FiHome />
                    <span>{home.name}</span>
                    {home.id === currentHomeId && <FiCheck className="check-icon" />}
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
        <button className="theme-toggle" onClick={toggle} title={dark ? 'Светлая тема' : 'Тёмная тема'}>
            {dark ? <FiSun /> : <FiMoon />}
        </button>
        <Link to="/profile" className="nav-link">
          <FiUser />
        </Link>
        <button className="logout-btn" onClick={logout}>
          <FiLogOut />
        </button>
      </div>
    </header>
  )
}