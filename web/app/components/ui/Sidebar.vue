<script setup lang="ts">
import type { SidebarItem } from '~/types/sidebar'

defineProps<{
  icon: string
  title: string
  items: SidebarItem[]
  bottomItems: SidebarItem[]
}>()
</script>

<template>
  <aside class="bg-100 bg-zinc-100 flex flex-col h-full w-82 dark:bg-zinc-900">
    <div class="px-6 flex shrink-0 gap-x-4 h-20 items-center">
      <Icon :name="icon" class="text-zinc-900 h-7 w-7 dark:text-zinc-50" />
      <span class="text-xl text-zinc-900 font-semibold dark:text-zinc-50">{{ title }}</span>
    </div>

    <nav class="p-4 flex flex-1 flex-col justify-between">
      <div class="space-y-4">
        <UiSidebarItem v-for="item in items" :key="item.title" :item="item" />
      </div>

      <div class="space-y-4">
        <UiSidebarItem v-for="item in bottomItems" :key="item.title" :item="item" />
      </div>
    </nav>
  </aside>
</template>
