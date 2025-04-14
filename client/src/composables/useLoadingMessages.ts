import { ref } from 'vue'

export function useLoadingMessages(customMessages: string[]) {
  const loadingMessage = ref(customMessages[0])
  let interval: ReturnType<typeof setInterval> | null = null

  const startRotating = () => {
    let i = 0
    interval = setInterval(() => {
      loadingMessage.value = customMessages[i % customMessages.length]
      i++
    }, 2500)
  }

  const stopRotating = () => {
    if (interval) clearInterval(interval)
  }

  return {
    loadingMessage,
    startRotating,
    stopRotating,
  }
}
