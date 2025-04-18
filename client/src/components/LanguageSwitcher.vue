<!-- LanguageSwitcher.vue -->
<template>
    <div class="flex justify-end gap-2">
      <button
        v-for="lang in langs"
        :key="lang"
        @click="setLang(lang)"
        :class="[
          'px-3 py-1 rounded border',
          locale === lang ? 'bg-blue-600 text-white' : 'bg-white text-blue-600'
        ]"
      >
        {{ lang.toUpperCase() }}
      </button>
    </div>
  </template>
  
  <script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  
  const { locale } = useI18n()
  const langs = [ 'en', 'de', 'es'] as const
  type Lang = typeof langs[number]
  
  const setLang = (lang: Lang) => {
    locale.value = lang
    localStorage.setItem('locale', lang)
  }
  </script>
  