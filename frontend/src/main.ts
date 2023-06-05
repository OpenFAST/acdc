import "bootstrap/dist/css/bootstrap.min.css"

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useProjectStore } from './project';

const pinia = createPinia()
const app = createApp(App)

app.use(router)
app.use(pinia)
app.mount('#app')

// Redirect to home path if project not loaded
const project = useProjectStore()
router.beforeEach((to, from) => {
    if (to.path != '/' && !project.loaded) return { path: '/' }
})