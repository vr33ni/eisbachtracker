<template>
<div id="app"
  class="min-h-screen bg-gradient-to-b from-blue-100 to-white dark:from-gray-900 dark:to-gray-800 px-0 sm:px-4 py-6 sm:py-10">

  <div class="w-full max-w-2xl sm:mx-auto space-y-6">

      <!-- Header -->
      <h1 class="text-4xl sm:text-5xl font-bold tracking-tight text-blue-600 dark:text-blue-300 text-center">
        Eisbach Tracker
      </h1>

      <!-- Water Data Card -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 sm:p-6">
        <WaterDataCard :water-level-text="waterLevelText" :water-flow-text="waterFlowText"
          :show-water-level-alert="showWaterLevelAlert" :water-temperature="waterTemperature"
          :water-temperature-loading="predictionLoading" :water-temperature-error="waterTemperatureError"
          :chart-labels="chartLabels" :chart-values="chartValues" :cached-age-minutes="cachedAgeMinutes" />


      </div>

      <!-- Surfer Spotter Card -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 sm:p-6 space-y-4">
        <h2 class="text-2xl font-semibold text-blue-700 dark:text-blue-300">🧍 Surfer Spotter</h2>

        <SurferPrediction :predictionLoading="predictionLoading" :predictionError="predictionError"
          :predictionLoadingMessage="predictionLoadingMessage" :currentHourPrediction="currentHourPrediction" />

        <SurferEntries :todaysEntries="todaysEntries" :historyEntries="historyEntries" :entriesLoading="entriesLoading"
          :entriesError="entriesError" :entriesLoadingMessage="entriesLoadingMessage" />

        <SurferEntriesForm :submitting="submitting" @submit="submitSurferCount" />
      </div>

      <!-- Refresh Button -->
      <div class="flex justify-center">
  <button @click="refreshEverything" :disabled="isRefreshing"
    class="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg text-lg font-medium shadow transition-all disabled:opacity-50 disabled:cursor-not-allowed">
    {{ isRefreshing ? 'Refreshing...' : 'Refresh Data' }}
  </button>
</div>


    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

import WaterDataCard from './components/WaterDataCard.vue'
import SurferEntries from './components/SurferEntries.vue'
import SurferPrediction from './components/SurferPrediction.vue'
import SurferEntriesForm from './components/SurferEntriesForm.vue'

import { useTemperature } from '@/composables/useWaterTemperatureData'
import { useSurferEntries } from '@/composables/useSurferEntries'
import { useWaterLevelData } from '@/composables/useWaterLevelData'

const surferCountRaw = ref('')
const currentHourPrediction = ref<number | null>(null)
const submitting = ref(false)

const { waterTemperature, waterTemperatureLoading, waterTemperatureError, cachedAgeMinutes, ensureTemperature} = useTemperature()

const {
  entriesLoading,
  errorEntries: entriesError,
  entriesLoadingMessage,
  fetchEntries,
  addEntry,
  getPredictionForHour,
  predictionLoading,
  predictionError,
  predictionLoadingMessage,
  todaysEntries,
  historyEntries,
} = useSurferEntries()

const {
  waterLevelText,
  waterFlowText,
  showWaterLevelAlert,
  chartLabels,
  chartValues,
  fetchWaterData,
  loading: waterLevelLoading,
} = useWaterLevelData()

const isRefreshing = computed(() => waterTemperatureLoading.value || entriesLoading.value || waterLevelLoading.value)

const fetchPrediction = async () => {
  try {
    await ensureTemperature() // Will use cache or fetch if needed

    const now = new Date()
    const hour = now.getHours()

    const prediction = await getPredictionForHour(hour, waterTemperature.value || undefined)

    if (prediction) {
      currentHourPrediction.value = prediction.prediction
      // Optionally update temp from prediction if backend returns something newer
      waterTemperature.value = prediction.water_temperature
    } else {
      currentHourPrediction.value = null
    }
  } catch (err) {
    console.error('Failed to fetch prediction', err)
    currentHourPrediction.value = null
  }
}


const surferCount = computed(() => Number(surferCountRaw.value))


const submitSurferCount = async () => {
  const count = surferCount.value
  if (!isNaN(count) && count >= 0) {
    submitting.value = true
    await addEntry(count, undefined, Number(waterLevelText.value.replace(' cm', '')), Number(waterFlowText.value.replace(' m³/s', '')))
    await fetchEntries()
    surferCountRaw.value = ''
    submitting.value = false
  }
}


const refreshEverything = async () => {
  await ensureTemperature()     
  await fetchWaterData()
  await fetchEntries()
  await fetchPrediction()
}

onMounted(refreshEverything)
</script>

<style scoped>
html,
body,
#app {
  height: 100%;
  margin: 0;
}
</style>
