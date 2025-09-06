<template>
  <PopupModal
    :show="show"
    title="Persetujuan Pengeluaran"
    :closeable="!isLoading"
    @close="$emit('close')"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700">Status</label>
          <fieldset class="mt-2">
            <legend class="sr-only">Status Persetujuan</legend>
            <div class="flex items-center space-x-4">
              <div class="flex items-center">
                <input
                  id="approve"
                  value="approved"
                  v-model="form.status"
                  type="radio"
                  name="status"
                  class="w-4 h-4 text-indigo-600 border-gray-300"
                />
                <label for="approve" class="ml-2 block text-sm text-gray-900">Disetujui</label>
              </div>
              <div class="flex items-center">
                <input
                  id="reject"
                  value="rejected"
                  v-model="form.status"
                  type="radio"
                  name="status"
                  class="w-4 h-4 text-indigo-600 border-gray-300"
                />
                <label for="reject" class="ml-2 block text-sm text-gray-900">Ditolak</label>
              </div>
            </div>
          </fieldset>
        </div>

        <div>
          <label for="notes" class="block text-sm font-medium text-gray-700">Catatan</label>
          <textarea
            id="notes"
            v-model="form.notes"
            rows="4"
            class="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md0"
          ></textarea>
        </div>

        <div v-if="errorMessage" class="p-3 text-sm text-red-700 bg-red-100 rounded-md">
          {{ errorMessage }}
        </div>

        <div class="flex justify-end pt-4 space-x-3">
          <button
            type="button"
            @click="$emit('close')"
            :disabled="isLoading"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="isLoading || !isFormValid"
            class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 disabled:bg-indigo-300"
          >
            {{ isLoading ? 'Loading...' : 'Simpan' }}
          </button>
        </div>
      </div>
    </form>
  </PopupModal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useExpenseStore } from '@/stores/expense'
import PopupModal from '@/components/ui/PopupModal.vue'
import { extractErrorMessage } from '@/utils/error'

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
  expenseId: {
    type: Number,
    required: true,
  },
})
const emit = defineEmits(['close', 'success'])

const expenseStore = useExpenseStore()

const form = reactive<{
  status: 'approved' | 'rejected' | ''
  notes: string | null
}>({
  status: '',
  notes: null,
})
const isLoading = ref(false)
const errorMessage = ref<string | null>(null)

const isFormValid = computed(() => form.status !== '')

const resetForm = () => {
  form.status = ''
  form.notes = null
  errorMessage.value = null
}

const handleSubmit = async () => {
  if (!isFormValid.value) return

  isLoading.value = true
  errorMessage.value = null

  try {
    const payload = {
      id: props.expenseId,
      notes: form.notes,
    }

    if (form.status === 'approved') {
      await expenseStore.approveExpense(payload)
    } else if (form.status === 'rejected') {
      await expenseStore.rejectExpense(payload)
    }

    emit('success')
    emit('close')
  } catch (error) {
    errorMessage.value = extractErrorMessage(error)
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

watch(
  () => props.show,
  (newVal) => {
    if (!newVal) {
      setTimeout(resetForm, 300)
    }
  },
)
</script>
