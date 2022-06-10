import { createStore } from 'vuex';
import VuexPersistence from 'vuex-persist'

const vuexLocal = new VuexPersistence({
  storage: window.localStorage
})

const store = createStore({
  state () {
    return {
      loggedIn: false,
      user: {},
      token: null,
    }
  },
  mutations: {
    setLoggedIn (state, isLoggedIn) {
      state.loggedIn = isLoggedIn
    },
    setUser (state, user) {
      state.user = user
    },
    setToken (state, token) {
      state.token = token
    }
  },
  getters: {
    isLoggedIn: (state) => state.loggedIn,
    user: (state) => state.user,
    token: (state) => state.token,
  },
  plugins: [vuexLocal.plugin]
})

export default store