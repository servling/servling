<script setup lang="ts">
import { cva } from 'cva'

const props = defineProps<{
  status: Application['status']
}>()

const statusStyles = cva(
  {
    base: 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
    variants: {
      status: {
        running: 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-300',
        stopped: 'bg-zinc-100 text-zinc-800 dark:bg-zinc-900/20 dark:text-zinc-300',
        starting: 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-300',
        stopping: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-300',
        deploying: 'bg-purple-100 text-purple-800 dark:bg-purple-900/20 dark:text-purple-300',
        error: 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-300',
      },
    },
    defaultVariants: {
      status: 'stopped',
    },
  },
)

const statusText = computed(() => {
  const statusMap = {
    running: 'Running',
    stopped: 'Stopped',
    starting: 'Starting',
    stopping: 'Stopping',
    deploying: 'Deploying',
    error: 'Error',
  }
  return statusMap[props.status]
})

const statusIcon = computed(() => {
  const iconMap = {
    running: 'ph:play-circle',
    stopped: 'ph:stop-circle',
    starting: 'ph:spinner',
    stopping: 'ph:spinner',
    deploying: 'ph:spinner',
    error: 'ph:warning-circle',
  }
  return iconMap[props.status]
})

const iconClass = computed(() => {
  return ['starting', 'stopping', 'deploying'].includes(props.status) ? 'animate-spin' : ''
})
</script>

<template>
  <span :class="statusStyles({ status })">
    <Icon :name="statusIcon" class="mr-1 h-3 w-3" :class="iconClass" />
    {{ statusText }}
  </span>
</template>
