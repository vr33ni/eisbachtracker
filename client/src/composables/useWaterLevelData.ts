import { computed, ref } from 'vue'

export function useWaterLevelData() {
  const showWaterLevelAlert = ref(false)

  const currentWaterLevel = ref<number | null>(null)
  const currentWaterFlow = ref<number | null>(null)
  const requestDate = ref<string>('')

  const labels = ref<string[]>([])
  const values = ref<number[]>([])

  const historyLabels = ref<string[]>([])
  const historyValues = ref<number[]>([])

  const waterDataLoading = ref(false)
  const error = ref<string | null>(null)

  const chartLabels = computed(() => [...historyLabels.value, ...labels.value])
  const chartValues = computed(() => [...historyValues.value, ...values.value])

  const notifyUser = (waterLevel: number) => {
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification('🌊 Eisbach Alert', {
        body: `Water level is currently ${waterLevel}cm — longboard only 🏄‍♀️`,
      })
    }
  }

  const fetchWaterData = async () => {
    waterDataLoading.value = true
    error.value = null

    try {
      const res = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/conditions/water`)
      if (!res.ok) throw new Error('Backend error')

      const data = await res.json()

      requestDate.value = new Date(data.request_date).toLocaleTimeString()
      currentWaterLevel.value = data.water_level
      currentWaterFlow.value = data.water_flow

      showWaterLevelAlert.value = currentWaterLevel.value !== null && currentWaterLevel.value <= 140

      if (showWaterLevelAlert.value && currentWaterLevel.value !== null) notifyUser(currentWaterLevel.value)

      const label = new Date().toLocaleTimeString()
      labels.value.push(label)
      if (currentWaterLevel.value !== null) {
        values.value.push(currentWaterLevel.value)
      }

      if (labels.value.length > 10) labels.value.shift()
      if (values.value.length > 10) values.value.shift()

    } catch (err) {
      console.error('Error fetching water data from backend:', err)
      error.value = 'Failed to fetch water data'
    } finally {
      waterDataLoading.value = false
    }
  }

  return {
    requestDate,
    currentWaterLevel,
    currentWaterFlow,
    showWaterLevelAlert,
    fetchWaterData,
    waterDataLoading,
    error,
    chartLabels,
    chartValues,
  }
}
