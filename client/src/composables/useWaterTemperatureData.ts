import { ref } from 'vue'
import axios from 'axios'

export function useTemperature() {
  const temperature = ref(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchTemperature = async () => {
    loading.value = true
    error.value = null

    try {
      const res = await axios.get('http://localhost:3000/api/temperature') // ← your backend endpoint
      temperature.value = res.data.temperature
    } catch (err) {
      error.value = (err instanceof Error ? err.message : 'Failed to fetch temperature')
    } finally {
      loading.value = false
    }
  }

  return {
    temperature,
    loading,
    error,
    fetchTemperature,
  }
}
