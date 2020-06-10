<template>
  <v-menu v-model="open" offset-y :close-on-content-click="false">
    <template v-slot:activator="{ on }">
      <v-icon :color="color" class="ma-1" medium v-on="on">mdi-circle</v-icon>
    </template>
    <v-card>
      <v-list dense subheader max-width="400">
        <v-subheader>Color</v-subheader>
        <v-container>
          <v-row class="px-2">
            <v-col
              v-for="(c, i) in colors"
              :key="i"
              cols="auto"
              class="pa-0"
            >
              <v-list-item class="px-0">
                <v-icon medium :color="c" @click.stop="select(c)">mdi-circle</v-icon>
              </v-list-item>
            </v-col>
          </v-row>
        </v-container>
      </v-list>
    </v-card>
  </v-menu>
</template>

<script>
export default {
  name: 'ColorPicker',
  props: {
    value: Boolean,
    color: String,
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
  data() {
    return {
      colors: ['red', 'green', 'blue', 'purple'],
    };
  },
  methods: {
    select(color) {
      this.$emit('select', color);
      this.open = false;
    },
  },
};
</script>
