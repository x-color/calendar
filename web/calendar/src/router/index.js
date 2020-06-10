import Vue from 'vue';
import VueRouter from 'vue-router';
import Store from '@/store/index';
import HomePage from '@/views/HomePage.vue';
import CalendarPage from '@/views/CalendarPage.vue';
import SigninPage from '@/views/SigninPage.vue';
import SignupPage from '@/views/SignupPage.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'HomePage',
    component: HomePage,
  },
  {
    path: '/calendar',
    name: 'CalendarPage',
    component: CalendarPage,
    meta: { requiresAuth: true },
  },
  {
    path: '/signup',
    name: 'SignupPage',
    component: SignupPage,
  },
  {
    path: '/signin',
    name: 'SigninPage',
    component: SigninPage,
  },
  {
    path: '*',
    redirect: '/',
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.query.redirect) {
    next({ path: to.query.redirect });
  }
  if (to.matched.some((record) => record.meta.requiresAuth) && !Store.state.user.user.signin) {
    next({ path: '/' });
  } else if ((to.path === '/signin' || to.path === '/signup') && Store.state.user.user.signin) {
    next({ path: '/calendar' });
  } else {
    next();
  }
});


export default router;
