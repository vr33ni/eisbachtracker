<template>
    <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white p-6 rounded shadow-lg w-96">
        <h3 class="text-lg font-semibold mb-4">{{ title }}</h3>
        <slot></slot>
        <button
          class="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          @click="close"
        >
          Close
        </button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, defineProps, defineEmits } from 'vue'
  
  const props = defineProps<{ title: string; visible: boolean }>()
  const emit = defineEmits(['close'])
  
  const close = () => emit('close')
  </script>