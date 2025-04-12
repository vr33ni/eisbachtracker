import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
// import router from './router'

const app = createApp(App)

app.use(createPinia())
// app.use(router)

console.log('🧭 Running in environment:', import.meta.env.MODE) // 'development' or 'production'

app.mount('#app')
