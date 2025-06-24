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
    'nuxt-auth-utils',
    '@hey-api/nuxt',
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

  heyApi: {
    config: {
      input: {
        path: '../schema/openapi.json',
        validate_EXPERIMENTAL: true,
      },
      plugins: [
        '@hey-api/schemas',
        {
          name: '@hey-api/sdk',
          transformer: true,
          validator: false,
          auth: true,
        },
        {
          enums: 'javascript',
          name: '@hey-api/typescript',
        },
        '@hey-api/transformers',
        'zod',
      ],
    },
  },

  pinia: {
    storesDirs: ['./app/stores'],
  },

  piniaPluginPersistedstate: {
    storage: 'cookies',
  },
})
