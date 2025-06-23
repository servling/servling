<template>
  <div class="px-4 py-12 bg-white flex flex-col min-h-screen items-center justify-center dark:bg-zinc-950">
    <div class="max-w-md w-full">
      <div class="flex gap-2 items-center justify-center">
        <Icon name="ph:package-bold" class="text-blue-600 h-10 w-10" /> <span class="text-xl font-semibold">Servling</span>
      </div>
      <slot />
    </div>
  </div>
</template>
