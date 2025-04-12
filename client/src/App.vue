<template>
  <div id="app"
    class="min-h-screen bg-gradient-to-b from-blue-100 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-md space-y-8 text-center">
      <!-- Title -->
      <h1 class="text-4xl sm:text-5xl font-bold tracking-tight text-blue-600 dark:text-blue-300">
        Eisbach Tracker
      </h1>

      <!-- Water Data -->
      <div class="space-y-2">
        <p class="text-lg text-gray-700 dark:text-gray-300">{{ waterLevelText }}</p>
        <p class="text-lg text-gray-700 dark:text-gray-300">{{ waterFlowText }}</p>
        <p class="text-lg text-gray-700 dark:text-gray-300">Current Water Temperature:</p>
        <div v-if="isLoading">{{ loadingMessage }}</div>
        <div v-else-if="error">❌ {{ error }}</div>
        <div v-else>🌡️ {{ temperature }} °C</div>
      </div>

      <!-- Alert -->
      <div v-if="showWaterLevelAlert" class="text-red-600 font-semibold text-base">
        🚨 Water level exceeds threshold!
      </div>

      <!-- Button -->
      <button @click="refreshEverything" :disabled="loading"
        class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg text-lg font-medium shadow transition-all disabled:opacity-50 disabled:cursor-not-allowed">
        {{ loading ? 'Refreshing...' : 'Refresh Data' }}
      </button>

      <!-- Chart (moved inside box) -->
      <div class="pt-6">
        <WaterChart :labels="chartLabels" :values="chartValues" />
      </div>

      <!-- Surfer Count Tracker -->
      <div class="space-y-4 pt-8 border-t border-gray-300 dark:border-gray-600 mt-6 text-left">
        <h2 class="text-2xl font-semibold text-blue-700 dark:text-blue-300">🧍 Surfer Spotter</h2>

        <form @submit.prevent="submitSurferCount" class="flex items-center gap-2">
          <input v-model="surferCountRaw" @input="onInputNumeric" type="text" placeholder="Number of surfers"
            inputmode="numeric" pattern="[0-9]*"
            class="px-3 py-2 rounded border dark:bg-gray-800 dark:border-gray-600" />

          <button type="submit" :disabled="submitting || surferCount === null"
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded disabled:opacity-50 disabled:cursor-not-allowed">
            Submit
          </button>
        </form>

        <div v-if="surfersLoading">📡 Loading surfer data...</div>

        <div v-else-if="surfersError" class="text-red-500">❌ {{ surfersError }}</div>

        <!-- 👇 Updated part below -->
        <div v-else>
          <!-- Today's Entries -->
          <ul v-if="todaysEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1">
            <li v-for="entry in todaysEntries" :key="entry.timestamp">
              {{ new Date(entry.timestamp).toLocaleTimeString() }} — {{ entry.count }} surfers
            </li>
          </ul>
          <p v-else class="text-gray-500 text-sm">No entries yet today</p>

          <!-- Prediction -->
          <div class="mt-4 text-left">
            <h3 class="font-semibold text-blue-700 dark:text-blue-300">📊 Prediction</h3>
            <p v-if="currentHourPrediction !== null">
              Based on the last hour: <strong>{{ currentHourPrediction }}</strong> surfers
            </p>
            <p v-else class="text-gray-500">Not enough data to predict crowd</p>
          </div>
        </div>


      </div>
      <br>
    </div>

  </div>
</template>




<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import WaterChart from './components/WaterChart.vue'
import { useTemperature } from '@/composables/useWaterTemperatureData'
import { useSurferEntries } from '@/composables/useSurferEntries'

const chartLabels = ref<string[]>([])
const chartValues = ref<number[]>([])

const waterLevelText = ref('Loading...')
const waterFlowText = ref('Loading...')
const showWaterLevelAlert = ref(false)
const { temperature, loading, error, loadingMessage, fetchTemperature } = useTemperature()

const {
  entries: surferEntries,
  addEntry,
  fetchEntries,
  getPredictionForHour,
  loading: surfersLoading,
  error: surfersError,
} = useSurferEntries()


const currentHour = new Date().getHours()
const predictedSurfers = ref<number | null>(null)


// 🔎 Entries from today
const todaysEntries = computed(() =>
  surferEntries.value.filter(entry => {
    const date = new Date(entry.timestamp)
    const today = new Date()
    return date.toDateString() === today.toDateString()
  })
)

// 📊 Prediction from past hour
const currentHourPrediction = computed(() => {
  const now = new Date()
  const hourAgo = new Date(now.getTime() - 60 * 60 * 1000)

  const lastHourEntries = surferEntries.value.filter(entry => {
    const time = new Date(entry.timestamp)
    return time > hourAgo && time <= now
  })

  if (!lastHourEntries.length) return null

  const avg = lastHourEntries.reduce((sum, e) => sum + e.count, 0) / lastHourEntries.length
  return Math.round(avg)
})


const surferCountRaw = ref('')
const surferCount = computed(() => Number(surferCountRaw.value))

const onInputNumeric = (e: Event) => {
  const target = e.target as HTMLInputElement
  target.value = target.value.replace(/[^0-9]/g, '') // Only keep digits
  surferCountRaw.value = target.value
}
const submitting = ref(false)


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

const apiUrl = import.meta.env.VITE_PEGEL_API_URL

onMounted(async () => {
  fetchTemperature()
  fetchEntries()
  const temp = temperature.value // might be null
  const prediction = await getPredictionForHour(currentHour, temp ?? undefined)
  predictedSurfers.value = prediction
})

const refreshEverything = () => {
  checkWaterLevel()
  fetchTemperature()
}


const isLoading = computed(() => temperature.value === null && !error.value)


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


const fetchWaterData = async () => {
  try {
    const response = await fetch(apiUrl)
    const data = await response.json()

    if (data?.payload?.stations?.length > 0) {
      const station = data.payload.stations[0]
      const waterLevel = station.data[0]?.value
      const waterFlow = station.data[1]?.value

      waterLevelText.value = `Current Water Level: ${waterLevel} cm`
      waterFlowText.value = `Current Water Flow: ${waterFlow} m³/s`

      if (waterLevel <= 140) {
        showWaterLevelAlert.value = true
        notifyUser(waterLevel)
      } else {
        showWaterLevelAlert.value = false
      }

      // ✅ Push to chart data
      chartLabels.value.push(new Date().toLocaleTimeString())
      chartValues.value.push(waterLevel)

      // ✅ Optional: Keep it short
      if (chartLabels.value.length > 10) chartLabels.value.shift()
      if (chartValues.value.length > 10) chartValues.value.shift()

    } else {
      console.error('Stations data is missing or undefined')
    }
  } catch (error) {
    console.error('Error fetching water level data:', error)
  }
}


const checkWaterLevel = () => {
  fetchWaterData()
}



fetchWaterData()
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
