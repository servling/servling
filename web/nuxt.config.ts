// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/icon',
    '@nuxt/image',
    '@nuxt/eslint',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    '@nuxtjs/color-mode',
    '@unocss/nuxt',
    'pinia-plugin-persistedstate/nuxt',
    '@nuxtjs/seo',
    'reka-ui/nuxt',
    '@bicou/nuxt-urql',
    'nuxt-auth-utils',
  ],
  devtools: { enabled: true },

  colorMode: {
    classSuffix: '',
  },
  future: {
    compatibilityVersion: 4,
  },

  experimental: {
    renderJsonPayloads: true,
    typedPages: true,
    buildCache: true,
  },
  compatibilityDate: '2025-05-15',

  nitro: {
    esbuild: {
      options: {
        target: 'esnext',
      },
    },
    prerender: {
      crawlLinks: false,
    },
  },

  eslint: {
    config: {
      standalone: false,
      nuxt: {
        sortConfigKeys: true,
      },
    },
  },

  pinia: {
    storesDirs: ['./app/stores'],
  },

  piniaPluginPersistedstate: {
    storage: 'cookies',
  },

  urql: {
    endpoint: 'http://localhost:8080/graphql',
  },
})
