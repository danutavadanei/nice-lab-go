import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import VueAxios from 'vue-axios'
import store from './store'

const app = createApp(App)

app.use(store)
app.use(router)

app.use(VueAxios, axios)
app.provide('axios', app.config.globalProperties.axios)
app.provide('apiBaseUrl', import.meta.env.VITE_API_BASE_URL)
app.mount('#app')
