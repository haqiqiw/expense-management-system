<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900">Pengeluaran Saya</h1>
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mt-6">
      <div class="flex items-center mt-4 sm:mt-0 space-4">
        <select
          v-model="selectedFilter"
          @change="applyFilter"
          class="block w-full py-2 pl-3 pr-10 text-base bg-grey-600 border-gray-300 rounded-md focus:outline-none sm:text-sm border border-gray-300"
        >
          <option value="all">Semua Status</option>
          <option value="awaiting_approval">Menunggu Persetujuan</option>
          <option value="approved">Disetujui</option>
          <option value="rejected">Ditolak</option>
          <option value="auto_approved">Disetujui Otomatis</option>
        </select>
      </div>
      <button
        @click="isModalOpen = true"
        type="button"
        class="inline-flex items-center justify-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
      >
        + Tambah
      </button>
    </div>

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
                    Nominal
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
                <tr v-if="expenseStore.isLoading">
                  <td colspan="7" class="px-6 py-4 text-center text-gray-500">Loading...</td>
                </tr>
                <tr v-else-if="expenseStore.expenses.length === 0">
                  <td colspan="7" class="px-6 py-4 text-center text-gray-500">Data kosong</td>
                </tr>
                <tr v-for="expense in expenseStore.expenses" :key="expense.id">
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
          <span class="font-medium">{{ expenseStore.filters.offset + 1 }}</span>
          hingga
          <span class="font-medium">{{
            Math.min(expenseStore.filters.offset + expenseStore.filters.limit, expenseStore.total)
          }}</span>
          dari
          <span class="font-medium">{{ expenseStore.total }}</span>
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

  <AddExpenseModal :show="isModalOpen" @close="isModalOpen = false" />
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useExpenseStore } from '@/stores/expense'
import { formatRupiah, formatDate } from '@/utils/formatter'
import { getStatusClass, getStatusText } from '@/utils/status'
import { ArrowRightCircleIcon } from '@heroicons/vue/24/outline'
import type { ExpenseFiltersStatus } from '@/types'
import AddExpenseModal from '@/components/expense/AddExpenseModal.vue'

const expenseStore = useExpenseStore()
const selectedFilter = ref<string>('all')
const isModalOpen = ref(false)

onMounted(() => {
  expenseStore.fetchExpenses()
})

const applyFilter = () => {
  if (selectedFilter.value === 'all') {
    expenseStore.setFilters({ status: null, auto_approved: false })
  } else if (selectedFilter.value === 'auto_approved') {
    expenseStore.setFilters({ auto_approved: true })
  } else {
    expenseStore.setFilters({ status: selectedFilter.value as ExpenseFiltersStatus })
  }
}

const currentPage = computed(() => {
  return Math.floor(expenseStore.filters.offset / expenseStore.filters.limit) + 1
})

const totalPages = computed(() => {
  return Math.ceil(expenseStore.total / expenseStore.filters.limit)
})

const changePage = (page: number) => {
  if (page > 0 && page <= totalPages.value) {
    expenseStore.setPage(page)
  }
}
</script>
