<script setup lang="ts">
import type { VariantProps } from 'cva'
import { cva } from 'cva'

const props = withDefaults(defineProps<{
  intent?: ButtonProps['intent']
  size?: ButtonProps['size']
  loading?: boolean
  as?: 'button' | 'a' | 'NuxtLink'
  leftIcon?: string
  rightIcon?: string
}>(), {
  as: 'button',
  size: 'md',
})

// 1. CVA RECIPE: Updated sizes to include font-size for better scaling.
const buttonStyles = cva(
  {
    base: 'inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none',
    variants: {
      intent: {
        primary: 'bg-blue-600 text-white hover:bg-blue-700',
        secondary: 'bg-zinc-100 text-zinc-900 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-50 dark:hover:bg-zinc-700',
        danger: 'bg-red-600 text-white hover:bg-red-700',
        ghost: 'hover:bg-zinc-100 hover:text-zinc-900 dark:hover:bg-zinc-800 dark:hover:text-zinc-50',
        link: 'text-blue-600 underline-offset-4 hover:underline dark:text-blue-400',
      },
      size: {
        sm: 'h-9 px-3 text-sm',
        md: 'h-10 px-4 py-2 text-base',
        lg: 'h-11 px-8 text-lg',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      intent: 'primary',
      size: 'md',
    },
  },
)

type ButtonProps = VariantProps<typeof buttonStyles>

const iconSizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'h-4 w-4'
    case 'lg':
      return 'h-6 w-6'
    case 'md':
    case 'icon':
    default:
      return 'h-5 w-5'
  }
})
</script>

<template>
  <component
    :is="as"
    :class="buttonStyles({ intent, size })"
    :disabled="loading"
  >
    <div v-if="loading" class="flex items-center justify-center">
      <Icon name="ph:spinner-bold" :class="iconSizeClass" class="animate-spin" />
    </div>

    <div v-else class="flex gap-x-2 items-center justify-center">
      <Icon v-if="leftIcon" :name="leftIcon" :class="iconSizeClass" />
      <slot v-else name="left-icon" />

      <slot />

      <Icon v-if="rightIcon" :name="rightIcon" :class="iconSizeClass" />
      <slot v-else name="right-icon" />
    </div>
  </component>
</template>
