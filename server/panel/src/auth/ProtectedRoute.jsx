import { Navigate } from 'react-router-dom'
import { Spinner } from 'react-bootstrap'
import { useAuth } from './AuthContext'

export function ProtectedRoute({ children }) {
  const { ready, isAuthenticated } = useAuth()

  if (!ready) {
    return (
      <div className="d-flex min-vh-100 align-items-center justify-content-center">
        <Spinner animation="border" variant="primary" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return children
}
