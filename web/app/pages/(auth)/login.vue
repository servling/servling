<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

definePageMeta({
  layout: 'auth',
  auth: false,
})

const validationSchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Username is required'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
  }),
)

const { handleSubmit, isSubmitting } = useForm({
  validationSchema,
})

const apiErrors = ref<string[]>([])

const onSubmit = handleSubmit(async (values) => {
  apiErrors.value = []
  try {
    const result = await login({
      composable: '$fetch',
      body: {
        username: values.name,
        password: values.password,
      },
    })

    if (!result) {
      apiErrors.value.push('Error fetching data.')
      return
    }

    await updateSession({
      accessToken: result.accessToken,
      accessTokenExpiresAt: result.accessTokenExpiresAt,
      refreshToken: result.refreshToken,
      refreshTokenExpiresAt: result.refreshTokenExpiresAt,
      user: {
        id: result.user.id,
        name: result.user.name,
        createdAt: result.user.createdAt,
        updatedAt: result.user.updatedAt,
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
    <UiAlert
      v-if="apiErrors.length > 0"
      variant="error"
      title="Login Error"
      class="mb-4"
      @dismiss="apiErrors = []"
    >
      <ul class="pl-5 list-disc">
        <li v-for="(error, index) in apiErrors" :key="index">
          {{ error }}
        </li>
      </ul>
    </UiAlert>

    <div class="mt-4 text-center">
      <h2 class="text-2xl text-zinc-900 leading-9 tracking-tight font-bold dark:text-white">
        Sign in to your account
      </h2>
      <p class="text-sm text-zinc-500 leading-6 mt-2">
        Don't have an account?
        <NuxtLink to="/register" class="text-blue-600 font-semibold dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
          Sign up
        </NuxtLink>
      </p>
    </div>

    <form class="mt-8" @submit="onSubmit">
      <UiCard>
        <div class="space-y-6">
          <!-- Email Field -->
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

          <!-- Password Field -->
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
        </div>

        <template #footer>
          <UiButton type="submit" class="w-full" :loading="isSubmitting">
            Sign In
          </UiButton>
        </template>
      </UiCard>
    </form>
  </div>
</template>
