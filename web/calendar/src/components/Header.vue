<template>
  <div>
    <v-app-bar color="primary" dense dark>
      <v-app-bar-nav-icon v-if="appMode" @click="$emit('input', !value)"></v-app-bar-nav-icon>
      <v-toolbar-title v-else>Sample Calendar Application</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-menu left bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>

        <v-list>
          <v-list-item
            v-for="(menu, i) in menuList"
            :key="i"
            :disabled="!menu.active"
            @click="menu.action"
          >
            <v-list-item-title>{{ menu.title }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'Header',
  props: {
    value: Boolean,
    appMode: Boolean,
  },
  data() {
    return {
      menuList: [
        {
          title: 'top',
          action: () => { this.$router.push('/', () => {}); },
          active: true,
        },
        {
          title: 'home',
          action: () => { this.$router.push('/calendar', () => {}); },
          active: this.$store.state.user.user.signin,
        },
        {
          title: 'signup',
          action: () => { this.$router.push('/signup', () => {}); },
          active: !this.$store.state.user.user.signin,
        },
        {
          title: 'signin',
          action: () => { this.$router.push('/signin', () => {}); },
          active: !this.$store.state.user.user.signin,
        },
        {
          title: 'signout',
          action: () => {
            this.signout();
            this.$router.push('/', () => {});
          },
          active: this.$store.state.user.user.signin,
        },
      ],
    };
  },
  methods: {
    ...mapActions({
      signout: 'user/signout',
    }),
  },
};
</script>

<style>
</style>
