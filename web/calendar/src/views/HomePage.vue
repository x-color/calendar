<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Sample Calendar Application</h1>
        <p>This application is created with Vue.js and Go.</p>
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
    <v-row>
      <v-col cols="12" align="center">
        <v-alert type="info" border="top" colored-border elevation="2" max-width="600">
        <p class="text-center">
          サンプルユーザが用意されています。<br>
          以下を利用すればサインアップする必要はありません<br>
          Username: testuser<br>
          Password: P@ssw0rd
        </p>
        </v-alert>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" align="center">
        <v-alert type="info" border="top" colored-border elevation="2" max-width="600">
          <p class="text-left">
            このウェブアプリケーションは学習用に作成したものとなります。<br>
            そのため、重要なデータなどは入力しないようお願いします。<br>
            また、データは不定期に削除します。
          </p>
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'HomePage',
};
</script>
