declare module '#auth-utils' {
  interface User {
    id: string
    name: string
    createdAt: string
    updatedAt: string
  }

  interface UserSession {
    accessToken: string
    accessTokenExpiresAt: string
    refreshToken: string
    refreshTokenExpiresAt: string
  }

  interface SecureSessionData {
  }
}

export {}
