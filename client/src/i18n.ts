import { createI18n } from 'vue-i18n'
import type { MessageSchema } from './types/i18n'
import en from './locales/en.json'
import es from './locales/es.json'

const i18n = createI18n<[MessageSchema], 'en' | 'es'>({
  legacy: false,
  locale: localStorage.getItem('locale') || 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    es
  }
})

export default i18n
