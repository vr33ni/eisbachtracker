<template>
  <form @submit.prevent="onSubmit" class="flex flex-col sm:flex-row items-stretch gap-2 w-full">
    <input
      v-model="countRaw"
      @input="onInputNumeric"
      type="text"
      :placeholder="t('surferCountPlaceholder')"
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
      <span>{{ submitting ? t('submitting') : t('submit') }}</span>
    </button>
  </form>
</template>


<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps<{
  submitting: boolean
  modelValue: string  
}>()

const emit = defineEmits<{
  (e: 'submit', count: number): void
  (e: 'update:modelValue', value: string): void  
}>()

const countRaw = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})
const count = computed(() => {
  return countRaw.value === '' ? null : Number(countRaw.value)
})

const onInputNumeric = (e: Event) => {
  const target = e.target as HTMLInputElement
  target.value = target.value.replace(/[^0-9]/g, '')
  countRaw.value = target.value
}

const onSubmit = () => {
  if (count.value != null && count.value >= 0) {
    emit('submit', count.value)
    countRaw.value = ''
  }
}

</script>
