<template>
<div class="space-y-2">
  <!-- Alert -->
  <div v-if="showWaterLevelAlert" class="text-red-600 text-sm font-semibold">
    🚨 {{ t('lowTideAlert') }}
  </div>

  <!-- Water Level -->
  <p class="text-md text-gray-700 dark:text-gray-300">
    💧 {{ t('waterLevel') }}: 
    <span v-if="waterDataLoading || currentWaterLevel === null">Loading...</span>
    <span v-else>{{ currentWaterLevel }} cm</span>
  </p>

  <!-- Water Flow -->
  <p class="text-md text-gray-700 dark:text-gray-300">
    🌊 {{ t('waterFlow') }}: 
    <span v-if="waterDataLoading || currentWaterFlow === null">Loading...</span>
    <span v-else>{{ currentWaterFlow }} m³/s</span>
  </p>

  <!-- Water Data Timestamp -->
  <p v-if="requestDate" class="text-xs text-gray-500">
    ⏱️ {{ t('requestDate') }}: {{ requestDate }}
  </p>

  <!-- Optional Separator -->
  <hr class="border-t border-gray-300 dark:border-gray-700 my-2" />

  <!-- Temperature -->
  <p class="text-md text-gray-700 dark:text-gray-300">
    🌡️ {{ t('waterTemp') }}:
    <span v-if="waterTemperatureLoading" class="animate-pulse text-blue-600">
      {{ t('loading') }}
    </span>
    <span v-else-if="waterTemperatureError" class="text-red-600">
      ❌ {{ waterTemperatureError }}
    </span>
    <span v-else>
      {{ waterTemperature }} °C
    </span>
  </p>

  <!-- Cached Age -->
  <p v-if="cachedAgeMinutes !== null && cachedAgeMinutes >= 0" class="text-xs text-gray-500">
    💾 {{ t('cachedAgo', { minutes: cachedAgeMinutes }) }}
  </p>

  <!-- Expandable Chart -->
  <ExpandableCard :title="`📈 ${t('waterChartTitle')}`">
    <WaterChartCard :labels="chartLabels" :values="chartValues" />
  </ExpandableCard>
</div>

</template>

<script setup lang="ts">
import ExpandableCard from './ExpandableCard.vue'
import WaterChartCard from './WaterChartCard.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  currentWaterLevel: number | null
  currentWaterFlow: number | null
  requestDate: Date | string
  showWaterLevelAlert: boolean
  waterTemperature: number | null
  waterTemperatureLoading: boolean
  waterTemperatureError: string | null
  cachedAgeMinutes: number | null
  chartLabels: string[]
  chartValues: number[]  
  waterDataLoading: boolean
}>()

</script>
