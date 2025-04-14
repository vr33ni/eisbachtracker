<template>
  <div class="space-y-4">
    <!-- Today's Entries -->
    <h3 class="font-semibold text-blue-700 dark:text-blue-300">📅 Entries Today</h3>

    <div v-if="entriesLoading">{{ entriesLoadingMessage }}</div>
    <div v-else-if="entriesError" class="text-red-500">❌ {{ entriesError }}</div>
    <div v-else>
      <ul v-if="todaysEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1">
        <li v-for="entry in todaysEntries" :key="entry.timestamp">
          {{ new Date(entry.timestamp).toLocaleTimeString() }} — {{ entry.count }} surfers
        </li>
      </ul>
      <p v-else class="text-gray-500 text-sm">No entries yet today</p>
    </div>

    <!-- History Entries Expandable -->
    <ExpandableCard title="🕰️ Entry History">
      <ul v-if="historyEntries.length" class="text-sm text-gray-700 dark:text-gray-300 space-y-1 mt-2">
        <li v-for="entry in historyEntries" :key="entry.timestamp">
          {{ new Date(entry.timestamp).toLocaleDateString() }}
          {{ new Date(entry.timestamp).toLocaleTimeString() }} —
          {{ entry.count }} surfers
        </li>
      </ul>
      <p v-else class="text-gray-500 text-sm mt-2">No historic entries yet</p>
    </ExpandableCard>
  </div>
</template>

<script setup lang="ts">
import ExpandableCard from './ExpandableCard.vue'

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
