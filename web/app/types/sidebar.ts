export interface SidebarItem {
  title: string
  icon?: string
  href?: string
  elevated?: boolean
  action?: () => void
}
