<!-- GlobalLoadingOverlay.vue -->
<template>
    <div
      v-if="isRefreshing"
      class="fixed inset-0 bg-white/90 dark:bg-black/80 z-50 flex items-center justify-center cursor-pointer"
      @click="cancelRefresh"
    >
      <div class="text-center space-y-2 px-4">
        <div class="text-2xl animate-pulse text-blue-700 dark:text-blue-300">
          {{ rotatingMessage }}
        </div>
        <div class="text-xs mt-2 text-gray-400 italic">
          {{ $t('clickToCancel') }}
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  defineProps<{
    isRefreshing: boolean
    rotatingMessage: string
    cancelRefresh: () => void
  }>()
  
  import { useI18n } from 'vue-i18n'
  const { t } = useI18n()
  </script>
  