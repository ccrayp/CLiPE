import { http } from './http'

function unwrapData(response) {
  return response?.data?.data ?? {}
}

export const authApi = {
  async login(payload) {
    const response = await http.post('/auth/login', payload)
    return unwrapData(response)
  },

  async logout(refreshToken) {
    return http.post('/auth/logout', { refresh_token: refreshToken })
  },
}
