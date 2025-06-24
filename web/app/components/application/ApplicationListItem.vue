<script setup lang="ts">
const props = defineProps<{
  application: Application
}>()

const applicationStore = useApplicationStore()

const isDeleteModalOpen = ref(false)

function openDeleteModal() {
  isDeleteModalOpen.value = true
}

function closeDeleteModal() {
  isDeleteModalOpen.value = false
}

async function handleDeleteConfirm() {
  const deletedApp = await deleteApplication({
    composable: '$fetch',
    path: {
      id: props.application.id,
    },
  })
  applicationStore.applications = applicationStore.applications.filter(app => app.id !== deletedApp.id)
  closeDeleteModal()
}

async function handleStartApplication() {
  await startApplication({
    composable: '$fetch',
    path: { id: props.application.id },
  })
}

async function handleStopApplication() {
  await stopApplication({
    composable: '$fetch',
    path: { id: props.application.id },
  })
}

const canStart = computed(() => {
  return ['stopped', 'error'].includes(props.application.status)
})

const canStop = computed(() => {
  return ['running', 'starting', 'deploying'].includes(props.application.status)
})
</script>

<template>
  <div class="p-4 rounded-lg bg-zinc-100 flex items-center justify-between dark:bg-zinc-900">
    <div class="flex items-center space-x-4">
      <div class="p-2 rounded-md bg-zinc-200 flex-shrink-0 dark:bg-zinc-800">
        <Icon name="ph:squares-four" class="text-zinc-600 h-5 w-5 dark:text-zinc-400" />
      </div>
      <div>
        <h3 class="text-black font-medium dark:text-white">
          {{ application.name }}
        </h3>
        <div class="mt-1 flex gap-2 items-center">
          <ApplicationStatusBadge :status="application.status" />
          <p class="text-sm text-zinc-600 dark:text-zinc-400">
            {{ application.description }}
          </p>
        </div>
        <div class="text-xs text-zinc-600 mt-1 flex items-center">
          <span>Created {{ new Date(application.createdAt).toLocaleDateString() }}</span>
          <span v-if="application.services.length > 0" class="ml-3 flex items-center">
            <Icon name="ph:cube" class="mr-1 h-3 w-3" />
            {{ application.services.length }} service{{ application.services.length > 1 ? 's' : '' }}
          </span>
        </div>
      </div>
    </div>
    <div class="flex space-x-2">
      <UiButton size="sm" intent="secondary" left-icon="ph:eye" @click="$router.push(`/apps/${application.id}`)">
        View
      </UiButton>
      <UiButton
        v-if="canStart"
        size="sm"
        intent="primary"
        left-icon="ph:play"
        :loading="application.status === 'starting'"
        @click="handleStartApplication"
      >
        Start
      </UiButton>
      <UiButton
        v-if="canStop"
        size="sm"
        intent="danger"
        left-icon="ph:stop"
        :loading="application.status === 'stopping'"
        @click="handleStopApplication"
      >
        Stop
      </UiButton>
      <UiButton size="sm" intent="danger" left-icon="ph:trash" @click="openDeleteModal">
        Delete
      </UiButton>
    </div>

    <UiConfirmationModal
      :is-open="isDeleteModalOpen"
      title="Delete Application"
      :message="`Are you sure you want to delete '${application.name}'? This action cannot be undone.`"
      confirm-button-text="Delete"
      confirm-button-intent="danger"
      @confirm="handleDeleteConfirm"
      @cancel="closeDeleteModal"
    />
  </div>
</template>
