import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from '../auth/AuthContext'

export function AppProviders({ children }) {
  return (
    <BrowserRouter>
      <AuthProvider>{children}</AuthProvider>
    </BrowserRouter>
  )
}
