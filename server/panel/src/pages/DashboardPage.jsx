import { useEffect, useState } from 'react'
import { Alert, Badge, Button, Card, Col, Row } from 'react-bootstrap'
import { Link } from 'react-router-dom'
import { apiMap } from '../api'
import { getApiErrorMessage } from '../api/http'
import { PageHeader } from '../components/common/PageHeader'
import { entityConfigs } from '../config/entities'
import { formatDateTime } from '../utils/formatters'

export function DashboardPage() {
  const [metrics, setMetrics] = useState({})
  const [recentRequests, setRecentRequests] = useState([])
  const [recentDecisions, setRecentDecisions] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    async function loadDashboard() {
      setLoading(true)
      setError('')

      try {
        const metricEntries = await Promise.all(
          entityConfigs.map(async (config) => {
            const result = await apiMap[config.key].search({}, 1, 0)
            return [config.key, result.count]
          }),
        )

        const [requestsResult, decisionsResult] = await Promise.all([
          apiMap.requests.search({}, 5, 0),
          apiMap.decisions.search({}, 5, 0),
        ])

        if (!mounted) {
          return
        }

        setMetrics(Object.fromEntries(metricEntries))
        setRecentRequests(requestsResult.items)
        setRecentDecisions(decisionsResult.items)
      } catch (requestError) {
        if (mounted) {
          setError(getApiErrorMessage(requestError, 'Не удалось загрузить dashboard'))
        }
      } finally {
        if (mounted) {
          setLoading(false)
        }
      }
    }

    loadDashboard()

    return () => {
      mounted = false
    }
  }, [])

  return (
    <>
      <PageHeader
        eyebrow="Домашняя страница"
        title="Панель централизованного управления доступом"
        description="Оперативный обзор основных сущностей, журналов и состояния административного контура CLiPE."
        actions={
          <div className="d-flex justify-content-lg-end">
            <Button as={Link} to="/policies" variant="dark">
              Перейти к политикам
            </Button>
          </div>
        }
      />

      {error ? (
        <Alert variant="danger" className="mb-4" dismissible>
          {error}
        </Alert>
      ) : null}

      <Row className="g-4 mb-4">
        {entityConfigs.map((config) => (
          <Col key={config.key} md={6} xl={4}>
            <Card className="surface-card metric-card border-0 h-100">
              <Card.Body className="d-flex flex-column justify-content-between">
                <div>
                  <div className="section-label mb-2">{config.eyebrow}</div>
                  <div className="h5 fw-semibold text-dark">{config.title}</div>
                  <p className="text-secondary small mb-0">{config.description}</p>
                </div>
                <div className="d-flex justify-content-between align-items-end mt-4">
                  <div className="metric-value">
                    {loading ? '...' : metrics[config.key] ?? 0}
                  </div>
                  <Button as={Link} to={`/${config.route}`} variant="outline-dark" size="sm">
                    Открыть
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </Col>
        ))}
      </Row>

      <Row className="g-4">
        <Col xl={6}>
          <Card className="surface-card h-100">
            <Card.Header className="py-3">
              <div className="fw-semibold">Последние запросы</div>
            </Card.Header>
            <Card.Body>
              {recentRequests.length === 0 ? (
                <div className="text-secondary">Данные еще не загружены или отсутствуют.</div>
              ) : (
                <div className="d-flex flex-column gap-3">
                  {recentRequests.map((request) => (
                    <div
                      key={request.request_id}
                      className="d-flex justify-content-between align-items-start gap-3"
                    >
                      <div>
                        <div className="fw-semibold text-dark">Request #{request.request_id}</div>
                        <div className="small text-secondary">
                          User ID: {request.user_id ?? '—'}
                        </div>
                      </div>
                      <Badge bg="light" text="dark">
                        {formatDateTime(request.timestamp)}
                      </Badge>
                    </div>
                  ))}
                </div>
              )}
            </Card.Body>
          </Card>
        </Col>

        <Col xl={6}>
          <Card className="surface-card h-100">
            <Card.Header className="py-3">
              <div className="fw-semibold">Последние решения</div>
            </Card.Header>
            <Card.Body>
              {recentDecisions.length === 0 ? (
                <div className="text-secondary">Данные еще не загружены или отсутствуют.</div>
              ) : (
                <div className="d-flex flex-column gap-3">
                  {recentDecisions.map((decision) => (
                    <div
                      key={decision.decision_id}
                      className="d-flex justify-content-between align-items-start gap-3"
                    >
                      <div>
                        <div className="fw-semibold text-dark">
                          Decision #{decision.decision_id}
                        </div>
                        <div className="small text-secondary">
                          Request #{decision.request_id}, Policy #{decision.policy_id ?? '—'}
                        </div>
                      </div>
                      <Badge bg={decision.result ? 'success' : 'danger'}>
                        {decision.result ? 'Allow' : 'Deny'}
                      </Badge>
                    </div>
                  ))}
                </div>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </>
  )
}
