// composables/useGlobalLoadingOverlay.ts
import { ref } from 'vue'
import { useLoadingMessages } from './useLoadingMessages'

export function useGlobalLoadingOverlay(messages: string[]) {
  const isRefreshing = ref(false)

  const {
    loadingMessage,
    startRotating,
    stopRotating
  } = useLoadingMessages(ref(messages))

  const start = () => {
    isRefreshing.value = true
    startRotating()
  }

  const stop = () => {
    stopRotating()
    isRefreshing.value = false
  }

  const cancel = () => stop() // alias for UI clarity

  return {
    isRefreshing,
    rotatingMessage: loadingMessage,
    startGlobalRefresh: start,
    stopGlobalRefresh: stop,
    cancelGlobalRefresh: cancel
  }
}
