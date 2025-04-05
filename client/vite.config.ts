import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  base: process.env.NODE_ENV === 'production' ? '/eisbachtracker-pwa/' : '/',
    plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
    VitePWA({ 
      manifest: {
        name: 'EisbachTracker PWA',
        short_name: 'EisbachTracker',
        description: 'Your app description',
        start_url: '/eisbachtracker-pwa/',  
        icons: [
          {
            src: '/eisbachtracker-pwa/pwa-icon.svg',
            sizes: '192x192',
            type: 'image/svg+xml'
          },
        ]
      },
      registerType: 'autoUpdate',
      devOptions: {
        enabled: true
      }
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
