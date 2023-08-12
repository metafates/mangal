import './assets/main.css'
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import mitt from 'mitt'

import { VueMasonryPlugin } from 'vue-masonry'

import App from './App.vue'
import router from './router'

const emitter = mitt()
const app = createApp(App)
app.config.globalProperties.emitter = emitter

app.use(createPinia())
app.use(VueMasonryPlugin)
app.use(router)

app.mount('#app')
