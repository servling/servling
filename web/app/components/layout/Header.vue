<script setup lang="ts">
const { user, clear } = useUserSession()
const router = useRouter()

async function handleSignOut() {
  await clear()
  await router.push('/login')
}
</script>

<template>
  <UiHeader
    :user="{
      name: user?.name ?? 'Error fetching name',
    }"
    @settings="() => router.push('/user/settings/')"
    @sign-out="handleSignOut"
  />
</template>
