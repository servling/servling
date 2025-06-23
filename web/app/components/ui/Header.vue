<script setup lang="ts">
interface User {
  name: string
  email?: string
  avatar?: string
}

interface Props {
  user: User
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'settings'): void
  (e: 'signOut'): void
}>()

const items = useBreadcrumbItems({ hideRoot: false, hideNonExisting: true })
</script>

<template>
  <div class="px-4 py-3 flex w-full items-center justify-between">
    <nav aria-label="Breadcrumb">
      <ol class="text-sm flex items-center">
        <li v-for="(item, index) in items" :key="item.label" class="flex items-center">
          <Icon
            v-if="index > 0"
            name="ph:caret-right-bold"
            class="text-zinc-500 mx-2 flex-shrink-0 h-4 w-4"
          />
          <NuxtLink
            v-if="!item.current"
            :to="item.to"
            class="text-zinc-400 flex gap-x-1.5 transition-colors items-center hover:text-white"
          >
            <Icon v-if="item.icon" :name="item.icon" class="flex-shrink-0 h-4 w-4" />
            <span>{{ item.label }}</span>
          </NuxtLink>
          <div
            v-else
            class="text-white font-semibold flex gap-x-1.5 items-center"
            aria-current="page"
          >
            <Icon v-if="item.icon" :name="item.icon" class="flex-shrink-0 h-4 w-4" />
            <span>{{ item.label }}</span>
          </div>
        </li>
      </ol>
    </nav>

    <UiUserDropdown :name="user.name" :email="user.email" :avatar="user.avatar" @sign-out="emit('signOut')" @settings="emit('settings')" />
  </div>
</template>
