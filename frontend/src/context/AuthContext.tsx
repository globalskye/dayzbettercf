import { createContext, useCallback, useContext, useEffect, useState } from 'react'
import { login as apiLogin, fetchAuthMe } from '../api/client'

export type AuthUser = { id: number; username: string; role: string }

type AuthContextValue = {
  token: string | null
  user: AuthUser | null
  loading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

const TOKEN_KEY = 'dayzsmartcf_token'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY))
  const [user, setUser] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)

  const loadUser = useCallback(async (t: string) => {
    try {
      const u = await fetchAuthMe(t)
      setUser(u)
    } catch {
      setToken(null)
      setUser(null)
      localStorage.removeItem(TOKEN_KEY)
    }
  }, [])

  useEffect(() => {
    if (!token) {
      setUser(null)
      setLoading(false)
      return
    }
    loadUser(token).finally(() => setLoading(false))
  }, [token, loadUser])

  const login = useCallback(async (username: string, password: string) => {
    const { token: t, user: u } = await apiLogin(username, password)
    localStorage.setItem(TOKEN_KEY, t)
    setToken(t)
    setUser(u)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ token, user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
