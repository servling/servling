import type { Ref } from 'vue'

interface UseSseOptions {
  autoReconnect?: boolean
  reconnectDelay?: number
}

interface UseSseReturn<T> {
  event: Readonly<Ref<T | null>>
  status: Readonly<Ref<'Idle' | 'Connecting' | 'Processing' | 'Finished' | 'Error'>>
  error: Readonly<Ref<Error | null>>
  connect: () => Promise<void>
  close: () => void
}

export function useSSE<T = unknown>(
  streamFactory: () => Promise<ReadableStream<Uint8Array>>,
  options: UseSseOptions = {},
): UseSseReturn<T> {
  const { autoReconnect = false, reconnectDelay = 3000 } = options

  const event = ref<T | null>(null)
  const status = ref<'Idle' | 'Connecting' | 'Processing' | 'Finished' | 'Error'>('Idle')
  const error = ref<Error | null>(null)

  let reader: ReadableStreamDefaultReader<Uint8Array> | undefined
  let leftover = ''
  let isClosedByUser = false
  let reconnectTimeoutId: ReturnType<typeof setTimeout> | undefined

  let handleReconnect: () => void

  const connect = async (): Promise<void> => {
    if (import.meta.server)
      return
    if (['Processing', 'Connecting'].includes(status.value)) {
      return
    }

    isClosedByUser = false
    clearTimeout(reconnectTimeoutId)
    status.value = 'Connecting'
    error.value = null
    event.value = null

    try {
      const stream = await streamFactory()
      reader = stream.getReader()
      const decoder = new TextDecoder()
      status.value = 'Processing'

      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          status.value = 'Finished'
          handleReconnect()
          break
        }
        // ... (rest of the processing logic is unchanged)
        const textChunk = leftover + decoder.decode(value, { stream: true })
        const lines = textChunk.split(/\r?\n/)
        leftover = lines.pop() ?? ''
        for (const line of lines) {
          if (line.trim().startsWith('data:')) {
            const jsonData = line.substring(line.indexOf(':') + 1).trim()
            if (jsonData) {
              try {
                event.value = JSON.parse(jsonData) as T
              }
              catch (e) {
                error.value = e instanceof Error ? e : new Error('SSE stream data is not valid.')
                console.error('SSE JSON parsing error:', error.value)
              }
            }
          }
        }
      }
    }
    catch (e: unknown) {
      // Improved Catch Block: If the error is not a user-initiated abort, try to reconnect.
      if (e instanceof Error && (e.name === 'AbortError' || e.message.includes('cancelled'))) {
        status.value = 'Finished'
      }
      else {
        console.error('An error occurred during SSE processing:', e)
        error.value = e instanceof Error ? e : new Error(String(e))
        status.value = 'Error'
        handleReconnect()
      }
    }
    finally {
      if (reader) {
        reader.releaseLock()
        reader = undefined
      }
    }
  }

  handleReconnect = () => {
    if (!autoReconnect || isClosedByUser)
      return

    status.value = 'Connecting'
    reconnectTimeoutId = setTimeout(() => {
      if (!isClosedByUser) {
        void connect()
      }
    }, reconnectDelay)
  }

  const close = (): void => {
    isClosedByUser = true
    clearTimeout(reconnectTimeoutId)

    if (reader) {
      reader.cancel('Stream closed by user.').catch(() => {})
      reader = undefined
      status.value = 'Finished'
    }
  }

  return {
    event: readonly(event) as Readonly<Ref<T | null>>,
    status: readonly(status),
    error: readonly(error),
    connect,
    close,
  }
}
