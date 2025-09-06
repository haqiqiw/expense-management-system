import 'vue-router'

export type LayoutKey = 'default' | 'empty'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth: boolean
    roles?: string[]
    layout?: LayoutKey
  }
}
