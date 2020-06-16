<template>
  <div>
    <v-app-bar color="primary" dense dark>
      <v-app-bar-nav-icon v-if="appMode" @click="$emit('input', !value)"></v-app-bar-nav-icon>
      <v-toolbar-title v-else>
      <v-btn text @click="$router.push('/', () => {})">
        Sample Calendar Application
      </v-btn>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-menu v-if="$store.state.user.user.signin" left bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>

        <v-list>
          <v-list-item v-for="(menu, i) in menuList" :key="i" @click="menu.action">
            <v-list-item-icon>
              <v-icon v-text="menu.icon"></v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ menu.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
    <UserInfo v-model="openInfo" />
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import UserInfo from '@/components/UserInfo.vue';

export default {
  name: 'Header',
  components: {
    UserInfo,
  },
  props: {
    value: Boolean,
    appMode: Boolean,
  },
  data() {
    return {
      menuList: [
        {
          title: 'account',
          icon: 'mdi-account',
          action: () => {
            this.openInfo = true;
          },
        },
        {
          title: 'sync',
          icon: 'mdi-sync',
          action: () => {
            this.getCalendars();
          },
        },
        {
          title: 'top page',
          icon: 'mdi-home',
          action: () => {
            this.$router.push('/', () => {});
          },
        },
        {
          title: 'signout',
          icon: 'mdi-exit-to-app',
          action: () => {
            this.signout();
            this.$router.push('/', () => {});
          },
        },
      ],
      openInfo: false,
    };
  },
  methods: {
    ...mapActions({
      signout: 'user/signout',
      getCalendars: 'calendars/getCalendars',
    }),
  },
};
</script>

<style>
</style>
