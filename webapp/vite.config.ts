import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
    '@': path.resolve(__dirname, './src')
    },
  },
  plugins: [vue()],
  base: 'http://localhost:8080',
  server: {
    cors: true,
    origin: 'http://127.0.0.1:8080',
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  }
})
