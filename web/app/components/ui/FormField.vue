<script setup lang="ts">
import { useField } from 'vee-validate'

const props = defineProps<{
  name: string
}>()

// The `name` is returned in a function because we want to make sure it stays reactive
// If the name changes you want `useField` to be able to pick it up
const { value, errorMessage, handleChange } = useField(() => props.name)

// Ensure we always have a string value for UiInput to avoid undefined warnings
const inputValue = computed(() => value.value === undefined ? '' : value.value)
</script>

<template>
  <slot
    :value="inputValue"
    :error-message="errorMessage"
    :handle-change="handleChange"
  />
</template>
