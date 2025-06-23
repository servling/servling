<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'
import { RegisterDocument } from '~/gql/mutations/register'

definePageMeta({
  layout: 'auth',
  auth: false,
})

const { executeMutation: mutateRegister } = useMutation(RegisterDocument)

const validationSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Name is required'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string().min(8, 'Please confirm your password'),
  }).refine(data => data.password === data.confirmPassword, {
    message: 'Passwords don\'t match',
    path: ['confirmPassword'],
  }),
)

const { handleSubmit, isSubmitting } = useForm({
  validationSchema,
})

const apiErrors = ref<string[]>([])

const onSubmit = handleSubmit(async (values) => {
  apiErrors.value = []
  try {
    const result = await mutateRegister({
      username: values.name,
      password: values.password,
    })

    if (!result.data) {
      apiErrors.value.push('An unknown error occurred')
      return
    }

    await updateSession({
      accessToken: result.data.register.accessToken,
      accessTokenExpiresAt: result.data.register.accessTokenExpiresAt,
      refreshToken: result.data.register.refreshToken,
      refreshTokenExpiresAt: result.data.register.refreshTokenExpiresAt,
      user: {
        id: result.data.register.user.id,
        name: result.data.register.user.name,
        createdAt: result.data.register.user.createdAt,
        updatedAt: result.data.register.user.updatedAt,
      },
    })

    await navigateTo('/')
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
  <div>
    <div class="mt-4 text-center">
      <h2 class="text-2xl text-zinc-900 leading-9 tracking-tight font-bold dark:text-white">
        Create an account
      </h2>
      <p class="text-sm text-zinc-500 leading-6 mt-2">
        Already have an account?
        <NuxtLink to="/login" class="text-blue-600 font-semibold dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
          Sign in
        </NuxtLink>
      </p>
    </div>

    <form class="mt-8" @submit="onSubmit">
      <UiCard>
        <div class="space-y-6">
          <UiAlert
            v-if="apiErrors.length > 0"
            variant="error"
            title="Registration Error"
            class="mb-4"
            @dismiss="apiErrors = []"
          >
            <ul class="pl-5 list-disc">
              <li v-for="(error, index) in apiErrors" :key="index">
                {{ error }}
              </li>
            </ul>
          </UiAlert>
          <UiFormField v-slot="{ value, errorMessage, handleChange }" name="name">
            <UiInput
              id="name"
              :model-value="value"
              :error="errorMessage"
              label="Username"
              type="text"
              placeholder="johndoe"
              left-icon="ph:user-bold"
              @update:model-value="handleChange"
            />
          </UiFormField>

          <UiFormField v-slot="{ value, errorMessage, handleChange }" name="password">
            <UiInput
              id="password"
              :model-value="value"
              :error="errorMessage"
              label="Password"
              type="password"
              placeholder="••••••••"
              left-icon="ph:lock-key-bold"
              @update:model-value="handleChange"
            />
          </UiFormField>

          <UiFormField v-slot="{ value, errorMessage, handleChange }" name="confirmPassword">
            <UiInput
              id="confirmPassword"
              :model-value="value"
              :error="errorMessage"
              label="Confirm Password"
              type="password"
              placeholder="••••••••"
              left-icon="ph:lock-key-bold"
              @update:model-value="handleChange"
            />
          </UiFormField>
        </div>

        <template #footer>
          <UiButton type="submit" class="w-full" :loading="isSubmitting">
            Create Account
          </UiButton>
        </template>
      </UiCard>
    </form>
  </div>
</template>
