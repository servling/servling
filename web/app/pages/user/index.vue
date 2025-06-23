<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:gear',
    ariaLabel: 'Settings',
    label: 'Settings',
  },
})

const { user } = useUserSession()

const name = ref(user.value?.name ?? '')
const password = ref('')

watchEffect(() => {
  name.value = user.value?.name ?? ''
})
</script>

<template>
  <main>
    <UiContainer constrained>
      <UiPageHeader
        title="Instance Settings"
        description="Manage your Servling instance."
      />

      <div class="mt-8 space-y-10">
        <UiCard>
          <template #header>
            <UiPageHeader
              :level="2"
              title="Profile"
              description="This information will be used for your account."
            />
          </template>

          <div class="max-w-xl space-y-6">
            <UiInput
              id="name"
              v-model="name"
              label="Full Name"
              type="text"
              left-icon="ph:user"
            />
            <UiInput
              id="password"
              v-model="password"
              label="Password"
              type="password"
              placeholder="••••••••"
              left-icon="ph:at-bold"
            />
          </div>

          <template #footer>
            <div class="flex justify-end">
              <UiButton>Save Changes</UiButton>
            </div>
          </template>
        </UiCard>

        <UiCard>
          <template #header>
            <UiPageHeader
              :level="2"
              title="Danger Zone"
            />
          </template>

          <div class="flex flex-col gap-4 items-start sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h3 class="text-base text-zinc-900 leading-6 font-semibold dark:text-white">
                Delete Your Account
              </h3>
              <p class="text-sm text-zinc-500 mt-1 max-w-2xl">
                Once you delete your account, all of your data will be permanently removed. This action cannot be undone.
              </p>
            </div>
            <UiButton intent="danger">
              Delete Account
            </UiButton>
          </div>
        </UiCard>
      </div>
    </UiContainer>
  </main>
</template>
