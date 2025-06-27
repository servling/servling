export default defineNuxtRouteMiddleware(async (to) => {
  const { loggedIn } = useUserSession()

  if (to.meta.auth === 'ignore') {
    return
  }

  if (to.meta.auth === false && loggedIn.value) {
    return navigateTo('/')
  }

  if ((to.meta.auth === true || to.meta.auth === undefined) && !loggedIn.value) {
    return navigateTo('/login')
  }
})
