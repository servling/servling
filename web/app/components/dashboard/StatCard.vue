<script setup lang="ts">
import { cva } from 'cva'

withDefaults(defineProps<{
  title: string
  value: any
  subtitle?: string
  icon?: string
  color?: string
  showProgressBar?: boolean
  total?: string
  percentage?: number
}>(), {
  subtitle: '',
  color: 'primary',
  showProgressBar: false,
  total: '',
  percentage: 0,
})

const valueStyles = cva(
  {
    base: 'text-3xl font-bold',
    variants: {
      color: {
        primary: 'text-blue-500',
        success: 'text-green-500',
        warning: 'text-amber-500',
        danger: 'text-red-500',
      },
    },
    defaultVariants: {
      color: 'primary',
    },
  },
)

const iconStyles = cva(
  {
    base: 'h-6 w-6',
    variants: {
      color: {
        primary: 'text-blue-500',
        success: 'text-green-500',
        warning: 'text-amber-500',
        danger: 'text-red-500',
      },
    },
    defaultVariants: {
      color: 'primary',
    },
  },
)

const progressBarStyles = cva(
  {
    base: 'h-2 mt-2 rounded-full bg-zinc-700 overflow-hidden',
    variants: {
      color: {
        primary: '',
        success: '',
        warning: '',
        danger: '',
      },
    },
    defaultVariants: {
      color: 'primary',
    },
  },
)

const progressFillStyles = cva(
  {
    base: 'h-full rounded-full',
    variants: {
      color: {
        primary: 'bg-blue-500',
        success: 'bg-green-500',
        warning: 'bg-amber-500',
        danger: 'bg-red-500',
      },
    },
    defaultVariants: {
      color: 'primary',
    },
  },
)
</script>

<template>
  <UiCard>
    <div class="flex flex-col h-full">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-sm text-zinc-500 font-medium dark:text-zinc-400">
          {{ title }}
        </h3>
        <div v-if="icon">
          <Icon :name="icon" :class="iconStyles({ color })" />
        </div>
      </div>

      <div class="flex-grow">
        <div class="flex items-baseline justify-between">
          <div :class="valueStyles({ color })">
            {{ value }}
          </div>
          <div v-if="total" class="text-sm text-zinc-500 dark:text-zinc-400">
            of {{ total }}
          </div>
        </div>

        <div v-if="subtitle" class="text-sm text-zinc-500 mt-1 dark:text-zinc-400">
          {{ subtitle }}
        </div>
      </div>

      <div v-if="showProgressBar" class="mt-auto">
        <div :class="progressBarStyles({ color })">
          <div
            :class="progressFillStyles({ color })"
            :style="{ width: `${percentage}%` }"
          />
        </div>
        <div class="text-xs text-zinc-500 mt-1 text-right dark:text-zinc-400">
          {{ percentage }}% used
        </div>
      </div>
    </div>
  </UiCard>
</template>
