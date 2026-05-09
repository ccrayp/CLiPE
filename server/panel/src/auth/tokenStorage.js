const STORAGE_KEY = 'clipe-session'

export function getStoredSession() {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function saveStoredSession(session) {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(session))
}

export function clearStoredSession() {
  window.localStorage.removeItem(STORAGE_KEY)
}

export function getAccessToken() {
  return getStoredSession()?.accessToken ?? ''
}

export function getRefreshToken() {
  return getStoredSession()?.refreshToken ?? ''
}

export function updateStoredTokens(tokens) {
  const previous = getStoredSession() ?? {}

  saveStoredSession({
    ...previous,
    accessToken: tokens.accessToken,
    refreshToken: tokens.refreshToken,
  })
}
