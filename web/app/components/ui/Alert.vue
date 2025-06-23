<script setup lang="ts">
import type { VariantProps } from 'cva'
import { cva } from 'cva'

const props = withDefaults(defineProps<{
  variant?: AlertProps['variant']
  title?: string
  icon?: string
  dismissible?: boolean
}>(), {
  variant: 'default',
  dismissible: false,
})

const emit = defineEmits(['dismiss'])

const alertStyles = cva(
  {
    base: 'flex w-full items-center justify-between rounded-md border p-4 text-sm',
    variants: {
      variant: {
        default: 'bg-zinc-50 text-zinc-900 border-zinc-200 dark:bg-zinc-800/50 dark:text-zinc-50 dark:border-zinc-700',
        error: 'bg-red-50 text-red-900 border-red-200 dark:bg-red-900/20 dark:text-red-300 dark:border-red-800',
        success: 'bg-green-50 text-green-900 border-green-200 dark:bg-green-900/20 dark:text-green-300 dark:border-green-800',
        warning: 'bg-yellow-50 text-yellow-900 border-yellow-200 dark:bg-yellow-900/20 dark:text-yellow-300 dark:border-yellow-800',
        info: 'bg-blue-50 text-blue-900 border-blue-200 dark:bg-blue-900/20 dark:text-blue-300 dark:border-blue-800',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  },
)

type AlertProps = VariantProps<typeof alertStyles>

const iconMap = {
  default: 'ph:info-bold',
  error: 'ph:warning-bold',
  success: 'ph:check-circle-bold',
  warning: 'ph:warning-bold',
  info: 'ph:info-bold',
}

const iconColor = computed(() => {
  switch (props.variant) {
    case 'error': return 'text-red-500 dark:text-red-400'
    case 'success': return 'text-green-500 dark:text-green-400'
    case 'warning': return 'text-yellow-500 dark:text-yellow-400'
    case 'info': return 'text-blue-500 dark:text-blue-400'
    default: return 'text-zinc-500 dark:text-zinc-400'
  }
})

const displayIcon = computed(() => props.icon || iconMap[props.variant])
</script>

<template>
  <div :class="alertStyles({ variant })">
    <div class="flex gap-3 items-center">
      <Icon v-if="displayIcon" :name="displayIcon" class="h-5 w-5" :class="iconColor" />

      <div>
        <div v-if="title" class="font-medium">
          {{ title }}
        </div>
        <div class="text-sm">
          <slot />
        </div>
      </div>
    </div>

    <button
      v-if="dismissible"
      class="ml-auto rounded-md flex h-6 w-6 items-center justify-center hover:bg-zinc-200/50 dark:hover:bg-zinc-700/50"
      @click="emit('dismiss')"
    >
      <Icon name="ph:x-bold" class="text-zinc-500 h-4 w-4 dark:text-zinc-400" />
    </button>
  </div>
</template>
