<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="pb-5 flex flex-col gap-4 items-start justify-between sm:flex-row sm:gap-0 sm:items-center">
    <div class="flex flex-col">
      <h1 class="text-2xl text-zinc-900 leading-7 font-bold sm:text-3xl dark:text-white sm:tracking-tight sm:truncate">
        {{ title }}
      </h1>
      <p v-if="description" class="text-sm text-zinc-500 mt-1 dark:text-zinc-400">
        {{ description }}
      </p>
    </div>
    <div class="flex shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>
