import axios from 'axios'
import {
  clearStoredSession,
  getAccessToken,
  getRefreshToken,
  updateStoredTokens,
} from '../auth/tokenStorage'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api/v1'

export const http = axios.create({
  baseURL: API_BASE_URL,
})

let onUnauthorized = () => {}
let refreshRequest = null

function unwrapTokens(response) {
  const data = response?.data?.data ?? {}

  return {
    accessToken: data.access_token,
    refreshToken: data.refresh_token,
  }
}

async function refreshTokens() {
  const refreshToken = getRefreshToken()

  if (!refreshToken) {
    throw new Error('Сессия истекла')
  }

  const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
    refresh_token: refreshToken,
  })

  const tokens = unwrapTokens(response)

  if (!tokens.accessToken || !tokens.refreshToken) {
    throw new Error('Сервер не вернул обновленные токены')
  }

  updateStoredTokens(tokens)
  return tokens
}

export function setUnauthorizedHandler(handler) {
  onUnauthorized = typeof handler === 'function' ? handler : () => {}
}

export function getApiErrorMessage(error, fallback = 'Не удалось выполнить запрос') {
  return (
    error?.response?.data?.message ||
    error?.response?.data?.error ||
    error?.message ||
    fallback
  )
}

http.interceptors.request.use((config) => {
  const token = getAccessToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    const isAuthEndpoint =
      originalRequest?.url?.includes('/auth/login') ||
      originalRequest?.url?.includes('/auth/refresh')

    if (
      error.response?.status === 401 &&
      !isAuthEndpoint &&
      !originalRequest?._retry &&
      getRefreshToken()
    ) {
      originalRequest._retry = true

      try {
        refreshRequest ??= refreshTokens()
        const tokens = await refreshRequest
        refreshRequest = null
        originalRequest.headers.Authorization = `Bearer ${tokens.accessToken}`
        return http(originalRequest)
      } catch (refreshError) {
        refreshRequest = null
        clearStoredSession()
        onUnauthorized()
        return Promise.reject(refreshError)
      }
    }

    if (error.response?.status === 401 && !isAuthEndpoint) {
      clearStoredSession()
      onUnauthorized()
    }

    return Promise.reject(error)
  },
)
