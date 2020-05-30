<template>
  <div class="home">
    <Header v-model="open" />
    <v-content>
      <v-container>
        <v-row>
          <SideMenu v-model="open" />
          <Calendar />
        </v-row>
      </v-container>
    </v-content>
  </div>
</template>

<script>
import Header from '@/components/Header.vue';
import Calendar from '@/components/Calendar.vue';
import SideMenu from '@/components/SideMenu.vue';

export default {
  name: 'Home',
  components: {
    Header,
    Calendar,
    SideMenu,
  },
  data() {
    return {
      open: false,
    };
  },
};
</script>
