<template>
  <v-content>
    <SideMenu v-model="open" />
    <v-container class="py-0">
      <v-row class="py-0">
        <Calendar />
      </v-row>
    </v-container>
  </v-content>
</template>

<script>
import { mapMutations } from 'vuex';
import Calendar from '@/components/Calendar.vue';
import SideMenu from '@/components/SideMenu.vue';

export default {
  name: 'CalendarPage',
  components: {
    Calendar,
    SideMenu,
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
