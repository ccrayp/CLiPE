import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    allowedHosts: [
      'clipe',
      'clipe.local',
      'localhost',
      '127.0.0.1',
      '.clipe'
    ],
    port: 5173,
    strictPort: false,
    proxy: {
      '/api/v1': {
        target: 'https://clipe',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
