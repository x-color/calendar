<template>
  <v-app>
    <Header v-model="open" :appMode="$route.name === 'CalendarPage'" />
    <router-view />
  </v-app>
</template>

<script>
import Header from '@/components/Header.vue';
import { mapMutations } from 'vuex';

export default {
  name: 'App',
  components: {
    Header,
  },
  computed: {
    open: {
      get() {
        return this.$store.state.user.sideMenu;
      },
      set(v) {
        this.setSideMenu(v);
      },
    },
  },
  methods: {
    ...mapMutations({
      setSideMenu: 'user/setSideMenu',
    }),
  },
};
</script>
