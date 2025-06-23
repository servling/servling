import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: {
    path: '../schema/openapi.json',
    validate_EXPERIMENTAL: true,
  },
  output: {
    path: 'app/client',
  },
  plugins: ['@hey-api/client-nuxt'],
})
