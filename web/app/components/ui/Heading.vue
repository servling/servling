<script setup lang="ts">
import { cva } from 'cva'

const props = withDefaults(defineProps<{
  title: string
  description?: string
  level?: 1 | 2 | 3 | 4
}>(), {
  level: 1,
})

const headingTag = computed(() => `h${props.level}`)

const headingStyles = cva(
  {
    base: 'font-bold leading-tight tracking-tight text-zinc-900 dark:text-white',
    variants: {
      level: {
        1: 'text-3xl',
        2: 'text-2xl',
        3: 'text-xl',
        4: 'text-lg',
      },
    },
    defaultVariants: {
      level: 1,
    },
  },
)
</script>

<template>
  <div class="flex flex-col gap-4 items-start justify-between sm:flex-row sm:gap-0 sm:items-center">
    <div class="flex flex-col">
      <component :is="headingTag" :class="headingStyles({ level })">
        {{ title }}
      </component>

      <p v-if="description" class="text-base text-zinc-500 mt-2 dark:text-zinc-400">
        {{ description }}
      </p>
    </div>

    <div class="flex shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>
