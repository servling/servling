import type { UserSession } from '#auth-utils'

export default eventHandler(async (event) => {
  const oldSession = await getUserSession(event)
  const req = await readBody<Partial<UserSession>>(event)

  return setUserSession(event, { ...oldSession, ...req })
})
