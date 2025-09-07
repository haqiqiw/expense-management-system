<template>
  <div>
    <div v-if="isLoading" class="px-4 py-10 text-center bg-white border border-gray-200 rounded-md">
      <p class="text-gray-500">Loading...</p>
    </div>

    <div
      v-else-if="errorMessage"
      class="px-4 py-5 text-center bg-white border border-gray-200 rounded-md"
    >
      <h3 class="text-xl font-semibold text-gray-800">Gagal Memuat Data</h3>
      <p class="mt-2 text-gray-600">{{ errorMessage }}</p>
      <RouterLink
        to="/expenses"
        class="mt-4 inline-block px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
      >
        Kembali
      </RouterLink>
    </div>

    <div v-else-if="expense" class="space-y-8">
      <div class="bg-white overflow-hidden sm:rounded-lg border border-gray-200 rounded-md">
        <div class="px-4 py-5 sm:px-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900">Detail Pengeluaran</h3>
        </div>
        <div class="border-t border-gray-200">
          <dl class="grid grid-cols-1 sm:grid-cols-3">
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Tanggal</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ formatDate(expense.created_at) }}</dd>
            </div>
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Nama</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ expense.user.name }}</dd>
            </div>
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Nominal</dt>
              <dd class="mt-1 text-sm font-medium text-gray-600">
                {{ formatRupiah(expense.amount_idr) }}
              </dd>
            </div>
            <div class="px-4 py-5 sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500">Deskripsi</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ expense.description }}</dd>
            </div>
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Status</dt>
              <dd class="mt-1 text-sm text-gray-900">
                <span
                  :class="getStatusClass(expense.status)"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                >
                  {{ getStatusText(expense.status) }}
                </span>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <div class="bg-white overflow-hidden sm:rounded-lg border border-gray-200 rounded-md">
        <div class="px-4 py-5 sm:px-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900">Struk / Nota</h3>
        </div>
        <div class="border-t border-gray-200 px-4 py-5 sm:p-6">
          <div v-if="expense.receipt_url">
            <img
              :src="expense.receipt_url"
              alt="Receipt image"
              class="max-h-96 w-auto rounded-lg border"
            />
          </div>
          <div v-else>
            <p id="receipt-empty" class="text-sm text-gray-500">
              Tidak ada struk/nota yang dilampirkan.
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white overflow-hidden sm:rounded-lg border border-gray-200 rounded-md">
        <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
          <h3 class="text-lg leading-6 font-medium text-gray-900">Persetujuan</h3>
          <button
            id="btn-approval"
            v-if="showApprovalAction"
            @click="isApprovalModalOpen = true"
            type="button"
            class="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700"
          >
            Aksi
          </button>
        </div>
        <div class="border-t border-gray-200">
          <div v-if="showApprovalNotes" class="px-4 pt-4">
            <p id="approval-warning" class="text-sm text-yellow-600 bg-yellow-50 p-2 rounded-md">
              Pengeluaran dengan nominal &lt; Rp 1.000.000 tidak membutuhkan persetujuan dari
              manager dan disetujui secara otomatis oleh sistem.
            </p>
          </div>
          <dl class="grid grid-cols-1 sm:grid-cols-3">
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Nama</dt>
              <dd class="mt-1 text-sm text-gray-900">
                {{ expense.approval?.approver_name || '-' }}
              </dd>
            </div>
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Tanggal</dt>
              <dd class="mt-1 text-sm text-gray-900">
                {{ expense.approval ? formatDate(expense.approval.created_at) : '-' }}
              </dd>
            </div>
            <div class="px-4 py-5 sm:col-span-1">
              <dt class="text-sm font-medium text-gray-500">Catatan</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ expense.approval?.notes || '-' }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <ApprovalActionModal
      v-if="expense"
      :show="isApprovalModalOpen"
      :expense-id="expense.id"
      @close="isApprovalModalOpen = false"
      @success="fetchExpenseDetails"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useExpenseStore } from '@/stores/expense'
import { useAuthStore } from '@/stores/auth'
import { formatRupiah, formatDate } from '@/utils/formatter'
import { getStatusClass, getStatusText } from '@/utils/status'
import ApprovalActionModal from '@/components/expense/ApprovalActionModal.vue'
import type { AxiosError } from 'axios'
import type { ApiError } from '@/utils/error'

const route = useRoute()
const expenseStore = useExpenseStore()
const authStore = useAuthStore()

const errorMessage = ref<string | null>(null)
const isApprovalModalOpen = ref(false)

const expense = computed(() => expenseStore.currentExpense)
const isLoading = computed(() => expenseStore.isLoading)
const currentUser = computed(() => authStore.currentUser)

const showApprovalAction = computed(() => {
  if (!expense.value || !currentUser.value) return false

  return (
    currentUser.value.role === 'manager' &&
    expense.value.status === 'awaiting_approval' &&
    expense.value.user.id !== currentUser.value.id
  )
})
const showApprovalNotes = computed(() => {
  return (expense.value?.amount_idr ?? 0) < 1000000
})

onMounted(async () => {
  await fetchExpenseDetails()
})

const fetchExpenseDetails = async () => {
  const expenseId = route.params.id as string
  if (expenseId) {
    try {
      await expenseStore.fetchExpenseById(expenseId)
    } catch (error) {
      const axError = error as AxiosError<ApiError>
      const status = axError.response?.status

      if (status === 403) {
        errorMessage.value = 'Anda tidak memiliki izin untuk melihat detail pengeluaran ini.'
      } else if (status === 404) {
        errorMessage.value = 'Data pengeluaran tidak ditemukan.'
      } else {
        errorMessage.value = 'Terjadi kesalahan saat mengambil data.'
      }
    }
  }
}

defineExpose({ isApprovalModalOpen })
</script>
