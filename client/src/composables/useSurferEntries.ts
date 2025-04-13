import { ref } from 'vue'
import type { SurferEntryDto } from '@/dto/surfer-entry.dto'
import type { PredictionResponseDto } from '@/dto/prediction-response.dto'

const entries = ref<SurferEntryDto[]>([])

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

export function useSurferEntries() {
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchEntries = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE_URL}/surfers`)
      if (!res.ok) throw new Error('Failed to fetch')
      const data: SurferEntryDto[] = await res.json()
      entries.value = data
    } catch (err: any) {
      error.value = err.message || 'Unexpected error'
    } finally {
      loading.value = false
    }
  }

  const addEntry = async (count: number, time?: string) => {
    try {
      const body = {
        count,
        timestamp: time || new Date().toISOString(),
      }

      const res = await fetch(`${API_BASE_URL}/surfers`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
      })

      if (!res.ok) throw new Error('Failed to add entry')

      await fetchEntries()
    } catch (err: any) {
      error.value = err.message || 'Failed to submit entry'
    }
  }

  const getPredictionForHour = async (hour: number, temp?: number) => {
    try {
      const url = new URL(`${API_BASE_URL}/surfers/predict`)
      url.searchParams.set('hour', hour.toString())
  
      const res = await fetch(url.toString())
      if (!res.ok) throw new Error('Failed to fetch prediction')
  
      const data = await res.json() as PredictionResponseDto
  
      return data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to get prediction'
      return null
    }
  }
  

  return {
    entries,
    loading,
    error,
    fetchEntries,
    addEntry,
    getPredictionForHour,
  }
}
