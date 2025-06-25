import { client } from '#hey-api/client.gen'
import { ofetch } from 'ofetch'

export const useApplicationStore = defineStore('applications', () => {
  const applications = ref<Application[]>([])
  const error = ref<string | null>(null)
  const fetching = ref<boolean>(true)
  const { session } = useUserSession()

  const {
    event,
    connect,
    close: closeSSEListener,
  } = useSSE<ApplicationStatusChangedMessage>(async () => ofetch(`${client.getConfig().baseURL}/applications/events`, {
    responseType: 'stream',
    headers: {
      Authorization: `Bearer ${session.value?.accessToken}`,
    },
  }))

  async function fetchApplications() {
    try {
      applications.value = await getApplications({
        composable: '$fetch',
      })
      fetching.value = false
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

  async function initializeSSEListener() {
    watch(event, (statusChangedEvent) => {
      if (statusChangedEvent == null) {
        return
      }
      const index = applications.value.findIndex(app => app.id === statusChangedEvent.id)
      if (index !== -1 && applications.value[index]) {
        applications.value[index].status = statusChangedEvent.status
        applications.value[index].error = statusChangedEvent.error
      }
    })

    await connect()
  }

  return {
    applications,
    fetchApplications,
    initializeSSEListener,
    closeSSEListener,
    error,
    fetching,
  }
})
