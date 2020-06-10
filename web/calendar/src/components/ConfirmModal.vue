<template>
  <v-row justify="center">
    <v-dialog v-model="open" max-width="290" @click:outside="open = false">
      <v-card>
        <v-card-title class="headline">{{ title }}</v-card-title>
        <v-card-text>{{ text }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="red darken-1" text @click="open = false">Cancel</v-btn>
          <v-btn color="green darken-1" text @click="confirm">OK</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
export default {
  name: 'ModalConfirm',
  props: {
    value: Boolean,
    title: String,
    text: String,
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
  methods: {
    confirm() {
      this.$emit('confirm');
      this.open = false;
    },
  },
};
</script>
