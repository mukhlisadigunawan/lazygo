import {createRouter, createWebHistory} from 'vue-router'

const routes = [

    {
        path: '/',
        component: () => import('../pages/PageWithSidebar.vue'),
        children: [
            {
                path: '/{{LowerCase .EntityName}}',
                component: () => import('../pages/{{LowerCase .EntityName}}/ViewTable.vue'),
            },
        ],
    },

]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router