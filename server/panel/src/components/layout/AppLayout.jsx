import { useState } from 'react'
import { Container, Button, Offcanvas, Nav, Navbar } from 'react-bootstrap'
import { NavLink, Outlet } from 'react-router-dom'
import { entityConfigs } from '../../config/entities'
import { useAuth } from '../../auth/AuthContext'

const groups = [
  {
    title: 'Справочники',
    items: ['users', 'hosts', 'services'],
  },
  {
    title: 'Управление доступом',
    items: ['rules', 'policies', 'policyContents'],
  },
  {
    title: 'Аудит',
    items: ['requests', 'decisions'],
  },
]

export function AppLayout({ children }) {
  const { session, logout } = useAuth()
  const [showMenu, setShowMenu] = useState(false)

  const configByKey = entityConfigs.reduce((accumulator, item) => {
    accumulator[item.key] = item
    return accumulator
  }, {})

  return (
    <div className="app-shell">
      <Offcanvas show={showMenu} onHide={() => setShowMenu(false)} placement="start">
        <Offcanvas.Header closeButton>
          <Offcanvas.Title>Навигация</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
          <div className="mb-4">
            <div className="fw-semibold">CLiPE</div>
            <div className="small text-muted">
              Централизованное управление доступом
            </div>
          </div>

          <Nav className="flex-column gap-2 mb-4">
            <NavLink to="/" end className="nav-link app-drawer-link" onClick={() => setShowMenu(false)}>
              Домашняя страница
            </NavLink>
          </Nav>

          {groups.map((group) => (
            <div key={group.title} className="mb-4">
              <div className="small text-uppercase text-muted mb-2">{group.title}</div>
              <Nav className="flex-column gap-2">
                {group.items.map((key) => (
                  <NavLink
                    key={key}
                    to={`/${configByKey[key].route}`}
                    className="nav-link app-drawer-link"
                    onClick={() => setShowMenu(false)}
                  >
                    {configByKey[key].title}
                  </NavLink>
                ))}
              </Nav>
            </div>
          ))}
        </Offcanvas.Body>
      </Offcanvas>

      <Navbar className="app-topbar" sticky="top">
        <Container fluid className="px-3 px-md-4">
          <div className="d-flex align-items-center gap-3">
            <Button
              className="floating-menu-button"
              variant="outline-secondary"
              onClick={() => setShowMenu(true)}
            >
              Меню
            </Button>
            <Navbar.Brand className="fw-semibold text-dark mb-0">
              Панель управления CLiPE
            </Navbar.Brand>
          </div>
          <div className="d-flex align-items-center gap-3 ms-auto">
            <div className="text-end d-none d-md-block">
              <div className="small text-muted">Пользователь</div>
              <div className="fw-semibold">{session?.username ?? 'admin'}</div>
            </div>
            <Button variant="outline-secondary" onClick={logout}>
              Выйти
            </Button>
          </div>
        </Container>
      </Navbar>

      <main className="app-content">
        <Container className="content-container px-0">
          {children ?? <Outlet />}
        </Container>
      </main>
    </div>
  )
}
