<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

definePageMeta({
  breadcrumb: {
    icon: 'ph:plus',
    ariaLabel: 'New Application',
    label: 'New Application',
  },
})

const router = useRouter()
const applicationStore = useApplicationStore()

// Form state
const services = ref([
  {
    name: '',
    image: '',
    ports: [{ key: '', value: '' }],
    labels: [{ key: '', value: '' }],
    environment: [{ key: '', value: '' }],
  },
])

// Validation schema
const validationSchema = z.object({
  name: z.string().min(1, 'Application name is required'),
  description: z.string().min(1, 'Description is required'),
  start: z.boolean().default(false),
  services: z.array(
    z.object({
      name: z.string().min(1, 'Service name is required'),
      image: z.string().min(1, 'Docker image is required'),
      entrypoint: z.string(),
      ports: z.array(
        z.object({
          key: z.string(),
          value: z.string(),
        }),
      ),
      labels: z.array(
        z.object({
          key: z.string(),
          value: z.string(),
        }),
      ),
      environment: z.array(
        z.object({
          key: z.string(),
          value: z.string(),
        }),
      ),
    }),
  ),
})

type CreateApplicationSchema = z.infer<typeof validationSchema>

const { handleSubmit, isSubmitting: _isSubmitting } = useForm({
  validationSchema: toTypedSchema(validationSchema),
  initialValues: {
    name: '',
    description: '',
    start: false,
    services: [
      {
        name: '',
        image: '',
        entrypoint: '',
        ports: [{ key: '', value: '' }],
        labels: [{ key: '', value: '' }],
        environment: [{ key: '', value: '' }],
      },
    ],
  },
})

// Service management
function _addService() {
  services.value.push({
    name: '',
    image: '',
    ports: [{ key: '', value: '' }],
    labels: [{ key: '', value: '' }],
    environment: [{ key: '', value: '' }],
  })
}

function _removeService(index: number) {
  if (services.value.length > 1) {
    services.value.splice(index, 1)
  }
}

// Key-value pair management
function _addKeyValuePair(service: CreateApplicationRequest['services'][number] & Record<string, any[]>, field: string) {
  service[field]?.push({ key: '', value: '' })
}

function _removeKeyValuePair(service: CreateApplicationRequest['services'][number] & Record<string, any[]>, field: string, index: number) {
  if (service[field]?.length ?? 0 > 1) {
    service[field]?.splice(index, 1)
  }
}

function prepareSubmitData(values: CreateApplicationSchema): CreateApplicationRequest {
  return {
    name: values.name,
    description: values.description,
    start: values.start,
    services: values.services.map((service) => {
      const ports: Record<string, string> = {}
      const labels: Record<string, string> = {}
      const environment: Record<string, string> = {}

      service.ports.forEach((pair) => {
        if (pair.key && pair.value) {
          ports[pair.key] = pair.value
        }
      })

      service.labels.forEach((pair) => {
        if (pair.key && pair.value) {
          labels[pair.key] = pair.value
        }
      })

      service.environment.forEach((pair) => {
        if (pair.key && pair.value) {
          environment[pair.key] = pair.value
        }
      })

      return {
        name: service.name,
        image: service.image,
        entrypoint: service.entrypoint,
        ports,
        labels,
        environment,
      }
    }),
  }
}

const apiErrors = ref<string[]>([])

const onSubmit = handleSubmit(async (values) => {
  apiErrors.value = []

  try {
    const input = prepareSubmitData(values)
    const result = await createApplication({
      composable: '$fetch',
      body: input,
    })
    applicationStore.applications.push(result)

    if (result) {
      await router.push('/apps')
    }
  }
  catch (err) {
    if (err instanceof Error) {
      apiErrors.value.push(err.message)
    }
    else {
      apiErrors.value.push('An unknown error occurred')
    }
  }
})
</script>

