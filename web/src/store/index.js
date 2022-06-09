import { createStore } from 'vuex';

const store = createStore({
  state () {
    return {
      loggedIn: false,
      user: null,
    }
  },
  mutations: {
    setLoggedIn (state, isLoggedIn) {
      state.loggedIn = isLoggedIn
    }
  },
  getters: {
    isLoggedIn: (state) => state.loggedIn,
    user: (state) => state.user,
  }
})

export default store