<template>
  <div id="app"
    class="min-h-screen bg-gradient-to-b from-blue-100 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-md space-y-8 text-center">

      <h1 class="text-4xl sm:text-5xl font-bold tracking-tight text-blue-600 dark:text-blue-300">
        Eisbach Tracker
      </h1>

      <!-- Water Data -->
      <div class="space-y-2">
        <div v-if="showWaterLevelAlert" class="text-red-600 text-sm font-semibold">
          🚨 Low tide alert
        </div>
        <p class="text-lg text-gray-700 dark:text-gray-300">
          🌊 Water Level: {{ waterLevelText }}
        </p>

        <p class="text-lg text-gray-700 dark:text-gray-300">
          💧 Water Flow: {{ waterFlowText }}
        </p>

        <p class="text-lg text-gray-700 dark:text-gray-300">Current Water Temperature:</p>

        <div v-if="waterTemperatureLoading">{{ temperatureLoadingMessage }}</div>
        <div v-else-if="waterTemperatureError">❌ {{ waterTemperatureError }}</div>
        <div v-else>🌡️ {{ waterTemperature }} °C</div>
      </div>



      <button @click="refreshEverything" :disabled="isRefreshing"
        class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg text-lg font-medium shadow transition-all disabled:opacity-50 disabled:cursor-not-allowed">
        {{ isRefreshing ? 'Refreshing...' : 'Refresh Data' }}
      </button>

      <WaterChart :labels="chartLabels" :values="chartValues" />

      <div class="space-y-4 pt-8 border-t border-gray-300 dark:border-gray-600 mt-6 text-left">
        <h2 class="text-2xl font-semibold text-blue-700 dark:text-blue-300">🧍 Surfer Spotter</h2>

        <form @submit.prevent="submitSurferCount" class="flex items-center gap-2">
          <input v-model="surferCountRaw" @input="onInputNumeric" type="text" placeholder="Number of surfers"
            inputmode="numeric" pattern="[0-9]*"
            class="px-3 py-2 rounded border dark:bg-gray-800 dark:border-gray-600" />
          <button type="submit" :disabled="submitting || surferCount === null"
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 justify-center min-w-[120px]">
            <svg v-if="submitting" class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg"
              fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
            </svg>
            <span>{{ submitting ? 'Submitting...' : 'Submit' }}</span>
          </button>

        </form>

        <div v-if="entriesLoading">{{ entriesLoadingMessage }}</div>
        <div v-else-if="entriesError" class="text-red-500">❌ {{ entriesError }}</div>

        <div v-else>
          <ul v-if="todaysEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1">
            <li v-for="entry in todaysEntries" :key="entry.timestamp">
              {{ new Date(entry.timestamp).toLocaleTimeString() }} — {{ entry.count }} surfers
            </li>
          </ul>
          <p v-else class="text-gray-500 text-sm">No entries yet today</p>

          <div class="mt-4 text-left">
            <h3 class="font-semibold text-blue-700 dark:text-blue-300">📊 Prediction</h3>

            <div v-if="predictionLoading">{{ predictionLoadingMessage }}</div>
            <div v-else-if="predictionError">❌ {{ predictionError }}</div>
            <div v-else-if="currentHourPrediction !== null">
              Predicted surfers: <strong>{{ currentHourPrediction }}</strong>
            </div>
            <div v-else class="text-gray-500">
              Not enough data to predict crowd
            </div>
          </div>
          <br /> <br />

        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import WaterChart from './components/WaterChart.vue'
import { useTemperature } from '@/composables/useWaterTemperatureData'
import { useSurferEntries } from '@/composables/useSurferEntries'
import { useWaterLevelData } from '@/composables/useWaterLevelData'

const surferCountRaw = ref('')
const currentHourPrediction = ref<number | null>(null)
const submitting = ref(false)

const { waterTemperature, ensureTemperature, loading: waterTemperatureLoading, error: waterTemperatureError, loadingMessage: temperatureLoadingMessage } = useTemperature()
const { entries, addEntry, fetchEntries, getPredictionForHour, entriesLoading, errorEntries: entriesError, entriesLoadingMessage, predictionLoading, predictionError, predictionLoadingMessage } = useSurferEntries()
const { waterLevelText, waterFlowText, showWaterLevelAlert, chartLabels, chartValues, fetchWaterData, loading: waterLevelLoading } = useWaterLevelData()

const isRefreshing = computed(() => waterTemperatureLoading.value || entriesLoading.value || waterLevelLoading.value)

const todaysEntries = computed(() => entries.value.filter(e => new Date(e.timestamp).toDateString() === new Date().toDateString()))

const fetchPrediction = async () => {
  try {
    const now = new Date()
    const hour = now.getHours()

    //await ensureTemperature() // still smart here!

    const prediction = await getPredictionForHour(hour)

    if (prediction) {
      currentHourPrediction.value = prediction.prediction
      waterTemperature.value = prediction.water_temperature // optionally update cached temp
    } else {
      currentHourPrediction.value = null
    }
  } catch (err) {
    console.error('Failed to fetch prediction', err)
    currentHourPrediction.value = null
  }
}

const surferCount = computed(() => Number(surferCountRaw.value))

const onInputNumeric = (e: Event) => {
  const target = e.target as HTMLInputElement
  target.value = target.value.replace(/[^0-9]/g, '')
  surferCountRaw.value = target.value
}

const submitSurferCount = async () => {
  const count = surferCount.value
  if (!isNaN(count) && count >= 0) {
    submitting.value = true
    await addEntry(count)
    await fetchEntries()
    surferCountRaw.value = ''
    submitting.value = false
  }
}

const refreshEverything = async () => {
  await fetchWaterData()
  await fetchEntries()
  await fetchPrediction()
}

onMounted(refreshEverything)

const notifyUser = (waterLevel: number) => {
  if ("Notification" in window) {
    if (Notification.permission === "granted") {
      new Notification("🌊 Eisbach Alert", {
        body: `Water level is currently ${waterLevel}cm — longboard only 🏄‍♀️`,
      });
    } else if (Notification.permission !== "denied") {
      Notification.requestPermission().then((permission) => {
        if (permission === "granted") {
          new Notification("🌊 Eisbach Alert", {
            body: `Water level is currently ${waterLevel}cm — longboard only 🏄‍♀️`,
          });
        }
      });
    }
  }
};

</script>

<style scoped>
.alert {
  font-size: 1.25rem;
}

html,
body,
#app {
  height: 100%;
  margin: 0;
}
</style>
