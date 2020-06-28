<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Sample Calendar Application</h1>
        <p>This application is created with Vue.js and Go.</p>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" align="center">
        <h2>Signin Form</h2>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-if="isSigninFailed" cols="12" align="center">
        <p class="red--text">
          Signin Failed...
          <br />Please try again
        </p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col cols="auto" align="center">
        <v-text-field
          v-model="username"
          :rules="[rules.required]"
          type="text"
          name="input-10-2"
          label="User name"
        ></v-text-field>
        <v-text-field
          v-model="password"
          :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
          :rules="[rules.required]"
          :type="show ? 'text' : 'password'"
          name="input-10-2"
          label="Password"
          @click:append="show = !show"
        ></v-text-field>
      </v-col>
      <v-col cols="12" align="center">
        <v-btn
          x-large
          color="primary"
          :disabled="!username || !password"
          @click="SigninAndGoToPage"
        >Signin</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'Signin',
  data() {
    return {
      show: false,
      isSigninFailed: false,
      username: '',
      password: '',
      rules: {
        required: (value) => !!value || 'Required.',
      },
    };
  },
  methods: {
    ...mapActions({
      signin: 'user/signin',
    }),
    SigninAndGoToPage() {
      this.isSigninFailed = false;
      this.signin({
        username: this.username,
        password: this.password,
        callback: (loggedIn) => {
          if (loggedIn) {
            this.$router.push('/calendar', () => {});
          } else {
            this.isSigninFailed = true;
          }
        },
      });
    },
  },
};
</script>
