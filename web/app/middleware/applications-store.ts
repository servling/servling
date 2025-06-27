export default defineNuxtRouteMiddleware(async () => {
  const appStore = useApplicationStore()
  const { loggedIn } = useUserSession()

  if (loggedIn.value) {
    appStore.fetchApplications()
    appStore.initializeSSEListener()
  }

  onUnmounted(() => {
    appStore.closeSSEListener()
  })
})
