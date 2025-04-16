<template>
  <div class="space-y-2">
    <!-- Alert -->
    <div v-if="showWaterLevelAlert" class="text-red-600 text-sm font-semibold">
      🚨 {{ t('lowTideAlert') }}
    </div>

    <p class="text-md text-gray-700 dark:text-gray-300">
      🌊 {{ t('waterLevel') }}: {{ waterLevelText }}
    </p>
    <p class="text-md text-gray-700 dark:text-gray-300">
      💧 {{ t('waterFlow') }}: {{ waterFlowText }}
    </p>

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
        <span
          v-if="cachedAgeMinutes !== null && cachedAgeMinutes >= 0"
          class="text-xs text-gray-500 ml-2"
        >
          ({{ t('cachedAgo', { minutes: cachedAgeMinutes }) }})
        </span>
      </span>
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
  waterLevelText: string
  waterFlowText: string
  showWaterLevelAlert: boolean
  waterTemperature: number | null
  waterTemperatureLoading: boolean
  waterTemperatureError: string | null
  cachedAgeMinutes: number | null
  chartLabels: string[]
  chartValues: number[]
}>()

</script>
