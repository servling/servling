import type { UserSession } from '#auth-utils'

export async function updateSession(newSession: Omit<UserSession, 'id'>) {
  const { session } = useUserSession()
  const resp = await $fetch<UserSession>('/api/session', {
    method: 'POST',
    body: newSession,
  })
  session.value = resp
  return resp
}

export async function patchSession(sessionDelta: Partial<Omit<UserSession, 'id'>>) {
  const { session } = useUserSession()
  const resp = await $fetch<UserSession>('/api/session', {
    method: 'PATCH',
    body: sessionDelta,
  })
  session.value = resp
  return resp
}
