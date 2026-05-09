/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../api/authApi'
import { setUnauthorizedHandler } from '../api/http'
import {
  clearStoredSession,
  getStoredSession,
  saveStoredSession,
} from './tokenStorage'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const navigate = useNavigate()
  const [session, setSession] = useState(() => getStoredSession())
  const ready = true

  useEffect(() => {
    setUnauthorizedHandler(() => {
      setSession(null)
      navigate('/login', { replace: true })
    })
  }, [navigate])

  const login = async ({ username, password }) => {
    const data = await authApi.login({ username, password })
    const nextSession = {
      username,
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
    }

    saveStoredSession(nextSession)
    setSession(nextSession)
  }

  const logout = async () => {
    try {
      if (session?.refreshToken) {
        await authApi.logout(session.refreshToken)
      }
    } catch {
      // Even if backend logout fails, local session should still be cleared.
    } finally {
      clearStoredSession()
      setSession(null)
      navigate('/login', { replace: true })
    }
  }

  const value = {
    ready,
    session,
    isAuthenticated: Boolean(session?.accessToken),
    login,
    logout,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }

  return context
}
