import { computed, ref } from 'vue'
import { DateTime } from 'luxon'

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

   // merge here for clean chart usage
   const chartLabels = computed(() => [...historyLabels.value, ...currentLabels.value])
   const chartValues = computed(() => [...historyValues.value, ...currentValues.value])
   
  const pegelAlarmApiUrl = import.meta.env.VITE_PEGEL_ALARM_API_URL
 
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
      const response = await fetch(pegelAlarmApiUrl)
      const data = await response.json()

      const station = data?.payload?.stations?.[0]
      const waterLevel = station?.data?.[0]?.value
      const waterFlow = station?.data?.[1]?.value

      waterLevelText.value = `${waterLevel} cm`
      waterFlowText.value = `${waterFlow} m³/s`

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

  // const fetchWaterHistory = async () => {
  //   try {
  //     const res = await fetch(hndApiUrl)
  //     const csv = await res.text()
  
  //     const lines = csv.trim().split('\n').slice(1) // remove header
  //     const parsed = lines.map(line => {
  //       const [timestamp, value] = line.split(',')
      
  //       // Parse format: "14.04.2025T18:15:00+0200"
  //       const dt = DateTime.fromFormat(timestamp, 'dd.MM.yyyy\'T\'HH:mm:ssZZ')
  //       return {
  //         label: dt.isValid ? dt.toFormat('HH:mm') : 'Invalid',
  //         value: Number(value),
  //       }
  //     })

  //     console.log("parsed: ", parsed)
      
  
  //     historyLabels.value = parsed.map(p => p.label)
  //     historyValues.value = parsed.map(p => p.value)
  //   } catch (err) {
  //     console.error("Failed to fetch HND water history", err)
  //   }
  // }
  
  

  return {
    waterLevelText,
    waterFlowText,
    showWaterLevelAlert,
    fetchWaterData,
    // fetchWaterHistory,
    loading,
    error,
    chartLabels,
    chartValues,
  }
  
}
