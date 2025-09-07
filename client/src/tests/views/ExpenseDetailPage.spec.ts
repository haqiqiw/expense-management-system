import { describe, it, expect, beforeEach, vi, type Mocked } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ExpenseDetailPage from './../../views/ExpenseDetailPage.vue'
import { useAuthStore } from '@/stores/auth'
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

vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: {
      id: '15',
    },
  }),
}))

const mockExpense = {
  id: 1,
  amount_idr: 500000,
  description: 'Team Lunch',
  receipt_url: 'https://example.com/receipt.jpg',
  status: 'awaiting_approval' as const,
  requires_approval: true,
  auto_approved: false,
  created_at: '2025-09-06T10:00:00Z',
  processed_at: null,
  user: { id: 2, name: 'Budi', email: 'budi@mail.com' },
  approval: null,
}

describe('ExpenseDetailPage', () => {
  const mockedApi = apiClient as Mocked<typeof apiClient>

  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('render loading message', async () => {
    mockedApi.get.mockResolvedValue({
      data: { data: mockExpense },
    })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await wrapper.vm.$nextTick()

    expect(wrapper.find('p').text()).toBe('Loading...')
  })

  it('render error message when forbidden error', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 403 },
    })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('h3').text()).toBe('Gagal Memuat Data')
    expect(wrapper.find('p').text()).toBe(
      'Anda tidak memiliki izin untuk melihat detail pengeluaran ini.',
    )
  })

  it('render error message when not found error', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 404 },
    })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('h3').text()).toBe('Gagal Memuat Data')
    expect(wrapper.find('p').text()).toBe('Data pengeluaran tidak ditemukan.')
  })

  it('render error message when unexpected error', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 500 },
    })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('h3').text()).toBe('Gagal Memuat Data')
    expect(wrapper.find('p').text()).toBe('Terjadi kesalahan saat mengambil data.')
  })

  it('render error message when unexpected error', async () => {
    mockedApi.get.mockRejectedValue({
      response: { status: 500 },
    })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('h3').text()).toBe('Gagal Memuat Data')
    expect(wrapper.find('p').text()).toBe('Terjadi kesalahan saat mengambil data.')
  })

  it('render expense data without approval data', async () => {
    mockedApi.get.mockResolvedValue({ data: { data: mockExpense } })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.findAll('h3')[0].text()).toBe('Detail Pengeluaran')
    expect(wrapper.findAll('dd')[0].text()).toBe('6/9/2025 17:00')
    expect(wrapper.findAll('dd')[1].text()).toBe('Budi')
    expect(wrapper.findAll('dd')[2].text()).toBe('Rp\u00a0500.000')
    expect(wrapper.findAll('dd')[3].text()).toBe('Team Lunch')
    expect(wrapper.findAll('dd')[4].text()).toBe('Menunggu Persetujuan')

    expect(wrapper.findAll('h3')[1].text()).toBe('Struk / Nota')

    const receiptImage = wrapper.find('img[alt="Receipt image"]')
    expect(receiptImage.exists()).toBe(true)
    expect(receiptImage.attributes('src')).toBe('https://example.com/receipt.jpg')

    expect(wrapper.findAll('h3')[2].text()).toBe('Persetujuan')
    expect(wrapper.find('#approval-warning').text()).toBe(
      'Pengeluaran dengan nominal < Rp 1.000.000 tidak membutuhkan persetujuan dari manager dan disetujui secara otomatis oleh sistem.',
    )
    expect(wrapper.findAll('dd')[5].text()).toBe('-')
    expect(wrapper.findAll('dd')[6].text()).toBe('-')
    expect(wrapper.findAll('dd')[7].text()).toBe('-')
  })

  it('render expense data with approval data', async () => {
    const mockExpense = {
      id: 1,
      amount_idr: 500000,
      description: 'Team Lunch',
      receipt_url: null,
      status: 'completed' as const,
      requires_approval: true,
      auto_approved: false,
      created_at: '2025-09-06T10:00:00Z',
      processed_at: '2025-09-06T10:50:58Z',
      user: { id: 2, name: 'Budi', email: 'budi@mail.com' },
      approval: {
        id: 5,
        approver_id: 1,
        approver_email: 'john@mail.com',
        approver_name: 'John',
        status: 'approved',
        notes: 'Approve from me',
        created_at: '2025-09-06T10:40:58Z',
      },
    }

    mockedApi.get.mockResolvedValue({ data: { data: mockExpense } })

    const wrapper = mount(ExpenseDetailPage, {
      global: {
        stubs: {
          ApprovalActionModal: true,
          RouterLink: RouterLinkStub,
        },
      },
    })

    await flushPromises()

    expect(wrapper.findAll('h3')[0].text()).toBe('Detail Pengeluaran')
    expect(wrapper.findAll('dd')[0].text()).toBe('6/9/2025 17:00')
    expect(wrapper.findAll('dd')[1].text()).toBe('Budi')
    expect(wrapper.findAll('dd')[2].text()).toBe('Rp\u00a0500.000')
    expect(wrapper.findAll('dd')[3].text()).toBe('Team Lunch')
    expect(wrapper.findAll('dd')[4].text()).toBe('Selesai')

    expect(wrapper.findAll('h3')[1].text()).toBe('Struk / Nota')
    expect(wrapper.find('#receipt-empty').text()).toBe('Tidak ada struk/nota yang dilampirkan.')

    expect(wrapper.findAll('h3')[2].text()).toBe('Persetujuan')
    expect(wrapper.find('#approval-warning').text()).toBe(
      'Pengeluaran dengan nominal < Rp 1.000.000 tidak membutuhkan persetujuan dari manager dan disetujui secara otomatis oleh sistem.',
    )
    expect(wrapper.findAll('dd')[5].text()).toBe('John')
    expect(wrapper.findAll('dd')[6].text()).toBe('6/9/2025 17:40')
    expect(wrapper.findAll('dd')[7].text()).toBe('Approve from me')
  })

  describe('Approval Action Button Visibility', () => {
    it('show action button for manager viewing another user pending expense', async () => {
      mockedApi.get.mockResolvedValue({ data: { data: mockExpense } })

      const authStore = useAuthStore()
      authStore.user = {
        id: 1,
        name: 'John',
        role: 'manager',
        email: 'john@mail.com',
        created_at: '',
      }

      const wrapper = mount(ExpenseDetailPage, {
        global: {
          stubs: {
            ApprovalActionModal: true,
            RouterLink: RouterLinkStub,
          },
        },
      })

      await flushPromises()

      expect(wrapper.vm.isApprovalModalOpen).toBe(false)

      const actionButton = wrapper.find('#btn-approval')
      expect(actionButton.exists()).toBe(true)
      expect(actionButton.text()).toBe('Aksi')

      await actionButton.trigger('click')

      expect(wrapper.vm.isApprovalModalOpen).toBe(true)
    })

    it('hide action button if current user is not manager', async () => {
      mockedApi.get.mockResolvedValue({ data: { data: mockExpense } })

      const authStore = useAuthStore()
      authStore.user = {
        id: 3,
        name: 'wawan',
        role: 'employee',
        email: 'wawan@mail.com',
        created_at: '',
      }

      const wrapper = mount(ExpenseDetailPage, {
        global: {
          stubs: {
            ApprovalActionModal: true,
            RouterLink: RouterLinkStub,
          },
        },
      })

      await flushPromises()

      const actionButton = wrapper.find('#btn-approval')
      expect(actionButton.exists()).toBe(false)
    })

    it('hide action button if expense is already approved', async () => {
      const approvedExpense = { ...mockExpense, status: 'approved' as const }
      mockedApi.get.mockResolvedValue({ data: { data: approvedExpense } })

      const authStore = useAuthStore()
      authStore.user = {
        id: 1,
        name: 'John',
        role: 'manager',
        email: 'john@mail.com',
        created_at: '',
      }

      const wrapper = mount(ExpenseDetailPage, {
        global: {
          stubs: {
            ApprovalActionModal: true,
            RouterLink: RouterLinkStub,
          },
        },
      })

      await flushPromises()

      const actionButton = wrapper.find('#btn-approval')
      expect(actionButton.exists()).toBe(false)
    })
  })
})
