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

        <div v-if="waterTemperatureLoading">{{ loadingMessage }}</div>
        <div v-else-if="waterTemperatureError">❌ {{ waterTemperatureError }}</div>
        <div v-else>🌡️ {{ waterTemperature }} °C</div>
      </div>

      <!-- Alert -->
      <div v-if="showWaterLevelAlert" class="text-red-600 font-semibold text-base">
        🚨 Water level exceeds threshold!
      </div>

      <!-- Button -->
      <button @click="refreshEverything" :disabled="waterTemperatureLoading || waterLevelLoading"
        class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg text-lg font-medium shadow transition-all disabled:opacity-50 disabled:cursor-not-allowed">
        {{ waterTemperatureLoading || waterLevelLoading ? 'Refreshing...' : 'Refresh Data' }}
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

        <div v-if="isLoadingPrediction">📡 Loading surfer data...</div>

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

            <p v-if="isLoadingPrediction">Loading prediction...</p>
            <p v-else-if="currentHourPrediction !== null">
              Predicted surfers: <strong>{{ currentHourPrediction }}</strong>
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
import { useWaterLevelData } from '@/composables/useWaterLevelData'

const surferCountRaw = ref('')
const currentHourPrediction = ref<number | null>(null)
const isLoadingPrediction = ref(false)

// const isLoading = computed(() =>
//   waterTemperatureLoading.value || isLoadingPrediction.value || waterLevelLoading.value
// )


const {
  waterTemperature,
  ensureTemperature,
  loading: waterTemperatureLoading,
  error: waterTemperatureError,
  loadingMessage,
  fetchTemperature,
} = useTemperature()

const {
  entries: surferEntries,
  addEntry,
  fetchEntries,
  getPredictionForHour,
  loading: surfersLoading,
  error: surfersError,
} = useSurferEntries()

const {
  waterLevelText,
  waterFlowText,
  showWaterLevelAlert,
  chartLabels,
  chartValues,
  fetchWaterData: fetchWaterLevelData,
  loading: waterLevelLoading,
} = useWaterLevelData()


// 🔎 Entries from today
const todaysEntries = computed(() =>
  surferEntries.value.filter(entry => {
    const date = new Date(entry.timestamp)
    const today = new Date()
    return date.toDateString() === today.toDateString()
  })
)

const fetchPrediction = async () => {
  isLoadingPrediction.value = true
  try {
    const now = new Date()
    const hour = now.getHours()

    await ensureTemperature()

    const prediction = await getPredictionForHour(hour)
    if (prediction) {
      currentHourPrediction.value = prediction.prediction
      waterTemperature.value = prediction.water_temperature  // update local temp!
    } else {
      currentHourPrediction.value = null
    }
  } catch (err) {
    console.error("Failed to fetch prediction", err)
    currentHourPrediction.value = null
  } finally {
    isLoadingPrediction.value = false
  }
}

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

onMounted(async () => {
  await fetchWaterLevelData()
  await fetchEntries()
  await fetchPrediction()
})


const refreshEverything = () => {
  fetchWaterLevelData()
  fetchEntries()
  fetchPrediction()

}




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
