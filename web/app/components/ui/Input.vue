<script setup lang="ts">
import { cva } from 'cva'

const props = withDefaults(defineProps<{
  modelValue?: string | number | object | null
  label?: string
  hint?: string
  error?: string
  id?: string
  leftIcon?: string
  rightIcon?: string
}>(), {
  modelValue: '',
})

defineEmits<{
  'update:modelValue': [value: string]
}>()

const inputStyles = cva(
  {

    base:
          'block w-full rounded-md border-0 bg-white/5 py-2 px-3 shadow-sm ring-1 ring-inset placeholder:text-zinc-400 focus:ring-2 focus:ring-inset sm:text-sm sm:leading-6',
    variants: {
      state: {
        default: 'text-zinc-900 ring-zinc-300 focus:ring-blue-600 dark:text-white dark:ring-zinc-700',
        error: 'text-red-500 ring-red-500 focus:ring-red-600 placeholder:text-red-300',
      },
      hasLeftIcon: {
        true: 'pl-10',
      },
      hasRightIcon: {
        true: 'pr-10',
      },
    },
    defaultVariants: {
      state: 'default',
    },
  },
)

const inputId = props.id ?? useId()
</script>

<template>
  <div>
    <label
      v-if="label"
      :for="inputId"
      class="text-sm text-zinc-900 leading-6 font-medium block dark:text-zinc-200"
    >
      {{ label }}
    </label>

    <div class="relative" :class="[label ? 'mt-2' : '']">
      <div v-if="leftIcon" class="pl-3 flex pointer-events-none items-center inset-y-0 left-0 absolute">
        <Icon :name="leftIcon" class="text-zinc-400 h-5 w-5" />
      </div>

      <input
        :id="inputId"
        :value="modelValue"
        v-bind="$attrs"
        :class="inputStyles({
          state: error ? 'error' : 'default',
          hasLeftIcon: !!leftIcon,
          hasRightIcon: !!rightIcon,
        })"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      >

      <div v-if="rightIcon" class="pr-3 flex pointer-events-none items-center inset-y-0 right-0 absolute">
        <Icon :name="rightIcon" class="h-5 w-5" :class="[error ? 'text-red-500' : 'text-zinc-400']" />
      </div>
    </div>

    <p v-if="error" :id="`${inputId}-error`" class="text-sm text-red-600 mt-2">
      {{ error }}
    </p>
    <p v-else-if="hint" :id="`${inputId}-hint`" class="text-sm text-zinc-500 mt-2">
      {{ hint }}
    </p>
  </div>
</template>
