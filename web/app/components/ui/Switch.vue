<script setup lang="ts">
import type { VariantProps } from 'cva'
import { cva } from 'cva'

const props = withDefaults(defineProps<{
  modelValue?: boolean | object | null
  disabled?: boolean
  label?: string
  description?: string
  id?: string
}>(), {
  modelValue: false,
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const switchStyles = cva(
  {
    base: 'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
    variants: {
      checked: {
        true: 'bg-blue-600',
        false: 'bg-zinc-200 dark:bg-zinc-700',
      },
    },
    defaultVariants: {
      checked: false,
    },
  },
)

const toggleStyles = cva(
  {
    base: 'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
    variants: {
      checked: {
        true: 'translate-x-5',
        false: 'translate-x-0',
      },
    },
    defaultVariants: {
      checked: false,
    },
  },
)

type SwitchProps = VariantProps<typeof switchStyles>

const switchId = props.id ?? useId()

function toggle() {
  if (!props.disabled) {
    emit('update:modelValue', !props.modelValue)
  }
}
</script>

<template>
  <div class="flex items-center">
    <button
      :id="switchId"
      type="button"
      :disabled="disabled"
      :class="switchStyles({ checked: modelValue })"
      :aria-checked="modelValue"
      role="switch"
      @click="toggle"
    >
      <span
        aria-hidden="true"
        :class="toggleStyles({ checked: modelValue })"
      />
    </button>
    <div v-if="label || description" class="ml-3">
      <label
        v-if="label"
        :for="switchId"
        class="text-sm text-zinc-900 font-medium dark:text-zinc-200"
      >
        {{ label }}
      </label>
      <p
        v-if="description"
        class="text-sm text-zinc-500 dark:text-zinc-400"
      >
        {{ description }}
      </p>
    </div>
  </div>
</template>
