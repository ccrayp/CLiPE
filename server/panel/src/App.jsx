import AppRouter from './app/router'
import { AppProviders } from './app/providers'
import './App.css'

function App() {
  return (
    <AppProviders>
      <AppRouter />
    </AppProviders>
  )
}

export default App
