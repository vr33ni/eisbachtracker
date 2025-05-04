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
      <button class="ml-4 px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600" @click="showModal = true">
        {{ t('viewExplanation') }}
      </button>
    </div>

    <!-- 4. Fallback if no prediction -->
    <div v-else-if="predictionHasBeenFetched" class="text-gray-500 italic">
      {{ t('notEnoughData') }}
    </div>

    <!-- Modal for explanation -->
    <Modal v-if="showModal" :visible="showModal" title="Feature Contributions" @close="showModal = false">
      <ul>
        <li v-for="(value, feature) in explanation" :key="feature">
          {{ feature }}: {{ value.toFixed(2) }}
        </li>
      </ul>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { watch, onUnmounted, computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLoadingMessages } from '@/composables/useLoadingMessages'
import Modal from './Modal.vue'

const { t, tm } = useI18n()

const props = defineProps<{
  predictionLoading: boolean
  predictionError: string | null
  currentHourPrediction: number | null
  explanation: Record<string, number>
  predictionHasBeenFetched: boolean
}>()

const messages = computed(() => tm('loadingMessages') as string[])
const showModal = ref(false)

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
