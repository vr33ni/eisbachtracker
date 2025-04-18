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
  const historyTimestamps = ref<string[]>([])

  const waterDataLoading = ref(false)
  const error = ref<string | null>(null)

  const chartViewMode = ref<'hourly' | 'daily'>('hourly')

  const chartLabels = computed(() => {
    if (chartViewMode.value === 'hourly') {
      const today = new Date().toDateString()

      const filteredHistory = historyLabels.value.map((label, i) => {
        const date = new Date(historyTimestamps.value[i])
        return { label, value: historyValues.value[i], isToday: date.toDateString() === today }
      }).filter(d => d.isToday)

      const filteredLive = labels.value.map((label, i) => ({
        label,
        value: values.value[i],
      }))

      return [...filteredHistory.map(d => d.label), ...filteredLive.map(d => d.label)]
    }

    // 🗓️ Daily view: extract day portion for grouping
    const uniqueDays = new Set<string>()
    historyLabels.value.forEach((label) => {
      const day = label.split(',')[0]
      uniqueDays.add(day)
    })

    return Array.from(uniqueDays)
  })

  const chartValues = computed(() => {
    if (chartViewMode.value === 'hourly') {
      const today = new Date().toDateString()

      const filteredHistory = historyLabels.value.map((label, i) => {
        const date = new Date(historyTimestamps.value[i])
        return { label, value: historyValues.value[i], isToday: date.toDateString() === today }
      }).filter(d => d.isToday)

      return [...filteredHistory.map(d => d.value), ...values.value]
    }

    // 🧮 Daily average mode
    const grouped: Record<string, number[]> = {}

    historyLabels.value.forEach((label, i) => {
      const day = label.split(',')[0]
      if (!grouped[day]) grouped[day] = []
      grouped[day].push(historyValues.value[i])
    })

    return Object.values(grouped).map((vals) => {
      const sum = vals.reduce((a, b) => a + b, 0)
      return Math.round(sum / vals.length)
    })
  })

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

      if (showWaterLevelAlert.value && currentWaterLevel.value !== null)
        notifyUser(currentWaterLevel.value)

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

  const fetchHistoricalWaterData = async () => {
    try {
      const res = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/conditions/water/history`)
      if (!res.ok) throw new Error('Backend error')
  
      const data = await res.json()
  
      historyLabels.value = []
      historyValues.value = []
      historyTimestamps.value = []
  
      const parsedData = data
        .map((entry: any) => {
          const [datePart, timePart] = entry.DateTime.split(' ')
          const [day, month, year] = datePart.split('.')
  
          const isoString = `${year}-${month}-${day}T${timePart}:00`
  
          const parsedDate = new Date(isoString)
          if (isNaN(parsedDate.getTime())) return null
  
          return {
            label: parsedDate.toLocaleString(undefined, {
              day: '2-digit',
              month: '2-digit',
              hour: '2-digit',
              minute: '2-digit',
            }),
            timestamp: parsedDate.toISOString(),
            value: entry.Value,
          }
        })
        .filter(Boolean)
        .sort((a: any, b: any) => new Date(a!.timestamp).getTime() - new Date(b!.timestamp).getTime())
  
      for (const row of parsedData) {
        historyLabels.value.push(row!.label)
        historyValues.value.push(row!.value)
        historyTimestamps.value.push(row!.timestamp)
      }
    } catch (err) {
      console.error('❌ Failed to fetch historical water data:', err)
    }
  }
  

  return {
    requestDate,
    currentWaterLevel,
    currentWaterFlow,
    showWaterLevelAlert,
    fetchWaterData,
    fetchHistoricalWaterData,
    waterDataLoading,
    error,
    chartLabels,
    chartValues,
    chartViewMode, // 👈 so you can switch between hourly / daily from the UI
  }
}