<template>
  <main>
    <UiContainer constrained>
      <UiPageHeader
        title="New Application"
        description="Create a new application with services."
      />

      <form class="mt-8" @submit="onSubmit">
        <UiAlert
          v-if="apiErrors.length > 0"
          variant="error"
          class="mb-6"
        >
          <ul class="pl-5 list-disc">
            <li v-for="(error, index) in apiErrors" :key="index">
              {{ error }}
            </li>
          </ul>
        </UiAlert>

        <UiCard>
          <template #header>
            <UiPageHeader
              :level="2"
              title="Application Details"
              description="Provide basic information about your application."
            />
          </template>

          <div class="space-y-6">
            <UiFormField v-slot="{ value, errorMessage, handleChange }" name="name">
              <UiInput
                id="name"
                :model-value="value"
                :error="errorMessage"
                label="Application Name"
                type="text"
                placeholder="My Application"
                left-icon="ph:app-window"
                @update:model-value="handleChange"
              />
            </UiFormField>

            <UiFormField v-slot="{ value, errorMessage, handleChange }" name="description">
              <UiInput
                id="description"
                :model-value="value"
                :error="errorMessage"
                label="Description"
                type="text"
                placeholder="Describe your application"
                left-icon="ph:text-align-left"
                @update:model-value="handleChange"
              />
            </UiFormField>

            <UiFormField v-slot="{ value, handleChange }" name="start">
              <UiSwitch
                id="start"
                :model-value="value"
                label="Auto-start Application"
                description="Automatically start the application after creation"
                @update:model-value="handleChange"
              />
            </UiFormField>
          </div>
        </UiCard>

        <div class="mt-8 space-y-6">
          <UiCard v-for="(service, serviceIndex) in services" :key="serviceIndex">
            <template #header>
              <div class="flex items-center justify-between">
                <UiPageHeader
                  :level="2"
                  :title="`Service ${serviceIndex + 1}`"
                  description="Configure a service for your application."
                />
                <UiButton
                  v-if="services.length > 1"
                  intent="danger"
                  size="sm"
                  type="button"
                  @click="_removeService(serviceIndex)"
                >
                  Remove
                </UiButton>
              </div>
            </template>

            <div class="space-y-6">
              <div class="gap-4 grid grid-cols-2 w-full">
                <UiFormField
                  v-slot="{ value, errorMessage, handleChange }"
                  :name="`services[${serviceIndex}].name`"
                >
                  <UiInput
                    :id="`service-${serviceIndex}-name`"
                    :model-value="value"
                    :error="errorMessage"
                    label="Service Name"
                    type="text"
                    placeholder="api"
                    left-icon="ph:cube"
                    @update:model-value="handleChange"
                  />
                </UiFormField>

                <UiFormField
                  v-slot="{ value, errorMessage, handleChange }"
                  :name="`services[${serviceIndex}].image`"
                >
                  <UiInput
                    :id="`service-${serviceIndex}-image`"
                    :model-value="value"
                    :error="errorMessage"
                    label="Docker Image"
                    type="text"
                    placeholder="nginx:latest"
                    left-icon="bxl:docker"
                    @update:model-value="handleChange"
                  />
                </UiFormField>
              </div>

              <!-- Ports -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <label class="text-sm text-zinc-900 leading-6 font-medium block dark:text-zinc-200">
                    Ports
                  </label>
                  <UiButton
                    size="sm"
                    intent="secondary"
                    type="button"
                    @click="_addKeyValuePair(service, 'ports')"
                  >
                    Add Port
                  </UiButton>
                </div>

                <div
                  v-for="(port, portIndex) in service.ports"
                  :key="portIndex"
                  class="flex gap-2 items-center"
                >
                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].ports[${portIndex}].key`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-port-${portIndex}-key`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Container Port (e.g., 80)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].ports[${portIndex}].value`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-port-${portIndex}-value`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Host Port (e.g., 8080)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiButton
                    v-if="service.ports.length > 1"
                    size="sm"
                    intent="danger"
                    type="button"
                    @click="_removeKeyValuePair(service, 'ports', portIndex)"
                  >
                    Remove
                  </UiButton>
                </div>
              </div>

              <!-- Labels -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <label class="text-sm text-zinc-900 leading-6 font-medium block dark:text-zinc-200">
                    Labels
                  </label>
                  <UiButton
                    size="sm"
                    intent="secondary"
                    type="button"
                    @click="_addKeyValuePair(service, 'labels')"
                  >
                    Add Label
                  </UiButton>
                </div>

                <div
                  v-for="(label, labelIndex) in service.labels"
                  :key="labelIndex"
                  class="flex gap-2 items-center"
                >
                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].labels[${labelIndex}].key`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-label-${labelIndex}-key`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Key (e.g., app)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].labels[${labelIndex}].value`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-label-${labelIndex}-value`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Value (e.g., myapp)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiButton
                    v-if="service.labels.length > 1"
                    size="sm"
                    intent="danger"
                    type="button"
                    @click="_removeKeyValuePair(service, 'labels', labelIndex)"
                  >
                    Remove
                  </UiButton>
                </div>
              </div>

              <!-- Environment Variables -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <label class="text-sm text-zinc-900 leading-6 font-medium block dark:text-zinc-200">
                    Environment Variables
                  </label>
                  <UiButton
                    size="sm"
                    intent="secondary"
                    type="button"
                    @click="_addKeyValuePair(service, 'environment')"
                  >
                    Add Environment Variable
                  </UiButton>
                </div>

                <div
                  v-for="(env, envIndex) in service.environment"
                  :key="envIndex"
                  class="flex gap-2 items-center"
                >
                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].environment[${envIndex}].key`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-env-${envIndex}-key`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Key (e.g., NODE_ENV)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiFormField
                    v-slot="{ value, errorMessage, handleChange }"
                    :name="`services[${serviceIndex}].environment[${envIndex}].value`"
                  >
                    <UiInput
                      :id="`service-${serviceIndex}-env-${envIndex}-value`"
                      :model-value="value"
                      :error="errorMessage"
                      placeholder="Value (e.g., production)"
                      @update:model-value="handleChange"
                    />
                  </UiFormField>

                  <UiButton
                    v-if="service.environment.length > 1"
                    size="sm"
                    intent="danger"
                    type="button"
                    @click="_removeKeyValuePair(service, 'environment', envIndex)"
                  >
                    Remove
                  </UiButton>
                </div>
              </div>
            </div>
          </UiCard>

          <div class="flex justify-center">
            <UiButton
              intent="secondary"
              left-icon="ph:plus"
              type="button"
              @click="_addService"
            >
              Add Another Service
            </UiButton>
          </div>
        </div>

        <div class="mt-8 flex justify-end space-x-4">
          <UiButton
            intent="secondary"
            type="button"
            @click="router.push('/apps')"
          >
            Cancel
          </UiButton>
          <UiButton
            :loading="_isSubmitting"
            type="submit"
          >
            Create Application
          </UiButton>
        </div>
      </form>
    </UiContainer>
  </main>
</template>
