<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Sample Calendar Application</h1>
        <p>This application is created with Vue.js.</p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col cols="12" align="center">
        <v-btn large color="primary" @click="$router.push('/signup', () => {})">
          GO TO SIGN UP
        </v-btn>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col cols="12" align="center">
        Do you want to
        <router-link to="/signin">signin</router-link>?
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'HomePage',
};
</script>
