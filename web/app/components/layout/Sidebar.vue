<script setup lang="ts">
import type { SidebarItem } from '~/types/sidebar'
import { sidebarBottomItems, sidebarItems } from '~/contants'

const color = useColorMode()

function toggleDark() {
  color.preference = color.value === 'dark' ? 'light' : 'dark'
}

const bottomItems = computed<SidebarItem[]>(() => ([
  ...sidebarBottomItems,
  {
    title: color.value === 'dark' ? 'Dark Mode' : 'Light Mode',
    icon: color.value === 'dark' ? 'ph:moon' : 'ph:sun',
    action: toggleDark,
    elevated: true,
  },
]))
</script>

<template>
  <UiSidebar icon="ph:package-bold" title="Servling" :items="sidebarItems" :bottom-items />
</template>
