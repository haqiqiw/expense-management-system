<template>
  <PopupModal
    :show="show"
    title="Tambah Pengeluaran"
    :closeable="!isLoading"
    @close="$emit('close')"
  >
    <form @submit.prevent="handleSubmit">
      <div class="space-y-4">
        <div>
          <label for="amount" class="block text-sm font-medium text-gray-700">Nominal</label>
          <input
            id="amount"
            ref="inputRef"
            type="text"
            class="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md"
            placeholder="10.000"
          />
          <p class="mt-2 text-xs text-gray-500">
            Minimum: {{ formatRupiah(10000) }} &nbsp; / &nbsp; Maksimum:
            {{ formatRupiah(50000000) }}
          </p>
          <p v-if="isAmountHigh" class="mt-2 text-sm text-yellow-600 bg-yellow-50 p-2 rounded-md">
            Nominal >= Rp 1.000.000 membutuhkan persetujuan dari manager.
          </p>
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-gray-700">Deskripsi</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="4"
            required
            class="w-full px-3 py-2 mt-1 border border-gray-300 rounded-md0"
            placeholder="Deksripsi"
          ></textarea>
        </div>

        <div>
          <label for="receipt" class="block text-sm font-medium text-gray-700">Struk / Nota</label>
          <input
            id="receipt"
            type="file"
            @change="handleFileUpload"
            accept="image/png, image/jpeg, image/webp"
            class="w-full mt-1 text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100"
          />
          <p v-if="fileName" class="mt-2 text-sm text-gray-600">File: {{ fileName }}</p>
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
import { useCurrencyInput } from 'vue-currency-input'
import PopupModal from '@/components/ui/PopupModal.vue'
import { formatRupiah } from '@/utils/formatter'
import { extractErrorMessage } from '@/utils/error'

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
})
const emit = defineEmits(['close'])

const expenseStore = useExpenseStore()

const { inputRef, numberValue, setValue } = useCurrencyInput({
  currency: 'IDR',
  locale: 'id-ID',
  precision: 0,
  valueRange: {
    min: 10000,
    max: 50000000,
  },
  hideCurrencySymbolOnFocus: true,
  hideGroupingSeparatorOnFocus: false,
  hideNegligibleDecimalDigitsOnFocus: true,
})

const form = reactive<{
  description: string
  receipt_url: string | null
}>({
  description: '',
  receipt_url: null,
})
const fileName = ref('')
const isLoading = ref(false)
const errorMessage = ref<string | null>(null)

const isAmountHigh = computed(() => (numberValue.value ?? 0) >= 1000000)
const isFormValid = computed(() => (numberValue.value ?? 0) > 0 && form.description.trim() !== '')

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    fileName.value = target.files[0].name
    // mock the image upload flow
    // if a file is selected, send a dummy URL, otherwise send null
    form.receipt_url = 'https://placehold.co/500x700'
  } else {
    fileName.value = ''
    form.receipt_url = null
  }
}

const resetForm = () => {
  setValue(0)
  form.description = ''
  form.receipt_url = ''
  fileName.value = ''
  errorMessage.value = null
}

const handleSubmit = async () => {
  if (!isFormValid.value) return

  isLoading.value = true
  errorMessage.value = null

  try {
    await expenseStore.createExpense({
      amount_idr: numberValue.value ?? 0,
      description: form.description,
      receipt_url: form.receipt_url || null,
    })
    resetForm()
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
