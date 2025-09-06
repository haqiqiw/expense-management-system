<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900">Pengeluaran Butuh Persetujuan</h1>
    <div class="mt-4 flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="overflow-hidden border border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    ID
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Tanggal
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Nama
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Deskripsi
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Jumlah
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Status
                  </th>
                  <th
                    scope="col"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    Aksi
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-if="approvalStore.isLoading">
                  <td colspan="7" class="px-6 py-4 text-center text-gray-500">Loading...</td>
                </tr>
                <tr v-else-if="approvalStore.expenses.length === 0">
                  <td colspan="7" class="px-6 py-4 text-center text-gray-500">Data kosong</td>
                </tr>
                <tr v-for="expense in approvalStore.expenses" :key="expense.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ expense.id }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ formatDate(expense.created_at) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ expense.user.name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ expense.description }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ formatRupiah(expense.amount_idr) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span
                      :class="getStatusClass(expense.status)"
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                    >
                      {{ getStatusText(expense.status) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <RouterLink
                      :to="`/expenses/${expense.id}`"
                      class="text-indigo-600 hover:text-indigo-900"
                    >
                      <ArrowRightCircleIcon class="h-6 w-6" />
                    </RouterLink>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
    <div class="mt-4 flex justify-between items-center">
      <div>
        <p class="text-sm text-gray-700">
          Menampilkan
          <span class="font-medium">{{ approvalStore.filters.offset + 1 }}</span>
          hingga
          <span class="font-medium">{{
            Math.min(
              approvalStore.filters.offset + approvalStore.filters.limit,
              approvalStore.total,
            )
          }}</span>
          dari
          <span class="font-medium">{{ approvalStore.total }}</span>
          hasil
        </p>
      </div>
      <div class="flex space-x-2">
        <button
          @click="changePage(currentPage - 1)"
          :disabled="currentPage === 1 || totalPages < 1"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
        >
          Sebelumnya
        </button>
        <button
          @click="changePage(currentPage + 1)"
          :disabled="currentPage === totalPages || totalPages < 1"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
        >
          Berikutnya
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useApprovalStore } from '@/stores/approval'
import { formatRupiah, formatDate } from '@/utils/formatter'
import { getStatusClass, getStatusText } from '@/utils/status'
import { ArrowRightCircleIcon } from '@heroicons/vue/24/outline'

const approvalStore = useApprovalStore()

onMounted(() => {
  approvalStore.fetchExpenses()
})

const currentPage = computed(() => {
  return Math.floor(approvalStore.filters.offset / approvalStore.filters.limit) + 1
})

const totalPages = computed(() => {
  return Math.ceil(approvalStore.total / approvalStore.filters.limit)
})

const changePage = (page: number) => {
  if (page > 0 && page <= totalPages.value) {
    approvalStore.setPage(page)
  }
}
</script>
