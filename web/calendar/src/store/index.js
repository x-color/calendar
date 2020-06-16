import Vue from 'vue';
import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';
import calendars from '@/store/calendars';
import user from '@/store/user';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    calendars,
    user,
  },
  plugins: [
    createPersistedState({
      key: 'sample-calendar-application',
      paths: ['user.user'],
      storage: window.localStorage,
    }),
  ],
});
