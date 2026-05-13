import { useState } from 'react'
import { Button, Card, Container, Form } from 'react-bootstrap'
import { Navigate, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { getApiErrorMessage } from '../api/http'
import Notification from '../components/common/Notification'

export function LoginPage() {
  const navigate = useNavigate()
  const { isAuthenticated, login } = useAuth()
  const [credentials, setCredentials] = useState({
    username: '',
    password: '',
  })
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setError('')
    setBusy(true)

    try {
      await login(credentials)
      navigate('/', { replace: true })
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, 'Не удалось выполнить вход'))
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="login-page d-flex align-items-center">
      <Container className="py-5">
        <Card className="login-card mx-auto border-0">
          <Card.Body className="p-4 p-md-5">
            <h1 className="h2 fw-bold text-dark mb-3">Вход</h1>

            {error ? <Notification variant="danger" dismissible>{error}</Notification> : null}

            <Form onSubmit={handleSubmit}>
              <Form.Group className="mb-3" controlId="username">
                <Form.Label>Логин</Form.Label>
                <Form.Control
                  value={credentials.username}
                  onChange={(event) =>
                    setCredentials((previous) => ({
                      ...previous,
                      username: event.target.value,
                    }))
                  }
                  required
                />
              </Form.Group>

              <Form.Group className="mb-4" controlId="password">
                <Form.Label>Пароль</Form.Label>
                <Form.Control
                  type="password"
                  value={credentials.password}
                  onChange={(event) =>
                    setCredentials((previous) => ({
                      ...previous,
                      password: event.target.value,
                    }))
                  }
                  required
                />
              </Form.Group>

              <Button type="submit" className="w-100 py-2" disabled={busy}>
                {busy ? 'Выполняется вход...' : 'Войти'}
              </Button>
            </Form>
          </Card.Body>
        </Card>
      </Container>
    </div>
  )
}
