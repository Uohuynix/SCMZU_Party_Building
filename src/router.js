import { createRouter, createWebHistory } from 'vue-router';
import Login from './views/Login.vue';
import Register from './views/Register.vue';
import Dashboard from './views/Dashboard.vue';
import Development from './views/Development.vue';
import Training from './views/Training.vue';
import Organization from './views/Organization.vue';

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/dashboard', component: Dashboard },
  { path: '/development', component: Development },
  { path: '/training', component: Training },
  { path: '/organization', component: Organization },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;