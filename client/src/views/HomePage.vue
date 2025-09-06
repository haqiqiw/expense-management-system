<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-800">Selamat datang!</h1>
    <p v-if="authStore.currentUser" class="mt-2 text-lg text-gray-600">
      Anda masuk sebagai <span class="underline">{{ authStore.currentUser.name }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>
