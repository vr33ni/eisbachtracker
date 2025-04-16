<template>
  <div class="space-y-4">
    <!-- Today's Entries -->
    <h3 class="font-semibold text-blue-700 dark:text-blue-300">📅 {{ t('entriesToday') }}</h3>

    <div v-if="entriesLoading">{{ t('loadingEntries') }}</div>
    <div v-else-if="entriesError" class="text-red-500">❌ {{ entriesError }}</div>
    <div v-else>
      <ul v-if="todaysEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1">
        <li v-for="entry in todaysEntries" :key="entry.timestamp">
          {{ new Date(entry.timestamp).toLocaleTimeString() }} — {{ entry.count }} {{ t('surfers') }}
        </li>
      </ul>
      <p v-else class="text-gray-500 text-sm">{{ t('noEntriesToday') }}</p>
    </div>

    <!-- History Entries Expandable -->
    <ExpandableCard :title="t('entryHistoryTitle')">
      <ul v-if="historyEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1 mt-2">
        <li v-for="entry in historyEntries" :key="entry.timestamp">
          {{ new Date(entry.timestamp).toLocaleDateString() }}
          {{ new Date(entry.timestamp).toLocaleTimeString() }} —
          {{ entry.count }} {{ t('surfers') }}
        </li>
      </ul>
      <p v-else class="text-gray-500 text-sm mt-2">{{ t('noHistory') }}</p>
    </ExpandableCard>
  </div>
</template>

<script setup lang="ts">
import ExpandableCard from './ExpandableCard.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps<{
  todaysEntries: {
    timestamp: string
    count: number
  }[]
  historyEntries: {
    timestamp: string
    count: number
  }[]
  entriesLoading: boolean
  entriesError: string | null
  entriesLoadingMessage: string
}>()
</script>
