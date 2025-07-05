<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:squares-four',
    ariaLabel: 'Application',
    label: 'Application',
  },
  middleware: ['applications-store'],
})

const route = useRoute()

const applicationStore = useApplicationStore()

const application = computed(() => {
  return applicationStore.applications.find(application => application.id === route.params.id)
})

async function handleStartApplication() {
  await startApplication({
    composable: '$fetch',
    path: {
      id: route.params.id as string,
    },
  })
}

async function handleStopApplication() {
  await stopApplication({
    composable: '$fetch',
    path: {
      id: route.params.id as string,
    },
  })
}

const canStart = computed(() => {
  if (!application.value)
    return false
  return ['stopped', 'error'].includes(application.value.status)
})

const canStop = computed(() => {
  if (!application.value)
    return false
  return ['running', 'starting', 'deploying'].includes(application.value.status)
})

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString()
}

function formatPorts(ports: Record<string, any>) {
  if (!ports || Object.keys(ports).length === 0) {
    return []
  }

  return Object.entries(ports).map(([containerPort, hostPort]) => ({
    container: containerPort,
    host: hostPort,
  }))
}

function formatEnvironment(env: Record<string, any>) {
  if (!env || Object.keys(env).length === 0) {
    return []
  }

  return Object.entries(env).map(([key, value]) => ({
    key,
    value: String(value),
  }))
}

function formatLabels(labels: Record<string, any>) {
  if (!labels || Object.keys(labels).length === 0) {
    return []
  }

  return Object.entries(labels).map(([key, value]) => ({
    key,
    value: String(value),
  }))
}

function isEmpty(obj: Record<string, any> | undefined | null) {
  return !obj || Object.keys(obj).length === 0
}
</script>

<template>
  <main>
    <UiContainer constrained>
      <div v-if="applicationStore.error" class="py-8">
        <UiAlert variant="error">
          {{ applicationStore.error }}
        </UiAlert>
      </div>

      <div v-else-if="!application" class="py-8">
        <UiAlert variant="error">
          Application not found
        </UiAlert>
      </div>

      <template v-else>
        <!-- Application Header -->
        <div class="flex items-center justify-between">
          <UiPageHeader
            :title="application.name"
            :description="application.description"
          />
          <div class="flex gap-4 items-center">
            <ApplicationStatusBadge :status="application.status" />
            <UiButton
              v-if="canStart"
              intent="primary"
              left-icon="ph:play-bold"
              :loading="application.status === 'starting'"
              @click="handleStartApplication"
            >
              Start
            </UiButton>
            <UiButton
              v-if="canStop"
              intent="danger"
              left-icon="ph:stop-bold"
              :loading="application.status === 'stopping'"
              @click="handleStopApplication"
            >
              Stop
            </UiButton>
          </div>
        </div>

        <UiAlert
          v-if="application.error"
          variant="error"
          class="mt-6"
        >
          {{ application.error }}
        </UiAlert>

        <div class="mt-8">
          <h2 class="text-xl text-black font-bold mb-4 dark:text-white">
            Services
          </h2>

          <div v-if="application.services.length === 0" class="py-6 text-center rounded-lg bg-zinc-100 dark:bg-zinc-900">
            <p class="text-zinc-600 dark:text-zinc-400">
              No services defined for this application.
            </p>
          </div>

          <div v-else class="space-y-4">
            <UiCard v-for="service in application.services" :key="service.id">
              <div class="space-y-4">
                <div class="flex items-center justify-between">
                  <h3 class="text-lg font-semibold">
                    {{ service.name }}
                  </h3>
                  <div class="text-sm text-zinc-500">
                    {{ service.image }}
                  </div>
                </div>

                <!-- Service Details -->
                <div class="gap-4 grid grid-cols-1 md:grid-cols-3">
                  <!-- Ports -->
                  <div class="space-y-2">
                    <h4 class="text-sm text-zinc-700 font-medium dark:text-zinc-300">
                      Ports
                    </h4>
                    <div v-if="isEmpty(service.ports)" class="text-xs text-zinc-500 italic">
                      No ports configured
                    </div>
                    <div v-else class="p-3 rounded-md bg-zinc-800 max-h-40 overflow-auto">
                      <div
                        v-for="(port, index) in formatPorts(service.ports)" :key="index"
                        class="py-1 border-b border-zinc-700 flex items-center justify-between last:border-0"
                      >
                        <div class="text-xs text-emerald-400">
                          {{ port.container }}
                        </div>
                        <div class="text-xs text-zinc-400 flex items-center">
                          <Icon name="ph:arrow-right" class="mx-2" />
                          <span class="text-blue-400">{{ port.host }}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Labels -->
                  <div class="space-y-2">
                    <h4 class="text-sm text-zinc-700 font-medium dark:text-zinc-300">
                      Labels
                    </h4>
                    <div v-if="isEmpty(service.labels)" class="text-xs text-zinc-500 italic">
                      No labels defined
                    </div>
                    <div v-else class="p-3 rounded-md bg-zinc-800 max-h-40 overflow-auto">
                      <div
                        v-for="(label, index) in formatLabels(service.labels)" :key="index"
                        class="mb-1 last:mb-0"
                      >
                        <div class="text-xs text-violet-400 font-medium">
                          {{ label.key }}
                        </div>
                        <div class="text-xs text-zinc-300 ml-1 pl-2 border-l-2 border-zinc-700">
                          {{ label.value }}
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Environment -->
                  <div class="space-y-2">
                    <h4 class="text-sm text-zinc-700 font-medium dark:text-zinc-300">
                      Environment
                    </h4>
                    <div v-if="isEmpty(service.environment)" class="text-xs text-zinc-500 italic">
                      No environment variables set
                    </div>
                    <div v-else class="p-3 rounded-md bg-zinc-800 max-h-40 overflow-auto">
                      <div
                        v-for="(env, index) in formatEnvironment(service.environment)" :key="index"
                        class="mb-1 last:mb-0"
                      >
                        <div class="text-xs text-amber-400 font-medium">
                          {{ env.key }}
                        </div>
                        <div class="text-xs text-zinc-300 ml-1 pl-2 border-l-2 border-zinc-700">
                          {{ env.value }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="text-xs text-zinc-500">
                  Created: {{ formatDate(service.createdAt) }} | Updated: {{ formatDate(service.updatedAt) }}
                </div>
              </div>
            </UiCard>
          </div>
        </div>

        <!-- Application Metadata -->
        <div class="text-sm text-zinc-500 mt-8">
          <div>ID: {{ application.id }}</div>
          <div>Created: {{ formatDate(application.createdAt) }}</div>
          <div>Last Updated: {{ formatDate(application.updatedAt) }}</div>
        </div>
      </template>
    </UiContainer>
  </main>
</template>
