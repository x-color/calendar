<template>
  <v-navigation-drawer width="350" v-model="drawer" absolute temporary>
    <v-list-item>
      <v-list-item-avatar>
        <!-- <v-img src="https://randomuser.me/api/portraits/men/78.jpg"></v-img> -->
      </v-list-item-avatar>

      <v-list-item-content>
        <v-list-item-title>John Leider</v-list-item-title>
      </v-list-item-content>
    </v-list-item>

    <v-divider></v-divider>

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
      </template>
      <v-list dense>
        <v-list-item v-for="(cal, id) in this.getMyCalendars" :key="id" link>
          <v-checkbox
            :label="cal.name"
            :color="cal.color"
            :input-value="cal.active"
            @change="activeCal(cal)"
            hide-details
            class="my-0"
          ></v-checkbox>
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
        <v-list-item v-for="(cal, id) in this.getSharedCalendars" :key="id" link>
          <v-checkbox
            :label="cal.name"
            :color="cal.color"
            :input-value="cal.active"
            @change="activeCal(cal)"
            hide-details
            class="my-0"
          ></v-checkbox>
        </v-list-item>
      </v-list>
    </v-list-group>
  </v-navigation-drawer>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import MinCalendar from '@/components/MinCalendar.vue';

export default {
  components: {
    MinCalendar,
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
    drawer: {
      get() {
        return this.value;
      },
      set(v) {
        this.$emit('input', v);
      },
    },
  },
  methods: {
    ...mapActions({
      editCalendar: 'calendars/editCalendar',
    }),
    activeCal(cal) {
      // eslint-disable-next-line no-param-reassign
      cal.active = !cal.active;
      this.editCalendar(cal);
    },
  },
};
</script>
