<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { Chart } from 'chart.js/auto'

const props = defineProps<{
  labels: string[] // history + live
  values: number[] // history + live
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

onMounted(() => {
  if (canvasRef.value) {
    chart = new Chart(canvasRef.value, {
      type: 'line',
      data: {
        labels: [...props.labels],
        datasets: [
          {
            label: 'Water Level (cm)',
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

watch(
  () => [props.labels.length, props.values.length],
  () => {
    if (chart) {
      chart.data.labels = [...props.labels]  // Already history + live combined
      chart.data.datasets[0].data = [...props.values]
      chart.update()
    }
  }
)

</script>

<template>
  <div class="w-full h-64 mt-10">
    <canvas ref="canvasRef"></canvas>
  </div>
</template>
