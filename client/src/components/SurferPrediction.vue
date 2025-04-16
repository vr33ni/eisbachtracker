<template>
  <div class="mt-4 text-left">
    <h3 class="font-semibold text-blue-700 dark:text-blue-300">📊 {{ t('predictionHeading') }}</h3>

    <div v-if="predictionLoading">{{ t('predictionLoadingMessage') }}</div>
    <div v-else-if="predictionError" class="text-red-500">❌ {{ predictionError }}</div>

    <div v-else-if="currentHourPrediction !== null">
      {{ t('predictedSurfers') }} <strong>{{ currentHourPrediction }}</strong>
    </div>

    <div v-else class="text-gray-500">
      {{ t('notEnoughData') }}
    </div>
  </div>
</template>


<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  predictionLoading: boolean
  predictionError: string | null
  predictionLoadingMessage: string
  currentHourPrediction: number | null
}>()
</script>