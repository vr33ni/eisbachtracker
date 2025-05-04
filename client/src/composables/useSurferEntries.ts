import { ref, computed } from 'vue'
import axios from 'axios'
import type { SurferEntryDto } from '@/dto/surfer-entry.dto'
import type { PredictionResponseDto } from '@/dto/prediction-response.dto'

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

export function useSurferEntries() {
  const entries = ref<SurferEntryDto[]>([])
  const entriesLoading = ref(false)
  const errorEntries = ref<string | null>(null)

  const predictionLoading = ref(false)
  const predictionError = ref<string | null>(null)
  const predictionHasBeenFetched = ref(false)
  const currentHourPrediction = ref<number | null>(null)
  const explanation = ref<Record<string, number> | null>(null) 

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

  const addEntry = async (
    count: number,
    time?: string,
    waterLevel?: number,
    waterFlow?: number,
    waterTemperature?: number,
  ) => {
    try {
      const body: any = {
        count,
        timestamp: time || new Date().toISOString(),
        water_level: waterLevel,
        water_flow: waterFlow,
      }

      if (waterTemperature !== undefined) {
        body.water_temperature = waterTemperature
      }

      const res = await axios.post(`${API_BASE_URL}/surfers`, body)
      if (!res.status.toString().startsWith('2')) throw new Error('Failed to add entry')

      await fetchEntries()
    } catch (err) {
      errorEntries.value = err instanceof Error ? err.message : 'Failed to submit entry'
    }
  }

  const fetchPrediction = async (hour: number, waterTemperature?: number, delay = 2500) => {
    predictionLoading.value = true
    predictionError.value = null
    predictionHasBeenFetched.value = false
    explanation.value = null  


    try {
      const url = new URL(`${API_BASE_URL}/surfers/predict`)
      url.searchParams.set('hour', hour.toString())
      if (waterTemperature !== undefined) {
        url.searchParams.set('water_temperature', waterTemperature.toString())
      }

      const [res] = await Promise.all([
        axios.get(url.toString()),
        new Promise(resolve => setTimeout(resolve, delay)), // optional UI delay
      ])

      const data = res.data as PredictionResponseDto
      currentHourPrediction.value = data.prediction
      explanation.value = data.explanation 
      return data
    } catch (err) {
      predictionError.value = err instanceof Error ? err.message : 'Failed to fetch prediction'
      currentHourPrediction.value = null
      explanation.value = null 
      return null
    } finally {
      predictionHasBeenFetched.value = true
      predictionLoading.value = false
    }
  }

  const todaysEntries = computed(() =>
    entries.value.filter((e) => new Date(e.timestamp).toDateString() === new Date().toDateString())
  )

  const historyEntries = computed(() =>
    entries.value.filter((e) => new Date(e.timestamp).toDateString() !== new Date().toDateString())
  )

  return {
    entries,
    entriesLoading,
    errorEntries,
    entriesLoadingMessage: 'Loading entries...',
    fetchEntries,
    addEntry,
    fetchPrediction,
    predictionLoading,
    predictionError,
    predictionHasBeenFetched,
    currentHourPrediction,
    explanation, 
    todaysEntries,
    historyEntries,
  }
}
