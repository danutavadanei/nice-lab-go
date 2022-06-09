import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/IndexView.vue'),
      meta: {
        requiresAuth: true
      },
    },
    {
      path: '/lab/:lab',
      name: 'connect',
      component: () => import('../views/ConnectView.vue'),
      meta: {
        requiresAuth: true
      },
    },
    {
      path: '/session/:session',
      name: 'monitor',
      component: () => import('../views/MonitorView.vue'),
      meta: {
        requiresAuth: true
      },
    },
      {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.getters.isLoggedIn) {
      next({ name: 'login' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
