interface UseSseOptions {
  /**
   * Enables automatic reconnection if the connection is lost.
   * @default false
   */
  autoReconnect?: boolean
  /**
   * Delay in milliseconds before attempting to reconnect.
   * @default 3000
   */
  reconnectDelay?: number
  /**
   * Hook that is called just before a reconnection attempt.
   * Can be an async function.
   */
  onReconnect?: () => Promise<void> | void
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
  const { autoReconnect = false, reconnectDelay = 3000, onReconnect } = options

  const event = ref<T | null>(null)
  const status = ref<'Idle' | 'Connecting' | 'Processing' | 'Finished' | 'Error'>('Idle')
  const error = ref<Error | null>(null)

  let abortController: AbortController | undefined
  let reconnectTimeoutId: ReturnType<typeof setTimeout> | undefined

  // Using function declarations to allow hoisting, which prevents
  // 'no-use-before-define' ESLint errors since connect and scheduleReconnect
  // call each other.
  async function connect(): Promise<void> {
    if (import.meta.server || status.value === 'Processing')
      return

    clearTimeout(reconnectTimeoutId)
    reconnectTimeoutId = undefined

    abortController?.abort()
    abortController = new AbortController()
    const { signal } = abortController

    status.value = 'Connecting'
    error.value = null

    let reader: ReadableStreamDefaultReader<Uint8Array> | undefined

    try {
      const stream = await streamFactory()
      if (signal.aborted)
        return

      reader = stream.getReader()
      status.value = 'Processing'
      let leftover = ''
      const decoder = new TextDecoder()

      while (true) {
        const { done, value } = await reader.read()

        if (done || signal.aborted) {
          if (!signal.aborted)
            status.value = 'Finished'
          scheduleReconnect()
          break
        }

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
      if (!signal.aborted) {
        console.error('An error occurred during SSE processing:', e)
        error.value = e instanceof Error ? e : new Error(String(e))
        status.value = 'Error'
        scheduleReconnect()
      }
    }
    finally {
      if (reader && !signal.aborted)
        reader.releaseLock()
    }
  }

  function scheduleReconnect(): void {
    if (!autoReconnect || reconnectTimeoutId)
      return

    status.value = 'Connecting'

    reconnectTimeoutId = setTimeout(() => {
      // This wrapper avoids the 'no-misused-promises' ESLint error.
      const reconnect = async () => {
        reconnectTimeoutId = undefined
        if (onReconnect) {
          try {
            await onReconnect()
          }
          catch (e) {
            console.error('Error in onReconnect hook:', e)
          }
        }
        if (status.value === 'Finished')
          return
        await connect()
      }
      void reconnect()
    }, reconnectDelay)
  }

  function close(): void {
    clearTimeout(reconnectTimeoutId)
    reconnectTimeoutId = undefined
    abortController?.abort()
    status.value = 'Finished'
  }

  return {
    event: readonly(event) as Readonly<Ref<T | null>>,
    status: readonly(status),
    error: readonly(error),
    connect,
    close,
  }
}
