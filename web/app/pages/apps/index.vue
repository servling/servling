<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:squares-four',
    ariaLabel: 'Applications',
    label: 'Applications',
  },
})

const applicationStore = useApplicationStore()
</script>

<template>
  <main>
    <UiContainer constrained>
      <UiPageHeader
        title="Applications"
        description="Manage your deployed applications."
      >
        <template #actions>
          <UiButton left-icon="ph:plus" @click="$router.push('/apps/new')">
            New Application
          </UiButton>
        </template>
      </UiPageHeader>

      <div v-if="applicationStore.fetching" class="mt-6 flex justify-center">
        <Icon name="ph:spinner-bold" class="text-primary h-6 w-6 animate-spin" />
      </div>

      <div v-else-if="applicationStore.error" class="mt-6">
        <UiAlert variant="error">
          {{ applicationStore.error }}
        </UiAlert>
      </div>

      <div v-else-if="applicationStore.applications.length === 0" class="mt-6 text-center">
        <p class="text-zinc-500 dark:text-zinc-400">
          No applications yet. Create your first one to get started.
        </p>
        <UiButton class="mt-4" left-icon="ph:plus" @click="$router.push('/apps/new')">
          New Application
        </UiButton>
      </div>

      <div v-else class="mt-6 space-y-4">
        <ApplicationListItem
          v-for="app in applicationStore.applications"
          :key="app.id"
          :application="app"
        />
      </div>
    </UiContainer>
  </main>
</template>
