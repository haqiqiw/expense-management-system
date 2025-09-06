<template>
  <div class="flex flex-col items-center justify-center min-h-[50vh] text-center">
    <h1 class="text-6xl font-bold text-indigo-600">404</h1>
    <p class="mt-4 text-2xl font-semibold text-gray-800">Halaman tidak ditemukan</p>
    <p class="mt-2 text-gray-500">Maaf, kami tidak bisa menemukan halaman yang anda cari.</p>
    <div class="mt-6">
      <RouterLink
        to="/home"
        class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
      >
        Kembali ke beranda
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>
