import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ApprovalActionModal from './../../../components/expense/ApprovalActionModal.vue'
import { useExpenseStore } from '@/stores/expense'

const PopupModalStub = {
  name: 'PopupModal',
  template: '<div><slot v-if="show" /></div>',
  props: ['show'],
}

describe('ApprovalActionModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('call approve action success', async () => {
    const expenseStore = useExpenseStore()

    const approveSpy = vi.spyOn(expenseStore, 'approveExpense').mockResolvedValue()
    const rejectSpy = vi.spyOn(expenseStore, 'rejectExpense')

    const wrapper = mount(ApprovalActionModal, {
      props: {
        show: true,
        expenseId: 42,
      },
      global: {
        stubs: { PopupModal: PopupModalStub },
      },
    })

    await wrapper.find('input#approve').setValue(true)
    await wrapper.find('textarea#notes').setValue('Looks good')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(approveSpy).toHaveBeenCalledWith({
      id: 42,
      notes: 'Looks good',
    })
    expect(rejectSpy).not.toHaveBeenCalled()

    expect(wrapper.emitted('success')).toBeTruthy()
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('call reject action success', async () => {
    const expenseStore = useExpenseStore()

    const approveSpy = vi.spyOn(expenseStore, 'approveExpense')
    const rejectSpy = vi.spyOn(expenseStore, 'rejectExpense').mockResolvedValue()

    const wrapper = mount(ApprovalActionModal, {
      props: {
        show: true,
        expenseId: 42,
      },
      global: {
        stubs: { PopupModal: PopupModalStub },
      },
    })

    await wrapper.find('input#reject').setValue(true)
    await wrapper.find('textarea#notes').setValue('Invalid reason')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(rejectSpy).toHaveBeenCalledWith({
      id: 42,
      notes: 'Invalid reason',
    })
    expect(approveSpy).not.toHaveBeenCalled()

    expect(wrapper.emitted('success')).toBeTruthy()
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('show error message when error', async () => {
    const expenseStore = useExpenseStore()

    const approveSpy = vi.spyOn(expenseStore, 'approveExpense')
    const rejectSpy = vi
      .spyOn(expenseStore, 'rejectExpense')
      .mockRejectedValue(new Error('Unexpected error'))

    const wrapper = mount(ApprovalActionModal, {
      props: {
        show: true,
        expenseId: 42,
      },
      global: {
        stubs: { PopupModal: PopupModalStub },
      },
    })

    await wrapper.find('input#reject').setValue(true)
    await wrapper.find('textarea#notes').setValue('Invalid reason')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(rejectSpy).toHaveBeenCalledWith({
      id: 42,
      notes: 'Invalid reason',
    })
    expect(approveSpy).not.toHaveBeenCalled()

    const errorMessage = wrapper.find('div.bg-red-100')
    expect(errorMessage.exists()).toBe(true)
  })
})
