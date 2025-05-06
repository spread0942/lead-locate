import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Maps from '../views/Maps.vue';
import Gowler from '../views/Gowler.vue';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/maps',
        name: 'Maps',
        component: Maps,
    },
    {
        path: '/gowler/:website?',
        name: 'Gowler',
        component: Gowler
    }
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
});

export default router;