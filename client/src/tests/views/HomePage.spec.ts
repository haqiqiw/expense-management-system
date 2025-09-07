import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import HomePage from './../../views/HomePage.vue'
import { useAuthStore } from '@/stores/auth'

describe('HomePage', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('displays the welcome message and user name', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: 1,
      name: 'John Doe',
      email: 'john@mail.com',
      role: 'manager',
      created_at: '2025-01-01T00:00:00Z',
    }

    const wrapper = mount(HomePage)

    const heading = wrapper.find('h1')
    expect(heading.exists()).toBe(true)
    expect(heading.text()).toBe('Selamat datang!')

    const paragraph = wrapper.find('p')
    expect(paragraph.exists()).toBe(true)
    expect(paragraph.text()).toContain('Anda masuk sebagai John Doe')
  })
})
