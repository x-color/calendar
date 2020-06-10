<template>
  <v-row justify="space-between">
    <div>
      <v-date-picker
        v-model="focus"
        no-title
        :events="plans"
        event-color="green lighten-1"
      ></v-date-picker>
    </div>
  </v-row>
</template>

<script>
import moment from 'moment';
import { mapGetters, mapMutations } from 'vuex';

export default {
  data: () => ({
    today: moment().format('YYYY-MM-DD'),
  }),
  computed: {
    ...mapGetters({
      getActivePlans: 'calendars/getActivePlans',
    }),
    focus: {
      get() {
        return this.$store.state.calendars.focusDate;
      },
      set(v) {
        this.setFocusDate(v);
      },
    },
    plans() {
      const plans = this.getActivePlans().map((plan) => {
        const start = plan.start.hour(0).minute(0).second(0).millisecond(0);
        const end = plan.end.hour(23).minute(0).second(0).millisecond(0);

        const list = [];
        let date = start.clone().hour(1);
        while (date.isBetween(start, end)) {
          list.push(date);
          date = date.clone().add(1, 'days');
        }

        return list;
      }).flat();

      return plans.map((e) => e.format('YYYY-MM-DD'));
    },
  },
  methods: {
    ...mapMutations({
      setFocusDate: 'calendars/setFocusDate',
    }),
  },
};
</script>
