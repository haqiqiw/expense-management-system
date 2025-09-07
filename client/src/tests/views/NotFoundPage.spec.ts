import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import NotFoundPage from './../../views/NotFoundPage.vue'

const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a :to="to"><slot /></a>',
  props: ['to'],
}

describe('NotFoundPage.vue', () => {
  it('render all elements', () => {
    const wrapper = mount(NotFoundPage, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    const heading = wrapper.find('h1')
    expect(heading.exists()).toBe(true)
    expect(heading.text()).toBe('404')

    const subHeading = wrapper.find('p.text-2xl')
    expect(subHeading.exists()).toBe(true)
    expect(subHeading.text()).toBe('Halaman tidak ditemukan')

    const homeLink = wrapper.findComponent(RouterLinkStub)
    expect(homeLink.exists()).toBe(true)

    expect(homeLink.text()).toBe('Kembali ke beranda')
    expect(homeLink.props('to')).toBe('/home')
  })
})
