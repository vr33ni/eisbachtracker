import { ref } from 'vue'
import axios from 'axios'
import type { WaterTemperatureDto } from '@/dto/water-temperature.dto'

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL

export function useTemperature() {
  const waterTemperature = ref<number | null>(null)
  const waterTemperatureLoading = ref(false)
  const waterTemperatureError = ref<string | null>(null)

  const fetchTemperature = async () => {
    // Prevent re-fetch if already loaded
    if (waterTemperature.value !== null) return

    waterTemperatureLoading.value = true
    waterTemperatureError.value = null

    try {
      const res = await axios.get(`${API_BASE_URL}/conditions/water-temperature`)
      const data: WaterTemperatureDto = res.data
      waterTemperature.value = data.water_temperature
    } catch (err) {
      waterTemperatureError.value = err instanceof Error ? err.message : 'Failed to fetch temperature'
    } finally {
      waterTemperatureLoading.value = false
    }
  }

  const ensureTemperature = async () => {
    if (waterTemperature.value === null) {
      await fetchTemperature()
    }
  }

  return {
    waterTemperature,
    waterTemperatureLoading,
    waterTemperatureError,
    fetchTemperature,
    ensureTemperature,
  }
}
