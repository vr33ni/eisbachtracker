import { ref } from 'vue'
import axios from 'axios'
import type { SurferEntryDto } from '@/dto/surfer-entry.dto'
import type { PredictionResponseDto } from '@/dto/prediction-response.dto'

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

const predictionMessages = [
  '🔮 Watching the river...',
  '🤔 Estimating surf hype...',
  '📈 Crunching the numbers...',
  '🐟 Asking fish for advice...',
  '📷 Checking river cams...',
  '🌊 Feeling the vibes...',
  '⏳ Counting wetsuits...',
]

export function useSurferEntries() {
  const entries = ref<SurferEntryDto[]>([])

  const entriesLoading = ref(false)
  const errorEntries = ref<string | null>(null)

  const predictionLoading = ref(false)
  const predictionError = ref<string | null>(null)
  const predictionLoadingMessage = ref(predictionMessages[0])
  let interval: ReturnType<typeof setInterval> | null = null

 
  const fetchEntries = async () => {
    entriesLoading.value = true
    errorEntries.value = null

    try {
      const res = await axios.get(`${API_BASE_URL}/surfers`)
      entries.value = res.data
    } catch (err) {
      errorEntries.value = err instanceof Error ? err.message : 'Failed to fetch entries'
    } finally {
      entriesLoading.value = false
    }
  }

  const addEntry = async (count: number, time?: string) => {
    try {
      const body = {
        count,
        timestamp: time || new Date().toISOString(),
      }

      const res = await axios.post(`${API_BASE_URL}/surfers`, body)
      if (!res.status.toString().startsWith('2')) throw new Error('Failed to add entry')

      await fetchEntries()
    } catch (err) {
      errorEntries.value = err instanceof Error ? err.message : 'Failed to submit entry'
    }
  }

  const getPredictionForHour = async (hour: number, waterTemperature?: number) => {
    predictionLoading.value = true
    predictionError.value = null

    let i = 0
    interval = setInterval(() => {
      predictionLoadingMessage.value = predictionMessages[i % predictionMessages.length]
      i++
    }, 2500)

    try {
      const url = new URL(`${API_BASE_URL}/surfers/predict`)
      url.searchParams.set('hour', hour.toString())
      if (waterTemperature !== undefined) {
        url.searchParams.set('water_temperature', waterTemperature.toString())
      }

      const res = await axios.get(url.toString())
      return res.data as PredictionResponseDto
    } catch (err) {
      predictionError.value = err instanceof Error ? err.message : 'Failed to fetch prediction'
      return null
    } finally {
      predictionLoading.value = false
      if (interval) clearInterval(interval)
    }
  }

  return {
    entries,
    entriesLoading,
    errorEntries,
    entriesLoadingMessage: 'Loading entries...',
    fetchEntries,
    addEntry,
    getPredictionForHour,
    predictionLoading,
    predictionError,
    predictionLoadingMessage,
   }
}
