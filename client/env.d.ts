/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// Add your environment variable typings here
interface ImportMetaEnv {
  readonly VITE_ENV: string
  readonly VITE_BACKEND_API_URL: string
  readonly MODE: string
  readonly BASE_URL: string
  // Add more as needed
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
