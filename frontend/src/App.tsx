import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from './components/Layout'
import { AuthProvider, useAuth } from './context/AuthContext'
import { PlayerModalProvider } from './context/PlayerModalContext'
import { HomePage } from './pages/HomePage'
import { BasePage } from './pages/BasePage'
import { GroupsPage } from './pages/GroupsPage'
import { TrackedPage } from './pages/TrackedPage'
import { SettingsPage } from './pages/SettingsPage'
import { LoginPage } from './pages/LoginPage'
import './App.css'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { token, loading } = useAuth()
  if (loading) return <div className="loading-full">Загрузка...</div>
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const { token, user, loading } = useAuth()
  if (loading) return <div className="loading-full">Загрузка...</div>
  if (!token) return <Navigate to="/login" replace />
  if (user?.role !== 'admin') return <div className="admin-denied">Доступ только для администратора</div>
  return <>{children}</>
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <PlayerModalProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<Layout />}>
              <Route index element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
              <Route path="base" element={<ProtectedRoute><BasePage /></ProtectedRoute>} />
              <Route path="groups" element={<ProtectedRoute><GroupsPage /></ProtectedRoute>} />
              <Route path="tracked" element={<ProtectedRoute><TrackedPage /></ProtectedRoute>} />
              <Route path="settings" element={<AdminRoute><SettingsPage /></AdminRoute>} />
            </Route>
          </Routes>
        </PlayerModalProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
