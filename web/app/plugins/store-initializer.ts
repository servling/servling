import { useApplicationStore } from '~/stores/application'

export default defineNuxtPlugin(async () => {
  const appStore = useApplicationStore()

  await appStore.fetchApplications()
  appStore.initializeSseListener()
})
