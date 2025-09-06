<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white border-b border-gray-200">
      <div class="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-auto min-h-16 py-2 flex-wrap">
          <div class="flex items-center flex-wrap">
            <div class="flex-shrink-0 mr-6">
              <span class="text-xl font-bold">EMS</span>
            </div>
            <div class="flex items-center space-x-8">
              <RouterLink to="/home" :class="getLinkClass('/home')"> Beranda </RouterLink>
              <RouterLink to="/expenses" :class="getLinkClass('/expenses')">Pengeluaran</RouterLink>
              <RouterLink
                v-if="authStore.userRole === 'manager'"
                to="/approvals"
                :class="getLinkClass('/approvals')"
              >
                Persetujuan
              </RouterLink>
            </div>
          </div>

          <div class="ml-auto pl-4">
            <button
              @click="handleLogout"
              type="button"
              class="px-3 py-2 text-sm font-medium text-gray-500 bg-white rounded-md hover:text-gray-700"
            >
              Keluar
            </button>
          </div>
        </div>
      </div>
    </nav>

    <main class="py-10">
      <div class="px-4 mx-auto max-w-7xl sm:px-6 lg:px-8">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const getLinkClass = (path: string) => {
  const baseClasses = 'inline-flex items-center px-1 pt-1 text-sm font-medium'
  if (route.path === path) {
    return `${baseClasses} text-gray-900`
  }
  return `${baseClasses} text-gray-500 hover:text-gray-700`
}

const handleLogout = async () => {
  await authStore.logout()
  router.push({ name: 'login' })
}
</script>
