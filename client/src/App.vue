<template>
  <component :is="layout">
    <RouterView />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import DefaultLayout from '@/components/layout/DefaultLayout.vue'
import EmptyLayout from '@/components/layout/EmptyLayout.vue'

const route = useRoute()

const layouts = {
  default: DefaultLayout,
  empty: EmptyLayout,
}

const layout = computed(() => {
  return layouts[route.meta.layout || 'empty']
})
</script>
