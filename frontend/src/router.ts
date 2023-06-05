import { createRouter, createWebHashHistory } from 'vue-router'

export default createRouter({
    history: createWebHashHistory(),
    linkActiveClass: "active",
    routes: [
        {
            path: '/',
            component: () => import('./components/Project.vue'),
        },
        {
            path: '/turbine',
            component: () => import('./components/Turbine.vue'),
        },
        {
            path: '/analyze',
            component: () => import('./components/Analyze.vue'),
        },
        {
            path: '/results',
            component: () => import('./components/Results.vue'),
        },
    ],
})
