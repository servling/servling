<script setup lang="ts">
defineSlots<{
  default: () => any
  header: () => any
  footer: () => any
}>()
const slots = useSlots()
const hasHeader = computed(() => !!slots.header)
const hasFooter = computed(() => !!slots.footer)
</script>

<template>
  <div class="rounded-lg bg-white ring-1 ring-zinc-200 shadow overflow-hidden dark:bg-zinc-900 dark:ring-zinc-800">
    <div
      v-if="hasHeader"
      class="px-4 py-5 border-b border-zinc-200 sm:px-6 dark:border-zinc-800"
    >
      <slot name="header" />
    </div>

    <div class="px-4 py-5 sm:p-6">
      <slot />
    </div>

    <div
      v-if="hasFooter"
      class="px-4 py-5 border-t border-zinc-200 sm:px-6 dark:border-zinc-800"
    >
      <slot name="footer" />
    </div>
  </div>
</template>
