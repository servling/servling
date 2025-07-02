<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  width?: string | number
  height?: string | number
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '1rem',
  class: '',
})

const formatUnit = (value: string | number) => (typeof value === 'number' ? `${value}px` : value)

const styles = computed(() => ({
  width: formatUnit(props.width),
  height: formatUnit(props.height),
}))
</script>

<template>
  <div
    class="rounded-md bg-zinc-200 animate-pulse dark:bg-zinc-800" :class="props.class"
    :style="styles"
  />
</template>
