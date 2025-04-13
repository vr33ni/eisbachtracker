import { ref } from 'vue'

export function useWaterLevelData() {
  const waterLevelText = ref('Loading...')
  const waterFlowText = ref('Loading...')
  const showWaterLevelAlert = ref(false)
  const chartLabels = ref<string[]>([])
  const chartValues = ref<number[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const apiUrl = import.meta.env.VITE_PEGEL_API_URL

  const notifyUser = (waterLevel: number) => {
    if ('Notification' in window) {
      if (Notification.permission === 'granted') {
        new Notification('🌊 Eisbach Alert', {
          body: `Water level is currently ${waterLevel}cm — longboard only 🏄‍♀️`,
        })
      }
    }
  }

  const fetchWaterData = async () => {
    loading.value = true
    try {
      const response = await fetch(apiUrl)
      const data = await response.json()

      const station = data?.payload?.stations?.[0]
      const waterLevel = station?.data?.[0]?.value
      const waterFlow = station?.data?.[1]?.value

      waterLevelText.value = `Current Water Level: ${waterLevel} cm`
      waterFlowText.value = `Current Water Flow: ${waterFlow} m³/s`

      showWaterLevelAlert.value = waterLevel <= 140

      if (showWaterLevelAlert.value) notifyUser(waterLevel)

      chartLabels.value.push(new Date().toLocaleTimeString())
      chartValues.value.push(waterLevel)

      if (chartLabels.value.length > 10) chartLabels.value.shift()
      if (chartValues.value.length > 10) chartValues.value.shift()
    } catch (err) {
      console.error('Error fetching water data:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    waterLevelText,
    waterFlowText,
    showWaterLevelAlert,
    chartLabels,
    chartValues,
    fetchWaterData,
    loading,
    error,
  }
}
