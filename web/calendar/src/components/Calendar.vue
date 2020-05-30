<template>
  <v-container>
    <v-row class="fill-height">
      <v-col>
        <v-sheet height="64">
          <v-toolbar flat color="white">
            <v-btn
              outlined
              class="mr-4"
              color="grey darken-2"
              @click="setToday"
            >Today</v-btn>
            <v-btn fab text small color="grey darken-2" @click="prev">
              <v-icon small>mdi-chevron-left</v-icon>
            </v-btn>
            <v-btn fab text small color="grey darken-2" @click="next">
              <v-icon small>mdi-chevron-right</v-icon>
            </v-btn>
            <v-toolbar-title>{{ title }}</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-menu bottom right>
              <template v-slot:activator="{ on }">
                <v-btn outlined color="grey darken-2" v-on="on">
                  <span>{{ typeToLabel[type] }}</span>
                  <v-icon right>mdi-menu-down</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item @click="type = 'day'">
                  <v-list-item-title>Day</v-list-item-title>
                </v-list-item>
                <v-list-item @click="type = 'week'">
                  <v-list-item-title>Week</v-list-item-title>
                </v-list-item>
                <v-list-item @click="type = 'month'">
                  <v-list-item-title>Month</v-list-item-title>
                </v-list-item>
                <v-list-item @click="type = '4day'">
                  <v-list-item-title>4 days</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-toolbar>
        </v-sheet>
        <v-sheet height="600">
          <v-calendar
            ref="calendar"
            v-model="focus"
            color="primary"
            :events="plans"
            :event-color="getPlanColor"
            :now="today"
            :type="type"
            @click:event="showPlan"
            @click:more="viewDay"
            @click:date="viewDay"
            @change="updateRange"
            @mousedown:day="tmpFunc"
          ></v-calendar>
          <Plan
            v-model="selectedOpen"
            :id="selectedPlan.id"
            :element="selectedElement"
          />
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import moment from 'moment';
import Plan from '@/components/Plan.vue';
import { mapGetters, mapMutations } from 'vuex';

export default {
  components: {
    Plan,
  },
  data: () => ({
    today: moment().format('YYYY-MM-DD'),
    type: 'month',
    typeToLabel: {
      month: 'Month',
      week: 'Week',
      day: 'Day',
      '4day': '4 Days',
    },
    start: null,
    end: null,
    selectedPlan: {},
    selectedElement: null,
    selectedOpen: false,
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
    title() {
      const { start, end } = this;
      if (!start || !end) {
        return '';
      }

      const startMonth = this.monthFormatter(start);
      const endMonth = this.monthFormatter(end);
      const suffixMonth = startMonth === endMonth ? '' : endMonth;

      const startYear = start.year;
      const endYear = end.year;
      const suffixYear = startYear === endYear ? '' : endYear;

      const startDay = start.day + this.nth(start.day);
      const endDay = end.day + this.nth(end.day);

      switch (this.type) {
        case 'month':
          return `${startMonth} ${startYear}`;
        case 'week':
        case '4day':
          return `${startMonth} ${startDay} ${startYear} - ${suffixMonth} ${endDay} ${suffixYear}`;
        case 'day':
          return `${startMonth} ${startDay} ${startYear}`;
        default:
          return '';
      }
    },
    monthFormatter() {
      return this.$refs.calendar.getFormatter({
        timeZone: 'UTC', month: 'long',
      });
    },
    plans() {
      return this.getActivePlans.map((p) => {
        const plan = { ...p };
        plan.start = this.formatDate(plan.start, !p.allday);
        plan.end = this.formatDate(plan.end, !p.allday);
        return plan;
      });
    },
  },
  mounted() {
    this.$refs.calendar.checkChange();
  },
  methods: {
    ...mapMutations({
      setFocusDate: 'calendars/setFocusDate',
    }),
    tmpFunc(v) {
      console.log(v);
    },
    viewDay({ date }) {
      this.focus = date;
      this.type = 'day';
    },
    getPlanColor(plan) {
      return plan.color;
    },
    setToday() {
      this.focus = this.today;
    },
    prev() {
      this.$refs.calendar.prev();
    },
    next() {
      this.$refs.calendar.next();
    },
    showPlan({ nativeEvent, event }) {
      const open = () => {
        this.selectedPlan = event;
        this.selectedElement = nativeEvent.target;
        setTimeout(() => { this.selectedOpen = true; }, 10);
      };

      if (this.selectedOpen) {
        this.selectedOpen = false;
        setTimeout(open, 10);
      } else {
        open();
      }

      nativeEvent.stopPropagation();
    },
    updateRange({ start, end }) {
      this.start = start;
      this.end = end;
    },
    nth(d) {
      return d > 3 && d < 21
        ? 'th'
        : ['th', 'st', 'nd', 'rd', 'th', 'th', 'th', 'th', 'th', 'th'][d % 10];
    },
    rnd(a, b) {
      return Math.floor((b - a + 1) * Math.random()) + a;
    },
    formatDate(t, withTime) {
      return withTime
        ? t.format('YYYY-MM-DD HH:mm')
        : t.format('YYYY-MM-DD');
    },
  },
};
</script>
