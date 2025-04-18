import { ref, computed } from 'vue'
import axios from 'axios'
import type { WaterTemperatureDto } from '@/dto/water-temperature.dto'

const API_BASE_URL = import.meta.env.VITE_BACKEND_API_URL
const STORAGE_KEY = 'cachedWaterTemperature'

export function useTemperature() {
  const waterTemperature = ref<number | null>(null)
  const waterTemperatureLoading = ref(false)
  const waterTemperatureError = ref<string | null>(null)
  const cacheTimestamp = ref<number | null>(null)

  const cachedAgeMinutes = computed(() => {
    if (!cacheTimestamp.value) return null
    return Math.floor((Date.now() - cacheTimestamp.value) / 60000)
  })

  const fetchTemperature = async () => {
    const cached = JSON.parse(localStorage.getItem(STORAGE_KEY) || 'null')
  
    if (cached && Date.now() - cached.timestamp < 1000 * 60 * 60) {
      waterTemperature.value = cached.temperature
      cacheTimestamp.value = cached.timestamp
      return
    }
  
    waterTemperatureLoading.value = true
    waterTemperatureError.value = null
  
    try {
      const res = await axios.get(`${API_BASE_URL}/conditions/water/temperature`)
      const data: WaterTemperatureDto = res.data
      waterTemperature.value = data.water_temperature
      cacheTimestamp.value = Date.now()
  
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        temperature: data.water_temperature,
        timestamp: cacheTimestamp.value,
      }))
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
    cachedAgeMinutes, 
  }
}
