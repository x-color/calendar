<template>
  <v-navigation-drawer width="350" v-model="open" absolute temporary>
    <v-row justify="center">
      <v-col cols="auto">
        <MinCalendar />
      </v-col>
    </v-row>

    <v-divider></v-divider>

    <v-list-group no-action sub-group value="true">
      <template v-slot:activator>
        <v-list-item-content>
          <v-list-item-title>My Calendars</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-btn icon small @click.stop="isOpenedNewCal = true">
            <v-icon medium>mdi-plus</v-icon>
          </v-btn>
        </v-list-item-action>
      </template>
      <v-list dense>
        <v-list-item v-for="(cal, id) in getMyCalendars" :key="id" link>
          <v-list-item-content>
            <v-checkbox
              :label="cal.name"
              :color="cal.color"
              :input-value="cal.active"
              @change="activeCal(cal)"
              hide-details
              class="my-0"
            ></v-checkbox>
          </v-list-item-content>
          <v-list-item-action>
            <v-btn icon small @click.stop="openEditor(cal)">
              <v-icon small>mdi-pencil</v-icon>
            </v-btn>
          </v-list-item-action>
          <v-list-item-action>
            <v-btn icon small @click.stop="openRemoveModal(cal)">
              <v-icon small>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
        <v-list-item v-if="isOpenedNewCal">
          <v-list-item-content class="py-0 pl-8">
            <v-text-field
              v-model="newCalName"
              dense
              clearable
              hide-details
              autofocus
              placeholder="New calendar"
              @blur="isOpenedNewCal = false"
              @keydown.enter="addNewCalendar()"
            ></v-text-field>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-list-group>

    <v-list-group no-action sub-group value="true">
      <template v-slot:activator>
        <v-list-item-content>
          <v-list-item-title>Other Calendars</v-list-item-title>
        </v-list-item-content>
      </template>
      <v-list dense>
        <v-list-item v-for="(cal, id) in getSharedCalendars" :key="id" link>
          <v-list-item-content>
            <v-checkbox
              :label="cal.name"
              :color="cal.color"
              :input-value="cal.active"
              @change="activeCal(cal)"
              hide-details
              class="my-0"
            ></v-checkbox>
          </v-list-item-content>
          <v-list-item-action>
            <v-btn icon small @click.stop="openRemoveModal(cal)">
              <v-icon small>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-list-group>
    <CalendarEditor v-model="isOpenedEditor" :calendar="selectedCal" @save="save" />
    <ConfirmModal
      v-model="isOpenedRemoveModal"
      :title="removeModalTitle"
      @confirm="removeCal(selectedCal)"
    />
  </v-navigation-drawer>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import generateUuid from '@/utils/uuid';
import MinCalendar from '@/components/MinCalendar.vue';
import CalendarEditor from '@/components/CalendarEditor.vue';
import ConfirmModal from '@/components/ConfirmModal.vue';

export default {
  components: {
    MinCalendar,
    CalendarEditor,
    ConfirmModal,
  },
  props: {
    value: Boolean,
  },
  computed: {
    ...mapGetters({
      getCalendarByID: 'calendars/getCalendarByID',
      getMyCalendars: 'calendars/getMyCalendars',
      getSharedCalendars: 'calendars/getSharedCalendars',
    }),
    open: {
      get() {
        return this.value;
      },
      set(v) {
        this.$emit('input', v);
      },
    },
    removeModalTitle() {
      return `Delete '${this.selectedCal.name}'?`;
    },
  },
  data() {
    return {
      isOpenedEditor: false,
      isOpenedRemoveModal: false,
      selectedCal: {},
      isOpenedNewCal: false,
      newCalName: '',
    };
  },
  methods: {
    ...mapActions({
      addCalendar: 'calendars/addCalendar',
      editCalendar: 'calendars/editCalendar',
      removeCalendar: 'calendars/removeCalendar',
    }),
    activeCal(cal) {
      // eslint-disable-next-line no-param-reassign
      cal.active = !cal.active;
      this.editCalendar(cal);
    },
    removeCal(cal) {
      this.removeCalendar(cal);
    },
    openEditor(cal) {
      this.selectedCal = cal;
      this.isOpenedEditor = true;
    },
    openRemoveModal(cal) {
      this.selectedCal = cal;
      this.isOpenedRemoveModal = true;
    },
    save(newCal) {
      // Todo: Cancel to change calendar and set old calendar if calling API is failed.
      this.editCalendar(newCal);
    },
    addNewCalendar() {
      if (!this.newCalName) {
        return;
      }
      const newCal = {
        id: generateUuid(),
        user_id: this.$store.state.user.user.id,
        active: true,
        name: this.newCalName,
        color: 'red',
        shares: [],
        plans: [],
      };
      this.addCalendar(newCal);
      this.isOpenedNewCal = false;
    },
  },
};
</script>
