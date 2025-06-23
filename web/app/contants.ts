import type { SidebarItem } from '~/types/sidebar'

export const sidebarItems: SidebarItem[] = [
  {
    title: 'Dashboard',
    icon: 'ph:house',
    href: '/',
  },
  {
    title: 'Applications',
    icon: 'ph:squares-four',
    href: '/apps',
  },
  {
    title: 'App Store',
    icon: 'ph:storefront',
    href: '/apps/store',
  },
]

export const sidebarBottomItems: SidebarItem[] = [
  {
    title: 'Instance Settings',
    icon: 'ph:gear',
    href: '/settings',
  },
]
