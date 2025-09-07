import { describe, it, expect, beforeEach, vi, type Mocked } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import LoginPage from './../../views/LoginPage.vue'
import apiClient from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
  },
  setAuthToken: vi.fn(),
}))

const mockRouterPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockRouterPush,
  }),
}))

describe('LoginPage', () => {
  const mockedApi = apiClient as Mocked<typeof apiClient>

  beforeEach(() => {
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('redirect to home when success', async () => {
    mockedApi.post.mockResolvedValueOnce({
      data: { data: { access_token: 'fake-jwt-token' } },
    })
    mockedApi.get.mockResolvedValueOnce({
      data: { data: { id: 1, name: 'John Doe', email: 'john@mail.com', role: 'manager' } },
    })

    const wrapper = mount(LoginPage)

    await wrapper.find('input#email').setValue('john@mail.com')
    await wrapper.find('input#password').setValue('password123')

    await wrapper.find('form').trigger('submit.prevent')

    const button = wrapper.find('button[type="submit"]')
    expect(button.text()).toBe('Loading...')
    expect(button.attributes('disabled')).toBeDefined()

    await flushPromises()

    expect(mockRouterPush).toHaveBeenCalledWith({ name: 'home' })
  })

  it('show error message when error', async () => {
    mockedApi.post.mockRejectedValue(new Error('Unexpected error'))

    const wrapper = mount(LoginPage)

    await wrapper.find('input#email').setValue('john@mail.com')
    await wrapper.find('input#password').setValue('password123')

    await wrapper.find('form').trigger('submit.prevent')

    await flushPromises()

    const errorMessage = wrapper.find('div.bg-red-100')
    expect(errorMessage.exists()).toBe(true)
    expect(mockRouterPush).not.toHaveBeenCalled()
  })
})
