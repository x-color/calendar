<template>
  <v-menu v-model="open" :close-on-content-click="false" :activator="element" offset-x>
    <v-card color="grey lighten-4" min-width="350px" flat>
      <v-toolbar :color="plan.color" dark>
        <v-toolbar-title v-html="plan.name"></v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn
          icon
          @click.stop="openEditor = true"
          :disabled="plan.owner_id !== $store.state.user.user.id"
        >
          <v-icon>mdi-pencil</v-icon>
        </v-btn>
        <v-btn icon @click.stop="openRemoveModal = true">
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
    <PlanEditor v-model="openEditor" :plan="plan" @save="save" />
    <ConfirmModal v-model="openRemoveModal" :title="removeModalTitle" @confirm="remove" />
  </v-menu>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import ConfirmModal from '@/components/ConfirmModal.vue';
import PlanEditor from '@/components/PlanEditor.vue';

export default {
  name: 'Plan',
  components: {
    PlanEditor,
    ConfirmModal,
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
    plan() {
      const p = this.getPlanByID(this.id);
      if (p) {
        return p;
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
    removeModalTitle() {
      return `Delete ${this.plan.name}?`;
    },
  },
  data() {
    return {
      openEditor: false,
      openRemoveModal: false,
    };
  },
  methods: {
    ...mapActions({
      editPlan: 'calendars/editPlan',
      removePlan: 'calendars/removePlan',
    }),
    save(newPlan) {
      this.editPlan(newPlan);
    },
    remove() {
      this.open = false;
      this.removePlan(this.plan);
    },
  },
};
</script>

<style>
</style>
