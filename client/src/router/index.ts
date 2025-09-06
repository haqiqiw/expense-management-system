import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import HomePage from '@/views/HomePage.vue'
import LoginPage from '@/views/LoginPage.vue'
import ExpenseListPage from '@/views/ExpenseListPage.vue'
import ApprovalListPage from '@/views/ApprovalListPage.vue'
import ExpenseDetailPage from '@/views/ExpenseDetailPage.vue'
import NotFoundPage from '@/views/NotFoundPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/home',
      meta: {
        requiresAuth: true,
        layout: 'default',
      },
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: {
        requiresAuth: false,
        layout: 'empty',
      },
    },
    {
      path: '/home',
      name: 'home',
      component: HomePage,
      meta: {
        requiresAuth: true,
        layout: 'default',
      },
    },
    {
      path: '/expenses',
      name: 'expenses',
      component: ExpenseListPage,
      meta: {
        requiresAuth: true,
        layout: 'default',
      },
    },
    {
      path: '/approvals',
      name: 'approvals',
      component: ApprovalListPage,
      meta: {
        requiresAuth: true,
        roles: ['manager'],
        layout: 'default',
      },
    },
    {
      path: '/expenses/:id',
      name: 'expense-detail',
      component: ExpenseDetailPage,
      meta: {
        requiresAuth: true,
        layout: 'default',
      },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFoundPage,
      meta: {
        requiresAuth: false,
        layout: 'empty',
      },
    },
  ],
})

router.beforeEach(async (to, _, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth

  if (authStore.token && !authStore.user) {
    try {
      await authStore.fetchUser()
    } catch (error) {
      console.error('Router failed to fetch user:', error)
      authStore.clearUser()
      next({ name: 'login' })
    }
  }

  const isAuthenticated = authStore.isAuthenticated
  const userRole = authStore.userRole ?? ''
  const toRoles = to.meta.roles || []

  if (requiresAuth && !isAuthenticated) {
    next({ name: 'login' })
  } else if (toRoles.length > 0 && !toRoles.includes(userRole)) {
    next({ name: 'home' })
  } else if (to.name === 'login' && isAuthenticated) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
