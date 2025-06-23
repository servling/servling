<script setup lang="ts">
withDefaults(defineProps<{
  constrained?: boolean
}>(), {
  constrained: false,
})
</script>

<template>
  <div
    class="py-6 w-full sm:py-8" :class="[
      constrained
        ? 'mx-auto max-w-7xl px-4 sm:px-6 lg:px-8'
        : 'px-4 sm:px-6 lg:px-8',
    ]"
  >
    <slot />
  </div>
</template>
