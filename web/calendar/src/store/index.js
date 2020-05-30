import Vue from 'vue';
import Vuex from 'vuex';
import calendars from '@/store/calendars';
import user from '@/store/user';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    calendars,
    user,
  },
});
