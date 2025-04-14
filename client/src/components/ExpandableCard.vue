<template>
  <div class="rounded-lg shadow bg-white dark:bg-gray-800">
    <!-- Header -->
    <div
      @click="toggle"
      class="flex items-center justify-between px-6 py-4 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
    >
      <h3 class="text-lg font-medium text-gray-800 dark:text-gray-200">
        {{ title }}
      </h3>

      <svg
        class="h-5 w-5 text-gray-500 transform transition-transform duration-300"
        :class="{ 'rotate-180': isOpen }"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
      </svg>
    </div>

    <!-- Expandable Content -->
    <div
      v-show="isOpen"
      class="px-6 pb-4 pt-2 border-t border-gray-200 dark:border-gray-700 transition-all"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  title: string
  initiallyOpen?: boolean
}>()

const isOpen = ref(props.initiallyOpen ?? false)
const toggle = () => {
  isOpen.value = !isOpen.value
}
</script>
