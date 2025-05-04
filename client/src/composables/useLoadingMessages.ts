import { ref, type Ref, watchEffect } from 'vue'

export function useLoadingMessages(messages: Ref<string[] | undefined>, intervalMs = 1000) {
  const loadingMessage = ref('')
  let index = 0
  let timer: ReturnType<typeof setInterval> | null = null

  const getMessages = () => messages.value ?? []

  const updateMessage = () => {
    const msgs = getMessages()
    if (msgs.length > 0) {
      loadingMessage.value = msgs[index % msgs.length]
    }
  }

  const startRotating = () => {
    stopRotating()
    updateMessage()
    timer = setInterval(() => {
      const msgs = getMessages()
      if (msgs.length === 0) return
      index = (index + 1) % msgs.length
      loadingMessage.value = msgs[index]
    }, intervalMs)
  }

  const stopRotating = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    index = 0
    const msgs = getMessages()
    loadingMessage.value = msgs[0] ?? ''
  }

  // Live update when language or message set changes
  watchEffect(() => {
    const msgs = getMessages()
    if (msgs.length > 0) {
      // console.log('🌍 Locale switched or messages changed:', msgs)
      updateMessage()
    }
  })

  // Set something immediately
  updateMessage()

  return {
    loadingMessage,
    startRotating,
    stopRotating,
  }
}
