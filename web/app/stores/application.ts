export const useApplicationStore = defineStore('applications', () => {
  const applications = ref<Application[]>([])
  const eventSource = ref<EventSource | null>(null)
  const error = ref<string | null>(null)
  const { session } = useUserSession()

  async function fetchApplications() {
    try {
      applications.value = await getApplications({
        composable: '$fetch',
      })
    }
    catch (err) {
      if (err instanceof Error) {
        error.value = err.message
      }
      else {
        throw err
      }
    }
  }

  function initializeSseListener() {
    if (import.meta.server) {
      return
    }

    if (eventSource.value) {
      return
    }

    eventSource.value = new EventSource(`http://localhost:9999/applications/events?token=${session.value?.accessToken}`)

    eventSource.value.addEventListener('application.status-changed', (event) => {
      // eslint-disable-next-line ts/no-unsafe-argument
      const changedStatus = JSON.parse(event.data) as ApplicationStatusChangedMessage

      const index = applications.value.findIndex(app => app.id === changedStatus.id)
      if (index !== -1 && applications.value[index]) {
        applications.value[index].status = changedStatus.status
        applications.value[index].error = changedStatus.error
      }
    })

    eventSource.value.onerror = (error) => {
      console.error('SSE Error:', error)
      closeSseListener()
    }
  }

  function closeSseListener() {
    if (eventSource.value) {
      eventSource.value.close()
      eventSource.value = null
    }
  }

  return {
    applications,
    fetchApplications,
    initializeSseListener,
    closeSseListener,
    error,
  }
})
