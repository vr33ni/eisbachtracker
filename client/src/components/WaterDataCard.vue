<template>
  <div class="space-y-2">
    <!-- Alert -->
    <div v-if="showWaterLevelAlert" class="text-red-600 text-sm font-semibold">
      🚨 Low tide alert
    </div>

    <p class="text-md text-gray-700 dark:text-gray-300">🌊 Water Level: {{ waterLevelText }}</p>
    <p class="text-md text-gray-700 dark:text-gray-300">💧 Water Flow: {{ waterFlowText }}</p>

    <!-- Temperature -->
    <p class="text-md text-gray-700 dark:text-gray-300">
      🌡️ Water Temperature:
      <span v-if="waterTemperatureLoading" class="animate-pulse text-blue-600">Loading...</span>
      <span v-else-if="waterTemperatureError" class="text-red-600">❌ {{ waterTemperatureError }}</span>
      <span v-else>{{ waterTemperature }} °C</span>
    </p>

    <!-- Expandable Chart -->
    <ExpandableCard title="📈 Water Level History">
  <WaterChartCard :labels="chartLabels" :values="chartValues" />
</ExpandableCard>

  </div>
</template>

<script setup lang="ts">
import ExpandableCard from './ExpandableCard.vue'
import WaterChartCard from './WaterChartCard.vue'

defineProps<{
  waterLevelText: string
  waterFlowText: string
  showWaterLevelAlert: boolean
  waterTemperature: number | null
  waterTemperatureLoading: boolean
  waterTemperatureError: string | null
  chartLabels: string[]
  chartValues: number[]
}>()
</script>
