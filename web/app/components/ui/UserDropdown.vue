<script setup lang="ts">
const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'settings'): void
  (e: 'signOut'): void
}>()
const router = useRouter()
const { clear } = useUserSession()

interface Props {
  name: string
  email?: string
  avatar?: string
}

const initials = computed(() =>
  props.name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase(),
)
</script>

<template>
  <DropdownMenuRoot>
    <DropdownMenuTrigger as-child>
      <button
        class="rounded-full flex focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-white dark:focus:ring-offset-zinc-900"
      >
        <span class="sr-only">Open user menu</span>

        <AvatarRoot class="align-middle rounded-full inline-flex h-9 w-9 select-none items-center justify-center overflow-hidden">
          <AvatarImage
            class="h-full w-full object-cover"
            :src="avatar ?? ''"
            :alt="name"
          />
          <AvatarFallback
            class="text-sm text-zinc-700 font-semibold bg-zinc-200 flex h-full w-full items-center justify-center dark:text-zinc-200 dark:bg-zinc-700"
            :delay-ms="300"
          >
            {{ initials }}
          </AvatarFallback>
        </AvatarRoot>
      </button>
    </DropdownMenuTrigger>

    <DropdownMenuContent
      :side-offset="5"
      class="data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=open]:animate-in data-[state=closed]:animate-out rounded-md bg-white w-64 ring-1 ring-black/5 shadow-lg origin-top-right z-10 focus:outline-none dark:bg-zinc-900 dark:ring-white/10"
    >
      <DropdownMenuLabel class="font-normal px-4 py-3 border-b border-zinc-200 dark:border-zinc-800">
        <p class="text-sm text-zinc-500 dark:text-zinc-400">
          Signed in as
        </p>
        <p class="text-sm text-zinc-800 font-medium truncate dark:text-white">
          {{ email ?? name }}
        </p>
      </DropdownMenuLabel>
      <div class="py-1">
        <DropdownMenuItem
          class="text-sm text-zinc-700 px-4 py-2 flex gap-3 w-full cursor-pointer items-center dark:text-zinc-300 data-[highlighted]:outline-none data-[highlighted]:bg-zinc-100 dark:data-[highlighted]:bg-zinc-800"
          @select="emit('settings')"
        >
          <Icon name="ph:gear-six-bold" class="text-zinc-400 h-5 w-5" />
          <span>User Settings</span>
        </DropdownMenuItem>
      </div>
      <DropdownMenuSeparator class="bg-zinc-200 h-px dark:bg-zinc-800" />
      <div class="py-1">
        <DropdownMenuItem
          class="text-sm text-red-600 px-4 py-2 flex gap-3 w-full cursor-pointer items-center dark:text-red-500 data-[highlighted]:outline-none data-[highlighted]:bg-zinc-100 dark:data-[highlighted]:bg-zinc-800"
          @select="emit('signOut')"
        >
          <Icon name="ph:sign-out-bold" class="h-5 w-5" />
          <span>Sign out</span>
        </DropdownMenuItem>
      </div>
    </DropdownMenuContent>
  </DropdownMenuRoot>
</template>
