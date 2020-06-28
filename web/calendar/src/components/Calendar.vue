<template>
  <v-container>
    <v-row>
      <v-col class="py-0">
        <v-sheet height="64">
          <v-toolbar flat color="white">
            <v-btn
              fab
              text
              small
              @click="setToday"
            >
              <v-icon medium>mdi-calendar</v-icon>
            </v-btn>
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
                <v-btn outlined color="grey darken-2" v-on="on" small>
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
        <v-sheet :height="calendarHeight">
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
            @click:day="openEditor"
            @click:time="openEditor"
            @change="updateRange"
            :class="{'small-cal': $vuetify.breakpoint.xs}"
          ></v-calendar>
          <Plan
            v-model="selectedOpen"
            :id="selectedPlan.id"
            :element="selectedElement"
          />
          <PlanEditor v-model="open" :plan="plan" @save="save" />
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import moment from 'moment';
import Plan from '@/components/Plan.vue';
import PlanEditor from '@/components/PlanEditor.vue';
import { mapGetters, mapActions, mapMutations } from 'vuex';

export default {
  components: {
    Plan,
    PlanEditor,
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
    open: false,
    plan: {
      id: '',
      calendar_id: '',
      user_id: '',
      name: '',
      memo: '',
      color: 'red',
      private: false,
      shares: [],
      start: null,
      end: null,
      allday: false,
    },
    calendarHeight: 600,
  }),
  computed: {
    ...mapGetters({
      getActivePlans: 'calendars/getActivePlans',
      getMyCalendars: 'calendars/getMyCalendars',
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
      const startYear = start.year;
      return `${startMonth} ${startYear}`;
    },
    monthFormatter() {
      return this.$refs.calendar.getFormatter({
        timeZone: 'UTC', month: 'long',
      });
    },
    plans() {
      return this.getActivePlans().map((p) => {
        const plan = { ...p };
        plan.start = this.formatDate(plan.start, !p.allday);
        plan.end = this.formatDate(plan.end, !p.allday);
        return plan;
      });
    },
  },
  async mounted() {
    this.calendarHeight = window.innerHeight - 150;
    this.$refs.calendar.checkChange();
    await this.getCalendars().catch((e) => {
      if (e.message === 'AuthError') {
        this.setUser({ signin: false });
        this.$router.push('/', () => {});
      }
    });
    setInterval(() => {
      this.getCalendars();
    }, 60 * 1000);
  },
  methods: {
    ...mapMutations({
      setFocusDate: 'calendars/setFocusDate',
      setUser: 'user/setUser',
    }),
    ...mapActions({
      getCalendars: 'calendars/getCalendars',
      addCalendar: 'calendars/addCalendar',
      addPlan: 'calendars/addPlan',
    }),
    openEditor(date) {
      this.plan = {
        id: '',
        calendar_id: '',
        user_id: this.$store.state.user.user.id,
        name: '',
        memo: '',
        color: 'red',
        private: false,
        shares: [],
        start: moment([date.year, date.month - 1, date.day, date.hour, 0]),
        end: moment([date.year, date.month - 1, date.day, date.hour + 1, 0]),
        allday: date.time === '',
      };
      this.open = true;
    },
    save(newPlan) {
      this.addPlan(newPlan);
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

<style>
.small-cal .v-btn {
  height: 40px;
  width: 40px;
}
</style>
