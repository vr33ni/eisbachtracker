import { createI18n } from 'vue-i18n'
import type { MessageSchema } from './types/i18n'
import en from './locales/en.json'
import es from './locales/es.json'
import de from './locales/de.json' // <-- add this

const i18n = createI18n<[MessageSchema], 'en' | 'es' | 'de'>({
  legacy: false,
  locale: localStorage.getItem('locale') || 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    es,
    de 
  }
})

export default i18n
