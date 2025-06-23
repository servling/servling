<template>
  <div class="text-black bg-white h-full w-full transition-500 transition-colors dark:text-white dark:bg-zinc-950">
    <NuxtRouteAnnouncer />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
html,
body,
#__nuxt {
  height: 100%;
  margin: 0;
  padding: 0;
}

html.dark {
  color-scheme: dark;
}
</style>
