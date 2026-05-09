import { Navigate, Outlet, Route, Routes } from 'react-router-dom'
import { AppLayout } from '../components/layout/AppLayout'
import { ProtectedRoute } from '../auth/ProtectedRoute'
import { DashboardPage } from '../pages/DashboardPage'
import { EntityPage } from '../pages/EntityPage'
import { LoginPage } from '../pages/LoginPage'
import { entityConfigs } from '../config/entities'

function ProtectedLayout() {
  return (
    <ProtectedRoute>
      <AppLayout>
        <Outlet />
      </AppLayout>
    </ProtectedRoute>
  )
}

function AppRouter() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<ProtectedLayout />}>
        <Route index element={<DashboardPage />} />
        {entityConfigs.map((config) => (
          <Route
            key={config.key}
            path={config.route}
            element={<EntityPage config={config} />}
          />
        ))}
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default AppRouter
