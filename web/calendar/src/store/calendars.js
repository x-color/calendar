import moment from 'moment';

// calendars: []calendar
// calendar: { id, user_id, name, color, []plan }
// plan: { id, name, memo, color, private, all_day, begin, end }
const state = () => ({
  focusDate: moment().format('YYYY-MM-DD'),
  calendars: [
    {
      id: 'calendar01',
      user_id: 'user01',
      active: true,
      name: 'my calendar',
      color: 'red',
      shares: [],
      plans: [
        {
          id: 'plan01',
          calendar_id: 'calendar01',
          owner_id: 'user01',
          name: 'my plan',
          memo: 'sample text',
          color: 'blue',
          private: false,
          shares: [],
          start: moment([2020, 5 - 1, 18, 9, 0, 0]).toISOString(),
          end: moment([2020, 5 - 1, 19, 18, 0, 0]).toISOString(),
          allday: false,
        },
      ],
    },
    {
      id: 'calendar02',
      user_id: 'user02',
      active: true,
      name: 'other\'s calendar',
      color: 'green',
      shares: ['user01'],
      plans: [
        {
          id: 'plan02',
          calendar_id: 'calendar02',
          owner_id: 'user02',
          name: 'other\'s plan',
          memo: 'sample text',
          color: 'purple',
          private: false,
          shares: ['calendar01'],
          start: moment([2020, 5 - 1, 25, 0, 0, 0]).toISOString(),
          end: moment([2020, 5 - 1, 25, 0, 0, 0]).toISOString(),
          allday: true,
        },
      ],
    },
  ],
});

const getters = {
  getCalendarByID: ({ calendars }) => (id) => calendars.find((calendar) => calendar.id === id),
  // eslint-disable-next-line max-len
  getPlanByID: ({ calendars }) => (id) => calendars.map((calendar) => {
    const p = calendar.plans.find((plan) => id === plan.id);
    if (!p) {
      return null;
    }
    const plan = { ...p };
    plan.start = moment(p.start);
    plan.end = moment(p.end);
    return plan;
  }).filter((plan) => plan)[0],
  // eslint-disable-next-line max-len
  getMyCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id === rootState.user.user.id),
  // eslint-disable-next-line max-len
  getSharedCalendars: ({ calendars }, _, rootState) => calendars.filter((calendar) => calendar.user_id !== rootState.user.user.id),
  getActiveCalendars: ({ calendars }) => calendars.filter((calendar) => calendar.active),
  getActivePlans: ({ calendars }) => () => calendars.map((calendar) => {
    if (!calendar.active) {
      return [];
    }
    return calendar.plans.map((p) => {
      const plan = { ...p };
      plan.start = moment(p.start);
      plan.end = moment(p.end);
      return plan;
    });
  }).flat(),
};

const actions = {
  // getCalendars({ commit }) {
  //   Call API
  //   commit('setCalendars', calendars)
  // },
  addCalendar({ commit }, calendar) {
    commit('addCalendar', calendar);
  },
  removeCalendar({ commit }, { id }) {
    commit('removeCalendar', id);
  },
  editCalendar({ commit }, calendar) {
    commit('setCalendar', calendar);
  },
  addPlan({ commit }, p) {
    const plan = { ...p };
    plan.start = p.start.toISOString();
    plan.end = p.end.toISOString();
    commit('addPlan', plan);
  },
  removePlan({ commit }, { id }) {
    commit('removePlan', id);
  },
  editPlan({ commit }, p) {
    const plan = { ...p };
    plan.start = p.start.toISOString();
    plan.end = p.end.toISOString();
    commit('setPlan', plan);
  },
};

const mutations = {
  setFocusDate(st, date) {
    // eslint-disable-next-line no-param-reassign
    st.focusDate = date;
  },
  setCalendars(st, calendars) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = calendars;
  },
  setCalendar(st, calendar) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (calendar.id === cal.id) {
        return calendar;
      }
      return cal;
    });
  },
  addCalendar(st, calendar) {
    st.calendars.push(calendar);
  },
  removeCalendar(st, id) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.filter((calendar) => calendar.id !== id);
  },
  setPlan(st, plan) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (cal.id === plan.calendar_id) {
        // eslint-disable-next-line no-param-reassign
        cal.plans = cal.plans.map((p) => {
          if (p.id === plan.id) {
            return plan;
          }
          return p;
        });
      }
      return cal;
    });
  },
  addPlan(st, plan) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      if (cal.id === plan.calendar_id) {
        cal.plans.push(plan);
      }
      return cal;
    });
  },
  removePlan(st, id) {
    // eslint-disable-next-line no-param-reassign
    st.calendars = st.calendars.map((cal) => {
      // eslint-disable-next-line no-param-reassign
      cal.plans = cal.plans.filter((p) => p.id !== id);
      return cal;
    });
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
