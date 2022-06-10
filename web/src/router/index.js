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
    },
    {
      path: '/logout',
      name: 'logout',
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.name === 'logout') {
    store.commit('logout')
    next({ name: 'login' })
    return
  }

  if (
    to.matched.some(record => record.meta.requiresAuth)
    && !store.getters.isLoggedIn
  ) {
    next({name: 'login'})
    return
  }

  next()
})

export default router
