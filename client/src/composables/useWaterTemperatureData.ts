import { ref } from 'vue'
import axios from 'axios'

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

export function useTemperature() {
  const temperature = ref(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const apiUrl = import.meta.env.VITE_BACKEND_API_URL
  const loadingMessage = ref(messages[0])
  let interval: ReturnType<typeof setInterval> | null = null

  const fetchTemperature = async () => {
    loading.value = true
    error.value = null
    temperature.value = null

    // Rotate loading message
    let i = 0
    interval = setInterval(() => {
      loadingMessage.value = messages[i % messages.length]
      i++
    }, 3000)

    try {
      const res = await axios.get(apiUrl)
      temperature.value = res.data.temperature
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch temperature'
    } finally {
      loading.value = false
      if (interval) clearInterval(interval)
    }
  }

  return {
    temperature,
    loading,
    error,
    loadingMessage,
    fetchTemperature,
  }
}
