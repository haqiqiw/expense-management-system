import { describe, it, expect, beforeEach, vi, type Mocked } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ExpenseListPage from './../../views/ExpenseListPage.vue'
import { useExpenseStore } from '@/stores/expense'
import apiClient from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn(),
  },
}))

const RouterLinkStub = {
  template: '<a><slot /></a>',
  props: ['to'],
}

const mockExpense = {
  id: 1,
  amount_idr: 500000,
  description: 'Team Lunch',
  receipt_url: null,
  status: 'awaiting_approval' as const,
  requires_approval: true,
  auto_approved: false,
  created_at: '2025-09-06T10:00:00Z',
  user: { id: 2, name: 'Budi', email: 'budi@mail.com' },
}

describe('ExpenseListPage', () => {
  const mockedApi = apiClient as Mocked<typeof apiClient>

  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('render loading message', async () => {
    mockedApi.get.mockResolvedValue({
      data: { data: [], meta: { total: 0 } },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.find('td').text()).toBe('Loading...')
  })

  it('render empty message when no expenses', async () => {
    mockedApi.get.mockResolvedValue({
      data: { data: [], meta: { total: 0 } },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('td').text()).toBe('Data kosong')
  })

  it('render table with expense data', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        data: [mockExpense],
        meta: {
          total: 1,
          limit: 5,
          offset: 0,
        },
      },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.html()).toContain('Pengeluaran Saya')

    const tableRows = wrapper.findAll('tbody tr')
    expect(tableRows.length).toBe(1)

    const cells = tableRows[0].findAll('td')
    expect(cells[0].text()).toBe('1')
    expect(cells[1].text()).toBe('6/9/2025 17:00')
    expect(cells[2].text()).toBe('Budi')
    expect(cells[4].text()).toBe('Rp\u00a0500.000')
  })

  it('call setPage action when next button is clicked', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        data: [mockExpense],
        meta: {
          total: 10,
          limit: 5,
          offset: 0,
        },
      },
    })

    const approvalStore = useExpenseStore()
    const setPageSpy = vi.spyOn(approvalStore, 'setPage')

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    const nextButton = wrapper.find('#btn-next')
    await nextButton.trigger('click')

    expect(setPageSpy).toHaveBeenCalledTimes(1)
    expect(setPageSpy).toHaveBeenCalledWith(2)
  })

  it('go to detail page when detail icon is clicked', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        data: [mockExpense],
        meta: {
          total: 1,
          limit: 5,
          offset: 0,
        },
      },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    const detailLink = wrapper.findComponent(RouterLinkStub)
    expect(detailLink.props('to')).toBe('/expenses/1')
  })

  it('set filter when select option is changed', async () => {
    const expenseStore = useExpenseStore()
    const setFiltersSpy = vi.spyOn(expenseStore, 'setFilters')

    mockedApi.get.mockResolvedValue({
      data: {
        data: [mockExpense],
        meta: {
          total: 1,
          limit: 5,
          offset: 0,
        },
      },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    const select = wrapper.find('#select-filter')
    await select.setValue('approved')

    expect(setFiltersSpy).toHaveBeenCalledTimes(1)
    expect(setFiltersSpy).toHaveBeenCalledWith({ status: 'approved' })
  })

  it('set isModalOpen to true when add button is clicked', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        data: [mockExpense],
        meta: {
          total: 1,
          limit: 5,
          offset: 0,
        },
      },
    })

    const wrapper = mount(ExpenseListPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.vm.isModalOpen).toBe(false)

    const addButton = wrapper.find('#btn-add')
    await addButton.trigger('click')

    expect(wrapper.vm.isModalOpen).toBe(true)
  })
})
