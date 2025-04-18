<script setup lang="ts">
import { onMounted, watch, ref, computed } from 'vue'
import { Chart } from 'chart.js/auto'
import { useI18n } from 'vue-i18n'

const { locale, t } = useI18n()

const props = defineProps<{
  labels: string[] // pre-processed based on mode in parent
  values: number[]
  mode: 'hourly' | 'daily'
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null


const emit = defineEmits<{
  (e: 'update:mode', value: 'hourly' | 'daily'): void
}>()


// Watch for i18n changes
watch(locale, () => {
  if (chart) {
    chart.data.datasets[0].label = t('chart.waterLevelLabel')
    chart.update()
  }
})

// Watch for chart data changes
watch(
  () => [props.labels, props.values],
  () => {
    if (chart) {
      chart.data.labels = [...props.labels]
      chart.data.datasets[0].data = [...props.values]
      chart.update()
    }
  },
  { deep: true }
)

onMounted(() => {
  if (canvasRef.value) {
    chart = new Chart(canvasRef.value, {
      type: 'line',
      data: {
        labels: [...props.labels],
        datasets: [
          {
            label: t('chart.waterLevelLabel'),
            data: [...props.values],
            fill: true,
            backgroundColor: 'rgba(59, 130, 246, 0.2)',
            borderColor: 'rgba(59, 130, 246, 1)',
            tension: 0.3,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: false,
          },
        },
      },
    })
  }
})
</script>

<template>
  <div class="w-full mt-8 space-y-4">
    <!-- Mode Switch -->
    <div class="flex justify-end">
      <select :value="props.mode" @change="emit('update:mode', ($event.target as HTMLSelectElement)?.value as 'hourly' | 'daily')"
        class="text-sm rounded border p-1 bg-white dark:bg-gray-700 dark:text-white">
        <option value="hourly">🕐 {{ $t('chart.hourly') }}</option>
        <option value="daily">📆 {{ $t('chart.daily') }}</option>
      </select>

    </div>

    <!-- Chart -->
    <div class="w-full h-64">
      <canvas ref="canvasRef"></canvas>
    </div>
  </div>
</template>
