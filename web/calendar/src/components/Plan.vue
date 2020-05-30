<template>
  <v-menu
    v-model="open"
    :close-on-content-click="false"
    :activator="element"
    offset-x
  >
    <v-card color="grey lighten-4" min-width="350px" flat>
      <v-toolbar :color="plan.color" dark>
        <v-toolbar-title v-html="plan.name"></v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn icon @click="openEditor = true">
          <v-icon>mdi-pencil</v-icon>
        </v-btn>
        <v-btn icon>
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </v-toolbar>
      <v-card-text>
        <span v-html="plan.memo"></span>
        <span v-html="plan.start"></span>
        <span v-html="plan.end"></span>
      </v-card-text>
      <v-card-actions>
        <v-btn text color="secondary" @click="open = false">Cancel</v-btn>
      </v-card-actions>
    </v-card>
    <PlanEditor v-model="openEditor" :start="s" :plan="plan"/>
  </v-menu>
</template>

<script>
import moment from 'moment';
import { mapGetters } from 'vuex';
import PlanEditor from '@/components/PlanEditor.vue';

export default {
  name: 'Plan',
  components: {
    PlanEditor,
  },
  props: {
    value: Boolean,
    id: String,
    element: HTMLDivElement,
  },
  computed: {
    ...mapGetters({
      getPlanByID: 'calendars/getPlanByID',
    }),
    s() {
      const a = moment([2020, 5 - 1, 24, 15, 0, 0]);
      console.log(a.format());
      return a;
      // return moment([2020, 5 - 1, 24, 15, 0, 0]);
    },
    plan() {
      if (this.id) {
        return this.getPlanByID(this.id);
      }
      return {};
    },
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
      openEditor: false,
    };
  },
};
</script>

<style>
</style>
