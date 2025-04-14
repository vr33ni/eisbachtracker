<template>
    <div class="mt-4 text-left">
      <h3 class="font-semibold text-blue-700 dark:text-blue-300">📊 Prediction for the next hour</h3>
  
      <div v-if="predictionLoading">{{ predictionLoadingMessage }}</div>
      <div v-else-if="predictionError" class="text-red-500">❌ {{ predictionError }}</div>
      <div v-else-if="currentHourPrediction !== null">
        Predicted surfers: <strong>{{ currentHourPrediction }}</strong>
      </div>
      <div v-else class="text-gray-500">
        Not enough data to predict crowd
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  defineProps<{
    predictionLoading: boolean
    predictionError: string | null
    predictionLoadingMessage: string
    currentHourPrediction: number | null
  }>()
  </script>
  