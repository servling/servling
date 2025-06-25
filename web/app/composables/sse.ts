import type { Ref } from 'vue'
import { readonly, ref } from 'vue'

interface UseSseReturn<T> {
  event: Readonly<Ref<T | null>>
  status: Readonly<Ref<'Idle' | 'Connecting' | 'Processing' | 'Finished' | 'Error'>>
  error: Readonly<Ref<Error | null>>
  connect: () => Promise<void>
  close: () => Promise<void>
}

export function useSSE<T = unknown>(
  streamFactory: () => Promise<ReadableStream<Uint8Array>>,
): UseSseReturn<T> {
  const event = ref<T | null>(null)
  const status = ref<'Idle' | 'Connecting' | 'Processing' | 'Finished' | 'Error'>('Idle')
  const error = ref<Error | null>(null)

  let reader: ReadableStreamDefaultReader<Uint8Array> | undefined
  let leftover = '' // Buffer for incomplete lines

  const connect = async (): Promise<void> => {
    if (import.meta.server) {
      return
    }
    if (['Processing', 'Connecting'].includes(status.value)) {
      console.warn('SSE connection is already active or connecting.')
      return
    }

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
          break
        }

        // Prepend any leftover data from the previous chunk
        const textChunk = leftover + decoder.decode(value, { stream: true })
        const lines = textChunk.split(/\r?\n/)

        // The last line might be incomplete, so we save it for the next chunk
        leftover = lines.pop() ?? ''

        for (const line of lines) {
          if (line.trim().startsWith('data:')) {
            const jsonData = line.substring(line.indexOf(':') + 1).trim()
            if (jsonData) {
              try {
                event.value = JSON.parse(jsonData) as T
              }
              catch (e) {
                const err = e instanceof Error ? e : new Error('SSE stream data is not valid.')
                console.error('SSE JSON parsing error:', err)
                error.value = err
              }
            }
          }
        }
      }
    }
    catch (e: unknown) {
      if (e instanceof Error) {
        if (e.name === 'AbortError' || e.message.includes('cancelled')) {
          status.value = 'Finished'
        }
        else {
          console.error('An error occurred during SSE processing:', e)
          error.value = e
          status.value = 'Error'
        }
      }
    }
    finally {
      if (reader) {
        reader.releaseLock()
        reader = undefined
      }
    }
  }

  const close = async (): Promise<void> => {
    if (reader) {
      await reader.cancel('Stream closed by user.')
    }
    else {
      console.warn('No active SSE stream to close.')
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
