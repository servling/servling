<script setup lang="ts">
// Use `cva` from the 'cva' package (previously class-variance-authority)
import { cva } from 'cva'

interface Props {
  title: string
  icon?: string
  href?: string
  elevated?: boolean
  action?: () => void
}

defineProps<Props>()

// This CVA recipe is built according to the `cva@beta` documentation you provided.
const sidebarItem = cva(
  // `base`: Applies to all variants, as shown in the docs.
  {
    base: 'flex items-center w-full px-4 py-3 rounded-lg transition-colors duration-200 ease-in-out gap-x-4 font-medium text-base',
    // `variants`: The primary variants for our component.
    variants: {
      variant: {
        item: '', // Item-specific styles are handled in compoundVariants for clarity
        title: 'px-3 pt-6 pb-2 text-xs font-semibold uppercase tracking-wider text-zinc-400 dark:text-zinc-500',
      },
      active: {
        true: '',
        false: '',
      },
      elevated: {
        true: '',
        false: '',
      },
    },

    // `compoundVariants`: The key to fixing our style conflicts, as shown in your docs.
    // This allows us to define classes for specific *combinations* of variants.
    compoundVariants: [
      // STATE 1: Active (Highest precedence)
      {
        variant: 'item',
        active: true,
        class: 'bg-blue-600 text-white hover:bg-blue-500',
      },
      // STATE 2: Inactive AND Elevated
      {
        variant: 'item',
        active: false,
        elevated: true,
        class: 'bg-zinc-100 text-zinc-700 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-200 dark:hover:bg-zinc-700',
      },
      // STATE 3: Inactive AND Default (This now correctly includes the text color)
      {
        variant: 'item',
        active: false,
        elevated: false,
        class: 'text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800/50',
      },
    ],

    // `defaultVariants`: The default state of an item.
    defaultVariants: {
      variant: 'item',
      active: false,
      elevated: false,
    },
  },
)
</script>

<template>
  <NuxtLink
    v-if="href"
    v-slot="{ isActive, href, navigate }"
    :to="href"
    custom
  >
    <a :href="href" :class="sidebarItem({ elevated, active: isActive })" @click="navigate">
      <Icon v-if="icon" :name="icon" class="flex-shrink-0 h-5 w-5" />
      <span>{{ title }}</span>
    </a>
  </NuxtLink>

  <button
    v-else-if="action"
    :class="sidebarItem({ elevated, active: false })"
    @click="action"
  >
    <Icon v-if="icon" :name="icon" class="flex-shrink-0 h-5 w-5" />
    <span>{{ title }}</span>
  </button>

  <span v-else :class="sidebarItem({ variant: 'title' })">
    {{ title }}
  </span>
</template>
