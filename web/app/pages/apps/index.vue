<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:squares-four',
    ariaLabel: 'Applications',
    label: 'Applications',
  },
  middleware: ['applications-store'],
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

      <div class="mt-6">
        <ClientOnly>
          <div v-if="applicationStore.fetching" class="flex justify-center">
            <Icon name="ph:spinner-bold" class="text-primary h-6 w-6 animate-spin" />
          </div>

          <div v-else-if="applicationStore.error">
            <UiAlert variant="error">
              {{ applicationStore.error }}
            </UiAlert>
          </div>

          <div v-else-if="applicationStore.applications.length === 0" class="text-center">
            <p class="text-zinc-500 dark:text-zinc-400">
              No applications yet. Create your first one to get started.
            </p>
            <UiButton class="mt-4" left-icon="ph:plus" @click="$router.push('/apps/new')">
              New Application
            </UiButton>
          </div>

          <div v-else class="space-y-4">
            <ApplicationListItem
              v-for="app in applicationStore.applications"
              :key="app.id"
              :application="app"
            />
          </div>

          <template #fallback>
            <div class="mt-6 flex justify-center">
              <Icon name="ph:spinner-bold" class="text-primary h-6 w-6 animate-spin" />
            </div>
          </template>
        </ClientOnly>
      </div>
    </UiContainer>
  </main>
</template>
