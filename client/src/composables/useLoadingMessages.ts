import { ref } from 'vue'

export function useLoadingMessages(messages: string[], intervalMs = 1000) {
  const loadingMessage = ref(messages[0])
  let index = 0
  let timer: number | null = null

  const startRotating = () => {
    stopRotating() // prevent multiple timers
    timer = window.setInterval(() => {
      index = (index + 1) % messages.length
      loadingMessage.value = messages[index]
    }, intervalMs)
  }

  const stopRotating = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
      index = 0
      loadingMessage.value = messages[0]
    }
  }

  return {
    loadingMessage,
    startRotating,
    stopRotating
  }
}
