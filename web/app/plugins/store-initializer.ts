import { useApplicationStore } from '~/stores/application'

export default defineNuxtPlugin(() => {
  const appStore = useApplicationStore()

  appStore.fetchApplications()
  appStore.initializeSSEListener()
})
