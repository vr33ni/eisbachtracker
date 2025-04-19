<template>
  <div class="mt-4 text-left">
    <h3 class="font-semibold text-blue-700 dark:text-blue-300">📊 {{ t('predictionHeading') }}</h3>

    <!-- 1. Rotating message during loading -->
    <div v-if="predictionLoading" class="text-gray-500 italic">
      {{ rotatingMessage }}
    </div>

    <!-- 2. Show error -->
    <div v-else-if="predictionError" class="text-red-500">
      ❌ {{ predictionError }}
    </div>

    <!-- 3. Show prediction result -->
    <div v-else-if="currentHourPrediction !== null">
      {{ t('predictedSurfers') }} <strong>{{ currentHourPrediction }}</strong>
    </div>

    <!-- 4. Fallback if no prediction -->
    <div v-else-if="predictionHasBeenFetched" class="text-gray-500 italic">
      {{ t('notEnoughData') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLoadingMessages } from '@/composables/useLoadingMessages'

const { t, tm } = useI18n()

const props = defineProps<{
  predictionLoading: boolean
  predictionError: string | null
  currentHourPrediction: number | null
  predictionHasBeenFetched: boolean
}>()

const messages = computed(() => tm('loadingMessages') as string[])

const {
  loadingMessage: rotatingMessage,
  startRotating,
  stopRotating,
} = useLoadingMessages(messages, 3000)

// ⏱ Start / stop messages based on predictionLoading prop
watch(
  () => props.predictionLoading,
  (loading) => {
    if (loading) startRotating()
    else stopRotating()
  },
  { immediate: true }  
)


onUnmounted(stopRotating)
</script>
