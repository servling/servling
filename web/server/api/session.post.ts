import type { UserSession } from '#auth-utils'

export default eventHandler(async (event) => {
  const req = await readBody<UserSession>(event)

  return setUserSession(event, req)
})
