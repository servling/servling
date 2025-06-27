<script setup lang="ts">
definePageMeta({
  breadcrumb: {
    icon: 'ph:storefront',
    ariaLabel: 'App Store',
    label: 'App Store',
  },
  middleware: ['applications-store'],
})

interface Template {
  name: string
  description: string
  services: 0
  avatarUrl: string
}

const templates: Template[] = [{
  name: 'Plex',
  description: 'A client-server media player system to access and stream your media from your own server.',
  services: 0,
  avatarUrl: 'https://avatars.githubusercontent.com/u/324832',
}, {
  name: 'Nextcloud',
  description: 'A suite of client-server software for creating and using file hosting services and other productivity tools.',
  services: 0,
  avatarUrl: 'https://upload.wikimedia.org/wikipedia/commons/6/60/Nextcloud_Logo.svg',
}, {
  name: 'Home Assistant',
  description: 'Open source home automation that puts local control and privacy first.',
  services: 0,
  avatarUrl: 'https://avatars.githubusercontent.com/u/13844975',
}, {
  name: 'Bitwarden',
  description: 'A free and open-source password management service that stores sensitive information in an encrypted vault.',
  services: 0,
  avatarUrl: 'https://raw.githubusercontent.com/bitwarden/brand/main/icons/icon.svg',
}, {
  name: 'Jellyfin',
  description: 'A free software media system that lets you control and stream your media.',
  services: 0,
  avatarUrl: 'https://avatars.githubusercontent.com/u/45698031',
}, {
  name: 'PhotoPrism',
  description: 'A server-based application for Browse, organizing and sharing your personal photo collection.',
  services: 0,
  avatarUrl: 'https://dl.photoprism.app/img/logo/logo.svg',
}, {
  name: 'Gitea',
  description: 'A self-hosted Git service that is lightweight and easy to set up.',
  services: 0,
  avatarUrl: 'https://avatars.githubusercontent.com/u/12724356',
}, {
  name: 'Vaultwarden',
  description: 'An unofficial Bitwarden compatible server written in Rust, formerly known as bitwarden_rs.',
  services: 0,
  avatarUrl: 'https://cdn.shortpixel.ai/spai/q_lossy+ret_img+to_auto/linuxiac.com/wp-content/uploads/2023/07/vaultwarden-1024x768.jpg',
}, {
  name: 'AdGuard Home',
  description: 'Network-wide software for blocking ads and tracking on all your devices.',
  services: 0,
  avatarUrl: 'https://avatars.githubusercontent.com/u/8361145',
}, {
  name: 'Uptime Kuma',
  description: 'A fancy self-hosted monitoring tool to keep track of your services\' uptime.',
  services: 0,
  avatarUrl: 'https://raw.githubusercontent.com/louislam/uptime-kuma/master/public/icon.svg',
}]
</script>

<template>
  <main>
    <UiContainer constrained>
      <UiPageHeader
        title="App Store"
        description="Find applications in the app store to install on your server."
      />
      <div class="gap-4 grid grid-cols-2">
        <UiCard v-for="template in templates" :key="template.name">
          <div class="flex gap-4">
            <NuxtImg :src="template.avatarUrl" width="64px" height="64px" fit="contain" class="object-contain" />
            <div class="flex">
              <div>
                <h3 class="text-2xl font-bold">
                  {{ template.name }}
                </h3>
                <span class="text-zinc-700 dark:text-zinc-300">{{ template.description }}</span>
              </div>
              <UiButton intent="primary" size="sm">
                Install
              </UiButton>
            </div>
          </div>
        </UiCard>
      </div>
    </UiContainer>
  </main>
</template>
