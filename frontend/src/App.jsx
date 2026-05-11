import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import { ToastProvider } from './contexts/ToastContext'
import { ThemeProvider } from './contexts/ThemeContext'
import Layout from './components/Layout'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import DashboardPage from './pages/DashboardPage'
import HomePage from './pages/HomePage'
import DevicePage from './pages/DevicePage'
import ScenesPage from './pages/ScenesPage'
import AnomaliesPage from './pages/AnomaliesPage'
import ProfilePage from './pages/ProfilePage'

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <ToastProvider>
            <Routes>
              {/* Публичные страницы */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />

              {/* Защищённые страницы с Layout */}
              <Route element={<Layout />}>
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/homes/:homeId" element={<HomePage />} />
                <Route path="/homes/:homeId/scenes" element={<ScenesPage />} />
                <Route path="/homes/:homeId/anomalies" element={<AnomaliesPage />} />
                <Route path="/homes/:homeId/devices/:deviceId" element={<DevicePage />} />
                <Route path="/profile" element={<ProfilePage />} />
              </Route>

              <Route path="*" element={<Navigate to="/login" replace />} />
            </Routes>
          </ToastProvider>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App