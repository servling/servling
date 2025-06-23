// @ts-check
import antfu from '@antfu/eslint-config'
import nuxt from './.nuxt/eslint.config.mjs'

export default antfu(
  {
    unocss: true,
    rules: {
      'no-undef': 'off',
    },
    typescript: {
      tsconfigPath: 'tsconfig.json',
    },
  },
)
  .append(nuxt())
