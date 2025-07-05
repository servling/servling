<script setup lang="ts">
const color = useColorMode()

function toggleDark() {
  color.preference = color.value === 'dark' ? 'light' : 'dark'
}

const darkModeTitle = computed(() => color.value === 'dark' ? 'Dark Mode' : 'Light Mode')
const darkModeIcon = computed(() => color.value === 'dark' ? 'ph:moon' : 'ph:sun')
</script>

<template>
  <UiSidebar>
    <template #title>
      <Icon name="ph:package-bold" class="text-zinc-900 h-7 w-7 dark:text-zinc-50" />
      <span class="text-xl text-zinc-900 font-semibold dark:text-zinc-50">Servling</span>
    </template>
    <UiSidebarItem title="Dashboard" icon="ph:house" href="/" />
    <UiSidebarItem title="Applications" icon="ph:squares-four" href="/apps" />
    <UiSidebarItem title="App Store" icon="ph:storefront" href="/apps/store" />
    <UiSidebarItem title="Domains" icon="ph:globe-hemisphere-west" href="/domains" />
    <template #bottom>
      <ClientOnly>
        <UiSidebarItem :title="darkModeTitle" :icon="darkModeIcon" :action="toggleDark" elevated />
        <template #fallback>
          <div class="px-4 py-3 flex gap-x-4 w-full items-center">
            <UiSkeleton :width="20" :height="20" class="rounded-full flex-shrink-0" />
            <div class="w-full">
              <UiSkeleton height="1rem" width="75%" />
            </div>
          </div>
        </template>
      </ClientOnly>
      <UiSidebarItem title="Instance Settings" icon="ph:gear" href="/settings" />
    </template>
  </UiSidebar>
</template>
