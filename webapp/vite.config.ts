import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
    '@': path.resolve(__dirname, './src')
    },
  },
  plugins: [
    vue(),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
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
