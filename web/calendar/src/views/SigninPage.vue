<template>
  <div>
    <Signin />
    <v-row>
      <v-col cols="12" align="center">
        <v-alert type="info" border="top" colored-border elevation="2" max-width="600">
          <p class="text-center">
            サンプルユーザ<br>
            Username: testuser<br>
            Password: P@ssw0rd
          </p>
        </v-alert>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import Signin from '@/components/Signin.vue';

export default {
  name: 'SigninPage',
  components: {
    Signin,
  },
};
</script>
