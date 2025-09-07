import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AddExpenseModal from './../../../components/expense/AddExpenseModal.vue'
import { useExpenseStore } from '@/stores/expense'

const mockSetValue = vi.fn()
const mockNumberValue = { value: 10000 }

vi.mock('vue-currency-input', () => ({
  useCurrencyInput: () => ({
    inputRef: vi.fn(),
    numberValue: mockNumberValue,
    setValue: mockSetValue,
  }),
}))

const PopupModalStub = {
  name: 'PopupModal',
  template: '<div><slot v-if="show" /></div>',
  props: ['show'],
}

describe('AddExpenseModal', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('call create action success', async () => {
    const expenseStore = useExpenseStore()
    const createExpenseSpy = vi.spyOn(expenseStore, 'createExpense').mockResolvedValue()

    const wrapper = mount(AddExpenseModal, {
      props: {
        show: true,
      },
      global: {
        stubs: { PopupModal: PopupModalStub },
      },
    })

    await wrapper.find('textarea#description').setValue('Something')
    await wrapper.find('form').trigger('submit.prevent')

    await flushPromises()

    expect(createExpenseSpy).toHaveBeenCalledWith({
      amount_idr: 10000,
      description: 'Something',
      receipt_url: null,
    })

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('show error message when error', async () => {
    const expenseStore = useExpenseStore()
    vi.spyOn(expenseStore, 'createExpense').mockRejectedValue(new Error('Unexpected error'))

    const wrapper = mount(AddExpenseModal, {
      props: {
        show: true,
      },
      global: {
        stubs: { PopupModal: PopupModalStub },
      },
    })

    await wrapper.find('textarea#description').setValue('Something')
    await wrapper.find('form').trigger('submit.prevent')

    await flushPromises()

    const errorMessage = wrapper.find('div.bg-red-100')
    expect(errorMessage.exists()).toBe(true)
  })
})
