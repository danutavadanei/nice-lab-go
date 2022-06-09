import { createStore } from 'vuex';

const store = createStore({
  state () {
    return {
      loggedIn: true,
      user: {
        name: 'dan',
        // type: 'professor',
        type: 'student',
      },
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