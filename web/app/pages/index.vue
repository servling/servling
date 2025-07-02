<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:house',
    ariaLabel: 'Dashboard',
    label: 'Dashboard',
  },
  middleware: ['applications-store'],
})

const applicationStore = useApplicationStore()

const applicationStats = computed(() => {
  const count = applicationStore.applications.length
  const running = applicationStore.applications.filter(app => app.status === 'running').length

  return {
    count,
    running,
    percentage: count > 0 ? Math.round((running / count) * 100) : 0,
  }
})

const stats = {
  nodes: {
    count: 1,
    online: 1,
  },
  memory: {
    used: '64.2 GB',
    total: '128 GB',
    percentage: 50.2,
  },
  storage: {
    used: '1.2 TB',
    total: '2 TB',
    percentage: 60,
  },
}
</script>

<template>
  <main>
    <UiContainer>
      <div class="gap-8 grid grid-cols-1 lg:grid-cols-4 md:grid-cols-2">
        <ClientOnly>
          <DashboardStatCard
            title="Applications"
            :value="applicationStats.count"
            :subtitle="`${applicationStats.running} running`"
            icon="ph:app-window"
            color="primary"
            :show-progress-bar="true"
            :percentage="applicationStats.percentage"
          />
          <template #fallback>
            <UiCard>
              <div class="space-y-2">
                <UiSkeleton :height="32" width="75%" class="rounded-md" />
                <UiSkeleton :height="48" width="10%" class="rounded-md" />
                <UiSkeleton :height="24" width="100%" class="rounded-md" />
              </div>
            </UiCard>
          </template>
        </ClientOnly>

        <DashboardStatCard
          title="Nodes"
          :value="stats.nodes.count"
          :subtitle="`${stats.nodes.online} online`"
          icon="ph:circles-four"
          color="success"
        />
        <DashboardStatCard
          title="Memory Usage"
          :value="stats.memory.used"
          :total="stats.memory.total"
          :percentage="stats.memory.percentage"
          icon="ph:chart-pie"
          color="warning"
          :show-progress-bar="true"
        />
        <DashboardStatCard
          title="Storage Usage"
          :value="stats.storage.used"
          :total="stats.storage.total"
          :percentage="stats.storage.percentage"
          icon="ph:hard-drive"
          color="danger"
          :show-progress-bar="true"
        />
      </div>

      <div class="mt-8">
        <h2 class="text-xl text-black font-bold mb-4 dark:text-white">
          Application Status
        </h2>

        <ClientOnly>
          <div v-if="applicationStats.count === 0" class="py-6 text-center rounded-lg bg-zinc-100 dark:bg-zinc-900">
            <p class="text-zinc-600 dark:text-zinc-400">
              No applications yet. Create your first one to get started.
            </p>
            <UiButton class="mt-4" left-icon="ph:plus" @click="$router.push('/apps/new')">
              New Application
            </UiButton>
          </div>

          <div v-else class="space-y-4">
            <div v-for="app in applicationStore.applications" :key="app.id" class="p-4 rounded-lg bg-zinc-100 flex items-center justify-between dark:bg-zinc-900">
              <div class="flex items-center space-x-4">
                <div class="p-2 rounded-md bg-zinc-200 flex-shrink-0 dark:bg-zinc-800">
                  <Icon name="ph:squares-four" class="text-zinc-600 h-5 w-5 dark:text-zinc-400" />
                </div>
                <div>
                  <div class="flex gap-2 items-center">
                    <h3 class="text-black font-medium dark:text-white">
                      {{ app.name }}
                    </h3>
                    <ApplicationStatusBadge :status="app.status" />
                  </div>
                  <p class="text-sm text-zinc-600 dark:text-zinc-400">
                    {{ app.description }}
                  </p>
                </div>
              </div>
              <UiButton size="sm" intent="secondary" left-icon="ph:eye" @click="$router.push(`/apps/${app.id}`)">
                View
              </UiButton>
            </div>
          </div>

          <template #fallback>
            <div class="flex justify-center">
              <Icon name="ph:spinner-bold" class="text-primary h-6 w-6 animate-spin" />
            </div>
          </template>
        </ClientOnly>
      </div>
    </UiContainer>
  </main>
</template>
