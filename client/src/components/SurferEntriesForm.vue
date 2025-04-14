<template>
  <form @submit.prevent="onSubmit" class="flex items-center gap-2">
    <input
      v-model="countRaw"
      @input="onInputNumeric"
      type="text"
      placeholder="Number of surfers"
      inputmode="numeric"
      pattern="[0-9]*"
      class="px-3 py-2 rounded border dark:bg-gray-800 dark:border-gray-600"
    />

    <button
      type="submit"
      :disabled="submitting || count === null"
      class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 justify-center min-w-[120px]"
    >
      <svg
        v-if="submitting"
        class="animate-spin h-5 w-5 text-white"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        />
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
        />
      </svg>
      <span>{{ submitting ? 'Submitting...' : 'Submit' }}</span>
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  submitting: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', count: number): void
}>()

const countRaw = ref('')
const count = computed(() => Number(countRaw.value))

const onInputNumeric = (e: Event) => {
  const target = e.target as HTMLInputElement
  target.value = target.value.replace(/[^0-9]/g, '')
  countRaw.value = target.value
}

const onSubmit = () => {
  if (!isNaN(count.value) && count.value >= 0) {
    emit('submit', count.value)
    countRaw.value = ''
  }
}
</script>
