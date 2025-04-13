import { ref } from 'vue'
import axios from 'axios'
import type { WaterTemperatureDto } from '@/dto/water-temperature.dto'

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

const messages = [
  '🌐 Contacting the Bavarian Water Lords...',
  '📡 Negotiating cookie treaties...',
  '📬 Enqueueing top-secret data packet...',
  '🔄 Waiting for temperature to be deemed worthy...',
  '📦 Unzipping meteorological mysteries...',
  '📊 Decoding aquatic runes...',
  '🌡️ Extracting the sacred temperature...',
  '🧊 Counting water molecules...',
  '🐟 Interviewing local fish...',
]
let fetchedOnce = false

export function useTemperature() {
  const waterTemperature = ref<number | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const loadingMessage = ref('Loading...')

  const fetchTemperature = async () => {
    // Prevent double call
    if (fetchedOnce) return
    fetchedOnce = true

    loading.value = true
    error.value = null
    waterTemperature.value = null

    try {
      const res = await axios.get(`${API_BASE_URL}/conditions/water-temperature`)
      waterTemperature.value = res.data.water_temperature
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch temperature'
    } finally {
      loading.value = false
    }
  }

  async function ensureTemperature() {
    if (waterTemperature.value === null) {
      await fetchTemperature()
    }
  }

  return {
    waterTemperature,
    loading,
    error,
    loadingMessage,
    fetchTemperature,
    ensureTemperature,
  }
}
