<template>
  <v-row justify="center">
    <v-dialog v-model="open" max-width="350" @click:outside="open = false">
      <v-card>
        <v-card-title class="headline">{{ $store.state.user.user.name }}</v-card-title>
        <v-card-text>
            User ID: <br>
            {{ $store.state.user.user.id }}
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="green darken-1" text @click="open = false">OK</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
export default {
  name: 'UserInfo',
  props: {
    value: Boolean,
  },
  computed: {
    open: {
      get() {
        return this.value;
      },
      set(v) {
        this.$emit('input', v);
      },
    },
  },
};
</script>
