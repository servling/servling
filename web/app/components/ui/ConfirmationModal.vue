<script setup lang="ts">
const props = withDefaults(defineProps<{
  isOpen?: boolean
  title?: string
  message?: string
  confirmButtonText?: string
  cancelButtonText?: string
  confirmButtonIntent?: string
}>(), {
  isOpen: false,
  title: 'Confirm Action',
  message: 'Are you sure you want to proceed?',
  confirmButtonText: 'Confirm',
  cancelButtonText: 'Cancel',
  confirmButtonIntent: 'danger',
})

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

function onConfirm() {
  emit('confirm')
}

function onCancel() {
  emit('cancel')
}
</script>

<template>
  <DialogRoot :open="isOpen" modal>
    <DialogPortal>
      <DialogOverlay class="bg-black/50 inset-0 fixed" />
      <DialogContent class="p-6 rounded-lg bg-white max-w-md w-full shadow-lg transform left-1/2 top-1/2 fixed dark:bg-zinc-900 -translate-x-1/2 -translate-y-1/2">
        <DialogTitle class="text-lg text-zinc-900 font-medium dark:text-white">
          {{ title }}
        </DialogTitle>
        <DialogDescription class="text-sm text-zinc-500 mt-2 dark:text-zinc-400">
          {{ message }}
        </DialogDescription>

        <div class="mt-6 flex justify-end space-x-3">
          <DialogClose as-child>
            <UiButton size="sm" intent="secondary" @click="onCancel">
              {{ cancelButtonText }}
            </UiButton>
          </DialogClose>
          <DialogClose as-child>
            <UiButton size="sm" :intent="confirmButtonIntent" @click="onConfirm">
              {{ confirmButtonText }}
            </UiButton>
          </DialogClose>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
