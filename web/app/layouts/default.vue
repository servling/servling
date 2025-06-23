<template>
  <div class="flex h-full w-full">
    <LayoutSidebar />
    <div class="flex flex-col w-full overflow-auto">
      <LayoutHeader />
      <slot />
    </div>
  </div>
</template>
