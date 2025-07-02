<template>
  <aside class="bg-100 bg-zinc-100 flex flex-col h-full w-82 dark:bg-zinc-900">
    <div class="px-6 flex shrink-0 gap-x-4 h-20 items-center">
      <slot name="title" />
    </div>

    <nav class="p-4 flex flex-1 flex-col justify-between">
      <div class="space-y-4">
        <slot />
      </div>

      <div class="space-y-4">
        <slot name="bottom" />
      </div>
    </nav>
  </aside>
</template>
