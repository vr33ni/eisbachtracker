import { computed, ref } from 'vue'

export function useWaterLevelData() {
  const waterLevelText = ref('Loading...')
  const waterFlowText = ref('Loading...')
  const showWaterLevelAlert = ref(false)

  const currentLabels = ref<string[]>([])
  const currentValues = ref<number[]>([])

  const historyLabels = ref<string[]>([])
  const historyValues = ref<number[]>([])

  const loading = ref(false)
  const error = ref<string | null>(null)

  const chartLabels = computed(() => [...historyLabels.value, ...currentLabels.value])
  const chartValues = computed(() => [...historyValues.value, ...currentValues.value])

  const pegelAlarmApiUrl = import.meta.env.VITE_PEGEL_ALARM_API_URL

  const notifyUser = (waterLevel: number) => {
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification('🌊 Eisbach Alert', {
        body: `Water level is currently ${waterLevel}cm — longboard only 🏄‍♀️`,
      })
    }
  }

  const fetchWaterData = async () => {
    loading.value = true
    try {
      const response = await fetch(pegelAlarmApiUrl)
      const data = await response.json()

      const station = data?.payload?.stations?.[0]
      const waterLevel = station?.data?.[0]?.value
      const waterFlow = station?.data?.[1]?.value

      waterLevelText.value = `${waterLevel} cm`
      waterFlowText.value = `${waterFlow} m³/s`
      showWaterLevelAlert.value = waterLevel <= 140

      if (showWaterLevelAlert.value) notifyUser(waterLevel)

      const label = new Date().toLocaleTimeString()

      currentLabels.value.push(label)
      currentValues.value.push(waterLevel)

      if (currentLabels.value.length > 10) currentLabels.value.shift()
      if (currentValues.value.length > 10) currentValues.value.shift()
    } catch (err) {
      console.error('Error fetching water data:', err)
      error.value = 'Failed to fetch water data'
    } finally {
      loading.value = false
    }
  }

  return {
    waterLevelText,
    waterFlowText,
    showWaterLevelAlert,
    fetchWaterData,
    loading,
    error,
    chartLabels,
    chartValues,
  }
}
