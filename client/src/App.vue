<template>
  <div id="app"
    class="min-h-screen bg-gradient-to-b from-blue-100 to-white dark:from-gray-900 dark:to-gray-800 px-0 sm:px-4 py-6 sm:py-10">
    <!-- Global Refresh Overlay -->
    <GlobalLoadingOverlay :is-refreshing="isRefreshing" :rotating-message="rotatingMessage"
      :cancel-refresh="cancelRefresh" />


    <div class="w-full max-w-2xl sm:mx-auto space-y-6">

      <!-- Language Switcher -->
      <LanguageSwitcher />

      <!-- Header + last updated -->
      <div class="text-center space-y-1">
        <h1 class="text-4xl sm:text-5xl font-bold tracking-tight text-blue-600 dark:text-blue-300">
          {{ $t('title') }}
        </h1>
        <div v-if="lastRefreshTime && shouldShowLastUpdated" class="text-xs text-gray-500">
          🔄 {{ $t('lastUpdated') }}: {{ formatTimeAgo }}
        </div>


      </div>


      <!-- Water Data Card -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 sm:p-6">
        <WaterDataCard :current-water-level="currentWaterLevel" :current-water-flow="currentWaterFlow"
          :request-date="requestDate" :show-water-level-alert="showWaterLevelAlert"
          :water-temperature="waterTemperature" :water-temperature-loading="waterTemperatureLoading"
          :water-temperature-error="waterTemperatureError" :chart-labels="chartLabels" :chart-values="chartValues"
          :cached-age-minutes="cachedAgeMinutes" :water-data-loading="waterDataLoading" />
      </div>

      <!-- Surfer Spotter Card -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 sm:p-6 space-y-4">
        <h2 class="text-2xl font-semibold text-blue-700 dark:text-blue-300">🧍 {{ $t('surferSpotter') }}</h2>

        <SurferPrediction :prediction-loading="predictionLoading" :prediction-error="predictionError"
          :current-hour-prediction="currentHourPrediction" :prediction-has-been-fetched="predictionHasBeenFetched" />


        <SurferEntries :todaysEntries="todaysEntries" :historyEntries="historyEntries" :entriesLoading="entriesLoading"
          :entriesError="entriesError" :entriesLoadingMessage="entriesLoadingMessage" />

        <SurferEntriesForm :submitting="submitting" v-model="surferCountRaw" @submit="submitSurferCount" />
      </div>

      <!-- Refresh Button -->
      <div class="flex justify-center">
        <button @click="refreshEverything" :disabled="isRefreshing"
          class="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg text-lg font-medium shadow transition-all disabled:opacity-50 disabled:cursor-not-allowed">
          {{ isRefreshing ? t('refreshing') : t('refresh') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'

import WaterDataCard from './components/WaterDataCard.vue'
import SurferEntries from './components/SurferEntries.vue'
import SurferPrediction from './components/SurferPrediction.vue'
import SurferEntriesForm from './components/SurferEntriesForm.vue'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import GlobalLoadingOverlay from './components/GlobalLoadingOverlay.vue'

import { useTemperature } from '@/composables/useWaterTemperatureData'
import { useSurferEntries } from '@/composables/useSurferEntries'
import { useWaterLevelData } from '@/composables/useWaterLevelData'
import { useLoadingMessages } from '@/composables/useLoadingMessages'

const { t } = useI18n()


const surferCountRaw = ref('')
const submitting = ref(false)


const messages = [
  '🔮 Watching the river...',
  '📈 Crunching the numbers...',
  '🐟 Asking fish for advice...',
  '📷 Checking river cams...',
  '🌊 Feeling the vibes...',
]

const {
  loadingMessage: rotatingMessage,
  startRotating,
  stopRotating
} = useLoadingMessages(messages)

const isRefreshing = ref(false)

const cancelRefresh = () => {
  stopRotating()
  isRefreshing.value = false
}

const lastRefreshTime = ref<string | null>(null)
const now = ref(Date.now())

const formatTimeAgo = computed(() => {
  if (!lastRefreshTime.value) return ''
  const diffMs = now.value - new Date(lastRefreshTime.value).getTime()
  const mins = Math.floor(diffMs / 60000)
  return mins < 1
    ? t('lastUpdatedJustNow')
    : t('lastUpdatedMinutesAgo', { minutes: mins })
})
const shouldShowLastUpdated = computed(() => {
  if (!lastRefreshTime.value) return false
  const diffMs = now.value - new Date(lastRefreshTime.value).getTime()
  return diffMs > 10 * 60 * 1000 // 10 minutes
})




const refreshEverything = async () => {
  isRefreshing.value = true
  startRotating()

  const minDelay = new Promise(resolve => setTimeout(resolve, 3000)) // 3s guaranteed

  await Promise.all([
    fetchWaterData(),
    fetchEntries(),
    ensureTemperature(),
    minDelay
  ])

  const hour = new Date().getHours()
  await fetchPrediction(hour, waterTemperature.value ?? undefined, 3000)

  const now = new Date().toISOString()
  localStorage.setItem('lastRefresh', now)
  lastRefreshTime.value = now

  stopRotating()
  isRefreshing.value = false
}


const {
  waterTemperature,
  waterTemperatureLoading,
  waterTemperatureError,
  cachedAgeMinutes,
  ensureTemperature,
} = useTemperature()

const {
  entriesLoading,
  errorEntries: entriesError,
  entriesLoadingMessage,
  fetchEntries,
  addEntry,
  fetchPrediction,
  predictionLoading,
  predictionError,
  predictionHasBeenFetched,
  currentHourPrediction,
  todaysEntries,
  historyEntries,
} = useSurferEntries()



const {
  requestDate,
  currentWaterFlow,
  currentWaterLevel,
  showWaterLevelAlert,
  chartLabels,
  chartValues,
  fetchWaterData,
  waterDataLoading,
} = useWaterLevelData()


const surferCount = computed(() => Number(surferCountRaw.value))

const submitSurferCount = async () => {
  const count = surferCount.value
  if (!isNaN(count) && count >= 0) {
    submitting.value = true
    await fetchWaterData() // 🆕 Get fresh data before submitting
    await addEntry(
      count,
      undefined,
      currentWaterLevel.value ?? undefined,
      currentWaterFlow.value ?? undefined,
      waterTemperature.value ?? undefined
    )
    await fetchEntries()
    surferCountRaw.value = ''
    submitting.value = false
  }
}
onMounted(() => {
  // 1. Load last refresh from localStorage if available
  const saved = localStorage.getItem('lastRefresh')
  if (saved) {
    lastRefreshTime.value = saved
  }

  // 2. Run the actual full refresh
  refreshEverything()

  // 3. Update the "now" ref every minute so the "x min ago" label stays fresh
  setInterval(() => {
    now.value = Date.now()
  }, 60000)
})

</script>
