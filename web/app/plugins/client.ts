import { client } from '#hey-api/client.gen'

export default defineNuxtPlugin(() => {
  const { session, clear } = useUserSession()

  function isTokenExpired(expiresAt: string): boolean {
    return new Date() > new Date(expiresAt)
  }

  async function refreshTokens(): Promise<void> {
    if (!session.value)
      return
    const response = await refresh({
      composable: '$fetch',
      body: {
        refreshToken: session.value.refreshToken,
      },
    })

    if (response == null) {
      return
    }

    await patchSession({
      accessToken: response.accessToken,
      accessTokenExpiresAt: response.accessTokenExpiresAt,
    })
  }
  client.setConfig({
    async auth() {
      if (!session.value) {
        return undefined
      }

      if (isTokenExpired(session.value.accessTokenExpiresAt)) {
        if (isTokenExpired(session.value.refreshTokenExpiresAt)) {
          await clear()
          return undefined
        }

        await refreshTokens()
      }

      return session.value.accessToken ?? undefined
    },
    baseURL: 'http://localhost:9999',
    headers: {
      Authorization: '',
    },
  })
})
