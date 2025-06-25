// @ts-check
import antfu from '@antfu/eslint-config'
import nuxt from './.nuxt/eslint.config.mjs'

export default antfu(
  {
    unocss: true,
    typescript: {
      tsconfigPath: 'tsconfig.json',
    },
  },
  {
    rules: {
      'no-undef': 'off',
      'ts/no-unsafe-assignment': 'off',
      'ts/no-unsafe-member-access': 'off',
      'ts/no-unsafe-call': 'off',
    },
  },
)
  .append(nuxt())
