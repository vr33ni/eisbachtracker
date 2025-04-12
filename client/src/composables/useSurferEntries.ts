import { ref } from 'vue'

type SurferEntry = {
  timestamp: string
  count: number
}

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

export function useSurferEntries() {
  const entries = ref<SurferEntry[]>([]) // 🧠 typed!
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchEntries = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_BASE_URL}/surfers`)
      if (!res.ok) throw new Error('Failed to fetch')
      entries.value = await res.json()
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

  const getPredictionForHour = async (hour: number, temperature?: number) => {
    try {
      const url = new URL(`${API_BASE_URL}/surfers/predict`)
      url.searchParams.set('hour', hour.toString())
      if (temperature !== undefined) {
        url.searchParams.set('temperature', temperature.toString())
      }

      const res = await fetch(url.toString())
      if (!res.ok) throw new Error('Failed to fetch prediction')

      const data = await res.json()
      return data.prediction
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
